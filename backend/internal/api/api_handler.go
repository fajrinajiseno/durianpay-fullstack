package api

import (
	"net/http"

	ah "github.com/fajrinajiseno/mygolangapp/internal/auth/handler"
	"github.com/fajrinajiseno/mygolangapp/internal/openapigen"
	ph "github.com/fajrinajiseno/mygolangapp/internal/payment/handler"
)

type APIHandler struct {
	Auth    *ah.AuthHandler
	Payment *ph.PaymentHandler
}

var _ openapigen.ServerInterface = (*APIHandler)(nil)

func (h *APIHandler) PostDashboardV1AuthLogin(w http.ResponseWriter, r *http.Request) {
	h.Auth.PostDashboardV1AuthLogin(w, r)
}

func (h *APIHandler) GetDashboardV1Payments(w http.ResponseWriter, r *http.Request, body openapigen.GetDashboardV1PaymentsParams) {
	h.Payment.GetDashboardV1Payments(w, r, body)
}

func (h *APIHandler) PutDashboardV1PaymentIdReview(w http.ResponseWriter, r *http.Request, id string) {
	h.Payment.PutDashboardV1PaymentIdReview(w, r, id)
}
