package entity

import (
	"math/big"
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/dto"
)

type OrderType string
type OrderStatus string

const (
	OrderTypeBuy  OrderType = "BUY"
	OrderTypeSell OrderType = "SELL"

	OrderStatusOpen            OrderStatus = "OPEN"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCancelled       OrderStatus = "CANCELLED"
)

type Order struct {
	ID                string
	AccountID         string
	InstrumentID      string
	Type              OrderType
	Status            OrderStatus
	Price             *big.Float
	Quantity          *big.Float
	RemainingQuantity *big.Float
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ToEntity converts a CreateOrderRequest DTO to an Order entity.
func ToEntity(request dto.CreateOrderRequest) (*Order, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	return &Order{
		AccountID:         request.AccountID,
		InstrumentID:      request.InstrumentID,
		Type:              OrderType(request.Type),
		Status:            OrderStatusOpen,
		Price:             request.Price.Float,
		Quantity:          request.Quantity.Float,
		RemainingQuantity: request.Quantity.Float,
	}, nil
}

// ToDTO converts an Order entity to an OrderDTO.
func (o *Order) ToDTO() dto.OrderDTO {
	return dto.OrderDTO{
		ID:                o.ID,
		AccountID:         o.AccountID,
		InstrumentID:      o.InstrumentID,
		Type:              string(o.Type),
		Status:            string(o.Status),
		Price:             *o.Price,
		Quantity:          *o.Quantity,
		RemainingQuantity: *o.RemainingQuantity,
		CreatedAt:         o.CreatedAt,
		UpdatedAt:         o.UpdatedAt,
	}
}

// ToListDTO converts a slice of Order entities to a slice of OrderDTOs.
func ToListDTO(orders []Order) []dto.OrderDTO {
	dtos := make([]dto.OrderDTO, len(orders))
	for i, o := range orders {
		dtos[i] = o.ToDTO()
	}
	return dtos
}
