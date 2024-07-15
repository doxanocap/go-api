package rest

import (
	"auth-api/config"
	"auth-api/internal/interfaces"
	"auth-api/internal/pkg/metrics"
	"auth-api/logger"
	"auth-api/server/rest/middlewares"
	"context"
	"errors"
	"fmt"
	"github.com/doxanocap/pkg/router"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type REST struct {
	log     *logger.Logger
	config  *config.Config
	router  *gin.Engine
	server  *http.Server
	manager interfaces.IManager

	middlewares *middlewares.Middlewares
}

func InitREST(config *config.Config, manager interfaces.IManager, log *logger.Logger) *REST {
	m := metrics.NewAPIMetrics()
	return &REST{
		log:     log,
		config:  config,
		manager: manager,
		router:  router.InitGinRouter(config.ENV),

		middlewares: middlewares.InitMiddlewares(manager, m, log.WithModule("MIDDLEWARE")),
	}
}

func (r *REST) Run() {
	r.InitRoutes()

	r.server = &http.Server{
		Addr:           ":" + r.config.ServerPORT,
		Handler:        r.router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		r.log.Info(fmt.Sprintf("REST server running at: %s", r.config.ServerPORT))
		if err := r.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			r.log.Error(fmt.Sprintf("r.ListenAndServer: %v", err))
		}
	}()
}

func (r *REST) Shutdown(ctx context.Context) error {
	r.log.Info("REST graceful shut down...")
	return r.server.Shutdown(ctx)
}
