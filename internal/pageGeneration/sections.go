package pageGeneration

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

func (pg *PageGenerator) GenerateRecentQuadrantHTML() (template.HTML, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -DefaultDaysRange)

	charts, data, err := pg.generateSections(startDate, endDate, "icmp")
	if err != nil {
		pg.logger.Printf("Error generating sections: %v", err)
		return template.HTML(NoDataMessage), nil
	}

	if len(charts) == 0 {
		pg.logger.Println("No charts available in date range")
		return template.HTML(NoDataMessage), nil
	}

	return pg.renderRecentQuadrant(charts, data)
}

func (pg *PageGenerator) generateSections(startDate, endDate time.Time, testType string) ([]template.HTML, []template.HTML, error) {
	results, err := pg.repository.GetTestDataInRange(startDate, endDate, testType)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving test data: %w", err)
	}

	var chartSections, dataSections []template.HTML

	// Generate data sections
	if len(results) > 0 {
		var dataHTML template.HTML
		var err error

		switch testType {
		case "icmp":
			dataHTML, err = pg.GenerateICMPDataHTML(results[0].ICMP)
		case "download":
			dataHTML, err = pg.GenerateSpeedDataHTML(results[0].Download)
		case "upload":
			dataHTML, err = pg.GenerateSpeedDataHTML(results[0].Upload)
		}

		if err == nil {
			dataSections = append(dataSections, dataHTML)
		}
	}

	// Generate chart sections
	exists, filePath, err := pg.repository.GetChartInRange(startDate, endDate, testType)
	if err == nil && exists {
		chartHTML, err := pg.generateChartHTMLFromFile(filePath)
		if err == nil {
			chartSections = append(chartSections, chartHTML)
		}
	}

	return chartSections, dataSections, nil
}

func (pg *PageGenerator) renderRecentQuadrant(charts, data []template.HTML) (template.HTML, error) {
	sections := struct {
		ChartSection template.HTML
		DataSection  template.HTML
	}{
		ChartSection: template.HTML(strings.Join(htmlSliceToStrings(charts), "")),
		DataSection:  template.HTML(strings.Join(htmlSliceToStrings(data), "")),
	}

	return pg.executeTemplate("recentQuadrant", sections)
}
