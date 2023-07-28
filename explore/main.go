package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rnd/goerliscan/explore/config"
	ctl "rnd/goerliscan/explore/controller"
	"rnd/goerliscan/explore/logger"
	"rnd/goerliscan/explore/model"
	rt "rnd/goerliscan/explore/router"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	cf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// logger 초기화
	if err := logger.InitLogger(cf); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	logger.Debug("ready server....")
	
	//model 모듈 선언
	if mod, err := model.NewModel(cf); err != nil {
		panic(err)
	} else if controller, err := ctl.NewCTL(mod); err != nil {
		logger.Error(fmt.Errorf("controller.New > %v", err))
		panic(err)
	} else if rt, err := rt.NewRouter(controller); err != nil {
		logger.Error(fmt.Errorf("router.NewRouter > %v", err))
		panic(err)
	} else {
		mapi := &http.Server{
			Addr:           ":8080",
			Handler:        rt.Idx(),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		g.Go(func() error {
			return mapi.ListenAndServe()
		})	

		// Graceful shutdown
		// Wait for either an error or termination signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Warn("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := mapi.Shutdown(ctx); err != nil {
			logger.Error("Server Shutdown:", err)
		}

		select {
			case <-ctx.Done():
				logger.Info("timeout of 5 seconds.")
		}

		logger.Info("Server exiting")
	}

	if err := g.Wait(); err != nil {
		logger.Error(err)
	}

}