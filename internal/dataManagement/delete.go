package dataManagement

import "fmt"

// DeleteByDate removes all test results (and their associated charts via
// CASCADE) for the given date string (format: 2006-01-02). Historic charts
// have no result_id to cascade from (they aren't tied to a single test
// run), so a date can exist purely because of one of those — they're
// deleted explicitly here too.
func (r *Repository) DeleteByDate(date string) error {
	res, err := r.db.Exec(
		`DELETE FROM test_results WHERE strftime('%Y-%m-%d', timestamp) = ?`, date,
	)
	if err != nil {
		return fmt.Errorf("failed to delete records for date %s: %w", date, err)
	}
	n, _ := res.RowsAffected()

	res2, err := r.db.Exec(
		`DELETE FROM charts WHERE result_id IS NULL AND strftime('%Y-%m-%d', timestamp) = ?`, date,
	)
	if err != nil {
		return fmt.Errorf("failed to delete historic charts for date %s: %w", date, err)
	}
	n2, _ := res2.RowsAffected()

	if n == 0 && n2 == 0 {
		return fmt.Errorf("no data found for date %s", date)
	}

	return nil
}
