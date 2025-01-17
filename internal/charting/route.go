package charting

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateRouteAnalysisCharts(result *networkTesting.RouteTestResult) (*charts.Line, error) {
	pie, err := generateRouteRequestPathChart(result)
	if err != nil {
		return nil, err
	}

	return pie, nil
}

func generateRouteRequestPathChart(result *networkTesting.RouteTestResult) (*charts.Line, error) {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Route Test RTT by Hop",
		Subtitle: fmt.Sprintf("Target: %s, Test ran at: %v", result.Target, result.Timestamp),
	}))

	xAxis := make([]string, 0)
	rttData := make([]opts.LineData, 0)

	for _, hop := range result.Hops {
		hopLabel := fmt.Sprintf("%d", hop.Number)
		xAxis = append(xAxis, hopLabel)
		if hop.Lost {
			rttData = append(rttData, opts.LineData{Value: 0.0})
		} else {
			rttData = append(rttData, opts.LineData{Value: float64(hop.RTT) / 1e6})
		}
	}

	line.SetXAxis(xAxis).AddSeries("RTT (s)", rttData)
	return line, nil
}

func (g *Generator) GenerateHistoricRouteAnalysisCharts(results []*networkTesting.RouteTestResult) (*charts.Bar3D, error) {
	if results == nil {
		return nil, fmt.Errorf("function called with no results")
	}

	bar3d, err := generateRoute3DBar(results)
	if err != nil {
		return nil, err
	}

	return bar3d, nil
}

func generateRoute3DBar(results []*networkTesting.RouteTestResult) (*charts.Bar3D, error) {
	bar3d := charts.NewBar3D()

	hopNumbers := make(map[int]bool)
	var maxRTT time.Duration
	for _, result := range results {
		for _, hop := range result.Hops {
			hopNumbers[hop.Number] = true
			if !hop.Lost && hop.RTT > maxRTT {
				maxRTT = hop.RTT
			}
		}
	}

	uniqueHops := make([]int, 0, len(hopNumbers))
	for hop := range hopNumbers {
		uniqueHops = append(uniqueHops, hop)
	}
	sort.Ints(uniqueHops)

	// Set global options including color mapping
	bar3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Route Latency Analysis Over Time",
			Subtitle: fmt.Sprintf("Test data from: %v Days", len(results)),
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: opts.Bool(true),
			Max:        float32(maxRTT.Milliseconds()),
			Min:        0,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8",
					"#ffffbf", "#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026"},
			},
		}),
	)

	var data []opts.Chart3DData
	for resultIdx, result := range results {
		for _, hop := range result.Hops {
			rtt := 0.0
			if !hop.Lost {
				rtt = float64(hop.RTT.Milliseconds())
			}
			data = append(data, opts.Chart3DData{
				Value: []interface{}{
					resultIdx,  // Keep numeric index for x-axis
					hop.Number, // Keep numeric value for y-axis
					rtt,        // Keep numeric value for z-axis
				},
				Name: fmt.Sprintf("IP: %s ", hop.Address),
			})
		}
	}
	xAxis := make([]string, len(results))
	for i := range xAxis {
		xAxis[i] = results[i].Timestamp.Format("2006-01-02 15:04:05")
	}

	yAxis := make([]int, len(uniqueHops))
	for i := range yAxis {
		yAxis[i] = i + 1
	}

	bar3d.AddSeries(" ", data).
		SetSeriesOptions(
			charts.WithBar3DChartOpts(opts.Bar3DChart{
				Shading: "lambert",
			}),
		)

	maxHop := len(uniqueHops) // Get largest hop number
	// Set 3D axis options with explicit step size
	bar3d.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
		charts.WithXAxis3DOpts(opts.XAxis3D{
			Name: "Timestamp",
			Data: xAxis,
		}),
		charts.WithYAxis3DOpts(opts.YAxis3D{
			Show: opts.Bool(true),
			Data: yAxis,
			Max:  float32(maxHop),
			Min:  1,
			Name: "Hop",
			Type: "value",
		}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{
			Name: "RTT (ms)",
		}),
	)

	return bar3d, nil
}
