package config

import (
	"errors"
	"strings"
	"time"
)

type ReportWorkerConfig struct {
	FileTTL     time.Duration
	DeleteBatch int
}

func ParseReportWorkerConfig() (ReportWorkerConfig, error) {
	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	ttl, err := parseInt("FILE_TTL", false)
	add(err)

	db, err := parseInt("DELETE_BATCH", false)
	add(err)

	if len(errs) > 0 {
		return ReportWorkerConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return ReportWorkerConfig{FileTTL: time.Second * time.Duration(ttl), DeleteBatch: db}, nil
}
