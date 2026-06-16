package dataManagement

import "fmt"

// GetTestDirectories returns distinct dates that have test results or
// historic charts (historic charts aren't tied to a test run, so a date
// with only a generated historic chart wouldn't otherwise show up here),
// newest first.
func (r *Repository) GetTestDirectories() ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT d FROM (
			SELECT strftime('%Y-%m-%d', timestamp) AS d FROM test_results WHERE timestamp IS NOT NULL
			UNION
			SELECT strftime('%Y-%m-%d', timestamp) AS d FROM charts WHERE result_id IS NULL
		)
		ORDER BY d DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query test dates: %w", err)
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		dates = append(dates, d)
	}
	return dates, rows.Err()
}

// ListTestTypesInDateDir returns distinct test types present for the given
// date, from either test results or historic charts generated that day.
func (r *Repository) ListTestTypesInDateDir(date string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT t FROM (
			SELECT test_type AS t FROM test_results WHERE strftime('%Y-%m-%d', timestamp) = ?
			UNION
			SELECT test_type AS t FROM charts WHERE result_id IS NULL AND strftime('%Y-%m-%d', timestamp) = ?
		)
		ORDER BY t
	`, date, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query test types: %w", err)
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}
