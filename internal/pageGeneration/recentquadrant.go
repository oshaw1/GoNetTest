package pageGeneration

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/dataManagment"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

var (
	recentDataTemplate     *template.Template
	chartTemplate          *template.Template
	recentQuadrantTemplate *template.Template
)

func init() {
	// Parse all templates together
	templates, err := template.ParseGlob("internal/pageGeneration/templates/*.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template files: %v", err)
	}

	recentDataTemplate = templates.Lookup("recentData.tmpl")
	chartTemplate = templates.Lookup("chart.tmpl")
	recentQuadrantTemplate = templates.Lookup("recentquadrant.tmpl")

	if recentDataTemplate == nil || chartTemplate == nil || recentQuadrantTemplate == nil {
		log.Fatal("One or more templates not found after parsing")
	}

	log.Println("Templates parsed successfully")
}

func GenerateRecentQuadrantHTML(result *networkTesting.ICMBTestResult) (template.HTML, error) {
	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	chartHTML, err := generateChartHTML(conf.RecentDays)
	if err != nil {
		return "", fmt.Errorf("failed to generate chart html: %v", err)
	}

	dataHTML, err := generateDataSectionHTML(result, conf.RecentDays)
	if err != nil {
		return "", fmt.Errorf("failed to generate data section html: %v", err)
	}

	var buf bytes.Buffer
	data := struct {
		ChartSection template.HTML
		DataSection  template.HTML
	}{
		ChartSection: template.HTML(chartHTML),
		DataSection:  template.HTML(dataHTML),
	}

	err = recentQuadrantTemplate.ExecuteTemplate(&buf, "recentQuadrant", data)
	if err != nil {
		return "", fmt.Errorf("error executing full quadrant template: %v", err)
	}

	return template.HTML(buf.String()), nil
}

func generateChartHTML(daysBack int) (template.HTML, error) {
	dataExists, imagePath, err := dataManagment.CheckForRecentTestData("data/output", daysBack, ".jpg")
	if err != nil {
		return "", fmt.Errorf("failed to check recent test data: %v", err)
	}

	imagePath = filepath.ToSlash(imagePath)
	imagePath = strings.TrimPrefix(imagePath, "data/output/")

	var buf bytes.Buffer
	data := struct {
		HasData        bool
		ChartImagePath string
	}{
		HasData:        dataExists,
		ChartImagePath: imagePath,
	}

	err = chartTemplate.ExecuteTemplate(&buf, "chartSection", data)
	if err != nil {
		return "", fmt.Errorf("error executing chart template: %v", err)
	}
	return template.HTML(buf.String()), nil
}

func generateDataSectionHTML(result *networkTesting.ICMBTestResult, daysBack int) (template.HTML, error) {
	dataExists, _, err := dataManagment.CheckForRecentTestData("data/output", daysBack, ".json")
	if err != nil {
		log.Printf("Error checking for recent test data: %v", err)
		return "", fmt.Errorf("failed to check recent test data: %v", err)
	}

	var buf bytes.Buffer
	data := struct {
		HasData        bool
		ICMBTestResult *networkTesting.ICMBTestResult
		LossPercentage float64
	}{
		HasData:        dataExists,
		ICMBTestResult: result,
		LossPercentage: lossPercentage(result),
	}

	err = recentDataTemplate.ExecuteTemplate(&buf, "dataSection", data)
	if err != nil {
		log.Printf("Error executing data section template: %v", err)
		return "", fmt.Errorf("error executing data section template: %v", err)
	}

	return template.HTML(buf.String()), nil
}

func lossPercentage(tr *networkTesting.ICMBTestResult) float64 {
	if tr == nil || tr.Sent == 0 {
		return 0
	}
	return float64(tr.Lost) / float64(tr.Sent) * 100
}

func ReturnNoDataHTML() (template.HTML, error) {
	return template.HTML("<p>No recent test data available.</p>"), nil
}
