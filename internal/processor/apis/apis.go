package apis

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"go.uber.org/zap"
)

type APIs struct {
	authAPI interfaces.IAuthAPIProcessor
}

func NewAPIsProcessor(config *models.Config, log *zap.Logger) *APIs {
	return &APIs{
		authAPI: NewAuthAPIProcessor(config, log.Named("[AUTH_API]")),
	}
}

func (a *APIs) AuthAPI() interfaces.IAuthAPIProcessor {
	return a.authAPI
}
