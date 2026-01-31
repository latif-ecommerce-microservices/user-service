package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
)

type ServiceProvider interface {
	CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.CreateUserResponse, error)
	GetAllUsers(ctx context.Context, id uuid.UUID) (dto.ListUserResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.DetailUserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, request dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
