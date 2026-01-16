package auth

import (
	"context"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
)

type ServiceProvider interface {
	Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, request dto.LogoutRequest) error
}
