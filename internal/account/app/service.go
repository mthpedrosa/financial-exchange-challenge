package app

import (
	"context"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type Account interface {
	Create(ctx context.Context, request dto.CreateAccountRequest) (dto.CreateAccountResponse, error)
	FindByID(ctx context.Context, email string) (dto.AccountDTO, error)
	GetAccounts(ctx context.Context, filters dto.AccountFilter) ([]dto.AccountListDTO, error)
	DeleteByID(ctx context.Context, id string) error
	Update(ctx context.Context, id string, request dto.UpdateAccountRequest) (dto.AccountDTO, error)
}

type account struct {
	accountPort port.AccountRepository
}

func NewAccountApp(accountPort port.AccountRepository) Account {
	return &account{
		accountPort: accountPort,
	}
}

func (a *account) Create(ctx context.Context, request dto.CreateAccountRequest) (dto.CreateAccountResponse, error) {
	entityAccount, err := entity.ToEntity(request)
	if err != nil {
		return dto.CreateAccountResponse{}, err
	}

	existingAccounts, err := a.accountPort.GetAccounts(ctx, entity.AccountFilter{Email: entityAccount.Email})
	if err != nil {
		return dto.CreateAccountResponse{}, err
	}

	if len(existingAccounts) > 0 {
		return dto.CreateAccountResponse{}, ierr.ErrConflict
	}

	entityAccount.ID, err = a.accountPort.Create(ctx, *entityAccount)
	if err != nil {
		return dto.CreateAccountResponse{}, err
	}

	return dto.CreateAccountResponse{ID: entityAccount.ID}, nil
}

func (a *account) FindByID(ctx context.Context, id string) (dto.AccountDTO, error) {
	account, err := a.accountPort.FindByID(ctx, id)
	return account.ToDTO(), err
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

func (a *account) Update(ctx context.Context, id string, request dto.UpdateAccountRequest) (dto.AccountDTO, error) {
	entityAccount, err := entity.ToEntityUpdate(request)
	if err != nil {
		return dto.AccountDTO{}, err
	}

	account, err := a.accountPort.FindByID(ctx, id)
	if err != nil {
		return dto.AccountDTO{}, err
	}

	if !account.IsExisting() {
		return dto.AccountDTO{}, nil
	}

	entityAccount.ID = id

	err = a.accountPort.Update(ctx, *entityAccount)
	if err != nil {
		return dto.AccountDTO{}, err
	}

	return entityAccount.ToDTO(), nil
}
