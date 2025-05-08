package usecase

import (
	"context"
	"fmt"
	"tablelink/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) error
}

type authUseCase struct {
	userRepository repository.UserRepository
	redis          redis.Client
}

func NewAuthUseCase(userRepo repository.UserRepository, redis redis.Client) AuthUseCase {
	return &authUseCase{
		userRepository: userRepo,
		redis:          redis,
	}
}

func (u *authUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("Make sure you have provide valid email or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("Make sure you have provide valid email or password")
	}

	token := uuid.NewString()
	now := time.Now().UTC()
	user.LastAccess = &now
	if err := u.redis.Set(ctx, token, user, 24*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("Failed to save token to cache")
	}
	return token, nil
}

func (u *authUseCase) Logout(ctx context.Context, token string) error {
	return u.redis.Del(ctx, token).Err()
}
