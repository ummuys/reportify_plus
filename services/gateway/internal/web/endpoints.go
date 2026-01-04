package web

const (
	CreateReportPath = "reports"                   // POST
	ListReportsPath  = "reports"                   // GET
	ReportInfoPath   = "reports/:report_id"        // GET
	ReportStatusPath = "reports/:report_id/status" // GET
	DeleteReports    = "reports"                   // DELETE
	DeleteReport     = "reports/:report_id"        // DELETE
)

const (
	ListSchemasPath = "db/schemas"
	ListTablesPath  = "db/tables"
	ListColumnsPath = "db/columns"
)

const (
	ListsPath  = "/admin/users"
	CreatePath = "/admin/users"
	UpdatePath = "/admin/users"

	GetInfoPath = "/admin/users/:user_id"
	DeletePath  = "/admin/users/:user_id"
)

const (
	LoginPath        = "secure/login"
	RefreshTokenPath = "secure/refresh"
)
