package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required"`
	PhoneNumber          string `json:"phone_number" validate:"required"`
	AnnualIncomeID       int32  `json:"annual_income_id" validate:"required"`
	EntityBusinessTypeID int32  `json:"entity_business_type_id" validate:"required"`
	CreatorTypeID        int32  `json:"creator_type_id" validate:"required"`
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}
