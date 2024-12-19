package pageGeneration

import (
	"html/template"
)

const (
	DefaultDaysRange = 7
	NoDataMessage    = "<p>No recent test data available.</p>"
)

type PageData struct {
	HasData bool
	Content template.HTML
}

// ChartData represents the common structure for chart-specific data
type ChartData struct {
	PageData
	ChartHTML template.HTML
}
