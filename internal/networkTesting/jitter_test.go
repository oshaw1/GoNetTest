package networkTesting

import (
	"testing"

	"github.com/oshaw1/go-net-test/config"
)

func TestJitterTest(t *testing.T) {
	cfg := &config.Config{
		Tests: config.TestConfigs{
			JitterTest: config.JitterConfig{
				Target:      "8.8.8.8",
				PacketCount: 3,
			},
			ICMP: config.ICMPConfig{
				TimeoutSeconds: 1,
			},
		},
	}
	tester := NewNetworkTester(cfg)

	result, err := tester.RunJitterTest()
	if err != nil {
		t.Fatalf("runJitterTest returned unexpected error: %v", err)
	}

	if result.PacketLoss < 0 || result.PacketLoss > 100 {
		t.Errorf("Expected packet loss between 0-100%%, got %v", result.PacketLoss)
	}

	if len(result.RTTs) == 0 {
		t.Error("Expected at least one RTT measurement")
	}

	if result.AvgJitter < 0 {
		t.Errorf("Expected average jitter >= 0, got %v", result.AvgJitter)
	}

	if result.MaxJitter < result.MinJitter {
		t.Errorf("Max jitter %v less than min jitter %v", result.MaxJitter, result.MinJitter)
	}

	if result.Status != "SUCCESS" && result.Status != "FAILED" {
		t.Errorf("Invalid status: %v", result.Status)
	}
}
