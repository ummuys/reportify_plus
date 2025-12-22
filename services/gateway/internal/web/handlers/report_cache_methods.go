package handlers

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	rcv1 "github.com/ummuys/reportify/api/pb/report/cache/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/gateway/internal/webdto"
)

type reportCacheHandler struct {
	sc     rcv1.ReportCacheServiceClient
	logger zerolog.Logger
}

func NewReportCacheHandler(sc rcv1.ReportCacheServiceClient, baseLogger zerolog.Logger) ReportCacheHandler {
	logger := baseLogger.With().Str("component", "srv").Logger()
	return &reportCacheHandler{sc: sc, logger: logger}
}

func (ch *reportCacheHandler) Get(g *gin.Context) {
	ch.logger.Debug().Str("evt", "call Get").Msg("")

	userID := g.GetString("user_id")

	out, gErr := ch.sc.Get(g.Request.Context(), &rcv1.GetRequest{Key: userID})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}
		g.Set("msg", st.Message())
		g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
		return
	}

	response := append([]byte{'['}, bytes.Join(out.Values, []byte(","))...)
	response = append(response, ']')

	g.Set("msg", "user queries are returned")
	g.Data(http.StatusOK, "application/json; charset=", response)
}

func (ch *reportCacheHandler) Delete(g *gin.Context) {
	ch.logger.Debug().Str("evt", "call Delete").Msg("")

	userID := g.GetString("user_id")

	req := webdto.DeleteQueryRequest{}
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, webdto.BaseResponse{Message: err.Error()})
		return
	}

	_, gErr := ch.sc.Delete(g.Request.Context(), &rcv1.DeleteRequest{
		Key:   userID,
		Value: req.Query,
	})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}
		g.Set("msg", st.Message())
		g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
		return
	}

	g.Set("msg", "query deleted")
	g.JSON(http.StatusOK, webdto.BaseResponse{Message: "query deleted"})
}

func (ch *reportCacheHandler) DeleteAll(g *gin.Context) {
	ch.logger.Debug().Str("evt", "call DeleteAll").Msg("")

	userID := g.GetString("user_id")

	_, gErr := ch.sc.DeleteAll(g.Request.Context(), &rcv1.DeleteAllRequest{Key: userID})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}
		g.Set("msg", st.Message())
		g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
		return
	}

	g.Set("msg", "all user queries deleted")
	g.JSON(http.StatusOK, webdto.BaseResponse{Message: "all user queries deleted"})
}
