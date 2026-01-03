package webdto

type CreateReportRequest struct {
	Name   string `json:"name"`
	Comm   string `json:"comment"`
	Query  string `json:"query_sql"`
	Format string `json:"format"`
	CSVSep string `json:"csv_separator,omitempty"`
}

type CreateReportResponse struct {
	ReportID string `json:"report_id"`
	Status   string `json:"status"`
}
