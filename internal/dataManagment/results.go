package dataManagment

import (
	"encoding/json"
	"log"
	"os"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func ParseRecentTestJSON() (*networkTesting.ICMBTestResult, error) {
	dataExists, path, err := CheckForRecentTestData("/data/output", ".json")
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

	var result *networkTesting.ICMBTestResult
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
