package config

import (
	"errors"
	"strings"
)

type GatewayServiceConfig struct {
	AuthServiceAddr string
}

func ParseGatewayConfig() (GatewayServiceConfig, error) {
	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	name, err := parseStr("AUTH_APP_CONTAINER_NAME")
	add(err)

	port, err := parseStr("AUTH_SERVICE_IN_PORT")
	add(err)

	if len(errs) > 0 {
		return GatewayServiceConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return GatewayServiceConfig{
		AuthServiceAddr: name + ":" + port,
	}, nil
}
