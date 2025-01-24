package handler

import (
	"log"
	"net/http"

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

func (h *DashboardHandler) HandleTestQuadrant(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	testType := r.URL.Query().Get("type")

	data, err := h.generator.GenerateTestQuadrant(date, testType)
	if err != nil {
		handleError(w, "Error generating test data", err, 500)
		return
	}

	if testType == "" {
		h.generator.RenderTestSelection(w, data)
		h.generator.RenderTestResults(w, data)
	} else {
		h.generator.RenderTestResults(w, data)
	}
}

func (h *DashboardHandler) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	if err := h.generator.RenderDashboard(w); err != nil {
		handleError(w, "Error rendering dashboard", err, 500)
	}
}

func (h *DashboardHandler) ServeTestQuadrant(w http.ResponseWriter, r *http.Request) {
	data, err := h.generator.GenerateTestQuadrant("", "")
	if err != nil {
		handleError(w, "Error generating test quadrant", err, 500)
		return
	}
	h.generator.RenderTestQuadrant(w, data)
}

func (h *DashboardHandler) ServeGenerateQuadrant(w http.ResponseWriter, r *http.Request) {
	// Serve generate quadrant template with data
}

func (h *DashboardHandler) ServeControlQuadrant(w http.ResponseWriter, r *http.Request) {
	// Serve control quadrant template with data
}

func (h *DashboardHandler) ServeSchedulerQuadrant(w http.ResponseWriter, r *http.Request) {
	// Serve scheduler quadrant template with data
}
