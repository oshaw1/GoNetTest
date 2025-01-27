package pageGeneration

import (
	"fmt"
	"html/template"

	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Repository interface {
	GetTestDirectories() ([]string, error)
	ListTestTypesInDateDir(date string) ([]string, error)
	MapTestFilesByTimestamp(date, testType string) (map[string][]string, error)
}

type PageGenerator struct {
	repository Repository
	templates  *template.Template
	caser      cases.Caser
}

var requiredTemplates = []string{
	"base.gohtml",
	"control_quadrant.gohtml",
	"generate_quadrant.gohtml",
	"scheduler_quadrant.gohtml",
	"test_quadrant.gohtml",
}

func NewPageGenerator(templatePath string, repo *dataManagement.Repository) (*PageGenerator, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository cannot be nil")
	}
	templates, err := template.ParseGlob(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}
	pg := &PageGenerator{
		templates:  templates,
		repository: repo,
		caser:      cases.Title(language.English),
	}
	if err := pg.validateRequiredTemplates(templates); err != nil {
		return nil, err
	}
	return pg, nil
}

func (pg *PageGenerator) validateRequiredTemplates(templates *template.Template) error {
	for _, name := range requiredTemplates {
		if templates.Lookup(name) == nil {
			return fmt.Errorf("required template not found: %s", name)
		}
	}
	return nil
}
