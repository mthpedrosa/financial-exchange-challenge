package entity

import (
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/dto"
)

type Instrument struct {
	ID         string    `json:"id"`
	BaseAsset  string    `json:"base_asset"`
	QuoteAsset string    `json:"quote_asset"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type InstrumentFilter struct {
	BaseAsset  string `json:"base_asset"`
	QuoteAsset string `json:"quote_asset"`
}

func (i *Instrument) ToDTO() dto.InstrumentDTO {
	return dto.InstrumentDTO{
		ID:         i.ID,
		BaseAsset:  i.BaseAsset,
		QuoteAsset: i.QuoteAsset,
		CreatedAt:  i.CreatedAt,
		UpdatedAt:  i.UpdatedAt,
	}
}

func ToEntity(dto dto.CreateInstrumentRequest) (*Instrument, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	return &Instrument{
		BaseAsset:  dto.BaseAsset,
		QuoteAsset: dto.QuoteAsset,
	}, nil
}

func ToListDTO(accounts []Instrument) []dto.InstrumentListDTO {
	dtos := make([]dto.InstrumentListDTO, len(accounts))
	for i, a := range accounts {
		dtos[i] = dto.InstrumentListDTO{
			ID:         a.ID,
			BaseAsset:  a.BaseAsset,
			QuoteAsset: a.QuoteAsset,
		}
	}
	return dtos
}

func ToEntityFilter(dto dto.InstrumentFilter) *InstrumentFilter {
	return &InstrumentFilter{
		BaseAsset:  dto.BaseAsset,
		QuoteAsset: dto.QuoteAsset,
	}
}
