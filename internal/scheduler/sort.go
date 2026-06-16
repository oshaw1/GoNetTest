package scheduler

import "sort"

// TaskEntry pairs a schedule ID with its task so ordering can be
// controlled explicitly (map iteration order is otherwise arbitrary).
type TaskEntry struct {
	ID string
	*Task
}

// ActionType returns the kind of work the task performs, used for
// grouping/sorting by action type ("test" or "chart" plus its sub-type).
func (t *Task) ActionType() string {
	if t.TestType != "" {
		return "test:" + t.TestType
	}
	return "chart:" + t.ChartType
}

const (
	SortByDate     = "date"
	SortByType     = "type"
	SortByLastRan  = "last_ran"
)

// SortTasks returns the schedule as a slice ordered per sortBy.
// Unrecognised values fall back to SortByDate.
func SortTasks(schedule map[string]*Task, sortBy string) []TaskEntry {
	entries := make([]TaskEntry, 0, len(schedule))
	for id, task := range schedule {
		entries = append(entries, TaskEntry{ID: id, Task: task})
	}

	switch sortBy {
	case SortByType:
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].ActionType() != entries[j].ActionType() {
				return entries[i].ActionType() < entries[j].ActionType()
			}
			return entries[i].DateTime.Before(entries[j].DateTime)
		})
	case SortByLastRan:
		sort.Slice(entries, func(i, j int) bool {
			a, b := entries[i].LastRan, entries[j].LastRan
			switch {
			case a == nil && b == nil:
				return entries[i].DateTime.Before(entries[j].DateTime)
			case a == nil:
				return false // never-ran tasks sort after ones that have run
			case b == nil:
				return true
			default:
				return a.After(*b) // most recently ran first
			}
		})
	default: // SortByDate
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].DateTime.Before(entries[j].DateTime)
		})
	}

	return entries
}
