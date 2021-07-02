package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/greatfocus/gf-sframe/config"
	"github.com/greatfocus/gf-sframe/server"
)

// VaultHandler struct
type VaultHandler func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (v VaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v(w, r)
}

// Handler method routes to http methods supported
func (v *VaultHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		v.getConfig(w, r)
	default:
		err := errors.New("invalid Request")
		server.Error(w, http.StatusNotFound, err)
		return
	}
}

// requestMessage method creates a message request
func (v *VaultHandler) getConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		server.Error(w, http.StatusBadGateway, derr)
		return
	}

	vault := config.Impl{}
	err = json.Unmarshal(body, &vault)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		server.Error(w, http.StatusBadGateway, derr)
		return
	}

	// validate payload
	if vault.Application == "" || vault.Env == "" {
		derr := errors.New("invalid payload request")
		server.Error(w, http.StatusBadGateway, derr)
		return
	}

	var path = "./vault/" + vault.Application + "/" + vault.Impl + "/" + vault.Env + ".json"
	config, err := read(path)
	if err != nil {
		log.Printf("Error: %v\n", err)
		server.Error(w, http.StatusUnprocessableEntity, errors.New("error accessing key vault"))
		return
	}

	config.Env = vault.Env
	config.Impl = vault.Impl

	log.Printf("Received and sent vault config request for %s to %s service", vault.Env, vault.Impl)
	server.Success(w, http.StatusOK, config)
}

func read(file string) (config.Config, error) {
	var result = config.Config{}
	if len(file) < 1 {
		derr := errors.New("config file name is empty")
		return result, derr
	}

	jsonFile, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		derr := errors.New("cannot find location of config file")
		return result, derr
	}

	// read the config file
	byteContent, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		derr := errors.New("invalid config formate")
		return result, derr
	}

	// convert the config file bytes to json
	err = json.Unmarshal([]byte(byteContent), &result)
	if err != nil {
		derr := errors.New("invalid config formate")
		return result, derr
	}

	// the closing of our jsonFile so that we can parse it later on
	_ = jsonFile.Close()

	return result, nil
}
