package dto

import "time"

type CacheGetResult struct {
	Values [][]byte
}

type ReportCacheValue struct {
	Name      string
	Comm      string
	Query     string
	CSVSep    string
	Format    string
	Status    string
	FilePath  string
	CreatedAt time.Time
}
