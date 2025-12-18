package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/model"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Users, error)
}
