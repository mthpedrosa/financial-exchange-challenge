package dto

import "github.com/go-playground/validator/v10"

type AccountListDTO struct {
	ID   string
	Name string
}

type CreateAccountRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateAccountRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type AccountFilter struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type CreateAcountResponse struct {
	ID string `json:"id"`
}

func (c *CreateAccountRequest) Validate() error {
	return validator.New().Struct(c)
}

func (u *UpdateAccountRequest) Validate() error {
	return validator.New().Struct(u)
}
