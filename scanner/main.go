package main

import (
	"fmt"
	"os"
	"os/signal"
	"rnd/goerliscan/scanner/config"
	ctl "rnd/goerliscan/scanner/controller"
	"rnd/goerliscan/scanner/logger"
	"rnd/goerliscan/scanner/model"
	"syscall"
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

	// Graceful shutdown
	// Start scanner processing in a separate goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		err := ctl.StartScanner(cf, blockModel, transactionModel)
		if err != nil {
			logger.Warn("StartScanner error")
		}
	}()
	// Create a channel to receive termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Wait for a termination signal
	<-quit
	logger.Warn("Shutting down the scanner...")
	// Wait for the scanner goroutine to finish
	<-done
	logger.Info("Scanner exited gracefully.")
}
