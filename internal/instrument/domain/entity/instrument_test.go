package entity_test

import (
	"testing"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestToEntity_Valid(t *testing.T) {
	req := dto.CreateInstrumentRequest{
		BaseAsset:  "BTC",
		QuoteAsset: "USD",
	}
	inst, err := entity.ToEntity(req)
	assert.NoError(t, err)
	assert.Equal(t, "BTC", inst.BaseAsset)
	assert.Equal(t, "USD", inst.QuoteAsset)
}

func TestToEntity_Invalid(t *testing.T) {
	req := dto.CreateInstrumentRequest{
		BaseAsset:  "",
		QuoteAsset: "",
	}
	inst, err := entity.ToEntity(req)
	assert.Error(t, err)
	assert.Nil(t, inst)
}

func TestToListDTO(t *testing.T) {
	instruments := []entity.Instrument{
		{ID: "uuid-1", BaseAsset: "BTC", QuoteAsset: "USD"},
		{ID: "uuid-2", BaseAsset: "ETH", QuoteAsset: "USD"},
	}
	dtos := entity.ToListDTO(instruments)
	assert.Len(t, dtos, 2)
	assert.Equal(t, "uuid-1", dtos[0].ID)
	assert.Equal(t, "BTC", dtos[0].BaseAsset)
	assert.Equal(t, "USD", dtos[0].QuoteAsset)
	assert.Equal(t, "uuid-2", dtos[1].ID)
	assert.Equal(t, "ETH", dtos[1].BaseAsset)
}

func TestToEntityFilter(t *testing.T) {
	filter := dto.InstrumentFilter{
		BaseAsset:  "BTC",
		QuoteAsset: "USD",
	}
	entityFilter := entity.ToEntityFilter(filter)
	assert.Equal(t, "BTC", entityFilter.BaseAsset)
	assert.Equal(t, "USD", entityFilter.QuoteAsset)
}

func TestInstrument_DefaultValues(t *testing.T) {
	inst := entity.Instrument{}
	assert.Empty(t, inst.ID)
	assert.Empty(t, inst.BaseAsset)
	assert.Empty(t, inst.QuoteAsset)
	assert.True(t, inst.CreatedAt.IsZero())
	assert.True(t, inst.UpdatedAt.IsZero())
}
