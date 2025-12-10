package kafkacli

import (
	"github.com/IBM/sarama"
)

type producer struct {
	p sarama.SyncProducer
}
