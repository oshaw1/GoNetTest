package dataManagement

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestSetup struct {
	baseDir string
	repo    *Repository
}

func setupTest(t *testing.T) *TestSetup {
	// Create a temporary directory for test data
	baseDir, err := os.MkdirTemp("", "test-data-*")
	require.NoError(t, err)

	repo := &Repository{
		baseDir: baseDir,
	}

	return &TestSetup{
		baseDir: baseDir,
		repo:    repo,
	}
}

func teardownTest(setup *TestSetup) {
	os.RemoveAll(setup.baseDir)
}

func createTestData(t *testing.T, setup *TestSetup, date string, testType string, data interface{}) string {
	datePath := filepath.Join(setup.baseDir, date)
	testPath := filepath.Join(datePath, testType)
	require.NoError(t, os.MkdirAll(testPath, 0755))

	filePath := filepath.Join(testPath, "test.json")
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filePath, jsonData, 0644))

	return filePath
}

func TestGetTestData(t *testing.T) {
	setup := setupTest(t)
	defer teardownTest(setup)

	tests := []struct {
		name        string
		date        string
		testType    string
		setupData   func() string
		wantErr     bool
		validateRes func(*networkTesting.TestResult) bool
	}{
		{
			name:     "valid ICMP data",
			date:     "2024-03-25",
			testType: "icmp",
			setupData: func() string {
				icmpData := &networkTesting.ICMPTestResult{
					Timestamp: time.Now(),
					AvgRTT:    20,
				}
				return createTestData(t, setup, "2024-03-25", "icmp", icmpData)
			},
			validateRes: func(res *networkTesting.TestResult) bool {
				return res.ICMP != nil && res.ICMP.AvgRTT == 20
			},
		},
		{
			name:     "valid download data",
			date:     "2024-03-25",
			testType: "download",
			setupData: func() string {
				downloadData := &networkTesting.AverageSpeedTestResult{
					Timestamp:   time.Now(),
					AverageMbps: 100.5,
				}
				return createTestData(t, setup, "2024-03-25", "download", downloadData)
			},
			validateRes: func(res *networkTesting.TestResult) bool {
				return res.Download != nil && res.Download.AverageMbps == 100.5
			},
		},
		{
			name:     "invalid date format",
			date:     "invalid-date",
			testType: "icmp",
			wantErr:  true,
		},
		{
			name:     "non-existent date",
			date:     "2024-03-26",
			testType: "icmp",
			validateRes: func(res *networkTesting.TestResult) bool {
				return res == nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupData != nil {
				tt.setupData()
			}

			result, err := setup.repo.GetTestData(tt.date, tt.testType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.validateRes != nil {
				assert.True(t, tt.validateRes(result))
			}
		})
	}
}

func TestGetTestDataInRange(t *testing.T) {
	setup := setupTest(t)
	defer teardownTest(setup)

	// Create test data for a range of dates
	startDate := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC)

	// Setup test data
	dates := []string{"2024-03-20", "2024-03-22", "2024-03-25"}
	for _, date := range dates {
		icmpData := &networkTesting.ICMPTestResult{
			Timestamp: time.Now(),
			AvgRTT:    20,
		}
		createTestData(t, setup, date, "icmp", icmpData)
	}

	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		testType  string
		wantLen   int
		wantErr   bool
	}{
		{
			name:      "valid date range",
			startDate: startDate,
			endDate:   endDate,
			testType:  "icmp",
			wantLen:   3,
		},
		{
			name:      "invalid date range",
			startDate: endDate,
			endDate:   startDate,
			testType:  "icmp",
			wantLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := setup.repo.GetTestDataInRange(tt.startDate, tt.endDate, tt.testType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.wantLen)
		})
	}
}

func TestGetChart(t *testing.T) {
	setup := setupTest(t)
	defer teardownTest(setup)

	// Create test chart file
	date := "2024-03-25"
	datePath := filepath.Join(setup.baseDir, date)
	testPath := filepath.Join(datePath, "icmp")
	require.NoError(t, os.MkdirAll(testPath, 0755))
	chartPath := filepath.Join(testPath, "chart.html")
	require.NoError(t, os.WriteFile(chartPath, []byte("<html></html>"), 0644))

	tests := []struct {
		name       string
		date       string
		testType   string
		wantExists bool
		wantErr    bool
	}{
		{
			name:       "existing chart",
			date:       "2024-03-25",
			testType:   "icmp",
			wantExists: true,
		},
		{
			name:       "non-existent chart",
			date:       "2024-03-26",
			testType:   "icmp",
			wantExists: false,
		},
		{
			name:     "invalid date",
			date:     "invalid-date",
			testType: "icmp",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, path, err := setup.repo.GetChart(tt.date, tt.testType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantExists, exists)
			if tt.wantExists {
				assert.NotEmpty(t, path)
			}
		})
	}
}

func TestGetChartInRange(t *testing.T) {
	setup := setupTest(t)
	defer teardownTest(setup)

	// Create test chart files
	dates := []string{"2024-03-20", "2024-03-22", "2024-03-25"}
	for _, date := range dates {
		datePath := filepath.Join(setup.baseDir, date)
		testPath := filepath.Join(datePath, "icmp")
		require.NoError(t, os.MkdirAll(testPath, 0755))
		chartPath := filepath.Join(testPath, "chart.html")
		require.NoError(t, os.WriteFile(chartPath, []byte("<html></html>"), 0644))
	}

	startDate := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		startDate  time.Time
		endDate    time.Time
		testType   string
		wantExists bool
		wantErr    bool
	}{
		{
			name:       "valid date range",
			startDate:  startDate,
			endDate:    endDate,
			testType:   "icmp",
			wantExists: true,
		},
		{
			name:       "invalid date range",
			startDate:  endDate,
			endDate:    startDate,
			testType:   "icmp",
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, path, err := setup.repo.GetChartInRange(tt.startDate, tt.endDate, tt.testType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantExists, exists)
			if tt.wantExists {
				assert.NotEmpty(t, path)
			}
		})
	}
}
