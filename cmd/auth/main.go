package main

import (
	"context"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"tablelink/internal/adapter"
	"tablelink/internal/config"
	grpcHandler "tablelink/internal/delivery/grpc/handler"
	grpcInterceptor "tablelink/internal/delivery/grpc/interceptors"

	"tablelink/internal/repository"
	"tablelink/internal/usecase"
	"tablelink/proto/proto/authpb"
)

func main() {
	ctx := context.Background()
	logger := logrus.New()
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Fatal(err)
	}

	db := config.NewDB(ctx, logger, cfg)
	defer db.Close()

	redis := config.NewRedis(cfg)
	defer redis.Close()

	cacheAdapter := adapter.NewCacheAdapter(redis)

	userRepo := repository.NewUserRepository(db)
	rightsRepo := repository.NewRoleRightRepository(db)

	rightUC := usecase.NewRightUseCase(rightsRepo, logger)
	authUC := usecase.NewAuthUseCase(userRepo, cacheAdapter, logger)

	lis, err := net.Listen("tcp", ":"+cfg.PortAuth)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpcInterceptor.AuthInterceptor(authUC, rightUC)),
	)

	authpb.RegisterAuthServiceServer(server, grpcHandler.NewAuthGRPCHandler(authUC))

	log.Printf("Auth service listening on %s", cfg.PortAuth)
	if err := server.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
