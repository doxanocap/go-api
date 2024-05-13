package service

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/service/websocket"
	"go.uber.org/zap"
)

type Service struct {
	config  *models.Config
	manager interfaces.IManager

	workspace interfaces.IWorkspaceService
	websocket interfaces.IWebsocketService
}

func InitService(manager interfaces.IManager, config *models.Config, log *zap.Logger) *Service {
	return &Service{
		config:  config,
		manager: manager,

		workspace: InitWorkspaceService(manager, log),
		websocket: websocket.NewWebsocketService(manager, config, log.Named("[WS]")),
	}
}

func (s *Service) Workspace() interfaces.IWorkspaceService {
	return s.workspace
}

func (s *Service) Websocket() interfaces.IWebsocketService {
	return s.websocket
}
