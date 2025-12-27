package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
}

type UpdateUserRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=3"`
	Email       *string `json:"email" validate:"omitempty,email"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty"`
}

type UpdateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt string    `json:"updated_at"`
}

type DetailUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   *string   `json:"created_at"`
	UpdatedAt   *string   `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}
