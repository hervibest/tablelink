package grpc

import (
	"context"
	"fmt"
	"strings"
	"tablelink/internal/constant/message"
	"tablelink/internal/model"
	"tablelink/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(authUC usecase.AuthUseCase, rightUC usecase.RightUseCase) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// WHITELISTED METHODS
		whitelistedMethods := map[string]bool{
			"/proto.AuthService/Login":    true,
			"/proto.AuthService/Register": true,
		}

		if whitelistedMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no metadata found")
		}

		// 1. Cek Section
		section := md.Get("x-link-service")
		if len(section) == 0 || section[0] != "be" {
			return nil, status.Error(codes.PermissionDenied, "invalid section")
		}

		// 2. Cek Token
		authHeader := md.Get("authorization")
		print(authHeader)
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader[0], "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		user, err := authUC.Verify(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, message.ClientUnauthenticated)
		}

		auth := &model.AuthResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			LastAccess: user.LastAccess,
			Token:      token,
		}
		ctx = context.WithValue(ctx, model.AuthContextKey, auth)

		// 3. Extract Route dan Action
		fullMethod := info.FullMethod // contoh: "/proto.UserService/GetAllUser"
		parts := strings.Split(fullMethod, "/")
		if len(parts) < 3 {
			return nil, status.Error(codes.Internal, "malformed method name")
		}

		// "/proto.UserService/GetAllUser" → "UserService/GetAllUser"
		serviceParts := strings.Split(parts[1], ".")
		service := serviceParts[len(serviceParts)-1]
		method := parts[2]
		route := fmt.Sprintf("%s/%s", service, method)

		var action string
		switch {
		case strings.HasPrefix(strings.ToLower(method), "get"):
			action = "read"
		case strings.HasPrefix(strings.ToLower(method), "create"):
			action = "create"
		case strings.HasPrefix(strings.ToLower(method), "update"):
			action = "update"
		case strings.HasPrefix(strings.ToLower(method), "delete"):
			action = "delete"
		case strings.HasPrefix(strings.ToLower(method), "logout"):
			return handler(ctx, req)
		default:
			action = "read" // default aman
		}

		if err := rightUC.Authorize(ctx, user.RoleID, section[0], route, action); err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		return handler(ctx, req)
	}
}
