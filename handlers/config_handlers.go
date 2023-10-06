package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"institute-person-api/config"
)

type ConfigHandler struct {
	config *config.ConfigInterface
}

func NewConfigHandler(theConfig config.ConfigInterface) *ConfigHandler {
	this := &ConfigHandler{}
	this.config = &theConfig
	return this
}

func (h *ConfigHandler) GetConfig(responseWriter http.ResponseWriter, request *http.Request) {
	// transaction logging
	correltionId, _ := exec.Command("uuidgen").Output()
	stringId := strings.TrimSuffix(string(correltionId), "\n")
	log.Printf("TRANSACTION CID: %s Get Config Start", stringId)
	defer log.Printf("TRANSACTION CID: %s Get Config Complete", stringId)

	// Return the Config object as JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(h.config)
}
