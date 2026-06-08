package dataManagement

import (
	"testing"
	"time"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRepo(t *testing.T) *Repository {
	t.Helper()
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	return NewRepository(db, nil)
}

func TestGetTestData(t *testing.T) {
	repo := newTestRepo(t)

	icmpData := &networkTesting.ICMPTestResult{AvgRTT: 20}
	id, err := repo.SaveTestResult(icmpData, "icmp")
	require.NoError(t, err)
	require.Greater(t, id, int64(0))

	today := time.Now().UTC().Format(dateFormat)

	tests := []struct {
		name        string
		date        string
		testType    string
		wantErr     bool
		wantNil     bool
		validateRes func(*networkTesting.TestResult) bool
	}{
		{
			name:     "valid ICMP data",
			date:     today,
			testType: "icmp",
			validateRes: func(res *networkTesting.TestResult) bool {
				return res != nil && res.ICMP != nil && res.ICMP.AvgRTT == 20
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
			date:     "2000-01-01",
			testType: "icmp",
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetTestData(tt.date, tt.testType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, result)
				return
			}
			if tt.validateRes != nil {
				assert.True(t, tt.validateRes(result))
			}
		})
	}
}

func TestGetTestDataInRange(t *testing.T) {
	repo := newTestRepo(t)

	icmpData := &networkTesting.ICMPTestResult{AvgRTT: 20}
	_, err := repo.SaveTestResult(icmpData, "icmp")
	require.NoError(t, err)

	start := time.Now().UTC().AddDate(0, 0, -1)
	end := time.Now().UTC().AddDate(0, 0, 1)

	results, err := repo.GetTestDataInRange(start, end, "icmp")
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	// Inverted range returns nothing
	results, err = repo.GetTestDataInRange(end, start, "icmp")
	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetChart(t *testing.T) {
	repo := newTestRepo(t)

	today := time.Now().UTC().Format(dateFormat)

	// Save a test result and a linked chart
	id, err := repo.SaveTestResult(&networkTesting.ICMPTestResult{}, "icmp")
	require.NoError(t, err)
	_, err = repo.SaveChart(MockChart{}, "icmp", "distribution", id)
	require.NoError(t, err)

	exists, path, err := repo.GetChart(today, "icmp")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Contains(t, path, "/charts/view?id=")

	exists, _, err = repo.GetChart("2000-01-01", "icmp")
	assert.NoError(t, err)
	assert.False(t, exists)

	_, _, err = repo.GetChart("invalid-date", "icmp")
	assert.Error(t, err)
}

func TestMapTestsByTimestamp(t *testing.T) {
	repo := newTestRepo(t)

	icmpData := &networkTesting.ICMPTestResult{AvgRTT: 42}
	resultID, err := repo.SaveTestResult(icmpData, "icmp")
	require.NoError(t, err)

	_, err = repo.SaveChart(MockChart{}, "icmp", "distribution", resultID)
	require.NoError(t, err)

	today := time.Now().UTC().Format(dateFormat)
	records, err := repo.MapTestsByTimestamp(today, "icmp")
	require.NoError(t, err)
	require.Len(t, records, 1)

	for _, rec := range records {
		assert.NotEmpty(t, rec.TestJSON)
		assert.Contains(t, rec.ChartPaths, "distribution")
		assert.Contains(t, rec.ChartPaths["distribution"], "/charts/view?id=")
	}
}
