package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	reportv1 "github.com/ummuys/reportify/api/pb/report/v1"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logs, err := logger.InitLogger("report", "LOG_LEVEL_AUTH")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseReportServiceConfig()
	if err != nil {
		logs.Fatal().Err(err).Msg("config")
	}

	lis, err := net.Listen(cfg.Network, cfg.Port)
	if err != nil {
		logs.Fatal().Err(err).Msg("listener")
	}

	srv := grpc.NewServer()
	reportv1.RegisterReportServiceServer(srv, nil)

	wg := sync.WaitGroup{}

	wg.Go(func() {
		<-ctx.Done()
		srv.GracefulStop()
	})

	wg.Go(func() {
		logs.Info().Msg("run the grpc-server")
		if err := srv.Serve(lis); err != nil {
			logs.Fatal().Err(err).Msg("shutdown grpc-server")
		}
	})

	wg.Wait()
	logs.Info().Msg("graceful shutdown")
}
