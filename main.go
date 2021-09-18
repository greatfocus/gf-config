package main

import (
	"log"
	"net/http"
	"time"

	"github.com/greatfocus/gf-config/router"
)

// Entry point to service
func main() {
	// initialize variables
	mux := router.LoadRouter()
	addr := ":5000"

	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(200) * time.Second,
		WriteTimeout:   time.Duration(200) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	// create server connection
	log.Println("Listening to port HTTP", addr)
	log.Fatal(srv.ListenAndServe())
}
