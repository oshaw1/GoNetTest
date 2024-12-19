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

type speedTest struct {
	urls    []string
	measure func(string) AverageSpeedTestResult
}

func (t *NetworkTester) MeasureDownloadSpeed() (*AverageSpeedTestResult, error) {
	return t.runSpeedTest(speedTest{
		urls:    t.config.Tests.SpeedTestURLs.DownloadURLs,
		measure: t.measureSingleDownload,
	}, "download")
}

func (t *NetworkTester) MeasureUploadSpeed() (*AverageSpeedTestResult, error) {
	data := make([]byte, 10*1024*1024) // 10MB
	if _, err := rand.Read(data); err != nil {
		return nil, fmt.Errorf("failed to generate test data: %v", err)
	}

	return t.runSpeedTest(speedTest{
		urls: t.config.Tests.SpeedTestURLs.UploadURLs,
		measure: func(url string) AverageSpeedTestResult {
			return t.measureSingleUpload(url, data)
		},
	}, "upload")
}

func (t *NetworkTester) runSpeedTest(test speedTest, testType string) (*AverageSpeedTestResult, error) {
	var totalSpeed float64
	var totalBytes int64
	var totalTime time.Duration
	var successfulTests int
	var lastError error

	testedURLs := make(map[string]SpeedTestResult)

	for _, url := range test.urls {
		result := test.measure(url)

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
			Status:     fmt.Sprintf("All %s tests failed", testType),
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
		Status: fmt.Sprintf("Average %s: %.2f Mbps, Total: %s, Time: %s",
			testType,
			avgSpeed,
			formatBytes(totalBytes),
			formatDuration(totalTime)),
	}, nil
}

func createSpeedTestResult(url string, speedMbps float64, elapsed time.Duration, bytes int64, err error, statusFormat string, args ...interface{}) AverageSpeedTestResult {
	if err != nil {
		return AverageSpeedTestResult{
			Error: err,
			TestedURLs: map[string]SpeedTestResult{
				url: {Status: fmt.Sprintf(statusFormat, args...)},
			},
		}
	}

	return AverageSpeedTestResult{
		AverageMbps:   speedMbps,
		ElapsedTime:   elapsed,
		BytesReceived: bytes,
		TestedURLs: map[string]SpeedTestResult{
			url: {
				Speed:    speedMbps,
				Status:   fmt.Sprintf("%.2f Mbps (%s %s in %s)", speedMbps, formatBytes(bytes), args[0], formatDuration(elapsed)),
				Duration: elapsed,
				Bytes:    bytes,
			},
		},
	}
}

func (t *NetworkTester) measureSingleDownload(url string) AverageSpeedTestResult {
	start := time.Now()

	client, req, err := t.setupClient(url, nil, "GET")
	if err != nil {
		return createSpeedTestResult(url, 0, 0, 0, err, "FAILED - Request Creation Error")
	}

	resp, err := client.Do(req)
	if err != nil {
		return createSpeedTestResult(url, 0, 0, 0, err, "FAILED - Connection Error")
	}
	defer resp.Body.Close()

	bytes, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return createSpeedTestResult(url, 0, 0, bytes, err, "FAILED - Download Error")
	}

	elapsed := time.Since(start)
	speedMbps := float64(bytes*8) / (1024 * 1024) / elapsed.Seconds()

	return createSpeedTestResult(url, speedMbps, elapsed, bytes, nil, "%s",
		fmt.Sprintf("%.2f Mbps (%s downloaded in %s)",
			speedMbps,
			formatBytes(bytes),
			formatDuration(elapsed)))
}

func (t *NetworkTester) measureSingleUpload(url string, data []byte) AverageSpeedTestResult {
	start := time.Now()

	client, req, err := t.setupClient(url, data, "POST")
	if err != nil {
		return createSpeedTestResult(url, 0, 0, 0, err, "FAILED - Request Creation Error")
	}

	resp, err := client.Do(req)
	if err != nil {
		return createSpeedTestResult(url, 0, 0, 0, err, "FAILED - Connection Error")
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return createSpeedTestResult(url, 0, 0, 0, err, "FAILED - Response Error")
	}

	if resp.StatusCode != http.StatusOK {
		return createSpeedTestResult(url, 0, 0, 0,
			fmt.Errorf("upload failed with status: %d", resp.StatusCode),
			"FAILED - HTTP %d", resp.StatusCode)
	}

	elapsed := time.Since(start)
	bytes := int64(len(data))
	speedMbps := float64(bytes*8) / (1024 * 1024) / elapsed.Seconds()

	return createSpeedTestResult(url, speedMbps, elapsed, bytes, nil, "%s",
		fmt.Sprintf("%.2f Mbps (%s uploaded in %s)",
			speedMbps,
			formatBytes(bytes),
			formatDuration(elapsed)))
}

func (t *NetworkTester) setupClient(url string, data []byte, method string) (*http.Client, *http.Request, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	var req *http.Request
	var err error

	if method == "POST" {
		req, err = http.NewRequest(method, url, bytes.NewReader(data))
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("Content-Type", "application/octet-stream")
		req.ContentLength = int64(len(data))
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")
	return client, req, nil
}
