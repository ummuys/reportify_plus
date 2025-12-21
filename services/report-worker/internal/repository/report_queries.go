package repository

import (
	"fmt"
	"strings"

	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

const (
	ReportInfoQuery = `
	SELECT author_id, name, comment, query_sql, format, csv_separator FROM report_metadata.report_requests
	WHERE report_id = $1;
	`
)

func buildStatusQuery(in dto.SetReportStatusParams) (string, []any) {
	args := make([]any, 0, 5)
	set := make([]string, 0, 4)

	// updated_at
	set = append(set, "updated_at = NOW()")

	// status
	args = append(args, in.UpdateStatus)
	set = append(set, fmt.Sprintf("status = $%d", len(args)))

	if in.FilePath != nil {
		args = append(args, *in.FilePath)
		set = append(set, fmt.Sprintf("file_path = $%d", len(args)))
	}

	if in.ExpireAt != nil {
		args = append(args, *in.ExpireAt)
		set = append(set, fmt.Sprintf("expire_at = $%d", len(args)))
	}

	if in.ErrMsg != nil {
		args = append(args, *in.ErrMsg)
		set = append(set, fmt.Sprintf("error_message = $%d", len(args)))
	}

	// WHERE
	args = append(args, in.UUID)
	whereReportID := fmt.Sprintf("report_id = $%d", len(args))

	args = append(args, in.BeforeStatus)
	whereBeforeStatus := fmt.Sprintf("status = $%d", len(args))

	q := fmt.Sprintf(`
UPDATE report_metadata.report_requests
SET %s
WHERE %s
  AND %s`,
		strings.Join(set, ",\n    "),
		whereReportID,
		whereBeforeStatus,
	)

	return q, args
}
