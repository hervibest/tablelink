package grpc

import (
	"context"
	"tablelink/internal/constant/message"
	"tablelink/internal/model"
	"tablelink/internal/usecase"
	"tablelink/proto/proto/authpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGRPCHandler struct {
	authUC usecase.AuthUseCase
	authpb.UnimplementedAuthServiceServer
}

func NewAuthGRPCHandler(uc usecase.AuthUseCase) *AuthGRPCHandler {
	return &AuthGRPCHandler{
		authUC: uc,
	}
}

func (h *AuthGRPCHandler) Login(ctx context.Context, pbReq *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	request := &model.LoginRequest{
		Email:    pbReq.GetEmail(),
		Password: pbReq.GetPassword(),
	}

	token, err := h.authUC.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		Status:  true,
		Message: "Login successful",
		Data: &authpb.LoginData{
			AccessToken: token,
		},
	}, nil

}

func (h *AuthGRPCHandler) Logout(ctx context.Context, pbReq *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	authVal := ctx.Value(model.AuthContextKey)
	auth, ok := authVal.(*model.AuthResponse)
	if !ok || auth == nil {
		return nil, status.Error(codes.Internal, message.InternalUserAuthNotFound)
	}

	if err := h.authUC.Logout(ctx, auth.Token); err != nil {
		return nil, err
	}

	return &authpb.LogoutResponse{
		Status:  true,
		Message: "Logout successful",
	}, nil

}
