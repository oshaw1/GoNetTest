package dataManagement

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"io"

	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockChart struct {
	renderFunc func(w io.Writer) error
}

func (m MockChart) RenderContent() []byte {
	panic("unimplemented")
}

func (m MockChart) RenderSnippet() render.ChartSnippet {
	panic("unimplemented")
}

func (m MockChart) Render(w io.Writer) error {
	if m.renderFunc != nil {
		return m.renderFunc(w)
	}
	_, err := io.WriteString(w, "<html><body>Mock Chart</body></html>")
	return err
}

func TestSaveTestResult(t *testing.T) {
	// Create a temporary directory for test data
	baseDir, err := os.MkdirTemp("", "test-data-*")
	require.NoError(t, err)
	defer os.RemoveAll(baseDir)

	repo := &Repository{
		baseDir: baseDir,
	}

	tests := []struct {
		name     string
		data     interface{}
		testType string
		wantErr  bool
		setup    func()
		validate func(t *testing.T, path string)
	}{
		{
			name: "valid test result",
			data: map[string]interface{}{
				"speed": 100.5,
				"time":  time.Now(),
			},
			testType: "download",
			validate: func(t *testing.T, path string) {
				assert.FileExists(t, path)

				content, err := os.ReadFile(path)
				require.NoError(t, err)

				var result map[string]interface{}
				err = json.Unmarshal(content, &result)
				require.NoError(t, err)
				assert.Contains(t, result, "speed")
			},
		},
		{
			name:     "invalid data",
			data:     make(chan int), // channels can't be marshaled to JSON
			testType: "download",
			wantErr:  true,
		},
		{
			name: "existing file",
			data: map[string]interface{}{
				"speed": 100.5,
			},
			testType: "download",
			setup: func() {
				dir := filepath.Join(baseDir, time.Now().Format(dateFormat), "download")
				_ = os.MkdirAll(dir, 0755)
			},
			validate: func(t *testing.T, path string) {
				assert.FileExists(t, path)
				assert.True(t, strings.Contains(path, "download_test_"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			path, err := repo.SaveTestResult(tt.data, tt.testType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, path)

			if tt.validate != nil {
				tt.validate(t, path)
			}
		})
	}
}

func TestSaveChart(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "test-charts-*")
	require.NoError(t, err)
	defer os.RemoveAll(baseDir)

	repo := &Repository{
		baseDir: baseDir,
	}

	tests := []struct {
		name      string
		chart     MockChart
		testType  string
		chartType string
		wantErr   bool
		validate  func(t *testing.T, path string)
	}{
		{
			name:      "successful chart save",
			chart:     MockChart{},
			testType:  "latency",
			chartType: "line",
			validate: func(t *testing.T, path string) {
				assert.FileExists(t, path)
				content, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Contains(t, string(content), "Mock Chart")
			},
		},
		{
			name: "render error",
			chart: MockChart{
				renderFunc: func(w io.Writer) error {
					return fmt.Errorf("render failed")
				},
			},
			testType:  "latency",
			chartType: "line",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := repo.SaveChart(tt.chart, tt.testType, tt.chartType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, path)

			if tt.validate != nil {
				tt.validate(t, path)
			}
		})
	}
}

func TestGetUniqueFilename(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "test-filenames-*")
	require.NoError(t, err)
	defer os.RemoveAll(baseDir)

	repo := &Repository{
		baseDir: baseDir,
	}

	testTime := time.Date(2024, 3, 25, 10, 30, 0, 0, time.UTC)
	dir := filepath.Join(baseDir, "test-dir")
	require.NoError(t, os.MkdirAll(dir, 0755))

	tests := []struct {
		name     string
		testType string
		setup    func()
		validate func(t *testing.T, filename string)
	}{
		{
			name:     "unique filename",
			testType: "download",
			validate: func(t *testing.T, filename string) {
				assert.Contains(t, filename, "download_test_")
				assert.Contains(t, filename, "20240325_103000")
			},
		},
		{
			name:     "handle existing file",
			testType: "download",
			setup: func() {
				baseFilename := fmt.Sprintf("download_test_%s", testTime.Format("20060102_150405"))
				existingPath := filepath.Join(dir, baseFilename)
				_, _ = os.Create(existingPath)
			},
			validate: func(t *testing.T, filename string) {
				assert.Contains(t, filename, "download_test_")
				assert.True(t, strings.Contains(filename, "_1") || strings.Contains(filename, "_2"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			filename := repo.getUniqueFilename(tt.testType, dir, testTime)
			assert.NotEmpty(t, filename)

			if tt.validate != nil {
				tt.validate(t, filename)
			}
		})
	}
}
