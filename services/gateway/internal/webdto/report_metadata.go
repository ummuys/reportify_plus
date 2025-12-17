package webdto

import "time"

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

type ListSchemasResponse struct {
	Schemas []Schema `json:"schemas"`
}
type Schema struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type ListTablesResponse struct {
	Tables []Table `json:"tables"`
}
type Table struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type ListColumnsResponse struct {
	Columns []Column `json:"columns"`
}
type Column struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}
