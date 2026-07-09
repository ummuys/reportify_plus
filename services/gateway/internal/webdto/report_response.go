package webdto

import "time"

type CreateReportResponse struct {
	ReportID string `json:"report_id"`
	Status   string `json:"status"`
}

type RecreateReportResponse struct {
	ReportID string `json:"report_id"`
	Status   string `json:"status"`
}

type DeleteReportResponse struct {
	ReportID string `json:"report_id"`
}

type ListReportsResponse struct {
	Reports []ReportMetadata `json:"reports"`
}

type ReportInfoResponse struct {
	Report ReportMetadata `json:"report"`
}

type ReportMetadata struct {
	ReportID string `json:"report_id"`
	AuthorID string `json:"author_id"`
	Name     string `json:"name"`
	// NOTE: JSON key intentionally kept as "comm" (not "comment") to avoid
	// breaking existing API clients. Field renamed at the Go level only;
	// fix the wire key via a versioned/additive change if needed.
	Comment   string    `json:"comm"`
	Query     string    `json:"query"`
	Format    string    `json:"format"`
	CSVSep    string    `json:"csv_sep"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FilePath  string    `json:"file_path"`
	ErrMsg    string    `json:"err_msg"`
}

type ReportStatusResponse struct {
	ReportID string `json:"report_id"`
	Status   string `json:"status"`
}
