package grpc

import (
	"context"
	"tablelink/internal/constant/message"
	"tablelink/internal/model"
	"tablelink/internal/usecase"
	"tablelink/proto/proto/userpb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCHandler struct {
	userUC usecase.UserUseCase
	userpb.UnimplementedUserServiceServer
}

func NewUserHandler(uc usecase.UserUseCase) *UserGRPCHandler {
	return &UserGRPCHandler{userUC: uc}
}

func (h *UserGRPCHandler) GetAllUser(ctx context.Context, pbReq *userpb.GetAllUserRequest) (*userpb.GetAllUserResponse, error) {
	users, err := h.userUC.GetAllUser(ctx)
	if err != nil {
		return nil, err
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

	return &userpb.GetAllUserResponse{
		Status:  true,
		Message: message.Success,
		Users:   pbUsers,
	}, nil
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, pbReq *userpb.CreateUserRequest) (*userpb.CreateUserReponse, error) {
	request := &model.CreateUserRequest{
		RoleID:   int(pbReq.User.GetRoleId()),
		Name:     pbReq.User.GetName(),
		Email:    pbReq.User.GetEmail(),
		Password: pbReq.User.GetPassword(),
	}

	if err := h.userUC.CreateUser(ctx, request); err != nil {
		return nil, err
	}

	return &userpb.CreateUserReponse{
		Status:  true,
		Message: message.Success,
	}, nil
}

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, pbReq *userpb.UpdateUserRequest) (*userpb.UpdateUserReponse, error) {
	authVal := ctx.Value(model.AuthContextKey)
	auth, ok := authVal.(*model.AuthResponse)
	if !ok || auth == nil {
		return nil, status.Error(codes.Internal, message.InternalUserAuthNotFound)
	}

	request := &model.UpdateUserRequest{
		UserID: auth.ID,
		Name:   pbReq.GetName(),
	}

	if err := h.userUC.UpdateUser(ctx, request); err != nil {
		return nil, err
	}

	return &userpb.UpdateUserReponse{
		Status:  true,
		Message: message.Success,
	}, nil
}

func (h *UserGRPCHandler) DeleteUser(ctx context.Context, pbReq *userpb.DeleteUserRequest) (*userpb.DeleteUserReponse, error) {
	if err := h.userUC.DeleteUser(ctx, int(pbReq.GetUserId())); err != nil {
		return nil, err
	}

	return &userpb.DeleteUserReponse{
		Status:  true,
		Message: message.Success,
	}, nil

}
