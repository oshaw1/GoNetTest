package handler

import (
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
	charts     *charting.Generator
}

func NewNetworkTestHandler(tester *networkTesting.NetworkTester, repo *dataManagement.Repository) *NetworkTestHandler {
	return &NetworkTestHandler{
		tester:     tester,
		repository: repo,
		charts:     charting.NewGenerator(),
	}
}

// HandleNetworkTest handles the test execution request
func (h *NetworkTestHandler) HandleNetworkTest(w http.ResponseWriter, r *http.Request) {
	testType := r.URL.Query().Get("test")
	if testType == "" {
		handleError(w, "missing test type parameter: 'test'", nil, http.StatusBadRequest)
		return
	}

	result, err := h.runAndSaveTest(testType)
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

	startDate := r.URL.Query().Get("date")
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

func (h *NetworkTestHandler) runAndSaveTest(testType string) (interface{}, error) {
	result, err := h.tester.RunTest(testType)
	if err != nil {
		return nil, fmt.Errorf("test execution failed: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("no test results returned")
	}

	_, err = h.repository.SaveTestResult(result, testType)
	if err != nil {
		return nil, fmt.Errorf("failed to save test result: %w", err)
	}

	go func() {
		if err := h.generateAndSaveCharts(result, testType); err != nil {
			log.Printf("Chart generation failed: %v", err)
		}
	}()

	return result, nil
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
			pieChart, err := h.charts.GenerateICMPAnalysisCharts(icmpResult)
			if err != nil {
				return fmt.Errorf("failed to generate distribution chart: %w", err)
			}
			if _, err := h.repository.SaveChart(pieChart, "icmp", "distribution"); err != nil {
				log.Printf("Failed to save distribution chart: %v", err)
			}
		}
	case "download":
		if downloadResult, ok := result.(*networkTesting.AverageSpeedTestResult); ok {
			bar, err := h.charts.GenerateDownloadAnalysisCharts(downloadResult)
			if err != nil {
				return fmt.Errorf("failed to generate distribution chart: %w", err)
			}
			if _, err := h.repository.SaveChart(bar, "download", "speed"); err != nil {
				log.Printf("Failed to save route line chart: %v", err)
			}
		}
	case "upload":
		if uploadResult, ok := result.(*networkTesting.AverageSpeedTestResult); ok {
			bar, err := h.charts.GenerateUploadAnalysisCharts(uploadResult)
			if err != nil {
				return fmt.Errorf("failed to generate distribution chart: %w", err)
			}
			if _, err := h.repository.SaveChart(bar, "upload", "speed"); err != nil {
				log.Printf("Failed to save route line chart: %v", err)
			}
		}
	case "route":
		if routeResult, ok := result.(*networkTesting.RouteTestResult); ok {
			lineChart, err := h.charts.GenerateRouteAnalysisCharts(routeResult)
			if err != nil {
				return fmt.Errorf("failed to generate distribution chart: %w", err)
			}
			if _, err := h.repository.SaveChart(lineChart, "route", "path"); err != nil {
				log.Printf("Failed to save route line chart: %v", err)
			}
		}
	case "latency":
		if latencyResult, ok := result.(*networkTesting.LatencyTestResult); ok {
			lineChart, err := h.charts.GenerateLatencyAnalysisCharts(latencyResult)
			if err != nil {
				return fmt.Errorf("failed to generate distribution chart: %w", err)
			}
			if _, err := h.repository.SaveChart(lineChart, "latency", "path"); err != nil {
				log.Printf("Failed to save route line chart: %v", err)
			}
		}
	case "bandwidth":
		if bandwidthResult, ok := result.(*networkTesting.BandwidthTestResult); ok {
			bar3dSpeed, bar3dDuration, err := h.charts.GenerateBandwidthAnalysisCharts(bandwidthResult)
			if err != nil {
				return fmt.Errorf("failed to generate distribution chart: %w", err)
			}
			if _, err := h.repository.SaveChart(bar3dSpeed, "bandwidth", "speed"); err != nil {
				log.Printf("Failed to save route line chart: %v", err)
			}
			if _, err := h.repository.SaveChart(bar3dDuration, "bandwidth", "duration"); err != nil {
				log.Printf("Failed to save route line chart: %v", err)
			}
		}
	default:
		return fmt.Errorf("unsupported test type: %s", testType)
	}
	return nil
}
