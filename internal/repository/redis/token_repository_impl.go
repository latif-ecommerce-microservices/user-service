package redis

import (
	"context"
	"errors"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
	"github.com/redis/go-redis/v9"
	"time"
)

type tokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) repository.TokenRepositoryProvider {
	return &tokenRepository{client: client}
}

func (r *tokenRepository) SaveRefreshToken(ctx context.Context, userID string, token string, duration time.Duration) error {
	key := "refresh_token:" + userID
	return r.client.Set(ctx, key, token, duration).Err()
}

func (r *tokenRepository) BlacklistAccessToken(ctx context.Context, token string, duration time.Duration) error {
	key := "blacklist:" + token
	return r.client.Set(ctx, key, "1", duration).Err()
}

func (r *tokenRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := "blacklist:" + token

	_, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *tokenRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := "refresh_token:" + userID

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
