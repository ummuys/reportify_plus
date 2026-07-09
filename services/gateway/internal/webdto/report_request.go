package webdto

type CreateReportRequest struct {
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Query       string `json:"query_sql"`
	Format      string `json:"format"`
	CSVSep      string `json:"csv_separator,omitempty"`
	GraphicMode bool   `json:"graphic_mode"`
}

type ListReportsRequest struct {
	AuthorID string `json:"author_id"`
}
