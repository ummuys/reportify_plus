package dto

import "time"

type GetReportInfoParams struct {
	UUID string
}

type GetReportInfoResult struct {
	ReportID string
	AuthorID string
	Name     string
	Comm     string
	Query    string
	Format   string
	CSVSep   byte
}

type SetReportStatusParams struct {
	UUID         string
	UpdateStatus string
	BeforeStatus string
	FilePath     *string
	ErrMsg       *string
	ExpireAt     *time.Time
}
