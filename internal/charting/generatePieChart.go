package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func IcmpPieChart(result networkTesting.ICMPTestResult) (*charts.Pie, error) {
	pie := charts.NewPie()

	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "ICMP Packet Distribution",
		Subtitle: fmt.Sprintf("Test ran at: %v", result.Timestamp),
	}))

	// Add data to the pie chart
	pie.AddSeries("Packet Distribution", []opts.PieData{
		{Name: "Received", Value: result.Received},
		{Name: "Lost", Value: result.Lost},
	}).SetSeriesOptions(charts.WithLabelOpts(opts.Label{
		Formatter: "{b}: {c} ({d}%)",
	}))

	return pie, nil
}

// func GetIcmpPieHTML(pie *charts.Pie) template.HTML {
// 	pie.Renderer = newSnippetRenderer(pie, pie.Validate)
// 	return renderToHtml(pie)
// }
