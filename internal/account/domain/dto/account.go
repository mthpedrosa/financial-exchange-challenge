package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type AccountDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountListDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateAccountRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateAccountRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type AccountFilter struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type CreateAccountResponse struct {
	ID string `json:"id"`
}

func (c *CreateAccountRequest) Validate() error {
	return validator.New().Struct(c)
}

func (u *UpdateAccountRequest) Validate() error {
	return validator.New().Struct(u)
}
