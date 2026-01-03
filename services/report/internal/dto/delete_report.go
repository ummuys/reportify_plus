package dto

type DeleteUserReportsParams struct {
	AuthorID string
}

type DeleteUserReportParams struct {
	AuthorID string
	ReportID string
}

type DeleteUserReportResult struct {
	ReportID string
}
