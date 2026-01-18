package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"

	"github.com/latif-ecommerce-microservices/user-service/internal/config"

	authpb "github.com/latif-ecommerce-microservices/user-service/pkg/pb/auth"
	userpb "github.com/latif-ecommerce-microservices/user-service/pkg/pb/user"

	grpchandler "github.com/latif-ecommerce-microservices/user-service/internal/transport/grpc"
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
		Addr:    fmt.Sprintf(":%s", cfg.GRPCPort),
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

func (s *Server) GRPCServer() *grpc.Server {
	return s.grpcSrv
}
