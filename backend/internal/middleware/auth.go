package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fajrinajiseno/mygolangapp/internal/config"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx context.Context, in *openapi3filter.AuthenticationInput) error {
	req := in.RequestValidationInput.Request

	sub, err := GetTokenSub(req)
	if err != nil {
		return in.NewError(err)
	}

	newReq := req.WithContext(context.WithValue(req.Context(), config.ContextUserID, sub))
	in.RequestValidationInput.Request = newReq

	return nil
}

func GetTokenSub(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("missing Authorization header")
	}

	const authLength = 2
	parts := strings.SplitN(auth, " ", authLength)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errors.New("invalid Authorization header")
	}
	tokenString := parts[1]

	tkn, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return config.JwtSecret, nil
	})
	if err != nil || !tkn.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", errors.New("token missing sub")
	}
	return sub, nil
}

func GetUserID(ctx context.Context) string {
	v := ctx.Value(config.ContextUserID)
	if v == nil {
		return ""
	}
	id, _ := v.(string)
	return id
}
