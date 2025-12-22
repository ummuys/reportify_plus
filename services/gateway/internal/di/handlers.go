package di

import (
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/services/gateway/internal/web/handlers"
)

type RESTHandlers struct {
	AuthService       handlers.AuthHandler
	ReportService     handlers.ReportServiceHandler
	DatasourceService handlers.DatasourceHandler
	ReportCache       handlers.ReportCacheHandler
}

func NewRESTHandlers(scs GRPCSC, baseLogger zerolog.Logger) RESTHandlers {
	auth := handlers.NewAuthHandler(scs.AuthService, baseLogger)
	reportSvc := handlers.NewReportServiceHandler(scs.ReportService, baseLogger)
	datasourceSvc := handlers.NewDatasourceHandler(scs.DatasourceService, baseLogger)
	reportCch := handlers.NewReportCacheHandler(scs.ReportCache, baseLogger)

	return RESTHandlers{
		AuthService:       auth,
		ReportService:     reportSvc,
		DatasourceService: datasourceSvc,
		ReportCache:       reportCch,
	}
}
