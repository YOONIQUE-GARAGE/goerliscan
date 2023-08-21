package model

import (
	"context"
	"encoding/json"
	"fmt"
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
	colHeader 			*mongo.Collection
	colBlock      	*mongo.Collection
	colTransaction 	*mongo.Collection
	db            	*mongo.Database
}

type Header struct {
	BlockNumber  	uint64      `bson:"blockNumber"`
	ParentHash		string      `bson:"parentHash"`
	Bloom					string 			`bson:"bloom"`
	Time         	string      `bson:"time"`
	Nonce        	string      `bson:"nonce"`
}

type Block struct {
	BlockNumber  	uint64      `bson:"blockNumber"`
	Miner  				string	    `bson:"miner"`
	BlockSize			uint64      `bson:"blockSize"`
	GasUsed      	uint64      `bson:"gasUsed"`
	GasLimit     	uint64      `bson:"gasLimit"`
	BaseFeePerGas uint64      `bson:"baseFeePerGas"`
	BurntFees			uint64		  `bson:"burntFees"`
	ExtraData			string		  `bson:"extraData"`
	BlockHash    	string      `bson:"blockHash"`
	StateRoot     string			`bson:"stateRoot"`
	Transactions 	[]string 		`bson:"transactions"`
}

type Transaction struct {
	Hash        		string  	 `bson:"hash"`
	Status      		string  	 `bson:"status"`
	Time        		string  	 `bson:"time"`
	From        		string  	 `bson:"from"`
	To          		string  	 `bson:"to"` // return nil for contract
	Value      			uint64	 	 `bson:"amount"`
	TransactionFee 	uint64 	   `bson:"transactionFee"`
	GasPrice    		uint64     `bson:"gasPrice"`
	GasUsed					uint64 	 	 `bson:"gasUsed"`
	BlockHash   		string  	 `bson:"blockHash"`
	BlockNumber 		uint64  	 `bson:"blockNumber"`
}

// Allinfo struct
type AllData struct {
	Headers 			[]Header
	Blocks      	[]Block
	Transactions	[]Transaction
}

type OneBlock struct {
	Header		Header
	Block   	Block
}

// Setting Model
func NewModel(config *config.Config) (*Model, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Database.Host))
	if err != nil {
		logger.Debug("Can't get mongoClient")
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
	
func (m *Model) Check(c *gin.Context) {
	m.RespOK(c, 0)
}

func (m *Model) RespOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

// get Blocks and Transactions
func (m *Model) GetAll() (AllData, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{}
	
	// Header
	cursor, err := m.colHeader.Find(ctx, filter)
	if err != nil {
		logger.Debug("Can't get allHeaders")
	}
	// header into headers 
	var headers []Header
	for cursor.Next(ctx) {
		var header Header
		if err := cursor.Decode(&header); err != nil {
			logger.Debug("Header: Can't unmarshalling")
		}
		headers = append(headers, header)
		_, err := json.MarshalIndent(header, "prefix string", " ")
		if err != nil {
			logger.Debug("Header: Can't MarshalIndent")
		}
	}
	// Block
	cursor, err = m.colBlock.Find(ctx, filter)
	if err != nil {
		logger.Debug("Can't get allBlocks")
	}
	// block into blocks 
	var blocks []Block
	for cursor.Next(ctx) {
		var block Block
		if err := cursor.Decode(&block); err != nil {
			logger.Debug("Block: Can't unmarshalling")
		}
		blocks = append(blocks, block)
		_, err := json.MarshalIndent(block, "prefix string", " ")
		if err != nil {
			logger.Debug("Block: Can't MarshalIndent")
		}
	}
	// Transaction
	cursor, err = m.colTransaction.Find(ctx, filter)
	if err != nil {
		logger.Debug("Can't get allTxs")
	}
	// tx into txs 
	var txs []Transaction
	for cursor.Next(ctx) {
		var tx Transaction
		if err := cursor.Decode(&tx); err != nil {
			logger.Debug("Transaction: Can't unmarshalling")
		}
		txs = append(txs, tx)
		_, err := json.MarshalIndent(tx, "prefix string", " ")
		if err != nil {
			logger.Debug("Transaction: Can't MarshalIndent")
		}
	}

	result := AllData{
		Headers: headers,
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
		// header & block 
		cursor, err := m.colHeader.Find(ctx, filter)
		if err != nil {
			logger.Debug("Cant't get more header data")
		}
		// block into blocks 
		var blocks []Block
		for cursor.Next(ctx) {
			var block Block
			if err := cursor.Decode(&block); err != nil {
				logger.Debug("Can't get more block datas")
			}
			blocks = append(blocks, block)
			_, err := json.MarshalIndent(block, "prefix string", " ")
			if err != nil {
				logger.Debug("Can't MarshalIndent")
			}
		}
		return blocks, nil
	} else if elem == "txs" {
		cursor, err := m.colTransaction.Find(ctx, filter)
		if err != nil {
			logger.Debug("Cant't get more tx data")
		}
		// tx into txs 
		var txs []Transaction
		for cursor.Next(ctx) {
			var tx Transaction
			if err := cursor.Decode(&tx); err != nil {
				logger.Debug("Cant't get more tx data")
			}
			txs = append(txs, tx)
			_, err := json.MarshalIndent(tx, "prefix string", " ")
			if err != nil {
				logger.Debug("Can't MarshalIndent")
			}
		}
		return txs, nil
	} else {
		logger.Debug("transaction with hash %s not found", elem)
		return fmt.Errorf("transaction with hash %s not found", elem), nil
	}
}

// MainPage에서 blockNumber 클릭 시
func (m *Model) GetOneBlcok(elem string) (OneBlock, error) {
	opts := []*options.FindOneOptions{}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// parameter에 다른 분기
	elemToint, err := strconv.ParseInt(elem, 10, 64)
	if err != nil {
		logger.Debug("Can't parsing elem")
	}
	filter := bson.M{"blockNumber": elemToint}
		
	var header Header
	if err := m.colHeader.FindOne(ctx, filter, opts...).Decode(&header); err != nil {
		logger.Debug("No Header Documents")
	} 

	var block Block
	if err := m.colBlock.FindOne(ctx, filter, opts...).Decode(&block); err != nil {
		logger.Debug("No Block Documents")
	}
	result := OneBlock{
		Header: header,
		Block:  block,
	}
	return result, err
}

// MainPage에서 txHash 클릭 시
func (m *Model) GetOneTransaction(elem string) (Transaction, error) {
	opts := []*options.FindOneOptions{}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"hash": elem}

	var tx Transaction
	if err := m.colTransaction.FindOne(ctx, filter, opts...).Decode(&tx); err != nil {
		logger.Debug("No Tx Documents")
		return tx, err
	} else {
		return tx, nil
	}
}

