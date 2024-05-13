package interfaces

import (
	"auth-api/internal/models"
	"context"
)

type IRepository interface {
	Workspace() IWorkspaceRepository
	WorkspaceUsers() IWorkspaceUsersRepository
	WorkspaceMessage() IWorkspaceMessagesRepository
}

type IWorkspaceRepository interface {
	Create(ctx context.Context, w *models.Workspace) (err error)
	GetAll(ctx context.Context) (ws []models.Workspace, err error)
	FindByID(ctx context.Context, id int64) (w *models.Workspace, err error)
	FindByIDCode(ctx context.Context, idcode string) (w *models.Workspace, err error)
	UpdateByID(ctx context.Context, w *models.Workspace) (err error)
	UpdateBodyByID(ctx context.Context, w *models.Workspace) (err error)
	DeleteByID(ctx context.Context, id int64) (err error)
}

type IWorkspaceUsersRepository interface {
	Create(ctx context.Context, wu *models.WorkspaceUser) (err error)
	GetAll(ctx context.Context) (ws []models.WorkspaceUser, err error)
	GetAllByWorkspaceID(ctx context.Context, workspaceID int64) (ws []models.WorkspaceUser, err error)
	GetAllByUserID(ctx context.Context, userID string) (w []models.Workspace, err error)
	SetUserPresent(ctx context.Context, workspaceID int64, userID string) (err error)
	SetUserAbsent(ctx context.Context, workspaceID int64, userID string) (err error)
	DeleteByID(ctx context.Context, id int64) (err error)
}

type IWorkspaceMessagesRepository interface {
	Create(ctx context.Context, wm *models.WorkspaceMessage) (err error)
	GetAllByWorkspaceID(ctx context.Context, workspaceID int64) (ws []models.WorkspaceMessage, err error)
	DeleteByWorkspaceID(ctx context.Context, id int64) (err error)
}

type IWorkspaceBodyCache interface {
	Get(ctx context.Context, workspaceID int64) (value string, err error)
	Set(ctx context.Context, workspaceID int64, body string) (err error)
}
