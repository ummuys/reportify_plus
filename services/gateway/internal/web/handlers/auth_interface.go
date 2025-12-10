package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(refreshTime time.Duration) gin.HandlerFunc
	CreateUser(g *gin.Context)
	UpdateUser(g *gin.Context)
	DeleteUser(g *gin.Context)
	RefreshToken(g *gin.Context)
	ListUsers(g *gin.Context)
}
