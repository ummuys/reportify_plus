package kafkacli

import "context"

type ReportCreateConsumer interface {
	Run(ctx context.Context) error
}
