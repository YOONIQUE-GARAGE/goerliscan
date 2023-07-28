package main

import (
	"fmt"
	"rnd/goerliscan/scanner/config"
	ctl "rnd/goerliscan/scanner/controller"
	"rnd/goerliscan/scanner/logger"
	"rnd/goerliscan/scanner/model"
)

func main() {
	cf, err := config.LoadCofig()
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	// logger 초기화
	if err := logger.InitLogger(cf); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	logger.Debug("ready server....")

	// Create separate models for blocks and transactions
	blockModel, err := model.NewModel(cf)
	if err != nil {
		logger.Error(err)
		panic(err)	
	}

	transactionModel, err := model.NewModel(cf)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	// Start scanner processing
	ctl.StartScanner(cf, blockModel, transactionModel)
}
