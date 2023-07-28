package controller

import (
	"context"
	"math/big"
	"rnd/goerliscan/scanner/config"
	"rnd/goerliscan/scanner/logger"
	"rnd/goerliscan/scanner/model"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func StartScanner(cf *config.Config, blockModel *model.Model, transactionModel *model.Model) error{
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
				b := model.GetBlockData(header, block)
				
				txs := block.Transactions()
				if len(txs) > 0 {
					for _, tx := range txs {
						// Setting Transaction Data
						t := model.GetTxsData(client, header, tx, block)
						
						b.Transactions = append(b.Transactions, t.Hash)
						
						// Save the transaction to the database
						err = transactionModel.SaveTransaction(&t)
						if err != nil {
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


