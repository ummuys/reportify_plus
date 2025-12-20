package dto

import (
	"io"
	"time"
)

type PutReportIn struct {
	Reader      io.Reader
	ObjectName  string
	FileName    string
	Bucket      string
	ContentType string
	Expire      time.Duration
}
