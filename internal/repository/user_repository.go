package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/model"
)

type UserRepositoryProvider interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Users, error)
	GetUserByEmail(ctx context.Context, email string) (*model.Users, error)
	CreateUser(ctx context.Context, user model.Users) (model.Users, error)
	UpdateUser(ctx context.Context, user model.Users) (model.Users, error)
}
