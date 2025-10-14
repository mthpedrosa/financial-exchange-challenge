package dto

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/go-playground/validator/v10"
)

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

type CreateOrderRequest struct {
	AccountID    string   `json:"account_id" validate:"required"`
	InstrumentID string   `json:"instrument_id" validate:"required"`
	Type         string   `json:"type" validate:"required,oneof=BUY SELL"`
	Price        BigFloat `json:"price" validate:"required"`
	Quantity     BigFloat `json:"quantity" validate:"required"`
}

type CreateOrderResponse struct {
	ID string `json:"id"`
}

func (r *CreateOrderRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type OrderDTO struct {
	ID                string    `json:"id"`
	AccountID         string    `json:"account_id"`
	InstrumentID      string    `json:"instrument_id"`
	Type              string    `json:"type"`
	Status            string    `json:"status"`
	Price             big.Float `json:"price"`
	Quantity          big.Float `json:"quantity"`
	RemainingQuantity big.Float `json:"remaining_quantity"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
