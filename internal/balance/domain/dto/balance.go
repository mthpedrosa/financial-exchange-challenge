package dto

import (
	"encoding/json"
	"math/big"

	"github.com/go-playground/validator/v10"
)

type BalanceDTO struct {
	AccountID string    `json:"account_id"`
	Amount    big.Float `json:"amount"`
	Asset     string    `json:"asset"`
}

type GetBalanceRequest struct {
	AccountID string `json:"account_id" binding:"required"`
}

type GetBalanceResponse struct {
	AccountID string    `json:"account_id"`
	Amount    big.Float `json:"amount"`
	Asset     string    `json:"asset"`
}

type CreateBalanceRequest struct {
	AccountID string    `json:"account_id" validate:"required"`
	Asset     string    `json:"asset" validate:"required"`
	Amount    *BigFloat `json:"amount" validate:"required"`
}

type CreateBalanceResponse struct {
	ID string `json:"id"`
}

type UpdateBalanceRequest struct {
	Amount *big.Float `json:"amount" validate:"required"`
}

type BalanceListDTO struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Asset     string    `json:"asset"`
	Amount    big.Float `json:"amount"`
}

func (r *CreateBalanceRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *UpdateBalanceRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type BigFloat struct {
	*big.Float
}

func (b *BigFloat) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		f, _, err := big.ParseFloat(s, 10, 256, big.ToNearestEven)
		if err != nil {
			return err
		}
		b.Float = f
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	b.Float = big.NewFloat(f)
	return nil
}
