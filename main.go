package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/greatfocus/gf-config/router"
)

// Entry point to the solution
func main() {
	// Get arguments
	var env, encryptionKey, cert, key string
	env = os.Args[1]
	encryptionKey = os.Args[2]
	cert = os.Args[3]
	key = os.Args[4]

	if env == "" {
		panic("Configure the env")
	}

	if encryptionKey == "" {
		panic("Failed to set encryption Key")
	}

	if cert == "" || key == "" {
		panic("Failed to set security")
	}

	serve(router.LoadRouter())
}

// serve creates server instance
func serve(mux *http.ServeMux) {
	// initialize variables
	addr := ":5000"
	cert := os.Args[3]
	key := os.Args[4]

	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(200) * time.Second,
		WriteTimeout:   time.Duration(200) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
		TLSConfig:      TLSClientConfig(),
	}

	// create server connection
	log.Println("Listening to port secure HTTPS", addr)
	log.Fatal(srv.ListenAndServeTLS(cert, key))

}

// TLSClientConfig update cert and key
func TLSClientConfig() *tls.Config {
	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		ClientCAs:          caCertPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}
	return config
}
