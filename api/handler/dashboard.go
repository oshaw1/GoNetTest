package handler

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/dataManagment"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
	"github.com/oshaw1/go-net-test/internal/pageGeneration"
)

type PageHandler struct {
}

func (h PageHandler) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("web", "static", "dashboard.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// change this in the future it should:
// get the data from the datamanagment service
// return that html
func (h PageHandler) GetRecentTestData(w http.ResponseWriter, r *http.Request) {

	dataExists, err := dataManagment.CheckForRecentTestData("data/output", 7)

	if err != nil {
		log.Fatalf("Failed to check recent test data: %v", err)
	}

	if !dataExists {
		html, err := pageGeneration.ReturnNoDataHTML()
		if err != nil {
			http.Error(w, "Error returning no data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}
	// change this to retrieve the result
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

	html, err := pageGeneration.GenerateRecentTestResultHTML(result)
	if err != nil {
		http.Error(w, "Error generating result", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// change this in the future it should:
// check if there is recent charts
// if not return no data html to frontend
// get the charts from the datamanagment service
// return that html
// func (h PageHandler) GetRecentTestCharts(w http.ResponseWriter, r *http.Request) {

// }

func (h PageHandler) RunICMBTest(w http.ResponseWriter, r *http.Request) {
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

	html, err := pageGeneration.GenerateRecentTestResultHTML(result)
	if err != nil {
		http.Error(w, "Error generating result", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
