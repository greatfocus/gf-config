package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-config/controllers"
)

// LoadRouter is exported and used in main.go
func LoadRouter() *http.ServeMux {
	mux := http.NewServeMux()
	loadHandlers(mux)
	log.Println("Created routes with controllers")
	return mux
}

// configRoute created all routes and handlers relating to controller
func loadHandlers(mux *http.ServeMux) {
	// Initialize controller
	vaultController := controllers.VaultController{}

	// Initialize routes
	mux.HandleFunc("/config/vaul", vaultController.Handler)
}
