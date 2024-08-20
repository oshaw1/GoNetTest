package dataManagment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func SaveICMBTestData(data *networkTesting.ICMBTestResult) error {
	now := time.Now()

	dir, err := generateFilePath("data/output/", now, "icmb")
	if err != nil {
		return err
	}

	filename := getUniqueFilename("icmb", dir, now)
	fullPath := filepath.Join(dir, filename)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	fmt.Printf("Data saved successfully to: %s\n", fullPath)
	return nil
}

func CheckForRecentTestData(rootDir string, days int, fileExtension string) (bool, error) {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	return walkRootDirectory(rootDir, cutoffTime, fileExtension)
}

func walkRootDirectory(rootDir string, cutoffTime time.Time, fileExtension string) (bool, error) {
	hasRecentData := false
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && isDateFolder(info.Name()) {
			recentData, err := checkDateFolder(path, cutoffTime, fileExtension)
			if err != nil {
				return err
			}
			if recentData {
				hasRecentData = true
				return filepath.SkipAll
			}
		}
		return nil
	})

	if err != nil && err != filepath.SkipAll {
		return false, err
	}
	return hasRecentData, nil
}

func checkDateFolder(path string, cutoffTime time.Time, fileExtension string) (bool, error) {
	folderDate, err := time.Parse("2006-01-02", filepath.Base(path))
	if err != nil {
		return false, err
	}
	if folderDate.After(cutoffTime) {
		return checkForFiles(path, fileExtension)
	}
	return false, nil
}

func checkForFiles(path string, fileExtension string) (bool, error) {
	hasJSON := false
	err := filepath.Walk(path, func(subPath string, subInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !subInfo.IsDir() && filepath.Ext(subInfo.Name()) == fileExtension {
			hasJSON = true
			return filepath.SkipAll
		}
		return nil
	})
	return hasJSON, err
}

func isDateFolder(name string) bool {
	_, err := time.Parse("2006-01-02", name)
	return err == nil
}

func generateFilePath(baseDir string, now time.Time, folder string) (string, error) {
	dateFolder := now.Format("2006-01-02")
	dir := filepath.Join(baseDir, dateFolder, folder)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory structure: %w", err)
	}

	return dir, nil
}

func getUniqueFilename(typeOfTest string, dir string, now time.Time) string {
	base := fmt.Sprintf(typeOfTest+"_test_%s", now.Format("20060102_150405"))
	filename := base + ".json"

	for i := 1; ; i++ {
		if _, err := os.Stat(filepath.Join(dir, filename)); os.IsNotExist(err) {
			return filename
		}
		filename = fmt.Sprintf("%s_%d.json", base, i)
	}
}
