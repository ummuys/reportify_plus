package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	authv1 "github.com/ummuys/reportify/api/pb/auth/service/v1"
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

// Login godoc
// @Summary Login
// @Description Authenticates user and sets refresh token cookie. Returns access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body webdto.LoginRequest true "payload"
// @Success 200 {object} webdto.LoginResponse "OK"
// @Failure 400 {object} webdto.ErrResponse "Validation error"
// @Failure 404 {object} webdto.ErrResponse "User not found"
// @Failure 500 {object} webdto.ErrResponse "Internal error"
// @Router /secure/login [post]
func (a *authHandler) Login(refreshTime time.Duration) gin.HandlerFunc {
	return func(g *gin.Context) {
		a.logger.Debug().Str("evt", "call Login").Msg("")
		var req webdto.LoginRequest
		if err := g.ShouldBindJSON(&req); err != nil {
			g.Set("msg", err.Error())
			g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: errs.ErrInvalidJSON.Error()})
			return
		}

		req.Username = strings.TrimSpace(req.Username)
		req.Password = strings.TrimSpace(req.Password)

		if req.Username == "" || req.Password == "" {
			g.Set("msg", errs.ErrInvalidPaylod.Error())
			g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: errs.ErrInvalidPaylod.Error()})
			return
		}

		out, gErr := a.sc.Login(g.Request.Context(), &authv1.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		})

		if gErr != nil {
			st, ok := errs.GRPCtoREST(gErr)
			if !ok {
				g.Set("msg", gErr.Error())
				g.AbortWithStatusJSON(http.StatusInternalServerError,
					webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
				return
			}

			g.Set("msg", st.Message())
			var (
				code int
				resp webdto.ErrResponse
			)
			switch st.Code() {
			case codes.Unauthenticated:
				code = http.StatusUnauthorized
				resp.Error = st.Message()
			default:
				code = http.StatusInternalServerError
				resp.Error = errs.ErrServerInternal.Error()
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

// CreateUser godoc
// @Summary Create user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body webdto.CreateUserRequest true "payload"
// @Success 201 {object} webdto.CreateUserResponse
// @Failure 400 {object} webdto.ErrResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 403 {object} webdto.ErrResponse
// @Failure 409 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /admin/users [post]
func (a *authHandler) CreateUser(g *gin.Context) {
	a.logger.Debug().Str("evt", "call CreateUser").Msg("")
	var req webdto.CreateUserRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: errs.ErrInvalidJSON.Error()})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Role = strings.TrimSpace(req.Role)

	if req.Username == "" || req.Password == "" || req.Role == "" {
		g.Set("msg", errs.ErrInvalidPaylod.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: errs.ErrInvalidPaylod.Error()})
		return
	}

	out, gErr := a.sc.CreateUser(g.Request.Context(), &authv1.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	})

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("non-gprc-msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp webdto.ErrResponse
		)
		switch st.Code() {
		case codes.AlreadyExists:
			code = http.StatusConflict
			resp.Error = st.Message()
		case codes.InvalidArgument:
			code = http.StatusBadRequest
			resp.Error = st.Message()
		default:
			code = http.StatusInternalServerError
			resp.Error = errs.ErrServerInternal.Error()
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "user created")
	g.JSON(http.StatusCreated, webdto.CreateUserResponse{UserID: out.UserId})
}

// UpdateUser godoc
// @Summary Update user
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body webdto.UpdateUserRequest true "payload"
// @Success 200 {object} webdto.UpdateUserResponse
// @Failure 400 {object} webdto.ErrResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 403 {object} webdto.ErrResponse
// @Failure 404 {object} webdto.ErrResponse
// @Failure 409 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /admin/users [put]
func (a *authHandler) UpdateUser(g *gin.Context) {
	a.logger.Debug().Str("evt", "call UpdateUser").Msg("")
	var req webdto.UpdateUserRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Set("msg", err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: errs.ErrInvalidJSON.Error()})
		return
	}

	req.UserID = strings.TrimSpace(req.UserID)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Role = strings.TrimSpace(req.Role)

	if req.UserID == "" || req.Username == "" || req.Password == "" || req.Role == "" {
		g.Set("msg", errs.ErrInvalidPaylod.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, webdto.ErrResponse{Error: errs.ErrInvalidPaylod.Error()})
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
				webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp webdto.ErrResponse
		)
		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp.Error = st.Message()
		case codes.AlreadyExists:
			code = http.StatusConflict
			resp.Error = st.Message()
		default:
			code = http.StatusInternalServerError
			resp.Error = errs.ErrServerInternal.Error()
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "user updated")
	g.JSON(http.StatusOK, webdto.UpdateUserResponse{UserID: out.UserId, Username: out.Username, Role: out.Role, IsActive: out.IsActive})
}

// DeleteUser godoc
// @Summary Delete user
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} webdto.DeleteUserResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 403 {object} webdto.ErrResponse
// @Failure 404 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /admin/users/{user_id} [delete]
func (a *authHandler) DeleteUser(g *gin.Context) {
	a.logger.Debug().Str("evt", "call DeleteUser").Msg("")

	userID := g.Param("user_id")

	// TODO: fix validate answer
	// mark at (04.01.2026) deadline at (10.01.2026)
	if userID == "" {
		g.Set("msg", errs.ErrInvalidPaylod.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest,
			webdto.ErrResponse{Error: errs.ErrInvalidPaylod.Error()})
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
				webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp webdto.ErrResponse
		)
		switch st.Code() {
		case codes.NotFound:
			code = http.StatusNotFound
			resp.Error = st.Message()
		case codes.PermissionDenied:
			code = http.StatusForbidden
			resp.Error = st.Message()
		default:
			code = http.StatusInternalServerError
			resp.Error = errs.ErrServerInternal.Error()
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "user deleted")
	g.JSON(http.StatusOK, webdto.DeleteUserResponse{UserID: out.UserId})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Uses refresh_token cookie to issue a new access token.
// @Tags auth
// @Produce json
// @Success 200 {object} webdto.RefreshTokenResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /secure/auth/refresh [get]
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
				webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp webdto.ErrResponse
		)
		switch st.Code() {
		case codes.Unauthenticated:
			code = http.StatusUnauthorized
			resp.Error = st.Message()
		default:
			code = http.StatusInternalServerError
			resp.Error = errs.ErrServerInternal.Error()
		}
		g.AbortWithStatusJSON(code, resp)
		return
	}

	g.Set("msg", "access token refreshed")
	g.JSON(http.StatusOK, webdto.RefreshTokenResponse{AccessToken: out.AccessToken})
}

// ListUsers godoc
// @Summary List users
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} webdto.ListUsersResponse
// @Failure 401 {object} webdto.ErrResponse
// @Failure 403 {object} webdto.ErrResponse
// @Failure 500 {object} webdto.ErrResponse
// @Router /admin/users [get]
func (a *authHandler) ListUsers(g *gin.Context) {
	a.logger.Debug().Str("evt", "call ListUsers").Msg("")

	out, gErr := a.sc.ListUsers(g.Request.Context(), nil)

	if gErr != nil {
		st, ok := errs.GRPCtoREST(gErr)
		if !ok {
			g.Set("msg", gErr.Error())
			g.AbortWithStatusJSON(http.StatusInternalServerError,
				webdto.ErrResponse{Error: errs.ErrServerInternal.Error()})
			return
		}

		g.Set("msg", st.Message())
		var (
			code int
			resp webdto.ErrResponse
		)
		switch st.Code() {
		default:
			code = http.StatusInternalServerError
			resp.Error = errs.ErrServerInternal.Error()
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
