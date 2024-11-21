package dataManagement

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (r *Repository) GetRecentICMPTestResult(test string) (*networkTesting.ICMPTestResult, error) {
	dataExists, path, err := r.CheckRecentData(test, ".json")
	if err != nil {
		return nil, fmt.Errorf("failed to check recent data: %w", err)
	}
	if !dataExists {
		log.Printf("No recent data exists for test: %s", test)
		return nil, nil
	}
	log.Printf("Recent data exists at: %s", path)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var result *networkTesting.ICMPTestResult
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

func ReturnRecentTestDataPath(rootDir string, testType string, fileExtension string) (bool, string, error) {
	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cutoffTime := time.Now().AddDate(0, 0, -conf.RecentDays)
	return walkRootDirectoryForRecent(rootDir, cutoffTime, fileExtension, testType)
}

func walkRootDirectoryForRecent(rootDir string, cutoffTime time.Time, fileExtension string, testType string) (bool, string, error) {
	hasRecentData := false
	var mostRecentPath string
	var mostRecentTime time.Time

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && isDateFolder(info.Name()) {
			recentData, filePath, err := checkDateFolder(path, cutoffTime, fileExtension, testType)
			if err != nil {
				return err
			}
			if !recentData {
				return nil
			}
			hasRecentData = true
			fileTime := info.ModTime()
			if fileTime.After(mostRecentTime) {
				mostRecentTime = fileTime
				mostRecentPath = filePath
			}
		}
		return nil
	})

	if err != nil {
		return false, "", err
	}
	return hasRecentData, mostRecentPath, nil
}

func checkDateFolder(path string, cutoffTime time.Time, fileExtension string, testType string) (bool, string, error) {
	folderDate, err := time.Parse(dateFormat, filepath.Base(path))
	if err != nil {
		return false, "", err
	}
	if folderDate.After(cutoffTime) {
		return checkForFiles(filepath.Join(path, testType), fileExtension)
	}
	return false, "", nil
}

func checkForFiles(path string, fileExtension string) (bool, string, error) {
	var mostRecentPath string
	var mostRecentTime time.Time

	err := filepath.Walk(path, func(subPath string, subInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !subInfo.IsDir() && filepath.Ext(subInfo.Name()) == fileExtension {
			fileTime := subInfo.ModTime()
			if mostRecentPath == "" || fileTime.After(mostRecentTime) {
				mostRecentTime = fileTime
				mostRecentPath = subPath
			}
		}
		return nil
	})

	return mostRecentPath != "", mostRecentPath, err
}

func (r *Repository) GetTestResults(date string, testType string) ([]interface{}, error) {
	log.Printf("GetTestResults called with date: %s, testType: %s", date, testType)
	targetDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	dir := filepath.Join(r.baseDir, targetDate.Format(dateFormat), testType)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist", dir)
		return nil, nil
	}

	var results []interface{}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			log.Printf("Found file: %s", path)
			result, err := r.readJSONFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v", path, err)
				return nil // Continue to next file
			}
			results = append(results, result)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	log.Printf("GetTestResults returned %d results", len(results))
	return results, nil
}

// GetRecentTestResult returns the most recent test result of a given test type
func (r *Repository) GetRecentTestResult(testType string) (interface{}, error) {
	log.Printf("GetRecentTestResult called with testType: %s", testType)
	cutoffTime := time.Now().AddDate(0, 0, -r.config.RecentDays)
	var mostRecentPath string
	var mostRecentTime time.Time

	err := filepath.Walk(r.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && isDateFolder(info.Name()) {
			folderDate, err := time.Parse(dateFormat, info.Name())
			if err == nil && folderDate.After(cutoffTime) {
				testPath := filepath.Join(path, testType)
				if stat, err := os.Stat(testPath); err == nil && stat.IsDir() {
					log.Printf("Found directory: %s", testPath)
					filepath.Walk(testPath, func(filePath string, fileInfo os.FileInfo, err error) error {
						if err != nil {
							return nil
						}
						if !fileInfo.IsDir() && filepath.Ext(filePath) == ".json" {
							log.Printf("Found file: %s", filePath)
							if mostRecentPath == "" || fileInfo.ModTime().After(mostRecentTime) {
								mostRecentPath = filePath
								mostRecentTime = fileInfo.ModTime()
							}
						}
						return nil
					})
				} else {
					log.Printf("Error getting stat for directory: %s, error: %v", testPath, err)
				}
			} else {
				log.Printf("Skipping folder: %s, error: %v", info.Name(), err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error finding recent file: %w", err)
	}

	if mostRecentPath == "" {
		log.Printf("No recent test result found for testType: %s", testType)
		return nil, nil
	}

	log.Printf("Most recent test result found at: %s", mostRecentPath)
	return r.readJSONFile(mostRecentPath)
}

func (r *Repository) CheckRecentData(testType, fileExtension string) (bool, string, error) {
	log.Printf("CheckRecentData called with testType: %s, fileExtension: %s", testType, fileExtension)
	cutoffTime := time.Now().AddDate(0, 0, -r.config.RecentDays)
	var mostRecentPath string
	var mostRecentTime time.Time

	err := filepath.Walk(r.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && isDateFolder(info.Name()) {
			folderDate, err := time.Parse(dateFormat, info.Name())
			if err == nil && folderDate.After(cutoffTime) {
				testPath := filepath.Join(path, testType)
				if stat, err := os.Stat(testPath); err == nil && stat.IsDir() {
					log.Printf("Found directory: %s", testPath)
					filepath.Walk(testPath, func(filePath string, fileInfo os.FileInfo, err error) error {
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
				} else {
					log.Printf("Error getting stat for directory: %s, error: %v", testPath, err)
				}
			} else {
				log.Printf("Skipping folder: %s, error: %v", info.Name(), err)
			}
		}
		return nil
	})

	if err != nil {
		return false, "", fmt.Errorf("error finding recent file: %w", err)
	}

	log.Printf("CheckRecentData returned: %t, %s", mostRecentPath != "", mostRecentPath)
	return mostRecentPath != "", mostRecentPath, nil
}
