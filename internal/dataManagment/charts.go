package dataManagment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
)

func SavePieChart(chart *charts.Pie, testType string, testDataFilename string) error {
	dir := filepath.Dir(testDataFilename)
	baseFilename := filepath.Base(testDataFilename)
	chartFilename := testType + "_pie_chart_" + strings.TrimPrefix(baseFilename, "icmp_test_")
	chartFilename = strings.TrimSuffix(chartFilename, filepath.Ext(chartFilename)) + ".html"
	fullPath := filepath.Join(dir, chartFilename)

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	err = chart.Render(f)
	if err != nil {
		return fmt.Errorf("failed to render chart: %w", err)
	}

	fmt.Printf("Chart saved successfully to: %s\n", fullPath)
	return nil
}
