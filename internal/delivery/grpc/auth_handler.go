package grpc

import (
	"context"
	"tablelink/internal/usecase"
	"tablelink/proto/proto/authpb"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
	authpb.UnimplementedAuthServiceServer
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUC: uc,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	token, err := h.authUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &authpb.LoginResponse{
			Status:  false,
			Message: err.Error(),
		}, nil
	}

	return &authpb.LoginResponse{
		Status:  true,
		Message: "Login successful",
		Data: &authpb.LoginData{
			AccessToken: token,
		},
	}, nil

}

func (h *AuthHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	if err := h.authUC.Logout(ctx, req.GetAccessToken()); err != nil {
		return &authpb.LogoutResponse{
			Status:  false,
			Message: err.Error(),
		}, nil
	}

	return &authpb.LogoutResponse{
		Status:  true,
		Message: "Logout successful",
	}, nil

}
