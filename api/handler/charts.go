package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/charting"
	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

type ChartHandler struct {
	repository *dataManagement.Repository
	charts     *charting.Generator
	config     *config.Config
}

func NewChartHandler(repo *dataManagement.Repository, conf *config.Config) *ChartHandler {
	return &ChartHandler{
		repository: repo,
		charts:     charting.NewGenerator(),
		config:     conf,
	}
}

func (h *ChartHandler) GenerateChart(w http.ResponseWriter, r *http.Request) {
	testType := r.URL.Query().Get("test")
	if testType == "" {
		handleError(w, "missing test type parameter: 'test'", nil, http.StatusBadRequest)
		return
	}
	date := r.URL.Query().Get("date")
	if date == "" {
		handleError(w, "missing date parameter", nil, http.StatusBadRequest)
		return
	}
	result, err := h.repository.GetTestData(date, testType)
	if err != nil {
		handleError(w, "error retrieving data", err, http.StatusInternalServerError)
		return
	}
	if result == nil {
		handleError(w, "nil data returned", err, http.StatusInternalServerError)
		return
	}

	chartPath, err := h.generateAndSaveCharts(result, testType)
	if err != nil {
		handleError(w, "Chart generation failed", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(chartPath))
}

func (h *ChartHandler) GenerateHistoricChart(w http.ResponseWriter, r *http.Request) {
	testType := r.URL.Query().Get("test")
	if testType == "" {
		handleError(w, "missing test type parameter: 'test'", nil, http.StatusBadRequest)
		return
	}

	days, err := strconv.Atoi(r.URL.Query().Get("days"))
	if err != nil {
		handleError(w, "invalid days parameter", err, http.StatusBadRequest)
		return
	}

	start := time.Now().AddDate(0, 0, -days)

	results, err := h.repository.GetTestDataInRange(start, time.Now(), testType)
	if err != nil {
		handleError(w, "error retrieving data", err, http.StatusInternalServerError)
		return
	}

	chartPath, err := h.generateAndSaveHistoricCharts(results, testType)
	if err != nil {
		handleError(w, "error while generating charts", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(chartPath))
}

func (h *ChartHandler) generateAndSaveHistoricCharts(results []*networkTesting.TestResult, testType string) (string, error) {
	chartPath := ""
	switch testType {
	case "icmp":

	case "download":
		downloadResults := make([]*networkTesting.AverageSpeedTestResult, len(results))
		for i, r := range results {
			downloadResults[i] = r.Download
		}
		bar, err := h.charts.GenerateHistoricDownloadAnalysisCharts(downloadResults)
		if err != nil {
			return "", fmt.Errorf("failed to generate download chart: %w", err)
		}
		chartPath, _ = h.repository.SaveChart(bar, "download", "speed_ot")
	case "upload":
		uploadResults := make([]*networkTesting.AverageSpeedTestResult, len(results))
		for i, r := range results {
			uploadResults[i] = r.Upload
		}
		bar, err := h.charts.GenerateHistoricUploadAnalysisCharts(uploadResults)
		if err != nil {
			return "", fmt.Errorf("failed to generate upload chart: %w", err)
		}
		chartPath, _ = h.repository.SaveChart(bar, "upload", "speed_ot")
	case "route":
		result := make([]*networkTesting.RouteTestResult, len(results))
		for i, r := range results {
			result[i] = r.Route
		}
		barline, err := h.charts.GenerateHistoricRouteAnalysisCharts(result)
		if err != nil {
			return "", fmt.Errorf("failed to generate upload chart: %w", err)
		}
		chartPath, _ = h.repository.SaveChart(barline, "route", "rtt_ot")
	case "latency":
		latencyResults := make([]*networkTesting.LatencyTestResult, len(results))
		for i, r := range results {
			latencyResults[i] = r.Latency
		}
		bar, err := h.charts.GenerateHistoricLatencyAnalysisCharts(latencyResults)
		if err != nil {
			return "", fmt.Errorf("failed to generate latency chart: %w", err)
		}
		chartPath, _ = h.repository.SaveChart(bar, "latency", "latency_ot")
	case "bandwidth":
		bandwidthResult := make([]*networkTesting.BandwidthTestResult, len(results))
		for i, r := range results {
			bandwidthResult[i] = r.Bandwidth
		}
		speedBar, durationBar, err := h.charts.GenerateHistoricBandwidthAnalysisCharts(bandwidthResult)
		if err != nil {
			return "", fmt.Errorf("failed to generate download chart: %w", err)
		}
		chartPath, _ = h.repository.SaveChart(speedBar, "bandwidth", "bandwidth_speed_ot")
		chartPath2, _ := h.repository.SaveChart(durationBar, "bandwidth", "bandwidth_duration_ot")
		chartPath = chartPath + chartPath2
	default:
		return "", fmt.Errorf("unsupported test type: %s", testType)
	}
	return chartPath, nil
}

func (h *ChartHandler) generateAndSaveCharts(result *networkTesting.TestResult, testType string) (string, error) {
	chartPath := ""
	switch testType {
	case "icmp":
		pieChart, err := h.charts.GenerateICMPAnalysisCharts(result.ICMP)
		if err != nil {
			return "", fmt.Errorf("failed to generate ICMP chart: %w", err)
		}
		chartPath, err = h.repository.SaveChart(pieChart, "icmp", "distribution")
		if err != nil {
			return "", fmt.Errorf("failed to save ICMP chart: %w", err)
		}
	case "download":
		bar, err := h.charts.GenerateDownloadAnalysisCharts(result.Download)
		if err != nil {
			return "", fmt.Errorf("failed to generate download chart: %w", err)
		}
		chartPath, err = h.repository.SaveChart(bar, "download", "speed")
		if err != nil {
			return "", fmt.Errorf("failed to save download chart: %w", err)
		}
	case "upload":
		bar, err := h.charts.GenerateUploadAnalysisCharts(result.Upload)
		if err != nil {
			return "", fmt.Errorf("failed to generate upload chart: %w", err)
		}
		chartPath, err = h.repository.SaveChart(bar, "upload", "speed")
		if err != nil {
			return "", fmt.Errorf("failed to save upload chart: %w", err)
		}
	case "route":
		lineChart, err := h.charts.GenerateRouteAnalysisCharts(result.Route)
		if err != nil {
			return "", fmt.Errorf("failed to generate route chart: %w", err)
		}
		chartPath, err = h.repository.SaveChart(lineChart, "route", "path")
		if err != nil {
			return "", fmt.Errorf("failed to save route chart: %w", err)
		}
	case "latency":
		lineChart, err := h.charts.GenerateLatencyAnalysisCharts(result.Latency)
		if err != nil {
			return "", fmt.Errorf("failed to generate latency chart: %w", err)
		}
		chartPath, err = h.repository.SaveChart(lineChart, "latency", "path")
		if err != nil {
			return "", fmt.Errorf("failed to save latency chart: %w", err)
		}
	case "bandwidth":
		bar3dSpeed, bar3dDuration, err := h.charts.GenerateBandwidthAnalysisCharts(result.Bandwidth)
		if err != nil {
			return "", fmt.Errorf("failed to generate bandwidth charts: %w", err)
		}
		chartPath, err = h.repository.SaveChart(bar3dSpeed, "bandwidth", "speed")
		if err != nil {
			return "", fmt.Errorf("failed to save bandwidth speed chart: %w", err)
		}
		chartPath2, err := h.repository.SaveChart(bar3dDuration, "bandwidth", "duration")
		if err != nil {
			return "", fmt.Errorf("failed to save bandwidth duration chart: %w", err)
		}
		chartPath = chartPath + " " + chartPath2
	default:
		return "", fmt.Errorf("unsupported test type: %s", testType)
	}
	return chartPath, nil
}
