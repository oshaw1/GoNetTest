package scheduler

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExportImportSchedules(t *testing.T) {
	s := NewScheduler("http://test.com", "test_schedules.json")
	testFile := "test_schedules.json"
	defer os.Remove(testFile)

	testTime := time.Now().UTC() // Use UTC for consistent timezone
	schedule := &Task{
		Name:      "test123",
		TestType:  "jitter",
		DateTime:  testTime,
		Recurring: true,
		Interval:  "daily",
		Active:    true,
	}

	s.Schedule[schedule.Name] = schedule

	err := s.ExportSchedule(testFile)
	assert.NoError(t, err)

	s2 := NewScheduler("http://test.com", "test_schedules.json")
	err = s2.ImportSchedule(testFile)
	assert.NoError(t, err)

	assert.Equal(t, s.Schedule["test123"], s2.Schedule["test123"])
}

func TestExportSchedulesError(t *testing.T) {
	s := NewScheduler("http://test.com", "")

	err := s.ExportSchedule("")
	assert.Error(t, err, "Expected an error when exporting to an invalid path")
}

func TestImportSchedulesError(t *testing.T) {
	s := NewScheduler("http://test.com", "test_schedules.json")
	err := s.ImportSchedule("nonexistent.json")
	assert.Error(t, err)
}

func TestDeleteSchedule(t *testing.T) {
	s := NewScheduler("http://test.com", "test_schedules.json")
	schedule := &Task{
		Name:      "test123",
		TestType:  "jitter",
		DateTime:  time.Now(),
		Recurring: true,
		Interval:  "daily",
		Active:    true,
	}
	s.Schedule[schedule.Name] = schedule

	err := s.DeleteSchedule("test123")
	assert.NoError(t, err)
	assert.Nil(t, s.Schedule["test123"])

	err = s.DeleteSchedule("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
