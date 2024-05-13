package models

import "time"

type Workspace struct {
	ID           int64      `json:"workspace_id" db:"workspace_id"`
	IDCode       string     `json:"workspace_idcode" db:"workspace_idcode"`
	Body         string     `json:"body" db:"body"`
	LastEditorID string     `json:"last_editor_id,omitempty" db:"last_editor_id"`
	CreatedAt    *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (w Workspace) ToDTO() *WorkspaceDTO {
	return nil
}

type WorkspaceDTO struct {
	IDCode       string     `json:"id" db:"idcode"`
	Body         string     `json:"body" db:"body"`
	LastEditorID string     `json:"last_editor_id,omitempty" db:"last_editor_id"`
	CreatedAt    *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type WorkspaceUser struct {
	WorkspaceIDRef int64  `json:"workspace_idref" db:"workspace_idref"`
	UserIDRef      string `json:"user_idref" db:"user_idref"`
	IsPresent      bool   `json:"is_present" db:"is_present"`
}

type WorkspaceMessage struct {
	MessageID      string     `json:"message_id" db:"message_id"`
	WorkspaceIDRef int64      `json:"workspace_idref" db:"workspace_idref"`
	UserIDRef      string     `json:"user_idref" db:"user_idref"`
	Content        string     `json:"content" db:"content"`
	Timestamp      *time.Time `json:"timestamp" db:"ts"`
}
