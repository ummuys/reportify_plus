package web

const (
	CreateReportPath    = "reports"                   // POST
	ListUserReportsPath = "reports"                   // GET
	ReportInfoPath      = "reports/:report_id"        // GET
	ReportStatusPath    = "reports/:report_id/status" // GET
	DeleteUserReports   = "reports"                   // DELETE
	DeleteUserReport    = "reports/:report_id"        // DELETE
)

const (
	ListSchemasPath = "db/schemas"
	ListTablesPath  = "db/tables"
	ListColumnsPath = "db/columns"
)

const (
	ListUsersPath  = "/admin/users"
	CreateUserPath = "/admin/users"
	UpdateUserPath = "/admin/users"

	GetUserInfoPath = "/admin/users/:user_id"
	DeleteUserPath  = "/admin/users/:user_id"
)

const (
	LoginPath        = "secure/login"
	RefreshTokenPath = "secure/refresh"
)
