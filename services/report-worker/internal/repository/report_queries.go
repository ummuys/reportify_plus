package repository

const (
	ReportInfoQuery = `
	SELECT author_id, name, comment, query_sql, format, csv_separator FROM report_metadata.report_requests
	WHERE report_id = $1;
	`

	SetReportStatusQuery = `
	UPDATE report_metadata.report_requests
	SET 
		status = $1,
		updated_at = NOW()
	WHERE report_id = $2
		and status = $3
	`

	SetReportFailedStatusQuery = `
	UPDATE report_metadata.report_requests
	SET 
		status = 'FAILED',
		updated_at = NOW(),
		error_message = $1
	WHERE report_id = $2
		and status = $3
	`

	FinalizeReportQuery = `
	UPDATE report_metadata.report_requests
	SET 
		status = $1,
		file_path = $2,
		updated_at = NOW()
	WHERE report_id = $3
		and status = $4
	`
)
