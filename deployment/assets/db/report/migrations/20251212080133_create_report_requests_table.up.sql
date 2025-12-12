CREATE TYPE report_status_type AS ENUM ('CREATED', 'RUNNING', 'COMPLETED', 'FAILED');
CREATE TYPE report_output_format AS ENUM ('CSV', 'JSON', 'PDF', 'XLSX', 'DOCX');

CREATE TABLE IF NOT EXISTS report_metadata.report_requests (
    report_id UUID PRIMARY KEY,
    author_id UUID NOT NULL,
    name TEXT NOT NULL,
    comment TEXT,
    query_sql TEXT NOT NULL,
    format report_output_format NOT NULL,
    csv_separator VARCHAR(1) NOT NULL DEFAULT ',',
    status report_status_type NOT NULL DEFAULT 'CREATED', 
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    file_path TEXT,
    error_message TEXT
);