package pageGeneration

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/oshaw1/go-net-test/internal/dataManagment"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

var (
	recentDataTemplate     *template.Template
	chartTemplate          *template.Template
	recentQuadrantTemplate *template.Template
)

func InitTemplates() {
	// Parse all templates together
	templates, err := template.ParseGlob("internal/pageGeneration/templates/*.tmpl")
	if err != nil {
		log.Panicf("Error parsing template files: %v", err)
	}

	recentDataTemplate = templates.Lookup("recentData.tmpl")
	chartTemplate = templates.Lookup("chart.tmpl")
	recentQuadrantTemplate = templates.Lookup("recentquadrant.tmpl")

	if recentDataTemplate == nil || chartTemplate == nil || recentQuadrantTemplate == nil {
		log.Fatal("One or more templates not found after parsing")
	}

	log.Println("Templates parsed successfully")
}

func GenerateRecentQuadrantHTML(results []interface{}) (template.HTML, error) {
	var chartSections []template.HTML
	var dataSections []template.HTML
	for _, result := range results {
		chartHTML, err := returnChartSectionHTML(result)
		if err != nil {
			log.Printf("Error generating chart section for result type %T: %v", result, err)
			continue
		}
		dataHTML, err := returnDataSectionHTML(result)
		if err != nil {
			log.Printf("Error generating data section for result type %T: %v", result, err)
			continue
		}
		chartSections = append(chartSections, chartHTML)
		dataSections = append(dataSections, dataHTML)
	}

	if len(dataSections) == 0 {
		return "", fmt.Errorf("no valid data sections generated")
	}

	var buf bytes.Buffer
	sections := struct {
		ChartSection template.HTML
		DataSection  template.HTML
	}{
		ChartSection: template.HTML(strings.Join(convertHTMLSliceToStringSlice(chartSections), "")),
		DataSection:  template.HTML(strings.Join(convertHTMLSliceToStringSlice(dataSections), "")),
	}

	err := recentQuadrantTemplate.ExecuteTemplate(&buf, "recentQuadrant", sections)
	if err != nil {
		return "", fmt.Errorf("error executing full quadrant template: %v", err)
	}

	return template.HTML(buf.String()), nil
}

func convertHTMLSliceToStringSlice(htmlSlice []template.HTML) []string {
	stringSlice := make([]string, len(htmlSlice))
	for i, html := range htmlSlice {
		stringSlice[i] = string(html)
	}
	return stringSlice
}

func returnChartSectionHTML(result interface{}) (template.HTML, error) {
	var chartHtml template.HTML
	var err error

	switch result.(type) {
	case *networkTesting.ICMPTestResult:
		chartHtml, err = generateICMPChartHTML()
	// case *networkTesting.TCPTestResult:
	//     dataHTML, err = generateTCPDataSectionHTML(v)
	// case *networkTesting.UDPTestResult:
	//     dataHTML, err = generateUDPDataSectionHTML(v)
	// Add more cases for other test types as needed
	default:
		return "", fmt.Errorf("unsupported test result type: %T", result)
	}

	if err != nil {
		return "", fmt.Errorf("failed to generate data section html: %v", err)
	}

	return chartHtml, nil
}

func generateICMPChartHTML() (template.HTML, error) {
	dataExists, filePath, err := dataManagment.ReturnRecentTestDataPath("data/output", "icmp", ".html")
	if err != nil {
		return "", fmt.Errorf("failed to check recent test data: %v", err)
	}

	if !dataExists {
		var buf bytes.Buffer
		data := struct {
			HasData   bool
			ChartHTML template.HTML
		}{
			HasData: false,
		}
		err = chartTemplate.ExecuteTemplate(&buf, "chartSection", data)
		if err != nil {
			return "", fmt.Errorf("error executing chart template: %v", err)
		}
		return template.HTML(buf.String()), nil
	}

	// Read the chart HTML file
	chartContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read chart file: %v", err)
	}

	// Prepare the data for the template
	var buf bytes.Buffer
	data := struct {
		HasData   bool
		ChartHTML template.HTML
	}{
		HasData:   true,
		ChartHTML: template.HTML(chartContent),
	}

	// Execute the template
	err = chartTemplate.ExecuteTemplate(&buf, "chartSection", data)
	if err != nil {
		return "", fmt.Errorf("error executing chart template: %v", err)
	}

	return template.HTML(buf.String()), nil
}

func returnDataSectionHTML(result interface{}) (template.HTML, error) {
	var dataHTML template.HTML
	var err error

	switch v := result.(type) {
	case *networkTesting.ICMPTestResult:
		dataHTML, err = generateICMPDataSectionHTML(v)
	// case *networkTesting.TCPTestResult:
	// 	dataHTML, err = generateTCPDataSectionHTML(v)
	// case *networkTesting.UDPTestResult:
	// 	dataHTML, err = generateUDPDataSectionHTML(v)
	// Add more cases for other test types as needed
	default:
		return "", fmt.Errorf("unsupported test result type: %T", result)
	}

	if err != nil {
		return "", fmt.Errorf("failed to generate data section html: %v", err)
	}

	return dataHTML, nil
}

func generateICMPDataSectionHTML(result *networkTesting.ICMPTestResult) (template.HTML, error) {
	dataExists, _, err := dataManagment.ReturnRecentTestDataPath("data/output", "icmp", ".json")
	if err != nil {
		log.Printf("Error checking for recent test data: %v", err)
		return "", fmt.Errorf("failed to check recent test data: %v", err)
	}

	var buf bytes.Buffer
	data := struct {
		HasData        bool
		ICMPTestResult *networkTesting.ICMPTestResult
		LossPercentage float64
	}{
		HasData:        dataExists,
		ICMPTestResult: result,
		LossPercentage: lossPercentage(result),
	}

	err = recentDataTemplate.ExecuteTemplate(&buf, "dataSection", data)
	if err != nil {
		log.Printf("Error executing data section template: %v", err)
		return "", fmt.Errorf("error executing data section template: %v", err)
	}

	return template.HTML(buf.String()), nil
}

func lossPercentage(tr *networkTesting.ICMPTestResult) float64 {
	if tr == nil || tr.Sent == 0 {
		return 0
	}
	return float64(tr.Lost) / float64(tr.Sent) * 100
}

func ReturnNoDataHTML() (template.HTML, error) {
	return template.HTML("<p>No recent test data available.</p>"), nil
}
