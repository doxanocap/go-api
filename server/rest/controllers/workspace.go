package controllers

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/metrics"
	"auth-api/internal/pkg/tools"
	"github.com/doxanocap/pkg/ctxholder"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Workspace struct {
	log     *zap.Logger
	config  *models.Config
	metrics *metrics.APIMetrics
	service interfaces.IService
}

func InitWorkspaceController(
	config *models.Config,
	service interfaces.IService,
	metrics *metrics.APIMetrics,
	log *zap.Logger) *Workspace {
	return &Workspace{
		log:     log,
		config:  config,
		metrics: metrics,
		service: service,
	}
}

func (ctl *Workspace) Create(c *gin.Context) {
	userID := ctxholder.GetUserID(c)
	if userID == "" {
		errs.SetGinError(c, models.HttpUnauthorized)
		return
	}

	w, err := ctl.service.Workspace().Create(c, userID)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, w)
}

func (ctl *Workspace) GetAllByUserID(c *gin.Context) {
	userID := ctxholder.GetUserID(c)
	if !tools.IsUUID(userID) {
		errs.SetGinError(c, models.HttpBadRequest)
		return
	}

	response, err := ctl.service.Workspace().GetAllByUserID(c, userID)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctl *Workspace) GetByID(c *gin.Context) {
	workspaceID := c.Param("id")
	if !tools.IsUUID(workspaceID) {
		errs.SetGinError(c, models.HttpBadRequest)
		return
	}

	response, err := ctl.service.Workspace().GetByIDCode(c, workspaceID)
	if err != nil {
		errs.SetGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ToDTO())
}
