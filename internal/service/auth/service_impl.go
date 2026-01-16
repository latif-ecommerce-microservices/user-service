package auth

import (
	"context"
	"fmt"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
	"github.com/latif-ecommerce-microservices/user-service/internal/helper/passwordhelper"
	"github.com/latif-ecommerce-microservices/user-service/internal/helper/tokenhelper"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
	"github.com/latif-ecommerce-microservices/user-service/pkg/customerror"
	"github.com/latif-ecommerce-microservices/user-service/pkg/logging"
	"time"
)

type service struct {
	userRepo  repository.UserRepositoryProvider
	tokenRepo repository.TokenRepositoryProvider
	logger    *logging.Logger
}

func (s service) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetActiveUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, customerror.InternalServerError.
			WithCause(err).
			WithLocator(customerror.WhereAmI())
	}

	if user == nil {
		return nil, customerror.InvalidCredentialError.
			WithCause(err).
			WithLocator(customerror.WhereAmI())
	}

	if err := passwordhelper.ComparePassword(user.Password, request.Password); err != nil {
		return nil, customerror.InvalidCredentialError.
			WithCause(err).
			WithLocator(customerror.WhereAmI())
	}

	accessToken, err := tokenhelper.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, customerror.InternalServerError.
			WithCause(err).
			WithLocator(customerror.WhereAmI())
	}

	refreshToken, err := tokenhelper.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, customerror.InternalServerError.
			WithCause(err).
			WithLocator(customerror.WhereAmI())
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, user.ID.String(), refreshToken, 7*24*time.Hour)
	if err != nil {
		return nil, customerror.InternalServerError.
			WithCause(err).
			WithLocator(customerror.WhereAmI())
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s service) Logout(ctx context.Context, request dto.LogoutRequest) error {
	userID := request.UserID

	err := s.tokenRepo.DeleteRefreshToken(ctx, userID)
	if err != nil {
		return customerror.InternalServerError.
			WithCause(err).
			WithLocator(customerror.WhereAmI())
	}

	if request.AccessToken != "" {
		err = s.tokenRepo.BlacklistAccessToken(ctx, request.AccessToken, 15*time.Minute)
		if err != nil {
			s.logger.Error(fmt.Sprintf("failed to blacklist token: %v", err))
		}
	}

	return nil
}

func NewService(userRepo repository.UserRepositoryProvider, tokenRepo repository.TokenRepositoryProvider) ServiceProvider {
	return &service{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}
