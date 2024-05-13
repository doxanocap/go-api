package consts

import (
	"auth-api/internal/models"
	"time"
)

// App constants
const (
	CacheWorkspaceBodyTTL    = 30 * time.Minute
	CacheWorkspaceBodyPrefix = "wspace"

	DateFormat        = "2006-01-02"
	EmailRegexp       = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	PhoneNumberRegexp = `^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`
)

const (
	PingMessage models.MsgType = iota + 1
	SendMessage
	RegisterMessage
	UnregisterMessage
)
