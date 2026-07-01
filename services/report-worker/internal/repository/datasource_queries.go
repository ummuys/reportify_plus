package repository

const (
	setStatementTimeout = `
		SET LOCAL statement_timeout = '90s';
	`

	setSearchPath = `
		SET LOCAL search_path = university;
	`
)