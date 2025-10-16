package entity_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/entity"
	"github.com/stretchr/testify/assert"
)

func newBigFloat(s string) *dto.BigFloat {
	f, _, _ := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	return &dto.BigFloat{Float: f}
}

func TestToEntity(t *testing.T) {
	t.Run("should convert valid DTO to entity successfully", func(t *testing.T) {
		// arrange
		price := newBigFloat("100.50")
		quantity := newBigFloat("10.0")

		request := dto.CreateOrderRequest{
			AccountID:    "acc-123",
			InstrumentID: "inst-456",
			Type:         "BUY",
			Price:        price,
			Quantity:     quantity,
		}

		// act
		order, err := entity.ToEntity(request)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, "acc-123", order.AccountID)
		assert.Equal(t, "inst-456", order.InstrumentID)
		assert.Equal(t, entity.OrderTypeBuy, order.Type)
		assert.Equal(t, entity.OrderStatusOpen, order.Status)

		// compare big.Float safely
		assert.Zero(t, price.Float.Cmp(order.Price), "Price should match")
		assert.Zero(t, quantity.Float.Cmp(order.Quantity), "Quantity should match")
		assert.Zero(t, quantity.Float.Cmp(order.RemainingQuantity), "RemainingQuantity should be initialized with Quantity")
	})

	t.Run("should return error for invalid DTO", func(t *testing.T) {
		// arrange
		request := dto.CreateOrderRequest{
			AccountID: "", // Campo inválido
		}

		// act
		order, err := entity.ToEntity(request)

		// assert
		assert.Error(t, err)
		assert.Nil(t, order)
	})
}

func TestOrder_ToDTO(t *testing.T) {
	// Arrange
	price, _ := new(big.Float).SetString("200.0")
	quantity, _ := new(big.Float).SetString("50.0")
	remaining, _ := new(big.Float).SetString("20.0")

	order := &entity.Order{
		ID:                "order-789",
		AccountID:         "acc-123",
		InstrumentID:      "inst-456",
		Type:              entity.OrderTypeSell,
		Status:            entity.OrderStatusPartiallyFilled,
		Price:             price,
		Quantity:          quantity,
		RemainingQuantity: remaining,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// act
	orderDTO := order.ToDTO()

	// assert
	assert.Equal(t, "order-789", orderDTO.ID)
	assert.Equal(t, "acc-123", orderDTO.AccountID)
	assert.Equal(t, "SELL", orderDTO.Type)
	assert.Equal(t, "PARTIALLY_FILLED", orderDTO.Status)
	assert.Zero(t, price.Cmp(&orderDTO.Price))
	assert.Zero(t, quantity.Cmp(&orderDTO.Quantity))
	assert.Zero(t, remaining.Cmp(&orderDTO.RemainingQuantity))
	assert.Equal(t, order.CreatedAt, orderDTO.CreatedAt)
}

func TestToListDTO(t *testing.T) {
	t.Run("should convert a slice of entities to a slice of DTOs", func(t *testing.T) {
		// arrange
		price, _ := new(big.Float).SetString("100.0")
		quantity, _ := new(big.Float).SetString("10.0")

		orders := []entity.Order{
			{ID: "order-1", Price: price, Quantity: quantity, RemainingQuantity: quantity},
			{ID: "order-2", Price: price, Quantity: quantity, RemainingQuantity: quantity},
		}

		// act
		dtos := entity.ToListDTO(orders)

		// assert
		assert.NotNil(t, dtos)
		assert.Len(t, dtos, 2)
		assert.Equal(t, "order-1", dtos[0].ID)
		assert.Equal(t, "order-2", dtos[1].ID)
	})

	t.Run("should return an empty slice for an empty input slice", func(t *testing.T) {
		// arrange
		orders := []entity.Order{}

		// act
		dtos := entity.ToListDTO(orders)

		// assert
		assert.NotNil(t, dtos)
		assert.Empty(t, dtos)
	})
}
