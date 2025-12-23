package repository

const (
	createReportQuery = `
	INSERT INTO report_metadata.report_requests 
	(report_id, author_id, name, comment, query_sql, format, csv_separator)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING status;
	`

	getReportStatusQuery = `
	SELECT status, error_message, file_path FROM report_metadata.report_requests
	where report_id = $1;
	`

	listUserReportsQuery = `
	SELECT * FROM report_metadata.report_requests
	WHERE author_id = $1;
	`
)

const (
	getAllReportsQuery = `
		SELECT
		author_id,
		name,
		comment,
		query_sql,
		csv_separator,
		format::text,
		status::text,
		file_path,
		created_at
	FROM report_metadata.report_requests
	ORDER BY created_at ASC;
	`
)

// COLUMNS QUERY
const (
	columnsWithCommentQuery = `
SELECT
  a.attname AS column_name,
  COALESCE(col_description(a.attrelid, a.attnum), '') AS comment
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
JOIN pg_attribute a ON a.attrelid = c.oid
WHERE n.nspname = $1            -- schema
  AND c.relname = $2            -- table
  AND c.relkind IN ('r','p','v','m','f')  -- table, partitioned table, view, matview, foreign table
  AND a.attnum > 0
  AND NOT a.attisdropped;
	`
)

// TABLE QUERY
const (
	tablesWithCommentQuery = `
SELECT
  c.relname AS table_name,
  COALESCE(obj_description(c.oid, 'pg_class'), '') AS comment
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = $1
  AND c.relkind IN ('r','p','v','m','f')  -- отфильтровали лишнее;
	`
)

// SCHEMA QUERY
const (
	schemaWithCommentQuery = `
SELECT
  n.nspname AS schema_name,
  COALESCE(obj_description(n.oid, 'pg_namespace'), '') AS comment
FROM pg_namespace n
WHERE n.nspname NOT IN ('pg_toast', 'pg_catalog', 'information_schema')
  AND n.nspname NOT LIKE 'pg\_%' ESCAPE '\';
	`
)

const (
	setCacheQuery = `
  INSERT INTO identity.queries (user_id, list)
  VALUES ($1, $2)
  ON CONFLICT (user_id) DO UPDATE
  SET list = EXCLUDED.list;`
	getCacheQuery = `SELECT * FROM identity.queries`
)
