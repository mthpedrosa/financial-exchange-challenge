package repository

import (
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/entity"
)

type InstrumentModel struct {
	ID         string    `json:"id"`
	BaseAsset  string    `json:"base_asset"`
	QuoteAsset string    `json:"quote_asset"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToModel(instrument *entity.Instrument) *InstrumentModel {
	return &InstrumentModel{
		ID:         instrument.ID,
		BaseAsset:  instrument.BaseAsset,
		QuoteAsset: instrument.QuoteAsset,
		CreatedAt:  instrument.CreatedAt,
		UpdatedAt:  instrument.UpdatedAt,
	}
}

func ToEntity(model *InstrumentModel) *entity.Instrument {
	return &entity.Instrument{
		ID:         model.ID,
		BaseAsset:  model.BaseAsset,
		QuoteAsset: model.QuoteAsset,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}
