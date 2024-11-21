package pageGeneration

import (
	"html/template"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func convertToStringSlice(htmlSlice []template.HTML) []string {
	stringSlice := make([]string, len(htmlSlice))
	for i, html := range htmlSlice {
		stringSlice[i] = string(html)
	}
	return stringSlice
}

func icmpLossPercentage(tr *networkTesting.ICMPTestResult) float64 {
	if tr == nil || tr.Sent == 0 {
		return 0
	}
	return float64(tr.Lost) / float64(tr.Sent) * 100
}

func ReturnNoDataHTML() (template.HTML, error) {
	return template.HTML("<p>No recent test data available.</p>"), nil
}
