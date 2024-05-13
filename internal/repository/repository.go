package repository

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/repository/pg"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	config *models.Config
	cache  interfaces.ICacheProcessor

	workspace        interfaces.IWorkspaceRepository
	workspaceUsers   interfaces.IWorkspaceUsersRepository
	workspaceMessage interfaces.IWorkspaceMessagesRepository

	workspaceBodyCache interfaces.IWorkspaceBodyCache
}

func InitRepository(db *sqlx.DB, config *models.Config, cacheProcessor interfaces.ICacheProcessor) *Repository {
	return &Repository{
		config: config,

		workspace:        pg.InitWorkspaceRepository(db),
		workspaceUsers:   pg.InitWorkspaceUsersRepository(db),
		workspaceMessage: pg.InitWorkspaceMessagesRepository(db),
	}
}

func (r *Repository) Workspace() interfaces.IWorkspaceRepository {
	return r.workspace
}

func (r *Repository) WorkspaceUsers() interfaces.IWorkspaceUsersRepository {
	return r.workspaceUsers
}

func (r *Repository) WorkspaceMessage() interfaces.IWorkspaceMessagesRepository {
	return r.workspaceMessage
}

func (r *Repository) WorkspaceBodyCache() interfaces.IWorkspaceBodyCache {
	return r.workspaceBodyCache
}
