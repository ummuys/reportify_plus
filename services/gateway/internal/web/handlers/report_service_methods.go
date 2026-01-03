package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	reportservicev1 "github.com/ummuys/reportify/api/pb/report/service/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/gateway/internal/webdto"
	"google.golang.org/grpc/codes"
)

type reportServiceHandler struct {
	sc     reportservicev1.ReportServiceClient
	logger zerolog.Logger
}

func NewReportServiceHandler(sc reportservicev1.ReportServiceClient, baseLogger zerolog.Logger) ReportServiceHandler {
	logger := baseLogger.With().Str("component", "srv").Logger()
	return &reportServiceHandler{sc: sc, logger: logger}
}

func (rsh *reportServiceHandler) CreateReport(g *gin.Context) {
	rsh.logger.Debug().Str("evt", "call CreateReport").Msg("")

	userID := g.GetString("user_id")

	req := webdto.CreateReportRequest{}
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	out, gErr := rsh.sc.CreateReport(g.Request.Context(), &reportservicev1.CreateReportRequest{
		AuthorId: userID,
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
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())

		code := http.StatusInternalServerError
		resp := any(webdto.ErrResponse{Error: errs.ErrInternal.Error()})

		switch st.Code() {
		default:
		}

		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "task to create report created")
	g.JSON(http.StatusCreated, webdto.CreateReportResponse{
		ReportID: out.ReportId,
		Status:   out.Status,
	})
}

func (rsh *reportServiceHandler) ListUserReports(g *gin.Context) {
	rsh.logger.Debug().Str("evt", "call ListUserReports").Msg("")

	userID := g.GetString("user_id")

	out, gErr := rsh.sc.ListUserReports(g.Request.Context(), &reportservicev1.ListUserReportsRequest{
		AuthorId: userID,
	})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())

		code := http.StatusInternalServerError
		resp := any(webdto.ErrResponse{Error: errs.ErrInternal.Error()})

		switch st.Code() {
		default:
		}

		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "list of reports returned")

	resp := webdto.ListUserReportsResponse{
		Reports: make([]webdto.ReportMetadata, 0, len(out.Reports)),
	}

	for _, pbReport := range out.Reports {
		metadata := webdto.ReportMetadata{
			ReportID: pbReport.ReportId,
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

		resp.Reports = append(resp.Reports, metadata)
	}

	g.JSON(http.StatusOK, resp)
}

func (rsh *reportServiceHandler) ReportStatus(g *gin.Context) {
	rsh.logger.Debug().Str("evt", "call ReportStatus").Msg("")

	reportID := g.Param("report_id")
	userID := g.GetString("user_id")

	out, gErr := rsh.sc.ReportStatus(g.Request.Context(), &reportservicev1.ReportStatusRequest{
		AuthorId: userID,
		ReportId: reportID,
	})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())

		code := http.StatusInternalServerError
		resp := any(webdto.ErrResponse{Error: errs.ErrInternal.Error()})

		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp = webdto.BaseResponse{Message: "can't find report by report_id"}
		default:
		}

		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "report status returned")
	g.JSON(http.StatusOK, webdto.ReportStatusResponse{
		ReportID: out.ReportId,
		Status:   out.Status,
	})
}

func (rsh *reportServiceHandler) ReportInfo(g *gin.Context) {
	rsh.logger.Debug().Str("evt", "call ReportInfo").Msg("")

	reportID := g.Param("report_id")
	userID := g.GetString("user_id")

	out, gErr := rsh.sc.ReportInfo(g.Request.Context(), &reportservicev1.ReportInfoRequest{
		AuthorId: userID,
		ReportId: reportID,
	})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())

		code := http.StatusInternalServerError
		resp := any(webdto.ErrResponse{Error: errs.ErrInternal.Error()})

		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp = webdto.BaseResponse{Message: "can't find report by report_id"}
		default:
		}

		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "report info returned")

	pbReport := out.Report

	metadata := webdto.ReportMetadata{
		ReportID: pbReport.ReportId,
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

	g.JSON(http.StatusOK, webdto.ReportInfoResponse{
		Report: metadata,
	})
}

func (rsh *reportServiceHandler) DeleteUserReports(g *gin.Context) {
	rsh.logger.Debug().Str("evt", "call DeleteUserReports").Msg("")

	userID := g.GetString("user_id")

	_, gErr := rsh.sc.DeleteUserReports(g.Request.Context(), &reportservicev1.DeleteUserReportsRequest{
		AuthorId: userID,
	})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())

		code := http.StatusInternalServerError
		resp := any(webdto.ErrResponse{Error: errs.ErrInternal.Error()})

		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp = webdto.BaseResponse{Message: "can't find reports for user"}
		default:
		}

		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "user reports deleted")
	g.JSON(http.StatusOK, webdto.BaseResponse{Message: "ok"})
}

func (rsh *reportServiceHandler) DeleteUserReport(g *gin.Context) {
	rsh.logger.Debug().Str("evt", "call DeleteUserReport").Msg("")

	reportID := g.Param("report_id")
	userID := g.GetString("user_id")

	out, gErr := rsh.sc.DeleteUserReport(g.Request.Context(), &reportservicev1.DeleteUserReportRequest{
		AuthorId: userID,
		ReportId: reportID,
	})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())

		code := http.StatusInternalServerError
		resp := any(webdto.ErrResponse{Error: errs.ErrInternal.Error()})

		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp = webdto.BaseResponse{Message: "can't find report by report_id"}
		default:
		}

		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "user report deleted")
	g.JSON(http.StatusOK, webdto.DeleteUserReportResponse{
		ReportID: out.ReportId,
	})
}
