package pageGeneration

import (
	"fmt"
	"html/template"

	"github.com/oshaw1/go-net-test/internal/dataManagement"
)

type PageGenerator struct {
	templates  *template.Template
	repository *dataManagement.Repository
}

func NewPageGenerator(templatePath string, repo *dataManagement.Repository) (*PageGenerator, error) {
	templates, err := template.ParseGlob(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	required := []string{"recentData.tmpl", "chart.tmpl", "recentQuadrant.tmpl"}
	for _, name := range required {
		if templates.Lookup(name) == nil {
			return nil, fmt.Errorf("required template not found: %s", name)
		}
	}

	return &PageGenerator{
		templates:  templates,
		repository: repo,
	}, nil
}
