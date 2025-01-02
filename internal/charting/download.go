package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateDownloadAnalysisCharts(result *networkTesting.AverageSpeedTestResult) (*charts.Bar, error) {
	if result == nil {
		return nil, fmt.Errorf("function called with no results")
	}

	bar, err := generateDownloadSpeedBar(result)
	if err != nil {
		return nil, err
	}

	return bar, nil
}

func (g *Generator) GenerateHistoricDownloadAnalysisCharts(results []*networkTesting.AverageSpeedTestResult) (*charts.Bar, error) {
	if results == nil {
		return nil, fmt.Errorf("function called with no results")
	}

	barOverTime, err := generateDownloadOverTimeBar(results)
	if err != nil {
		return nil, err
	}

	return barOverTime, nil
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
			Title:    "Download Speeds by URL",
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

func generateDownloadOverTimeBar(results []*networkTesting.AverageSpeedTestResult) (*charts.Bar, error) {
	bar := charts.NewBar()
	line := charts.NewLine()

	var xAxis []string
	var speeds []float64

	max := 0.0
	for _, result := range results {
		xAxis = append(xAxis, result.Timestamp.Format("2006-01-02 15:04:05"))
		speeds = append(speeds, result.AverageMbps)
		if result.AverageMbps > max {
			max = result.AverageMbps
		}
	}

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Download Speeds Over Time",
			Subtitle: fmt.Sprintf("Test data from: %v Days", len(results)),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Trigger:   "axis",
			Formatter: "{b}: {c} Mbps",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: opts.Bool(true),
			Max:        float32(max + 25),
			Min:        0,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#a50026", "#d73027", "#f46d43", "#fdae61", "#fee090",
					"#ffffbf", "#e0f3f8", "#abd9e9", "#74add1", "#4575b4", "#313695"},
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:         "Speed (Mbps)",
			NameLocation: "middle",
			NameGap:      35,
			AxisLabel: &opts.AxisLabel{
				Show:         opts.Bool(true),
				Rotate:       90,
				ShowMaxLabel: opts.Bool(true),
			},
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
		AddSeries("Download Speed (Mbps)", generateBarItems(speeds)).
		SetSeriesOptions(
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:   "#4169E1",
				Opacity: 1,
			}),
		)

	line.AddSeries("Download Speed (Mbps)", generateLineItems(speeds)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(false)}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Type:  "solid",
				Width: 2,
				Color: "#FF4500",
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color: "#FF4500",
			}),
		)

	bar.Overlap(line)

	return bar, nil
}
