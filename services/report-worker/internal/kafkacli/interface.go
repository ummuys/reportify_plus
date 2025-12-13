package kafkacli

import "context"

type KafkaConsumer interface {
	Run(ctx context.Context) error
}
