package model

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"rnd/goerliscan/scanner/config"
	"rnd/goerliscan/scanner/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	colBlock      *mongo.Collection
	colTransaction *mongo.Collection
	db            *mongo.Database
}

type Block struct {
	BlockNumber  	uint64        `bson:"blockNumber"`
	Time         	string        `bson:"timestamp"`
	FeeRecipient  string	      `bson:"feeRecipient"`
	BlockSize			uint64        `bson:"blockSize"`
	GasUsed      	uint64        `bson:"gasUsed"`
	GasLimit     	uint64        `bson:"gasLimit"`
	BaseFeePerGas *big.Int    `bson:"baseFeePerGas"`
	BurntFees			*big.Int		`bson:"burntFees"`
	ExtraData			string				`bson:"extraData"`
	BlockHash    	string        `bson:"blockHash"`
	ParentHash		string        `bson:"parentHash"`
	StateRoot     string				`bson:"stateRoot"`
	Nonce        	string        `bson:"nonce"`
	Transactions 	[]string `bson:"transactions"`
}

type Transaction struct {
	Hash        		string  	 `bson:"hash"`
	Status      		string  	 `bson:"status"`
	Time        		string  	 `bson:"timestamp"`
	From        		string  	 `bson:"from"`
	To          		string  	 `bson:"to"` // return nil for contract
	Value      			*big.Int	 `bson:"amount"`
	TransactionFee 	*big.Int `bson:"transactionFee"`
	GasPrice    		*big.Int   `bson:"gasPrice"`
	GasUsed					uint64	 `bson:"gasUsed"`
	GasLimit    		uint64  	 `bson:"gasLimit"`
	BlockHash   		string  	 `bson:"blockHash"`
	BlockNumber 		uint64  	 `bson:"blockNumber"`
}

// Setting Model
func NewModel(config *config.Config) (*Model, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Database.Host))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db := client.Database(config.Database.Name)
	colBlock := db.Collection("block")
	colTransaction := db.Collection("transaction")
	model := &Model{
		colBlock:      colBlock,
		colTransaction: colTransaction,
		db:            db,
	}
	return model, nil
}

// Get LatestBlock
func (m *Model) GetLatestBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.FindOne().SetSort(bson.M{"blockNumber": -1})

	var result Block
	err := m.colBlock.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// No block found in the database, return 0 as the latest block number
		return 0, nil
	} else if err != nil {
		// log.Fatal(err)
		return 0, err
	}

	return result.BlockNumber, nil
}

// Save BlcokInfo
func (m *Model) SaveBlock(block *Block) error {
	_, err := m.colBlock.InsertOne(context.Background(), block)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// Save TransactionInfo
func (m *Model) SaveTransaction(transaction *Transaction) error {
	result, err := m.colTransaction.InsertOne(context.Background(), transaction)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Info(fmt.Sprintf("insertId: %s", result.InsertedID))
	return nil
}

func (m *Model) RemoveTxs(blockNumber uint64) error {
	opts := []*options.DeleteOptions{} 
	filter := bson.M{"blockNumber": blockNumber}
	
	result, err := m.colTransaction.DeleteMany(context.Background(), filter, opts...)
	if err != nil {
		return err
	} 

	logger.Info(fmt.Sprintf("DeleteCount: %d", result.DeletedCount))
	return nil
}
