package dto_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/dto"
)

func TestCreateAccountRequest_Validate_Valid(t *testing.T) {
    req := &dto.CreateAccountRequest{
        Name:  "Alice",
        Email: "alice@email.com",
    }
    err := req.Validate()
    assert.NoError(t, err)
}

func TestCreateAccountRequest_Validate_MissingName(t *testing.T) {
    req := &dto.CreateAccountRequest{
        Name:  "",
        Email: "alice@email.com",
    }
    err := req.Validate()
    assert.Error(t, err)
}

func TestCreateAccountRequest_Validate_InvalidEmail(t *testing.T) {
    req := &dto.CreateAccountRequest{
        Name:  "Alice",
        Email: "invalid-email",
    }
    err := req.Validate()
    assert.Error(t, err)
}

func TestUpdateAccountRequest_Validate_Valid(t *testing.T) {
    req := &dto.UpdateAccountRequest{
        Name:  "Bob",
        Email: "bob@email.com",
    }
    err := req.Validate()
    assert.NoError(t, err)
}

func TestUpdateAccountRequest_Validate_MissingEmail(t *testing.T) {
    req := &dto.UpdateAccountRequest{
        Name:  "Bob",
        Email: "",
    }
    err := req.Validate()
    assert.Error(t, err)
}

func TestUpdateAccountRequest_Validate_InvalidEmail(t *testing.T) {
    req := &dto.UpdateAccountRequest{
        Name:  "Bob",
        Email: "not-an-email",
    }
    err := req.Validate()
    assert.Error(t, err)
}