package entity

import (
	"math/big"
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/dto"
)

type Balance struct {
	ID        string
	AccountID string
	Asset     string
	Amount    *big.Float
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToListDTO(balances []Balance) []dto.BalanceListDTO {
	dtos := make([]dto.BalanceListDTO, len(balances))
	for i, a := range balances {
		dtos[i] = dto.BalanceListDTO{
			ID:        a.ID,
			AccountID: a.AccountID,
			Asset:     a.Asset,
			Amount:    *a.Amount,
		}
	}
	return dtos
}

func ToEntity(request dto.CreateBalanceRequest) (*Balance, error) {
	if err := request.Validate(); err != nil {
		return &Balance{}, err
	}

	return &Balance{
		AccountID: request.AccountID,
		Asset:     request.Asset,
		Amount:    request.Amount.Float,
	}, nil
}

func ToEntityUpdate(request dto.UpdateBalanceRequest) (*Balance, error) {
	if err := request.Validate(); err != nil {
		return &Balance{}, err
	}

	return &Balance{
		Amount: request.Amount,
	}, nil
}
