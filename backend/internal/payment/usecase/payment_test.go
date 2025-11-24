package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	am "github.com/fajrinajiseno/mygolangapp/internal/auth/repository/mock"
	"github.com/fajrinajiseno/mygolangapp/internal/config"
	"github.com/fajrinajiseno/mygolangapp/internal/entity"
	pm "github.com/fajrinajiseno/mygolangapp/internal/payment/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPayment_ListPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*entity.Payment{
		{ID: "1",
			Merchant:  "merchant 1",
			Status:    "completed",
			Amount:    100,
			CreatedAt: time.Now()},
	}

	mockPaymentRepo := pm.NewMockPaymentRepository(ctrl)
	mockUserRepo := am.NewMockUserRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo.EXPECT().
			GetPayments("completed", "created_at", 10, 1).
			Return(expected, 1, 2, 3, nil)

		u := NewPaymentUsecase(mockPaymentRepo, mockUserRepo)

		items, total, totalSuccess, totalFailed, err := u.ListPayment("completed", "created_at", 10, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, items)
		assert.Equal(t, 1, total)
		assert.Equal(t, 2, totalSuccess)
		assert.Equal(t, 3, totalFailed)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockPaymentRepo.EXPECT().
			GetPayments("completed", "created_at", 10, 1).
			Return(nil, 0, 0, 0, errors.New("db fail"))

		u := NewPaymentUsecase(mockPaymentRepo, mockUserRepo)

		_, _, _, _, err := u.ListPayment("completed", "created_at", 10, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db fail")
	})
}

func TestPayment_ReviewPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentRepo := pm.NewMockPaymentRepository(ctrl)
	mockUserRepo := am.NewMockUserRepository(ctrl)

	t.Run("GetUserById middleware return empty", func(t *testing.T) {
		u := NewPaymentUsecase(mockPaymentRepo, mockUserRepo)

		message, err := u.ReviewPayment(context.Background(), "1")
		assert.Equal(t, "", message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
	t.Run("GetUserById Repo Error", func(t *testing.T) {
		mockUserRepo.EXPECT().
			GetUserById("1").
			Return(nil, errors.New("user not found"))

		u := NewPaymentUsecase(mockPaymentRepo, mockUserRepo)

		ctx := context.WithValue(context.Background(), config.ContextUserID, "1")
		message, err := u.ReviewPayment(ctx, "1")
		assert.Equal(t, "", message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("GetUserById not operation role", func(t *testing.T) {
		mockUserRepo.EXPECT().
			GetUserById("1").
			Return(&entity.User{
				ID:           "u1",
				Email:        "alice@example.com",
				PasswordHash: "123456",
				Role:         "cs",
			}, nil)

		u := NewPaymentUsecase(mockPaymentRepo, mockUserRepo)

		ctx := context.WithValue(context.Background(), config.ContextUserID, "1")
		message, err := u.ReviewPayment(ctx, "1")
		assert.Equal(t, "", message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user forbidden")
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.EXPECT().
			GetUserById("1").
			Return(&entity.User{
				ID:           "u1",
				Email:        "alice@example.com",
				PasswordHash: "123456",
				Role:         "operation",
			}, nil)
		mockPaymentRepo.EXPECT().
			Review("123").
			Return("Success Review", nil)

		u := NewPaymentUsecase(mockPaymentRepo, mockUserRepo)

		ctx := context.WithValue(context.Background(), config.ContextUserID, "1")
		message, err := u.ReviewPayment(ctx, "123")
		assert.NoError(t, err)
		assert.Equal(t, "Success Review", message)
	})
}
