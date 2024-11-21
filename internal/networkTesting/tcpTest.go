package networkTesting

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type TCPTestResult struct {
	Host               string
	Timestamp          time.Time
	TestType           string
	PortResults        []*PortResult
	SuccessfulPorts    int
	FailedPorts        int
	MinConnectTime     time.Duration
	MaxConnectTime     time.Duration
	AverageConnectTime time.Duration
	totalConnectTime   time.Duration // unexported, for internal calculation
}

type PortResult struct {
	Port        int
	Connected   bool
	ConnectTime time.Duration
	Error       string
}

type tcpResponse struct {
	port        int
	connectTime time.Duration
	connected   bool
	err         error
}

func (t *NetworkTester) runTCPTest(host string) (*TCPTestResult, error) {
	result := &TCPTestResult{
		Host:      host,
		Timestamp: time.Now(),
		TestType:  "tcp",
	}

	if len(t.config.Tests.TCP.Ports) == 0 {
		return nil, fmt.Errorf("no TCP ports configured for testing")
	}

	responses := make(chan *tcpResponse, len(t.config.Tests.TCP.Ports))
	var wg sync.WaitGroup

	// Test all ports concurrently
	for _, port := range t.config.Tests.TCP.Ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			t.testTCPPort(host, p, responses)
		}(port)
	}

	go func() {
		wg.Wait()
		close(responses)
	}()

	// Process responses
	for resp := range responses {
		portResult := &PortResult{
			Port:        resp.port,
			Connected:   resp.connected,
			ConnectTime: resp.connectTime,
		}

		if resp.err != nil {
			portResult.Error = resp.err.Error()
			result.FailedPorts++
		} else if resp.connected {
			result.SuccessfulPorts++
		}

		t.updateConnectTimeStats(result, resp.connectTime)
		result.PortResults = append(result.PortResults, portResult)
	}

	// Calculate average connect time
	if result.SuccessfulPorts > 0 {
		result.AverageConnectTime = result.totalConnectTime / time.Duration(result.SuccessfulPorts)
	}

	return result, nil
}

func (t *NetworkTester) testTCPPort(host string, port int, responses chan<- *tcpResponse) {
	start := time.Now()
	response := &tcpResponse{
		port: port,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr,
		time.Duration(t.config.Tests.TCP.TimeoutSeconds)*time.Second)

	response.connectTime = time.Since(start)

	if err != nil {
		response.connected = false
		response.err = err
		responses <- response
		return
	}

	response.connected = true
	conn.Close()
	responses <- response
}

func (t *NetworkTester) updateConnectTimeStats(result *TCPTestResult, connectTime time.Duration) {
	if connectTime > 0 {
		if result.MinConnectTime == 0 || connectTime < result.MinConnectTime {
			result.MinConnectTime = connectTime
		}
		if connectTime > result.MaxConnectTime {
			result.MaxConnectTime = connectTime
		}
		result.totalConnectTime += connectTime
	}
}
