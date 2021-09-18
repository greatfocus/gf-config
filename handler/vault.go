package handler

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

// Vault struct
type Vault struct {
	VaultHandler func(http.ResponseWriter, *http.Request)
}

func (v Vault) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		v.getConfig(w, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET, POST, PUT, DELETE")
}

// requestMessage method creates a message request
func (v *Vault) getConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		server.Error(w, http.StatusBadRequest, derr)
		return
	}

	vault := config.Impl{}
	err = json.Unmarshal(body, &vault)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		server.Error(w, http.StatusBadRequest, derr)
		return
	}

	// validate payload
	if vault.Service == "" || vault.Env == "" {
		derr := errors.New("invalid payload request")
		server.Error(w, http.StatusBadRequest, derr)
		return
	}

	if vault.VaultUser != os.Getenv("VALT_USER") || vault.VaultPass != os.Getenv("VALT_PASS") {
		derr := errors.New("invalid credentials")
		server.Error(w, http.StatusUnauthorized, derr)
		return
	}

	var path = os.Getenv("VALT_PATH") + "/vault/" + vault.Service + "/" + vault.Env + ".json"
	config, err := read(path)
	if err != nil {
		log.Printf("Error: %v\n", err)
		server.Error(w, http.StatusUnprocessableEntity, errors.New("error accessing vault"))
		return
	}

	config.Env = vault.Env
	config.Impl = vault.Service

	log.Printf("Received and sent vault config request for %s to %s service", vault.Env, vault.Service)
	server.Success(w, http.StatusOK, config)
}

func read(file string) (config.Config, error) {
	var result = config.Config{}
	if len(file) < 1 {
		derr := errors.New("config is missing")
		return result, derr
	}

	jsonFile, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		derr := errors.New("config is missing")
		return result, derr
	}

	// read the config file
	byteContent, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		derr := errors.New("invalid config formart")
		return result, derr
	}

	// convert the config file bytes to json
	err = json.Unmarshal([]byte(byteContent), &result)
	if err != nil {
		derr := errors.New("invalid config formart")
		return result, derr
	}

	// the closing of our jsonFile so that we can parse it later on
	_ = jsonFile.Close()

	return result, nil
}
