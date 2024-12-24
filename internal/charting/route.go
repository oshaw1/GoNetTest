package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateRouteTestChart(result *networkTesting.RouteTestResult) (*charts.Line, error) {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Route Test RTT by Hop",
		Subtitle: fmt.Sprintf("Target: %s, Test ran at: %v", result.Target, result.Timestamp),
	}))

	xAxis := make([]string, 0)
	rttData := make([]opts.LineData, 0)

	for _, hop := range result.Hops {
		hopLabel := fmt.Sprintf("%d", hop.Number)
		if hop.Address != "" {
			hopLabel = fmt.Sprintf("%d\n%s", hop.Number, hop.Address)
		}
		xAxis = append(xAxis, hopLabel)
		if hop.Lost {
			rttData = append(rttData, opts.LineData{Value: 0.0})
		} else {
			rttData = append(rttData, opts.LineData{Value: float64(hop.RTT) / 1e6})
		}
	}

	line.SetXAxis(xAxis).AddSeries("RTT (ms)", rttData)
	return line, nil
}
