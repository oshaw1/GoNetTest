package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateICMPDistributionPie(result *networkTesting.ICMPTestResult) (*charts.Pie, error) {
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

// GenerateICMPRTTLine creates a line chart showing RTT times
func (g *Generator) GenerateICMPRTTLine(result *networkTesting.ICMPTestResult) (*charts.Line, error) {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "ICMP Round Trip Times",
		Subtitle: fmt.Sprintf("Test ran at: %v", result.Timestamp),
	}))

	xAxis := []string{"Min RTT", "Avg RTT", "Max RTT"}
	yAxis := []opts.LineData{
		{Value: result.MinRTT.Milliseconds()},
		{Value: result.AvgRTT.Milliseconds()},
		{Value: result.MaxRTT.Milliseconds()},
	}

	line.SetXAxis(xAxis)
	line.AddSeries("Time (ms)", yAxis)

	return line, nil
}
