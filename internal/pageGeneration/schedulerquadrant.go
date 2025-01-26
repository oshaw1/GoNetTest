package pageGeneration

import "github.com/oshaw1/go-net-test/internal/scheduler"

type SchedulerQuadrantData struct {
	QuadrantData
	Schedule map[string]*scheduler.Task
}

func (pg *PageGenerator) GenerateSchedulerQuadrant() (*SchedulerQuadrantData, error) {
	return &SchedulerQuadrantData{
		QuadrantData: QuadrantData{Title: "Scheduler"},
	}, nil
}
