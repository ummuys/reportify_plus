package webdto

type CreateReportRequest struct {
	Name   string `json:"name"`
	Comm   string `json:"comment"`
	Query  string `json:"query_sql"`
	Format string `json:"format"`
	CSVSep string `json:"csv_separator,omitempty"`
}

type ListReportsRequest struct {
	AuthorID string `json:"author_id"`
}
