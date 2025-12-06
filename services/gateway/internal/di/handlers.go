package di

import (
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/services/gateway/internal/web/handlers"
)

type RESTHandlers struct {
	Auth handlers.AuthHandler
}

func NewRESTHandlers(scs GRPCSC, baseLogger zerolog.Logger) RESTHandlers {
	auth := handlers.NewAuthHandler(scs.Auth, baseLogger)
	return RESTHandlers{Auth: auth}
}
