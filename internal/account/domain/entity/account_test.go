package entity_test

import (
	"testing"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestToEntity_Valid(t *testing.T) {
	req := dto.CreateAccountRequest{
		Name:  "Alice",
		Email: "alice@email.com",
	}
	acc, err := entity.ToEntity(req)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", acc.Name)
	assert.Equal(t, "alice@email.com", acc.Email)
}

func TestToEntity_Invalid(t *testing.T) {
	req := dto.CreateAccountRequest{
		Name:  "",
		Email: "invalid-email",
	}
	acc, err := entity.ToEntity(req)
	assert.Error(t, err)
	assert.Empty(t, acc.Name)
	assert.Empty(t, acc.Email)
}

func TestToEntityUpdate_Valid(t *testing.T) {
	req := dto.UpdateAccountRequest{
		Name:  "Bob",
		Email: "bob@email.com",
	}
	acc, err := entity.ToEntityUpdate(req)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", acc.Name)
	assert.Equal(t, "bob@email.com", acc.Email)
}

func TestToEntityUpdate_Invalid(t *testing.T) {
	req := dto.UpdateAccountRequest{
		Name:  "",
		Email: "not-an-email",
	}
	acc, err := entity.ToEntityUpdate(req)
	assert.Error(t, err)
	assert.Empty(t, acc.Name)
	assert.Empty(t, acc.Email)
}

func TestToEntityFilter(t *testing.T) {
	filter := dto.AccountFilter{Name: "Alice", Email: "alice@email.com"}
	entityFilter := entity.ToEntityFilter(filter)
	assert.Equal(t, "Alice", entityFilter.Name)
	assert.Equal(t, "alice@email.com", entityFilter.Email)
}

func TestToListDTO(t *testing.T) {
	accounts := []entity.Account{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
	}
	dtos := entity.ToListDTO(accounts)
	assert.Len(t, dtos, 2)
	assert.Equal(t, "1", dtos[0].ID)
	assert.Equal(t, "Alice", dtos[0].Name)
	assert.Equal(t, "2", dtos[1].ID)
	assert.Equal(t, "Bob", dtos[1].Name)
}

func TestAccount_IsExisting(t *testing.T) {
	acc := entity.Account{ID: "123"}
	assert.True(t, acc.IsExisting())

	accEmpty := entity.Account{}
	assert.False(t, accEmpty.IsExisting())
}

func TestAccount_ZeroValues(t *testing.T) {
	acc := entity.Account{}
	assert.Empty(t, acc.ID)
	assert.Empty(t, acc.Name)
	assert.Empty(t, acc.Email)
	assert.True(t, acc.CreatedAt.IsZero())
	assert.True(t, acc.UpdatedAt.IsZero())
}
