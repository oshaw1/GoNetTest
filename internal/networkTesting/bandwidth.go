package networkTesting

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type BandwidthTestResult struct {
	StartTime     time.Time
	EndTime       time.Time
	Steps         []ConnectionStep
	OptimalConns  int
	MaxThroughput float64
	FailurePoint  int
	TotalData     int64
}

type ConnectionStep struct {
	Connections int
	ConnResults []ConnectionResult
	TotalBytes  int64
	AvgSpeed    float64
	Duration    time.Duration `json:"duration,string"`
	Failed      bool
}

type ConnectionResult struct {
	ID        int
	BytesRecv int64
	Duration  time.Duration
	Speed     float64 // Store individual connection speed in Mbps
	Error     error
}

func (t *NetworkTester) RunBandwidthTest() (*BandwidthTestResult, error) {
	fmt.Printf("Starting bandwidth test\n")
	result := &BandwidthTestResult{
		StartTime: time.Now(),
		Steps:     make([]ConnectionStep, 0),
	}

	testURL := t.config.Tests.Bandwidth.DownloadURL
	if testURL == "" {
		return nil, fmt.Errorf("no download URL configured")
	}

	for users := t.config.Tests.Bandwidth.InitialConnections; users <= t.config.Tests.Bandwidth.MaxConnections; users += t.config.Tests.Bandwidth.StepSize {
		totalTestMB := users * 100 //100MB per user

		fmt.Printf("\nTesting bandwidth with %d concurrent users (%dMB total)\n", users, totalTestMB)
		step := t.runStep(users, testURL)
		result.Steps = append(result.Steps, step)
		result.TotalData += step.TotalBytes

		if step.AvgSpeed > result.MaxThroughput {
			result.MaxThroughput = step.AvgSpeed
			result.OptimalConns = users
			fmt.Printf("New max user bandwidth: %.2f Mbps (%.2f MB/s) with %d users\n",
				step.AvgSpeed,
				step.AvgSpeed/8,
				users)
		}

		if t.shouldStopTest(step, result) {
			result.FailurePoint = users
			fmt.Printf("Bandwidth degradation detected at %d users (%dMB)\n", users, totalTestMB)
			break
		}
	}

	result.EndTime = time.Now()
	fmt.Printf("\nBandwidth test completed in %s\n", result.EndTime.Sub(result.StartTime).Round(time.Second))
	fmt.Printf("Optimal user count: %d (%dMB)\n",
		result.OptimalConns,
		result.OptimalConns*100)
	fmt.Printf("Maximum user bandwidth: %.2f Mbps (%.2f MB/s)\n",
		result.MaxThroughput,
		result.MaxThroughput/8)
	return result, nil
}

func (t *NetworkTester) runStep(users int, url string) ConnectionStep {
	stepStart := time.Now()
	step := ConnectionStep{
		Connections: users,
		ConnResults: make([]ConnectionResult, 0),
	}

	var wg sync.WaitGroup
	results := make(chan ConnectionResult, users)

	for i := 0; i < users; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			connStart := time.Now()
			downloaded, err := t.downloadWithProgress(url)
			duration := time.Since(connStart)
			speed := calculateMbps(downloaded, duration)

			results <- ConnectionResult{
				ID:        id,
				BytesRecv: downloaded,
				Duration:  duration,
				Speed:     speed,
				Error:     err,
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	successfulUsers := 0
	for result := range results {
		if result.BytesRecv > 0 {
			step.ConnResults = append(step.ConnResults, result)
			step.TotalBytes += result.BytesRecv
			successfulUsers++

			fmt.Printf("User %d bandwidth: %.2f Mbps (duration: %s)\n",
				result.ID,
				result.Speed,
				result.Duration.Round(time.Millisecond))
		} else if result.Error != nil {
			fmt.Printf("User %d error: %v\n", result.ID, result.Error)
		}
	}

	step.Duration = time.Since(stepStart).Round(time.Millisecond)

	if successfulUsers > 0 {
		step.AvgSpeed = calculateMbps(step.TotalBytes/int64(successfulUsers), step.Duration)
		fmt.Printf("\nStep complete: %d/%d users successful\n", successfulUsers, users)
		fmt.Printf("Average bandwidth per user: %.2f Mbps (%.2f MB/s)\n",
			step.AvgSpeed,
			step.AvgSpeed/8)
	}

	return step
}
func (t *NetworkTester) shouldStopTest(step ConnectionStep, result *BandwidthTestResult) bool {
	if len(result.Steps) <= 1 {
		return false
	}

	//compare against the maximum throughput
	if result.MaxThroughput == 0 || step.AvgSpeed == 0 {
		return false
	}

	dropPct := (result.MaxThroughput - step.AvgSpeed) / result.MaxThroughput * 100
	fmt.Printf("Bandwidth change: %.1f%% from peak\n", -dropPct)

	fmt.Printf("MaxThroughput: %.2f, Current: %.2f, Drop%%: %.1f, Threshold: %.1f, Will Stop: %v\n",
		result.MaxThroughput,
		step.AvgSpeed,
		dropPct,
		t.config.Tests.Bandwidth.FailThreshold,
		dropPct >= t.config.Tests.Bandwidth.FailThreshold)

	return dropPct >= t.config.Tests.Bandwidth.FailThreshold
}

func (t *NetworkTester) downloadWithProgress(url string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, req, err := t.setupClient(url, nil, "GET")
	if err != nil {
		return 0, fmt.Errorf("failed to setup request: %v", err)
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var received int64
	buf := make([]byte, 256*1024)

	for {
		if ctx.Err() != nil {
			return received, ctx.Err()
		}

		n, err := resp.Body.Read(buf)
		if n > 0 {
			atomic.AddInt64(&received, int64(n))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return received, err
		}
	}

	return received, nil
}

func calculateMbps(bytes int64, duration time.Duration) float64 {
	return float64(bytes*8) / (1024 * 1024) / duration.Seconds()
}
