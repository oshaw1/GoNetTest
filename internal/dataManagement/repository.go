package dataManagement

import (
	"database/sql"

	"github.com/oshaw1/go-net-test/config"
)

const dateFormat = "2006-01-02"

type Repository struct {
	db     *sql.DB
	config *config.Config
}

func NewRepository(db *sql.DB, config *config.Config) *Repository {
	return &Repository{
		db:     db,
		config: config,
	}
}
