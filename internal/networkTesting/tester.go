package networkTesting

import (
	"fmt"

	"github.com/oshaw1/go-net-test/config"
)

type NetworkTester struct {
	config *config.Config
}

func NewNetworkTester(config *config.Config) *NetworkTester {
	return &NetworkTester{
		config: config,
	}
}

type TestResult struct {
	ICMP      *ICMPTestResult         `json:"ICMP,omitempty"`
	Download  *AverageSpeedTestResult `json:"Download,omitempty"`
	Upload    *AverageSpeedTestResult `json:"Upload,omitempty"`
	Route     *RouteTestResult        `json:"Route,omitempty"`
	Jitter    *JitterTestResult       `json:"Jitter,omitempty"`
	Bandwidth *BandwidthTestResult    `json:"Bandwidth,omitempty"`
}

func (t *NetworkTester) RunTest(testType string) (any, error) {

	var result any
	var err error

	switch testType {
	case "icmp":
		result, err = t.runICMPTest()
	case "download":
		result, err = t.MeasureDownloadSpeed()
	case "upload":
		result, err = t.MeasureUploadSpeed()
	case "route":
		result, err = t.RunRouteTest()
	case "latency":
		result, err = t.RunJitterTest()
	case "bandwidth":
		result, err = t.RunBandwidthTest()
	default:
		err = fmt.Errorf("unsupported test type: %s", testType)
	}

	if err != nil {
		return result, err
	}

	return result, nil
}
