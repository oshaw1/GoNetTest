package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oshaw1/go-net-test/internal/dataManagment"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

type NetworkTestHandler struct {
}

// this should use the data managment to save the results/the chart when that is done
func (h NetworkTestHandler) HandleICMPNetworkTest(w http.ResponseWriter, r *http.Request) {
	host := networkTesting.GetHost(r)

	result, err := networkTesting.TestNetwork(host)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error performing network test: %v", err), http.StatusInternalServerError)
		return
	}

	err = dataManagment.SaveTestData(result, "icmb")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving test result: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
