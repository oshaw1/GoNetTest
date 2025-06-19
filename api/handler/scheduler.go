package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oshaw1/go-net-test/internal/scheduler"
)

type SchedulerHandler struct {
	scheduler *scheduler.Scheduler
}

func NewSchedulerHandler(scheduler *scheduler.Scheduler) *SchedulerHandler {
	return &SchedulerHandler{scheduler: scheduler}
}

func (h *SchedulerHandler) HandleCreateSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task scheduler.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	task.CreatedOn = time.Now()

	h.scheduler.Mu.Lock()
	h.scheduler.Schedule[id] = &task
	h.scheduler.Mu.Unlock()

	response := map[string]*scheduler.Task{id: &task}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *SchedulerHandler) HandleGetSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.scheduler.Mu.RLock()
	defer h.scheduler.Mu.RUnlock()

	json.NewEncoder(w).Encode(h.scheduler.Schedule) // Returns map[id]*Task directly
}

func (h *SchedulerHandler) HandleExportSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filepath := r.URL.Query().Get("filepath")
	if filepath == "" {
		filepath = "data/schedules.json"
	}

	if err := h.scheduler.ExportSchedule(filepath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Schedules exported successfully to %s", filepath)))
}

func (h *SchedulerHandler) HandleImportSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filepath := r.URL.Query().Get("filepath")
	if filepath == "" {
		filepath = "data/schedules.json"
	}

	if err := h.scheduler.ImportSchedule(filepath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Schedules imported successfully from %s", filepath)))
}

func (h *SchedulerHandler) HandleDeleteSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing schedule ID", http.StatusBadRequest)
		return
	}

	if err := h.scheduler.DeleteSchedule(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SchedulerHandler) HandleEditSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	var updatedTask scheduler.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	editedTask, err := h.scheduler.EditTask(id, updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(editedTask)
}
