package entity

import (
	"time"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/dto"
)

type Account struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountFilter struct {
	Name  string
	Email string
}

func (a Account) ToDTO() dto.AccountDTO {
	return dto.AccountDTO{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func ToEntityFilter(dto dto.AccountFilter) AccountFilter {
	return AccountFilter{
		Name:  dto.Name,
		Email: dto.Email,
	}
}

func ToListDTO(accounts []Account) []dto.AccountListDTO {
	dtos := make([]dto.AccountListDTO, len(accounts))
	for i, a := range accounts {
		dtos[i] = dto.AccountListDTO{
			ID:   a.ID,
			Name: a.Name,
		}
	}
	return dtos
}

func ToEntity(request dto.CreateAccountRequest) (*Account, error) {
	if err := request.Validate(); err != nil {
		return &Account{}, err
	}

	return &Account{
		Name:  request.Name,
		Email: request.Email,
	}, nil
}

func ToEntityUpdate(request dto.UpdateAccountRequest) (*Account, error) {
	if err := request.Validate(); err != nil {
		return &Account{}, err
	}

	return &Account{
		Name:  request.Name,
		Email: request.Email,
	}, nil
}

func (a *Account) IsExisting() bool {
	return a.ID != ""
}
