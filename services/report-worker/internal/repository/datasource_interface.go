package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type DatasourceDB interface {
	GetData(ctx context.Context, in dto.GetDataParams) (dto.GetDataResult, error)
	Close()
}
