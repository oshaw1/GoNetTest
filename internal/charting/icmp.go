package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateICMPAnalysisCharts(result *networkTesting.ICMPTestResult) (*charts.Pie, error) {
	pie, err := generateICMPDistributionPie(result)
	if err != nil {
		return nil, err
	}

	return pie, nil
}

func generateICMPDistributionPie(result *networkTesting.ICMPTestResult) (*charts.Pie, error) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "ICMP Packet Distribution",
		Subtitle: fmt.Sprintf("Test ran at: %v", result.Timestamp),
	}))

	pie.AddSeries("Packet Distribution", []opts.PieData{
		{Name: "Received", Value: result.Received},
		{Name: "Lost", Value: result.Lost},
	}).SetSeriesOptions(charts.WithLabelOpts(opts.Label{
		Formatter: "{b}: {c} ({d}%)",
	}))

	return pie, nil
}
