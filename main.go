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

// Order of the args
// /etc/sysconfig/respect.env
// 1. Vault Url
// 2. Application
// 3. Env
// 4. Encryption key
// 5. CA
// 6. Server Cert
// 7. Server Key
// 8. Client Cert
// 9. Client Key

// Entry point to service
func main() {
	//declare variables
	var env, encryptionKey, cert, key string

	// assign arguments
	env = os.Args[3]
	encryptionKey = os.Args[4]
	cert = os.Args[6]
	key = os.Args[7]

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
	start(env, encryptionKey, cert, key)
}

// serve creates server instance
func start(env, encryptionKey, cert, key string) {
	// initialize variables
	mux := router.LoadRouter()
	addr := ":5000"

	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(200) * time.Second,
		WriteTimeout:   time.Duration(200) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
		TLSConfig:      TLSServerConfig(),
	}

	// create server connection
	log.Println("Listening to port secure HTTPS", addr)
	log.Fatal(srv.ListenAndServeTLS(cert, key))

}

// TLSServerConfig provides config with cert and key
func TLSServerConfig() *tls.Config {
	caCertPEM, err := ioutil.ReadFile(os.Args[5])
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		log.Fatalf("error reading CA certificate: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(os.Args[6], os.Args[7])
	if err != nil {
		log.Fatalf("error reading certificate: %v", err)
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          roots,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}
}
