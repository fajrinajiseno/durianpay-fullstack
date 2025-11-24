package http

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fajrinajiseno/mygolangapp/internal/middleware"
	"github.com/fajrinajiseno/mygolangapp/internal/openapigen"
	"github.com/fajrinajiseno/mygolangapp/internal/usecase"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	oapinethttpmw "github.com/oapi-codegen/nethttp-middleware"
)

type Server struct {
	paymentUC usecase.PaymentUsecase
	authUC    usecase.AuthUsecase
	router    http.Handler
}

const (
	readTimeout  = 10
	writeTimeout = 10
	idleTimeout  = 60
)

func NewServer(paymentUC usecase.PaymentUsecase, authUC usecase.AuthUsecase) *Server {
	swagger, err := openapigen.GetSwagger()
	if err != nil {
		log.Fatalf("failed to load swagger: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(oapinethttpmw.OapiRequestValidatorWithOptions(
		swagger,
		&oapinethttpmw.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: middleware.AuthMiddleware,
			},
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)

				resp := struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				}{
					Code:    statusCode,
					Message: message,
				}

				err := json.NewEncoder(w).Encode(resp)
				if err != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
			},
			SilenceServersWarning: true,
		},
	))
	r.Use(middleware.ContextMiddleware)

	apiHandler := &apiHandler{
		authUC:    authUC,
		paymentUC: paymentUC,
	}

	openapigen.HandlerFromMux(apiHandler, r)

	return &Server{
		paymentUC: paymentUC,
		authUC:    authUC,
		router:    r,
	}
}

func (s *Server) Start(addr string) {
	log.Printf("listening on %s", addr)
	service := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
	}
	err := service.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
