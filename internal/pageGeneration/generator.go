package pageGeneration

import (
	"fmt"
	"html/template"
	"log"

	"github.com/oshaw1/go-net-test/internal/dataManagement"
)

// PageGenerator handles HTML page generation
type PageGenerator struct {
	templates  *template.Template
	repository *dataManagement.Repository
	logger     *log.Logger
}

var requiredTemplates = []string{
	"recentData.tmpl",
	"chart.tmpl",
	"recentQuadrant.tmpl",
}

// NewPageGenerator creates a new instance of PageGenerator with validation
func NewPageGenerator(templatePath string, repo *dataManagement.Repository) (*PageGenerator, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository cannot be nil")
	}

	templates, err := template.ParseGlob(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	if err := validateRequiredTemplates(templates); err != nil {
		return nil, err
	}

	return &PageGenerator{
		templates:  templates,
		repository: repo,
	}, nil
}

// validateRequiredTemplates ensures all required templates are present
func validateRequiredTemplates(templates *template.Template) error {
	for _, name := range requiredTemplates {
		if templates.Lookup(name) == nil {
			return fmt.Errorf("required template not found: %s", name)
		}
	}
	return nil
}
