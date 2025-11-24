package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fajrinajiseno/mygolangapp/internal/entity"
	"github.com/fajrinajiseno/mygolangapp/internal/openapigen"
	"github.com/fajrinajiseno/mygolangapp/internal/transport"
	"github.com/fajrinajiseno/mygolangapp/internal/usecase"
)

type apiHandler struct {
	paymentUC usecase.PaymentUsecase
	authUC    usecase.AuthUsecase
}

func (a *apiHandler) PostDashboardV1AuthLogin(w http.ResponseWriter, r *http.Request) {
	var req openapigen.PostDashboardV1AuthLoginJSONBody
	if !decodeJSONBody(w, r, &req) {
		return
	}
	token, user, err := a.authUC.Login(req.Email, req.Password)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(openapigen.LoginResponse{Email: &user.Email, Role: &user.Role, Token: &token})
	if err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
		return
	}
}

func (a *apiHandler) GetDashboardV1Payments(w http.ResponseWriter, r *http.Request, body openapigen.GetDashboardV1PaymentsParams) {
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
		status = *body.Sort
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

func (a *apiHandler) PutDashboardV1PaymentIdReview(w http.ResponseWriter, r *http.Request, id string) {
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

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	if r.Body == nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("empty body"))
		return false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("failed to read body"))
		return false
	}

	if err := json.Unmarshal(body, dst); err != nil {
		transport.WriteAppError(w, entity.ErrorBadRequest("invalid json: "+err.Error()))
		return false
	}
	return true
}
