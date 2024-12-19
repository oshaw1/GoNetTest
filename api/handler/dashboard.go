package handler

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/pageGeneration"
)

type DashboardHandler struct {
	repository *dataManagement.Repository
	generator  *pageGeneration.PageGenerator
}

func NewDashboardHandler(repo *dataManagement.Repository, templatePath string) *DashboardHandler {
	if repo == nil {
		log.Fatalf("Repository cannot be nil")
	}

	generator, err := pageGeneration.NewPageGenerator(templatePath, repo)
	if err != nil {
		log.Fatalf("Failed to create page generator: %v", err)
	}

	return &DashboardHandler{
		repository: repo,
		generator:  generator,
	}
}

func (h *DashboardHandler) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("web", "static", "dashboard.html"))
	if err != nil {
		handleError(w, "Error parsing dashboard template", err, 500)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		handleError(w, "Error rendering dashboard", err, 500)
	}
}

func (h *DashboardHandler) GetRecentQuadrant(w http.ResponseWriter, r *http.Request) {
	html, err := h.generator.GenerateRecentQuadrantHTML()
	if err != nil {
		handleError(w, "Error generating recent quadrant", err, 500)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (h *DashboardHandler) GetChart(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		handleError(w, "Date parameter is required", errors.New("missing date parameter"), http.StatusBadRequest)
		return
	}
	testStr := r.URL.Query().Get("test")
	if testStr == "" {
		handleError(w, "Test type parameter is required", errors.New("request missing 'test' parameter"), http.StatusBadRequest)
		return
	}

	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		handleError(w, "Invalid date format", err, http.StatusBadRequest)
		return
	}

	chartHtml, err := h.generator.GenerateICMPChartHTMLGivenDate(dateStr, testStr)
	if err != nil {
		handleError(w, "Error generating chart", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(chartHtml))
}

func (h *DashboardHandler) GetData(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		handleError(w, "Date parameter is required", errors.New("missing date parameter"), http.StatusBadRequest)
		return
	}

	testType := r.URL.Query().Get("type")
	if testType == "" {
		testType = "icmp"
	}

	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		handleError(w, "Invalid date format", err, http.StatusBadRequest)
		return
	}

	result, err := h.repository.GetTestData(dateStr, testType)
	if err != nil {
		handleError(w, "Error occured fetching data", err, http.StatusInternalServerError)
		return
	}

	var dataHTML template.HTML
	switch testType {
	case "icmp":
		dataHTML, err = h.generator.GenerateICMPDataHTML(result.ICMP)
	case "download":
		dataHTML, err = h.generator.GenerateSpeedDataHTML(result.Download)
	case "upload":
		dataHTML, err = h.generator.GenerateSpeedDataHTML(result.Upload)
	default:
		handleError(w, "Invalid test type", fmt.Errorf("unsupported test type: %s", testType), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error generating %s data: %v", testType, err)
	} else {
		log.Printf("Generated %s data HTML: %s", testType, dataHTML)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dataHTML))
}
