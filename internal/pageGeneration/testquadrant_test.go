package pageGeneration

import (
	"errors"
	"html/template"
	"net/http/httptest"
	"testing"
)

// MockRepository implements required methods for testing
type MockRepository struct {
	dates     []string
	testTypes []string
	fileMap   map[string][]string
	err       error
}

func (m *MockRepository) GetTestDirectories() ([]string, error) {
	return m.dates, m.err
}

func (m *MockRepository) ListTestTypesInDateDir(date string) ([]string, error) {
	return m.testTypes, m.err
}

func (m *MockRepository) MapTestFilesByTimestamp(date, testType string) (map[string][]string, error) {
	return m.fileMap, m.err
}

func TestGenerateTestQuadrant(t *testing.T) {
	tests := []struct {
		name         string
		selectedDate string
		selectedType string
		mockRepo     *MockRepository
		wantErr      bool
		wantData     *TestQuadrantData
	}{
		// {
		// 	name:         "successful_generation",
		// 	selectedDate: "2024-01-26",
		// 	selectedType: "performance",
		// 	mockRepo: &MockRepository{
		// 		dates:     []string{"2024-01-26", "2024-01-25"},
		// 		testTypes: []string{"performance", "unit"},
		// 		fileMap: map[string][]string{
		// 			"timestamp1": {
		// 				filepath.Join("testdata", "test.json"),
		// 				filepath.Join("testdata", "test_chart_type1_123456_1234.html"),
		// 			},
		// 		},
		// 	},
		// 	wantErr: false,
		// 	wantData: &TestQuadrantData{
		// 		QuadrantData: QuadrantData{Title: "Tests"},
		// 		Dates:        []string{"2024-01-26", "2024-01-25"},
		// 		TestTypes:    []string{"performance", "unit"},
		// 		SelectedDate: "2024-01-26",
		// 		SelectedType: "performance",
		// 	},
		// },
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
			},
			wantErr: false,
			wantData: &TestQuadrantData{
				SelectedDate: "2024-01-26",
				Dates:        []string{"2024-01-26"},
				TestTypes:    []string{"performance"},
			},
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

func TestProcessTestFiles(t *testing.T) {
	tests := []struct {
		name    string
		files   []string
		wantErr bool
	}{
		// {
		// 	name: "process_valid_files",
		// 	files: []string{
		// 		filepath.Join("testdata", "test.json"),
		// 		filepath.Join("testdata", "test_chart_type1_123456_1234.html"),
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "process_invalid_json",
			files: []string{
				"nonexistent.json",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &TestGroup{
				ChartPaths: make(map[string]string),
			}

			err := processTestFiles(tt.files, group)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTestFiles() error = %v, wantErr %v", err, tt.wantErr)
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
		{
			name:     "test_quadrant",
			template: "test_quadrant",
			wantErr:  false,
		},
		{
			name:     "test_selection",
			template: "test_selection",
			wantErr:  false,
		},
		{
			name:     "test_results",
			template: "test_results",
			wantErr:  false,
		},
	}

	tmpl := template.Must(template.New("test").Parse(`{{define "test_quadrant"}}{{end}}
		{{define "test_selection"}}{{end}}
		{{define "test_results"}}{{end}}`))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PageGenerator{
				templates: tmpl,
			}

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
