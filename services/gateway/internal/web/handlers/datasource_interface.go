package handlers

import "github.com/gin-gonic/gin"

type DatasourceHandler interface {
	ListSchemas(g *gin.Context)
	ListTables(g *gin.Context)
	ListColumns(g *gin.Context)
}
