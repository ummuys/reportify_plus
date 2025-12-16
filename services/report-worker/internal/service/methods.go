package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report-worker/internal/convert"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
	"github.com/ummuys/reportify/services/report-worker/internal/miniocli"
	"github.com/ummuys/reportify/services/report-worker/internal/repository"
	"golang.org/x/sync/errgroup"
)

type publish struct {
	logger   zerolog.Logger
	dataDB   repository.DataDB
	reportDB repository.ReportDB
	conv     convert.ReportConvert
	minioCli miniocli.MinIOClient
}

type prResMsg struct {
	path string
	err  error
}

func NewPublishService(dataDB repository.DataDB, reportDB repository.ReportDB,
	convert convert.ReportConvert, minioCli miniocli.MinIOClient, baseLogger zerolog.Logger) (PublishService, error) {
	logger := baseLogger.With().Str("component", "svc").Logger()
	return &publish{
		dataDB:   dataDB,
		reportDB: reportDB,
		conv:     convert,
		logger:   logger,
		minioCli: minioCli,
	}, nil
}

func (p *publish) CreateReport(ctx context.Context, in dto.KafkaMessage) error {

	// Change status: created -> runnig
	if err := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
		UUID:         in.UUID,
		UpdateStatus: repository.StatusRunnig,
		BeforeStatus: repository.StatusCreated,
	}); err != nil {
		// Change status: created -> failed
		p.stepFailed(ctx, in.UUID, err, repository.StatusCreated)
		return errs.ParsePgError(err)
	}

	// Get report info
	info, err := p.reportDB.GetReportInfo(ctx, dto.GetReportInfoParams(in))

	if err != nil {
		p.stepFailed(ctx, in.UUID, err, repository.StatusRunnig)
		return errs.ParsePgError(err)
	}

	// Get data from query
	data, err := p.dataDB.GetData(ctx, dto.GetDataParams{
		Query: info.Query,
	})

	if err != nil {
		// Change status: created -> failed
		p.stepFailed(ctx, in.UUID, err, repository.StatusRunnig)
		return errs.ParsePgError(err)
	}

	pr, pw := io.Pipe()

	eg, gctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer pw.Close()
		params := dto.ConvParams{
			Writer: pw,
			Colums: data.Columns,
			Rows:   data.Rows,
			Sep:    info.CSVSep,
		}
		var err error
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
		default:
			err = fmt.Errorf("unsupported format: %s", info.Format)
		}

		if err != nil {
			_ = pw.CloseWithError(err)
		}
		return err
	})

	// Read data from pipe reader
	var path string
	eg.Go(func() error {
		defer pr.Close()
		pth, perr := p.minioCli.UploadAndPresign(gctx, dto.PutReportIn{
			Reader:      pr,
			FileName:    in.UUID,
			Bucket:      "report",
			ContentType: convert.ContentTypeByFormat(info.Format),
			Expire:      time.Hour * 1,
		})

		if perr != nil {
			_ = pr.CloseWithError(perr)
			return perr
		}
		path = pth
		return nil
	})

	if err := eg.Wait(); err != nil {
		p.stepFailed(ctx, in.UUID, err, repository.StatusRunnig)
		return err
	}

	// Change status: running -> completed & save file path
	if err := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
		UUID:         in.UUID,
		UpdateStatus: repository.StatusCompleted,
		BeforeStatus: repository.StatusRunnig,
		FilePath:     &path,
	}); err != nil {
		// Change status: created -> failed
		p.stepFailed(ctx, in.UUID, err, repository.StatusRunnig)
		return errs.ParsePgError(err)
	}

	return nil
}

func (p *publish) stepFailed(ctx context.Context, uuid string, err error, befStat string) {
	errMsg := err.Error()

	fctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if ferr := p.reportDB.SetReportStatus(fctx, dto.SetReportStatusParams{
		UUID:         uuid,
		ErrMsg:       &errMsg,
		UpdateStatus: repository.StatusFailed,
		BeforeStatus: befStat,
	}); ferr != nil {
		p.logger.Error().Err(ferr).Msg("can't change to failed")
	}
}
