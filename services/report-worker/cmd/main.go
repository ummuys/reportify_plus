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
	"github.com/ummuys/reportify/services/report-worker/internal/cache"
	"github.com/ummuys/reportify/services/report-worker/internal/convert"
	"github.com/ummuys/reportify/services/report-worker/internal/kafkacli"
	"github.com/ummuys/reportify/services/report-worker/internal/miniocli"
	"github.com/ummuys/reportify/services/report-worker/internal/repository"
	"github.com/ummuys/reportify/services/report-worker/internal/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logs, err := logger.InitLogger("report-worker", "LOG_LEVEL_REPORT_WORKER")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseReportWorkerConfig()
	if err != nil {
		log.Fatal(err)
	}

	reportDB, err := repository.NewReportDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("report-db")
	}
	defer reportDB.Close()

	reportCache, err := cache.NewReportCache(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("report-cache")
	}

	dataDB, err := repository.NewDatasourceDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("data-db")
	}
	defer dataDB.Close()

	minioCli, err := miniocli.NewMinIOCli(logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("miniocli")
	}

	conv := convert.NewReportConvert(logs)
	svc, err := service.NewPublishService(dataDB, reportDB, reportCache, conv, minioCli, cfg.FileTTL, cfg.DeleteBatch, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("service")
	}

	kafkacli, err := kafkacli.NewKafkaConsumer(svc, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("kafka-consumer")
	}

	SDChan := make(chan errs.SDMsg, 2)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		logs.Info().Msg("Worker is running")
		err := kafkacli.Run(ctx)
		if err != nil {
			SDChan <- errs.SDMsg{
				Err:  err,
				From: "kafka-consumer",
			}
		}
	})

	// HERE
	wg.Go(func() {
		t := time.NewTicker(cfg.FileTTL)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				svc.CleanOldReports(ctx)
			}
		}
	})

	wg.Wait()
	close(SDChan)
	errs.ShutdownStatus(logs, SDChan)
}
