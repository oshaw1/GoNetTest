package handler

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/pageGeneration"
)

type DashboardHandler struct {
	repository *dataManagement.Repository
	generator  *pageGeneration.PageGenerator
}

func NewDashboardHandler(repo *dataManagement.Repository, templatePath string) *DashboardHandler {
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
