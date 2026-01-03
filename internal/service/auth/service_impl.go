package auth

import (
	"context"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
	"github.com/latif-ecommerce-microservices/user-service/internal/helper/passwordhelper"
	"github.com/latif-ecommerce-microservices/user-service/internal/helper/tokenhelper"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
	"github.com/latif-ecommerce-microservices/user-service/pkg/customerror"
)

type service struct {
	userRepo repository.UserRepositoryProvider
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

	return &dto.LoginResponse{
		AuthToken:    accessToken,
		RefreshToken: refreshToken,
	}, nil
}

//func (s service) Logout(ctx context.Context, token string) error {
//	//TODO implement me
//	panic("implement me")
//}

func NewService(userRepo repository.UserRepositoryProvider) ServiceProvider {
	return &service{userRepo: userRepo}
}
