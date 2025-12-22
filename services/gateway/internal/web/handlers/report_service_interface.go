package handlers

import "github.com/gin-gonic/gin"

type ReportServiceHandler interface {
	CreateReport(g *gin.Context)
	ReportStatus(g *gin.Context)
	ListUserReports(g *gin.Context)
}
