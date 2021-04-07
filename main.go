package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/greatfocus/gf-config/router"
	_ "github.com/greatfocus/pq"
)

// Entry point to the solution
func main() {
	// Get arguments
	var env, user, pass, cert, key, port string
	env = os.Args[1]
	user = os.Args[2]
	pass = os.Args[3]
	cert = os.Args[4]
	key = os.Args[5]
	port = os.Args[6]
	if env == "" {
		panic("Configure the env")
	}

	if user == "" || pass == "" {
		panic("Failed to start config service")
	}

	if cert == "" || key == "" {
		panic("Failed to set security")
	}

	if port == "" {
		panic("Failed to start port")
	}

	serve(router.LoadRouter())
}

// serve creates server instance
func serve(mux *http.ServeMux) {
	// initialize variables
	cert := os.Args[4]
	key := os.Args[5]
	addr := ":" + os.Args[6]
	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(200) * time.Second,
		WriteTimeout:   time.Duration(200) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	// create server connection
	log.Println("Listening to port secure HTTPS", addr)
	log.Fatal(srv.ListenAndServeTLS(cert, key))

}
