package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func NetworkTest(w http.ResponseWriter, r *http.Request) {
	conf, err := config.Load("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	host, err := config.SetupHost(conf, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	port, err := config.SetupPort(conf, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := networkTesting.TestNetwork(host, port)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error performing network test: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
