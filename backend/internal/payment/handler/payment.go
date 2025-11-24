package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fajrinajiseno/mygolangapp/internal/entity"
	"github.com/fajrinajiseno/mygolangapp/internal/openapigen"
	"github.com/fajrinajiseno/mygolangapp/internal/payment/usecase"
	"github.com/fajrinajiseno/mygolangapp/internal/transport"
)

type PaymentHandler struct {
	paymentUC usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUC usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUC: paymentUC,
	}
}

func (a *PaymentHandler) GetDashboardV1Payments(w http.ResponseWriter, r *http.Request, body openapigen.GetDashboardV1PaymentsParams) {
	limit := 10
	offset := 0
	sort := "-created_at"
	status := ""

	if body.Limit != nil {
		limit = *body.Limit
	}

	if body.Offset != nil {
		offset = *body.Offset
	}

	if body.Sort != nil {
		sort = *body.Sort
	}

	if body.Status != nil {
		status = *body.Status
	}

	payments, total, totalSuccess, totalFailed, err := a.paymentUC.ListPayment(status, sort, limit, offset)
	if err != nil {
		transport.WriteError(w, err)
		return
	}
	genPayments := make([]openapigen.Payment, len(payments))
	for i, item := range payments {
		amountStr := fmt.Sprint(item.Amount)
		genPayments[i] = openapigen.Payment{
			Id:        &item.ID,
			Amount:    &amountStr,
			CreatedAt: &item.CreatedAt,
			Merchant:  &item.Merchant,
			Status:    &item.Status,
		}
	}
	err = json.NewEncoder(w).Encode(openapigen.PaymentListResponse{Meta: &openapigen.PaginationMeta{
		Limit:  body.Limit,
		Offset: body.Offset,
		Total:  &total,
	}, Summary: &openapigen.PaymentSummary{
		Failed:  &totalFailed,
		Success: &totalSuccess,
	}, Payments: &genPayments})
	if err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
		return
	}
}

func (a *PaymentHandler) PutDashboardV1PaymentIdReview(w http.ResponseWriter, r *http.Request, id string) {
	message, err := a.paymentUC.ReviewPayment(r.Context(), id)
	if err != nil {
		transport.WriteError(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(openapigen.PaymentReviewResponse{Message: &message})
	if err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
		return
	}
}
