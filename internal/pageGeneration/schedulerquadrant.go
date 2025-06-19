package pageGeneration

import (
	"bytes"
	"net/http"

	"github.com/oshaw1/go-net-test/internal/scheduler"
)

type SchedulerQuadrantData struct {
	QuadrantData
	Schedule map[string]*scheduler.Task
}

func (pg *PageGenerator) GenerateSchedulerQuadrant() (*SchedulerQuadrantData, error) {
	return &SchedulerQuadrantData{
		QuadrantData: QuadrantData{Title: "Scheduler"},
		Schedule:     make(map[string]*scheduler.Task), // Initialize empty map
	}, nil
}

func (pg *PageGenerator) RenderSchedule(w http.ResponseWriter, data *SchedulerQuadrantData) error {
	var buf bytes.Buffer
	err := pg.templates.ExecuteTemplate(&buf, "schedule", data)
	if err != nil {
		return err
	}

	w.Write(buf.Bytes())
	return nil
}
