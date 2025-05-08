package grpc

import (
	"context"
	"tablelink/internal/usecase"
	"tablelink/proto/proto/userpb"
	"time"
)

type UserHandler struct {
	userUC usecase.UserUseCase
	userpb.UnimplementedUsersServiceServer
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUC: uc}
}

func (h *UserHandler) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	users, err := h.userUC.ListUser(ctx, int(req.GetRoleId()), req.GetSection(), req.GetRoute())
	if err != nil {
		return &userpb.ListUsersResponse{
			Status:  false,
			Message: err.Error(),
		}, nil
	}

	var pbUsers []*userpb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &userpb.User{
			Id:         int32(u.ID),
			Name:       u.Name,
			Email:      u.Email,
			RoleId:     int32(u.RoleID),
			LastAccess: u.LastAccess.Format(time.RFC3339),
		})
	}

	return &userpb.ListUsersResponse{
		Status:  true,
		Message: "Successfully get list users",
		Users:   pbUsers,
	}, nil
}
