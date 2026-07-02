package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type ReportDB interface {
	GetReportInfo(ctx context.Context, in dto.GetReportInfoParams) (dto.GetReportInfoResult, error)
	SetReportStatus(ctx context.Context, in dto.SetReportStatusParams) error
	PickAndMarkDeletingFile(ctx context.Context, in dto.PickAndMarkDeletingFileParams) (dto.PickAndMarkDeletingFileResult, error)
	MarkArchived(ctx context.Context, in dto.MarkArchivedParams) error
	Close()
}
