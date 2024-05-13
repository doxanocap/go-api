package tools

import (
	"auth-api/internal/models/consts"
	"github.com/google/uuid"
	"regexp"
	"time"
)

var (
	phoneNumberRegexpFn = regexp.MustCompile(consts.PhoneNumberRegexp)
	emailRegexpFn       = regexp.MustCompile(consts.EmailRegexp)
)

func CurrTimePtr() *time.Time {
	t := time.Now()
	return &t
}

func GetPtr[T any](v T) *T {
	return &v
}

func IsUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}
