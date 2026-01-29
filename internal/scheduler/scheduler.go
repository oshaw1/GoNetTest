package scheduler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Task struct {
	Name       string     `json:"name"`
	TestType   string     `json:"test_type,omitempty"`
	ChartType  string     `json:"chart_type,omitempty"`
	RecentDays int        `json:"recent_days,omitempty"`
	DateTime   time.Time  `json:"datetime"`
	Recurring  bool       `json:"recurring"`
	Interval   string     `json:"interval,omitempty"`
	Active     bool       `json:"active"`
	LastRan    *time.Time `json:"last_ran,omitempty"`
	CreatedOn  time.Time  `json:"created_on"`
}

type Scheduler struct {
	Schedule     map[string]*Task
	client       *http.Client
	baseURL      string
	Mu           sync.RWMutex
	done         chan struct{}
	schedulePath string
}

func NewScheduler(baseURL string, schedulePath string) *Scheduler {
	s := &Scheduler{
		Schedule:     make(map[string]*Task),
		client:       &http.Client{Timeout: 30 * time.Second},
		baseURL:      baseURL,
		done:         make(chan struct{}),
		schedulePath: schedulePath,
	}
	return s
}

func (s *Scheduler) Start() {
	if err := s.ImportSchedule(s.schedulePath); err != nil {
		log.Printf("Scheduler could not import schedule due to: %s", err)
	} else {
		tasksJSON, err := json.MarshalIndent(s.Schedule, "", "  ")
		if err != nil {
			log.Printf("Error formatting tasks: %v", err)
			return
		}
		log.Printf("Scheduler Successfully Loaded %d tasks from %s \nScheduled Tasks: %s ", len(s.Schedule), s.schedulePath, tasksJSON)
	}

	go s.run()
}

func (s *Scheduler) Stop() {
	select {
	case <-s.done:
		return
	default:
		s.ExportSchedule(s.schedulePath)
		close(s.done)
	}
}

func (s *Scheduler) run() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndExecuteSchedule()
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) checkAndExecuteSchedule() {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	now := time.Now().Format(time.RFC3339)
	nowTime, _ := time.Parse(time.RFC3339, now)

	for _, schedule := range s.Schedule {
		if !schedule.Active || !schedule.DateTime.Before(nowTime) {
			continue
		}

		schedule.LastRan = &nowTime

		if schedule.ChartType != "" {
			if schedule.RecentDays >= 0 {
				go s.executeHistoricChart(schedule)
			} else {
				go s.executeChart(schedule)
			}
		}
		if schedule.TestType != "" {
			go s.executeTest(schedule)
		}

		if schedule.Recurring {
			s.updateNextRunTime(schedule)
		} else {
			schedule.Active = false
		}
		s.ExportSchedule(s.schedulePath)
	}
}

func (s *Scheduler) updateNextRunTime(schedule *Task) {
	now := time.Now()

	// First bring the date up to current if it's in the past
	for schedule.DateTime.Before(now) {
		switch schedule.Interval {
		case "daily":
			schedule.DateTime = schedule.DateTime.AddDate(0, 0, 1)
		case "weekly":
			schedule.DateTime = schedule.DateTime.AddDate(0, 0, 7)
		case "monthly":
			schedule.DateTime = schedule.DateTime.AddDate(0, 1, 0)
		case "bimonthly":
			schedule.DateTime = schedule.DateTime.AddDate(0, 2, 0)
		case "biannually":
			schedule.DateTime = schedule.DateTime.AddDate(0, 6, 0)
		case "annually":
			schedule.DateTime = schedule.DateTime.AddDate(1, 0, 0)
		}
	}
}

func (s *Scheduler) executeTest(schedule *Task) error {
	url := fmt.Sprintf("%s/networktest?test=%s", s.baseURL, schedule.TestType)
	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to execute test: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (s *Scheduler) executeChart(schedule *Task) error {
	currentDate := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("%s/charts/generate?test=%s&date=%s", s.baseURL, schedule.ChartType, currentDate)

	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to generate chart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chart generation failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Scheduler) executeHistoricChart(schedule *Task) error {
	url := fmt.Sprintf("%s/charts/generate-historic?test=%s&days=%d", s.baseURL, schedule.ChartType, schedule.RecentDays)

	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to generate chart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chart generation failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
