package dto

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
