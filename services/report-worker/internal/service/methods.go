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
		UpdateStatus: repository.StatusRunnig,
		BeforeStatus: repository.StatusCreated,
	}); err != nil {
		if ferr := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
			UUID:         in.UUID,
			UpdateStatus: repository.StatusFailed,
			BeforeStatus: repository.StatusCreated,
		}); ferr != nil {
			p.logger.Error().Err(ferr).Msg("can't change to failed")
		}
		return errs.ParsePgError(err)
	}

	info, err := p.reportDB.GetReportInfo(ctx, dto.GetReportInfoParams{
		UUID: in.UUID,
	})

	if err != nil {
		if ferr := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
			UUID:         in.UUID,
			UpdateStatus: repository.StatusFailed,
			BeforeStatus: repository.StatusRunnig,
		}); ferr != nil {
			p.logger.Error().Err(ferr).Msg("can't change to failed")
		}
		return errs.ParsePgError(err)
	}

	data, err := p.dataDB.GetData(ctx, dto.GetDataParams{
		Query: info.Query,
	})

	if err != nil {
		if ferr := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
			UUID:         in.UUID,
			UpdateStatus: repository.StatusFailed,
			BeforeStatus: repository.StatusRunnig,
		}); ferr != nil {
			p.logger.Error().Err(ferr).Msg("can't change to failed")
		}
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
		if ferr := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
			UUID:         in.UUID,
			UpdateStatus: repository.StatusFailed,
			BeforeStatus: repository.StatusRunnig,
		}); ferr != nil {
			p.logger.Error().Err(ferr).Msg("can't change to failed")
		}
		return err
	}

	p.minioCli.LoadFile()

	p.reportDB.FinalizeReport(ctx, dto.FinalizeReportParams{
		UUID:         in.UUID,
		UpdateStatus: repository.StatusCompleted,
		BeforeStatus: repository.StatusRunnig,
		FilePath:     "this must be a file path :)",
	})

	return nil
}
