package apis

import (
	"auth-api/config"
	"auth-api/logger"
)

type APIs struct {
}

func NewAPIsProcessor(config *config.Config, log *logger.Logger) *APIs {
	return &APIs{}
}
