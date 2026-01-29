package service

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

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
		zerolog.Nop(),
	)

	return p, reportDB, datasourceDB, conv, minioCli, cache
}

func TestPublish_CreateReport_Success_AllFormats(t *testing.T) {
	cases := []struct {
		name       string
		format     string
		expectConv func(conv *mocks.MockReportConvert)
	}{
		{
			name:   "PDF",
			format: "PDF",
			expectConv: func(conv *mocks.MockReportConvert) {
				conv.EXPECT().ToPDF(mock.Anything).Return(nil).Once()
			},
		},
		{
			name:   "XLSX",
			format: "XLSX",
			expectConv: func(conv *mocks.MockReportConvert) {
				conv.EXPECT().ToXLSX(mock.Anything).Return(nil).Once()
			},
		},
		{
			name:   "JSON",
			format: "JSON",
			expectConv: func(conv *mocks.MockReportConvert) {
				conv.EXPECT().ToJSON(mock.Anything).Return(nil).Once()
			},
		},
		{
			name:   "DOCX",
			format: "DOCX",
			expectConv: func(conv *mocks.MockReportConvert) {
				conv.EXPECT().ToDOCX(mock.Anything).Return(nil).Once()
			},
		},
		{
			name:   "CSV",
			format: "CSV",
			expectConv: func(conv *mocks.MockReportConvert) {
				conv.EXPECT().ToCSV(mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
			ctx := context.Background()

			in := dto.KafkaMessage{ReportID: "r1"}

			reportDB.EXPECT().
				SetReportStatus(mock.Anything, mock.Anything).
				Return(nil).
				Once()

			info := dto.GetReportInfoResult{
				ReportID: "r1",
				Name:     "rep",
				Format:   tc.format,
				CSVSep:   ';',
				Query:    "select 1",
			}
			reportDB.EXPECT().
				GetReportInfo(mock.Anything, mock.Anything).
				Return(info, nil).
				Once()

			data := dto.GetDataResult{
				Columns: []string{"c1"},
				Rows:    [][]any{{"v1"}},
			}
			datasourceDB.EXPECT().
				GetData(mock.Anything, dto.GetDataParams{Query: "select 1"}).
				Return(data, nil).
				Once()

			tc.expectConv(conv)

			minioCli.EXPECT().
				UploadAndPresign(mock.Anything, mock.Anything).
				Run(func(ctx context.Context, in dto.PutReportIn) {
					_, _ = io.ReadAll(in.Reader)
				}).
				Return("https://signed/url", nil).
				Once()

			reportDB.EXPECT().
				SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
					return p.ReportID == "r1" &&
						p.UpdateStatus == repository.StatusCompleted &&
						p.BeforeStatus == repository.StatusRunnig &&
						p.FilePath != nil && *p.FilePath == "https://signed/url" &&
						p.ExpireAt != nil
				})).
				Return(nil).
				Once()

			cache.EXPECT().Set(mock.Anything, "r1", repository.StatusCompleted).Return(nil).Once()

			err := p.CreateReport(ctx, in)
			require.NoError(t, err)

		})
	}
}

func TestPublish_CreateReport_SetRunningError_ReturnsParsedPgError_AndStepFailed(t *testing.T) {
	p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
	ctx := context.Background()

	in := dto.KafkaMessage{ReportID: "r1"}

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.Anything).
		Return(dbErr).
		Once()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusCreated &&
				p.ErrMsg != nil
		})).
		Return(nil).
		Once()

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)

	datasourceDB.AssertNotCalled(t, "GetData", mock.Anything, mock.Anything)
	conv.AssertNotCalled(t, "ToCSV", mock.Anything)
	minioCli.AssertNotCalled(t, "UploadAndPresign", mock.Anything, mock.Anything)
}

func TestPublish_CreateReport_GetInfoError_ReturnsParsedPgError_AndStepFailed(t *testing.T) {
	p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
	ctx := context.Background()

	in := dto.KafkaMessage{ReportID: "r1"}

	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.Anything).
		Return(nil).
		Once()

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	reportDB.EXPECT().
		GetReportInfo(mock.Anything, mock.Anything).
		Return(dto.GetReportInfoResult{}, dbErr).
		Once()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusRunnig &&
				p.ErrMsg != nil
		})).
		Return(nil).
		Once()

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)

	datasourceDB.AssertNotCalled(t, "GetData", mock.Anything, mock.Anything)
	conv.AssertNotCalled(t, "ToCSV", mock.Anything)
	minioCli.AssertNotCalled(t, "UploadAndPresign", mock.Anything, mock.Anything)
}

func TestPublish_CreateReport_GetDataError_ReturnsParsedPgError_AndStepFailed(t *testing.T) {
	p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
	ctx := context.Background()

	in := dto.KafkaMessage{ReportID: "r1"}

	reportDB.EXPECT().SetReportStatus(mock.Anything, mock.Anything).Return(nil).Once()

	info := dto.GetReportInfoResult{ReportID: "r1", Name: "rep", Format: "CSV", CSVSep: ';', Query: "select 1"}
	reportDB.EXPECT().GetReportInfo(mock.Anything, mock.Anything).Return(info, nil).Once()

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	datasourceDB.EXPECT().GetData(mock.Anything, dto.GetDataParams{Query: "select 1"}).Return(dto.GetDataResult{}, dbErr).Once()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusRunnig &&
				p.ErrMsg != nil
		})).
		Return(nil).
		Once()

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)

	conv.AssertNotCalled(t, "ToCSV", mock.Anything)
	minioCli.AssertNotCalled(t, "UploadAndPresign", mock.Anything, mock.Anything)
}

func TestPublish_CreateReport_ConvertError_ReturnsError_AndStepFailed(t *testing.T) {
	p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
	ctx := context.Background()

	in := dto.KafkaMessage{ReportID: "r1"}

	reportDB.EXPECT().SetReportStatus(mock.Anything, mock.Anything).Return(nil).Once()

	info := dto.GetReportInfoResult{ReportID: "r1", Name: "rep", Format: "CSV", CSVSep: ';', Query: "select 1"}
	reportDB.EXPECT().GetReportInfo(mock.Anything, mock.Anything).Return(info, nil).Once()

	data := dto.GetDataResult{Columns: []string{"c1"}, Rows: [][]any{{"v1"}}}
	datasourceDB.EXPECT().GetData(mock.Anything, dto.GetDataParams{Query: "select 1"}).Return(data, nil).Once()

	convErr := errors.New("convert err")
	conv.EXPECT().ToCSV(mock.Anything).Return(convErr).Once()

	minioCli.EXPECT().
		UploadAndPresign(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, in dto.PutReportIn) {
			_, _ = io.ReadAll(in.Reader)
		}).
		Return("", nil).
		Maybe()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusRunnig &&
				p.ErrMsg != nil
		})).
		Return(nil).
		Once()

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, convErr)
}

func TestPublish_CreateReport_UploadError_ReturnsError_AndStepFailed(t *testing.T) {
	p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
	ctx := context.Background()

	in := dto.KafkaMessage{ReportID: "r1"}

	reportDB.EXPECT().SetReportStatus(mock.Anything, mock.Anything).Return(nil).Once()

	info := dto.GetReportInfoResult{ReportID: "r1", Name: "rep", Format: "CSV", CSVSep: ';', Query: "select 1"}
	reportDB.EXPECT().GetReportInfo(mock.Anything, mock.Anything).Return(info, nil).Once()

	data := dto.GetDataResult{Columns: []string{"c1"}, Rows: [][]any{{"v1"}}}
	datasourceDB.EXPECT().GetData(mock.Anything, dto.GetDataParams{Query: "select 1"}).Return(data, nil).Once()

	conv.EXPECT().ToCSV(mock.Anything).Return(nil).Once()

	upErr := errors.New("upload err")
	minioCli.EXPECT().
		UploadAndPresign(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, in dto.PutReportIn) {
			_, _ = io.ReadAll(in.Reader)
		}).
		Return("", upErr).
		Once()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusRunnig &&
				p.ErrMsg != nil
		})).
		Return(nil).
		Once()

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, upErr)
}

func TestPublish_CreateReport_SetCompletedError_ReturnsParsedPgError_AndStepFailed(t *testing.T) {
	p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
	ctx := context.Background()

	in := dto.KafkaMessage{ReportID: "r1"}

	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.Anything).
		Return(nil).
		Once()

	info := dto.GetReportInfoResult{ReportID: "r1", Name: "rep", Format: "CSV", CSVSep: ';', Query: "select 1"}
	reportDB.EXPECT().GetReportInfo(mock.Anything, mock.Anything).Return(info, nil).Once()

	data := dto.GetDataResult{Columns: []string{"c1"}, Rows: [][]any{{"v1"}}}
	datasourceDB.EXPECT().GetData(mock.Anything, dto.GetDataParams{Query: "select 1"}).Return(data, nil).Once()

	conv.EXPECT().ToCSV(mock.Anything).Return(nil).Once()

	minioCli.EXPECT().
		UploadAndPresign(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, in dto.PutReportIn) {
			_, _ = io.ReadAll(in.Reader)
		}).
		Return("https://signed/url", nil).
		Once()

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusCompleted &&
				p.BeforeStatus == repository.StatusRunnig
		})).
		Return(dbErr).
		Once()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusRunnig &&
				p.ErrMsg != nil
		})).
		Return(nil).
		Once()

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
}

func TestPublish_CreateReport_UnsupportedFormat_ReturnsError_AndStepFailed(t *testing.T) {
	p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
	ctx := context.Background()

	in := dto.KafkaMessage{ReportID: "r1"}

	reportDB.EXPECT().SetReportStatus(mock.Anything, mock.Anything).Return(nil).Once()

	info := dto.GetReportInfoResult{
		ReportID: "r1",
		Name:     "rep",
		Format:   "AVI",
		CSVSep:   ';',
		Query:    "select 1",
	}
	reportDB.EXPECT().GetReportInfo(mock.Anything, mock.Anything).Return(info, nil).Once()

	data := dto.GetDataResult{Columns: []string{"c1"}, Rows: [][]any{{"v1"}}}
	datasourceDB.EXPECT().GetData(mock.Anything, dto.GetDataParams{Query: "select 1"}).Return(data, nil).Once()

	minioCli.EXPECT().
		UploadAndPresign(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, in dto.PutReportIn) {
			_, _ = io.ReadAll(in.Reader)
		}).
		Return("", nil).
		Maybe()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusRunnig &&
				p.ErrMsg != nil
		})).
		Return(nil).
		Once()

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
	p, reportDB, datasourceDB, conv, minioCli, cache := newPublish(t)
	ctx := context.Background()

	in := dto.KafkaMessage{ReportID: "r1"}

	reportDB.EXPECT().SetReportStatus(mock.Anything, mock.Anything).Return(nil).Once()

	info := dto.GetReportInfoResult{
		ReportID: "r1",
		Name:     "rep",
		Format:   "CSV",
		CSVSep:   ';',
		Query:    "select 1",
	}
	reportDB.EXPECT().GetReportInfo(mock.Anything, mock.Anything).Return(info, nil).Once()

	data := dto.GetDataResult{Columns: []string{"c1"}, Rows: [][]any{{"v1"}}}
	datasourceDB.EXPECT().GetData(mock.Anything, dto.GetDataParams{Query: "select 1"}).Return(data, nil).Once()

	conv.EXPECT().ToCSV(mock.Anything).Return(nil).Once()

	minioCli.EXPECT().
		UploadAndPresign(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, in dto.PutReportIn) { _, _ = io.ReadAll(in.Reader) }).
		Return("https://signed/url", nil).
		Once()

	reportDB.EXPECT().SetReportStatus(mock.Anything, mock.Anything).Return(nil).Once()

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

	reportDB.EXPECT().SetReportStatus(mock.Anything, mock.Anything).Return(dbErr).Once()

	cacheSetErr := errors.New("cache set failed err")
	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(cacheSetErr).Once()

	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusCreated &&
				p.ErrMsg != nil
		})).
		Return(nil).
		Once()

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

	reportDB.EXPECT().SetReportStatus(mock.Anything, mock.Anything).Return(dbErr).Once()

	cache.EXPECT().Set(mock.Anything, "r1", repository.StatusFailed).Return(nil).Once()

	ferr := errors.New("set status failed failed err")
	reportDB.EXPECT().
		SetReportStatus(mock.Anything, mock.MatchedBy(func(p dto.SetReportStatusParams) bool {
			return p.ReportID == "r1" &&
				p.UpdateStatus == repository.StatusFailed &&
				p.BeforeStatus == repository.StatusCreated &&
				p.ErrMsg != nil
		})).
		Return(ferr).
		Once()

	err := p.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
}
