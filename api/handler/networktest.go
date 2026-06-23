package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/oshaw1/go-net-test/internal/charting"
	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

const dateFormat = "2006-01-02"

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

	resultID, err := h.repository.SaveTestResult(result, testType)
	if err != nil {
		return nil, fmt.Errorf("failed to save test result: %w", err)
	}

	// Generated synchronously (not fire-and-forget) so the chart is already
	// saved by the time the response goes back — callers that trigger a UI
	// refresh on completion (e.g. the dashboard's quick-test button) would
	// otherwise race the chart write and refresh too early.
	if err := h.generateAndSaveCharts(result, testType, resultID); err != nil {
		log.Printf("Chart generation failed: %v", err)
	}

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

func (h *NetworkTestHandler) generateAndSaveCharts(result interface{}, testType string, resultID int64) error {
	switch testType {
	case "icmp":
		if icmpResult, ok := result.(*networkTesting.ICMPTestResult); ok {
			pieChart, err := h.charts.GenerateICMPAnalysisCharts(icmpResult)
			if err != nil {
				return fmt.Errorf("failed to generate ICMP chart: %w", err)
			}
			if _, err := h.repository.SaveChart(pieChart, "icmp", "distribution", resultID); err != nil {
				log.Printf("Failed to save ICMP chart: %v", err)
			}
		}
	case "download":
		if downloadResult, ok := result.(*networkTesting.AverageSpeedTestResult); ok {
			bar, err := h.charts.GenerateDownloadAnalysisCharts(downloadResult)
			if err != nil {
				return fmt.Errorf("failed to generate download chart: %w", err)
			}
			if _, err := h.repository.SaveChart(bar, "download", "speed", resultID); err != nil {
				log.Printf("Failed to save download chart: %v", err)
			}
		}
	case "upload":
		if uploadResult, ok := result.(*networkTesting.AverageSpeedTestResult); ok {
			bar, err := h.charts.GenerateUploadAnalysisCharts(uploadResult)
			if err != nil {
				return fmt.Errorf("failed to generate upload chart: %w", err)
			}
			if _, err := h.repository.SaveChart(bar, "upload", "speed", resultID); err != nil {
				log.Printf("Failed to save upload chart: %v", err)
			}
		}
	case "route":
		if routeResult, ok := result.(*networkTesting.RouteTestResult); ok {
			lineChart, err := h.charts.GenerateRouteAnalysisCharts(routeResult)
			if err != nil {
				return fmt.Errorf("failed to generate route chart: %w", err)
			}
			if _, err := h.repository.SaveChart(lineChart, "route", "path", resultID); err != nil {
				log.Printf("Failed to save route chart: %v", err)
			}
		}
	case "latency":
		if latencyResult, ok := result.(*networkTesting.LatencyTestResult); ok {
			lineChart, err := h.charts.GenerateLatencyAnalysisCharts(latencyResult)
			if err != nil {
				return fmt.Errorf("failed to generate latency chart: %w", err)
			}
			if _, err := h.repository.SaveChart(lineChart, "latency", "path", resultID); err != nil {
				log.Printf("Failed to save latency chart: %v", err)
			}
		}
	case "bandwidth":
		if bandwidthResult, ok := result.(*networkTesting.BandwidthTestResult); ok {
			bar3dSpeed, bar3dDuration, err := h.charts.GenerateBandwidthAnalysisCharts(bandwidthResult)
			if err != nil {
				return fmt.Errorf("failed to generate bandwidth charts: %w", err)
			}
			if _, err := h.repository.SaveChart(bar3dSpeed, "bandwidth", "speed", resultID); err != nil {
				log.Printf("Failed to save bandwidth speed chart: %v", err)
			}
			if _, err := h.repository.SaveChart(bar3dDuration, "bandwidth", "duration", resultID); err != nil {
				log.Printf("Failed to save bandwidth duration chart: %v", err)
			}
		}
	default:
		return fmt.Errorf("unsupported test type: %s", testType)
	}
	return nil
}

func (h *NetworkTestHandler) HandleDeleteTests(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Missing date property", http.StatusBadRequest)
		return
	}

	if err := h.repository.DeleteByDate(date); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *NetworkTestHandler) HandleDeleteTestResult(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Missing or invalid id property", http.StatusBadRequest)
		return
	}

	if err := h.repository.DeleteByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *NetworkTestHandler) HandleDeleteCharts(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		http.Error(w, "Missing ids property", http.StatusBadRequest)
		return
	}

	parts := strings.Split(idsParam, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ids property", http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	if err := h.repository.DeleteChartsByIDs(ids); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
