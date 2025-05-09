package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tablelink/internal/constant/message"
	"tablelink/internal/helper"
	"tablelink/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RightUseCase interface {
	Authorize(ctx context.Context, roleID int, section, route, action string) error
}

type rightUseCase struct {
	rightRepo repository.RoleRightRepository
	logger    helper.Logger
}

func NewRightUseCase(rightRepo repository.RoleRightRepository, logger helper.Logger) RightUseCase {
	return &rightUseCase{
		rightRepo: rightRepo,
		logger:    logger,
	}
}

func (u *rightUseCase) Authorize(ctx context.Context, roleID int, section, route, action string) error {
	rights, err := u.rightRepo.CheckPermission(ctx, roleID, section, route)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.logger.Printf("user does not have any permission in route : %s and role id : %d and action : %s", route, roleID, action)
			return status.Error(codes.PermissionDenied, message.ClientPermissionDenied)
		}
		return status.Error(codes.Internal, fmt.Sprintf("failed to check permission in database : %v", err))
	}

	var allowed int
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
	if allowed != 1 {
		u.logger.Printf("user have a permision but not with the action in route : %s and role id : %d and action : %s", route, roleID, action)
		return status.Error(codes.PermissionDenied, message.ClientPermissionDenied)
	}
	return nil
}
