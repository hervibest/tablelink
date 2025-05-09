package usecase_test

import (
	"context"
	"errors"
	"testing"

	"tablelink/internal/domain"
	mockhelper "tablelink/internal/helper/mock"
	"tablelink/internal/model"
	mockrepo "tablelink/internal/repository/mock"
	"tablelink/internal/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepo.NewMockUserRepository(ctrl)
	mockRoleRepo := mockrepo.NewMockRoleRightRepository(ctrl)
	mockLogger := mockhelper.NewMockLogger(ctrl)

	userUC := usecase.NewUserUseCase(mockUserRepo, mockRoleRepo, mockLogger)
	ctx := context.Background()

	t.Run("GetAllUser success", func(t *testing.T) {
		mockUsers := []*domain.User{
			{ID: 1, Name: "User1", Email: "u1@example.com", RoleID: 2},
		}
		mockUserRepo.EXPECT().ListAll(ctx).Return(mockUsers, nil)

		users, err := userUC.GetAllUser(ctx)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "User1", users[0].Name)
	})

	t.Run("GetAllUser error", func(t *testing.T) {
		mockUserRepo.EXPECT().ListAll(ctx).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		users, err := userUC.GetAllUser(ctx)
		assert.Nil(t, users)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("CreateUser already exists", func(t *testing.T) {
		req := &model.CreateUserRequest{Name: "X", Email: "test@example.com", Password: "123", RoleID: 1}
		mockUserRepo.EXPECT().CountByEmail(ctx, req.Email).Return(1, nil)

		err := userUC.CreateUser(ctx, req)
		assert.Equal(t, codes.AlreadyExists, status.Code(err))
	})

	t.Run("CreateUser success", func(t *testing.T) {
		req := &model.CreateUserRequest{Name: "X", Email: "new@example.com", Password: "123", RoleID: 1}
		mockUserRepo.EXPECT().CountByEmail(ctx, req.Email).Return(0, nil)
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes() // in case bcrypt fails

		err := userUC.CreateUser(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("CreateUser repo error", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()) // ← Tambahkan ini
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()) // ← Tambahkan ini
		req := &model.CreateUserRequest{Name: "Y", Email: "y@example.com", Password: "123", RoleID: 1}
		mockUserRepo.EXPECT().CountByEmail(ctx, req.Email).Return(0, nil)
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("insert fail"))

		err := userUC.CreateUser(ctx, req)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("UpdateUser success", func(t *testing.T) {
		req := &model.UpdateUserRequest{UserID: 2, Name: "Updated"}
		mockUserRepo.EXPECT().UpdateName(ctx, gomock.Any()).Return(nil)

		err := userUC.UpdateUser(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("UpdateUser error", func(t *testing.T) {
		req := &model.UpdateUserRequest{UserID: 3, Name: "X"}
		mockUserRepo.EXPECT().UpdateName(ctx, gomock.Any()).Return(errors.New("update error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		err := userUC.UpdateUser(ctx, req)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("DeleteUser success", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()) // ← Tambahkan ini
		mockUserRepo.EXPECT().Delete(ctx, 5).Return(nil)

		err := userUC.DeleteUser(ctx, 5)
		assert.NoError(t, err)
	})

	t.Run("DeleteUser error", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()) // ← Tambahkan ini
		mockUserRepo.EXPECT().Delete(ctx, 6).Return(errors.New("delete error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		err := userUC.DeleteUser(ctx, 6)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}
