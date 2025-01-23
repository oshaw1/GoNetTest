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

	var schedule scheduler.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	schedule.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	schedule.Active = true

	h.scheduler.Mu.Lock()
	h.scheduler.Schedules[schedule.ID] = &schedule
	h.scheduler.Mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

func (h *SchedulerHandler) HandleGetSchedules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.scheduler.Mu.RLock()
	schedules := make([]*scheduler.Schedule, 0, len(h.scheduler.Schedules))
	for _, s := range h.scheduler.Schedules {
		schedules = append(schedules, s)
	}
	h.scheduler.Mu.RUnlock()

	json.NewEncoder(w).Encode(schedules)
}
