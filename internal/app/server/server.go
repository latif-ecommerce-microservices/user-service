package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/latif-ecommerce-microservices/user-service/internal/config"
	"github.com/latif-ecommerce-microservices/user-service/internal/transport/http/handler"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

type Server struct {
	srv    *http.Server
	router chi.Router
	cfg    *config.Config
	logger *logging.Logger

	InternalConnection *InternalConnection
}

func NewAppServer(cfg *config.Config, logger *logging.Logger) *Server {
	router := chi.NewMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppHTTPPort),
		Handler: router,
	}

	return &Server{
		cfg:    cfg,
		router: router,
		srv:    srv,
		logger: logger,
	}
}

func (s *Server) BeforeStart(ctx context.Context) error {
	internalClient, err := NewInternalConnection(ctx, s.logger, &s.cfg.Database)
	if err != nil {
		return err
	}
	repository := NewRepository(internalClient)
	service := NewService(repository)
	newValidator := validator.New()

	userHandler := handler.NewUserHandler(service.UserService, newValidator, s.logger)
	userHandler.RegisterRoutes(s.router)

	authHandler := handler.NewAuthHandler(service.AuthService, newValidator, s.logger)
	authHandler.RegisterRoutes(s.router)

	s.InternalConnection = &internalClient
	return nil
}

func (s *Server) AfterStart(ctx context.Context) error {
	if s.InternalConnection != nil {
		err := s.InternalConnection.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) HTTPServer() *http.Server {
	return s.srv
}
