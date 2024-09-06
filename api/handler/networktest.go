package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oshaw1/go-net-test/internal/dataManagment"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

type NetworkTestHandler struct {
}

func (h NetworkTestHandler) HandleICMPNetworkTest(w http.ResponseWriter, r *http.Request) {
	host := networkTesting.GetHost(r)

	result, err := networkTesting.TestNetwork(host)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error performing network test: %v", err), http.StatusInternalServerError)
		return
	}

	err = dataManagment.SaveTestData(result, "icmp")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving test result: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h NetworkTestHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	testType := r.URL.Query().Get("test")
	if testType == "" {
		log.Print("type of test not found in request")
		http.Error(w, "Invalid request. Please specify a type of test e.g. icmp", http.StatusBadRequest)
	}
	if date == "" {
		now := time.Now()
		date = now.Format("2006-01-02")
		log.Printf("No date in GetICMPResults request, using current date: %s", date)
	}
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Invalid date format in request: %s", date)
		http.Error(w, "Invalid date format. Please use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	results, err := dataManagment.GetTestResults(date, testType)
	if err != nil {
		log.Printf("Error getting test results from %v caused by: %s", date, err)
		http.Error(w, "Error getting test results consult server log for information.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
