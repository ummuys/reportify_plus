package dto

import (
	"io"
	"time"
)

type PutReportIn struct {
	Reader      io.Reader
	FileName    string
	Bucket      string
	ContentType string
	Expire      time.Duration
}
