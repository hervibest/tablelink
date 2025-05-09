package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	mockadapter "tablelink/internal/adapter/mock"
	"tablelink/internal/domain"
	mockhelper "tablelink/internal/helper/mock"
	"tablelink/internal/model"
	mockrepo "tablelink/internal/repository/mock"
	"tablelink/internal/usecase"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepo.NewMockUserRepository(ctrl)
	mockCache := mockadapter.NewMockCacheAdapter(ctrl)
	mockLogger := mockhelper.NewMockLogger(ctrl)

	authUC := usecase.NewAuthUseCase(mockUserRepo, mockCache, mockLogger)

	ctx := context.Background()
	password := "secret"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		ID:       1,
		Name:     "Hervi",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		RoleID:   2,
	}

	t.Run("invalid email", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(ctx, user.Email).Return(nil, sql.ErrNoRows)
		mockLogger.EXPECT().Printf(gomock.Any(), user.Email)

		token, err := authUC.Login(ctx, &model.LoginRequest{
			Email:    user.Email,
			Password: password,
		})
		assert.Empty(t, token)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("wrong password", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(ctx, user.Email).Return(user, nil)
		mockLogger.EXPECT().Printf(gomock.Any(), user.Email)

		token, err := authUC.Login(ctx, &model.LoginRequest{
			Email:    user.Email,
			Password: "wrong-password",
		})
		assert.Empty(t, token)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("update last access failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(ctx, user.Email).Return(user, nil)
		mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes() // from bcrypt comparison

		mockUserRepo.EXPECT().UpdateLastAccess(ctx, gomock.Any()).Return(errors.New("db-error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		token, err := authUC.Login(ctx, &model.LoginRequest{
			Email:    user.Email,
			Password: password,
		})
		assert.Empty(t, token)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	// t.Run("marshal failed", func(t *testing.T) {
	// 	mockUserRepo.EXPECT().GetByEmail(ctx, user.Email).Return(user, nil).AnyTimes()
	// 	mockUserRepo.EXPECT().UpdateLastAccess(ctx, gomock.Any()).Return(nil)
	// 	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes()

	// 	// inject invalid user to cause marshal failure
	// 	invalidUser := make(chan int) // cannot be marshaled
	// 	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

	// 	uc := authUC.(usecase.AuthUseCase)
	// 	token := uuid.NewString()
	// 	user.LastAccess = nil

	// 	// force marshal error
	// 	oldMarshal := sonic.ConfigFastest.Marshal
	// 	sonic.ConfigFastest.Marshal = func(val interface{}) ([]byte, error) {
	// 		return nil, errors.New("marshal error")
	// 	}
	// 	defer func() { sonic.ConfigFastest.Marshal = oldMarshal }()

	// 	result, err := uc.Login(ctx, &model.LoginRequest{
	// 		Email:    user.Email,
	// 		Password: password,
	// 	})
	// 	assert.Empty(t, result)
	// 	assert.Equal(t, codes.Internal, status.Code(err))
	// })

	t.Run("cache set failed", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByEmail(ctx, user.Email).Return(user, nil).AnyTimes()
		mockUserRepo.EXPECT().UpdateLastAccess(ctx, gomock.Any()).Return(nil).AnyTimes()
		mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes()

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		mockCache.EXPECT().Set(ctx, gomock.Any(), gomock.Any(), 24*time.Hour).Return(errors.New("redis down"))

		result, err := authUC.Login(ctx, &model.LoginRequest{
			Email:    user.Email,
			Password: password,
		})
		assert.Empty(t, result)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	// t.Run("success login", func(t *testing.T) {
	// 	mockUserRepo.EXPECT().GetByEmail(ctx, user.Email).Return(user, nil)
	// 	mockUserRepo.EXPECT().UpdateLastAccess(ctx, gomock.Any()).Return(nil)
	// 	mockCache.EXPECT().Set(ctx, gomock.Any(), gomock.Any(), 24*time.Hour).Return(nil)
	// 	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes()

	// 	token, err := authUC.Login(ctx, &model.LoginRequest{
	// 		Email:    user.Email,
	// 		Password: password,
	// 	})
	// 	assert.NotEmpty(t, token)
	// 	assert.NoError(t, err)
	// })
}

func TestVerify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepo.NewMockUserRepository(ctrl)
	mockCache := mockadapter.NewMockCacheAdapter(ctrl)
	mockLogger := mockhelper.NewMockLogger(ctrl)

	authUC := usecase.NewAuthUseCase(mockUserRepo, mockCache, mockLogger)
	ctx := context.Background()

	mockLogger.EXPECT().Printf(gomock.Any(), gomock.Any()).AnyTimes()
	token := "valid-token"
	user := &domain.User{
		ID:     1,
		Name:   "Hervi",
		Email:  "test@example.com",
		RoleID: 2,
	}
	userJSON, _ := sonic.ConfigFastest.Marshal(user)

	t.Run("token not found in cache", func(t *testing.T) {
		mockCache.EXPECT().Get(ctx, "auth:token:"+token).Return("", redis.Nil)

		result, err := authUC.Verify(ctx, token)
		assert.Nil(t, result)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("redis error", func(t *testing.T) {
		mockCache.EXPECT().Get(ctx, "auth:token:"+token).Return("", errors.New("redis down"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		result, err := authUC.Verify(ctx, token)
		assert.Nil(t, result)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("unmarshal failed", func(t *testing.T) {
		mockCache.EXPECT().Get(ctx, "auth:token:"+token).Return("not-json", nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		result, err := authUC.Verify(ctx, token)
		assert.Nil(t, result)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("success", func(t *testing.T) {
		mockCache.EXPECT().Get(ctx, "auth:token:"+token).Return(string(userJSON), nil)

		result, err := authUC.Verify(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, result.Email)
	})
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepo.NewMockUserRepository(ctrl)
	mockCache := mockadapter.NewMockCacheAdapter(ctrl)
	mockLogger := mockhelper.NewMockLogger(ctrl)

	authUC := usecase.NewAuthUseCase(mockUserRepo, mockCache, mockLogger)
	ctx := context.Background()
	token := "valid-token"

	t.Run("success", func(t *testing.T) {
		mockCache.EXPECT().Del(ctx, "auth:token:"+token).Return(nil)

		err := authUC.Logout(ctx, token)
		assert.NoError(t, err)
	})

	t.Run("redis delete error", func(t *testing.T) {
		mockCache.EXPECT().Del(ctx, "auth:token:"+token).Return(errors.New("redis down"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		err := authUC.Logout(ctx, token)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}
