package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
	"github.com/ummuys/reportify/services/report-worker/internal/mocks"
	"github.com/ummuys/reportify/services/report-worker/internal/repository"
)

func newPublish(t *testing.T) (PublishService,
	*mocks.MockReportDB,
	*mocks.MockDatasourceDB,
	*mocks.MockReportConvert,
	*mocks.MockMinIOClient,
	*mocks.MockReportCache,
) {
	t.Helper()

	reportDB := mocks.NewMockReportDB(t)
	datasourceDB := mocks.NewMockDatasourceDB(t)
	conv := mocks.NewMockReportConvert(t)
	minioCli := mocks.NewMockMinIOClient(t)
	cache := mocks.NewMockReportCache(t)

	p, _ := NewPublishService(
		datasourceDB,
		reportDB,
		cache,
		conv,
		minioCli,
		time.Hour,
		50,
		zerolog.Nop(),
	)

	return p, reportDB, datasourceDB, conv, minioCli, cache
}

func defaultInfo(format string) dto.GetReportInfoResult {
	return dto.GetReportInfoResult{
		ReportID: "r1",
		Name:     "rep",
		Format:   format,
		CSVSep:   ';',
		Query:    "select 1",
	}
}

func defaultData() dto.GetDataResult {
	return dto.GetDataResult{
		Columns: []string{"c1"},
		Rows:    [][]any{{"v1"}},
	}
}

func expectSetRunningOK(reportDB *mocks.MockReportDB) {
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusRunning &&
				p.BeforeStatus == repository.StatusCreated
		})).
		Return(nil).
		Once()
}

func expectSetRunningErr(reportDB *mocks.MockReportDB, err error) {
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.Anything).
		Return(err).
		Once()
}

func expectGetInfoOK(reportDB *mocks.MockReportDB, info dto.GetReportInfoResult) {
	reportDB.EXPECT().
		GetReportInfo(mock.Anything, mock.Anything).
		Return(info, nil).
		Once()
}

func expectGetInfoErr(reportDB *mocks.MockReportDB, err error) {
	reportDB.EXPECT().
		GetReportInfo(mock.Anything, mock.Anything).
		Return(dto.GetReportInfoResult{}, err).
		Once()
}

func expectGetDataOK(ds *mocks.MockDatasourceDB, query string, data dto.GetDataResult) {
	ds.EXPECT().
		GetData(mock.Anything, dto.GetDataParams{Query: query}).
		Return(data, nil).
		Once()
}

func expectGetDataErr(ds *mocks.MockDatasourceDB, query string, err error) {
	ds.EXPECT().
		GetData(mock.Anything, dto.GetDataParams{Query: query}).
		Return(dto.GetDataResult{}, err).
		Once()
}

func expectConvert(conv *mocks.MockReportConvert, format string, retErr error) {
	switch format {
	case "PDF":
		conv.EXPECT().ToPDF(mock.Anything).Return(retErr).Once()
	case "XLSX":
		conv.EXPECT().ToXLSX(mock.Anything).Return(retErr).Once()
	case "JSON":
		conv.EXPECT().ToJSON(mock.Anything).Return(retErr).Once()
	case "CSV":
		conv.EXPECT().ToCSV(mock.Anything).Return(retErr).Once()
	case "DOCX":
		conv.EXPECT().ToDOCX(mock.Anything).Return(retErr).Once()
	default:
	}
}

func expectUpload(minio *mocks.MockMinIOClient, url string, retErr error) {
	minio.EXPECT().
		UploadAndPresign(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, in dto.PutReportIn) {
			_, _ = io.ReadAll(in.Reader)
		}).
		Return(url, retErr).
		Once()
}

func expectSetCompletedOK(reportDB *mocks.MockReportDB, url string) {
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusCompleted &&
				p.BeforeStatus == repository.StatusRunning &&
				p.FilePath != nil && *p.FilePath == url &&
				p.ExpireAt != nil
		})).
		Return(nil).
		Once()
}

func expectSetCompletedErr(reportDB *mocks.MockReportDB, err error) {
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusCompleted &&
				p.BeforeStatus == repository.StatusRunning
		})).
		Return(err).
		Once()
}

func expectStepFailedOK(t *testing.T, reportDB *mocks.MockReportDB, cache *mocks.MockReportCache, beforeStatus string) {
	t.Helper()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()

	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == beforeStatus &&
				p.ErrMsg != nil && *p.ErrMsg != ""
		})).
		Return(nil).
		Once()
}

func expectStepFailed_CacheSetErr(t *testing.T, reportDB *mocks.MockReportDB, cache *mocks.MockReportCache, beforeStatus string) {
	t.Helper()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(errors.New("cache set failed err")).Once()

	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == beforeStatus &&
				p.ErrMsg != nil && *p.ErrMsg != ""
		})).
		Return(nil).
		Once()
}

func expectStepFailed_DbSetErr(t *testing.T, reportDB *mocks.MockReportDB, cache *mocks.MockReportCache, beforeStatus string) {
	t.Helper()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()

	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == beforeStatus &&
				p.ErrMsg != nil && *p.ErrMsg != ""
		})).
		Return(errors.New("db set status failed failed err")).
		Once()
}

func TestPublish_CreateReport_Success_AllFormats(t *testing.T) {
	cases := []struct {
		name   string
		format string
	}{
		{"PDF", "PDF"},
		{"XLSX", "XLSX"},
		{"JSON", "JSON"},
		{"DOCX", "DOCX"},
		{"CSV", "CSV"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, reportDB, ds, conv, minio, cache := newPublish(t)
			ctx := context.Background()
			in := dto.KafkaMessage{ReportID: "r1"}

			info := defaultInfo(tc.format)
			data := defaultData()

			expectSetRunningOK(reportDB)
			expectGetInfoOK(reportDB, info)
			expectGetDataOK(ds, info.Query, data)
			expectConvert(conv, info.Format, nil)
			expectUpload(minio, "https://signed/url", nil)
			expectSetCompletedOK(reportDB, "https://signed/url")
			cache.EXPECT().Set(mock.Anything, "r1", repository.StatusCompleted).Return(nil).Once()

			err := p.CreateReport(ctx, in)
			require.NoError(t, err)
		})
	}
}

func TestPublish_CreateReport_SetRunningError_ReturnsParsedPgError_AndStepFailed(t *testing.T) {
	p, reportDB, ds, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	expectSetRunningErr(reportDB, dbErr)
	expectStepFailedOK(t, reportDB, cache, repository.StatusCreated)

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)

	ds.AssertNotCalled(t, "GetData", mock.Anything, mock.Anything)
	conv.AssertNotCalled(t, "ToCSV", mock.Anything)
	minio.AssertNotCalled(t, "UploadAndPresign", mock.Anything, mock.Anything)
}

func TestPublish_CreateReport_GetInfoError_ReturnsParsedPgError_AndStepFailed(t *testing.T) {
	p, reportDB, ds, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	expectSetRunningOK(reportDB)
	expectGetInfoErr(reportDB, dbErr)
	expectStepFailedOK(t, reportDB, cache, repository.StatusRunning)

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)

	ds.AssertNotCalled(t, "GetData", mock.Anything, mock.Anything)
	conv.AssertNotCalled(t, "ToCSV", mock.Anything)
	minio.AssertNotCalled(t, "UploadAndPresign", mock.Anything, mock.Anything)
}

func TestPublish_CreateReport_GetDataError_ReturnsParsedPgError_AndStepFailed(t *testing.T) {
	_, reportDB, _, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	info := defaultInfo("CSV")

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	expectSetRunningOK(reportDB)
	expectGetInfoOK(reportDB, info)

	ds := mocks.NewMockDatasourceDB(t)
	p2, _ := NewPublishService(ds, reportDB, cache, conv, minio, time.Hour, 50, zerolog.Nop())

	expectGetDataErr(ds, info.Query, dbErr)
	expectStepFailedOK(t, reportDB, cache, repository.StatusRunning)

	err := p2.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)

	conv.AssertNotCalled(t, "ToCSV", mock.Anything)
	minio.AssertNotCalled(t, "UploadAndPresign", mock.Anything, mock.Anything)
}

func TestPublish_CreateReport_ConvertError_ReturnsError_AndStepFailed(t *testing.T) {
	p, reportDB, ds, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	info := defaultInfo("CSV")
	data := defaultData()

	expectSetRunningOK(reportDB)
	expectGetInfoOK(reportDB, info)
	expectGetDataOK(ds, info.Query, data)

	convErr := errors.New("convert err")
	expectConvert(conv, info.Format, convErr)

	expectUpload(minio, "", nil)

	expectStepFailedOK(t, reportDB, cache, repository.StatusRunning)

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, convErr)
}

func TestPublish_CreateReport_UploadError_ReturnsError_AndStepFailed(t *testing.T) {
	p, reportDB, ds, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	info := defaultInfo("CSV")
	data := defaultData()

	expectSetRunningOK(reportDB)
	expectGetInfoOK(reportDB, info)
	expectGetDataOK(ds, info.Query, data)

	expectConvert(conv, info.Format, nil)

	upErr := errors.New("upload err")
	expectUpload(minio, "", upErr)

	expectStepFailedOK(t, reportDB, cache, repository.StatusRunning)

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, upErr)
}

func TestPublish_CreateReport_SetCompletedError_ReturnsParsedPgError_AndStepFailed(t *testing.T) {
	p, reportDB, ds, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	info := defaultInfo("CSV")
	data := defaultData()

	expectSetRunningOK(reportDB)
	expectGetInfoOK(reportDB, info)
	expectGetDataOK(ds, info.Query, data)

	expectConvert(conv, info.Format, nil)
	expectUpload(minio, "https://signed/url", nil)

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	expectSetCompletedErr(reportDB, dbErr)

	expectStepFailedOK(t, reportDB, cache, repository.StatusRunning)

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
}

func TestPublish_CreateReport_UnsupportedFormat_ReturnsError_AndStepFailed(t *testing.T) {
	p, reportDB, ds, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	info := defaultInfo("AVI")
	data := defaultData()

	expectSetRunningOK(reportDB)
	expectGetInfoOK(reportDB, info)
	expectGetDataOK(ds, info.Query, data)

	expectUpload(minio, "", nil)

	expectStepFailedOK(t, reportDB, cache, repository.StatusRunning)

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "unsupported format"))

	conv.AssertNotCalled(t, "ToPDF", mock.Anything)
	conv.AssertNotCalled(t, "ToXLSX", mock.Anything)
	conv.AssertNotCalled(t, "ToJSON", mock.Anything)
	conv.AssertNotCalled(t, "ToCSV", mock.Anything)
	conv.AssertNotCalled(t, "ToDOCX", mock.Anything)
}

func TestPublish_CreateReport_CacheSetCompletedError_Ignored(t *testing.T) {
	p, reportDB, ds, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	info := defaultInfo("CSV")
	data := defaultData()

	expectSetRunningOK(reportDB)
	expectGetInfoOK(reportDB, info)
	expectGetDataOK(ds, info.Query, data)
	expectConvert(conv, info.Format, nil)
	expectUpload(minio, "https://signed/url", nil)
	expectSetCompletedOK(reportDB, "https://signed/url")

	cacheErr := errors.New("cache set completed err")
	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusCompleted).Return(cacheErr).Once()

	err := p.CreateReport(ctx, in)
	require.NoError(t, err)
}

func TestPublish_CreateReport_StepFailed_CacheSetFailedError_BranchCovered(t *testing.T) {
	p, reportDB, _, _, _, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	expectSetRunningErr(reportDB, dbErr)
	expectStepFailed_CacheSetErr(t, reportDB, cache, repository.StatusCreated)

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
}

func TestPublish_CreateReport_StepFailed_DbSetStatusFailedError_BranchCovered(t *testing.T) {
	p, reportDB, _, _, _, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	expectSetRunningErr(reportDB, dbErr)
	expectStepFailed_DbSetErr(t, reportDB, cache, repository.StatusCreated)

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
}

func TestPublish_CreateReport_ConvertError_IsNotParsed(t *testing.T) {
	p, reportDB, ds, conv, minio, cache := newPublish(t)
	ctx := context.Background()
	in := dto.KafkaMessage{ReportID: "r1"}

	info := defaultInfo("CSV")
	data := defaultData()

	expectSetRunningOK(reportDB)
	expectGetInfoOK(reportDB, info)
	expectGetDataOK(ds, info.Query, data)

	convErr := fmt.Errorf("plain convert err")
	expectConvert(conv, info.Format, convErr)
	expectUpload(minio, "", nil)
	expectStepFailedOK(t, reportDB, cache, repository.StatusRunning)

	err := p.CreateReport(ctx, in)
	require.ErrorIs(t, err, convErr)
}
