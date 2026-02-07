package grpc

import (
	"context"
	"errors"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
	"github.com/latif-ecommerce-microservices/user-service/pkg/customerror"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/latif-ecommerce-microservices/user-service/internal/service/auth"

	authpb "github.com/latif-ecommerce-microservices/user-service/pkg/pb/auth"

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
		var ce customerror.Error
		if errors.As(err, &ce) {
			st := status.New(mapGRPCCode(ce), ce.Error())

			st, _ = st.WithDetails(
				&errdetails.ErrorInfo{
					Reason: ce.Code(),
				},
			)

			return nil, st.Err()
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &authpb.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (h *AuthHandler) Logout(
	ctx context.Context,
	req *authpb.LogoutRequest,
) (*authpb.LogoutResponse, error) {

	logoutDto := dto.LogoutRequest{
		RefreshToken: req.RefreshToken,
	}

	err := h.authService.Logout(ctx, logoutDto)
	if err != nil {
		var ce customerror.Error
		if errors.As(err, &ce) {
			st := status.New(mapGRPCCode(ce), ce.Error())

			st, _ = st.WithDetails(
				&errdetails.ErrorInfo{
					Reason: ce.Code(),
				},
			)

			return nil, st.Err()
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &authpb.LogoutResponse{}, nil
}
