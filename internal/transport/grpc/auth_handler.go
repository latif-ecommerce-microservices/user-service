package grpc

import (
	"context"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"

	"github.com/latif-ecommerce-microservices/user-service/internal/service/auth"

	authpb "github.com/latif-ecommerce-microservices/user-service/internal/generated/pb/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	authService auth.ServiceProvider
}

func NewAuthHandler(authService auth.ServiceProvider) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {

	loginDto := dto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.authService.Login(ctx, loginDto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &authpb.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
