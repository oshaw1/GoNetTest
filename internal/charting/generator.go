package charting

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

type ChartGenerator interface {
	GenerateICMPAnalysisCharts(result *networkTesting.ICMPTestResult) (*charts.Pie, error)
	GenerateJitterAnalysisCharts(result *networkTesting.JitterTestResult) (*charts.Line, error)
	GenerateRouteAnalysisCharts(result *networkTesting.RouteTestResult) (*charts.Line, error)
	GenerateDownloadAnalysisCharts(result *networkTesting.AverageSpeedTestResult) (*charts.Line, error)
	GenerateUploadAnalysisCharts(result *networkTesting.AverageSpeedTestResult) (*charts.Line, error)
	GenerateBandwidthAnalysisCharts(result *networkTesting.BandwidthTestResult) (*charts.Line, error)
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}
