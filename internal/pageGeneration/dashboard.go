package pageGeneration

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (pg *PageGenerator) GenerateRecentQuadrantHTML() (template.HTML, error) {
	var chartSections []template.HTML
	var dataSections []template.HTML

	// Get data for last 7 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	// Get most recent test data
	results, err := pg.repository.GetTestDataInRange(startDate, endDate, "icmp")
	if err != nil {
		log.Printf("Error retrieving ICMP results: %v", err)
	} else if len(results) > 0 {
		icmpResult := results[0]
		log.Printf("Retrieved ICMP result: %+v", icmpResult)

		dataHTML, err := pg.GenerateICMPDataHTML(icmpResult)
		if err != nil {
			log.Printf("Error generating ICMP data: %v", err)
		} else {
			log.Printf("Generated ICMP data HTML: %s", dataHTML)
			dataSections = append(dataSections, dataHTML)
		}
	}

	exists, filePath, err := pg.repository.GetChartInRange(startDate, endDate, "icmp")
	if err != nil {
		log.Printf("Error retrieving chart in range: %v", err)
	} else if exists {
		chartHTML, err := pg.generateChartHTMLFromFile(filePath)
		if err == nil {
			log.Printf("Generated ICMP chart HTML")
			chartSections = append(chartSections, chartHTML)
		} else {
			log.Printf("Error generating chart HTML: %v", err)
		}
	}

	if len(chartSections) == 0 {
		log.Println("No successful chart generations in date range")
		return ReturnNoDataHTML()
	}

	quadrantHTML, err := pg.renderQuadrant(chartSections, dataSections)
	if err != nil {
		log.Printf("Error rendering quadrant: %v", err)
		return "", err
	}

	return quadrantHTML, nil
}

func (pg *PageGenerator) generateChartHTMLFromFile(filePath string) (template.HTML, error) {
	chartContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read chart file: %w", err)
	}

	var buf bytes.Buffer
	data := struct {
		HasData   bool
		ChartHTML template.HTML
	}{
		HasData:   true,
		ChartHTML: template.HTML(chartContent),
	}

	if err := pg.templates.ExecuteTemplate(&buf, "chartSection", data); err != nil {
		return "", fmt.Errorf("error executing chart template: %w", err)
	}

	return template.HTML(buf.String()), nil
}

func (pg *PageGenerator) renderQuadrant(chartSections, dataSections []template.HTML) (template.HTML, error) {
	var buf bytes.Buffer
	sections := struct {
		ChartSection template.HTML
		DataSection  template.HTML
	}{
		ChartSection: template.HTML(strings.Join(convertToStringSlice(chartSections), "")),
		DataSection:  template.HTML(strings.Join(convertToStringSlice(dataSections), "")),
	}

	if err := pg.templates.ExecuteTemplate(&buf, "recentQuadrant", sections); err != nil {
		return "", fmt.Errorf("error executing quadrant template: %w", err)
	}

	return template.HTML(buf.String()), nil
}

func (pg *PageGenerator) GenerateICMPChartHTML(date string) (template.HTML, error) {
	exists, filePath, err := pg.repository.CheckData(date, "icmp", ".html")
	if err != nil {
		return "", fmt.Errorf("failed to check test data for date %s: %w", date, err)
	}

	log.Printf("ICMP chart data exists: %t, file path: %s", exists, filePath)

	var chartContent []byte
	if exists {
		chartContent, err = os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read chart file: %w", err)
		}
	}

	var buf bytes.Buffer
	data := struct {
		HasData   bool
		ChartHTML template.HTML
	}{
		HasData:   exists,
		ChartHTML: template.HTML(chartContent),
	}

	if err := pg.templates.ExecuteTemplate(&buf, "chartSection", data); err != nil {
		return "", fmt.Errorf("error executing chart template: %w", err)
	}

	return template.HTML(buf.String()), nil
}

func (pg *PageGenerator) GenerateICMPDataHTML(icmpResult *networkTesting.ICMPTestResult) (template.HTML, error) {
	log.Printf("Generating ICMP data HTML with result: %+v", icmpResult)

	var lossPercentage float64
	if icmpResult != nil {
		lossPercentage = icmpLossPercentage(icmpResult)
	}

	var buf bytes.Buffer
	data := struct {
		HasData        bool
		ICMPTestResult *networkTesting.ICMPTestResult
		LossPercentage float64
	}{
		HasData:        icmpResult != nil,
		ICMPTestResult: icmpResult,
		LossPercentage: lossPercentage,
	}

	if err := pg.templates.ExecuteTemplate(&buf, "dataSection", data); err != nil {
		return "", fmt.Errorf("error executing data template: %w", err)
	}

	renderedHTML := buf.String()
	log.Printf("Rendered ICMP data HTML: %s", renderedHTML)

	return template.HTML(renderedHTML), nil
}
