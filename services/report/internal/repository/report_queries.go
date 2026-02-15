package repository

const (
	createReportQuery = `
	INSERT INTO report_metadata.report_requests 
	(report_id, author_id, name, comment, query_sql, format, csv_separator)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING status;
	`

	recreateReportQuery = `
	UPDATE report_metadata.report_requests
	SET status = 'CREATED',
		file_status = 'NOT_EXISTS',
		updated_at = now()
	WHERE author_id = $1
	AND report_id = $2
	AND file_status IN ('EXISTS', 'DELETED', 'DELETE_FAILED')
	RETURNING report_id;
`

	reportStatusQuery = `
	SELECT status FROM report_metadata.report_requests
	where author_id = $1
	and report_id = $2
	and status != 'ARCHIVED';
	`

	listUserReportsQuery = `
	SELECT 
		report_id, 
		name, 
		comment, 
		query_sql, 
		format, 
		csv_separator, 
		status, 
		created_at,
		updated_at, 
		file_path, 
		error_message
	FROM report_metadata.report_requests
	WHERE author_id = $1
	and status != 'ARCHIVED'
	ORDER BY created_at DESC;
	`

	reportInfoQuery = `
	SELECT  
		name, 
		comment, 
		query_sql, 
		format, 
		csv_separator, 
		status, 
		created_at, 
		file_path, 
		error_message
	FROM report_metadata.report_requests
	WHERE author_id = $1
	and report_id = $2
	and status != 'ARCHIVED';
	`

	deleteUserReportsQuery = `
	UPDATE report_metadata.report_requests
	SET status = 'ARCHIVED'
	WHERE author_id = $1;
	`

	deleteUserReportQuery = `
	UPDATE report_metadata.report_requests
	SET status = 'ARCHIVED'
	WHERE author_id = $1
	and report_id = $2;
	`
)
