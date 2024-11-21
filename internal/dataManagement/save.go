package dataManagement

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-echarts/go-echarts/v2/render"
)

func (r *Repository) SaveTestResult(data interface{}, testType string) (string, error) {
	now := time.Now()

	dir, err := r.generateFilePath(testType, now)
	if err != nil {
		return "", err
	}

	filename := r.getUniqueFilename(testType, dir, now) + ".json"
	fullPath := filepath.Join(dir, filename)

	dataMap, err := r.convertToMap(data)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	if err := os.WriteFile(fullPath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write data to file: %w", err)
	}

	log.Printf("Data saved successfully to: %s", fullPath)
	return fullPath, nil
}

func (r *Repository) SaveChart(chart render.Renderer, testType string, chartType string) (string, error) {
	now := time.Now()

	dir, err := r.generateFilePath(testType, now)
	if err != nil {
		return "", err
	}

	filename := r.getUniqueFilename(testType, dir, now) + "_" + chartType + ".html"
	fullPath := filepath.Join(dir, filename)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory structure: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	err = chart.Render(f)
	if err != nil {
		return "", fmt.Errorf("failed to render chart: %w", err)
	}

	fmt.Printf("Chart saved successfully to: %s\n", fullPath)
	return fullPath, nil
}

func (r *Repository) generateFilePath(testType string, now time.Time) (string, error) {
	dateFolder := now.Format(dateFormat)
	dir := filepath.Join(r.baseDir, dateFolder, testType)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory structure: %w", err)
	}

	return dir, nil
}

func (r *Repository) getUniqueFilename(testType string, dir string, now time.Time) string {
	base := fmt.Sprintf("%s_test_%s", testType, now.Format("20060102_150405"))

	for i := 1; ; i++ {
		if _, err := os.Stat(filepath.Join(dir, base)); os.IsNotExist(err) {
			return base
		}
		base = fmt.Sprintf("%s_%d.json", base, i)
	}
}
