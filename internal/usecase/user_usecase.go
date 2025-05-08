package usecase

import (
	"context"
	"errors"
	"fmt"
	"tablelink/internal/domain"
	"tablelink/internal/repository"
)

type UserUseCase interface {
	ListUser(ctx context.Context, roleID int, section, route string) ([]*domain.User, error)
	CreateUser(ctx context.Context, roleID int, section, route string, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, roleID int, section, route string, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, roleID int, section, route string, userID int) error
}

type userUseCase struct {
	userRepo  repository.UserRepository
	rightRepo repository.RoleRightRepository
}

func NewUserUseCase(userRepo repository.UserRepository, rightRepo repository.RoleRightRepository) UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		rightRepo: rightRepo,
	}
}

func (u *userUseCase) authorize(ctx context.Context, roleID int, section, route, action string) error {
	rights, err := u.rightRepo.CheckPermission(ctx, roleID, section, route)
	if err != nil {
		return fmt.Errorf("Failed to check permission with err %v", err)
	}

	var allowed bool
	switch action {
	case "create":
		allowed = rights.RCreate
	case "read":
		allowed = rights.RRead
	case "update":
		allowed = rights.RUpdate
	case "delete":
		allowed = rights.RDelete
	}
	if !allowed {
		return errors.New("permission denie")
	}
	return nil

}

func (u *userUseCase) ListUser(ctx context.Context, roleID int, section, route string) ([]*domain.User, error) {
	if err := u.authorize(ctx, roleID, section, route, "read"); err != nil {
		return nil, err
	}

	return u.userRepo.ListAll(ctx)
}

func (u *userUseCase) CreateUser(ctx context.Context, roleID int, section, route string, user *domain.User) (*domain.User, error) {
	if err := u.authorize(ctx, roleID, section, route, "create"); err != nil {
		return nil, err
	}

	return u.userRepo.Create(ctx, user)

}

func (u *userUseCase) UpdateUser(ctx context.Context, roleID int, section, route string, user *domain.User) (*domain.User, error) {
	if err := u.authorize(ctx, roleID, section, route, "updaet"); err != nil {
		return nil, err
	}

	return u.userRepo.Update(ctx, user)
}

func (u *userUseCase) DeleteUser(ctx context.Context, roleID int, section, route string, userID int) error {
	if err := u.authorize(ctx, roleID, section, route, "delete"); err != nil {
		return err
	}

	return u.userRepo.Delete(ctx, userID)
}
