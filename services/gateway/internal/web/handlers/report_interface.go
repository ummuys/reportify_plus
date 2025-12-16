package handlers

import "github.com/gin-gonic/gin"

type ReportHandler interface {
	CreateReport(g *gin.Context)
	ReportStatus(g *gin.Context)
	ListUserReports(g *gin.Context)
	ListSchemas(g *gin.Context)
	ListTables(g *gin.Context)
	ListColumns(g *gin.Context)
}
