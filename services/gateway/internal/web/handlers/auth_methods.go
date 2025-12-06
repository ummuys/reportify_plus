package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	"github.com/ummuys/reportify/services/gateway/internal/webdto"
)

type authHandler struct {
	sc     authv1.AuthServiceClient
	logger zerolog.Logger
}

func NewAuthHandler(sc authv1.AuthServiceClient, baseLogger zerolog.Logger) AuthHandler {
	logger := baseLogger.With().Str("component", "srv").Logger()
	return &authHandler{sc: sc, logger: logger}
}

func (a *authHandler) Login(g *gin.Context) {
	a.logger.Debug().Str("evt", "call Login").Msg("")
	var req webdto.LoginRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	out, err := a.sc.Login(g.Request.Context(), &authv1.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	g.Set("msg", "login succsessful")
	g.JSON(http.StatusOK, webdto.LoginResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})

}

func (a *authHandler) CreateUser(g *gin.Context)   {}
func (a *authHandler) UpdateUser(g *gin.Context)   {}
func (a *authHandler) DeleteUser(g *gin.Context)   {}
func (a *authHandler) RefreshToken(g *gin.Context) {}
func (a *authHandler) ListUsers(g *gin.Context)    {}
