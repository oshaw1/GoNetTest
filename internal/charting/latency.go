package charting

import (
	"fmt"
	"math"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateLatencyAnalysisCharts(result *networkTesting.LatencyTestResult) (*charts.Line, error) {
	line, err := generateLatencyLineChart(result)
	if err != nil {
		return nil, err
	}

	return line, nil
}

func (g *Generator) GenerateHistoricLatencyAnalysisCharts(results []*networkTesting.LatencyTestResult) (*charts.Bar, error) {
	if results == nil {
		return nil, fmt.Errorf("function called with no results")
	}

	barOverTime, err := generateLatencyOverTimeBar(results)
	if err != nil {
		return nil, err
	}

	return barOverTime, nil
}

func generateLatencyLineChart(result *networkTesting.LatencyTestResult) (*charts.Line, error) {
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

func generateLatencyOverTimeBar(results []*networkTesting.LatencyTestResult) (*charts.Bar, error) {
	bar := charts.NewBar()
	var xAxis []string
	var avgLatency []float64
	var minLatency []float64
	var maxLatency []float64

	for _, result := range results {
		var Latencys []float64
		for i := 1; i < len(result.RTTs); i++ {
			Latency := math.Abs(float64(result.RTTs[i]-result.RTTs[i-1])) / 1000000
			Latencys = append(Latencys, Latency)
		}

		xAxis = append(xAxis, result.Timestamp.Format("2006-01-02 15:04:05"))
		avgLatency = append(avgLatency, calculateAverage(Latencys))
		minLatency = append(minLatency, findMin(Latencys))
		maxLatency = append(maxLatency, findMax(Latencys))
	}

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Latency Over Time",
			Subtitle: fmt.Sprintf("Test data from: %v Days", len(results)),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    opts.Bool(true),
			Trigger: "axis",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:         "RTT (ms)",
			NameLocation: "middle",
			NameGap:      35,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Show:         opts.Bool(true),
				Rotate:       45,
				ShowMaxLabel: opts.Bool(true),
			},
		}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithGridOpts(opts.Grid{
			Bottom: "20%",
			Top:    "10%",
		}),
	)

	bar.SetXAxis(xAxis).
		AddSeries("Min RTT", generateBarItems(minLatency)).
		AddSeries("Average RTT", generateBarItems(avgLatency)).
		AddSeries("Max RTT", generateBarItems(maxLatency))

	return bar, nil
}
