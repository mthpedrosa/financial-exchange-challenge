package entity_test

import (
	"math/big"
	"testing"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestToListDTO(t *testing.T) {
	balances := []entity.Balance{
		{
			ID:        "id-1",
			AccountID: "acc-1",
			Asset:     "BTC",
			Amount:    big.NewFloat(10.5),
		},
		{
			ID:        "id-2",
			AccountID: "acc-2",
			Asset:     "ETH",
			Amount:    big.NewFloat(20.0),
		},
	}
	dtos := entity.ToListDTO(balances)
	assert.Len(t, dtos, 2)
	assert.Equal(t, "id-1", dtos[0].ID)
	assert.Equal(t, "acc-1", dtos[0].AccountID)
	assert.Equal(t, "id-2", dtos[1].ID)
	assert.Equal(t, "acc-2", dtos[1].AccountID)
}

func TestToEntity_Valid(t *testing.T) {
	req := dto.CreateBalanceRequest{
		AccountID: "acc-uuid",
		Asset:     "BTC",
		Amount:    dto.BigFloat{big.NewFloat(123.45)},
	}
	b, err := entity.ToEntity(req)
	assert.NoError(t, err)
	assert.Equal(t, "acc-uuid", b.AccountID)
	assert.Equal(t, "BTC", b.Asset)
	assert.Equal(t, "123.45", b.Amount.Text('f', 2))
}

func TestToEntity_Invalid(t *testing.T) {
	req := dto.CreateBalanceRequest{}
	b, err := entity.ToEntity(req)
	assert.Error(t, err)
	assert.NotNil(t, b)
}

func TestToEntityUpdate_Valid(t *testing.T) {
	req := dto.UpdateBalanceRequest{
		Amount: *big.NewFloat(99.99),
	}
	b, err := entity.ToEntityUpdate(req)
	assert.NoError(t, err)
	assert.Equal(t, "99.99", b.Amount.Text('f', 2))
}

func TestToEntityUpdate_Invalid(t *testing.T) {
	req := dto.UpdateBalanceRequest{}
	b, err := entity.ToEntityUpdate(req)
	assert.Error(t, err)
	assert.NotNil(t, b)
}
