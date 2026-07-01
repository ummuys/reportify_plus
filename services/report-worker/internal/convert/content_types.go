package convert

import "strings"

func ContentTypeByFormat(format string) string {
	switch strings.ToUpper(format) {
	case "PDF":
		return "application/pdf"
	case "XLSX":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "JSON":
		return "application/json"
	case "CSV":
		return "text/csv"
	case "DOCX":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}
