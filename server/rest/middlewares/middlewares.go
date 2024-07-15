package middlewares

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/pkg/metrics"
	"auth-api/logger"
	"github.com/doxanocap/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"time"
)

type Middlewares struct {
	manager interfaces.IManager
	metrics *metrics.APIMetrics
	log     *logger.Logger
}

func InitMiddlewares(manager interfaces.IManager, metrics *metrics.APIMetrics, log *logger.Logger) *Middlewares {
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

		args := slog.Group("request",
			slog.String("route", c.FullPath()),
			slog.Duration("latency", latency),
			slog.Int("status", status))

		privateErr := errs.GetGinPrivateErr(c)
		if privateErr != nil {
			m.metrics.ErrorHttpRequests.Inc()
			m.log.Error(privateErr.Error(), args)
			if c.Writer.Status() == http.StatusOK {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		// public error will be set automatically if errs.SetGinError is called with it
		m.metrics.SuccessfulHttpRequests.Inc()
		m.log.Info("ok", args)
	}
}

func (m *Middlewares) GinMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}
