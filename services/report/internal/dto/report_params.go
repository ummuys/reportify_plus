package dto

type CreateReportParams struct {
	AuthorID    string
	Name        string
	Comment     string
	Query       string
	Format      string
	CSVSep      string
	GraphicMode bool
}

type RecreateReportParams struct {
	AuthorID string
	ReportID string
}

type DeleteReportsParams struct {
	AuthorID string
}

type DeleteReportParams struct {
	AuthorID string
	ReportID string
}

type ListReportsParams struct {
	AuthorID string
}

type ReportInfoParams struct {
	AuthorID string
	ReportID string
}

type ReportStatusParams struct {
	AuthorID string
	ReportID string
}
