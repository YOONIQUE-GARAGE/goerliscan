package model_test

import (
	"rnd/goerliscan/scanner/config"
	"rnd/goerliscan/scanner/model"
	"testing"

	"github.com/go-playground/assert"
)

// func TestNewModel(t *testing.T) {
// 	// 모의 MongoDB 커넥션 생성
// 	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
// 	assert.Equal(t, nil, err)
// 	defer client.Disconnect(context.Background())

// 	err = client.Ping(context.Background(), readpref.Primary())
// 	assert.Equal(t, nil, err)

// 	cf := &config.Config{
// 			Database: config.DatabaseConfig{
// 					Host: "mongodb://localhost:27017",
// 					Name: "goerliscan",
// 			},
// 	}

// 	db := client.Database(cf.Database.Name)
// 	md, err := model.NewModel(cf)
// 	assert.Equal(t, nil, err)
// 	assert.NotEqual(t, nil, md)

// 	// collection 이름 확인
// 	assert.Equal(t, "header", db.Collection("header").Name())
// 	assert.Equal(t, "block", db.Collection("block").Name())
// 	assert.Equal(t, "transaction", db.Collection("transaction").Name())
// }

// func GetClient() *ethclient.Client{
// 	client, err := ethclient.Dial("wss://goerli.infura.io/ws/v3/6b97fa21eed84c918b8283b684674cff")
// 	if err != nil {
// 		logger.Error(err)
// 		panic(err)
// 	}
// 	return client
// }

// func GetTypes() (*types.Header, *types.Block, *types.Transaction) {
// 	client := GetClient()
// 	// Make a go channel
// 	header := make(chan *types.Header)
// 	sub, err := client.SubscribeNewHead(context.Background(), header)
// 	if err != nil {
// 		logger.Error(err)
// 		panic(err)
// 	}

// 	select {
// 	case err := <-sub.Err():
// 		logger.Warn(err)
// 	case header := <-header:
// 		block, _ := client.BlockByHash(context.Background(), header.Hash())
// 		tx, _ := client.TransactionInBlock(context.Background(), block.Hash(), 0)
// 		return header, block, tx
// 	}
// 	return nil, nil, nil
// }

// func TestGetHeaderData(t *testing.T) {
// 	header, block, _ := GetTypes()
// 	timestamp := int64(block.Time())
// 	timeUTC := time.Unix(timestamp, 0).UTC()
// 	utcTimeFormatted := timeUTC.Format("2006-01-02 15:04:05 AM MST")
// 	nonce := binary.BigEndian.Uint64(header.Nonce[:])
// 	hexNonce := fmt.Sprintf("0x%016x", nonce)
// 	blockNumber := block.Number().Uint64()
// 	parentHash := block.ParentHash().Hex()
// 	bloom := header.Bloom[:]
// 	time := utcTimeFormatted

// 	// Create a channel to receive data from GetHeaderData
// 	h := make(chan model.Header)
// 	go model.GetHeaderData(header, block, h)
// 	data := <-h
// 	getHeaderData := data

// 	assert.Equal(t, hexNonce, getHeaderData.Nonce)
// 	assert.Equal(t, time, getHeaderData.Time)
// 	assert.Equal(t, parentHash, getHeaderData.ParentHash)
// 	assert.Equal(t, blockNumber, getHeaderData.BlockNumber)
// 	assert.Equal(t, bloom, getHeaderData.Bloom)
// }

// func TestGetBlockData(t *testing.T) {
// 	header, block, _ := GetTypes()
// 	gasUsed := block.GasUsed()
// 	baseFeePerGas := block.BaseFee()
// 	burntFees := new(big.Int).Mul(baseFeePerGas, new(big.Int).SetUint64(gasUsed))
// 	extraData := hex.EncodeToString(block.Extra())
// 	BlockNumber := block.Number().Uint64()
// 	FeeRecipient:= block.Coinbase().Hex()
// 	BlockSize:= block.Size()
// 	GasUsed:= gasUsed
// 	GasLimit:= block.GasLimit()
// 	BaseFeePerGas:= baseFeePerGas
// 	BurntFees:= burntFees
// 	ExtraData:= extraData
// 	BlockHash:= block.Hash().Hex()
// 	StateRoot:= block.Root().Hex()

// 	b := make(chan model.Block)
// 	go model.GetBlockData(header, block, b)
// 	data := <-b
// 	GetBlockData := data

// 	assert.Equal(t, BlockNumber, GetBlockData.BlockNumber)
// 	assert.Equal(t, FeeRecipient, GetBlockData.FeeRecipient)
// 	assert.Equal(t, BlockSize, GetBlockData.BlockSize)
// 	assert.Equal(t, GasUsed, GetBlockData.GasUsed)
// 	assert.Equal(t, GasLimit, GetBlockData.GasLimit)
// 	assert.Equal(t, BaseFeePerGas, GetBlockData.BaseFeePerGas)
// 	assert.Equal(t, BurntFees, GetBlockData.BurntFees)
// 	assert.Equal(t, ExtraData, GetBlockData.ExtraData)
// 	assert.Equal(t, BlockHash, GetBlockData.BlockHash)
// 	assert.Equal(t, StateRoot, GetBlockData.StateRoot)
// }

// func TestGetTxData(t *testing.T) {
// 	client := GetClient()
// 	header, block, tx := GetTypes()

// 	timestamp := int64(block.Time())
// 	timeUTC := time.Unix(timestamp, 0).UTC()
// 	utcTimeFormatted := timeUTC.Format("2006-01-02 15:04:05 AM MST")

// 	Hash := tx.Hash().Hex()
// 	Time := utcTimeFormatted
// 	To := tx.To().Hex()
// 	Value := tx.Value()
// 	GasPrice := tx.GasPrice()
// 	GasLimit := tx.Gas()
// 	BlockHash := block.Hash().Hex()
// 	BlockNumber := block.Number().Uint64()

// 	txData, err := model.GetTxsData(client, header, tx, block)
// 	if err!=nil {
// 		t.Error(err)
// 	}

// 	assert.Equal(t, Hash, txData.Hash)
// 	assert.Equal(t, Time, txData.Time)
// 	assert.Equal(t, To, txData.To)
// 	assert.Equal(t, Value, txData.Value)
// 	assert.Equal(t, GasPrice, txData.GasPrice)
// 	assert.Equal(t, GasLimit, txData.GasLimit)
// 	assert.Equal(t, BlockHash, txData.BlockHash)
// 	assert.Equal(t, BlockNumber, txData.BlockNumber)
// }

// func TestGetLatestBlockNumber(t *testing.T){
// 	cf := &config.Config{
// 		Database: config.DatabaseConfig{
// 				Host: "mongodb://localhost:27017",
// 				Name: "goerliscan",
// 		},
// 	}

//   md, err := model.NewModel(cf)
// 	assert.Equal(t, nil, err)
// 	assert.NotEqual(t, nil, md)

// 	originBlockNum := 9491504
// 	latestBlockNumber := uint64(originBlockNum)
// 	blockNumber, _ := md.GetLatestBlockNumber()
// 	assert.Equal(t, latestBlockNumber, blockNumber)
// }


func TestSaveHeader(t *testing.T) {
	cf := &config.Config{
		Database: config.DatabaseConfig{
				Host: "mongodb://localhost:27017",
				Name: "goerliscan",
		},
	}
  md, err := model.NewModel(cf)
	if err != nil {
		t.Errorf("Error creating Model instance: %v", err)
	}
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, md)
	oringinBlockNumber := 2
	blockNumber := uint64(oringinBlockNumber)
	
	h := model.Header{
		BlockNumber: blockNumber,  	
		ParentHash:	"parenthash",	 
		Bloom: []byte(""),
		Time: "", 
		Nonce: "",  
	}    
	md.SaveHeader(&h)

	blockNumberDB, _ := model.FindBlockNumber(md, blockNumber)
	assert.Equal(t, blockNumberDB, blockNumber)
}
