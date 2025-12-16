package web

// BASE PATH

// REPORT
const (
	CreateReportPath    = "report"
	ListUserReportsPath = "report"
	ReportStatusPath    = "report/:report_id"
)

// Metadata
const (
	ListSchemasPath = "db/schemas"
	ListTablesPath  = "db/tables"
	ListColumnsPath = "db/columns"
)

// ADMIN
const (

	// Collection
	ListUsersPath  = "/admin/users" // GET
	CreateUserPath = "/admin/users" // POST
	UpdateUserPath = "/admin/users" // PATCH

	// Item
	GetUserInfoPath = "/admin/users/:user_id" // GET
	DeleteUserPath  = "/admin/users/:user_id" // DELETE
)

const (
	LoginPath        = "secure/login"
	RefreshTokenPath = "secure/refresh"
)
