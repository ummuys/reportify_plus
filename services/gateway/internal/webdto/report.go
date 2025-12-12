package webdto

type CreateReportRequest struct {
	AuthorID string `json:"author_id"`
	Name     string `json:"name"`
	Comm     string `json:"comment"`
	Query    string `json:"query_sql"`
	Format   string `json:"format"`
	CSVSep   string `json:"csv_separator,omitempty"`
}

type CreateReportResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type ReportStatusRequest struct {
	UUID string `json:"uuid"`
}

type ReportStatusResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}
