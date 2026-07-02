package kafkacli

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

type producer struct {
	cli    *kgo.Client
	logger zerolog.Logger
	tunnel <-chan dto.KafkaMessage

	topic    string
	topicDLQ string
}

func NewKafkaProducer(tunnelReader chan dto.KafkaMessage, baseLogger zerolog.Logger) (KafkaProducer, error) {
	cfg, err := config.ParseKafkaProducerConfig()
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

	return &producer{
		cli: cli, logger: logger, tunnel: tunnelReader,
		topic: cfg.Topic, topicDLQ: cfg.TopicDLQ,
	}, nil
}

func (p *producer) Run(ctx context.Context) error {
	for {
		select {

		case <-ctx.Done():
			return nil
		case msg, ok := <-p.tunnel:
			if !ok {
				return errors.New("kafka tunnel is closed")
			}
			rec := &kgo.Record{
				Topic: p.topic,
				Key:   msg.Key,
				Value: msg.Value,
			}

			res := p.cli.ProduceSync(ctx, rec)
			if err := res.FirstErr(); err != nil {
				p.logger.Error().Err(err)
				continue
			}
		}
	}
}

func (p *producer) Close() {
	p.cli.Close()
}
