package models

type PoolReq struct {
	UserID  string `json:"user_id"`
	BoardID int    `json:"board_id"`
}
