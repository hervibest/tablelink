package usecase

import "context"

type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) error
}
