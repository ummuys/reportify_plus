package config

import (
	"errors"
	"strings"
)

type MinIOCliConfig struct {
	EndPoint        string
	AccessKeyID     string
	SecretAccessKey string
}

func ParseMinIOCliConfig() (MinIOCliConfig, error) {
	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	name, err := parseStr("REPORT_MINIO_CONTAINER_NAME")
	add(err)

	port, err := parseStr("REPORT_MINIO_IN_API_PORT")
	add(err)

	user, err := parseStr("REPORT_MINIO_USER")
	add(err)

	pass, err := parseStr("REPORT_MINIO_PASSWORD")
	add(err)

	if len(errs) > 0 {
		return MinIOCliConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return MinIOCliConfig{
		EndPoint:        name + ":" + port,
		AccessKeyID:     user,
		SecretAccessKey: pass,
	}, nil
}
