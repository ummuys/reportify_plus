package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/logger"
	pkg "github.com/ummuys/reportify/pkg/tm"
	"github.com/ummuys/reportify/services/gateway/internal/di"
	"github.com/ummuys/reportify/services/gateway/internal/web"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logs, err := logger.InitLogger("gateway", "LOG_LEVEL_GATEWAY")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseGatewayConfig()
	if err != nil {
		logs.Fatal().Err(err).Msg("config")
	}

	// CLIENTS
	sc, err := di.NewGRPCServiceClients(cfg.AuthServiceAddr)
	if err != nil {
		logs.Fatal().Err(err).Msg("clients")
	}

	rh := di.NewRESTHandlers(sc, logs)

	tm, err := pkg.NewTokenManager()
	if err != nil {
		logs.Fatal().Err(err).Msg("token-manager")
	}

	// START SERVER
	server := web.CreateServer(cfg, rh, tm, logs)
	errsCh := make(chan error, 4)
	srvOff := make(chan struct{})

	var wg sync.WaitGroup

	wg.Go(func() {
		<-ctx.Done()
		defer func() { _ = server.Close() }()

		sdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(sdCtx); err != nil {
			errsCh <- fmt.Errorf("server shutdown: %w", err)
		}
		close(srvOff)
	})

	wg.Go(func() {
		logs.Info().Msg("run the rest-server")
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
			logs.Error().Err(err).Send()
		}
	}

	if hadErr {
		logs.Error().Msg("graceful shutdown completed with errors")
	} else {
		logs.Info().Msg("graceful shutdown")
	}
}
