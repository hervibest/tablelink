package main

import (
	"context"
	"log"
	"net"
	"tablelink/internal/cache"
	"tablelink/internal/config"
	"tablelink/internal/repository"
	"tablelink/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.PgURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	rdb := cache.NewRedis(cfg.RedisAddr)
	defer rdb.Close()

	userRepo := repository.NewUserRepository(pool)
	rightRepo := repository.NewRoleRightRepository(pool)
	userUC := usecase.NewUserUseCase(userRepo, rightRepo)

	lis, err := net.Listen("tcp", ":"+cfg.PortUsers)

}
