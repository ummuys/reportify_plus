package dto

type ReportStatusParams struct {
	UUID string
}

type ReportStatusResult struct {
	UUID   string
	Status string
}
