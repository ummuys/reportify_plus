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
	ReportID string
	Status   string
}
