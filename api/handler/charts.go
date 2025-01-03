package handler

import (
	"fmt"
	"log"
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

	go func() {
		if err := h.generateAndSaveCharts(result, testType); err != nil {
			log.Printf("Chart generation failed: %v", err)
			handleError(w, "Chart generation failed", err, http.StatusInternalServerError)
		}
	}()
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

	err = h.generateAndSaveHistoricCharts(results, testType)
	if err != nil {
		handleError(w, "error while generating charts", err, http.StatusInternalServerError)
	}
}

func (h *ChartHandler) generateAndSaveHistoricCharts(results []*networkTesting.TestResult, testType string) error {
	switch testType {
	case "icmp":

	case "download":
		downloadResults := make([]*networkTesting.AverageSpeedTestResult, len(results))
		for i, r := range results {
			downloadResults[i] = r.Download
		}
		bar, err := h.charts.GenerateHistoricDownloadAnalysisCharts(downloadResults)
		if err != nil {
			return fmt.Errorf("failed to generate download chart: %w", err)
		}
		h.repository.SaveChart(bar, "download", "speed_ot")
	case "upload":
		uploadResults := make([]*networkTesting.AverageSpeedTestResult, len(results))
		for i, r := range results {
			uploadResults[i] = r.Upload
		}
		bar, err := h.charts.GenerateHistoricUploadAnalysisCharts(uploadResults)
		if err != nil {
			return fmt.Errorf("failed to generate upload chart: %w", err)
		}
		h.repository.SaveChart(bar, "upload", "speed_ot")
	case "route":
		result := make([]*networkTesting.RouteTestResult, len(results))
		for i, r := range results {
			result[i] = r.Route
		}
		barline, err := h.charts.GenerateHistoricRouteAnalysisCharts(result)
		if err != nil {
			return fmt.Errorf("failed to generate upload chart: %w", err)
		}
		h.repository.SaveChart(barline, "route", "rtt_ot")
	case "latency":
		jitterReults := make([]*networkTesting.JitterTestResult, len(results))
		for i, r := range results {
			jitterReults[i] = r.Jitter
		}
		bar, err := h.charts.GenerateHistoricJitterAnalysisCharts(jitterReults)
		if err != nil {
			return fmt.Errorf("failed to generate download chart: %w", err)
		}
		h.repository.SaveChart(bar, "latency", "rtt_ot")
	case "bandwidth":

	default:
		return fmt.Errorf("unsupported test type: %s", testType)
	}
	return nil
}

func (h *ChartHandler) generateAndSaveCharts(result *networkTesting.TestResult, testType string) error {
	switch testType {
	case "icmp":
		pieChart, err := h.charts.GenerateICMPAnalysisCharts(result.ICMP)
		if err != nil {
			return fmt.Errorf("failed to generate ICMP chart: %w", err)
		}
		if _, err := h.repository.SaveChart(pieChart, "icmp", "distribution"); err != nil {
			log.Printf("Failed to save ICMP chart: %v", err)
		}
	case "download":
		bar, err := h.charts.GenerateDownloadAnalysisCharts(result.Download)
		if err != nil {
			return fmt.Errorf("failed to generate download chart: %w", err)
		}
		if _, err := h.repository.SaveChart(bar, "download", "speed"); err != nil {
			log.Printf("Failed to save download chart: %v", err)
		}
	case "upload":
		bar, err := h.charts.GenerateUploadAnalysisCharts(result.Upload)
		if err != nil {
			return fmt.Errorf("failed to generate upload chart: %w", err)
		}
		if _, err := h.repository.SaveChart(bar, "upload", "speed"); err != nil {
			log.Printf("Failed to save upload chart: %v", err)
		}
	case "route":
		lineChart, err := h.charts.GenerateRouteAnalysisCharts(result.Route)
		if err != nil {
			return fmt.Errorf("failed to generate route chart: %w", err)
		}
		if _, err := h.repository.SaveChart(lineChart, "route", "path"); err != nil {
			log.Printf("Failed to save route chart: %v", err)
		}
	case "latency":
		lineChart, err := h.charts.GenerateJitterAnalysisCharts(result.Jitter)
		if err != nil {
			return fmt.Errorf("failed to generate latency chart: %w", err)
		}
		if _, err := h.repository.SaveChart(lineChart, "latency", "path"); err != nil {
			log.Printf("Failed to save latency chart: %v", err)
		}
	case "bandwidth":
		bar3dSpeed, bar3dDuration, err := h.charts.GenerateBandwidthAnalysisCharts(result.Bandwidth)
		if err != nil {
			return fmt.Errorf("failed to generate bandwidth charts: %w", err)
		}
		if _, err := h.repository.SaveChart(bar3dSpeed, "bandwidth", "speed"); err != nil {
			log.Printf("Failed to save bandwidth speed chart: %v", err)
		}
		if _, err := h.repository.SaveChart(bar3dDuration, "bandwidth", "duration"); err != nil {
			log.Printf("Failed to save bandwidth duration chart: %v", err)
		}

	default:
		return fmt.Errorf("unsupported test type: %s", testType)
	}
	return nil
}
