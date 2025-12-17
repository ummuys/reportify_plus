package dto

import "time"

type ListUserReportsParams struct {
	AuthorID string
}

type ListReportsResult struct {
	Reports []ReportMetadata
}

type ReportMetadata struct {
	ReportID  string
	AuthorID  string
	Name      string
	Comm      string
	Query     string
	Format    string
	CSVSep    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	FilePath  string
	ErrMsg    string
}

type ListTablesParams struct {
	Schema string
}
type ListTablesResult struct {
	Tables []Table
}

type Table struct {
	Name    string
	Comment string
}

type ListSchemasResult struct {
	Schemas []Schema
}

type Schema struct {
	Name    string
	Comment string
}

type ListColumnsParams struct {
	Schema string
	Table  string
}
type ListColumnsResult struct {
	Columns []Column
}

type Column struct {
	Name    string
	Comment string
}
