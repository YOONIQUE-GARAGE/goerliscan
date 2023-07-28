package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"rnd/goerliscan/explore/config"
	"rnd/goerliscan/explore/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	BlockNumber  	uint64      `bson:"blockNumber"`
	Time         	string      `bson:"timestamp"`
	FeeRecipient  string	    `bson:"feeRecipient"`
	BlockSize			uint64      `bson:"blockSize"`
	GasUsed      	uint64      `bson:"gasUsed"`
	GasLimit     	uint64      `bson:"gasLimit"`
	BaseFeePerGas *big.Int    `bson:"baseFeePerGas"`
	BurntFees			*big.Int		`bson:"burntFees"`
	ExtraData			string		  `bson:"extraData"`
	BlockHash    	string      `bson:"blockHash"`
	ParentHash		string      `bson:"parentHash"`
	StateRoot     string			`bson:"stateRoot"`
	Nonce        	string      `bson:"nonce"`
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

// Allinfo struct
type Result struct {
	Blocks      []Block
	Transactions []Transaction
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
	
func (m *Model) Check(c *gin.Context) {
	m.RespOK(c, 0)
}

func (m *Model) RespOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

// get Blocks and Transactions
func (m *Model) GetAll() (Result, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// findBlocks with fiter
	filter := bson.D{}
	cursor, err := m.colBlock.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	// block into blocks 
	var blocks []Block
	for cursor.Next(ctx) {
		var block Block
		if err := cursor.Decode(&block); err != nil {
			panic(err)
		}
		blocks = append(blocks, block)
		_, err := json.MarshalIndent(block, "prefix string", " ")
		if err != nil {
			panic(err)
		}
	}

	cursor, err = m.colTransaction.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	// tx into txs 
	var txs []Transaction
	for cursor.Next(ctx) {
		var tx Transaction
		if err := cursor.Decode(&tx); err != nil {
			panic(err)
		}
		txs = append(txs, tx)
		_, err := json.MarshalIndent(tx, "prefix string", " ")
		if err != nil {
			panic(err)
		}
	}

	result := Result{
		Blocks:      blocks,
		Transactions: txs,
	}
	return result, err
}

func (m *Model) GetMore(elem string) (interface{}, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{}
	// elem에 따른 데이터 분기
	if elem == "blocks" {
		cursor, err := m.colBlock.Find(ctx, filter)
		if err != nil {
			panic(err)
		}
		// block into blocks 
		var blocks []Block
		for cursor.Next(ctx) {
			var block Block
			if err := cursor.Decode(&block); err != nil {
				panic(err)
			}
			blocks = append(blocks, block)
			_, err := json.MarshalIndent(block, "prefix string", " ")
			if err != nil {
				logger.Error(err)
				panic(err)
			}
		}
		return blocks, nil
	} else if elem == "txs" {
		cursor, err := m.colTransaction.Find(ctx, filter)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
		// tx into txs 
		var txs []Transaction
		for cursor.Next(ctx) {
			var tx Transaction
			if err := cursor.Decode(&tx); err != nil {
				logger.Error(err)
				panic(err)
			}
			txs = append(txs, tx)
			output, err := json.MarshalIndent(tx, "prefix string", " ")
			if err != nil {
				logger.Error(err)
				panic(err)
			}
			logger.Debug(output)
		}
		return txs, nil
	} else {
		logger.Error("transaction with hash %s not found")
		return fmt.Errorf("transaction with hash %s not found", elem), nil
	}
}

// MainPage에서 blockNumber 클릭 시
func (m *Model) GetOneBlcok(elem string) (Block, error) {
	opts := []*options.FindOneOptions{}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// parameter에 다른 분기
	elemToint, err := strconv.ParseInt(elem, 10, 64)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"blockNumber": elemToint}
		
	var block Block
	if err := m.colBlock.FindOne(ctx, filter, opts...).Decode(&block); err != nil {
		return block, err
	} else {
		return block, nil
	}
}

// MainPage에서 txHash 클릭 시
func (m *Model) GetOneTransaction(elem string) (Transaction, error) {
	opts := []*options.FindOneOptions{}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"hash": elem}

	var tx Transaction
	if err := m.colTransaction.FindOne(ctx, filter, opts...).Decode(&tx); err != nil {
		return tx, err
	} else {
		return tx, nil
	}
}

// from이나 to의 address로 조회시 
// func (b *Model) GetTransactions(elem string) ([]Transaction, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	// findBlocks with filter
// 	filter := bson.M{
// 		"transactions": bson.M{
// 			"$elemMatch": bson.M{
// 				"$or": []bson.M{
// 					{"from": elem},
// 					{"to": elem},
// 				},
// 			},
// 		},
// 	}
// 	cursor, err := b.colBlock.Find(ctx, filter)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// transactions into slice
// 	var transactions []Transaction
// 	for cursor.Next(ctx) {
// 		var block Block
// 		if err := cursor.Decode(&block); err != nil {
// 			panic(err)
// 		}

// 		// Filter transactions to include only those that match the condition
// 		for _, tx := range block.Transactions {
// 			if tx.From == elem || tx.To == elem {
// 				transactions = append(transactions, tx)
// 			}
// 		}

// 		output, err := json.MarshalIndent(block, "prefix string", " ")
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("%s\n", output)
// 	}
// 	return transactions, err
// }
