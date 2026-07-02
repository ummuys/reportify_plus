package kafkacli

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
	"github.com/ummuys/reportify/services/report-worker/internal/service"
)

type consumer struct {
	svc      service.PublishService
	logger   zerolog.Logger
	cli      *kgo.Client
	topic    string
	topicDLQ string
	cWorkers int
}

func NewReportCreateConsumer(svc service.PublishService, baseLogger zerolog.Logger) (ReportCreateConsumer, error) {
	cfg, err := config.ParseReportCreateConsumerConfig()
	if err != nil {
		return nil, err
	}

	cli, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Broker),
		kgo.ConsumerGroup(cfg.Group),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.ClientID(cfg.ClientID),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return nil, err
	}

	logger := baseLogger.With().Str("component", "kafka").Logger()

	return &consumer{
		cli:      cli,
		logger:   logger,
		svc:      svc,
		topic:    cfg.Topic,
		topicDLQ: cfg.TopicDLQ,
		cWorkers: 3,
	}, nil
}

func (c *consumer) Run(ctx context.Context) error {
	c.logger.Debug().Str("evt", "call Run").Msg("")

	records := make(chan *kgo.Record, 512)
	defer close(records)

	wg := sync.WaitGroup{}
	for i := 0; i < c.cWorkers; i++ {
		wg.Go(func() {
			for r := range records {
				func() {
					wctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
					defer cancel()

					var payload dto.KafkaMessage
					if err := json.Unmarshal(r.Value, &payload); err != nil {
						c.toDLQ(ctx, payload.ReportID, r, err)
						return
					}

					c.logger.Info().
						Str("topic", r.Topic).
						Str("report_id", payload.ReportID).
						Msg("catch new message")

					msg := dto.KafkaMessage{ReportID: payload.ReportID}

					var err error
					if payload.Recreating {
						err = c.svc.RecreateReport(wctx, msg)
					} else {
						err = c.svc.CreateReport(wctx, msg)
					}

					if err != nil {
						c.toDLQ(ctx, payload.ReportID, r, err)
						return
					}

					c.cli.MarkCommitRecords(r)

					c.logger.Info().
						Str("topic", r.Topic).
						Str("report_id", payload.ReportID).
						Msg("report created")
				}() 
			}
		})
	}

	for {
		fetches := c.cli.PollFetches(ctx)
		if fetches.IsClientClosed() {
			wg.Done()
			return nil
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, fe := range errs {
				c.logger.Warn().
					Err(fe.Err).
					Str("topic", fe.Topic).
					Int32("partition", fe.Partition).
					Msg("poll fetch error")
			}
			continue
		}

		fetches.EachRecord(func(r *kgo.Record) {
			records <- r
		})
	}
}

func (c *consumer) toDLQ(ctx context.Context, reportID string, r *kgo.Record, err error) {
	rec := &kgo.Record{
		Topic: c.topicDLQ,
		Key:   nil,
		Value: r.Value,
		Headers: append(
			r.Headers,
			kgo.RecordHeader{Key: "dlq_error", Value: []byte(err.Error())},
			kgo.RecordHeader{Key: "src_topic", Value: []byte(r.Topic)},
		),
	}

	res := c.cli.ProduceSync(ctx, rec)
	if perr := res.FirstErr(); perr != nil {
		c.logger.Error().
			Err(perr).
			Str("op", "dlq.produce").
			Str("topic", c.topicDLQ).
			Str("report_id", reportID).
			Msg("send to dlq failed")
		return
	}

	c.logger.Warn().
		Err(err).
		Str("topic", r.Topic).
		Str("report_id", reportID).
		Int32("partition", r.Partition).
		Msg("sent to dlq")
}
