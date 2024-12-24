package networkTesting

import (
	"testing"

	"github.com/oshaw1/go-net-test/config"
)

func TestRouteTest(t *testing.T) {
	cfg := &config.Config{
		Tests: config.TestConfigs{
			RouteTest: config.RouteConfig{
				Target:  "8.8.8.8",
				MaxHops: 30,
			},
			ICMP: config.ICMPConfig{
				TimeoutSeconds: 1,
			},
		},
	}
	tester := NewNetworkTester(cfg)

	result, err := tester.RunRouteTest()
	if err != nil {
		t.Fatalf("runRouteTest returned unexpected error: %v", err)
	}

	if len(result.Hops) == 0 {
		t.Error("Expected at least one hop")
	}

	if len(result.Hops) > cfg.Tests.RouteTest.MaxHops {
		t.Errorf("Got %d hops, more than max %d", len(result.Hops), cfg.Tests.RouteTest.MaxHops)
	}

	for i, hop := range result.Hops {
		if hop.Number != i+1 {
			t.Errorf("Hop %d has incorrect number %d", i+1, hop.Number)
		}

		if !hop.Lost && hop.RTT <= 0 {
			t.Errorf("Hop %d has invalid RTT: %v", i+1, hop.RTT)
		}

		if !hop.Lost && hop.Address == "" {
			t.Errorf("Hop %d missing address", i+1)
		}
	}

	if result.Status != "SUCCESS" && result.Status != "INCOMPLETE" {
		t.Errorf("Invalid status: %v", result.Status)
	}
}
