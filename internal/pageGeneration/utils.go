package pageGeneration

import (
	"bytes"
	"fmt"
	"html/template"
)

func ReturnNoDataHTML() (template.HTML, error) {
	return template.HTML("<p>No recent test data available.</p>"), nil
}

func (pg *PageGenerator) executeTemplate(name string, data interface{}) (template.HTML, error) {
	var buf bytes.Buffer
	if err := pg.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("error executing template %s: %w", name, err)
	}
	return template.HTML(buf.String()), nil
}

func htmlSliceToStrings(htmlSlice []template.HTML) []string {
	result := make([]string, len(htmlSlice))
	for i, html := range htmlSlice {
		result[i] = string(html)
	}
	return result
}
