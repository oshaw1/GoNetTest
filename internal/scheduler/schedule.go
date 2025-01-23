package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	s.Mu.Lock()

	task, exists := s.Schedule[taskID]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", taskID)
	}

	if updatedTask.Name != "" {
		task.Name = updatedTask.Name
	}
	if updatedTask.ChartType != "" {
		task.ChartType = updatedTask.ChartType
	}
	if updatedTask.RecentDays != 0 {
		task.RecentDays = updatedTask.RecentDays
	}
	if updatedTask.TestType != "" {
		task.TestType = updatedTask.TestType
	}
	if !updatedTask.DateTime.IsZero() {
		task.DateTime = updatedTask.DateTime
	}
	task.Recurring = updatedTask.Recurring
	if updatedTask.Interval != "" {
		task.Interval = updatedTask.Interval
	}
	task.Active = updatedTask.Active
	s.Mu.Unlock()
	if err := s.ExportSchedule(s.schedulePath); err != nil {
		return nil, fmt.Errorf("failed to export schedule: %w", err)
	}

	return task, nil
}
