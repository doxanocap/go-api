package repository

import (
	"auth-api/config"
	"auth-api/internal/interfaces"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	config *config.Config
	cache  interfaces.ICacheProcessor
}

func InitRepository(db *sqlx.DB, config *config.Config, cacheProcessor interfaces.ICacheProcessor) *Repository {
	return &Repository{
		config: config,
	}
}
