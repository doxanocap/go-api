package server

import (
	"auth-api/config"
	"auth-api/internal/interfaces"
	"auth-api/logger"
	"auth-api/server/rest"
	"sync"
)

type Server struct {
	log     *logger.Logger
	config  *config.Config
	manager interfaces.IManager

	restServer       interfaces.IRESTServer
	restServerRunner sync.Once
}

func InitServer(manager interfaces.IManager, config *config.Config, log *logger.Logger) *Server {
	return &Server{
		log:     log,
		config:  config,
		manager: manager,

		restServer: rest.InitREST(config, manager, log.WithModule("REST")),
	}
}

func (p *Server) REST() interfaces.IRESTServer {
	return p.restServer
}
