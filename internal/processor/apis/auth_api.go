package apis

import (
	"auth-api/internal/models"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/doxanocap/pkg/gohttp"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type AuthAPI struct {
	log    *zap.Logger
	config *models.Config
	client *http.Client
}

func NewAuthAPIProcessor(config *models.Config, log *zap.Logger) *AuthAPI {
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

func (a *AuthAPI) GetUserByID(ctx context.Context, userID string) (user *models.UserDTO, err error) {
	user = &models.UserDTO{}
	route := fmt.Sprintf(a.newRoute("/v1/auth/user/%s"), userID)

	response, err := gohttp.NewRequest(a.client).
		SetURL(route).
		SetMethod(http.MethodGet).
		Execute(ctx)
	if err != nil {
		return nil, errs.Wrap("auth_api.GetUserByID", err)
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, models.ErrUserNotFound
	}

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		httpErr := &errs.HttpError{}
		_ = json.Unmarshal(raw, httpErr)
		return nil, httpErr
	}

	if err = json.Unmarshal(raw, user); err != nil {
		return nil, errs.Wrap("auth_api.GetUserByID", err)
	}

	return user, nil
}

func (a *AuthAPI) VerifySession(ctx context.Context, token string) (us *models.UserSession, err error) {
	us = &models.UserSession{}

	response, err := gohttp.NewRequest(a.client).
		SetURL(a.newRoute("/v1/auth/verify")).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetMethod(http.MethodGet).
		Execute(ctx)
	if err != nil {
		return nil, errs.Wrap("auth_api.VerifySession", err)
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, models.HttpUnauthorized
	}

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		httpErr := &errs.HttpError{}
		_ = json.Unmarshal(raw, httpErr)
		return nil, httpErr
	}

	if err = json.Unmarshal(raw, us); err != nil {
		return nil, errs.Wrap("auth_api.VerifySession", err)
	}

	return us, nil
}

//func (a *AuthAPI) GetUserByID(ctx context.Context, userID string) (user *models.UserDTO, err error) {
//	user = &models.UserDTO{}
//	route := fmt.Sprintf(a.newRoute("/v1/auth/user/%s"), userID)
//
//	response, err := gohttp.NewRequest(a.client).
//		SetURL(route).
//		SetMethod(http.MethodGet).
//		Execute(ctx)
//	if err != nil {
//		return nil, errs.Wrap("auth_api.GetUserByID", err)
//	}
//	if response.StatusCode == http.StatusNotFound {
//		return nil, models.ErrUserNotFound
//	}
//
//	raw, err := io.ReadAll(response.Body)
//	if err != nil {
//		return
//	}
//
//	if response.StatusCode != http.StatusOK {
//		httpErr := &errs.HttpError{}
//		_ = json.Unmarshal(raw, httpErr)
//		return nil, httpErr
//	}
//
//	if err = json.Unmarshal(raw, user); err != nil {
//		return nil, errs.Wrap("auth_api.GetUserByID", err)
//	}
//
//	return user, nil
//}

func (a *AuthAPI) newRoute(r string) string {
	return fmt.Sprintf("%s%s", a.config.APIs.AuthAPI, r)
}
