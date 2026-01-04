package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	datasourceservicev1 "github.com/ummuys/reportify/api/pb/datasource/service/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/gateway/internal/webdto"
)

type datasourceHandler struct {
	sc     datasourceservicev1.DatasourceServiceClient
	logger zerolog.Logger
}

func NewDatasourceHandler(sc datasourceservicev1.DatasourceServiceClient, baseLogger zerolog.Logger) DatasourceHandler {
	logger := baseLogger.With().Str("component", "srv").Logger()
	return &datasourceHandler{sc: sc, logger: logger}
}

// ListSchemas godoc
// @Summary List schemas
// @Description Returns list of available schemas
// @Tags datasource
// @Produce json
// @Security BearerAuth
// @Success 200 {object} webdto.ListSchemasResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /schemas [get]
func (dh *datasourceHandler) ListSchemas(g *gin.Context) {
	dh.logger.Debug().Str("evt", "call ListSchemas").Msg("")

	out, gErr := dh.sc.ListSchemas(g.Request.Context(), &datasourceservicev1.ListSchemasRequest{})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
		return
	}

	resp := webdto.ListSchemasResponse{
		Schemas: make([]webdto.Schema, 0, len(out.Schemas)),
	}

	for _, s := range out.Schemas {
		resp.Schemas = append(resp.Schemas, webdto.Schema{
			Name:    s.SchemaName,
			Comment: s.SchemaComm,
		})
	}

	g.Set("msg", "list of schemas returned")
	g.JSON(http.StatusOK, resp)
}

// ListTables godoc
// @Summary List tables in schema
// @Description Returns list of tables for given schema
// @Tags datasource
// @Produce json
// @Security BearerAuth
// @Param schema query string true "Schema name"
// @Success 200 {object} webdto.ListTablesResponse
// @Failure 400 {object} webdto.ErrResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /tables [get]
func (dh *datasourceHandler) ListTables(g *gin.Context) {
	dh.logger.Debug().Str("evt", "call ListTables").Msg("")

	schema := g.Query("schema")
	if schema == "" {
		g.Set("msg", "schema is required")
		g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: "schema is required"})
		return
	}

	out, gErr := dh.sc.ListTables(g.Request.Context(), &datasourceservicev1.ListTablesRequest{
		SchemaName: schema,
	})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
		return
	}

	resp := webdto.ListTablesResponse{
		Tables: make([]webdto.Table, 0, len(out.Tables)),
	}

	for _, t := range out.Tables {
		resp.Tables = append(resp.Tables, webdto.Table{
			Name:    t.TableName,
			Comment: t.TableComm,
		})
	}

	g.Set("msg", "list of tables returned")
	g.JSON(http.StatusOK, resp)
}

// ListTables godoc
// @Summary List tables in schema
// @Description Returns list of tables for given schema
// @Tags datasource
// @Produce json
// @Security BearerAuth
// @Param schema query string true "Schema name"
// @Success 200 {object} webdto.ListTablesResponse
// @Failure 400 {object} webdto.ErrResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /tables [get]
func (dh *datasourceHandler) ListColumns(g *gin.Context) {
	dh.logger.Debug().Str("evt", "call ListColumns").Msg("")

	schema := g.Query("schema")
	table := g.Query("table")
	if schema == "" || table == "" {
		g.Set("msg", "schema and table are required")
		g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: "schema and table are required"})
		return
	}

	out, gErr := dh.sc.ListColumns(g.Request.Context(), &datasourceservicev1.ListColumnsRequest{
		SchemaName: schema,
		TableName:  table,
	})
	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
		return
	}

	resp := webdto.ListColumnsResponse{
		Columns: make([]webdto.Column, 0, len(out.Columns)),
	}

	for _, c := range out.Columns {
		resp.Columns = append(resp.Columns, webdto.Column{
			Name:    c.ColumnName,
			Comment: c.ColumnComm,
		})
	}

	g.Set("msg", "list of columns returned")
	g.JSON(http.StatusOK, resp)
}
