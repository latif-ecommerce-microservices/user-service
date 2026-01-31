package user

import (
	"context"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"
	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/model"
	"github.com/latif-ecommerce-microservices/user-service/internal/helper/passwordhelper"
	"github.com/latif-ecommerce-microservices/user-service/pkg/customerror"
	"github.com/latif-ecommerce-microservices/user-service/pkg/ptr"
	"time"

	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
)

type service struct {
	userRepo repository.UserRepositoryProvider
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, request dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	user, err := s.userRepo.FindActiveByID(ctx, id)

	if err != nil {
		return nil, customerror.InternalServerError.
			WithCause(err).
			WithStackTrace()
	}

	if user == nil {
		return nil, customerror.UserNotFoundError.
			WithLocator(customerror.WhereAmI())
	}

	if request.Email != nil && *request.Email != user.Email {
		existing, err := s.userRepo.GetActiveUserByEmail(ctx, *request.Email)
		if err != nil {
			return nil, customerror.InternalServerError.
				WithCause(err).
				WithStackTrace()
		}
		if existing != nil {
			return nil, customerror.EmailAlreadyExistError
		}
		user.Email = *request.Email
	}

	if request.Name != nil {
		user.Name = *request.Name
	}

	if request.PhoneNumber != nil {
		user.PhoneNumber = request.PhoneNumber
	}

	now := time.Now()
	user.UpdatedAt = &now

	updatedUser, err := s.userRepo.UpdateUser(ctx, *user)
	if err != nil {
		return nil, customerror.InternalServerError.
			WithCause(err).
			WithStackTrace()
	}

	return &dto.UpdateUserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}, nil
}

func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return customerror.InternalServerError.
			WithCause(err).
			WithStackTrace().
			WithLocator(customerror.WhereAmI())
	}

	if user == nil {
		return customerror.UserNotFoundError.
			WithLocator(customerror.WhereAmI())
	}

	if user.DeletedAt != nil {
		return nil
	}

	now := time.Now()
	user.DeletedAt = &now

	if _, err := s.userRepo.UpdateUser(ctx, *user); err != nil {
		return err
	}

	return nil
}

func (s *service) CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	currentUser, err := s.userRepo.GetActiveUserByEmail(ctx, request.Email)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	if currentUser != nil {
		return dto.CreateUserResponse{}, customerror.EmailAlreadyExistError
	}

	newID, _ := uuid.NewV7()
	newUser := model.Users{
		ID:          newID,
		Name:        request.Name,
		Email:       request.Email,
		PhoneNumber: ptr.String(request.PhoneNumber),
		Password:    passwordhelper.HashPassword(request.Password),
		CreatedAt:   time.Now(),
	}

	createdUser, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	return dto.CreateUserResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*dto.DetailUserResponse, error) {
	user, err := s.userRepo.FindActiveByID(ctx, id)
	if err != nil {
		return nil, customerror.InternalServerError.WithCause(err).WithStackTrace()
	}

	if user == nil {
		return nil, customerror.UserNotFoundError.WithLocator(customerror.WhereAmI())
	}

	return &dto.DetailUserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func (s *service) GetAllUsers(ctx context.Context, id uuid.UUID) (dto.ListUserResponse, error) {
	// check user admin or not
	// SOOONNN

	users, err := s.userRepo.FindAllActiveUser(ctx)
	if err != nil {
		return nil, customerror.InternalServerError.WithCause(err).WithStackTrace()
	}

	result := make(dto.ListUserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, dto.DetailUserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}

	return result, nil
}

func NewService(userRepo repository.UserRepositoryProvider) ServiceProvider {
	return &service{userRepo: userRepo}
}
