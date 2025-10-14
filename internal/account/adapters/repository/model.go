package repository

import (
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
)

type AccountModel struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToModel(account entity.Account) *AccountModel {
	return &AccountModel{
		ID:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}

func ToEntity(model AccountModel) entity.Account {
	return entity.Account{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
