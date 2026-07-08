package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
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

// CreateReport godoc
// @Summary Create report
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body webdto.CreateReportRequest true "payload"
// @Success 201 {object} webdto.CreateReportResponse
// @Failure 400 {object} webdto.ErrResponse
// @Router /reports [post]
func (rsh *reportServiceHandler) CreateReport(g *gin.Context) {
    rsh.logger.Debug().Str("evt", "call CreateReport").Msg("")

    userID := g.GetString("user_id")
    if _, err := uuid.Parse(userID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidUserID.Error()})
        return
    }

    req := webdto.CreateReportRequest{}
    if err := g.ShouldBindJSON(&req); err != nil {
        g.Set("msg", err.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: err.Error()})
        return
    }

    out, gErr := rsh.sc.CreateReport(g.Request.Context(), &reportservicev1.CreateReportRequest{
        AuthorId:    userID,
        Name:        req.Name,
        Comm:        req.Comm,
        Query:       req.Query,
        Format:      req.Format,
        CsvSep:      req.CSVSep,
        GraphicMode: req.GraphicMode,
    })

    if gErr != nil {
        st, ok := errs.GRPCtoREST(gErr)
        if !ok {
            g.Set("msg", gErr.Error())
            g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
            return
        }

        g.Set("msg", st.Message())

        code := http.StatusInternalServerError
        resp := webdto.ErrResponse{Error: errs.ErrServerInternal.Error()}

        // switch st.Code() {
        // default:
        // }

        g.AbortWithStatusJSON(code, resp)
        return
    }

    g.Set("msg", "task to create report created")
    g.JSON(http.StatusCreated, webdto.CreateReportResponse{
        ReportID: out.ReportId,
        Status:   out.Status,
    })
}

func (rsh *reportServiceHandler) RecreateReport(g *gin.Context) {
    rsh.logger.Debug().Str("evt", "call RecreateReport").Msg("")

    reportID := g.Param("report_id")
    if _, err := uuid.Parse(reportID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidReportID.Error()})
        return
    }

    userID := g.GetString("user_id")
    if _, err := uuid.Parse(userID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidUserID.Error()})
        return
    }

    out, gErr := rsh.sc.RecreateReport(g.Request.Context(), &reportservicev1.RecreateReportRequest{
        AuthorId: userID,
        ReportId: reportID,
    })

    if gErr != nil {
        st, ok := errs.GRPCtoREST(gErr)
        if !ok {
            g.Set("msg", gErr.Error())
            g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
            return
        }

        g.Set("msg", st.Message())

        code := http.StatusInternalServerError
        resp := webdto.ErrResponse{Error: errs.ErrServerInternal.Error()}

        // switch st.Code() {
        // default:
        // }

        g.AbortWithStatusJSON(code, resp)
        return
    }

    g.Set("msg", "report recreation started")
    g.JSON(http.StatusCreated, webdto.RecreateReportResponse{
        ReportID: out.ReportId,
        Status:   out.Status,
    })
}

// ListReports godoc
// @Summary List user reports
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} webdto.ListReportsResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /reports [get]
func (rsh *reportServiceHandler) ListReports(g *gin.Context) {
    rsh.logger.Debug().Str("evt", "call ListReports").Msg("")

    userID := g.GetString("user_id")
    if _, err := uuid.Parse(userID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidUserID.Error()})
        return
    }

    out, gErr := rsh.sc.ListReports(g.Request.Context(), &reportservicev1.ListReportsRequest{
        AuthorId: userID,
    })

    if gErr != nil {
        st, ok := errs.GRPCtoREST(gErr)
        if !ok {
            g.Set("msg", gErr.Error())
            g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
            return
        }

        g.Set("msg", st.Message())

        code := http.StatusInternalServerError
        resp := webdto.ErrResponse{Error: errs.ErrServerInternal.Error()}

        // switch st.Code() {
        // default:
        // }

        g.AbortWithStatusJSON(code, resp)
        return
    }

    g.Set("msg", "list of reports returned")

    resp := webdto.ListReportsResponse{
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

        if pbReport.UpdatedAt != nil {
            metadata.UpdatedAt = pbReport.UpdatedAt.AsTime()
        }

        resp.Reports = append(resp.Reports, metadata)
    }

    g.JSON(http.StatusOK, resp)
}

// ReportStatus godoc
// @Summary Get report status
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param report_id path string true "Report ID"
// @Success 200 {object} webdto.ReportStatusResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 404 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /reports/{report_id}/status [get]
func (rsh *reportServiceHandler) ReportStatus(g *gin.Context) {
    rsh.logger.Debug().Str("evt", "call ReportStatus").Msg("")

    reportID := g.Param("report_id")
    if _, err := uuid.Parse(reportID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidReportID.Error()})
        return
    }

    userID := g.GetString("user_id")
    if _, err := uuid.Parse(userID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidUserID.Error()})
        return
    }

    out, gErr := rsh.sc.ReportStatus(g.Request.Context(), &reportservicev1.ReportStatusRequest{
        AuthorId: userID,
        ReportId: reportID,
    })
    if gErr != nil {
        st, ok := errs.GRPCtoREST(gErr)
        if !ok {
            g.Set("msg", gErr.Error())
            g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
            return
        }

        g.Set("msg", st.Message())

        code := http.StatusInternalServerError
        var resp webdto.ErrResponse

        switch st.Code() {
        case codes.NotFound:
            code = http.StatusNotFound
            resp.Error = "can't find report by report_id"
        default:
            resp.Error = errs.ErrServerInternal.Error()
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

// ReportInfo godoc
// @Summary Get report info
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param report_id path string true "Report ID"
// @Success 200 {object} webdto.ReportInfoResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 404 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /reports/{report_id} [get]
func (rsh *reportServiceHandler) ReportInfo(g *gin.Context) {
    rsh.logger.Debug().Str("evt", "call ReportInfo").Msg("")

    reportID := g.Param("report_id")
    if _, err := uuid.Parse(reportID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidReportID.Error()})
        return
    }

    userID := g.GetString("user_id")
    if _, err := uuid.Parse(userID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidUserID.Error()})
        return
    }

    out, gErr := rsh.sc.ReportInfo(g.Request.Context(), &reportservicev1.ReportInfoRequest{
        AuthorId: userID,
        ReportId: reportID,
    })
    if gErr != nil {
        st, ok := errs.GRPCtoREST(gErr)
        if !ok {
            g.Set("msg", gErr.Error())
            g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
            return
        }

        g.Set("msg", st.Message())

        code := http.StatusInternalServerError
        var resp webdto.ErrResponse

        switch st.Code() {
        case codes.NotFound:
            code = http.StatusNotFound
            resp.Error = "can't find report by report_id"
        default:
            resp.Error = errs.ErrServerInternal.Error()
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

// DeleteReports godoc
// @Summary Delete all user reports
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} webdto.BaseResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 404 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /reports [delete]
func (rsh *reportServiceHandler) DeleteReports(g *gin.Context) {
    rsh.logger.Debug().Str("evt", "call DeleteReports").Msg("")

    userID := g.GetString("user_id")
    if _, err := uuid.Parse(userID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidUserID.Error()})
        return
    }

    _, gErr := rsh.sc.DeleteReports(g.Request.Context(), &reportservicev1.DeleteReportsRequest{
        AuthorId: userID,
    })
    if gErr != nil {
        st, ok := errs.GRPCtoREST(gErr)
        if !ok {
            g.Set("msg", gErr.Error())
            g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
            return
        }

        g.Set("msg", st.Message())

        code := http.StatusInternalServerError
        var resp webdto.ErrResponse

        switch st.Code() {
        case codes.NotFound:
            code = http.StatusNotFound
            resp.Error = "can't find reports by report_id"
        default:
            resp.Error = errs.ErrServerInternal.Error()
        }

        g.AbortWithStatusJSON(code, resp)
        return
    }

    g.Set("msg", "reports deleted")
    g.JSON(http.StatusOK, webdto.BaseResponse{Message: "ok"})
}

// DeleteReport godoc
// @Summary Delete one report
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param report_id path string true "Report ID"
// @Success 200 {object} webdto.DeleteReportResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 404 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /reports/{report_id} [delete]
func (rsh *reportServiceHandler) DeleteReport(g *gin.Context) {
    rsh.logger.Debug().Str("evt", "call DeleteReport").Msg("")

    reportID := g.Param("report_id")
    if _, err := uuid.Parse(reportID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidReportID.Error()})
        return
    }

    ID := g.GetString("user_id")
    if _, err := uuid.Parse(ID); err != nil {
        g.Set("msg", errs.ErrInvalidPaylod.Error())
        g.AbortWithStatusJSON(http.StatusBadRequest,
            webdto.ErrResponse{Error: errs.ErrInvalidUserID.Error()})
        return
    }

    out, gErr := rsh.sc.DeleteReport(g.Request.Context(), &reportservicev1.DeleteReportRequest{
        AuthorId: ID,
        ReportId: reportID,
    })
    if gErr != nil {
        st, ok := errs.GRPCtoREST(gErr)
        if !ok {
            g.Set("msg", gErr.Error())
            g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
            return
        }

        g.Set("msg", st.Message())

        code := http.StatusInternalServerError
        var resp webdto.ErrResponse

        switch st.Code() {
        case codes.NotFound:
            code = http.StatusNotFound
            resp.Error = "can't find reports by report_id"
        default:
            resp.Error = errs.ErrServerInternal.Error()
        }

        g.AbortWithStatusJSON(code, resp)
        return
    }

    g.Set("msg", " report deleted")
    g.JSON(http.StatusOK, webdto.DeleteReportResponse{
        ReportID: out.ReportId,
    })
}
