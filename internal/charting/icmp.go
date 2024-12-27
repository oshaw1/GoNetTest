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
