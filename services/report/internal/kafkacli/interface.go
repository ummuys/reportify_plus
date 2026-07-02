package kafkacli

import "context"

type KafkaProducer interface {
	Run(ctx context.Context) error
	Close()
}
