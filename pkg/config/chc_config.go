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

	name, err := parseStr("REPORT_CACHE_CONTAINER_NAME")
	add(err)

	port, err := parseStr("REPORT_CACHE_IN_PORT")
	add(err)

	pass, err := parseStr("REPORT_CACHE_PASSWORD")
	add(err)

	db, err := parseInt("REPORT_CACHE_DB_IDX", true)
	add(err)

	ttl, err := parseInt("REPORT_CACHE_TTL", true)
	add(err)

	if len(errs) > 0 {
		return ReportCacheConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return ReportCacheConfig{
		Addr:     name + ":" + port,
		Password: pass,
		DB:       db,
		TTL:      time.Duration(int64(ttl)),
	}, nil
}
