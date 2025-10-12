package app

import (
	"context"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/port"
)

type Account interface {
	Create(ctx context.Context, request dto.CreateAccountRequest) (dto.CreateAcountResponse, error)
	FindByID(ctx context.Context, email string) (entity.Account, error)
	GetAccounts(ctx context.Context, filters dto.AccountFilter) ([]dto.AccountListDTO, error)
	DeleteByID(ctx context.Context, id string) error
	Update(ctx context.Context, id string, request dto.UpdateAccountRequest) (entity.Account, error)
}

type account struct {
	accountPort port.AccountRepository
}

func NewAccountApp(accountPort port.AccountRepository) Account {
	return &account{
		accountPort: accountPort,
	}
}

func (a *account) Create(ctx context.Context, request dto.CreateAccountRequest) (dto.CreateAcountResponse, error) {
	entityAccount, err := entity.ToEntity(request)
	if err != nil {
		return dto.CreateAcountResponse{}, err
	}

	account, err := a.accountPort.GetAccounts(ctx, entity.AccountFilter{Email: entityAccount.Email})
	if err != nil {
		return dto.CreateAcountResponse{}, err
	}

	if len(account) == 0 && account[0].IsExisting() {
		return dto.CreateAcountResponse{}, nil
	}

	entityAccount.ID, err = a.accountPort.Create(ctx, *entityAccount)
	if err != nil {
		return dto.CreateAcountResponse{}, err

	}

	return dto.CreateAcountResponse{ID: entityAccount.ID}, nil
}

func (a *account) FindByID(ctx context.Context, id string) (entity.Account, error) {
	return a.accountPort.FindByID(ctx, id)
}

func (a *account) GetAccounts(ctx context.Context, filters dto.AccountFilter) ([]dto.AccountListDTO, error) {
	accounts, err := a.accountPort.GetAccounts(ctx, entity.ToEntityFilter(filters))
	if err != nil {
		return nil, err
	}

	return entity.ToListDTO(accounts), nil
}

func (a *account) DeleteByID(ctx context.Context, id string) error {
	return a.accountPort.DeleteByID(ctx, id)
}

func (a *account) Update(ctx context.Context, id string, request dto.UpdateAccountRequest) (entity.Account, error) {
	entityAccount, err := entity.ToEntityUpdate(request)
	if err != nil {
		return entity.Account{}, err
	}

	account, err := a.accountPort.FindByID(ctx, id)
	if err != nil {
		return entity.Account{}, err
	}

	if !account.IsExisting() {
		return entity.Account{}, nil
	}

	entityAccount.ID = id

	err = a.accountPort.Update(ctx, *entityAccount)
	if err != nil {
		return entity.Account{}, err
	}

	return *entityAccount, nil
}
