package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"

	"github.com/latif-ecommerce-microservices/user-service/internal/config"

	authpb "github.com/latif-ecommerce-microservices/user-service/pkg/pb/auth"
	userpb "github.com/latif-ecommerce-microservices/user-service/pkg/pb/user"

	grpchandler "github.com/latif-ecommerce-microservices/user-service/internal/transport/grpc"
	httphandler "github.com/latif-ecommerce-microservices/user-service/internal/transport/http/handler"

	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

type Server struct {
	srv     *http.Server
	grpcSrv *grpc.Server
	router  chi.Router
	cfg     *config.Config
	logger  *logging.Logger

	InternalConnection *InternalConnection
}

func NewAppServer(cfg *config.Config, logger *logging.Logger) *Server {
	router := chi.NewMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.AppHTTPPort),
		Handler: router,
	}

	grpcSrv := grpc.NewServer()

	return &Server{
		cfg:     cfg,
		router:  router,
		srv:     srv,
		grpcSrv: grpcSrv,
		logger:  logger,
	}
}

func (s *Server) BeforeStart(ctx context.Context) error {
	db, err := NewDbConnection(ctx, s.logger, &s.cfg.Database)
	if err != nil {
		return err
	}

	rdb, err := NewRedisConnection(ctx, s.logger, &s.cfg.Redis)
	if err != nil {
		return err
	}

	internalClient := NewInternalConnection(db, rdb)

	repository := NewRepository(internalClient)
	service := NewService(repository)
	newValidator := validator.New()

	userHandler := httphandler.NewUserHandler(service.UserService, newValidator, s.logger)
	userHandler.RegisterRoutes(s.router)

	authHandler := httphandler.NewAuthHandler(service.AuthService, newValidator, s.logger)
	authHandler.RegisterRoutes(s.router)

	userGrpcHandler := grpchandler.NewUserHandler(service.UserService)
	authGrpcHandler := grpchandler.NewAuthHandler(service.AuthService)

	userpb.RegisterUserServiceServer(s.grpcSrv, userGrpcHandler)
	authpb.RegisterAuthServiceServer(s.grpcSrv, authGrpcHandler)

	//reflection.Register(s.grpcSrv)

	s.InternalConnection = &internalClient
	return nil
}

func (s *Server) AfterStart(ctx context.Context) error {

	if s.grpcSrv != nil {
		s.grpcSrv.GracefulStop()
	}

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

func (s *Server) GRPCServer() *grpc.Server {
	return s.grpcSrv
}
