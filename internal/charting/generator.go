package charting

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

// ChartGenerator defines the interface for generating charts
type ChartGenerator interface {
	// ICMP Charts
	GenerateICMPDistributionPie(result *networkTesting.ICMPTestResult) (*charts.Pie, error)
}

type Generator struct{}

func NewGenerator() *Generator { // Return concrete type
	return &Generator{}
}
