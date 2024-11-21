package pageGeneration

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func (pg *PageGenerator) GenerateRecentQuadrantHTML() (template.HTML, error) {
	var chartSections []template.HTML
	var dataSections []template.HTML

	icmpResult, err := pg.repository.GetRecentICMPTestResult("icmp")
	if err != nil {
		log.Printf("Error retrieving recent ICMP results: %v", err)
	} else {
		log.Printf("Retrieved ICMP result: %+v", icmpResult)
	}

	chartHTML, err := pg.generateICMPChartHTML()
	if err != nil {
		log.Printf("Error generating ICMP chart: %v", err)
	} else {
		log.Printf("Generated ICMP chart HTML: %s", chartHTML)
		chartSections = append(chartSections, chartHTML)
	}

	dataHTML, err := pg.generateICMPDataHTML(icmpResult)
	if err != nil {
		log.Printf("Error generating ICMP data: %v", err)
	} else {
		log.Printf("Generated ICMP data HTML: %s", dataHTML)
		dataSections = append(dataSections, dataHTML)
	}

	if len(dataSections) == 0 {
		log.Println("No data sections, returning blank HTML")
		return ReturnNoDataHTML()
	}

	quadrantHTML, err := pg.renderQuadrant(chartSections, dataSections)
	if err != nil {
		log.Printf("Error rendering quadrant: %v", err)
		return "", err
	}

	log.Printf("Generated quadrant HTML: %s", quadrantHTML)
	return quadrantHTML, nil
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

func (pg *PageGenerator) generateICMPChartHTML() (template.HTML, error) {
	dataExists, filePath, err := pg.repository.CheckRecentData("icmp", ".html")
	if err != nil {
		return "", fmt.Errorf("failed to check recent test data: %w", err)
	}

	log.Printf("ICMP chart data exists: %t, file path: %s", dataExists, filePath)

	var chartContent []byte
	if dataExists {
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
		HasData:   dataExists,
		ChartHTML: template.HTML(chartContent),
	}

	if err := pg.templates.ExecuteTemplate(&buf, "chartSection", data); err != nil {
		return "", fmt.Errorf("error executing chart template: %w", err)
	}

	return template.HTML(buf.String()), nil
}

func (pg *PageGenerator) generateICMPDataHTML(icmpResult *networkTesting.ICMPTestResult) (template.HTML, error) {
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
