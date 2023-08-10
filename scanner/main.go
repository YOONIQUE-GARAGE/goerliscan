package main

import (
	"context"
	"os"
	"os/signal"
	"rnd/goerliscan/scanner/config"
	ctl "rnd/goerliscan/scanner/controller"
	"rnd/goerliscan/scanner/logger"
	"rnd/goerliscan/scanner/model"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	cf, err := config.LoadCofig()
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	// logger 초기화
	if err := logger.InitLogger(cf); err != nil {
		logger.Error("init logger failed, err:%v\n", err)
		return
	}
	logger.Debug("ready server....") 

	// Create model 
	model, err := model.NewModel(cf)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	// Graceful shutdown
	// Start scanner processing in a separate goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		// err := ctl.StartScanner(cf, headerModel, blockModel, transactionModel, mode)
		err := ctl.StartScanner(cf, model)
		if err != nil {
			logger.Warn("StartScanner error")
		}
	}()
	// Create a channel to receive termination signals & 	// Wait for a termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("Shutting down the scanner...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-ctx.Done()
	logger.Info("Timeout of 5 seconds for graceful shutdown.")

	logger.Info("Server exiting")

	if err := g.Wait(); err != nil {
		logger.Error(err)
	}
}


