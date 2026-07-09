package dto

import "time"

type CreateReportResult struct {
	ReportID string
	Status   string
}

type RecreateReportResult struct {
	ReportID string
	Status   string
}

type DeleteReportResult struct {
	ReportID string
}

type ListReportsResult struct {
	Reports []ReportMetadata
}

type ReportMetadata struct {
	ReportID  string
	Name      string
	Comment   string
	Query     string
	Format    string
	CSVSep    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	FilePath  string
	ErrMsg    string
}

type ReportInfoResult struct {
	Report ReportMetadata
}

type ReportStatusResult struct {
	ReportID string
	Status   string
}
