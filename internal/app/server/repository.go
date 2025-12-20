package server

import (
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository/pgsql"
)

type Repository struct {
	UserRepository repository.UserRepositoryProvider
}

func NewRepository(internalConnection InternalConnection) Repository {
	return Repository{
		UserRepository: pgsql.NewUserRepository(internalConnection.DB),
	}
}
