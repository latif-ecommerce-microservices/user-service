package user

import (
	"context"
	"github.com/latif-ecommerce-microservices/user-service/internal/generated/go_check/public/model"

	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/repository"
)

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) Service {
	return &service{userRepo: userRepo}
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*model.Users, error) {
	return s.userRepo.FindByID(ctx, id)
}
