package handlers

import "github.com/gin-gonic/gin"

type ReportServiceHandler interface {
	CreateReport(g *gin.Context)
	RecreateReport(g *gin.Context)
	ListReports(g *gin.Context)
	ReportStatus(g *gin.Context)

	ReportInfo(g *gin.Context)

	DeleteReports(g *gin.Context)
	DeleteReport(g *gin.Context)
}
