package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fajrinajiseno/mygolangapp/internal/config"
	"github.com/fajrinajiseno/mygolangapp/internal/entity"
	am "github.com/fajrinajiseno/mygolangapp/internal/module/auth/repository/mock"
	pm "github.com/fajrinajiseno/mygolangapp/internal/module/payment/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPayment_ListPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*entity.Payment{
		{
			ID:        "1",
			Merchant:  "merchant 1",
			Status:    "completed",
			Amount:    100,
			CreatedAt: time.Now(),
		},
	}

	mockPaymentRepo := pm.NewMockPaymentRepository(ctrl)
	mockUserRepo := am.NewMockUserRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo.EXPECT().
			GetPayments("completed", "1", "created_at", 10, 1).
			Return(expected, &entity.PaymentSummary{
				TotalByFiler:   1,
				Total:          4,
				TotalCompleted: 2,
				TotalFailed:    1,
				TotalPending:   1,
			}, nil)

		u := NewPaymentUsecase(mockPaymentRepo, mockUserRepo)

		items, totalSummary, err := u.ListPayment("completed", "1", "created_at", 10, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, items)
		assert.Equal(t, 1, totalSummary.TotalByFiler)
		assert.Equal(t, 4, totalSummary.Total)
		assert.Equal(t, 2, totalSummary.TotalCompleted)
		assert.Equal(t, 1, totalSummary.TotalFailed)
		assert.Equal(t, 1, totalSummary.TotalPending)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockPaymentRepo.EXPECT().
			GetPayments("completed", "1", "created_at", 10, 1).
			Return(nil, nil, errors.New("db fail"))

		u := NewPaymentUsecase(mockPaymentRepo, mockUserRepo)

		_, _, err := u.ListPayment("completed", "1", "created_at", 10, 1)
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
