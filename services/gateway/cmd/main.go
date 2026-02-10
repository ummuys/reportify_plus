package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/pkg/logger"
	pkg "github.com/ummuys/reportify/pkg/tm"
	"github.com/ummuys/reportify/services/gateway/internal/di"
	"github.com/ummuys/reportify/services/gateway/internal/web"
)

// @title           Reportify API
// @version         1.0
// @description     API for report generation platform

// @host            127.0.0.1:8088
// @BasePath        /api/v1

// @schemes         http

// JWT Bearer
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization

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

	sc, err := di.NewGRPCServiceClients(cfg.AuthServiceAddr, cfg.ReportServiceAddr)
	if err != nil {
		logs.Fatal().Err(err).Msg("clients")
	}

	rh := di.NewRESTHandlers(sc, logs)

	tm, err := pkg.NewTokenManager()
	if err != nil {
		logs.Fatal().Err(err).Msg("token-manager")
	}

	server := web.CreateServer(cfg, rh, tm, logs)
	SDChan := make(chan errs.SDMsg, 3)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		<-ctx.Done()
		defer func() { _ = server.Close() }()

		sdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(sdCtx); err != nil {
			SDChan <- errs.SDMsg{
				Err:  err,
				From: "off rest-server",
			}
		}
	}()

	go func() {
		defer wg.Done()
		logs.Info().Msg("Service is running")
		if err := web.RunServer(server); err != nil {
			SDChan <- errs.SDMsg{
				Err:  err,
				From: "rest-server",
			}
		}
	}()

	wg.Wait()
	close(SDChan)
	errs.ShutdownStatus(logs, SDChan)
}
