package main

import (
	"context"

	"github.com/sirupsen/logrus"

	"tablelink/internal/config"
	"tablelink/internal/model"

	"tablelink/internal/repository"
	"tablelink/internal/usecase"
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

	userRepo := repository.NewUserRepository(db)
	rightsRepo := repository.NewRoleRightRepository(db)

	userUC := usecase.NewUserUseCase(userRepo, rightsRepo, logger)

	users := []*model.CreateUserRequest{
		{Name: "Hervi", Email: "hervipro@gmail.com", Password: "12345", RoleID: 1},
		{Name: "Nur", Email: "hervi1234@gmail.com", Password: "12345", RoleID: 2},
		{Name: "Rahmandien", Email: "hervi12345@gmail.com", Password: "12345", RoleID: 3},
	}

	for _, user := range users {
		if err := userUC.CreateUser(ctx, user); err != nil {
			logger.Warn(err)
		}
	}

	logger.Println("Successfully seed user datas")
}
