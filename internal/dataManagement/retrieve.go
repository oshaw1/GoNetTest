package dataManagement

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (r *Repository) GetTestDataInRange(startDate time.Time, endDate time.Time, testType string) ([]*networkTesting.TestResult, error) {
	log.Printf("GetTestDataInRange called with startDate: %s, endDate: %s, testType: %s",
		startDate.Format(dateFormat), endDate.Format(dateFormat), testType)

	var allResults []*networkTesting.TestResult

	// Try each date in the range
	for d := endDate; !d.Before(startDate); d = d.AddDate(0, 0, -1) {
		dateStr := d.Format(dateFormat)

		result, err := r.GetTestData(dateStr, testType)
		if err != nil {
			log.Printf("Error getting data for date %s: %v", dateStr, err)
			continue
		}

		if result != nil {
			allResults = append(allResults, result)
		}
	}

	sort.Slice(allResults, func(i, j int) bool {
		var timeI, timeJ time.Time

		switch testType {
		case "icmp":
			if allResults[i].ICMP != nil && allResults[j].ICMP != nil {
				timeI = allResults[i].ICMP.Timestamp
				timeJ = allResults[j].ICMP.Timestamp
			}
		case "download":
			if allResults[i].Download != nil && allResults[j].Download != nil {
				timeI = allResults[i].Download.Timestamp
				timeJ = allResults[j].Download.Timestamp
			}
		case "upload":
			if allResults[i].Upload != nil && allResults[j].Upload != nil {
				timeI = allResults[i].Upload.Timestamp
				timeJ = allResults[j].Upload.Timestamp
			}
		case "latency":
			if allResults[i].Jitter != nil && allResults[j].Jitter != nil {
				timeI = allResults[i].Jitter.Timestamp
				timeJ = allResults[j].Jitter.Timestamp
			}
		}

		return timeI.After(timeJ)
	})

	return allResults, nil
}

func (r *Repository) GetTestData(date string, testType string) (*networkTesting.TestResult, error) {
	log.Printf("GetTestDataOnDate called with date: %s, testType: %s", date, testType)

	exists, filePath, err := r.CheckData(date, testType, ".json")
	if err != nil {
		return nil, fmt.Errorf("failed to check data: %w", err)
	}
	if !exists {
		return nil, nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result := &networkTesting.TestResult{}

	switch testType {
	case "icmp":
		var icmpResult *networkTesting.ICMPTestResult
		if err := json.Unmarshal(content, &icmpResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ICMP JSON: %w", err)
		}
		result.ICMP = icmpResult
	case "download":
		var speedResult *networkTesting.AverageSpeedTestResult
		if err := json.Unmarshal(content, &speedResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Speed JSON: %w", err)
		}
		result.Download = speedResult
	case "upload":
		var speedResult *networkTesting.AverageSpeedTestResult
		if err := json.Unmarshal(content, &speedResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Speed JSON: %w", err)
		}
		result.Upload = speedResult
	case "latency":
		var jitterResult *networkTesting.JitterTestResult
		if err := json.Unmarshal(content, &jitterResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Speed JSON: %w", err)
		}
		result.Jitter = jitterResult
	default:
		return nil, fmt.Errorf("unsupported test type: %s", testType)
	}

	return result, nil
}

func (r *Repository) GetChartInRange(startDate time.Time, endDate time.Time, testType string) (bool, string, error) {
	log.Printf("GetChartInRange called with startDate: %s, endDate: %s, testType: %s",
		startDate.Format(dateFormat), endDate.Format(dateFormat), testType)

	// Try each date from newest to oldest
	for d := endDate; !d.Before(startDate); d = d.AddDate(0, 0, -1) {
		dateStr := d.Format(dateFormat)

		exists, filePath, err := r.GetChart(dateStr, testType)
		if err != nil {
			log.Printf("Error getting chart for date %s: %v", dateStr, err)
			continue
		}

		if exists {
			log.Printf("Found chart for date %s: %s", dateStr, filePath)
			return true, filePath, nil
		}
	}

	log.Printf("No charts found between %s and %s",
		startDate.Format(dateFormat), endDate.Format(dateFormat))
	return false, "", nil
}

func (r *Repository) GetChart(date string, testType string) (bool, string, error) {
	log.Printf("GetChartOnDate called with date: %s, testType: %s", date, testType)

	return r.CheckData(date, testType, ".html")
}

func (r *Repository) isValidDateStr(dateStr string) bool {
	_, err := time.Parse(dateFormat, dateStr)
	return err == nil
}

func (r *Repository) CheckData(date, testType, fileExtension string) (bool, string, error) {
	log.Printf("CheckData called with date: %s, testType: %s, fileExtension: %s",
		date, testType, fileExtension)

	targetDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return false, "", fmt.Errorf("invalid date format: %w", err)
	}

	// Construct the path for the specific date
	datePath := filepath.Join(r.baseDir, targetDate.Format(dateFormat))
	testPath := filepath.Join(datePath, testType)

	// Check if the directory exists
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		log.Printf("Directory does not exist: %s", testPath)
		return false, "", nil
	}

	var mostRecentPath string
	var mostRecentTime time.Time

	err = filepath.Walk(testPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !fileInfo.IsDir() && strings.HasSuffix(filePath, fileExtension) {
			log.Printf("Found file: %s", filePath)
			if mostRecentPath == "" || fileInfo.ModTime().After(mostRecentTime) {
				mostRecentPath = filePath
				mostRecentTime = fileInfo.ModTime()
			}
		}
		return nil
	})

	if err != nil {
		return false, "", fmt.Errorf("error finding file: %w", err)
	}

	log.Printf("CheckDataForDate returned: %t, %s", mostRecentPath != "", mostRecentPath)
	return mostRecentPath != "", mostRecentPath, nil
}
