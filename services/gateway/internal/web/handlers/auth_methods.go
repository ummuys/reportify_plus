package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/gateway/internal/webdto"
	"google.golang.org/grpc/codes"
)

type authHandler struct {
	sc     authv1.AuthServiceClient
	logger zerolog.Logger
}

func NewAuthHandler(sc authv1.AuthServiceClient, baseLogger zerolog.Logger) AuthHandler {
	logger := baseLogger.With().Str("component", "srv").Logger()
	return &authHandler{sc: sc, logger: logger}
}

func (a *authHandler) Login(refreshTime time.Duration) gin.HandlerFunc {
	return func(g *gin.Context) {
		a.logger.Debug().Str("evt", "call Login").Msg("")
		var req webdto.LoginRequest
		if err := g.ShouldBindJSON(&req); err != nil {
			g.Set("msg", err.Error())
			g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: errs.ErrInvalidJSON.Error()})
			return
		}

		out, gErr := a.sc.Login(g.Request.Context(), &authv1.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		})

		if gErr != nil {
			st, ok := errs.GRPCtoREST(gErr)
			if !ok {
				g.Set("non-gprc-msg", gErr.Error())
				g.AbortWithStatusJSON(http.StatusInternalServerError,
					webdto.ErrResponse{Error: errs.ErrInternal.Error()})
				return
			}

			g.Set("msg", st.Message())
			var (
				code int
				resp any
			)
			switch st.Code() {
			case codes.Unauthenticated:
				code = http.StatusNotFound
				resp = webdto.ErrResponse{Error: st.Message()}
			default:
				code = http.StatusInternalServerError
				resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
			}
			g.AbortWithStatusJSON(code, resp)
			return
		}

		g.Set("msg", "login succsessful")
		g.SetCookie("refresh_token", out.RefreshToken, int(refreshTime), "/", "", false, true)
		g.JSON(http.StatusOK, webdto.LoginResponse{
			AccessToken: out.AccessToken,
		})
	}
}

func (a *authHandler) CreateUser(g *gin.Context) {
	a.logger.Debug().Str("evt", "call CreateUser").Msg("")
	var req webdto.CreateUserRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	out, gErr := a.sc.CreateUser(g.Request.Context(), &authv1.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Password,
	})

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("non-gprc-msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp any
		)
		switch st.Code() {
		case codes.AlreadyExists:
			code = http.StatusConflict
			resp = webdto.ErrResponse{Error: st.Message()}
		default:
			code = http.StatusInternalServerError
			resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
		}
		g.AbortWithStatusJSON(code, resp)
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

	out, gErr := a.sc.UpdateUser(g.Request.Context(), &authv1.UpdateUserRequest{
		UserId:   req.UserID,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	})

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp any
		)
		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp = webdto.ErrResponse{Error: st.Message()}
		case codes.AlreadyExists:
			code = http.StatusConflict
			resp = webdto.ErrResponse{Error: st.Message()}
		default:
			code = http.StatusInternalServerError
			resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "user updated")
	g.JSON(http.StatusOK, webdto.UpdateUserResponse{UserID: out.UserId, Username: out.Username, Role: out.Role, IsActive: out.IsActive})

}
func (a *authHandler) DeleteUser(g *gin.Context) {
	a.logger.Debug().Str("evt", "call DeleteUser").Msg("")

	userID, err := strconv.ParseInt(g.Param("user_id"), 10, 64)
	if err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusInternalServerError,
			webdto.ErrResponse{Error: errs.ErrInternal.Error()})
		return
	}
	if userID < 0 {
		g.Set("msg", "invalid user_id")
		g.AbortWithStatusJSON(http.StatusInternalServerError,
			webdto.ErrResponse{Error: errs.ErrInternal.Error()})
		return
	}
	out, gErr := a.sc.DeleteUser(g.Request.Context(), &authv1.DeleteUserRequest{
		UserId: userID,
	})

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp any
		)
		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp = webdto.ErrResponse{Error: st.Message()}
		case codes.PermissionDenied:
			code = http.StatusForbidden
			resp = webdto.ErrResponse{Error: st.Message()}
		default:
			code = http.StatusInternalServerError
			resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "user deleted")
	g.JSON(http.StatusOK, webdto.DeleteUserResponse{UserID: out.UserId})
}
func (a *authHandler) RefreshToken(g *gin.Context) {
	a.logger.Debug().Str("evt", "call RefreshToken").Msg("")

	refreshToken, err := g.Cookie("refresh_token")
	if err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusUnauthorized, webdto.ErrResponse{Error: errs.ErrBadRefreshToken.Error()})
		return
	}

	out, gErr := a.sc.RefreshToken(g.Request.Context(), &authv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp any
		)
		switch st.Code() {
		case codes.Unauthenticated:
			code = http.StatusUnauthorized
			resp = webdto.ErrResponse{Error: st.Message()}
		default:
			code = http.StatusInternalServerError
			resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "access token refreshed")
	g.JSON(http.StatusOK, webdto.RefreshTokenResponse{AccessToken: out.AccessToken})

}

func (a *authHandler) ListUsers(g *gin.Context) {
	a.logger.Debug().Str("evt", "call ListUsers").Msg("")

	out, gErr := a.sc.ListUsers(g.Request.Context(), nil)

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp any
		)
		switch st.Code() {
		default:
			code = http.StatusInternalServerError
			resp = webdto.ErrResponse{Error: errs.ErrInternal.Error()}
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	var resp webdto.ListUsersResponse
	resp.Users = make([]webdto.User, 0, len(out.Users))
	for _, user := range out.Users {
		resp.Users = append(resp.Users, webdto.User{
			UserID:   user.UserId,
			Username: user.Username,
			Role:     user.Role,
		})
	}
	g.Set("msg", "users list returned")
	g.JSON(http.StatusOK, resp)
}
