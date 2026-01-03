package server

import (
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository/pgsql"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository/redis"
)

type Repository struct {
	UserRepository  repository.UserRepositoryProvider
	TokenRepository repository.TokenRepositoryProvider
}

func NewRepository(internalConnection InternalConnection) Repository {
	return Repository{
		UserRepository:  pgsql.NewUserRepository(internalConnection.DB),
		TokenRepository: redis.NewTokenRepository(internalConnection.RedisClient),
	}
}
