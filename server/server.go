package server

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/server/rest"
	"go.uber.org/zap"
	"sync"
)

type Server struct {
	log     *zap.Logger
	config  *models.Config
	manager interfaces.IManager

	restServer       interfaces.IRESTServer
	restServerRunner sync.Once
}

func InitServer(manager interfaces.IManager, config *models.Config, log *zap.Logger) *Server {
	return &Server{
		log:     log,
		config:  config,
		manager: manager,

		restServer: rest.InitREST(config, manager, log.Named("[REST]")),
	}
}

func (p *Server) REST() interfaces.IRESTServer {
	return p.restServer
}
