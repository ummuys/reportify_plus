package config

import (
	"errors"
	"strings"
)

type GatewayServiceConfig struct {
	ReportServiceAddr string
	AuthServiceAddr   string
	Host              string
	Port              string
}

func ParseGatewayConfig() (GatewayServiceConfig, error) {
	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	authName, err := parseStr("AUTH_SERVICE_CONTAINER_NAME")
	add(err)

	authPort, err := parseStr("AUTH_SERVICE_IN_PORT")
	add(err)

	reportName, err := parseStr("REPORT_SERVICE_CONTAINER_NAME")
	add(err)

	reportPort, err := parseStr("REPORT_SERVICE_IN_PORT")
	add(err)

	restHost, err := parseStr("GATEWAY_SERVICE_HOST")
	add(err)

	restPort, err := parseStr("GATEWAY_SERVICE_IN_PORT")
	add(err)

	if len(errs) > 0 {
		return GatewayServiceConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return GatewayServiceConfig{
		ReportServiceAddr: reportName + ":" + reportPort,
		AuthServiceAddr:   authName + ":" + authPort,
		Host:              restHost,
		Port:              restPort,
	}, nil
}
