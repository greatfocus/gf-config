package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/response"
)

// VaultController struct
type VaultController struct{}

// Handler method routes to http methods supported
func (c *VaultController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.getConfig(w, r)
	default:
		err := errors.New("invalid Request")
		response.Error(w, http.StatusNotFound, err)
		return
	}
}

// requestMessage method creates a message request
func (c *VaultController) getConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}

	vault := config.Vault{}
	err = json.Unmarshal(body, &vault)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}

	// validate payload
	if vault.Application == "" || vault.Env == "" || vault.User == "" || vault.Pass == "" {
		derr := errors.New("invalid payload request")
		response.Error(w, http.StatusBadGateway, derr)
		return
	}

	// validate user
	if os.Args[7] != vault.User || os.Args[8] != vault.Pass {
		derr := errors.New("authentication Failed")
		response.Error(w, http.StatusBadGateway, derr)
		return
	}

	var path = "./vault/" + vault.Application + "/" + vault.Impl + "/" + vault.Env + ".json"
	config, err := read(path)
	if err != nil {
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, errors.New("error accessing key vault"))
		return
	}

	config.Env = vault.Env
	config.Impl = vault.Impl

	log.Printf("Received and sent vault config request for %s to %s service", vault.Env, vault.Impl)
	response.Success(w, http.StatusOK, config)
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
