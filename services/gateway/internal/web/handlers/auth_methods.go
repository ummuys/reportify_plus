package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	"github.com/ummuys/reportify/pkg/errs"
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
		g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: errs.ErrInvalidJSON.Error()})
		return
	}

	out, err := a.sc.Login(g.Request.Context(), &authv1.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		g.Set("msg", err.Error())
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			g.AbortWithStatusJSON(http.StatusUnauthorized, webdto.ErrResponse{Error: err.Error()})
		default:
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
		}
		return
	}

	g.Set("msg", "login succsessful")
	g.JSON(http.StatusOK, webdto.LoginResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	})

}

func (a *authHandler) CreateUser(g *gin.Context) {
	a.logger.Debug().Str("evt", "call CreateUser").Msg("")
	var req webdto.CreateUserRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	out, err := a.sc.CreateUser(g.Request.Context(), &authv1.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Password,
	})

	if err != nil {
		g.Set("msg", err.Error())
		switch {
		case errors.Is(err, errs.ErrDuplicate):
			g.AbortWithStatusJSON(http.StatusUnauthorized, webdto.ErrResponse{Error: errs.ErrUserExists.Error()})
		default:
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
		}
		return
	}

	g.Set("msg", "user created")
	g.JSON(http.StatusCreated, webdto.CreateUserResponse{UserID: out.UserId})
}
func (a *authHandler) UpdateUser(g *gin.Context) {
	a.logger.Debug().Str("evt", "call UpdateUser").Msg("")
	var req webdto.UpdateUserRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	out, err := a.sc.UpdateUser(g.Request.Context(), &authv1.UpdateUserRequest{
		UserId:   req.UserID,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Password,
	})

	if err != nil {
		g.Set("msg", err.Error())
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			g.AbortWithStatusJSON(http.StatusUnauthorized, webdto.ErrResponse{Error: err.Error()})
		default:
			g.AbortWithStatusJSON(http.StatusInternalServerError, webdto.ErrResponse{Error: errs.ErrInternal.Error()})
		}
		return
	}

	g.Set("msg", "user updated")
	g.JSON(http.StatusOK, webdto.UpdateUserResponse{UserID: out.UserId, Username: out.Username, Role: out.Role, IsActive: out.IsActive})

}
func (a *authHandler) DeleteUser(g *gin.Context)   {}
func (a *authHandler) RefreshToken(g *gin.Context) {}
func (a *authHandler) ListUsers(g *gin.Context)    {}
