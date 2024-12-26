package charting

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateDownloadAnalysisCharts(result *networkTesting.AverageSpeedTestResult) (*charts.Line, error) {
	line, err := generateDownloadChart(result)
	if err != nil {
		return nil, err
	}

	return line, nil
}

func generateDownloadChart(result *networkTesting.AverageSpeedTestResult) (*charts.Line, error) {
	line := charts.NewLine()
	return line, nil
}
