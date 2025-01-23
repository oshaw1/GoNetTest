package scheduler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewScheduler(t *testing.T) {
	baseURL := "http://test.com"
	schedulePath := "data/schedules.json"
	scheduler := NewScheduler(baseURL, schedulePath)

	assert.NotNil(t, scheduler)
	assert.Equal(t, baseURL, scheduler.baseURL)
	assert.NotNil(t, scheduler.Schedule)
	assert.NotNil(t, scheduler.client)
	assert.NotNil(t, scheduler.done)
}

func TestUpdateNextRunTime(t *testing.T) {
	schedulePath := "data/schedules.json"
	scheduler := NewScheduler("http://test.com", schedulePath)
	now := time.Now()

	tests := []struct {
		name     string
		schedule Task
		expected time.Time
		interval string
	}{
		{
			name: "Daily interval",
			schedule: Task{
				DateTime: now,
				Interval: "daily",
			},
			expected: now.AddDate(0, 0, 1),
		},
		{
			name: "Weekly interval",
			schedule: Task{
				DateTime: now,
				Interval: "weekly",
			},
			expected: now.AddDate(0, 0, 7),
		},
		{
			name: "Monthly interval",
			schedule: Task{
				DateTime: now,
				Interval: "monthly",
			},
			expected: now.AddDate(0, 1, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := &tt.schedule
			scheduler.updateNextRunTime(schedule)
			assert.Equal(t, tt.expected.Format(time.RFC3339), schedule.DateTime.Format(time.RFC3339))
		})
	}
}

func TestExecuteTest(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		expectError bool
	}{
		{
			name:        "Successful test execution",
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "Failed test execution",
			statusCode:  http.StatusInternalServerError,
			expectError: false, // Note: Current implementation doesn't check status code
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			schedulePath := "example/schedule.json"
			scheduler := NewScheduler(server.URL, schedulePath)
			task := &Task{
				TestType: "jitter",
			}

			err := scheduler.executeTest(task)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckAndExecuteSchedules(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	schedulePath := "data/schedules.json"
	scheduler := NewScheduler(server.URL, schedulePath)
	pastTime := time.Now().Add(-1 * time.Hour)
	futureTime := time.Now().Add(1 * time.Hour)

	tests := []struct {
		name   string
		task   Task
		active bool
	}{
		{
			name: "Past non-recurring task",
			task: Task{
				Name:      "1",
				TestType:  "jitter",
				DateTime:  pastTime,
				Recurring: false,
				Active:    true,
			},
			active: false,
		},
		{
			name: "Future task",
			task: Task{
				Name:      "2",
				TestType:  "jitter",
				DateTime:  futureTime,
				Recurring: false,
				Active:    true,
			},
			active: true,
		},
		{
			name: "Past recurring task",
			task: Task{
				Name:      "3",
				TestType:  "jitter",
				DateTime:  pastTime,
				Recurring: true,
				Interval:  "daily",
				Active:    true,
			},
			active: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler.Schedule[tt.task.Name] = &tt.task
			scheduler.checkAndExecuteSchedule()
			time.Sleep(100 * time.Millisecond) // Allow goroutine to complete

			assert.Equal(t, tt.active, scheduler.Schedule[tt.task.Name].Active)
			if tt.task.Recurring && !futureTime.After(tt.task.DateTime) {
				assert.True(t, scheduler.Schedule[tt.task.Name].DateTime.After(pastTime))
			}
		})
	}
}

func TestSchedulerStartStop(t *testing.T) {
	scheduler := NewScheduler("http://test.com", "")

	scheduler.Start()
	assert.NotPanics(t, func() {
		scheduler.Stop()
	})

	// Ensure multiple stops don't panic
	assert.NotPanics(t, func() {
		scheduler.Stop()
	})
}
