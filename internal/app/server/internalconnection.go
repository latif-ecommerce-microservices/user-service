package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/latif-ecommerce-microservices/user-service/internal/config"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

type InternalConnection struct {
	DB *sql.DB
}

func NewInternalConnection(ctx context.Context, log *logging.Logger, cfg *config.DatabaseConfig) (InternalConnection, error) {
	db, err := NewDbConnection(ctx, log, cfg)
	if err != nil {
		return InternalConnection{}, err
	}

	return InternalConnection{DB: db}, nil
}

func NewDbConnection(ctx context.Context, log *logging.Logger, cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open database connection: %s", err.Error()))
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(fmt.Sprintf("failed to ping db: %s", err.Error()))
	}

	return db, nil
}

func (ic *InternalConnection) Close() error {
	return ic.DB.Close()
}
