package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/latif-ecommerce-microservices/user-service/internal/config"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
)

type InternalConnection struct {
	DB          *sql.DB
	RedisClient *redis.Client
}

func NewInternalConnection(db *sql.DB, rdb *redis.Client) InternalConnection {
	return InternalConnection{
		DB:          db,
		RedisClient: rdb,
	}
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

func NewRedisConnection(ctx context.Context, log *logging.Logger, cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to redis: %s", err.Error()))
		return nil, err
	}

	return rdb, nil
}

func (ic *InternalConnection) Close() error {
	return ic.DB.Close()
}
