package charting

import (
	"fmt"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateJitterAnalysisCharts(result *networkTesting.JitterTestResult) (*charts.Line, error) {
	line, err := generateJitterLineChart(result)
	if err != nil {
		return nil, err
	}

	return line, nil
}

func generateJitterLineChart(result *networkTesting.JitterTestResult) (*charts.Line, error) {
	var totalTime time.Duration
	for _, rtt := range result.RTTs {
		totalTime += rtt
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Packet RTT Over Time",
			Subtitle: fmt.Sprintf("Total Time: %.2fms  Target: %s", float64(totalTime.Microseconds())/1000, result.Target),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:         "Return time (Ms)",
			NameLocation: "middle",
			NameGap:      35,
			AxisLabel: &opts.AxisLabel{
				Show:         opts.Bool(true),
				Rotate:       90,
				ShowMaxLabel: opts.Bool(true),
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name:         "Package",
			NameLocation: "middle",
			NameGap:      35,
			AxisLabel: &opts.AxisLabel{
				Show:         opts.Bool(true),
				ShowMaxLabel: opts.Bool(true),
			},
		}),
	)

	xAxis := make([]int, len(result.RTTs))
	rttData := make([]opts.LineData, len(result.RTTs))
	for i, rtt := range result.RTTs {
		xAxis[i] = i + 1
		rttData[i] = opts.LineData{Value: float64(rtt.Milliseconds())}
	}

	line.SetXAxis(xAxis).AddSeries("RTT (ms)", rttData)
	return line, nil
}
