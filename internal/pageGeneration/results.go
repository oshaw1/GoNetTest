package pageGeneration

import (
	"fmt"
	"html/template"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

type ICMPData struct {
	PageData
	ICMPTestResult *networkTesting.ICMPTestResult
	LossPercentage float64
}

type SpeedData struct {
	PageData               PageData
	AverageSpeedTestResult *networkTesting.AverageSpeedTestResult
}

func (pg *PageGenerator) GenerateICMPDataHTML(result *networkTesting.ICMPTestResult) (template.HTML, error) {
	fmt.Printf("GenerateICMPDataHTML called with result: %+v\n", result)

	data := ICMPData{
		PageData: PageData{
			HasData: result != nil,
		},
		ICMPTestResult: result,
		LossPercentage: calculateLossPercentage(result),
	}

	html, err := pg.executeTemplate("dataSection", data)
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return html, nil
}

func (pg *PageGenerator) GenerateSpeedDataHTML(result *networkTesting.AverageSpeedTestResult) (template.HTML, error) {
	fmt.Printf("GenerateSpeedDataHTML called with result: %+v\n", result)

	data := SpeedData{
		PageData: PageData{
			HasData: result != nil,
		},
		AverageSpeedTestResult: result,
	}

	html, err := pg.executeTemplate("dataSection", data)
	if err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return html, nil
}

func calculateLossPercentage(tr *networkTesting.ICMPTestResult) float64 {
	if tr == nil || tr.Sent == 0 {
		return 0
	}
	return float64(tr.Lost) / float64(tr.Sent) * 100
}
