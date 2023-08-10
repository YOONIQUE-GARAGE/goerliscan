// main_test.go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestGracefulShutdown(t *testing.T) {

	// Create a channel to receive termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create a new errgroup.Group for testing
	var g errgroup.Group

	// Start the scanner processing in a separate goroutine
	done := make(chan struct{})
	g.Go(func() error {
		defer close(done)
		// Replace the following with your actual StartScanner function call
		err := StartScanner()
		return err
	})

	// Simulate sending a termination signal to the process
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	t.Logf("server shutdown")
	fmt.Println("server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Wait for the server to exit or timeout
	select{	
	case <-ctx.Done():
		t.Logf("Graceful shutdown completed")
	case <- time.After(10 * time.Second):
		t.Errorf("Graceful shutdown timeout")
	}
	

	// Wait for the errgroup to finish (optional)
	if err := g.Wait(); err != nil {
		t.Errorf("Error occurred during shutdown: %v", err)
	}
}

// Mock StartScanner function for testing
func StartScanner() error {
	return nil
}
