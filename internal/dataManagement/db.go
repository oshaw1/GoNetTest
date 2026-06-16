package dataManagement

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := initSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS test_results (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			test_type TEXT    NOT NULL,
			timestamp DATETIME NOT NULL,
			data      TEXT    NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_results_type_time
			ON test_results(test_type, timestamp);

		CREATE TABLE IF NOT EXISTS charts (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			result_id    INTEGER REFERENCES test_results(id) ON DELETE CASCADE,
			test_type    TEXT    NOT NULL,
			chart_type   TEXT    NOT NULL,
			timestamp    DATETIME NOT NULL,
			html_content TEXT    NOT NULL,
			source_data  TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_charts_result ON charts(result_id);
		CREATE INDEX IF NOT EXISTS idx_charts_type_time ON charts(test_type, timestamp);

		PRAGMA foreign_keys = ON;
	`)
	if err != nil {
		return err
	}

	// Migration for databases created before source_data existed — the
	// CREATE TABLE IF NOT EXISTS above is a no-op on an existing table, so
	// add the column explicitly, ignoring the error if it's already there.
	if _, err := db.Exec(`ALTER TABLE charts ADD COLUMN source_data TEXT`); err != nil &&
		!strings.Contains(err.Error(), "duplicate column") {
		return fmt.Errorf("failed to migrate charts table: %w", err)
	}

	return nil
}
