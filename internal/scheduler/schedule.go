package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (s *Scheduler) ExportSchedule(filename string) error {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	file := filepath.Clean(filename)
	dir := filepath.Dir(file)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	data, err := json.MarshalIndent(s.Schedule, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schedules: %w", err)
	}

	if err := os.WriteFile(file, data, 0644); err != nil {
		return fmt.Errorf("failed to write schedules file: %w", err)
	}

	return nil
}

func (s *Scheduler) ImportSchedule(filename string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	file := filepath.Clean(filename)
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read schedules file: %w", err)
	}

	var schedule map[string]*Task
	if err := json.Unmarshal(data, &schedule); err != nil {
		return fmt.Errorf("failed to unmarshal schedules: %w", err)
	}

	s.Schedule = schedule
	return nil
}

func (s *Scheduler) DeleteSchedule(id string) error {
	s.Mu.Lock()

	if _, exists := s.Schedule[id]; !exists {
		return fmt.Errorf("schedule with ID %s not found", id)
	}

	delete(s.Schedule, id)
	s.Mu.Unlock()
	if err := s.ExportSchedule(s.schedulePath); err != nil {
		return fmt.Errorf("failed to export schedule: %w", err)
	}
	return nil
}

func (s *Scheduler) EditTask(taskID string, updatedTask Task) (*Task, error) {
	task, exists := s.Schedule[taskID]
	if !exists {
		fmt.Printf("DEBUG: Task %s not found\n", taskID)
		return nil, fmt.Errorf("task with ID %s not found", taskID)
	}

	task.Name = updatedTask.Name
	task.Recurring = updatedTask.Recurring
	task.Active = updatedTask.Active
	task.Interval = updatedTask.Interval

	if updatedTask.TestType != "" {
		task.TestType = updatedTask.TestType
		task.ChartType = ""
		task.RecentDays = 0
	} else if updatedTask.ChartType != "" {
		task.ChartType = updatedTask.ChartType
		task.TestType = ""
		task.RecentDays = updatedTask.RecentDays
	}

	if !updatedTask.DateTime.IsZero() {
		year, month, day := updatedTask.DateTime.Date()
		hour, min, sec := updatedTask.DateTime.Clock()
		task.DateTime = time.Date(year, month, day, hour, min, sec, 0, time.Local)
	}

	fmt.Printf("DEBUG: Final task: %+v\n", task)

	if err := s.ExportSchedule(s.schedulePath); err != nil {
		return nil, fmt.Errorf("failed to export schedule: %w", err)
	}

	return task, nil
}
