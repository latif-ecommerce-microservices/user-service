package dto

type CreateUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateUserRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=3"`
	Email       *string `json:"email" validate:"omitempty,email"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty"`
}
