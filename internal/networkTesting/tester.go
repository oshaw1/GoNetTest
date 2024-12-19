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
	default:
		err = fmt.Errorf("unsupported test type: %s", testType)
	}

	if err != nil {
		return result, err
	}

	return result, nil
}
