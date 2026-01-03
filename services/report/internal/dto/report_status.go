package dto

type ReportStatusParams struct {
	AuthorID string
	ReportID string
}

type ReportStatusResult struct {
	ReportID string
	Status   string
}
