package handler

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

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
// get the result from the datamanagment service
func (h PageHandler) GetRecentQuadrant(w http.ResponseWriter, r *http.Request) {
	// change this to retrieve the result
	result, err := networkTesting.TestNetwork("localhost", 8080)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error performing network test: %v", err), http.StatusInternalServerError)
		return
	}

	html, err := pageGeneration.GenerateRecentQuadrantHTML(result)
	if err != nil {
		http.Error(w, "Error generating result", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (h PageHandler) RunICMBTest(w http.ResponseWriter, r *http.Request) {
	host, port := networkTesting.GetNetworkParams(r)

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

	html, err := pageGeneration.GenerateRecentQuadrantHTML(result)
	if err != nil {
		http.Error(w, "Error generating result", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
