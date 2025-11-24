package usecase

import (
	"context"

	authRepository "github.com/fajrinajiseno/mygolangapp/internal/auth/repository"
	"github.com/fajrinajiseno/mygolangapp/internal/entity"
	"github.com/fajrinajiseno/mygolangapp/internal/middleware"
	paymentRepository "github.com/fajrinajiseno/mygolangapp/internal/payment/repository"
)

//go:generate mockgen -source payment.go -destination mock/payment_mock.go -package=mock
type PaymentUsecase interface {
	ListPayment(status string, sortExpr string, limit int, offset int) ([]*entity.Payment, int, int, int, error)
	ReviewPayment(ctx context.Context, id string) (string, error)
}

type Payment struct {
	userRepo    authRepository.UserRepository
	paymentRepo paymentRepository.PaymentRepository
}

func NewPaymentUsecase(pr paymentRepository.PaymentRepository, ur authRepository.UserRepository) *Payment {
	return &Payment{paymentRepo: pr, userRepo: ur}
}

func (u *Payment) ListPayment(status string, sortExpr string, limit int, offset int) ([]*entity.Payment, int, int, int, error) {
	return u.paymentRepo.GetPayments(status, sortExpr, limit, offset)
}

func (u *Payment) ReviewPayment(ctx context.Context, id string) (string, error) {
	userId := middleware.GetUserID(ctx)
	if userId == "" {
		return "", entity.ErrorNotFound("user not found")
	}
	user, err := u.userRepo.GetUserById(userId)
	if err != nil {
		return "", entity.ErrorNotFound("user not found")
	}
	const OperationRole = "operation"
	if user.Role != OperationRole {
		return "", entity.ErrorForbidden("user forbidden")
	}
	return u.paymentRepo.Review(id)
}
