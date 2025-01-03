package charting

import (
	"fmt"
	"math"
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

func (g *Generator) GenerateHistoricJitterAnalysisCharts(results []*networkTesting.JitterTestResult) (*charts.Bar, error) {
	if results == nil {
		return nil, fmt.Errorf("function called with no results")
	}

	barOverTime, err := generateJitterOverTimeBar(results)
	if err != nil {
		return nil, err
	}

	return barOverTime, nil
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

func generateJitterOverTimeBar(results []*networkTesting.JitterTestResult) (*charts.Bar, error) {
	bar := charts.NewBar()
	var xAxis []string
	var avgJitter []float64
	var minJitter []float64
	var maxJitter []float64

	for _, result := range results {
		var jitters []float64
		for i := 1; i < len(result.RTTs); i++ {
			jitter := math.Abs(float64(result.RTTs[i]-result.RTTs[i-1])) / 1000000
			jitters = append(jitters, jitter)
		}

		xAxis = append(xAxis, result.Timestamp.Format("2006-01-02 15:04:05"))
		avgJitter = append(avgJitter, calculateAverage(jitters))
		minJitter = append(minJitter, findMin(jitters))
		maxJitter = append(maxJitter, findMax(jitters))
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
		AddSeries("Min RTT", generateBarItems(minJitter)).
		AddSeries("Average RTT", generateBarItems(avgJitter)).
		AddSeries("Max RTT", generateBarItems(maxJitter))

	return bar, nil
}
