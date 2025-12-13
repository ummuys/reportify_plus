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
	"github.com/ummuys/reportify/services/report/internal/adapter"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/kafkacli"
	"github.com/ummuys/reportify/services/report/internal/repository"
	"github.com/ummuys/reportify/services/report/internal/service"
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

	db, err := repository.NewReportDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("report-db")
	}

	svc := service.NewReportService(db, logs)

	srv := grpc.NewServer()

	tunnel := make(chan dto.KafkaMessage, 10)
	defer close(tunnel)

	kafkaCli, err := kafkacli.NewKafkaProducer(tunnel, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("kafkacli")
	}
	defer kafkaCli.Close()

	adapter := adapter.NewReportAdapter(svc, tunnel, logs)
	reportv1.RegisterReportServiceServer(srv, adapter)

	wg := sync.WaitGroup{}

	wg.Go(func() {
		<-ctx.Done()
		srv.GracefulStop()
	})

	wg.Go(func() {
		err := kafkaCli.Run(ctx)
		if err != nil {
			logs.Fatal().Err(err).Msg("shutdown kafka-producer")
		}
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
