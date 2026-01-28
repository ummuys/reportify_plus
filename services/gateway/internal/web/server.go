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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/ummuys/reportify/services/gateway/docs"
)

func CreateServer(cfg config.GatewayServiceConfig, rh di.RESTHandlers, tm pkg.TokenManager, baseLogger zerolog.Logger) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()

	logger := baseLogger.With().Str("component", "srv").Logger()

	g.NoRoute(func(g *gin.Context) {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := g.Group("/api/v1")
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8088"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api.Use(middleware.RequestLogger(logger))
	api.Use(gin.Recovery())

	admPass := []string{"admin"}
	basePass := []string{"user", "admin"}

	report := api.Group("")
	report.Use(middleware.CheckJWT(tm, basePass))
	report.POST(CreateReportPath, rh.ReportService.CreateReport)
	report.GET(ListReportsPath, rh.ReportService.ListReports)
	report.GET(ReportStatusPath, rh.ReportService.ReportStatus)
	report.GET(ReportInfoPath, rh.ReportService.ReportInfo)
	report.DELETE(DeleteReports, rh.ReportService.DeleteReports)
	report.DELETE(DeleteReport, rh.ReportService.DeleteReport)

	datasource := api.Group("")
	datasource.Use(middleware.CheckJWT(tm, basePass))
	datasource.GET(ListSchemasPath, rh.DatasourceService.ListSchemas)
	datasource.GET(ListTablesPath, rh.DatasourceService.ListTables)
	datasource.GET(ListColumnsPath, rh.DatasourceService.ListColumns)

	auth := api.Group("")
	auth.POST(LoginPath, rh.AuthService.Login(tm.GetRefreshLifetime()))
	auth.GET(RefreshTokenPath, rh.AuthService.RefreshToken)

	authAdm := api.Group("")
	authAdm.Use(middleware.CheckJWT(tm, admPass))
	authAdm.POST(CreateUserPath, rh.AuthService.CreateUser)
	authAdm.PUT(UpdateUserPath, rh.AuthService.UpdateUser)
	authAdm.GET(ListsUserPath, rh.AuthService.ListUsers)
	authAdm.DELETE(DeletePath, rh.AuthService.DeleteUser)

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
