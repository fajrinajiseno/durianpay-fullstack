package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fajrinajiseno/mygolangapp/internal/config"
	"github.com/fajrinajiseno/mygolangapp/internal/entity"
	srv "github.com/fajrinajiseno/mygolangapp/internal/service/http"
	"github.com/fajrinajiseno/mygolangapp/internal/usecase/mock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestProtectedEndpointWithoutToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock.NewMockAuthUsecase(ctrl)
	mockPaymentUC := mock.NewMockPaymentUsecase(ctrl)
	srv := srv.NewServer(mockPaymentUC, mockAuthUC)
	ts := httptest.NewServer(srv.Routes())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/dashboard/v1/payments")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestLoginAndAccessProtected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const hour24 = 24
	claims := jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().Add(hour24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(config.JwtSecret)

	mockAuthUC := mock.NewMockAuthUsecase(ctrl)
	mockAuthUC.EXPECT().
		Login("alice@example.com", "password").
		Return(signed, &entity.User{
			ID:           "1",
			Email:        "alice@example.com",
			PasswordHash: "password",
			Role:         "operation",
		}, nil)

	mockUserUC := mock.NewMockPaymentUsecase(ctrl)
	mockUserUC.EXPECT().
		ListPayment("completed", "-created_at", 10, 1).
		Return([]*entity.Payment{
			{ID: "1",
				Merchant:  "merchant 1",
				Status:    "completed",
				Amount:    100,
				CreatedAt: time.Now()},
		}, 5, 2, 3, nil)
	mockUserUC.EXPECT().
		ReviewPayment(gomock.Any(), "1").
		Return("Success Review", nil)

	srv := srv.NewServer(mockUserUC, mockAuthUC)
	ts := httptest.NewServer(srv.Routes())
	defer ts.Close()

	login := map[string]string{"email": "alice@example.com", "password": "password"}
	bl, _ := json.Marshal(login)
	resLogin, err := http.Post(ts.URL+"/dashboard/v1/auth/login", "application/json", bytes.NewReader(bl))
	require.NoError(t, err)
	defer resLogin.Body.Close()
	require.Equal(t, http.StatusOK, resLogin.StatusCode)

	var resp map[string]string
	json.NewDecoder(resLogin.Body).Decode(&resp)
	respToken := resp["token"]
	require.NotEmpty(t, respToken)

	client := &http.Client{}
	reqGetPayment, _ := http.NewRequest("GET", ts.URL+"/dashboard/v1/payments?status=completed&sort=-created_at&limit=10&offset=1", nil)
	reqGetPayment.Header.Set("Authorization", "Bearer "+respToken)
	resGetPayment, err := client.Do(reqGetPayment)
	require.NoError(t, err)
	defer resGetPayment.Body.Close()
	require.Equal(t, http.StatusOK, resGetPayment.StatusCode)

	client2 := &http.Client{}
	reqPaymentReview, _ := http.NewRequest("PUT", ts.URL+"/dashboard/v1/payment/1/review", nil)
	reqPaymentReview.Header.Set("Authorization", "Bearer "+respToken)
	resPaymentReview, err := client2.Do(reqPaymentReview)
	require.NoError(t, err)
	defer resPaymentReview.Body.Close()
	require.Equal(t, http.StatusOK, resPaymentReview.StatusCode)
}
