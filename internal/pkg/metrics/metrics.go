package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type APIMetrics struct {
	SuccessfulHttpRequests prometheus.Counter
	ErrorHttpRequests      prometheus.Counter

	SuccessSessionVerification prometheus.Counter
	ErrorSessionVerification   prometheus.Counter

	// ws
	WSPoolConnectionRequest prometheus.Counter
}

// NewAPIMetrics creates a new instance of APIMetrics with Prometheus counters initialized.
func NewAPIMetrics() *APIMetrics {
	serviceName := "auth_api"

	return &APIMetrics{
		SuccessfulHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests", serviceName),
			Help: "The total number of successful http requests",
		}),
		ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests", serviceName),
			Help: "The total number of unsuccessful http requests",
		}),
		SuccessSessionVerification: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_session_verification_requests", serviceName),
			Help: "The total number of successful session verification http requests",
		}),
		ErrorSessionVerification: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_session_verification_requests", serviceName),
			Help: "The total number of unsuccessful session verification http requests",
		}),
		WSPoolConnectionRequest: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_verify_ws_pool_connection_requests", serviceName),
			Help: "The total number of ws pool connection http requests",
		}),
	}
}
