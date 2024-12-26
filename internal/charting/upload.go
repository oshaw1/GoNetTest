package charting

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateUploadAnalysisCharts(result *networkTesting.AverageSpeedTestResult) (*charts.Line, error) {
	line, err := generateUploadChart(result)
	if err != nil {
		return nil, err
	}

	return line, nil
}

func generateUploadChart(result *networkTesting.AverageSpeedTestResult) (*charts.Line, error) {
	line := charts.NewLine()
	return line, nil
}
