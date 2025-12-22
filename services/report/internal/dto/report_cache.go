package dto

import "time"

type CacheGetResult struct {
	Values [][]byte
}

type ReportCacheValue struct {
	Name      string
	Comm      string
	CreatedAt time.Time
	Query     string
	CSVSep    string
}
