package dataManagement

import (
	"github.com/oshaw1/go-net-test/config"
)

const dateFormat = "2006-01-02"

type Repository struct {
	baseDir string
	config  *config.Config
}

func NewRepository(baseDir string, config *config.Config) *Repository {
	return &Repository{
		baseDir: baseDir,
		config:  config,
	}
}
