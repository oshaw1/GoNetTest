package handler

import (
	"log"
	"net/http"

	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/pageGeneration"
	"github.com/oshaw1/go-net-test/internal/scheduler"
)

type DashboardHandler struct {
	repository *dataManagement.Repository
	generator  *pageGeneration.PageGenerator
	scheduler  *scheduler.Scheduler
}

func NewDashboardHandler(repo *dataManagement.Repository, templatePath string, scheduler *scheduler.Scheduler) *DashboardHandler {
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
		scheduler:  scheduler,
	}
}

func (h *DashboardHandler) ServeTestQuadrant(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	testType := r.URL.Query().Get("type")

	isDateChange := date != "" && r.URL.Query().Get("refresh_dropdown") != "false"

	data, err := h.generator.GenerateTestQuadrant(date, testType)
	if err != nil {
		handleError(w, "Error generating test data", err, 500)
		return
	}

	if isDateChange {
		w.Write([]byte(`<div id="test-selection" hx-swap-oob="true">`))
		h.generator.RenderTestSelection(w, data)
		w.Write([]byte(`</div>`))
	}

	h.generator.RenderTestResults(w, data)
}

// ServeTestDatesSidebar returns just the date list fragment, used by the
// sidebar's own self-poll so new dates show up without re-rendering (and
// disrupting) the results pane.
func (h *DashboardHandler) ServeTestDatesSidebar(w http.ResponseWriter, r *http.Request) {
	dates, err := h.repository.GetTestDirectories()
	if err != nil {
		handleError(w, "Error retrieving test dates", err, 500)
		return
	}

	h.generator.RenderTestDatesSidebar(w, &pageGeneration.TestQuadrantData{Dates: dates})
}

// ServeTestQuadrantFull re-renders the entire test-quadrant (sidebar, type
// selector and results) from scratch with no date/type pinned, so it falls
// back to whatever is now the latest date. Used after a destructive change
// like deleting a date, where the previously-selected date/type may no
// longer exist.
func (h *DashboardHandler) ServeTestQuadrantFull(w http.ResponseWriter, r *http.Request) {
	data, err := h.generator.GenerateTestQuadrant("", "")
	if err != nil {
		handleError(w, "Error generating test data", err, 500)
		return
	}

	h.generator.RenderTestQuadrant(w, data)
}

func (h *DashboardHandler) ServeSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	schedulerData, err := h.generator.GenerateSchedulerQuadrant()
	if err != nil {
		handleError(w, "Error generating scheduler data", err, 500)
		return
	}

	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = scheduler.SortByDate
	}
	schedulerData.SortBy = sortBy

	if h.scheduler != nil {
		h.scheduler.Mu.RLock()
		schedulerData.Schedule = h.scheduler.Schedule
		schedulerData.Tasks = scheduler.SortTasks(h.scheduler.Schedule, sortBy)
		h.scheduler.Mu.RUnlock()

		h.generator.RenderSchedule(w, schedulerData)
	}
}

func (h *DashboardHandler) ServeDashboard(w http.ResponseWriter, r *http.Request) {
	if err := h.generator.RenderDashboard(w); err != nil {
		handleError(w, "Error rendering dashboard", err, 500)
	}
}

func (h *DashboardHandler) ServeControlQuadrant(w http.ResponseWriter, r *http.Request) {
	// Serve control quadrant template with data
}
