package networkTesting

import (
	"testing"
	"time"

	"github.com/oshaw1/go-net-test/config"
)

func TestBandwidthTest(t *testing.T) {
	cfg := &config.Config{
		Tests: config.TestConfigs{
			Bandwidth: config.BandwidthConfig{
				InitialConnections: 1,
				MaxConnections:     4,
				StepSize:           1,
				DownloadURL:        "http://ipv4.download.thinkbroadband.com/10MB.zip",
				FailThreshold:      50,
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

	if result.OptimalConns <= 0 {
		t.Error("Expected optimal user count > 0")
	}

	if len(result.Steps) == 0 {
		t.Error("Expected at least one test step")
	}

	lastSpeed := float64(0)
	for i, step := range result.Steps {
		t.Run("ValidateStep", func(t *testing.T) {
			// Check user counts
			if step.Connections <= 0 {
				t.Error("Invalid user count in step")
			}

			// Verify data transfer
			if step.TotalBytes <= 0 {
				t.Error("Expected bytes transferred > 0")
			}

			// Check individual user results
			for _, user := range step.ConnResults {
				if user.BytesRecv <= 0 {
					t.Error("User transferred no data")
				}
			}

			// Verify bandwidth measurements
			if step.AvgSpeed <= 0 {
				t.Error("Invalid bandwidth measurement")
			}

			// Check for reasonable bandwidth degradation
			if i > 0 && !step.Failed {
				dropPct := (lastSpeed - step.AvgSpeed) / lastSpeed * 100
				if dropPct > float64(cfg.Tests.Bandwidth.FailThreshold*2) {
					t.Errorf("Excessive bandwidth drop: %.2f%%", dropPct)
				}
			}
			lastSpeed = step.AvgSpeed
		})
	}

	// Validate test timing
	if result.EndTime.Before(result.StartTime) {
		t.Error("Invalid test duration")
	}

	duration := result.EndTime.Sub(result.StartTime)
	if duration < time.Second {
		t.Error("Test completed too quickly")
	}

	// Validate failure detection
	if result.FailurePoint > 0 {
		if result.FailurePoint <= result.OptimalConns {
			t.Error("Failure point should be after optimal connection count")
		}
	}
}
