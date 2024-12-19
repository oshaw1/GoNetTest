package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oshaw1/go-net-test/internal/charting"
	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

const dateFormat = "2006-01-02"

// NetworkTestHandler handles network test HTTP requests
type NetworkTestHandler struct {
	tester     *networkTesting.NetworkTester
	repository *dataManagement.Repository
	charts     *charting.Generator // Use concrete type instead of interface
}

func NewNetworkTestHandler(tester *networkTesting.NetworkTester, repo *dataManagement.Repository) *NetworkTestHandler {
	return &NetworkTestHandler{
		tester:     tester,
		repository: repo,
		charts:     charting.NewGenerator(), // Create it directly
	}
}

// HandleNetworkTest handles the test execution request
func (h *NetworkTestHandler) HandleNetworkTest(w http.ResponseWriter, r *http.Request) {
	testType, err := h.extractTestParams(r)
	if err != nil {
		handleError(w, "test execution", err, http.StatusBadRequest)
		return
	}

	result, err := h.runAndSaveTest(r.Context(), testType)
	if err != nil {
		handleError(w, "test execution", err, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

// GetResults handles the retrieval of test results
func (h *NetworkTestHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	testType, date, err := h.extractResultParams(r)
	if err != nil {
		handleError(w, "parameter validation", err, http.StatusBadRequest)
		return
	}

	startDate := r.URL.Query().Get("startDate")
	if startDate != "" {
		h.getResultsRange(w, testType, startDate, date)
		return
	}

	result, err := h.repository.GetTestData(date, testType)
	if err != nil {
		handleError(w, "retrieving results", err, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, result)
}

func (h *NetworkTestHandler) getResultsRange(w http.ResponseWriter, testType, startDate, endDate string) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		handleError(w, "invalid start date format", err, http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		handleError(w, "invalid end date format", err, http.StatusBadRequest)
		return
	}

	results, err := h.repository.GetTestDataInRange(start, end, testType)
	if err != nil {
		handleError(w, "retrieving results", err, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, results)
}

func (h *NetworkTestHandler) runAndSaveTest(ctx context.Context, testType string) (interface{}, error) {
	results, err := h.tester.RunTest(ctx, []string{testType})
	if err != nil {
		return nil, fmt.Errorf("test execution failed: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no test results returned")
	}

	_, err = h.repository.SaveTestResult(results[0], testType)
	if err != nil {
		return nil, fmt.Errorf("failed to save test result: %w", err)
	}

	go func() {
		if err := h.generateAndSaveCharts(results[0], testType); err != nil {
			log.Printf("Chart generation failed: %v", err)
		}
	}()

	return results[0], nil
}

func (h *NetworkTestHandler) extractTestParams(r *http.Request) (string, error) {
	testType := r.URL.Query().Get("test")
	if testType == "" {
		return "", fmt.Errorf("missing test type parameter: 'test'")
	}

	return testType, nil
}

func (h *NetworkTestHandler) extractResultParams(r *http.Request) (string, string, error) {
	testType := r.URL.Query().Get("test")
	if testType == "" {
		return "", "", fmt.Errorf("missing test type parameter")
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format(dateFormat)
	}

	if _, err := time.Parse(dateFormat, date); err != nil {
		return "", "", fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	return testType, date, nil
}

func (h *NetworkTestHandler) generateAndSaveCharts(result interface{}, testType string) error {
	switch testType {
	case "icmp":
		if icmpResult, ok := result.(*networkTesting.ICMPTestResult); ok {
			// Generate distribution pie
			pieChart, err := h.charts.GenerateICMPDistributionPie(icmpResult)
			if err != nil {
				return fmt.Errorf("failed to generate distribution chart: %w", err)
			}
			if _, err := h.repository.SaveChart(pieChart, "icmp", "distribution"); err != nil {
				log.Printf("Failed to save distribution chart: %v", err)
			}

			// Generate RTT line
			lineChart, err := h.charts.GenerateICMPRTTLine(icmpResult)
			if err != nil {
				return fmt.Errorf("failed to generate RTT chart: %w", err)
			}
			if _, err := h.repository.SaveChart(lineChart, "icmp", "rtt"); err != nil {
				log.Printf("Failed to save RTT chart: %v", err)
			}
		}
	case "download":
		// implement download test

	case "upload":
		// implement download test

	}
	return nil
}
