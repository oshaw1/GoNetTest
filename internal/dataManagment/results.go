package dataManagment

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func GetRecentICMPTestResult(test string) (*networkTesting.ICMPTestResult, error) {
	dataExists, path, err := ReturnRecentTestDataPath("data/output", test, ".json")
	if err != nil {
		return nil, err
	}
	if !dataExists {
		log.Print("no data exists")
		return nil, err
	}
	log.Print("data exists at" + path)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result *networkTesting.ICMPTestResult
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func GetTestResults(date string, test string) ([]interface{}, error) {
	dataExists, paths, err := ReturnDataPaths(date, test, "data/output", ".json")
	if err != nil {
		return nil, err
	}
	if !dataExists {
		log.Print("No data exists for the specified date")
		return nil, nil
	}

	var results []interface{}
	for _, path := range paths {
		log.Print("data exists at " + path)
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", path, err)
		}

		var result interface{}
		err = json.Unmarshal(content, &result)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling data from %s: %w", path, err)
		}

		results = append(results, result)
	}

	return results, nil
}
