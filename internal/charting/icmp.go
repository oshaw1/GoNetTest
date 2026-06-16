package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateICMPAnalysisCharts(result *networkTesting.ICMPTestResult) (*charts.Bar, error) {
	bar, err := generateICMPDistributionBar(result)
	if err != nil {
		return nil, err
	}

	return bar, nil
}

func (g *Generator) GenerateHistoricICMPAnalysisCharts(results []*networkTesting.ICMPTestResult) (*charts.Bar, error) {
	if results == nil {
		return nil, fmt.Errorf("function called with no results")
	}

	barOverTime, err := generateICMPOverTimeBar(results)
	if err != nil {
		return nil, err
	}

	return barOverTime, nil
}

func generateICMPDistributionBar(result *networkTesting.ICMPTestResult) (*charts.Bar, error) {
	bar := charts.NewBar()

	// Configure global options
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "ICMP Packet Distribution",
			Subtitle: fmt.Sprintf("Test ran at: %v", result.Timestamp),
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		// Set color theme globally
		charts.WithColorsOpts(opts.Colors{"#4169E1", "#FF0000"}),
	)

	// Set X axis with categories
	bar.SetXAxis([]string{"Packets"})

	// Add series for Received packets
	bar.AddSeries("Received", []opts.BarData{{Value: result.Received}}).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Position:  "top",
				Formatter: "{c}",
			}),
		)

	// Add series for Lost packets
	bar.AddSeries("Lost", []opts.BarData{{Value: result.Lost}}).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Position:  "top",
				Formatter: "{c}",
			}),
		)

	return bar, nil
}

func generateICMPOverTimeBar(results []*networkTesting.ICMPTestResult) (*charts.Bar, error) {
	bar := charts.NewBar()

	var xAxis []string
	var received []float64
	var lost []float64

	for _, result := range results {
		xAxis = append(xAxis, result.Timestamp.Format("2006-01-02 15:04:05"))
		received = append(received, float64(result.Received))
		lost = append(lost, float64(result.Lost))
	}

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "ICMP Packet Loss Over Time",
			Subtitle: fmt.Sprintf("Test data from: %v Days", len(results)),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    opts.Bool(true),
			Trigger: "axis",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Show:         opts.Bool(true),
				Rotate:       45,
				ShowMaxLabel: opts.Bool(true),
			},
		}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithColorsOpts(opts.Colors{"#4169E1", "#FF0000"}),
		charts.WithGridOpts(opts.Grid{
			Bottom: "20%",
			Top:    "10%",
		}),
	)

	bar.SetXAxis(xAxis).
		AddSeries("Received", generateBarItems(received)).
		AddSeries("Lost", generateBarItems(lost))

	return bar, nil
}
