package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	dsv1 "github.com/ummuys/reportify/api/pb/datasource/service/v1"
	rcv1 "github.com/ummuys/reportify/api/pb/report/cache/v1"
	rsv1 "github.com/ummuys/reportify/api/pb/report/service/v1"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/pkg/logger"
	"github.com/ummuys/reportify/services/report/internal/adapter"
	"github.com/ummuys/reportify/services/report/internal/cache"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/kafkacli"
	"github.com/ummuys/reportify/services/report/internal/repository"
	"github.com/ummuys/reportify/services/report/internal/service"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logs, err := logger.InitLogger("report", "LOG_LEVEL_REPORT")
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

	reportDB, err := repository.NewReportDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("report-db")
	}
	defer reportDB.Close()

	datasourceDB, err := repository.NewDatasourceDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("datasource-db")
	}
	defer datasourceDB.Close()

	reportCache, err := cache.NewReportCache(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("report-cache")
	}

	reportSvc := service.NewReportService(reportDB, reportCache, logs)
	reportCacheSvc := service.NewReportCacheService(reportCache, logs)
	datasourceSVC := service.NewDatasourceService(datasourceDB, logs)

	srv := grpc.NewServer()

	tunnel := make(chan dto.KafkaMessage, 10)
	defer close(tunnel)

	kafkaCli, err := kafkacli.NewKafkaProducer(tunnel, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("kafkacli")
	}
	defer kafkaCli.Close()

	reportAdapter := adapter.NewReportAdapter(reportSvc, tunnel, logs)
	datasourceAdapter := adapter.NewDatasourceAdapter(datasourceSVC, logs)
	reportCacheAdapter := adapter.NewReportCacheAdapter(reportCacheSvc, logs)

	rsv1.RegisterReportServiceServer(srv, reportAdapter)
	dsv1.RegisterDatasourceServiceServer(srv, datasourceAdapter)
	rcv1.RegisterReportCacheServiceServer(srv, reportCacheAdapter)

	var wg sync.WaitGroup
	wg.Add(3)

	SDChan := make(chan errs.SDMsg, 2)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		srv.GracefulStop()
	}()

	go func() {
		defer wg.Done()
		if err := kafkaCli.Run(ctx); err != nil {
			SDChan <- errs.SDMsg{
				Err:  err,
				From: "kafka-producer",
			}
		}
	}()

	go func() {
		defer wg.Done()
		logs.Info().Msg("run the grpc-server")
		if err := srv.Serve(lis); err != nil {
			SDChan <- errs.SDMsg{
				Err:  err,
				From: "grpc-server",
			}
		}
	}()

	wg.Wait()
	close(SDChan)
	errs.ShutdownStatus(logs, SDChan)
}
