package dataManagement

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

// TestRecord holds a test result and its associated chart URLs for a single test run.
type TestRecord struct {
	TestJSON   string
	ChartPaths map[string]string // chart_type -> "/charts/view?id=X"
}

func (r *Repository) GetTestDataInRange(startDate, endDate time.Time, testType string) ([]*networkTesting.TestResult, error) {
	log.Printf("GetTestDataInRange: type=%s start=%s end=%s", testType, startDate.Format(dateFormat), endDate.Format(dateFormat))

	rows, err := r.db.Query(`
		SELECT data FROM test_results
		WHERE test_type = ? AND strftime('%Y-%m-%d', timestamp) >= ? AND strftime('%Y-%m-%d', timestamp) <= ?
		ORDER BY timestamp DESC
	`, testType, startDate.Format(dateFormat), endDate.Format(dateFormat))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*networkTesting.TestResult
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		result, err := unmarshalTestResult([]byte(data), testType)
		if err != nil {
			log.Printf("skipping malformed result: %v", err)
			continue
		}
		results = append(results, result)
	}
	return results, rows.Err()
}

func (r *Repository) GetTestData(date, testType string) (*networkTesting.TestResult, error) {
	log.Printf("GetTestData: type=%s date=%s", testType, date)

	if _, err := time.Parse(dateFormat, date); err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	var data string
	err := r.db.QueryRow(`
		SELECT data FROM test_results
		WHERE test_type = ? AND strftime('%Y-%m-%d', timestamp) = ?
		ORDER BY timestamp DESC LIMIT 1
	`, testType, date).Scan(&data)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return unmarshalTestResult([]byte(data), testType)
}

func (r *Repository) GetChart(date, testType string) (bool, string, error) {
	if _, err := time.Parse(dateFormat, date); err != nil {
		return false, "", fmt.Errorf("invalid date format: %w", err)
	}

	var id int64
	err := r.db.QueryRow(`
		SELECT id FROM charts
		WHERE test_type = ? AND strftime('%Y-%m-%d', timestamp) = ?
		ORDER BY timestamp DESC LIMIT 1
	`, testType, date).Scan(&id)

	if err == sql.ErrNoRows {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return true, fmt.Sprintf("/charts/view?id=%d", id), nil
}

func (r *Repository) GetChartInRange(startDate, endDate time.Time, testType string) (bool, string, error) {
	var id int64
	err := r.db.QueryRow(`
		SELECT id FROM charts
		WHERE test_type = ? AND strftime('%Y-%m-%d', timestamp) >= ? AND strftime('%Y-%m-%d', timestamp) <= ?
		ORDER BY timestamp DESC LIMIT 1
	`, testType, startDate.Format(dateFormat), endDate.Format(dateFormat)).Scan(&id)

	if err == sql.ErrNoRows {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return true, fmt.Sprintf("/charts/view?id=%d", id), nil
}

// MapTestsByTimestamp returns test results grouped by timestamp for a given date and type.
func (r *Repository) MapTestsByTimestamp(date, testType string) (map[string]*TestRecord, error) {
	rows, err := r.db.Query(`
		SELECT tr.id, strftime('%H%M%S', tr.timestamp) AS ts_key, tr.data,
		       c.id AS chart_id, c.chart_type
		FROM test_results tr
		LEFT JOIN charts c ON c.result_id = tr.id
		WHERE tr.test_type = ? AND DATE(tr.timestamp) = ?
		ORDER BY tr.timestamp DESC
	`, testType, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make(map[string]*TestRecord)

	for rows.Next() {
		var resultID int64
		var tsKey, data string
		var chartID sql.NullInt64
		var chartType sql.NullString

		if err := rows.Scan(&resultID, &tsKey, &data, &chartID, &chartType); err != nil {
			return nil, err
		}

		if _, exists := records[tsKey]; !exists {
			records[tsKey] = &TestRecord{
				TestJSON:   data,
				ChartPaths: make(map[string]string),
			}
		}

		if chartID.Valid && chartType.Valid {
			records[tsKey].ChartPaths[chartType.String] = fmt.Sprintf("/charts/view?id=%d", chartID.Int64)
		}
	}

	return records, rows.Err()
}

func (r *Repository) GetChartHTML(id int64) (string, error) {
	var html string
	err := r.db.QueryRow(`SELECT html_content FROM charts WHERE id = ?`, id).Scan(&html)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("chart not found")
	}
	return html, err
}

func unmarshalTestResult(data []byte, testType string) (*networkTesting.TestResult, error) {
	result := &networkTesting.TestResult{}
	switch testType {
	case "icmp":
		var v networkTesting.ICMPTestResult
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ICMP JSON: %w", err)
		}
		result.ICMP = &v
	case "download":
		var v networkTesting.AverageSpeedTestResult
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal download JSON: %w", err)
		}
		result.Download = &v
	case "upload":
		var v networkTesting.AverageSpeedTestResult
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal upload JSON: %w", err)
		}
		result.Upload = &v
	case "latency":
		var v networkTesting.LatencyTestResult
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal latency JSON: %w", err)
		}
		result.Latency = &v
	case "bandwidth":
		var v networkTesting.BandwidthTestResult
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal bandwidth JSON: %w", err)
		}
		result.Bandwidth = &v
	case "route":
		var v networkTesting.RouteTestResult
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal route JSON: %w", err)
		}
		result.Route = &v
	default:
		return nil, fmt.Errorf("unsupported test type: %s", testType)
	}
	return result, nil
}
