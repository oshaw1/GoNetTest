package pageGeneration

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/scheduler"
)

type PageGenerator struct {
	templates  *template.Template
	repository *dataManagement.Repository
}

var requiredTemplates = []string{
	"base.gohtml",
	"control_quadrant.gohtml",
	"generate_quadrant.gohtml",
	"scheduler_quadrant.gohtml",
	"test_quadrant.gohtml",
}

type DashboardData struct {
	TestData      *TestQuadrantData
	GenerateData  *GenerateQuadrantData
	ControlData   *ControlQuadrantData
	SchedulerData *SchedulerQuadrantData
}

type QuadrantData struct {
	Title string
	Error error
}

type TestQuadrantData struct {
	QuadrantData
	Dates        []string
	TestTypes    []string
	SelectedDate string
	SelectedType string
	TestGroups   []TestGroup
}

type GenerateQuadrantData struct {
	QuadrantData
}

type ControlQuadrantData struct {
	QuadrantData
}

type SchedulerQuadrantData struct {
	QuadrantData
	Schedule map[string]*scheduler.Task
}

type TestGroup struct {
	TimeGroup  string // HRMINSECOND
	JsonPath   string
	TestResult interface{}       // From networkTesting.TestResult based on type
	ChartPaths map[string]string // "speed" -> path.html
}

// NewPageGenerator creates a new instance of PageGenerator with validation
func NewPageGenerator(templatePath string, repo *dataManagement.Repository) (*PageGenerator, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository cannot be nil")
	}

	templates, err := template.ParseGlob(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	pg := &PageGenerator{
		templates:  templates,
		repository: repo,
	}

	if err := pg.validateRequiredTemplates(templates); err != nil {
		return nil, err
	}

	return pg, nil
}

// validateRequiredTemplates ensures all required templates are present
func (pg *PageGenerator) validateRequiredTemplates(templates *template.Template) error {
	for _, name := range requiredTemplates {
		if templates.Lookup(name) == nil {
			return fmt.Errorf("required template not found: %s", name)
		}
	}
	return nil
}

func (g *PageGenerator) GenerateTestQuadrant(selectedDate, selectedType string) (*TestQuadrantData, error) {
	dates, err := g.repository.GetTestDirectories()
	if err != nil {
		return nil, fmt.Errorf("failed to get dates: %w", err)
	}
	if selectedDate == "" && len(dates) > 0 {
		selectedDate = dates[0]
	}
	testTypes, err := g.repository.ListTestTypesInDateDir(selectedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get test types: %w", err)
	}
	var testGroups []TestGroup
	if selectedType != "" {
		fileMap, err := g.repository.MapTestFilesByTimestamp(selectedDate, selectedType)
		if err != nil {
			return nil, err
		}
		for timeGroup, files := range fileMap {
			group := TestGroup{
				TimeGroup:  timeGroup,
				ChartPaths: make(map[string]string),
			}
			for _, file := range files {
				timestampRegex := regexp.MustCompile(`\d{6}`)
				match := timestampRegex.FindString(file)
				if match != "" {
					group.TimeGroup = match
				}
				if strings.HasSuffix(file, ".json") {
					group.JsonPath = file
					// Read JSON as string
					content, err := os.ReadFile(file)
					if err != nil {
						return nil, fmt.Errorf("failed to read JSON file: %w", err)
					}
					group.TestResult = string(content)

					// Extract timestamp using regex
				} else if strings.HasSuffix(file, ".html") {
					chartType := strings.TrimSuffix(strings.Split(filepath.Base(file), "_")[3], ".html")
					group.ChartPaths[chartType] = file
				}
			}
			testGroups = append(testGroups, group)
		}
	}
	return &TestQuadrantData{
		QuadrantData: QuadrantData{Title: "Tests"},
		Dates:        dates,
		TestTypes:    testTypes,
		SelectedDate: selectedDate,
		SelectedType: selectedType,
		TestGroups:   testGroups,
	}, nil
}

func (pg *PageGenerator) RenderDashboard(w http.ResponseWriter) error {
	testData, err := pg.GenerateTestQuadrant("", "")
	if err != nil {
		return err
	}

	generateData, err := pg.GenerateHistoryQuadrant()
	if err != nil {
		return err
	}

	controlData, err := pg.GenerateControlQuadrant()
	if err != nil {
		return err
	}

	schedulerData, err := pg.GenerateSchedulerQuadrant()
	if err != nil {
		return err
	}

	data := &DashboardData{
		TestData:      testData,
		GenerateData:  generateData,
		ControlData:   controlData,
		SchedulerData: schedulerData,
	}

	return pg.templates.ExecuteTemplate(w, "base", data)
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

func (pg *PageGenerator) GenerateHistoryQuadrant() (*GenerateQuadrantData, error) {
	return &GenerateQuadrantData{
		QuadrantData: QuadrantData{Title: "History"},
	}, nil
}

func (pg *PageGenerator) GenerateControlQuadrant() (*ControlQuadrantData, error) {
	return &ControlQuadrantData{
		QuadrantData: QuadrantData{Title: "Control"},
	}, nil
}

func (pg *PageGenerator) GenerateSchedulerQuadrant() (*SchedulerQuadrantData, error) {
	return &SchedulerQuadrantData{
		QuadrantData: QuadrantData{Title: "Scheduler"},
	}, nil
}
