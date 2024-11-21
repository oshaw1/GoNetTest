package charting

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (g *Generator) GenerateTCPStatusPie(result *networkTesting.TCPTestResult) (*charts.Pie, error) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "TCP Connection Status",
		Subtitle: fmt.Sprintf("Test ran at: %v", result.Timestamp),
	}))

	var successful, failed int
	for _, port := range result.PortResults {
		if port.Connected {
			successful++
		} else {
			failed++
		}
	}

	pie.AddSeries("Connections", []opts.PieData{
		{Name: "Successful", Value: successful},
		{Name: "Failed", Value: failed},
	})

	return pie, nil
}

// GenerateTCPTimesBar creates a bar chart showing connection times
func (g *Generator) GenerateTCPTimesBar(result *networkTesting.TCPTestResult) (*charts.Bar, error) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "TCP Connection Times",
		Subtitle: fmt.Sprintf("Test ran at: %v", result.Timestamp),
	}))

	var ports []string
	var times []opts.BarData

	for _, pr := range result.PortResults {
		if pr.Connected {
			ports = append(ports, fmt.Sprintf("Port %d", pr.Port))
			times = append(times, opts.BarData{Value: pr.ConnectTime.Milliseconds()})
		}
	}

	bar.SetXAxis(ports)
	bar.AddSeries("Time (ms)", times)

	return bar, nil
}
