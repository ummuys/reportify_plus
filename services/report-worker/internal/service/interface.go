package service

import (
	"context"

	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type PublishService interface {
	CreateReport(ctx context.Context, in dto.KafkaMessage) error
	CleanOldReports(ctx context.Context)
}
