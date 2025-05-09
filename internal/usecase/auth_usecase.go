package usecase

import (
	"context"
	"database/sql"
	"errors"
	"tablelink/internal/adapter"
	"tablelink/internal/constant/message"
	"tablelink/internal/domain"
	"tablelink/internal/helper"
	"tablelink/internal/model"
	"tablelink/internal/model/converter"
	"tablelink/internal/repository"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUseCase interface {
	Login(ctx context.Context, request *model.LoginRequest) (string, error)
	Logout(ctx context.Context, token string) error
	Verify(ctx context.Context, token string) (*model.UserResponse, error)
}

type authUseCase struct {
	userRepository repository.UserRepository
	cacheAdapter   adapter.CacheAdapter
	logger         helper.Logger
}

func NewAuthUseCase(userRepo repository.UserRepository, cacheAdapter adapter.CacheAdapter, logger helper.Logger) AuthUseCase {
	return &authUseCase{
		userRepository: userRepo,
		cacheAdapter:   cacheAdapter,
		logger:         logger,
	}
}

func (u *authUseCase) Login(ctx context.Context, request *model.LoginRequest) (string, error) {
	user, err := u.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.logger.Printf("User trying to login using invalid email : %s", request.Email)
			return "", status.Error(codes.InvalidArgument, message.ClientInvalidEmailOrPassword)
		}
		u.logger.Errorf("failed to get user by email in database : %v", err)
		return "", status.Error(codes.Internal, message.InternalGracefulError)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.logger.Printf("User trying to login using invalid password with email : %s", request.Email)
		return "", status.Error(codes.InvalidArgument, message.ClientInvalidEmailOrPassword)
	}

	token := uuid.NewString()
	now := time.Now().UTC()
	user.LastAccess = &now

	if err := u.userRepository.UpdateLastAccess(ctx, user); err != nil {
		u.logger.Errorf("failed to update last access in database : %v", err)
		return "", status.Error(codes.Internal, message.InternalGracefulError)
	}

	userByte, err := sonic.ConfigFastest.Marshal(user)
	if err != nil {
		u.logger.Errorf("failed to marshal user struct using sonic : %v", err)
		return "", status.Error(codes.Internal, message.InternalGracefulError)
	}

	if err := u.cacheAdapter.Set(ctx, "auth:token:"+token, userByte, 24*time.Hour); err != nil {
		u.logger.Errorf("failed to save token to cache : %v", err)
		return "", status.Error(codes.Internal, message.InternalGracefulError)
	}
	return token, nil
}

func (u *authUseCase) Logout(ctx context.Context, token string) error {
	if err := u.cacheAdapter.Del(ctx, "auth:token:"+token); err != nil {
		u.logger.Errorf("failed to delete token with value : %s from redis : %v", token, err)
		return status.Error(codes.Internal, message.InternalGracefulError)
	}
	return nil
}

func (u *authUseCase) Verify(ctx context.Context, token string) (*model.UserResponse, error) {
	u.logger.Printf("token : %s", token)
	user := new(domain.User)
	userByte, err := u.cacheAdapter.Get(ctx, "auth:token:"+token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			u.logger.Printf("User trying to access route using invalid token : %s", token)
			return nil, status.Error(codes.Unauthenticated, message.ClientUnauthenticated)
		}
		u.logger.Errorf("failed to get token with value : %s from redis : %v", token, err)
		return nil, status.Error(codes.Internal, message.InternalGracefulError)
	}

	if err := sonic.ConfigFastest.Unmarshal([]byte(userByte), user); err != nil {
		u.logger.Errorf("failed to unmarshal user byte using sonic : %v", err)
		return nil, status.Error(codes.Internal, message.InternalGracefulError)
	}

	return converter.UserToResponse(user), nil
}
