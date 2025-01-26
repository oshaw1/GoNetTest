package pageGeneration

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const logPrefix = "test_quadrant"

type TestQuadrantData struct {
	QuadrantData
	Dates        []string
	TestTypes    []string
	SelectedDate string
	SelectedType string
	TestGroups   []TestGroup
}

func (g *PageGenerator) GenerateTestQuadrant(selectedDate, selectedType string) (*TestQuadrantData, error) {
	log.Printf("%s: Starting generation with selectedDate: %s, selectedType: %s", logPrefix, selectedDate, selectedType)

	// Get Dates
	dates, err := g.repository.GetTestDirectories()
	if err != nil {
		return nil, fmt.Errorf("%s failed to get dates: %w", logPrefix, err)
	}

	// Handle empty date selection
	if selectedDate == "" && len(dates) > 0 {
		selectedDate = dates[0]
		log.Printf("%s: No date selected, defaulting to latest date: %s", logPrefix, selectedDate)
	}

	// Get Test Types
	testTypes, err := g.repository.ListTestTypesInDateDir(selectedDate)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get test types: %w", logPrefix, err)
	}
	log.Printf("%s: Retrieved %d test types for date %s", logPrefix, len(testTypes), selectedDate)

	var testGroups []TestGroup

	// Process Test Type if specified
	if selectedType != "" {
		log.Printf("%s: Processing test type: %s", logPrefix, selectedType)

		// Map files by timestamp
		fileMap, err := g.repository.MapTestFilesByTimestamp(selectedDate, selectedType)
		if err != nil {
			return nil, fmt.Errorf("%s failed to map test files: %v", logPrefix, err)
		}
		log.Printf("%s: Found %d time groups for test type %s", logPrefix, len(fileMap), selectedType)

		// Process each timestamp group
		for timestamp, files := range fileMap {
			group := TestGroup{
				TimeGroup:  timestamp,
				ChartPaths: make(map[string]string),
			}

			log.Printf("%s: Processing time group: %s with %d files", logPrefix, timestamp, len(files))

			if err := processTestFiles(files, &group); err != nil {
				return nil, fmt.Errorf("%s failed to process files for time group %s: %v", logPrefix, timestamp, err)
			}

			testGroups = append(testGroups, group)
		}
	}

	log.Printf("%s: Generated %d test groups", logPrefix, len(testGroups))

	return &TestQuadrantData{
		QuadrantData: QuadrantData{Title: "Tests"},
		Dates:        dates,
		TestTypes:    testTypes,
		SelectedDate: selectedDate,
		SelectedType: selectedType,
		TestGroups:   testGroups,
	}, nil
}

func processJsonFile(filename string, group *TestGroup) error {
	group.JsonPath = filename
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("%s: Failed to read JSON file %s: %v", logPrefix, filename, err)
		return fmt.Errorf("%s: failed to read JSON file: %w", logPrefix, err)
	}
	group.TestResult = string(content)
	log.Printf("%s: Added JSON path: %s", logPrefix, filename)
	return nil
}

func processHtmlFile(filename string, group *TestGroup) error {
	parts := strings.Split(filepath.Base(filename), "_")
	chartIndex := 3 // Skip first 3 parts
	if len(parts) <= chartIndex {
		return fmt.Errorf("invalid HTML filename format: %s", filename)
	}
	chartType := parts[chartIndex]
	group.ChartPaths[chartType] = filename
	return nil
}

func processTestFiles(files []string, group *TestGroup) error {
	for _, file := range files {
		if strings.HasSuffix(file, ".json") {
			if err := processJsonFile(file, group); err != nil {
				return fmt.Errorf("%s: %v", logPrefix, err)
			}
		} else if strings.HasSuffix(file, ".html") {
			if err := processHtmlFile(file, group); err != nil {
				return fmt.Errorf("%s: %v", logPrefix, err)
			}
		}
	}
	return nil
}

func (pg *PageGenerator) RenderTestQuadrant(w http.ResponseWriter, data *TestQuadrantData) error {
	return pg.templates.ExecuteTemplate(w, "test_quadrant", data)
}

func (pg *PageGenerator) RenderTestSelection(w http.ResponseWriter, data *TestQuadrantData) error {
	return pg.templates.ExecuteTemplate(w, "test_selection", data)
}

func (pg *PageGenerator) RenderTestResults(w http.ResponseWriter, data *TestQuadrantData) error {
	return pg.templates.ExecuteTemplate(w, "test_results", data)
}
