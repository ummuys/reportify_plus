package web

import (
	"net/http"
	// _ "github.com/ummuys/reportify/docs"
)

func CreateServer() *http.Server {
	return nil
}

func RunServer(server *http.Server) error {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
