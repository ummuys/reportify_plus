package handlers

import (
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(g *gin.Context)
	CreateUser(g *gin.Context)
	UpdateUser(g *gin.Context)
	DeleteUser(g *gin.Context)
	RefreshToken(g *gin.Context)
	ListUsers(g *gin.Context)
}
