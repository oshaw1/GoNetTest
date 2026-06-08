package dataManagement

import "fmt"

// DeleteByDate removes all test results (and their associated charts via CASCADE)
// for the given date string (format: 2006-01-02).
func (r *Repository) DeleteByDate(date string) error {
	res, err := r.db.Exec(
		`DELETE FROM test_results WHERE strftime('%Y-%m-%d', timestamp) = ?`, date,
	)
	if err != nil {
		return fmt.Errorf("failed to delete records for date %s: %w", date, err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("no data found for date %s", date)
	}

	return nil
}
