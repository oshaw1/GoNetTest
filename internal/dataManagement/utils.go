package dataManagement

import "fmt"

// GetTestDirectories returns distinct dates that have test results, newest first.
func (r *Repository) GetTestDirectories() ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT strftime('%Y-%m-%d', timestamp) AS d
		FROM test_results
		WHERE timestamp IS NOT NULL
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

// ListTestTypesInDateDir returns distinct test types present for the given date.
func (r *Repository) ListTestTypesInDateDir(date string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT test_type
		FROM test_results
		WHERE strftime('%Y-%m-%d', timestamp) = ?
		ORDER BY test_type
	`, date)
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
