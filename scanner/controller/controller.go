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

// 1. 함수 인자 확인(model, config)
func StartScanner(cf *config.Config, md *model.Model) error {
	// Initialize mode
	// Start :=   cf.Mode.Start
	// Restart := cf.Mode.Restart

	// Initialize ethclient
	client, err := ethclient.Dial(cf.Netowrk.URL)
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	// make go Channels
	header := make(chan *types.Header)
	// Subscribe to new block headers
	sub, err := client.SubscribeNewHead(context.Background(), header)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	for {
		select {
		case err := <-sub.Err():
			logger.Warn(err)
		case header := <-header:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				logger.Debug("BlockByHash: Can't get BlockByHash block")
				continue
			}
			// Get the latest block number from the database
			latestBlockNumberDB, err := md.GetLatestBlockNumber()
			if err != nil {
				logger.Debug("GetLatestBlockNumber: Can't get latestBlockNumberDB")
				continue
			}
			// Get the latest block number from the client
			latestBlockNumber := block.Number().Uint64()
		
			if latestBlockNumber > latestBlockNumberDB {
				SaveMissingBlock(client, header, latestBlockNumber, latestBlockNumberDB, md)
			}
		}
	}
}

func SaveMissingBlock(client *ethclient.Client, header *types.Header, latestBlockNumber uint64, latestBlockNumberDB uint64, md *model.Model) {
	if latestBlockNumber > latestBlockNumberDB {
	
	for i := latestBlockNumberDB + 1; i <= latestBlockNumber; i++ {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			// logger.Debug("BlockByNumber: Cant get BlockByNumber block")
			continue
		}

		// Process Header Data
		hChan := make(chan model.Header)
		go model.GetHeaderData(header, block, hChan)
		h := <-hChan
		clonedHeader := h
		// Process Block&Tx Data
		bChan := make(chan model.Block)
		go model.GetBlockData(header, block, bChan)
		b := <-bChan
		clonedBlock := b
		
		txs := block.Transactions()
		if len(txs) > 0 { 
			for _, tx := range txs {
				t, err := model.GetTxsData(client, header, tx, block)
				if err != nil {
					logger.Debug("GetTxsData: Can't get txsData")
					continue
				}
				clonedBlock.Transactions = append(clonedBlock.Transactions, t.Hash)
				// Save the transaction to the database
				err = md.SaveTransaction(&t)
				if err != nil {
					logger.Debug("SaveTransaction: Can't save the transaction")
					continue
				}
			}
			// Save the Header to the database	
			err = md.SaveHeader(&clonedHeader)
			if err != nil {
					logger.Debug("SaveHeader: Can't save header")
					continue
			}
			// Save the block to the database
			err = md.SaveBlock(&clonedBlock)
			if err != nil {
				logger.Debug("SaveBlock: Can't save the block")
				continue
			}
		}
		}
		// latestBlockNumber update
		latestBlockNumberDB = latestBlockNumber
	}
}

