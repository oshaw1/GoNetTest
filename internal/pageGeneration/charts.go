package pageGeneration

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
)

func (pg *PageGenerator) generateChartHTMLFromFile(filePath string) (template.HTML, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read chart file: %w", err)
	}

	data := ChartData{
		PageData: PageData{
			HasData: true,
		},
		ChartHTML: template.HTML(content),
	}

	return pg.executeTemplate("chartSection", data)
}

func (pg *PageGenerator) GenerateICMPChartHTMLGivenDate(date string, test string) (template.HTML, error) {
	exists, filePath, err := pg.repository.CheckData(date, test, ".html")
	if err != nil {
		return "", fmt.Errorf("failed to check test data for date %s: %w", date, err)
	}

	log.Printf("ICMP chart data exists: %t, file path: %s", exists, filePath)

	var content []byte
	if exists {
		content, err = os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read chart file: %w", err)
		}
	}

	var buf bytes.Buffer
	data := ChartData{
		PageData: PageData{
			HasData: true,
		},
		ChartHTML: template.HTML(content),
	}

	if err := pg.templates.ExecuteTemplate(&buf, "chartSection", data); err != nil {
		return "", fmt.Errorf("error executing chart template: %w", err)
	}

	return template.HTML(buf.String()), nil
}
