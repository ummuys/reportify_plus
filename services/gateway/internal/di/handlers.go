package di

import (
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/services/gateway/internal/web/handlers"
)

type RESTHandlers struct {
	Auth   handlers.AuthHandler
	Report handlers.ReportHandler
}

func NewRESTHandlers(scs GRPCSC, baseLogger zerolog.Logger) RESTHandlers {
	auth := handlers.NewAuthHandler(scs.Auth, baseLogger)
	report := handlers.NewReportHandler(scs.Report, baseLogger)
	return RESTHandlers{Auth: auth, Report: report}
}
