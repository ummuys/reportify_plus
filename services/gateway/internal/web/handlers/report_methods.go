package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	reportv1 "github.com/ummuys/reportify/api/pb/report/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/gateway/internal/webdto"
	"google.golang.org/grpc/codes"
)

type reportHandler struct {
	sc     reportv1.ReportServiceClient
	logger zerolog.Logger
}

func NewReportHandler(sc reportv1.ReportServiceClient, baseLogger zerolog.Logger) ReportHandler {
	logger := baseLogger.With().Str("component", "srv").Logger()
	return &reportHandler{sc: sc, logger: logger}
}

func (rh *reportHandler) CreateReport(g *gin.Context) {
	rh.logger.Debug().Str("evt", "call CreateReport").Msg("")

	id := g.GetString("user_id")

	req := webdto.CreateReportRequest{}
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	out, gErr := rh.sc.CreateReport(g.Request.Context(), &reportv1.CreateReportRequest{
		AuthorId: id,
		Name:     req.Name,
		Comm:     req.Comm,
		Query:    req.Query,
		Format:   req.Format,
		CsvSep:   req.CSVSep,
	})

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp any
		)
		switch st.Code() {
		default:
			code = http.StatusInternalServerError
			resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "task to create report created")
	g.JSON(http.StatusCreated, webdto.CreateReportResponse{
		UUID:   out.Uuid,
		Status: out.Status,
	})
}

func (rh *reportHandler) ListUserReports(g *gin.Context) {
	rh.logger.Debug().Str("evt", "call ListUserReports").Msg("")

	id := g.GetString("user_id")

	out, gErr := rh.sc.ListUserReports(g.Request.Context(), &reportv1.ListUserReportsRequest{
		AuthorId: id,
	})

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp any
		)
		switch st.Code() {
		default:
			code = http.StatusInternalServerError
			resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "list of reports returned")
	resp := webdto.ListUserReportsResponse{}
	resp.Reports = make([]webdto.ReportMetadata, 0, len(out.Reports))

	for _, pbReport := range out.Reports {
		metadata := webdto.ReportMetadata{
			ReportID: pbReport.ReportId,
			AuthorID: pbReport.AuthorId,
			Name:     pbReport.Name,
			Comm:     pbReport.Comm,
			Query:    pbReport.Query,
			Format:   pbReport.Format,
			CSVSep:   pbReport.CsvSep,
			Status:   pbReport.Status,
			FilePath: pbReport.FilePath,
			ErrMsg:   pbReport.ErrMsg,
		}

		if pbReport.CreatedAt != nil {
			metadata.CreatedAt = pbReport.CreatedAt.AsTime()
		}
		if pbReport.UpdatedAt != nil {
			metadata.UpdatedAt = pbReport.UpdatedAt.AsTime()
		}

		resp.Reports = append(resp.Reports, metadata)
	}

	g.JSON(http.StatusOK, resp)

}

func (rh *reportHandler) ReportStatus(g *gin.Context) {
	rh.logger.Debug().Str("evt", "call ReportStatus").Msg("")

	id := g.Param("report_id")

	out, gErr := rh.sc.ReportStatus(g.Request.Context(), &reportv1.ReportStatusRequest{
		Uuid: id,
	})

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp any
		)
		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp = webdto.BaseResponse{Message: "can't find report by uuid"}
		default:
			code = http.StatusInternalServerError
			resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "report status returned")
	g.JSON(http.StatusOK, webdto.ReportStatusResponse{
		UUID:   out.Uuid,
		Status: out.Status,
	})
}
