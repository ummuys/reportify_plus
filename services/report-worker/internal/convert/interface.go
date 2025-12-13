package convert

import (
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type ReportConvert interface {
	ToPDF(in dto.ConvParams) error
	ToXLSX(in dto.ConvParams) error
	ToJSON(in dto.ConvParams) error
	ToCSV(in dto.ConvParams) error
	ToDOCX(in dto.ConvParams) error

	// ToMD()

	// Maybe
	// ToTSV()
	// ToYAML()
	// ToParquet()
	// ToArrow()
	// ToSQLDump()

	// Maybe maybe
	// ToSuperSet()
}
