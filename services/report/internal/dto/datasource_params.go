package dto

type ListTablesParams struct {
	Schema string
}

type ListColumnsParams struct {
	Schema string
	Table  string
}
