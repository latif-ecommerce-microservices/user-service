package user

import (
	"context"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/model"
	"github.com/latif-ecommerce-microservices/user-service/internal/helper/passwordhelper"
	"github.com/latif-ecommerce-microservices/user-service/pkg/customerror"
	"time"

	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
)

type service struct {
	userRepo repository.UserRepositoryProvider
}

func NewService(userRepo repository.UserRepositoryProvider) ServiceProvider {
	return &service{userRepo: userRepo}
}

func (s *service) CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	currentUser, err := s.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	if currentUser != nil {
		return dto.CreateUserResponse{}, customerror.EmailAlreadyExistError
	}

	newID, _ := uuid.NewV7()
	newUser := model.Users{
		ID:        newID,
		Name:      request.Name,
		Email:     request.Email,
		Password:  passwordhelper.HashPassword(request.Password),
		CreatedAt: time.Now(),
	}

	createdUser, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	return dto.CreateUserResponse{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt.Format(time.DateTime),
	}, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*dto.DetailUserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.DetailUserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
