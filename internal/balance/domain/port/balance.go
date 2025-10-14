package port

import (
	"context"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/entity"
)

// BalanceRepository defines the contract for balance persistence.
type BalanceRepository interface {
	Create(ctx context.Context, balance entity.Balance) (string, error)
	FindByID(ctx context.Context, id string) (entity.Balance, error)
	FindByAccountAndAsset(ctx context.Context, accountID, asset string) (entity.Balance, error)
	Update(ctx context.Context, balance entity.Balance) error
	DeleteByID(ctx context.Context, id string) error
	GetAllByAccountID(ctx context.Context, accountID string) ([]entity.Balance, error)
}
