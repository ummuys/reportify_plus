package dto

type GetDataParams struct {
	Query string
}

type GetDataResult struct {
	Columns []string
	Rows    [][]any
}
