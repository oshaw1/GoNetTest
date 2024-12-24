package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateJitterAnalysisCharts(result *networkTesting.JitterTestResult) (*charts.Line, *charts.Bar, error) {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Packet RTT Over Time",
		Subtitle: fmt.Sprintf("Target: %s", result.Target),
	}))

	xAxis := make([]int, len(result.RTTs))
	rttData := make([]opts.LineData, len(result.RTTs))
	for i, rtt := range result.RTTs {
		xAxis[i] = i + 1
		rttData[i] = opts.LineData{Value: float64(rtt.Milliseconds())}
	}
	line.SetXAxis(xAxis).AddSeries("RTT (ms)", rttData)

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Jitter Statistics",
	}))

	statNames := []string{"Min Jitter", "Avg Jitter", "Max Jitter"}
	statValues := []float64{
		float64(result.MinJitter.Milliseconds()),
		float64(result.AvgJitter.Milliseconds()),
		float64(result.MaxJitter.Milliseconds()),
	}

	bar.SetXAxis(statNames).AddSeries("Jitter (ms)", generateBarItems(statValues))

	return line, bar, nil
}

func generateBarItems(values []float64) []opts.BarData {
	items := make([]opts.BarData, len(values))
	for i := range values {
		items[i] = opts.BarData{Value: values[i]}
	}
	return items
}
