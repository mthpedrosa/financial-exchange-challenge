package repository

import (
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/entity"
)

type OrderModel struct {
	ID                string    `json:"id"`
	AccountID         string    `json:"account_id"`
	InstrumentID      string    `json:"instrument_id"`
	Type              string    `json:"type"`
	Status            string    `json:"status"`
	Price             string    `json:"price"`
	Quantity          string    `json:"quantity"`
	RemainingQuantity string    `json:"remaining_quantity"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func ToModel(entity entity.Order) *OrderModel {
	return &OrderModel{
		ID:                entity.ID,
		AccountID:         entity.AccountID,
		InstrumentID:      entity.InstrumentID,
		Type:              string(entity.Type),
		Status:            string(entity.Status),
		Price:             entity.Price.String(),
		Quantity:          entity.Quantity.String(),
		RemainingQuantity: entity.RemainingQuantity.String(),
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}
}
