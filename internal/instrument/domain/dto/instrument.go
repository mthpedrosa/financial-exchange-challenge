package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type InstrumentDTO struct {
	ID         string    `json:"id"`
	BaseAsset  string    `json:"base_asset"`
	QuoteAsset string    `json:"quote_asset"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateInstrumentRequest struct {
	BaseAsset  string `json:"base_asset" validate:"required"`
	QuoteAsset string `json:"quote_asset" validate:"required"`
}

type InstrumentFilter struct {
	BaseAsset  string `query:"base_asset"`
	QuoteAsset string `query:"quote_asset"`
}

type InstrumentListDTO struct {
	ID         string `json:"id"`
	BaseAsset  string `json:"base_asset"`
	QuoteAsset string `json:"quote_asset"`
}

func (r *CreateInstrumentRequest) Validate() error {
	return validator.New().Struct(r)
}
