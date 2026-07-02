package miniocli

import (
	"context"

	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type MinIOClient interface {
	UploadAndPresign(ctx context.Context, in dto.PutReportIn) (string, error)
	DeleteFiles(ctx context.Context, in dto.DeleteExpiredFilesParams) error
}
