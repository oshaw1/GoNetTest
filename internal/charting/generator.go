package charting

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

// ChartGenerator defines the interface for generating charts
type ChartGenerator interface {
	// ICMP Charts
	GenerateICMPDistributionPie(result *networkTesting.ICMPTestResult) (*charts.Pie, error)
	GenerateICMPRTTLine(result *networkTesting.ICMPTestResult) (*charts.Line, error)

	// TCP Charts
	GenerateTCPStatusPie(result *networkTesting.TCPTestResult) (*charts.Pie, error)
	GenerateTCPTimesBar(result *networkTesting.TCPTestResult) (*charts.Bar, error)
}

type Generator struct{}

func NewGenerator() *Generator { // Return concrete type
	return &Generator{}
}
