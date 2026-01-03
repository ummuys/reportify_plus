package dto

import "time"

type ListUserReportsParams struct {
	AuthorID string
}

type ListReportsResult struct {
	Reports []ReportMetadata
}

type ReportMetadata struct {
	ReportID  string
	Name      string
	Comm      string
	Query     string
	Format    string
	CSVSep    string
	Status    string
	CreatedAt time.Time
	FilePath  string
	ErrMsg    string
}

type ReportInfoParams struct {
	AuthorID string
	ReportID string
}

type ReportInfoResult struct {
	Report ReportMetadata
}
