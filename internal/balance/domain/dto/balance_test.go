package dto_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/dto"
	"github.com/stretchr/testify/assert"
)

func TestCreateBalanceRequest_Validate_Valid(t *testing.T) {
	req := dto.CreateBalanceRequest{
		AccountID: "account-uuid",
		Asset:     "BTC",
		Amount:    dto.BigFloat{big.NewFloat(100.5)},
	}
	err := req.Validate()
	assert.NoError(t, err)
}

func TestCreateBalanceRequest_Validate_Invalid(t *testing.T) {
	req := dto.CreateBalanceRequest{
		AccountID: "",
		Asset:     "",
		Amount:    dto.BigFloat{big.NewFloat(0)},
	}
	err := req.Validate()
	assert.Error(t, err)
}

func TestUpdateBalanceRequest_Validate_Valid(t *testing.T) {
	req := dto.UpdateBalanceRequest{
		Amount: *big.NewFloat(200.0),
	}
	err := req.Validate()
	assert.NoError(t, err)
}

func TestUpdateBalanceRequest_Validate_Invalid(t *testing.T) {
	req := dto.UpdateBalanceRequest{}
	err := req.Validate()
	assert.Error(t, err)
}

func TestBigFloat_UnmarshalJSON_Number(t *testing.T) {
	var bf dto.BigFloat
	err := json.Unmarshal([]byte(`123.45`), &bf)
	assert.NoError(t, err)
	assert.Equal(t, "123.45", bf.Float.Text('f', 2))
}

func TestBigFloat_UnmarshalJSON_String(t *testing.T) {
	var bf dto.BigFloat
	err := json.Unmarshal([]byte(`"678.90"`), &bf)
	assert.NoError(t, err)
	assert.Equal(t, "678.90", bf.Float.Text('f', 2))
}
