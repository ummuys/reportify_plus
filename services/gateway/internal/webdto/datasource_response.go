package webdto

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
