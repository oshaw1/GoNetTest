package scheduler

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Schedule struct {
	ID        string    `json:"id"`
	TestType  string    `json:"test_type"`
	DateTime  time.Time `json:"datetime"`
	Recurring bool      `json:"recurring"`
	Interval  string    `json:"interval,omitempty"` // daily, weekly, monthly
	Active    bool      `json:"active"`
}

type Scheduler struct {
	Schedules map[string]*Schedule
	client    *http.Client
	baseURL   string
	Mu        sync.RWMutex
	done      chan struct{}
}

func NewScheduler(baseURL string) *Scheduler {
	return &Scheduler{
		Schedules: make(map[string]*Schedule),
		client:    &http.Client{Timeout: 30 * time.Second},
		baseURL:   baseURL,
		done:      make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go s.run()
}

func (s *Scheduler) Stop() {
	close(s.done)
}

func (s *Scheduler) run() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndExecuteSchedules()
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) checkAndExecuteSchedules() {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	now := time.Now()
	for _, schedule := range s.Schedules {
		if !schedule.Active {
			continue
		}

		if schedule.DateTime.Before(now) {
			go s.executeTest(schedule)

			if schedule.Recurring {
				s.updateNextRunTime(schedule)
			} else {
				schedule.Active = false
			}
		}
	}
}

func (s *Scheduler) updateNextRunTime(schedule *Schedule) {
	switch schedule.Interval {
	case "daily":
		schedule.DateTime = schedule.DateTime.AddDate(0, 0, 1)
	case "weekly":
		schedule.DateTime = schedule.DateTime.AddDate(0, 0, 7)
	case "monthly":
		schedule.DateTime = schedule.DateTime.AddDate(0, 1, 0)
	}
}

func (s *Scheduler) executeTest(schedule *Schedule) error {
	url := fmt.Sprintf("%s/networktest?test=%s", s.baseURL, schedule.TestType)
	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to execute test: %w", err)
	}
	defer resp.Body.Close()
	return nil
}
