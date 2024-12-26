package charting

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateBandwidthAnalysisCharts(result *networkTesting.BandwidthTestResult) (*charts.Line, error) {
	line, err := generateBandwidthChart(result)
	if err != nil {
		return nil, err
	}

	return line, nil
}

func generateBandwidthChart(result *networkTesting.BandwidthTestResult) (*charts.Line, error) {
	line := charts.NewLine()
	return line, nil
}
