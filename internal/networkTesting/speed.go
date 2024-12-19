package networkTesting

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type SpeedTestResult struct {
	AverageMbps   float64
	ElapsedTime   time.Duration
	BytesReceived int64
	Error         error
	TestedURLs    map[string]TestResult
}

type TestResult struct {
	Speed    float64
	Status   string
	Duration time.Duration
	Bytes    int64
}

func (t *NetworkTester) MeasureDownloadSpeed() (*SpeedTestResult, error) {
	var totalSpeed float64
	var totalBytes int64
	var totalTime time.Duration
	var successfulTests int
	var lastError error

	testedURLs := make(map[string]TestResult)

	for _, url := range t.config.Tests.SpeedTestURLs.DownloadURLs {
		result := t.measureSingleDownload(url)

		if result.Error != nil {
			lastError = result.Error
			testedURLs[url] = TestResult{Status: "FAILED"}
			continue
		}

		totalSpeed += result.AverageMbps
		totalBytes += result.BytesReceived
		totalTime += result.ElapsedTime
		testedURLs[url] = result.TestedURLs[url]
		successfulTests++
	}

	if successfulTests == 0 {
		return &SpeedTestResult{
			Error:      lastError,
			TestedURLs: testedURLs,
		}, lastError
	}

	return &SpeedTestResult{
		AverageMbps:   totalSpeed / float64(successfulTests),
		ElapsedTime:   totalTime / time.Duration(successfulTests),
		BytesReceived: totalBytes / int64(successfulTests),
		TestedURLs:    testedURLs,
	}, nil
}

func (t *NetworkTester) measureSingleDownload(url string) SpeedTestResult {
	start := time.Now()

	client, req, err := t.setupDownloadClient(url)
	if err != nil {
		return SpeedTestResult{
			Error: err,
			TestedURLs: map[string]TestResult{
				url: {Status: "FAILED - Request Creation Error"},
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return SpeedTestResult{
			Error: err,
			TestedURLs: map[string]TestResult{
				url: {Status: "FAILED - Connection Error"},
			},
		}
	}
	defer resp.Body.Close()

	bytes, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return SpeedTestResult{
			Error: err,
			TestedURLs: map[string]TestResult{
				url: {
					Status: "FAILED - Download Error",
					Bytes:  bytes,
				},
			},
		}
	}

	elapsed := time.Since(start)
	speedMbps := float64(bytes*8) / (1024 * 1024) / elapsed.Seconds()

	return SpeedTestResult{
		AverageMbps:   speedMbps,
		ElapsedTime:   elapsed,
		BytesReceived: bytes,
		TestedURLs: map[string]TestResult{
			url: {
				Speed:    speedMbps,
				Status:   fmt.Sprintf("%.2f Mbps (%.2f MB downloaded)", speedMbps, float64(bytes)/(1024*1024)),
				Duration: elapsed,
				Bytes:    bytes,
			},
		},
	}
}

func (t *NetworkTester) setupDownloadClient(url string) (*http.Client, *http.Request, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add browser-like headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")

	return client, req, nil
}
