package webdto

type ReportStatusResponse struct {
	UUID     string `json:"uuid"`
	Status   string `json:"status"`
	ErrMsg   string `json:"error_message"`
	FilePath string `json:"file_path"`
}
