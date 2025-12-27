package user

import (
	"context"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"

	"github.com/google/uuid"
)

type ServiceProvider interface {
	CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.CreateUserResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.DetailUserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, request dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
