package networkTesting

import (
	"testing"

	"github.com/oshaw1/go-net-test/config"
)

func TestLatencyTest(t *testing.T) {
	cfg := &config.Config{
		Tests: config.TestConfigs{
			LatencyTest: config.LatencyConfig{
				Target:      "8.8.8.8",
				PacketCount: 3,
			},
			ICMP: config.ICMPConfig{
				TimeoutSeconds: 1,
			},
		},
	}
	tester := NewNetworkTester(cfg)

	result, err := tester.RunLatencyTest()
	if err != nil {
		t.Fatalf("runLatencyTest returned unexpected error: %v", err)
	}

	if result.PacketLoss < 0 || result.PacketLoss > 100 {
		t.Errorf("Expected packet loss between 0-100%%, got %v", result.PacketLoss)
	}

	if len(result.RTTs) == 0 {
		t.Error("Expected at least one RTT measurement")
	}

	if result.AvgLatency < 0 {
		t.Errorf("Expected average Latency >= 0, got %v", result.AvgLatency)
	}

	if result.MaxLatency < result.MinLatency {
		t.Errorf("Max Latency %v less than min Latency %v", result.MaxLatency, result.MinLatency)
	}

	if result.Status != "SUCCESS" && result.Status != "FAILED" {
		t.Errorf("Invalid status: %v", result.Status)
	}
}
