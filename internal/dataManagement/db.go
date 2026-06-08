package dataManagement

import (
	"database/sql"
	"fmt"

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
			html_content TEXT    NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_charts_result ON charts(result_id);
		CREATE INDEX IF NOT EXISTS idx_charts_type_time ON charts(test_type, timestamp);

		PRAGMA foreign_keys = ON;
	`)
	return err
}
