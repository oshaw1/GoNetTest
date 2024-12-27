package charting

import (
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateBandwidthAnalysisCharts(result *networkTesting.BandwidthTestResult) (*charts.Bar3D, *charts.Bar3D, error) {
	bar3dSpeed, err := generateBandwidth3DBarSpeed(result)
	if err != nil {
		return nil, nil, err
	}
	bar3dDuration, err := generateBandwidth3DBarDuration(result)
	if err != nil {
		return nil, nil, err
	}

	return bar3dSpeed, bar3dDuration, nil
}

func generateBandwidth3DBarSpeed(result *networkTesting.BandwidthTestResult) (*charts.Bar3D, error) {
	bar3d := charts.NewBar3D()

	bar3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Bandwidth Analysis by Connection",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: opts.Bool(true),
			Max:        float32(result.MaxThroughput + 25),
			Min:        0,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8",
					"#ffffbf", "#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026"},
			},
		}),
	)

	var data []opts.Chart3DData

	for stepIdx, step := range result.Steps {
		// For each connection in the step
		for _, conn := range step.ConnResults {
			data = append(data, opts.Chart3DData{
				Value: []interface{}{
					stepIdx,     // X axis (step number)
					conn.ID + 1, // Y axis (connection number)
					conn.Speed,  // Z axis (speed)
				},
			})
		}
	}

	xAxis := make([]int, len(result.Steps))
	for i := range xAxis {
		xAxis[i] = i + 1
	}

	maxConns := 0
	for _, step := range result.Steps {
		if step.Connections > maxConns {
			maxConns = step.Connections
		}
	}

	yAxis := make([]int, maxConns)
	for i := range yAxis {
		yAxis[i] = i + 1
	}

	bar3d.AddSeries("Bandwidth", data).
		SetSeriesOptions(
			charts.WithBar3DChartOpts(opts.Bar3DChart{
				Shading: "lambert",
			}),
		)

	bar3d.SetGlobalOptions(
		charts.WithXAxis3DOpts(opts.XAxis3D{
			Name: "Step",
			Data: xAxis,
		}),
		charts.WithYAxis3DOpts(opts.YAxis3D{
			Name: "Connection",
			Data: yAxis,
		}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{
			Name: "Speed (Mbps)",
		}),
	)

	return bar3d, nil
}

func generateBandwidth3DBarDuration(result *networkTesting.BandwidthTestResult) (*charts.Bar3D, error) {
	bar3d := charts.NewBar3D()

	var maxDuration time.Duration
	for _, step := range result.Steps {
		for _, conn := range step.ConnResults {
			if conn.Duration > maxDuration {
				maxDuration = conn.Duration
			}
		}
	}

	bar3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Connection Duration Analysis",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: opts.Bool(true),
			Max:        float32(maxDuration.Seconds()),
			Min:        0,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8",
					"#ffffbf", "#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026"},
			},
		}),
	)

	var data []opts.Chart3DData

	for stepIdx, step := range result.Steps {
		for _, conn := range step.ConnResults {
			data = append(data, opts.Chart3DData{
				Value: []interface{}{
					stepIdx,                          // X axis (step number)
					conn.ID + 1,                      // Y axis (connection number)
					float64(conn.Duration.Seconds()), // Z axis (duration in seconds)
				},
			})
		}
	}

	xAxis := make([]int, len(result.Steps))
	for i := range xAxis {
		xAxis[i] = i + 1
	}

	maxConns := 0
	for _, step := range result.Steps {
		if step.Connections > maxConns {
			maxConns = step.Connections
		}
	}

	yAxis := make([]int, maxConns)
	for i := range yAxis {
		yAxis[i] = i + 1
	}

	bar3d.AddSeries("Duration", data).
		SetSeriesOptions(
			charts.WithBar3DChartOpts(opts.Bar3DChart{
				Shading: "lambert",
			}),
		)

	bar3d.SetGlobalOptions(
		charts.WithXAxis3DOpts(opts.XAxis3D{
			Name: "Step",
			Data: xAxis,
		}),
		charts.WithYAxis3DOpts(opts.YAxis3D{
			Name: "Connection",
			Data: yAxis,
		}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{
			Name: "Duration (seconds)",
		}),
	)

	return bar3d, nil
}
