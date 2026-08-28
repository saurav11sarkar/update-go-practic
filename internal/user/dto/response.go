package dto

import "time"

type UserResponse struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type UserLoginResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreateAt     time.Time `json:"create_at"`
	UpdateAt     time.Time `json:"update_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"-"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}
