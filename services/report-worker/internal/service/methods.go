package service

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report-worker/internal/convert"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
	"github.com/ummuys/reportify/services/report-worker/internal/miniocli"
	"github.com/ummuys/reportify/services/report-worker/internal/repository"
)

type publish struct {
	logger   zerolog.Logger
	dataDB   repository.DataDB
	reportDB repository.ReportDB
	conv     convert.ReportConvert
	minioCli miniocli.MinIOClient
}

func NewPublishService(dataDB repository.DataDB, reportDB repository.ReportDB, convert convert.ReportConvert, baseLogger zerolog.Logger) (PublishService, error) {
	logger := baseLogger.With().Str("component", "svc").Logger()
	return &publish{
		dataDB:   dataDB,
		reportDB: reportDB,
		conv:     convert,
		logger:   logger,
	}, nil
}

func (p *publish) CreateReport(ctx context.Context, in dto.KafkaMessage) error {
	if err := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
		UUID:         in.UUID,
		UpdateStatus: "RUNNING",
		BeforeStatus: "CREATED",
	}); err != nil {
		return errs.ParsePgError(err)
	}

	info, err := p.reportDB.GetReportInfo(ctx, dto.GetReportInfoParams{
		UUID: in.UUID,
	})

	if err != nil {
		return errs.ParsePgError(err)
	}

	data, err := p.dataDB.GetData(ctx, dto.GetDataParams{
		Query: info.Query,
	})

	if err != nil {
		return errs.ParsePgError(err)
	}

	var file *os.File
	params := dto.ConvParams{
		Colums: data.Columns,
		Rows:   data.Rows,
		File:   file,
		Sep:    info.CSVSep,
	}
	switch info.Format {
	case "PDF":
		err = p.conv.ToPDF(params)
	case "XLSX":
		err = p.conv.ToXLSX(params)
	case "JSON":
		err = p.conv.ToJSON(params)
	case "CSV":
		err = p.conv.ToCSV(params)
	case "DOCX":
		err = p.conv.ToDOCX(params)
	}

	if err != nil {
		return err
	}

	p.minioCli.LoadFile()

	p.reportDB.FinalizeReport(ctx, dto.FinalizeReportParams{
		UUID:         in.UUID,
		UpdateStatus: "COMPLETED",
		BeforeStatus: "RUNNIG",
		FilePath:     "this must be a file path :)",
	})

	return nil
}
