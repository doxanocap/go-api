package middlewares

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/pkg/metrics"
	"github.com/doxanocap/pkg/ctxholder"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Middlewares struct {
	manager interfaces.IManager
	metrics *metrics.APIMetrics
	log     *zap.Logger
}

const (
	keyRefreshToken = "refresh_token"
)

func InitMiddlewares(manager interfaces.IManager, metrics *metrics.APIMetrics, log *zap.Logger) *Middlewares {
	return &Middlewares{
		log:     log,
		metrics: metrics,
		manager: manager,
	}
}

func (m *Middlewares) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)

		route := c.FullPath()
		status := c.Writer.Status()
		if route == "" && status == http.StatusNotFound {
			return
		}

		fields := []zap.Field{
			zap.String("route", c.FullPath()),
			zap.Duration("latency", latency),
			zap.Int("status", c.Writer.Status()),
		}

		log := m.log.With(zap.Any("payload", fields))

		privateErr := errs.GetGinPrivateErr(c)
		if privateErr != nil {
			m.metrics.ErrorHttpRequests.Inc()
			log.Error(privateErr.Error())
			if c.Writer.Status() == http.StatusOK {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		// public error will be set automatically if errs.SetGinError is called with it
		m.metrics.SuccessfulHttpRequests.Inc()
		log.Info("ok")
	}
}

func (m *Middlewares) VerifySession() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := m.log.Named("[SESSION]")

		token := m.getAuthToken(c)
		if token == "" {
			log.Error("getAuthToken")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.HttpUnauthorized)
			return
		}

		uc, err := m.manager.Processor().APIs().AuthAPI().VerifySession(c, token)
		if err != nil {
			log.Error(err.Error())
			httpErr := errs.UnmarshalError(err)
			if httpErr.StatusCode != 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, httpErr)
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, models.HttpUnauthorized)
			}
			return
		}

		ctxholder.SetUserID(c, uc.UserIDCode)
		c.Next()
	}
}

func (m *Middlewares) GinMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}

func (m *Middlewares) getAuthToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}
