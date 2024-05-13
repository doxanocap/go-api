package interfaces

import (
	"auth-api/internal/models"
	"context"
	"github.com/gorilla/websocket"
)

type IService interface {
	Workspace() IWorkspaceService
	Websocket() IWebsocketService
}

type IWorkspaceService interface {
	Create(ctx context.Context, editorID string) (*models.Workspace, error)
	GetByIDCode(ctx context.Context, workspaceIDCode string) (*models.Workspace, error)
	GetAllByUserID(ctx context.Context, userID string) ([]models.Workspace, error)
	HandleNewClient(ctx context.Context, conn *websocket.Conn, clientID, workspaceID string) error
}

type IWebsocketService interface {
	HandleNewClient(clientID, workspaceID string, conn *websocket.Conn) error
	GetByIDCode(workspaceID string) *models.Workspace
	WriteConnErr(conn *websocket.Conn, err error)
}

type IWebsocketPoolService interface {
	SetWorkspaces(m map[string]*models.Workspace)
	SetWorkspaceClients(m map[string][]string)

	Register(userID, workspaceID string, ch chan []byte) error
	Unregister(userID string)
	Send(message *models.Message)
	Handle()

	GetByID(workspaceID string) *models.Workspace
}
