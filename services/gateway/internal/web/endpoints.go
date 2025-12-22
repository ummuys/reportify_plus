package web

const (
	CreateReportPath    = "report"
	ListUserReportsPath = "report"
	ReportStatusPath    = "report/:report_id"
)

const (
	ListSchemasPath = "db/schemas"
	ListTablesPath  = "db/tables"
	ListColumnsPath = "db/columns"
)

const (
	GetCacheQueriesPath       = "cache"
	DeleteAllCacheQueriesPath = "cache/all"
	DeleteCacheQueryPath      = "cache"
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
