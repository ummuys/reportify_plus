package dto

import "time"

type CreateReportParams struct {
	AuthorID string
	Name     string
	Comm     string
	Query    string
	Format   string
	CSVSep   string
}

type CreateReportResult struct {
	UUID   string
	Status string
}

type ReportStatusParams struct {
	UUID string
}

type ReportStatusResult struct {
	UUID   string
	Status string
}

type ListUserReportsParams struct {
	AuthorID string
}

type ListReportsResult struct {
	Reports []ReportMetadata
}

type ReportMetadata struct {
	ReportID  string
	AuthorID  string
	Name      string
	Comm      string
	Query     string
	Format    string
	CSVSep    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	FilePath  string
	ErrMsg    string
}
