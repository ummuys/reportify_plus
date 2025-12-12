package repository

const (
	CreateReportQuery = `
	INSERT INTO report_metadata.report_requests 
	(report_id, author_id, name, comment, query_sql, format, csv_separator)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING status;
	`
)
