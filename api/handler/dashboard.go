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

func (h PageHandler) GetRecentQuadrant(w http.ResponseWriter, r *http.Request) {
	result, err := dataManagment.ParseRecentTestJSON()
	if err != nil {
		http.Error(w, "Error parsing recent results", http.StatusInternalServerError)
		log.Printf("Error parsing recent results, Caused by: %v", err)
		return
	}

	html, err := pageGeneration.GenerateRecentQuadrantHTML(result)
	if err != nil {
		http.Error(w, "Error generating result", http.StatusInternalServerError)
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
