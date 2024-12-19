package networkTesting

import (
	"testing"

	"github.com/oshaw1/go-net-test/config"
)

func TestMeasureDownloadSpeed(t *testing.T) {
	cfg := &config.Config{
		Tests: config.TestConfigs{
			SpeedTestURLs: config.SpeedTestURLs{
				DownloadURLs: []string{
					"http://ipv4.download.thinkbroadband.com/10MB.zip",
					"http://ipv4.download.thinkbroadband.com/5MB.zip",
				},
			},
		},
	}
	tester := NewNetworkTester(cfg)

	result, err := tester.MeasureDownloadSpeed()
	if err != nil {
		t.Fatalf("MeasureDownloadSpeed returned unexpected error: %v", err)
	}

	if result.AverageMbps <= 0 {
		t.Errorf("Expected average download speed > 0, got %v", result.AverageMbps)
	}

	if result.ElapsedTime <= 0 {
		t.Errorf("Expected elapsed time > 0, got %v", result.ElapsedTime)
	}

	if result.BytesReceived <= 0 {
		t.Errorf("Expected bytes received > 0, got %v", result.BytesReceived)
	}

	expectedURLCount := len(cfg.Tests.SpeedTestURLs.DownloadURLs)
	if len(result.TestedURLs) != expectedURLCount {
		t.Errorf("Expected individual results for %d URLs, got %d", expectedURLCount, len(result.TestedURLs))
	}

	hasSuccess := false
	for _, result := range result.TestedURLs {
		if result.Speed > 0 {
			hasSuccess = true
			break
		}
	}
	if !hasSuccess {
		t.Error("Expected at least one successful speed test")
	}
}
