package dto

type KafkaMessage struct {
	ReportID   string `json:"report_id"`
	Recreating bool   `json:"recreating"`
}
