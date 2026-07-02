package config

import (
	"errors"
	"strings"
)

type KafkaProducerConfig struct {
	ClientID string
	Broker   string
	Topic    string
	TopicDLQ string
}

func ParseKafkaProducerConfig() (KafkaProducerConfig, error) {
	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	cid, err := parseStr("REPORT_KAFKA_PRODUCER_ID")
	add(err)

	name, err := parseStr("REPORT_KAFKA_CONTAINER_NAME")
	add(err)

	port, err := parseStr("REPORT_KAFKA_DOCKER_PORT")
	add(err)

	t, err := parseStr("REPORT_KAFKA_TOPIC_REPORT")
	add(err)

	tDLQ, err := parseStr("REPORT_KAFKA_TOPIC_DLQ")
	add(err)

	if len(errs) > 0 {
		return KafkaProducerConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return KafkaProducerConfig{
		ClientID: cid,
		Broker:   name + ":" + port,
		Topic:    t,
		TopicDLQ: tDLQ,
	}, nil
}

type ReportCreateConsumerConfig struct {
	ClientID string
	Group    string
	Broker   string
	Topic    string
	TopicDLQ string
}

func ParseReportCreateConsumerConfig() (ReportCreateConsumerConfig, error) {
	var errs []string
	add := func(err error) {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	cid, err := parseStr("REPORT_KAFKA_PRODUCER_ID")
	add(err)

	name, err := parseStr("REPORT_KAFKA_CONTAINER_NAME")
	add(err)

	port, err := parseStr("REPORT_KAFKA_DOCKER_PORT")
	add(err)

	group, err := parseStr("REPORT_KAFKA_GROUP")
	add(err)
	t, err := parseStr("REPORT_KAFKA_TOPIC_REPORT")
	add(err)

	tDLQ, err := parseStr("REPORT_KAFKA_TOPIC_DLQ")
	add(err)

	if len(errs) > 0 {
		return ReportCreateConsumerConfig{}, errors.New(strings.Join(errs, ", "))
	}

	return ReportCreateConsumerConfig{
		ClientID: cid,
		Broker:   name + ":" + port,
		Group:    group,
		Topic:    t,
		TopicDLQ: tDLQ,
	}, nil
}
