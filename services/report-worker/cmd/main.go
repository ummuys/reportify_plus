package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ummuys/reportify/pkg/logger"
	"github.com/ummuys/reportify/services/report-worker/internal/convert"
	"github.com/ummuys/reportify/services/report-worker/internal/kafkacli"
	"github.com/ummuys/reportify/services/report-worker/internal/repository"
	"github.com/ummuys/reportify/services/report-worker/internal/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logs, err := logger.InitLogger("report", "LOG_LEVEL_REPORT_WORKER")
	if err != nil {
		log.Fatal(err)
	}

	reportDB, err := repository.NewReportDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("report-db")
	}

	dataDB, err := repository.NewDataDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("data-db")
	}

	conv := convert.NewReportConvert(logs)

	svc, err := service.NewPublishService(dataDB, reportDB, conv, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("service")
	}

	kafkacli, err := kafkacli.NewKafkaConsumer(svc, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("kafka-consumer")
	}

	wg := sync.WaitGroup{}
	wg.Go(func() {
		logs.Info().Msg("run the report-worker")
		err := kafkacli.Run(ctx)
		if err != nil {
			logs.Error().Err(err).Msg("kafka-consumer")
		}
	})

	wg.Wait()
	logs.Info().Msg("graceful shutdown")

}
