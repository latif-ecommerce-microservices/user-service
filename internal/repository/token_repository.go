package repository

import (
	"context"
	"time"
)

type TokenRepositoryProvider interface {
	SaveRefreshToken(ctx context.Context, userID string, token string, duration time.Duration) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
	BlacklistAccessToken(ctx context.Context, token string, duration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}
