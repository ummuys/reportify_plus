package config

import (
	"errors"
	"strings"
	"time"
)

type ReportCacheConfig struct {
	Addr     string
	Password string
	DB       int
	TTL      time.Duration
}

func ParseReportCacheEnv() (ReportCacheConfig, error) {
	var errs []string

	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	addr, err := parseStr("REPORT_CACHE_ADDR")
	add(err)

	pass, err := parseStr("REPORT_CACHE_PASSWORD")
	add(err)

	db, err := parseInt("REPORT_CACHE_DB", true)
	add(err)

	ttl, err := parseInt("REPORT_CACHE_TTL", true)
	add(err)

	if len(errs) > 0 {
		return ReportCacheConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return ReportCacheConfig{
		Addr:     addr,
		Password: pass,
		DB:       db,
		TTL:      time.Duration(int64(ttl)),
	}, nil
}
