package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/dataManagment"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

type NetworkTestHandler struct {
}

// this should use the data managment to save the results/the chart when that is done
func (h NetworkTestHandler) HandleICMPNetworkTest(w http.ResponseWriter, r *http.Request) {
	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	host, err := networkTesting.SetupHost(conf, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	port, err := networkTesting.SetupPort(conf, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := networkTesting.TestNetwork(host, port)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error performing network test: %v", err), http.StatusInternalServerError)
		return
	}

	err = dataManagment.SaveICMBTestData(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving test result: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
