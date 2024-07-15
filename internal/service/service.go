package service

import (
	"auth-api/config"
	"auth-api/internal/interfaces"
	"auth-api/logger"
)

type Service struct {
	config  *config.Config
	manager interfaces.IManager
}

func InitService(manager interfaces.IManager, config *config.Config, log *logger.Logger) *Service {
	return &Service{
		config:  config,
		manager: manager,
	}
}
