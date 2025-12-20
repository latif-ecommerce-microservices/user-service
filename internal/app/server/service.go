package server

import "github.com/latif-ecommerce-microservices/user-service/internal/service/user"

type Service struct {
	UserService user.ServiceProvider
}

func NewService(repository Repository) Service {
	return Service{
		UserService: user.NewService(repository.UserRepository),
	}
}
