package handlers

import "github.com/gin-gonic/gin"

type ReportCacheHandler interface {
	Get(g *gin.Context)
	Delete(g *gin.Context)
	DeleteAll(g *gin.Context)
}
