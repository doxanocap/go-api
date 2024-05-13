package models

import (
	"net/http"

	"github.com/doxanocap/pkg/errs"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// default http error responses
var (
	HttpBadRequest          = errs.NewHttp(http.StatusBadRequest, "bad request")
	HttpNotFound            = errs.NewHttp(http.StatusNotFound, "not found")
	HttpInternalServerError = errs.NewHttp(http.StatusInternalServerError, "internal server error")
	HttpConflict            = errs.NewHttp(http.StatusConflict, "conflict")
	HttpUnauthorized        = errs.NewHttp(http.StatusUnauthorized, "unauthorized")
)

// custom errors for special cases
var (
	ErrInvalidRepoQuery = errs.New("invalid repository query")

	ErrIncorrectPassword = errs.NewHttp(http.StatusUnauthorized, "incorrect password")

	ErrWorkspaceNotFound    = errs.NewHttp(http.StatusNotFound, "workspace wih such id not found")
	ErrMsgFromDifferentUser = errs.NewHttp(http.StatusConflict, "receiving message from different user_id")
	ErrClientAlreadyExist   = errs.NewHttp(http.StatusConflict, "client already exist")
	ErrConnGracefullyClosed = errs.NewHttp(http.StatusConflict, "connection gracefully closed")

	ErrNotPermissiblePassword = errs.NewHttp(http.StatusBadRequest, "not permissible password: its too long/short")

	ErrUserAlreadyExist = errs.NewHttp(http.StatusConflict, "user already exist")
	ErrUserNotFound     = errs.NewHttp(http.StatusConflict, "user not found")

	ErrInvalidOAuthProvider = errs.New("invalid oauth provider")
)
