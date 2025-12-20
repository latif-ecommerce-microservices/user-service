package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/latif-ecommerce-microservices/user-service/internal/app/server"
	"github.com/latif-ecommerce-microservices/user-service/internal/config"
	"github.com/latif-ecommerce-microservices/user-service/pkg/container"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	logger := logging.CreateDefaultLogger(cfg.GetLogLevel().String())

	err = server.RunMigration(ctx, cfg, logger)
	if err != nil {
		panic(fmt.Sprintf("Migration error: %s", err.Error()))
	}

	logger.Info(fmt.Sprintf("[API] Starting API server at port %s...", cfg.AppHTTPPort))
	apiServer := server.NewAppServer(cfg, logger)
	servers := []*http.Server{apiServer.HTTPServer()}

	appContainer := container.New(
		apiServer,
		container.WithHTTPServer(servers),
	)
	err = appContainer.Start(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Server error: %s", err.Error()))
		cancel()
	} else {
		logger.Info("[API] API server shutting down.")
	}
}
