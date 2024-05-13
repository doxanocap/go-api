package controllers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/metrics"
	"errors"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

type Websocket struct {
	log     *zap.Logger
	config  *models.Config
	metrics *metrics.APIMetrics
	service interfaces.IService

	upgrader *websocket.Upgrader
}

func InitWebsocketController(config *models.Config, service interfaces.IService,
	metrics *metrics.APIMetrics, log *zap.Logger) *Websocket {
	return &Websocket{
		log:     log,
		config:  config,
		metrics: metrics,
		service: service,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (ws *Websocket) Pool(c *gin.Context) {
	ws.metrics.WSPoolConnectionRequest.Inc()
	log := ws.log

	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")

	header := http.Header{}
	log.Info(fmt.Sprintf("connection: user_id: %s | workspace_id: %s | IP: %s", userID, workspaceID, c.RemoteIP()))

	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, header)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	if err = ws.service.Workspace().HandleNewClient(c, conn, userID, workspaceID); err != nil {
		if errors.Is(err, models.ErrConnGracefullyClosed) {
			return
		}

		log.Error(fmt.Sprintf("HandleNewClient: %s", err))
		if err = conn.Close(); err != nil {
			log.Error(fmt.Sprintf("close: %s", err))
		}
		return
	}
}
