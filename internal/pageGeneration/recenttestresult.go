package pageGeneration

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

var testResultTemplate *template.Template

func init() {
	var err error
	testResultTemplate, err = template.ParseFiles("internal/pageGeneration/templates/recenttestresult.tmpl")
	if err != nil {
		log.Printf("Error parsing template file: %v", err)
		panic(err)
	}
	log.Println("Template parsed successfully")
}

func GenerateRecentTestResultHTML(result *networkTesting.ICMBTestResult) (template.HTML, error) {
	var buf bytes.Buffer

	data := struct {
		networkTesting.ICMBTestResult
		LossPercentage float64
	}{
		ICMBTestResult: *result,
		LossPercentage: lossPercentage(result),
	}

	err := testResultTemplate.ExecuteTemplate(&buf, "testResult", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return "", err
	}

	html := template.HTML(buf.String())
	log.Printf("Generated HTML: %s", html)

	if html == "" {
		return "", fmt.Errorf("generated HTML is empty")
	}

	return html, nil
}

func lossPercentage(tr *networkTesting.ICMBTestResult) float64 {
	if tr.Sent == 0 {
		return 0
	}
	return float64(tr.Lost) / float64(tr.Sent) * 100
}
