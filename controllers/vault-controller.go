package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
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
	case http.MethodGet:
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

	// validate user
	if os.Args[2] != vault.User || os.Args[3] != vault.Pass {
		derr := errors.New("authentication Failed")
		response.Error(w, http.StatusBadGateway, derr)
		return
	}

	var path = "./vault/" + vault.Application + "/" + vault.Env + ".json"
	config := read(path)
	response.Success(w, http.StatusOK, config)
}

func read(file string) config.Config {
	log.Println("Reading configuration file")
	if len(file) < 1 {
		log.Fatal(fmt.Println("config file name is empty"))
	}

	jsonFile, err := os.OpenFile(file, os.O_RDONLY, 0600)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(fmt.Println("cannot find location of config file", err))
	}

	// read the config file
	byteContent, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(fmt.Println("invalid config format", err))
	}

	// convert the config file bytes to json
	var result = config.Config{}
	err = json.Unmarshal([]byte(byteContent), &result)
	if err != nil {
		log.Fatal(fmt.Println("Invalid config format", err))
	}

	// the closing of our jsonFile so that we can parse it later on
	_ = jsonFile.Close()

	return result
}
