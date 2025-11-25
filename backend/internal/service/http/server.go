package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fajrinajiseno/mygolangapp/internal/config"
	"github.com/fajrinajiseno/mygolangapp/internal/middleware"
	"github.com/fajrinajiseno/mygolangapp/internal/openapigen"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	oapinethttpmw "github.com/oapi-codegen/nethttp-middleware"
	swgui "github.com/swaggest/swgui/v5"
	"sigs.k8s.io/yaml"
)

type Server struct {
	router http.Handler
}

const (
	readTimeout  = 10
	writeTimeout = 10
	idleTimeout  = 60
	corsMaxAge   = 300
)

func NewServer(apiHandler openapigen.ServerInterface, openapiYamlPath string) *Server {
	swagger, err := openapigen.GetSwagger()
	if err != nil {
		log.Fatalf("failed to load swagger: %v", err)
	}
	openapiJSON, err := loadOpenAPIAsJSON(openapiYamlPath)
	if err != nil {
		log.Fatalf("failed to loadOpenAPIAsJSON: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ContextMiddleware)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.Cors},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           corsMaxAge,
	}))

	r.Get("/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(openapiJSON)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})

	r.Handle("/docs/*", swgui.New("Dashboard API Docs", "/openapi.json", "/docs/"))
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusTemporaryRedirect)
	})

	r.Route("/", func(api chi.Router) {
		api.Use(oapinethttpmw.OapiRequestValidatorWithOptions(
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
				DoNotValidateServers:  true,
				SilenceServersWarning: true,
			},
		))
		openapigen.HandlerFromMux(apiHandler, api)
	})

	return &Server{
		router: r,
	}
}

func (s *Server) Start(addr string) {
	service := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
	}
	go func() {
		log.Printf("listening on %s", addr)
		err := service.ListenAndServe()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down gracefully...")

	// Timeout for shutdown
	const shutdownTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := service.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server stopped cleanly ✔")
}

func (s *Server) Routes() http.Handler {
	return s.router
}

func loadOpenAPIAsJSON(yamlPath string) ([]byte, error) {
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return nil, err
	}
	var pretty json.RawMessage
	if err := json.Unmarshal(jsonData, &pretty); err == nil {
		out, _ := json.MarshalIndent(pretty, "", "  ")
		return out, nil
	}
	return jsonData, nil
}
