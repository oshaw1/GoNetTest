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
	Expected int64 // Added to track expected file size
}

func (t *NetworkTester) measureSingleDownload(url string) SpeedTestResult {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return SpeedTestResult{
			Error: err,
			TestedURLs: map[string]TestResult{
				url: {Status: "FAILED - Connection Error"},
			},
		}
	}
	defer resp.Body.Close()

	expectedBytes := resp.ContentLength
	if expectedBytes <= 0 {
		return SpeedTestResult{
			Error: fmt.Errorf("invalid content length"),
			TestedURLs: map[string]TestResult{
				url: {
					Status:   "FAILED - Invalid Content Length",
					Expected: 0,
				},
			},
		}
	}

	bytes, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return SpeedTestResult{
			Error: err,
			TestedURLs: map[string]TestResult{
				url: {
					Status:   "FAILED - Download Error",
					Bytes:    bytes,
					Expected: expectedBytes,
				},
			},
		}
	}

	// Validate we got the expected amount of data
	if bytes < expectedBytes {
		return SpeedTestResult{
			Error: fmt.Errorf("incomplete download"),
			TestedURLs: map[string]TestResult{
				url: {
					Status: fmt.Sprintf("FAILED - Incomplete Download (Got: %.2f MB, Expected: %.2f MB)",
						float64(bytes)/(1024*1024),
						float64(expectedBytes)/(1024*1024)),
					Bytes:    bytes,
					Expected: expectedBytes,
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
				Expected: expectedBytes,
			},
		},
	}
}

func (t *NetworkTester) MeasureDownloadSpeed() (*SpeedTestResult, error) {
	var totalSpeed float64
	var totalBytes int64
	var totalTime time.Duration
	var successfulTests int
	var lastError error

	testedURLs := make(map[string]TestResult)

	for _, url := range t.config.Tests.SpeedTestURLs.URLs {
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
