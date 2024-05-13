package models

type MsgType uint8

type Message struct {
	UserID      string  `json:"user_id"`
	WorkspaceID string  `json:"workspace_id"`
	MsgType     MsgType `json:"-"`

	ErrorMessage string `json:"error_message"`
	Body         string `json:"body"`

	// TODO переделать на отправку дэльты
}
