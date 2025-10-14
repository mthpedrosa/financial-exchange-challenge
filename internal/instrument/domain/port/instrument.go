package port

import (
	"context"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/entity"
)

type InstrumentRepository interface {
	Create(ctx context.Context, instrument *entity.Instrument) (string, error)
	Update(ctx context.Context, instrument *entity.Instrument) error
	DeleteByID(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.Instrument, error)
	FindAll(ctx context.Context, filter *entity.InstrumentFilter) ([]*entity.Instrument, error)
	FindByAssets(ctx context.Context, baseAsset, quoteAsset string) (*entity.Instrument, error)
}
