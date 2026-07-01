package dto

type KafkaMessage struct {
	Key   []byte
	Value []byte
}

type KafkaValue struct {
	ReportID    string `json:"report_id"`
	Recreating  bool   `json:"recreating"`
	GraphicMode bool   `json:"graphic_mode"`
}
