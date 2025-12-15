package dto

type GetReportInfoParams struct {
	UUID string
}

type GetReportInfoResult struct {
	ReportID string
	AuthorID string
	Name     string
	Comm     string
	Query    string
	Format   string
	CSVSep   byte
}

type SetReportStatusParams struct {
	UUID         string
	UpdateStatus string
	BeforeStatus string
}

type SetReportFailedStatusParams struct {
	UUID         string
	Err          string
	BeforeStatus string
}

type FinalizeReportParams struct {
	UUID         string
	UpdateStatus string
	BeforeStatus string
	FilePath     string
}
