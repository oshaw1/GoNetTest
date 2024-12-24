package networkTesting

import (
	"testing"

	"github.com/oshaw1/go-net-test/config"
)

func TestBandwidthTest(t *testing.T) {
	cfg := &config.Config{
		Tests: config.TestConfigs{
			SpeedTestURLs: config.SpeedTestURLs{
				DownloadURLs: []string{
					"http://ipv4.download.thinkbroadband.com/10MB.zip",
				},
			},
			Bandwidth: config.BandwidthConfig{
				InitialConnections: 1,
				MaxConnections:     4,
				StepSize:           1,
				StepDuration:       "30s",
				FailThreshold:      20,
			},
		},
	}
	tester := NewNetworkTester(cfg)

	result, err := tester.RunBandwidthTest()
	if err != nil {
		t.Fatalf("RunBandwidthTest returned error: %v", err)
	}

	if result.MaxThroughput <= 0 {
		t.Error("Expected max throughput > 0")
	}

	if result.OptimalConns == 0 {
		t.Error("Expected optimal connection count > 0")
	}

	if len(result.Steps) == 0 {
		t.Error("Expected at least one test step")
	}

	for _, step := range result.Steps {
		if step.Connections <= 0 {
			t.Error("Invalid connection count in step")
		}
		if step.TotalBytes <= 0 {
			t.Error("Expected bytes transferred > 0")
		}
		for _, conn := range step.ConnResults {
			if conn.BytesRecv <= 0 {
				t.Error("Individual connection transferred no data")
			}
			if conn.Speed <= 0 {
				t.Error("Invalid speed measurement")
			}
		}
	}

	if result.EndTime.Before(result.StartTime) {
		t.Error("Invalid test duration")
	}
}
