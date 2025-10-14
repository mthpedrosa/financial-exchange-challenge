package port

import (
	"context"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
)

type AccountRepository interface {
	GetAccounts(ctx context.Context, filters entity.AccountFilter) ([]entity.Account, error)
	FindByID(ctx context.Context, id string) (entity.Account, error)
	Create(ctx context.Context, account entity.Account) (string, error)
	DeleteByID(ctx context.Context, id string) error
	Update(ctx context.Context, account entity.Account) error
}
