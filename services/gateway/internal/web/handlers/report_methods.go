package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	reportv1 "github.com/ummuys/reportify/api/pb/report/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/gateway/internal/webdto"
)

type reportHandler struct {
	sc     reportv1.ReportServiceClient
	logger zerolog.Logger
}

func NewReportHandler(sc reportv1.ReportServiceClient, baseLogger zerolog.Logger) ReportHandler {
	logger := baseLogger.With().Str("component", "srv").Logger()
	return &reportHandler{sc: sc, logger: logger}
}

func (r *reportHandler) CreateReport(g *gin.Context) {
	r.logger.Debug().Str("evt", "call CreateReport").Msg("")
	req := webdto.CreateReportRequest{}
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	out, gErr := r.sc.CreateReport(g.Request.Context(), &reportv1.CreateReportRequest{
		AuthorId: req.AuthorID,
		Name:     req.Name,
		Comm:     req.Comm,
		Query:    req.Comm,
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

func (rh *reportHandler) ReportStatus(g *gin.Context) {
	return
}
