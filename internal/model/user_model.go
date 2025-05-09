package model

import "time"

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

type UpdateUserRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type UserResponse struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	RoleID     int        `json:"role_id"`
	LastAccess *time.Time `json:"last_access"`
}
