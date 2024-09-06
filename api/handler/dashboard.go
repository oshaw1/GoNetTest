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
	testTypes := []string{"icmp"} // Add more test types as needed
	var results []interface{}

	for _, testType := range testTypes {
		if testType == "icmp" {
			result, err := dataManagment.GetRecentICMPTestResult(testType)
			if err != nil {
				log.Print("Error getting recent test result")
			}
			results = append(results, result)
			continue // Skip this test type if there's an error
		}
	}

	if len(results) == 0 {
		http.Error(w, "No recent test results available", http.StatusNotFound)
		return
	}

	html, err := pageGeneration.GenerateRecentQuadrantHTML(results)
	if err != nil {
		http.Error(w, "Error generating result", http.StatusInternalServerError)
		log.Printf("Error generating result: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (h PageHandler) RunICMPTest(w http.ResponseWriter, r *http.Request) {
	host := networkTesting.GetHost(r)

	result, err := networkTesting.TestNetwork(host)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error performing network test: %v", err), http.StatusInternalServerError)
		return
	}
	var resultInterface interface{} = result
	err = dataManagment.SaveTestData(result, "icmp")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving test result: %v", err), http.StatusInternalServerError)
		return
	}

	html, err := pageGeneration.GenerateRecentQuadrantHTML([]interface{}{resultInterface})
	if err != nil {
		http.Error(w, "Error generating result", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
