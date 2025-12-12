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
	GetSchemasPath       = "db/schemas"
	GetTablesPath        = "db/tables"
	GetColumnsPath       = "db/columns"
	GetAllQueriesPath    = "cache"
	DeleteAllQueriesPath = "cache/all"
	DeleteQueryPath      = "cache"
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
