package repository

import (
	"math/big"
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/entity"
)

type BalanceModel struct {
	ID        string
	AccountID string
	Amount    *big.Float
	Asset     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToModel(b entity.Balance) *BalanceModel {
	return &BalanceModel{
		ID:        b.ID,
		AccountID: b.AccountID,
		Amount:    b.Amount,
		Asset:     b.Asset,
	}
}

func (b *BalanceModel) ToEntity() entity.Balance {
	return entity.Balance{
		ID:        b.ID,
		AccountID: b.AccountID,
		Asset:     b.Asset,
		Amount:    b.Amount,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
