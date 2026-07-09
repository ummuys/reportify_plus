package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report-worker/internal/cache"
	"github.com/ummuys/reportify/services/report-worker/internal/convert"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
	"github.com/ummuys/reportify/services/report-worker/internal/miniocli"
	"github.com/ummuys/reportify/services/report-worker/internal/repository"
	"golang.org/x/sync/errgroup"
)

type publish struct {
	logger       zerolog.Logger
	datasourceDB repository.DatasourceDB
	reportDB     repository.ReportDB
	reportCache  cache.ReportCache
	conv         convert.ReportConvert
	minioCli     miniocli.MinIOClient

	reportTTL  time.Duration
	countBatch int
}

func NewPublishService(datasourceDB repository.DatasourceDB, reportDB repository.ReportDB, reportCache cache.ReportCache,
	convert convert.ReportConvert, minioCli miniocli.MinIOClient, reportTTL time.Duration, countBatch int, baseLogger zerolog.Logger,
) (PublishService, error) {
	logger := baseLogger.With().Str("component", "svc").Logger()
	return &publish{
		datasourceDB: datasourceDB,
		reportDB:     reportDB,
		reportCache:  reportCache,
		conv:         convert,
		logger:       logger,
		minioCli:     minioCli,
		reportTTL:    reportTTL,
		countBatch:   countBatch,
	}, nil
}

func (p *publish) CreateReport(ctx context.Context, in dto.KafkaMessage) error {
	p.logger.Debug().Str("evt", "call CreateReport").Str("report_id", in.ReportID).Msg("")

	if err := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
		ReportID:     in.ReportID,
		UpdateStatus: repository.StatusRunning,
		BeforeStatus: repository.StatusCreated,
	}); err != nil {
		p.logger.Error().
			Err(err).
			Str("db-method", "SetReportStatus").
			Str("report_id", in.ReportID).
			Msg("set status running failed")

		p.stepFailed(ctx, in.ReportID, err, repository.StatusCreated)
		return errs.ParsePgError(err)
	}

	info, err := p.reportDB.GetReportInfo(ctx, dto.GetReportInfoParams{ReportID: in.ReportID})
	if err != nil {
		p.logger.Error().
			Err(err).
			Str("db-method", "GetReportInfo").
			Str("report_id", in.ReportID).
			Msg("get report info failed")

		p.stepFailed(ctx, in.ReportID, err, repository.StatusRunning)
		return errs.ParsePgError(err)
	}

	data, err := p.datasourceDB.GetData(ctx, dto.GetDataParams{Query: info.Query})
	if err != nil {
		p.logger.Error().
			Err(err).
			Str("db-method", "GetData").
			Str("report_id", in.ReportID).
			Msg("get datasource data failed")

		p.stepFailed(ctx, in.ReportID, err, repository.StatusRunning)
		return errs.ParsePgError(err)
	}

	pr, pw := io.Pipe()

	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		defer func() { _ = pw.Close() }()

		params := dto.ConvParams{
			Writer:  pw,
			Columns: data.Columns,
			Rows:    data.Rows,
			Sep:     info.CSVSep,
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
			p.logger.Error().
				Err(err).
				Str("op", "convert").
				Str("report_id", in.ReportID).
				Str("format", info.Format).
				Msg("convert failed")

			_ = pw.CloseWithError(err)
		}

		return err
	})

	var path string

	var reportTTL time.Duration
	if in.GraphicMode {
		reportTTL = time.Second * 60
	}
	reportTTL = p.reportTTL

	eg.Go(func() error {
		defer func() { _ = pr.Close() }()

		pth, perr := p.minioCli.UploadAndPresign(gctx, dto.PutReportIn{
			Reader:      pr,
			ObjectName:  in.ReportID,
			FileName:    info.Name + "." + info.Format,
			Bucket:      "report",
			ContentType: convert.ContentTypeByFormat(info.Format),
			Expire:      reportTTL,
		})
		if perr != nil {
			p.logger.Error().
				Err(perr).
				Str("op", "minio.UploadAndPresign").
				Str("report_id", in.ReportID).
				Msg("upload failed")

			_ = pr.CloseWithError(perr)
			return perr
		}

		path = pth
		return nil
	})

	if err := eg.Wait(); err != nil {
		p.stepFailed(ctx, in.ReportID, err, repository.StatusRunning)
		return err
	}

	expireAt := time.Now().Add(reportTTL)

	if err := p.reportDB.SetReportStatus(ctx, dto.SetReportStatusParams{
		ReportID:     in.ReportID,
		UpdateStatus: repository.StatusCompleted,
		BeforeStatus: repository.StatusRunning,
		FilePath:     &path,
		ExpireAt:     &expireAt,
	}); err != nil {
		p.logger.Error().
			Err(err).
			Str("db-method", "SetReportStatus").
			Str("report_id", in.ReportID).
			Msg("set status completed failed")

		p.stepFailed(ctx, in.ReportID, err, repository.StatusRunning)
		return errs.ParsePgError(err)
	}

	if err := p.reportCache.Set(ctx, in.ReportID, repository.StatusCompleted); err != nil {
		p.logger.Warn().
			Err(err).
			Str("op", "cache.Set").
			Str("report_id", in.ReportID).
			Msg("cache set completed failed")
	}

	p.logger.Info().
		Str("report_id", in.ReportID).
		Str("format", info.Format).
		Msg("report completed")

	return nil
}

func (p *publish) RecreateReport(ctx context.Context, in dto.KafkaMessage) error {
	if err := p.minioCli.DeleteFiles(ctx, dto.DeleteExpiredFilesParams{Names: []string{in.ReportID}, Bucket: "report"}); err != nil {
		return err
	}

	if err := p.CreateReport(ctx, in); err != nil {
		return err
	}

	return nil
}

func (p *publish) stepFailed(ctx context.Context, reportID string, err error, befStat string) {
	errMsg := err.Error()

	fctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if ferr := p.reportCache.Set(ctx, reportID, repository.StatusFailed); ferr != nil {
		p.logger.Warn().
			Err(ferr).
			Str("op", "cache.Set").
			Str("report_id", reportID).
			Msg("cache set failed failed")
	}

	if ferr := p.reportDB.SetReportStatus(fctx, dto.SetReportStatusParams{
		ReportID:     reportID,
		ErrMsg:       &errMsg,
		UpdateStatus: repository.StatusFailed,
		BeforeStatus: befStat,
	}); ferr != nil {
		p.logger.Error().
			Err(ferr).
			Str("db-method", "SetReportStatus").
			Str("report_id", reportID).
			Msg("set status failed failed")
	}
}

func (p *publish) CleanOldReports(ctx context.Context) {
	out, err := p.reportDB.PickAndMarkDeletingFile(ctx, dto.PickAndMarkDeletingFileParams{TimeLife: p.reportTTL, CountBatch: p.countBatch})
	if err != nil {
		p.logger.Error().
			Err(err).
			Str("op", "reportDB.PickAndMarkDeletingFile").
			Msg("failed to pick and mark reports")
		return
	}

	derr := p.minioCli.DeleteFiles(ctx, dto.DeleteExpiredFilesParams{
		Names:  out.ReportsId,
		Bucket: "report",
	})
	if derr != nil {
		p.logger.Error().
			Err(derr).
			Str("op", "minioCli.DeleteExpiredFiles").
			Msg("failed to delete reports from MinIO")
	}

	if err = p.reportDB.MarkArchived(ctx, dto.MarkArchivedParams{
		ReportsId: out.ReportsId,
		Error:     derr,
	}); err != nil {
		p.logger.Error().
			Err(err).
			Str("op", "reportDB.MarkArchived").
			Msg("failed to pick reports")
		return
	}
}
