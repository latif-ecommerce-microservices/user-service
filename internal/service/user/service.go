package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/model"
)

type ServiceProvider interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Users, error)
}
