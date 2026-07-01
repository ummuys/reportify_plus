package repository

const (
	setSearchPath = `
		SET LOCAL statement_timeout = '90s';
	`

	setStatementTimeout = `
		SET LOCAL search_path = university;
	`
)