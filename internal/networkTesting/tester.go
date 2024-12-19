package networkTesting

import (
	"context"
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

func (t *NetworkTester) RunTest(ctx context.Context, testTypes []string) ([]any, error) {
	var results []any
	var errors []error

	for _, testType := range testTypes {
		var result any
		var err error

		switch testType {
		case "icmp":
			result, err = t.runICMPTest()
		case "download":
			// implement this
		case "upload":
			// implement this
		default:
			err = fmt.Errorf("unsupported test type: %s", testType)
		}

		if err != nil {
			errors = append(errors, fmt.Errorf("%s test failed: %w", testType, err))
			continue
		}

		results = append(results, result)
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("some tests failed: %v", errors)
	}
	return results, nil
}
