package kafkacli

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
	"github.com/ummuys/reportify/services/report-worker/internal/service"
)

type consumer struct {
	svc    service.PublishService
	logger zerolog.Logger
	cli    *kgo.Client

	topic    string
	topicDLQ string
	cWorkers int
}

func NewKafkaConsumer(svc service.PublishService, baseLogger zerolog.Logger) (KafkaConsumer, error) {
	cfg, err := config.ParseKafkaConsumerConfig()
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

	logger := baseLogger.With().Str("component", "kafka").Logger()

	if err != nil {
		return nil, err
	}

	return &consumer{cli: cli, logger: logger,
		svc:   svc,
		topic: cfg.Topic, topicDLQ: cfg.TopicDLQ,
		cWorkers: 3,
	}, nil
}

func (c *consumer) Run(ctx context.Context) error {
	records := make(chan *kgo.Record, 512)
	defer close(records)

	wg := sync.WaitGroup{}
	for i := 0; i < c.cWorkers; i++ {
		wg.Go(func() {
			for r := range records {
				c.logger.Info().Str("topic", r.Topic).Msg("catch new message")
				uuid := string(r.Value)
				if err := c.svc.CreateReport(ctx, dto.KafkaMessage{
					UUID: uuid,
				}); err != nil {
					c.toDLQ(ctx, r, err)
				}
				c.cli.MarkCommitRecords(r)
			}
		})
	}

	// Get a new message
	for {
		fetches := c.cli.PollFetches(ctx)
		if fetches.IsClientClosed() {
			wg.Done()
			return nil
		}
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Warn().Err(err.Err).Str("topic", err.Topic).Int32("partition", err.Partition).Msg("")
			}
			continue
		}

		fetches.EachRecord(func(r *kgo.Record) {
			records <- r
		})
	}
}

func (c *consumer) toDLQ(ctx context.Context, r *kgo.Record, err error) {
	c.logger.Warn().Err(err).Str("topic", r.Topic).Int32("partition", r.Partition).Msg("")
	rec := &kgo.Record{
		Topic: c.topicDLQ,
		Key:   nil,
		Value: r.Value,
		Headers: append(r.Headers,
			kgo.RecordHeader{Key: "dlq_error", Value: []byte(err.Error())},
			kgo.RecordHeader{Key: "src_topic", Value: []byte(r.Topic)},
		),
	}
	res := c.cli.ProduceSync(ctx, rec)
	if err := res.FirstErr(); err != nil {
		c.logger.Error().Err(err).Msg("send to dlq failed")
	}
}
