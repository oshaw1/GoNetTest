package pageGeneration

import (
	"errors"
	"html/template"
	"net/http/httptest"
	"testing"

	"github.com/oshaw1/go-net-test/internal/dataManagement"
)

type MockRepository struct {
	dates     []string
	testTypes []string
	recordMap map[string]*dataManagement.TestRecord
	err       error
}

func (m *MockRepository) GetTestDirectories() ([]string, error) {
	return m.dates, m.err
}

func (m *MockRepository) ListTestTypesInDateDir(date string) ([]string, error) {
	return m.testTypes, m.err
}

func (m *MockRepository) MapTestsByTimestamp(date, testType string) (map[string]*dataManagement.TestRecord, error) {
	return m.recordMap, m.err
}

func TestGenerateTestQuadrant(t *testing.T) {
	tests := []struct {
		name         string
		selectedDate string
		selectedType string
		mockRepo     *MockRepository
		wantErr      bool
	}{
		{
			name:         "error_getting_dates",
			selectedDate: "2024-01-26",
			selectedType: "performance",
			mockRepo: &MockRepository{
				err: errors.New("failed to get dates"),
			},
			wantErr: true,
		},
		{
			name:         "empty_date_selection",
			selectedDate: "",
			selectedType: "performance",
			mockRepo: &MockRepository{
				dates:     []string{"2024-01-26"},
				testTypes: []string{"performance"},
				recordMap: map[string]*dataManagement.TestRecord{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := &PageGenerator{
				repository: tt.mockRepo,
			}

			data, err := generator.GenerateTestQuadrant(tt.selectedDate, tt.selectedType)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTestQuadrant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && data == nil {
				t.Error("GenerateTestQuadrant() returned nil data when no error expected")
			}
		})
	}
}

func TestRenderTemplates(t *testing.T) {
	tests := []struct {
		name     string
		template string
		wantErr  bool
	}{
		{name: "test_quadrant", template: "test_quadrant"},
		{name: "test_selection", template: "test_selection"},
		{name: "test_results", template: "test_results"},
	}

	tmpl := template.Must(template.New("test").Parse(`{{define "test_quadrant"}}{{end}}
		{{define "test_selection"}}{{end}}
		{{define "test_results"}}{{end}}`))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PageGenerator{templates: tmpl}
			w := httptest.NewRecorder()
			data := &TestQuadrantData{}

			var err error
			switch tt.template {
			case "test_quadrant":
				err = pg.RenderTestQuadrant(w, data)
			case "test_selection":
				err = pg.RenderTestSelection(w, data)
			case "test_results":
				err = pg.RenderTestResults(w, data)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Render%s() error = %v, wantErr %v", tt.template, err, tt.wantErr)
			}
		})
	}
}
