package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ummuys/reportify/pkg/logger"
	"github.com/ummuys/reportify/services/gateway/internal/web"
)

func main() {
	mainCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := logger.InitLogger("gateway")
	if err != nil {
		log.Fatal(err)
	}

	// START SERVER
	server := web.CreateServer()
	errsCh := make(chan error, 4)
	srvOff := make(chan struct{})

	var wg sync.WaitGroup

	wg.Go(func() {
		<-mainCtx.Done()
		defer func() { _ = server.Close() }()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			errsCh <- fmt.Errorf("server shutdown: %w", err)
		}
		close(srvOff)
	})

	wg.Go(func() {
		if err := web.RunServer(server); err != nil {
			errsCh <- fmt.Errorf("server start failed: %w", err)
		}
	})

	wg.Wait()
	close(errsCh)

	var hadErr bool
	for err := range errsCh {
		if err != nil {
			hadErr = true
			logger.Error().Err(err).Send()
		}
	}

	if hadErr {
		logger.Error().Msg("graceful shutdown completed with errors")
	} else {
		logger.Info().Msg("shutdown successful")
	}
}
