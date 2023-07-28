package controller

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"rnd/goerliscan/scanner/config"
	"rnd/goerliscan/scanner/logger"
	"rnd/goerliscan/scanner/model"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func StartScanner(cf *config.Config, blockModel *model.Model, transactionModel *model.Model) {
	// Initialize ethclient
	client, err := ethclient.Dial(cf.Netowrk.URL)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	// Subscribe to new block headers
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	for {
		select {
		case err := <-sub.Err():
			logger.Error(err)
			panic(err)
		case header := <-headers:

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				logger.Error(err)
				panic(err)
			}

			//Get the latest block number from the database
			latestBlockNumberDB, err := blockModel.GetLatestBlockNumber()
			if err != nil {
				logger.Error(err)
				panic(err)
			}
			
			// Get the latest block number from the client
			latestBlockNumber := block.Number().Uint64()
		
			if latestBlockNumber > latestBlockNumberDB {
			
				// Save the missing blocks
			for i := latestBlockNumberDB + 1; i <= latestBlockNumber; i++ {
				block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
				if err != nil {
					logger.Error(err)
					continue
				}
				// Setting Block Data
				b := getBlockData(header, block)
				
				txs := block.Transactions()
				if len(txs) > 0 {
					for _, tx := range txs {
						// Setting Transaction Data
						t := getTxsData(client, header, tx, block)
						
						b.Transactions = append(b.Transactions, t.Hash)
						
						// Save the transaction to the database
						err = transactionModel.SaveTransaction(&t)
						if err != nil {
							err := transactionModel.RemoveTxs(block.Number().Uint64())
							logger.Error(err)
							if(err != nil){
								logger.Error(err)
							}
							panic(err)
						}
					}
				}

			// Save the block to the database
			err = blockModel.SaveBlock(&b)
			if err != nil {
						logger.Error(err)
						if(err != nil){
							logger.Error(err)
						}
				panic(err)
			}
			
			// latestBlockNumber update
			latestBlockNumberDB = latestBlockNumber
			}
		}
		}
	}
}

// get Block Data to th ethclient
func getBlockData(header *types.Header, block *types.Block) (model.Block) {
	timestamp := int64(block.Time())
	timeUTC := time.Unix(timestamp, 0).UTC()
	utcTimeFormatted := timeUTC.Format("2006-01-02 15:04:05 AM MST")
	gasUsed := block.GasUsed()
	baseFeePerGas := block.BaseFee()
	burntFees := new(big.Int).Mul(baseFeePerGas, new(big.Int).SetUint64(gasUsed))
	nonce := binary.BigEndian.Uint64(header.Nonce[:])
	hexNonce := fmt.Sprintf("0x%016x", nonce)
	extraData := hex.EncodeToString(block.Extra())
				// Create block structure
	b := model.Block{
		BlockNumber: block.Number().Uint64(),  	
		Time: utcTimeFormatted,         	
		FeeRecipient: block.Coinbase().Hex(),  
		BlockSize: block.Size(),			
		GasUsed: gasUsed,     	
		GasLimit: block.GasLimit(),     	
		BaseFeePerGas: baseFeePerGas,
		BurntFees: burntFees,
		ExtraData: extraData,		
		BlockHash: block.Hash().Hex(),
		ParentHash:	block.ParentHash().Hex(),	
		StateRoot: block.Root().Hex(), 
		Nonce: hexNonce,      	
		Transactions: make([]string, 0),
	}
	return b
}
// get Txs Data to th ethclient
func getTxsData(client *ethclient.Client, header *types.Header, tx *types.Transaction, block *types.Block) model.Transaction {
	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	var status string
	if receipt.Status == 1 {
		status = "Success"
	} else if receipt.Status == 0 {
		status = "Fail"
	} else {
		status = "Unknown"
	}

	value := tx.Value()
	gasUsed := receipt.GasUsed
	gasPrice := tx.GasPrice()
	transactionFee := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice)
	timestamp := int64(block.Time())
	timeUTC := time.Unix(timestamp, 0).UTC()
	utcTimeFormatted := timeUTC.Format("2006-01-02 15:04:05 AM MST")
	
	// Create transaction structure
	t:= model.Transaction{
		Hash: tx.Hash().Hex(),        		
		Status: status,      		
		Time: utcTimeFormatted,        		
		From: sender.Hex(),        	
		To: "",          		
		Value: value,      			
		TransactionFee: transactionFee, 	
		GasPrice: tx.GasPrice(),    		
		GasUsed: gasUsed,					
		GasLimit: tx.Gas(),    		
		BlockHash: block.Hash().Hex(),
		BlockNumber: block.Number().Uint64(), 		   		    
	}

	if tx.To() != nil {
		t.To = tx.To().Hex()
	}

	return t
}