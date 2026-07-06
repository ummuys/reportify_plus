package config

import (
	"os"
	"testing"
)

func TestParseReportCreateConsumerConfig_GroupIsRequired(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		setEnv    bool
		wantError bool
	}{
		{
			name:      "group is missing",
			setEnv:    false,
			wantError: true,
		},
		{
			name:      "group is empty",
			envValue:  "",
			setEnv:    true,
			wantError: true,
		},
		{
			name:      "group is set",
			envValue:  "report-worker-group",
			setEnv:    true,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем все переменные окружения
			_ = os.Unsetenv("REPORT_KAFKA_GROUP")
			_ = os.Unsetenv("REPORT_KAFKA_PRODUCER_ID")
			_ = os.Unsetenv("REPORT_KAFKA_CONTAINER_NAME")
			_ = os.Unsetenv("REPORT_KAFKA_DOCKER_PORT")
			_ = os.Unsetenv("REPORT_KAFKA_TOPIC_REPORT")
			_ = os.Unsetenv("REPORT_KAFKA_TOPIC_DLQ")

			// Устанавливаем обязательные переменные
			_ = os.Setenv("REPORT_KAFKA_PRODUCER_ID", "test-producer")
			_ = os.Setenv("REPORT_KAFKA_CONTAINER_NAME", "kafka")
			_ = os.Setenv("REPORT_KAFKA_DOCKER_PORT", "9092")
			_ = os.Setenv("REPORT_KAFKA_TOPIC_REPORT", "test-topic")
			_ = os.Setenv("REPORT_KAFKA_TOPIC_DLQ", "test-dlq")

			if tt.setEnv {
				_ = os.Setenv("REPORT_KAFKA_GROUP", tt.envValue)
			}

			_, err := ParseReportCreateConsumerConfig()

			if tt.wantError && err == nil {
				t.Errorf("expected error when %s, got nil", tt.name)
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error when %s: %v", tt.name, err)
			}
		})
	}
}
