package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/greatfocus/gf-config/router"
	_ "github.com/greatfocus/pq"
)

// Entry point to the solution
func main() {
	// Get arguments
	var env, user, pass, cert, key, port string
	env = os.Args[1]
	cert = os.Args[2]
	key = os.Args[3]
	port = os.Args[4]
	user = os.Args[5]
	pass = os.Args[6]
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
	addr := ":5001"
	ca := os.Args[2]
	cert := os.Args[3]
	key := os.Args[4]

	// load CA certificate file and add it to list of client CAs
	ca = filepath.Clean(ca)
	caCertFile, err := ioutil.ReadFile(filepath.Clean(ca))
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertFile)

	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(200) * time.Second,
		WriteTimeout:   time.Duration(200) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
		TLSConfig: &tls.Config{
			ClientAuth:         tls.RequireAndVerifyClientCert,
			ClientCAs:          certPool,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
		},
	}

	// create server connection
	log.Println("Listening to port secure HTTPS", addr)
	log.Fatal(srv.ListenAndServeTLS(cert, key))

}
