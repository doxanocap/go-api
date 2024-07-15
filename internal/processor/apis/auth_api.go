package apis

import (
	"auth-api/config"
	"auth-api/logger"
	"crypto/tls"
	"net/http"
	"time"
)

type AuthAPI struct {
	log    *logger.Logger
	config *config.Config
	client *http.Client
}

func NewAuthAPIProcessor(config *config.Config, log *logger.Logger) *AuthAPI {
	return &AuthAPI{
		log:    log,
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}
