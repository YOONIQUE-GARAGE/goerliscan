package model

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"rnd/goerliscan/scanner/config"
	"rnd/goerliscan/scanner/logger"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	colHeader 			*mongo.Collection
	colBlock      *mongo.Collection
	colTransaction *mongo.Collection
	db            *mongo.Database
}

type Header struct {
	BlockNumber  	uint64      `bson:"blockNumber"`
	ParentHash		string      `bson:"parentHash"`
	Bloom					[]byte 			`bson:"bloom"`
	Time         	string      `bson:"timestamp"`
	Nonce        	string      `bson:"nonce"`
}

type Block struct {
	BlockNumber  	uint64      `bson:"blockNumber"`
	FeeRecipient  string	    `bson:"feeRecipient"`
	BlockSize			uint64      `bson:"blockSize"`
	GasUsed      	uint64      `bson:"gasUsed"`
	GasLimit     	uint64      `bson:"gasLimit"`
	BaseFeePerGas *big.Int    `bson:"baseFeePerGas"`
	BurntFees			*big.Int		`bson:"burntFees"`
	ExtraData			string		  `bson:"extraData"`
	BlockHash    	string      `bson:"blockHash"`
	StateRoot     string			`bson:"stateRoot"`
	Transactions 	[]string 		`bson:"transactions"`
}

type Transaction struct {
	Hash        		string  	 `bson:"hash"`
	Status      		string  	 `bson:"status"`
	Time        		string  	 `bson:"timestamp"`
	From        		string  	 `bson:"from"`
	To          		string  	 `bson:"to"` // return nil for contract
	Value      			*big.Int	 `bson:"amount"`
	TransactionFee 	*big.Int 	 `bson:"transactionFee"`
	GasPrice    		*big.Int   `bson:"gasPrice"`
	GasUsed					uint64	 	 `bson:"gasUsed"`
	GasLimit    		uint64  	 `bson:"gasLimit"`
	BlockHash   		string  	 `bson:"blockHash"`
	BlockNumber 		uint64  	 `bson:"blockNumber"`
}

// Setting Model
func NewModel(config *config.Config) (*Model, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Database.Host))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	db := client.Database(config.Database.Name)
	colHeader := db.Collection("header")
	colBlock := db.Collection("block")
	colTransaction := db.Collection("transaction")
	model := &Model{
		colHeader: colHeader,
		colBlock:      colBlock,
		colTransaction: colTransaction,
		db:            db,
	}
	return model, nil
}

// Get Header Data to th ethclient
func GetHeaderData(header *types.Header, block *types.Block, c chan Header){
	logger.Debug("GetHeaderData: start")
	timestamp := int64(block.Time())
	timeUTC := time.Unix(timestamp, 0).UTC()
	utcTimeFormatted := timeUTC.Format("2006-01-02 15:04:05 AM MST")
	nonce := binary.BigEndian.Uint64(header.Nonce[:])
	hexNonce := fmt.Sprintf("0x%016x", nonce)
	h := Header{
		BlockNumber: block.Number().Uint64(),  	
		ParentHash:	block.ParentHash().Hex(),	 
		Bloom: header.Bloom[:],
		Time: utcTimeFormatted, 
		Nonce: hexNonce,  
	}     
	c <- h
}

// Get Block Data to th ethclient
func GetBlockData(header *types.Header, block *types.Block, c chan Block) {
	logger.Debug("GetBlockData: start")
	gasUsed := block.GasUsed()
	baseFeePerGas := block.BaseFee()
	burntFees := new(big.Int).Mul(baseFeePerGas, new(big.Int).SetUint64(gasUsed))
	
	extraData := hex.EncodeToString(block.Extra())
				// Create block structure
	b := Block{
		BlockNumber: block.Number().Uint64(),  		
		FeeRecipient: block.Coinbase().Hex(),  
		BlockSize: block.Size(),			
		GasUsed: gasUsed,     	
		GasLimit: block.GasLimit(),     	
		BaseFeePerGas: baseFeePerGas,
		BurntFees: burntFees,
		ExtraData: extraData,		
		BlockHash: block.Hash().Hex(),
		StateRoot: block.Root().Hex(), 
		Transactions: make([]string, 0),
	}
	c <- b
}
// Get Txs Data to th ethclient
func GetTxsData(client *ethclient.Client, header *types.Header, tx *types.Transaction, block *types.Block) (Transaction, error) {
	logger.Debug("GetTxsData: start")
	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		logger.Debug("GetTxsData: Can't get sender")
	}
	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		logger.Debug("GetTxsData: Can't get receipt")
	}
	// Int to string
	var status string
	if receipt.Status == 1 {
		status = "Success"
	} else if receipt.Status == 0 {
		status = "Fail"
	} else {
		status = "Unknown"
	}
	// Get data
	value := tx.Value()
	gasUsed := receipt.GasUsed
	gasPrice := tx.GasPrice()
	transactionFee := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice)
	timestamp := int64(block.Time())
	timeUTC := time.Unix(timestamp, 0).UTC()
	utcTimeFormatted := timeUTC.Format("2006-01-02 15:04:05 AM MST")
	
	// Create transaction structure
	t:= Transaction{
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

	return t, nil
}

// Get LatestBlock
func (m *Model) GetLatestBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOne().SetSort(bson.M{"blockNumber": -1})

	var result Block
	err := m.colBlock.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		logger.Debug("GetLatestBlockNumber: Nonexist Documents")
		return 0, nil
	} else if err != nil {
		logger.Debug("GetLatestBlockNumber: Can't get latestBlockNumber")
		return 0, err
	}
	logger.Debug(fmt.Sprintf("GetLatestBlockNumber: %d", result.BlockNumber))
	return result.BlockNumber, nil
}

// Save Header
func (m *Model) SaveHeader(header *Header) error {
	logger.Debug("SaveHeader: start")
	result, err := m.colHeader.InsertOne(context.Background(), header)
	if err != nil {
		logger.Debug("Can't Insert HeaderData")
		return err
	}
	logger.Debug(fmt.Sprintf("SaveHeader: HeaderInsertId %s", result.InsertedID))
	return nil
}

// Save Blcok
func (m *Model) SaveBlock(block *Block) error {
	logger.Debug("SaveBlock: start")
	result, err := m.colBlock.InsertOne(context.Background(), block)
	if err != nil {
		logger.Debug("SaveBlock: Can't Insert Block")
		return err
	}
	logger.Debug(fmt.Sprintf("SaveBlock: BlockInsertId %s", result.InsertedID))
	return nil
}

// Save Transaction
func (m *Model) SaveTransaction(transaction *Transaction) error {
	logger.Debug("SaveTransaction: start")
	filter := bson.D{{Key: "hash", Value: transaction.Hash}}
	opts := options.Replace().SetUpsert(true)
	result, err := m.colTransaction.ReplaceOne(context.Background(), filter, transaction, opts)
	if err != nil {
		logger.Debug("SaveTransaction: Can't Upsert Tx")
		return err
	}
	if result.MatchedCount > 0 {
		logger.Debug("SaveTransaction: Insert Done")
	} else {
		logger.Debug("SaveTransaction: Update Done")
	}
	return nil
}


