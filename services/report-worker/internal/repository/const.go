package repository

const (
	StatusCreated   = "CREATED"
	StatusRunning   = "RUNNING"
	StatusFailed    = "FAILED"
	StatusCompleted = "COMPLETED"
	StatusArchived  = "ARCHIVED"

	StatusFileExists        = "EXISTS"
	StatusFileDeleting      = "DELETING"
	StatusFileDeleted       = "DELETED"
	StatusFileDeletedFailed = "DELETE_FAILED"
)
