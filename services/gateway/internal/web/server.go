package web

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
	pkg "github.com/ummuys/reportify/pkg/tm"
	"github.com/ummuys/reportify/services/gateway/internal/di"
	"github.com/ummuys/reportify/services/gateway/internal/web/middleware"
)

func CreateServer(cfg config.GatewayServiceConfig, rh di.RESTHandlers, tm pkg.TokenManager, baseLogger zerolog.Logger) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()

	logger := baseLogger.With().Str("component", "srv").Logger()

	api := g.Group("/api/v1")
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8008"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "UPDATE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api.Use(middleware.RequestLogger(logger))
	api.Use(gin.Recovery())

	admPass := []string{"admin"}
	_ = []string{"user", "admin"}

	// AUTH
	auth := api.Group("")
	auth.POST(LoginPath, rh.Auth.Login(tm.GetRefreshLifetime()))
	auth.GET(RefreshTokenPath, rh.Auth.RefreshToken)

	authAdm := api.Group("")
	authAdm.Use(middleware.CheckJWT(tm, admPass))
	authAdm.POST(CreateUserPath, rh.Auth.CreateUser)
	authAdm.PUT(UpdateUserPath, rh.Auth.UpdateUser)
	authAdm.GET(ListUsersPath, rh.Auth.ListUsers)
	authAdm.DELETE(DeleteUserPath, rh.Auth.DeleteUser)

	server := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           g,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return server
}

func RunServer(server *http.Server) error {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
