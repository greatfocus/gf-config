package main

import (
	"log"
	"net/http"
	"time"

	"github.com/greatfocus/gf-config/router"
)

// Entry point to service
func main() {
	//declare variables
	var env, encryptionKey, cert, key string

	// validate arguments
	if env == "" {
		panic("Configure the env")
	}
	if encryptionKey == "" {
		panic("Failed to set encryption Key")
	}
	if cert == "" || key == "" {
		panic("Failed to set security")
	}

	// start
	start()
}

// serve creates server instance
func start() {
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
	log.Println("Listening to port secure HTTPS", addr)
	log.Fatal(srv.ListenAndServe())

}
