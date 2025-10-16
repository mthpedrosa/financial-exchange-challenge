package dto_test

import (
	"math/big"
	"testing"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/dto"
	"github.com/stretchr/testify/assert"
)

// newBigFloat
func newBigFloat(s string) *dto.BigFloat {
	f, _, _ := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	return &dto.BigFloat{Float: f}
}

func TestCreateOrderRequest_Validate_Valid(t *testing.T) {
	req := &dto.CreateOrderRequest{
		AccountID:    "acc-123",
		InstrumentID: "inst-456",
		Type:         "BUY",
		Price:        newBigFloat("150.50"),
		Quantity:     newBigFloat("10.5"),
	}
	err := req.Validate()
	assert.NoError(t, err)
}

func TestCreateOrderRequest_Validate_MissingAccountID(t *testing.T) {
	req := &dto.CreateOrderRequest{
		AccountID:    "", // invalid (missing)
		InstrumentID: "inst-456",
		Type:         "SELL",
		Price:        newBigFloat("150.50"),
		Quantity:     newBigFloat("10.5"),
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestCreateOrderRequest_Validate_MissingInstrumentID(t *testing.T) {
	req := &dto.CreateOrderRequest{
		AccountID:    "acc-123",
		InstrumentID: "", // invalid (missing)
		Type:         "BUY",
		Price:        newBigFloat("150.50"),
		Quantity:     newBigFloat("10.5"),
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestCreateOrderRequest_Validate_InvalidType(t *testing.T) {
	req := &dto.CreateOrderRequest{
		AccountID:    "acc-123",
		InstrumentID: "inst-456",
		Type:         "INVALID_TYPE", // invalid type
		Price:        newBigFloat("150.50"),
		Quantity:     newBigFloat("10.5"),
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestCreateOrderRequest_Validate_MissingType(t *testing.T) {
	req := &dto.CreateOrderRequest{
		AccountID:    "acc-123",
		InstrumentID: "inst-456",
		Type:         "", // invalid (missing)
		Price:        newBigFloat("150.50"),
		Quantity:     newBigFloat("10.5"),
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestCreateOrderRequest_Validate_MissingPrice(t *testing.T) {
	// missing price
	req := &dto.CreateOrderRequest{
		AccountID:    "acc-123",
		InstrumentID: "inst-456",
		Type:         "BUY",
		Quantity:     newBigFloat("10.5"),
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestCreateOrderRequest_Validate_MissingQuantity(t *testing.T) {
	// missing quantity
	req := &dto.CreateOrderRequest{
		AccountID:    "acc-123",
		InstrumentID: "inst-456",
		Type:         "SELL",
		Price:        newBigFloat("150.50"),
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestBigFloat_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name          string
		inputJSON     []byte
		expectedValue string
		expectError   bool
	}{
		{
			name:          "Valid String Input",
			inputJSON:     []byte(`"150.75"`),
			expectedValue: "150.75",
			expectError:   false,
		},
		{
			name:          "Valid Numeric Input",
			inputJSON:     []byte(`250.5`),
			expectedValue: "250.5",
			expectError:   false,
		},
		{
			name:        "Invalid String Input",
			inputJSON:   []byte(`"not-a-number"`),
			expectError: true,
		},
		{
			name:        "Invalid JSON Type",
			inputJSON:   []byte(`{"key":"value"}`), // Um objeto, não uma string ou número
			expectError: true,
		},
		{
			name:          "Valid Integer String",
			inputJSON:     []byte(`"1000"`),
			expectedValue: "1000",
			expectError:   false,
		},
		{
			name:          "Valid Integer Number",
			inputJSON:     []byte(`2000`),
			expectedValue: "2000",
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b dto.BigFloat

			err := b.UnmarshalJSON(tc.inputJSON)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, b.Float.String())
			}
		})
	}
}
