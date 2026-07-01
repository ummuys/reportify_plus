#!/bin/sh

set -eu

BOOTSTRAP="${KAFKA_BOOTSTRAP:-report-kafka:29092}"

TOPIC_REPORT="${REPORT_KAFKA_TOPIC_REPORT:?REPORT_KAFKA_TOPIC_REPORT is not set}"
TOPIC_DLQ="${REPORT_KAFKA_TOPIC_DLQ:?REPORT_KAFKA_TOPIC_DLQ is not set}"
PARTITIONS_REPORT="${PARTITIONS_REPORT:-12}"
PARTITIONS_DLQ="${PARTITIONS_DLQ:-3}"
RF="${RF:-1}"

echo "Creating topics..."

/opt/kafka/bin/kafka-topics.sh --bootstrap-server "$BOOTSTRAP" \
  --create --if-not-exists --topic "$TOPIC_DLQ" \
  --partitions "$PARTITIONS_DLQ" --replication-factor "$RF" \
  --config cleanup.policy=delete \
  --config retention.ms=1209600000 \
  --config min.insync.replicas=1

/opt/kafka/bin/kafka-topics.sh --bootstrap-server "$BOOTSTRAP" \
  --create --if-not-exists --topic "$TOPIC_REPORT" \
  --partitions "$PARTITIONS_REPORT" --replication-factor "$RF"

echo "Topics ready"