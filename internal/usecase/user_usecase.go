package usecase

import (
	"context"
	"tablelink/internal/constant/message"
	"tablelink/internal/domain"
	"tablelink/internal/helper"
	"tablelink/internal/model"
	"tablelink/internal/model/converter"
	"tablelink/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserUseCase interface {
	GetAllUser(ctx context.Context) ([]*model.UserResponse, error)
	CreateUser(ctx context.Context, request *model.CreateUserRequest) error
	UpdateUser(ctx context.Context, request *model.UpdateUserRequest) error
	DeleteUser(ctx context.Context, userID int) error
}

type userUseCase struct {
	userRepo  repository.UserRepository
	rightRepo repository.RoleRightRepository
	logger    helper.Logger
}

func NewUserUseCase(userRepo repository.UserRepository, rightRepo repository.RoleRightRepository, logger helper.Logger) UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		rightRepo: rightRepo,
		logger:    logger,
	}
}

func (u *userUseCase) GetAllUser(ctx context.Context) ([]*model.UserResponse, error) {
	users, err := u.userRepo.ListAll(ctx)
	if err != nil {
		u.logger.Errorf("failed to get list users in database : %v", err)
		return nil, status.Error(codes.Internal, message.InternalGracefulError)
	}

	return converter.UsersToResponses(users), nil
}

func (u *userUseCase) CreateUser(ctx context.Context, request *model.CreateUserRequest) error {
	total, err := u.userRepo.CountByEmail(ctx, request.Email)
	if err != nil {
		u.logger.Errorf("failed to count user from database : %v", err)
		return status.Error(codes.Internal, message.InternalGracefulError)
	}

	if total > 0 {
		return status.Error(codes.AlreadyExists, message.ClientUserAlreadyExist)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Errorf("failed to generated hashed user passwod using bcrypt : %v", err)
		return status.Error(codes.Internal, message.InternalGracefulError)
	}

	user := &domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		RoleID:   request.RoleID,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		u.logger.Errorf("failed to create user in database : %v", err)
		return status.Error(codes.Internal, message.InternalGracefulError)
	}

	return nil
}

func (u *userUseCase) UpdateUser(ctx context.Context, request *model.UpdateUserRequest) error {
	user := &domain.User{
		ID:   request.UserID,
		Name: request.Name,
	}

	if err := u.userRepo.UpdateName(ctx, user); err != nil {
		u.logger.Errorf("failed to update user in database : %v", err)
		return status.Error(codes.Internal, message.InternalGracefulError)
	}

	return nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, userID int) error {
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		u.logger.Errorf("failed to delete user in database : %v", err)
		return status.Error(codes.Internal, message.InternalGracefulError)
	}

	return nil
}
