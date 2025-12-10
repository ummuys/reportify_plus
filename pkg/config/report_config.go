package config

import (
	"errors"
	"strings"
)

type ReportServiceConfig struct {
	Network string
	Port    string
}

func ParseReportServiceConfig() (ReportServiceConfig, error) {

	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	port, err := parseStr("REPORT_SERVICE_IN_PORT")
	add(err)

	network, err := parseStr("REPORT_SERVICE_NETWORK")
	add(err)

	if len(errs) > 0 {
		return ReportServiceConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return ReportServiceConfig{Network: network, Port: ":" + port}, nil
}
