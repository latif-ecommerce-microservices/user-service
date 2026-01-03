package server

import (
	"context"

	"github.com/latif-ecommerce-microservices/user-service/internal/config"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigration(ctx context.Context, cfg *config.Config, log *logging.Logger) error {
	if !cfg.Database.EnableMigrations {
		return nil
	}

	db, err := NewDbConnection(ctx, log, &cfg.Database)
	if err != nil {
		log.Error(err.Error())
	}
	defer db.Close()

	log.Info("Running database migration...")

	migrationDir := "./internal/generated/migrations"
	if err := goose.Up(db, migrationDir); err != nil {
		return err
	}
	log.Info("Database migration completed successfully.")

	return nil
}
