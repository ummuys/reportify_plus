package kafkacli

import (
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/services/report-worker/internal/service"
)

type consumer struct {
	svc    service.PublishService
	logger zerolog.Logger
	cli    *kgo.Client

	topic    string
	topicDLQ string
}

func NewKafkaConsumer(svc service.PublishService, baseLogger zerolog.Logger) (KafkaConsumer, error) {
	cfg, err := config.ParseKafkaConsumerConfig()
	if err != nil {
		return nil, err
	}

	cli, err := kgo.NewClient(
		kgo.ClientID(cfg.ClientID),
		kgo.SeedBrokers(cfg.Broker),
	)

	logger := baseLogger.With().Str("component", "kafka").Logger()

	if err != nil {
		return nil, err
	}

	return &consumer{cli: cli, logger: logger,
		topic: cfg.Topic, topicDLQ: cfg.TopicDLQ,
	}, nil
}
