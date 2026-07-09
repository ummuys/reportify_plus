package dto

import "time"

type GetReportInfoParams struct {
	ReportID string
}

type GetReportInfoResult struct {
	ReportID string
	AuthorID string
	Name     string
	Comment  string
	Query    string
	Format   string
	CSVSep   byte
}

type SetReportStatusParams struct {
	ReportID     string
	UpdateStatus string
	BeforeStatus string
	FilePath     *string
	ErrMsg       *string
	ExpireAt     *time.Time
}

type PickAndMarkDeletingFileParams struct {
	TimeLife   time.Duration
	CountBatch int
}

type PickAndMarkDeletingFileResult struct {
	ReportsId []string
}

type MarkArchivedParams struct {
	ReportsId []string
	Error     error
}
