package dataManagement

import (
	"testing"
	"time"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTestDirectories(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.SaveTestResult(&networkTesting.ICMPTestResult{}, "icmp")
	require.NoError(t, err)

	dates, err := repo.GetTestDirectories()
	require.NoError(t, err)
	require.Len(t, dates, 1)
	assert.Equal(t, time.Now().UTC().Format(dateFormat), dates[0])
}

func TestListTestTypesInDateDir(t *testing.T) {
	repo := newTestRepo(t)

	today := time.Now().UTC().Format(dateFormat)

	_, err := repo.SaveTestResult(&networkTesting.ICMPTestResult{}, "icmp")
	require.NoError(t, err)
	_, err = repo.SaveTestResult(&networkTesting.LatencyTestResult{}, "latency")
	require.NoError(t, err)

	types, err := repo.ListTestTypesInDateDir(today)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"icmp", "latency"}, types)
}

func TestDeleteByDate(t *testing.T) {
	repo := newTestRepo(t)

	today := time.Now().UTC().Format(dateFormat)

	_, err := repo.SaveTestResult(&networkTesting.ICMPTestResult{}, "icmp")
	require.NoError(t, err)

	err = repo.DeleteByDate(today)
	assert.NoError(t, err)

	dates, err := repo.GetTestDirectories()
	require.NoError(t, err)
	assert.Empty(t, dates)

	// Deleting a date with no data returns an error
	err = repo.DeleteByDate("2000-01-01")
	assert.Error(t, err)
}
