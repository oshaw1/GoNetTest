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

func (h *ChartHandler) GenerateHistoricChart(w http.ResponseWriter, r *http.Request) {
	testType := r.URL.Query().Get("test")
	if testType == "" {
		handleError(w, "missing test type parameter: 'test'", nil, http.StatusBadRequest)
	}

	days, err := strconv.Atoi(r.URL.Query().Get("days"))
	if err != nil {
		handleError(w, "invalid days parameter", err, http.StatusBadRequest)
	}

	start := time.Now().AddDate(0, 0, -days)

	results, err := h.repository.GetTestDataInRange(start, time.Now(), testType)
	if err != nil {
		handleError(w, "error retrieving data", err, http.StatusInternalServerError)
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
