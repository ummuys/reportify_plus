package config

import (
	"errors"
	"strings"
)

type GatewayServiceConfig struct {
	AuthServiceAddr string
	Host            string
	Port            string
}

func ParseGatewayConfig() (GatewayServiceConfig, error) {
	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	name, err := parseStr("AUTH_SERVICE_CONTAINER_NAME")
	add(err)

	port, err := parseStr("AUTH_SERVICE_IN_PORT")
	add(err)

	restHost, err := parseStr("GATEWAY_SERVICE_HOST")
	add(err)

	restPort, err := parseStr("GATEWAY_SERVICE_IN_PORT")
	add(err)

	if len(errs) > 0 {
		return GatewayServiceConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return GatewayServiceConfig{
		AuthServiceAddr: name + ":" + port,
		Host:            restHost,
		Port:            restPort,
	}, nil
}
