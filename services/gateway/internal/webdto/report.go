package webdto

import "time"

type CreateReportRequest struct {
	Name   string `json:"name"`
	Comm   string `json:"comment"`
	Query  string `json:"query_sql"`
	Format string `json:"format"`
	CSVSep string `json:"csv_separator,omitempty"`
}

type CreateReportResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type ReportStatusResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type ListUserReportsRequest struct {
	AuthorID string `json:"author_id"`
}

type ListUserReportsResponse struct {
	Reports []ReportMetadata `json:"reports"`
}

type ReportMetadata struct {
	ReportID  string    `json:"report_id"`
	AuthorID  string    `json:"author_id"`
	Name      string    `json:"name"`
	Comm      string    `json:"comm"`
	Query     string    `json:"query"`
	Format    string    `json:"format"`
	CSVSep    string    `json:"csv_sep"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FilePath  string    `json:"file_path"`
	ErrMsg    string    `json:"err_msg"`
}
