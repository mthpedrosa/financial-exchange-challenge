package dto_test

import (
	"testing"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstrumentRequest_Validate_Valid(t *testing.T) {
	req := &dto.CreateInstrumentRequest{
		BaseAsset:  "BTC",
		QuoteAsset: "USD",
	}
	err := req.Validate()
	assert.NoError(t, err)
}

func TestCreateInstrumentRequest_Validate_MissingBaseAsset(t *testing.T) {
	req := &dto.CreateInstrumentRequest{
		BaseAsset:  "",
		QuoteAsset: "USD",
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestCreateInstrumentRequest_Validate_MissingQuoteAsset(t *testing.T) {
	req := &dto.CreateInstrumentRequest{
		BaseAsset:  "BTC",
		QuoteAsset: "",
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestCreateInstrumentRequest_Validate_MissingBoth(t *testing.T) {
	req := &dto.CreateInstrumentRequest{}
	err := req.Validate()
	assert.Error(t, err)
}
