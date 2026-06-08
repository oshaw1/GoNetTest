package dataManagement

import (
	"io"
	"testing"

	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockChart struct {
	renderFunc func(w io.Writer) error
}

func (m MockChart) RenderContent() []byte        { panic("unimplemented") }
func (m MockChart) RenderSnippet() render.ChartSnippet { panic("unimplemented") }
func (m MockChart) Render(w io.Writer) error {
	if m.renderFunc != nil {
		return m.renderFunc(w)
	}
	_, err := io.WriteString(w, "<html><body>Mock Chart</body></html>")
	return err
}

func TestSaveTestResult(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db, nil)

	tests := []struct {
		name     string
		data     interface{}
		testType string
		wantErr  bool
	}{
		{
			name:     "valid test result",
			data:     map[string]interface{}{"speed": 100.5},
			testType: "download",
		},
		{
			name:     "invalid data",
			data:     make(chan int),
			testType: "download",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.SaveTestResult(tt.data, tt.testType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Greater(t, id, int64(0))
		})
	}
}

func TestSaveChart(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db, nil)

	tests := []struct {
		name      string
		chart     MockChart
		testType  string
		chartType string
		wantErr   bool
	}{
		{
			name:      "successful chart save",
			chart:     MockChart{},
			testType:  "latency",
			chartType: "line",
		},
		{
			name: "render error",
			chart: MockChart{
				renderFunc: func(w io.Writer) error {
					return assert.AnError
				},
			},
			testType:  "latency",
			chartType: "line",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := repo.SaveChart(tt.chart, tt.testType, tt.chartType, 0)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Contains(t, path, "/charts/view?id=")
		})
	}
}
