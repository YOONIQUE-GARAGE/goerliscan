package controller_test

import (
	"context"
	"rnd/goerliscan/scanner/config"
	"rnd/goerliscan/scanner/controller"
	"rnd/goerliscan/scanner/logger"
	"rnd/goerliscan/scanner/model"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-playground/assert"
)

func GetClient() *ethclient.Client{
	client, err := ethclient.Dial("wss://goerli.infura.io/ws/v3/6b97fa21eed84c918b8283b684674cff")
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return client
}

func GetTypes() (*types.Header, *types.Block, *types.Transaction) {
	client := GetClient()
	// Make a go channel
	header := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), header)
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	select {
	case err := <-sub.Err():
		logger.Warn(err)
	case header := <-header:
		block, _ := client.BlockByHash(context.Background(), header.Hash())
		tx, _ := client.TransactionInBlock(context.Background(), block.Hash(), 0)
		return header, block, tx
	}
	return nil, nil, nil
}

func TestSaveMissingBlock(t *testing.T) {
	client := GetClient()
	header, _, _ := GetTypes()
	cf := &config.Config{
		Database: config.DatabaseConfig{
				Host: "mongodb://localhost:27017",
				Name: "goerliscan",
		},
	}
	md, err := model.NewModel(cf)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, md)
	originalBlockNumber := 9491504
	originalBlockNumberDB := 9491499
	latestBlockNumber := uint64(originalBlockNumber)
	latestBlockNumberDB := uint64(originalBlockNumberDB)
	controller.SaveMissingBlock(client, header, latestBlockNumber , latestBlockNumberDB, md)
	BlockNumberDB, _  := md.GetLatestBlockNumber()
	assert.Equal(t, latestBlockNumber, BlockNumberDB)
}