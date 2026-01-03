package dto

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
