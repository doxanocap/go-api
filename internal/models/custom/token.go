package custom

import (
	"fmt"
	"strings"
)

type Token string

func NewToken(token string) Token {
	return Token(strings.TrimSpace(token))
}

func (t *Token) String() string {
	return string(*t)
}

func (t *Token) StringPtr() *string {
	if t == nil {
		return nil
	}

	res := t.String()

	return &res
}

func (t *Token) Mask() string {
	var start int = 0

	if t == nil {
		return ""
	}

	if len(*t) <= 12 {
		return t.String()
	}

	if strings.Contains(t.String(), "Bearer") {
		start = 7
	}

	return fmt.Sprintf("%s...%s", t.String()[start:start+4], t.String()[len(t.String())-5:len(t.String())-1])
}
