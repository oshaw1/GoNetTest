package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateDownloadAnalysisCharts(result *networkTesting.AverageSpeedTestResult) (*charts.Bar, error) {
	bar, err := generateDownloadSpeedBar(result)
	if err != nil {
		return nil, err
	}

	return bar, nil
}

func generateDownloadSpeedBar(result *networkTesting.AverageSpeedTestResult) (*charts.Bar, error) {
	bar := charts.NewBar()

	// Prepare data
	var xAxis []string
	var speeds []float64

	for url, test := range result.TestedURLs {
		xAxis = append(xAxis, url)
		speeds = append(speeds, test.Speed)
	}

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Upload Speeds by URL",
			Subtitle: fmt.Sprintf("Test ran at: %v", result.Timestamp.Format("2006-01-02 15:04:05")),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Formatter: "{b}: {c} Mbps",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:         "Speed (Mbps)",
			NameLocation: "middle",
			NameGap:      35,
			AxisLabel: &opts.AxisLabel{
				Show:         opts.Bool(true),
				Rotate:       90, // Rotate labels 45 degrees
				ShowMaxLabel: opts.Bool(true),
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Show:         opts.Bool(true),
				Rotate:       45, // Rotate labels 45 degrees
				ShowMaxLabel: opts.Bool(true),
			},
		}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithColorsOpts(opts.Colors{"#4169E1"}),
		// Add some bottom margin to accommodate angled labels
		charts.WithGridOpts(opts.Grid{
			Bottom: "20%",
			Top:    "10%",
		}),
	)

	bar.SetXAxis(xAxis).
		AddSeries("Download Speed", generateBarItems(speeds))

	return bar, nil
}
