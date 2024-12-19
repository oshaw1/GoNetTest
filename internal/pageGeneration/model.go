package pageGeneration

import (
	"html/template"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

const (
	DefaultDaysRange = 7
	NoDataMessage    = "<p>No recent test data available.</p>"
)

var requiredTemplates = []string{
	"recentData.tmpl",
	"chart.tmpl",
	"recentQuadrant.tmpl",
}

type PageData struct {
	HasData bool
	Content template.HTML
}

// ChartData represents the common structure for chart-specific data
type ChartData struct {
	PageData
	ChartHTML template.HTML
}

// ICMPData represents the structure for ICMP-specific test results
type ICMPData struct {
	PageData
	ICMPTestResult *networkTesting.ICMPTestResult
	LossPercentage float64
}
