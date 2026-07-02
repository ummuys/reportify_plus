package config

import (
	"errors"
	"strings"
)

type AuthServiceConfig struct {
	AdmUsername string
	AdmPassword string
	Network     string
	Port        string
}

func ParseAuthServiceConfig() (AuthServiceConfig, error) {
	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	un, err := parseStr("AUTH_SERVICE_USERNAME")
	add(err)

	pw, err := parseStr("AUTH_SERVICE_PASSWORD")
	add(err)

	port, err := parseStr("AUTH_SERVICE_IN_PORT")
	add(err)

	network, err := parseStr("AUTH_SERVICE_NETWORK")
	add(err)

	if len(errs) > 0 {
		return AuthServiceConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return AuthServiceConfig{AdmUsername: un, AdmPassword: pw, Network: network, Port: ":" + port}, nil
}
