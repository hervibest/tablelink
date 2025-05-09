package model

import "time"

type LoginRequest struct {
	Email    string
	Password string
}

type AuthResponse struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	LastAccess *time.Time `json:"last_access"`
	Token      string     `json:"token"`
}

type contextKey string

const AuthContextKey = contextKey("auth")
