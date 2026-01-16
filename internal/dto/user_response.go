package dto

import "github.com/google/uuid"

type CreateUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UpdateUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type DetailUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   *string   `json:"created_at"`
	UpdatedAt   *string   `json:"updated_at"`
}
