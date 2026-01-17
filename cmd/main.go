package main

import (
	"context"
	"fmt"
	"net"

	"github.com/latif-ecommerce-microservices/user-service/internal/app/server"
	"github.com/latif-ecommerce-microservices/user-service/internal/config"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	logger := logging.CreateDefaultLogger(cfg.GetLogLevel().String())

	err = server.RunMigration(ctx, cfg, logger)
	if err != nil {
		panic(fmt.Sprintf("Migration error: %s", err.Error()))
	}

	appServer := server.NewAppServer(cfg, logger)

	if err := appServer.BeforeStart(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to init server: %s", err.Error()))
	}

	grpcPort := cfg.GRPCPort
	if grpcPort == "" {
		grpcPort = "8001"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to listen on port %s: %v", grpcPort, err))
	}

	logger.Info(fmt.Sprintf("[gRPC] User Service running on port %s...", grpcPort))

	if err := appServer.GRPCServer().Serve(lis); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to serve gRPC: %v", err))
	}

	if err := appServer.AfterStart(ctx); err != nil {
		logger.Error(fmt.Sprintf("Error during cleanup: %v", err))
	}
}
