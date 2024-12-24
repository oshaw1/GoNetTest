package networkTesting

import (
	"context"
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
	Duration    time.Duration
	Failed      bool
}

type ConnectionResult struct {
	ID        int
	BytesRecv int64
	Speed     float64
	Duration  time.Duration
	Error     error
}

func (t *NetworkTester) RunBandwidthTest() (*BandwidthTestResult, error) {
	result := &BandwidthTestResult{
		StartTime: time.Now(),
		Steps:     make([]ConnectionStep, 0),
	}

	for conns := t.config.Tests.Bandwidth.InitialConnections; conns <= t.config.Tests.Bandwidth.MaxConnections; conns += t.config.Tests.Bandwidth.StepSize {
		step := t.runConnectionStep(conns)
		result.Steps = append(result.Steps, step)
		result.TotalData += step.TotalBytes

		if step.AvgSpeed > result.MaxThroughput {
			result.MaxThroughput = step.AvgSpeed
			result.OptimalConns = conns
		}

		if t.shouldStopTest(step, result) {
			result.FailurePoint = conns
			break
		}
	}

	result.EndTime = time.Now()
	return result, nil
}

func (t *NetworkTester) shouldStopTest(step ConnectionStep, result *BandwidthTestResult) bool {
	if step.Failed {
		return true
	}

	if len(result.Steps) <= 1 {
		return false
	}

	prevSpeed := result.Steps[len(result.Steps)-2].AvgSpeed
	dropPct := (prevSpeed - step.AvgSpeed) / prevSpeed * 100
	return dropPct > t.config.Tests.Bandwidth.FailThreshold
}

func (t *NetworkTester) runConnectionStep(conns int) ConnectionStep {
	start := time.Now()
	step := ConnectionStep{
		Connections: conns,
		ConnResults: make([]ConnectionResult, 0),
	}

	results := t.runDownloads(conns)
	for result := range results {
		if t.isValidResult(result) {
			step.ConnResults = append(step.ConnResults, result)
			step.TotalBytes += result.BytesRecv
		}
	}

	step.Duration = time.Since(start)
	step.AvgSpeed = calculateMbps(step.TotalBytes, step.Duration)
	return step
}

func (t *NetworkTester) runDownloads(conns int) chan ConnectionResult {
	var wg sync.WaitGroup
	results := make(chan ConnectionResult, conns)

	for i := 0; i < conns; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			start := time.Now()
			downloaded, err := t.downloadWithProgress(t.config.Tests.SpeedTestURLs.DownloadURLs[0])
			results <- ConnectionResult{
				ID:        id,
				BytesRecv: downloaded,
				Duration:  time.Since(start),
				Speed:     calculateMbps(downloaded, time.Since(start)),
				Error:     err,
			}
		}(i)
	}

	go t.waitForDownloads(&wg, results)
	return results
}

func (t *NetworkTester) waitForDownloads(wg *sync.WaitGroup, results chan ConnectionResult) {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	duration := t.getStepDuration()
	select {
	case <-done:
	case <-time.After(duration):
	}
	close(results)
}

func (t *NetworkTester) downloadWithProgress(url string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var received int64
	buf := make([]byte, 32*1024)

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

func (t *NetworkTester) isValidResult(r ConnectionResult) bool {
	return r.BytesRecv > 0 && r.Duration > 0 && r.Speed > 0
}

func (t *NetworkTester) getStepDuration() time.Duration {
	duration, err := time.ParseDuration(t.config.Tests.Bandwidth.StepDuration)
	if err != nil {
		return 30 * time.Second
	}
	return duration
}

func calculateMbps(bytes int64, duration time.Duration) float64 {
	return float64(bytes*8) / (1024 * 1024) / duration.Seconds()
}
