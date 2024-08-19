package pageGeneration

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

var noDataTemplate *template.Template

func init() {
	var err error
	noDataTemplate, err = template.ParseFiles("internal/pageGeneration/templates/norecentdata.tmpl")
	if err != nil {
		log.Printf("Error parsing template file: %v", err)
		panic(err)
	}
	log.Println("Template parsed successfully")
}

func ReturnNoDataHTML() (template.HTML, error) {
	var buf bytes.Buffer
	err := noDataTemplate.ExecuteTemplate(&buf, "testResult", nil)
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
