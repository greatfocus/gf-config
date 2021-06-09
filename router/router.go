package router

import (
	"net/http"

	"github.com/greatfocus/gf-config/handlers"
	"github.com/greatfocus/gf-sframe/server"
)

// LoadRouter is exported and used in main.go
func LoadRouter() *http.ServeMux {
	mux := http.NewServeMux()
	loadHandlers(mux)
	return mux
}

// configRoute created all routes and handlers relating to controller
func loadHandlers(mux *http.ServeMux) {
	// Initialize routes
	mux.Handle("/config/vault", server.Use(new(handlers.VaultHandler),
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.WithoutAuth()))
}
