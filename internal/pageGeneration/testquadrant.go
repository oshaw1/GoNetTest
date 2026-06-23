package pageGeneration

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
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

	dates, err := g.repository.GetTestDirectories()
	if err != nil {
		return nil, fmt.Errorf("%s failed to get dates: %w", logPrefix, err)
	}

	if selectedDate == "" && len(dates) > 0 {
		selectedDate = dates[0]
		log.Printf("%s: No date selected, defaulting to latest date: %s", logPrefix, selectedDate)
	}

	testTypes, err := g.repository.ListTestTypesInDateDir(selectedDate)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get test types: %w", logPrefix, err)
	}
	log.Printf("%s: Retrieved %d test types for date %s", logPrefix, len(testTypes), selectedDate)

	var testGroups []TestGroup

	if len(testTypes) != 0 {
		if !contains(testTypes, selectedType) {
			selectedType = testTypes[0]
		}
	}

	if selectedType != "" {
		log.Printf("%s: Processing test type: %s", logPrefix, selectedType)

		recordMap, err := g.repository.MapTestsByTimestamp(selectedDate, selectedType)
		if err != nil {
			return nil, fmt.Errorf("%s failed to map tests by timestamp: %w", logPrefix, err)
		}
		log.Printf("%s: Found %d time groups for test type %s", logPrefix, len(recordMap), selectedType)

		for timestamp, record := range recordMap {
			chartIDs := make([]string, len(record.ChartIDs))
			for i, id := range record.ChartIDs {
				chartIDs[i] = strconv.FormatInt(id, 10)
			}

			testGroups = append(testGroups, TestGroup{
				TimeGroup:  timestamp,
				ResultID:   record.ResultID,
				ChartIDs:   strings.Join(chartIDs, ","),
				TestResult: record.TestJSON,
				ChartPaths: record.ChartPaths,
				Historic:   record.Historic,
			})
		}

		sort.Slice(testGroups, func(i, j int) bool {
			return testGroups[i].TimeGroup > testGroups[j].TimeGroup
		})
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

func contains(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func (pg *PageGenerator) RenderTestQuadrant(w http.ResponseWriter, data *TestQuadrantData) error {
	return pg.templates.ExecuteTemplate(w, "test_quadrant", data)
}

func (pg *PageGenerator) RenderTestSelection(w http.ResponseWriter, data *TestQuadrantData) error {
	return pg.templates.ExecuteTemplate(w, "test_selection", data)
}

func (pg *PageGenerator) RenderTestDatesSidebar(w http.ResponseWriter, data *TestQuadrantData) error {
	return pg.templates.ExecuteTemplate(w, "test_dates_sidebar", data)
}

func (pg *PageGenerator) RenderTestResults(w http.ResponseWriter, data *TestQuadrantData) error {
	return pg.templates.ExecuteTemplate(w, "test_results", data)
}
