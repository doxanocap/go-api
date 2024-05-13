package models

import "time"

type User struct {
	ID          int64      `db:"user_id"`
	IDCode      string     `db:"user_idcode"`
	Email       string     `db:"email"`
	PhoneNumber string     `db:"phone_number"`
	Activated   bool       `db:"activated"`
	Password    string     `db:"password"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (u *User) ToUserDTO() UserDTO {
	return UserDTO{
		IDCode:      u.IDCode,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Activated:   u.Activated,
		Password:    u.Password,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
	}
}

type UserDTO struct {
	IDCode      string     `json:"id"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Activated   bool       `json:"activated"`
	Password    string     `json:"-"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (u *UserDTO) ToUser() *User {
	return &User{
		IDCode:      u.IDCode,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Activated:   u.Activated,
		Password:    u.Password,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
	}
}

type UserSession struct {
	UserIDCode    string `json:"user_id" db:"user_idcode"`
	UserEmail     string `json:"email" db:"email"`
	UserActivated bool   `json:"activated" db:"activated"`
	StartedAt     int64  `json:"session_started_at" db:"started_at"`
}

func (u *UserSession) ToUser() *User {
	return &User{
		IDCode:    u.UserIDCode,
		Email:     u.UserEmail,
		Activated: u.UserActivated,
	}
}
