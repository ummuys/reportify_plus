package service

import (
	"context"
	"io"
	"time"

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
	pwCh := make(chan error, 1)
	prCh := make(chan prResMsg, 1)

	// Write data to pipe writter
	go func() {
		defer pw.Close()
		defer close(pwCh)
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
		}

		if err != nil {
			_ = pw.CloseWithError(err)
		}
		pwCh <- err
	}()

	// Read data from pipe reader
	go func() {
		defer pr.Close()
		defer close(prCh)
		var (
			path string
			err  error
		)
		path, err = p.minioCli.UploadAndPresign(ctx, dto.PutReportIn{
			Reader:      pr,
			FileName:    info.Name,
			Bucket:      "report",
			ContentType: convert.ContentTypeByFormat(info.Format),
			Expire:      time.Hour * 1,
		})

		if err != nil {
			_ = pr.CloseWithError(err)
		}

		prCh <- prResMsg{
			path: path,
			err:  err,
		}
	}()

	res := <-prCh
	if res.err != nil {
		// Change status: created -> failed
		p.stepFailed(ctx, in.UUID, res.err, repository.StatusRunnig)
		return res.err
	}

	if err := <-pwCh; err != nil {
		// Change status: created -> failed
		p.stepFailed(ctx, in.UUID, err, repository.StatusRunnig)
		return err
	}

	// Change status: running -> completed & save file path
	if err := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
		UUID:         in.UUID,
		UpdateStatus: repository.StatusCompleted,
		BeforeStatus: repository.StatusRunnig,
		FilePath:     &res.path,
	}); err != nil {
		// Change status: created -> failed
		p.stepFailed(ctx, in.UUID, err, repository.StatusRunnig)
		return err
	}

	return nil
}

func (p *publish) stepFailed(ctx context.Context, uuid string, err error, befStat string) {
	errMsg := err.Error()
	if ferr := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
		UUID:         uuid,
		ErrMsg:       &errMsg,
		BeforeStatus: befStat,
	}); ferr != nil {
		p.logger.Error().Err(ferr).Msg("can't change to failed")
	}
}
