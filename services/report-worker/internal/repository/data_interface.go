package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type DataDB interface {
	GetData(ctx context.Context, in dto.GetDataParams) (dto.GetDataResult, error)
}
