package dataManagement

import (
	"fmt"
	"strings"
)

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

// DeleteByID removes a single test result by its test_results.id. Its
// chart is removed automatically via the result_id ON DELETE CASCADE
// foreign key.
func (r *Repository) DeleteByID(id int64) error {
	res, err := r.db.Exec(`DELETE FROM test_results WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete result %d: %w", id, err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("no test result found with id %d", id)
	}

	return nil
}

// DeleteChartsByIDs removes historic charts directly from the charts table.
// Historic charts have no result_id, so there's no test result to delete
// that would cascade to them — they must be removed by chart id instead.
func (r *Repository) DeleteChartsByIDs(ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("no chart ids provided")
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	res, err := r.db.Exec(
		fmt.Sprintf(`DELETE FROM charts WHERE id IN (%s)`, strings.Join(placeholders, ",")),
		args...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete charts: %w", err)
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("no charts found for given ids")
	}

	return nil
}
