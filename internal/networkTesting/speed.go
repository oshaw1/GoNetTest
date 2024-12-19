package networkTesting

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AverageSpeedTestResult struct {
	Status        string
	AverageMbps   float64
	ElapsedTime   time.Duration
	BytesReceived int64
	TestedURLs    map[string]SpeedTestResult
	Error         error
}

type SpeedTestResult struct {
	Speed    float64
	Status   string
	Duration time.Duration
	Bytes    int64
}

func (t *NetworkTester) MeasureDownloadSpeed() (*AverageSpeedTestResult, error) {
	var totalSpeed float64
	var totalBytes int64
	var totalTime time.Duration
	var successfulTests int
	var lastError error

	testedURLs := make(map[string]SpeedTestResult)

	for _, url := range t.config.Tests.SpeedTestURLs.DownloadURLs {
		result := t.measureSingleDownload(url)

		if result.Error != nil {
			lastError = result.Error
			testedURLs[url] = SpeedTestResult{Status: "FAILED"}
			continue
		}

		totalSpeed += result.AverageMbps
		totalBytes += result.BytesReceived
		totalTime += result.ElapsedTime
		testedURLs[url] = result.TestedURLs[url]
		successfulTests++
	}

	if successfulTests == 0 {
		return &AverageSpeedTestResult{
			TestedURLs: testedURLs,
			Status:     "All tests failed",
			Error:      lastError,
		}, lastError
	}

	avgSpeed := totalSpeed / float64(successfulTests)
	avgTime := totalTime / time.Duration(successfulTests)
	avgBytes := totalBytes / int64(successfulTests)

	return &AverageSpeedTestResult{
		AverageMbps:   avgSpeed,
		ElapsedTime:   avgTime,
		BytesReceived: avgBytes,
		TestedURLs:    testedURLs,
		Status: fmt.Sprintf("Average: %.2f Mbps, Total: %s, Time: %s",
			avgSpeed,
			formatBytes(totalBytes),
			formatDuration(totalTime)),
	}, nil
}

func (t *NetworkTester) measureSingleDownload(url string) AverageSpeedTestResult {
	start := time.Now()

	client, req, err := t.setupDownloadClient(url)
	if err != nil {
		return AverageSpeedTestResult{
			Error: err,
			TestedURLs: map[string]SpeedTestResult{
				url: {Status: "FAILED - Request Creation Error"},
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return AverageSpeedTestResult{
			Error: err,
			TestedURLs: map[string]SpeedTestResult{
				url: {Status: "FAILED - Connection Error"},
			},
		}
	}
	defer resp.Body.Close()

	bytes, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return AverageSpeedTestResult{
			Error: err,
			TestedURLs: map[string]SpeedTestResult{
				url: {
					Status: "FAILED - Download Error",
					Bytes:  bytes,
				},
			},
		}
	}

	elapsed := time.Since(start)
	speedMbps := float64(bytes*8) / (1024 * 1024) / elapsed.Seconds()

	return AverageSpeedTestResult{
		AverageMbps:   speedMbps,
		ElapsedTime:   elapsed,
		BytesReceived: bytes,
		TestedURLs: map[string]SpeedTestResult{
			url: {
				Speed: speedMbps,
				Status: fmt.Sprintf("%.2f Mbps (%s downloaded in %s)",
					speedMbps,
					formatBytes(bytes),
					formatDuration(elapsed)),
				Duration: elapsed,
				Bytes:    bytes,
			},
		},
	}
}

func (t *NetworkTester) MeasureUploadSpeed() (*AverageSpeedTestResult, error) {
	var totalSpeed float64
	var totalBytes int64
	var totalTime time.Duration
	var successfulTests int
	var lastError error

	testedURLs := make(map[string]SpeedTestResult)

	data := make([]byte, 10*1024*1024) // 10MB
	if _, err := rand.Read(data); err != nil {
		return nil, fmt.Errorf("failed to generate test data: %v", err)
	}

	for _, url := range t.config.Tests.SpeedTestURLs.UploadURLs {
		result := t.measureSingleUpload(url, data)

		if result.Error != nil {
			lastError = result.Error
			testedURLs[url] = SpeedTestResult{Status: "FAILED"}
			continue
		}

		totalSpeed += result.AverageMbps
		totalBytes += result.BytesReceived
		totalTime += result.ElapsedTime
		testedURLs[url] = result.TestedURLs[url]
		successfulTests++
	}

	if successfulTests == 0 {
		return &AverageSpeedTestResult{
			TestedURLs: testedURLs,
			Status:     "All upload tests failed",
			Error:      lastError,
		}, lastError
	}

	avgSpeed := totalSpeed / float64(successfulTests)
	avgTime := totalTime / time.Duration(successfulTests)
	avgBytes := totalBytes / int64(successfulTests)

	return &AverageSpeedTestResult{
		AverageMbps:   avgSpeed,
		ElapsedTime:   avgTime,
		BytesReceived: avgBytes,
		TestedURLs:    testedURLs,
		Status: fmt.Sprintf("Average Upload: %.2f Mbps, Total: %s, Time: %s",
			avgSpeed,
			formatBytes(totalBytes),
			formatDuration(totalTime)),
	}, nil
}

func (t *NetworkTester) measureSingleUpload(url string, data []byte) AverageSpeedTestResult {
	start := time.Now()

	client, req, err := t.setupUploadClient(url, data)
	if err != nil {
		return AverageSpeedTestResult{
			Error: err,
			TestedURLs: map[string]SpeedTestResult{
				url: {Status: "FAILED - Request Creation Error"},
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return AverageSpeedTestResult{
			Error: err,
			TestedURLs: map[string]SpeedTestResult{
				url: {Status: "FAILED - Connection Error"},
			},
		}
	}
	defer resp.Body.Close()

	// Read response to ensure upload completed
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return AverageSpeedTestResult{
			Error: err,
			TestedURLs: map[string]SpeedTestResult{
				url: {Status: "FAILED - Response Error"},
			},
		}
	}

	if resp.StatusCode != http.StatusOK {
		return AverageSpeedTestResult{
			Error: fmt.Errorf("upload failed with status: %d", resp.StatusCode),
			TestedURLs: map[string]SpeedTestResult{
				url: {Status: fmt.Sprintf("FAILED - HTTP %d", resp.StatusCode)},
			},
		}
	}

	elapsed := time.Since(start)
	bytes := int64(len(data))
	speedMbps := float64(bytes*8) / (1024 * 1024) / elapsed.Seconds()

	return AverageSpeedTestResult{
		AverageMbps:   speedMbps,
		ElapsedTime:   elapsed,
		BytesReceived: bytes,
		TestedURLs: map[string]SpeedTestResult{
			url: {
				Speed: speedMbps,
				Status: fmt.Sprintf("%.2f Mbps (%s uploaded in %s)",
					speedMbps,
					formatBytes(bytes),
					formatDuration(elapsed)),
				Duration: elapsed,
				Bytes:    bytes,
			},
		},
	}
}

func (t *NetworkTester) setupUploadClient(url string, data []byte) (*http.Client, *http.Request, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = int64(len(data))

	return client, req, nil
}

func (t *NetworkTester) setupDownloadClient(url string) (*http.Client, *http.Request, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")

	return client, req, nil
}
