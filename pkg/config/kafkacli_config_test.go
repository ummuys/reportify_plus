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
			os.Unsetenv("REPORT_KAFKA_GROUP")
			os.Unsetenv("REPORT_KAFKA_PRODUCER_ID")
			os.Unsetenv("REPORT_KAFKA_CONTAINER_NAME")
			os.Unsetenv("REPORT_KAFKA_DOCKER_PORT")
			os.Unsetenv("REPORT_KAFKA_TOPIC_REPORT")
			os.Unsetenv("REPORT_KAFKA_TOPIC_DLQ")

			// Устанавливаем обязательные переменные
			os.Setenv("REPORT_KAFKA_PRODUCER_ID", "test-producer")
			os.Setenv("REPORT_KAFKA_CONTAINER_NAME", "kafka")
			os.Setenv("REPORT_KAFKA_DOCKER_PORT", "9092")
			os.Setenv("REPORT_KAFKA_TOPIC_REPORT", "test-topic")
			os.Setenv("REPORT_KAFKA_TOPIC_DLQ", "test-dlq")

			if tt.setEnv {
				os.Setenv("REPORT_KAFKA_GROUP", tt.envValue)
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