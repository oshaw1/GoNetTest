package dataManagement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-echarts/go-echarts/v2/render"
)

func (r *Repository) SaveTestResult(data interface{}, testType string) (int64, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	res, err := r.db.Exec(
		`INSERT INTO test_results (test_type, timestamp, data) VALUES (?, ?, ?)`,
		testType, time.Now().UTC().Format("2006-01-02 15:04:05"), string(jsonData),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to save test result: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Printf("Test result saved (id=%d, type=%s)", id, testType)
	return id, nil
}

func (r *Repository) SaveChart(chart render.Renderer, testType, chartType string, resultID int64) (string, error) {
	var buf bytes.Buffer
	if err := chart.Render(&buf); err != nil {
		return "", fmt.Errorf("failed to render chart: %w", err)
	}

	var rid interface{}
	if resultID > 0 {
		rid = resultID
	}

	res, err := r.db.Exec(
		`INSERT INTO charts (result_id, test_type, chart_type, timestamp, html_content) VALUES (?, ?, ?, ?, ?)`,
		rid, testType, chartType, time.Now().UTC().Format("2006-01-02 15:04:05"), buf.String(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to save chart: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("/charts/view?id=%d", id)
	log.Printf("Chart saved (id=%d, type=%s/%s)", id, testType, chartType)
	return path, nil
}
