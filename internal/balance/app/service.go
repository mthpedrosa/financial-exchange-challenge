package app

import (
	"context"
	"errors"
	"fmt"

	account "github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type Balance interface {
	Create(ctx context.Context, request dto.CreateBalanceRequest) (dto.CreateBalanceResponse, error)
	FindByID(ctx context.Context, id string) (entity.Balance, error)
	FindByAccountAndAsset(ctx context.Context, accountID, asset string) (entity.Balance, error)
	GetAllByAccountID(ctx context.Context, accountID string) ([]dto.BalanceListDTO, error)
	Update(ctx context.Context, id string, request dto.UpdateBalanceRequest) (entity.Balance, error)
	DeleteByID(ctx context.Context, id string) error
}

type balance struct {
	balancePort port.BalanceRepository
	accountPort account.AccountRepository
}

func NewBalanceApp(balancePort port.BalanceRepository, accountPort account.AccountRepository) Balance {
	return &balance{
		balancePort: balancePort,
		accountPort: accountPort,
	}
}

func (b *balance) Create(ctx context.Context, request dto.CreateBalanceRequest) (dto.CreateBalanceResponse, error) {
	balanceEntity, err := entity.ToEntity(request)
	if err != nil {
		return dto.CreateBalanceResponse{}, err
	}
	fmt.Println(balanceEntity)

	_, err = b.accountPort.FindByID(ctx, balanceEntity.AccountID)
	if err != nil {
		return dto.CreateBalanceResponse{}, errors.New("account not found")
	}

	// Check for duplicate (by account_id and asset)
	_, err = b.balancePort.FindByAccountAndAsset(ctx, balanceEntity.AccountID, balanceEntity.Asset)
	if err == nil {
		return dto.CreateBalanceResponse{}, ierr.ErrConflict
	}
	if err != nil && err != ierr.ErrNotFound {
		return dto.CreateBalanceResponse{}, err
	}

	id, err := b.balancePort.Create(ctx, *balanceEntity)
	if err != nil {
		return dto.CreateBalanceResponse{}, err
	}

	return dto.CreateBalanceResponse{ID: id}, nil
}

func (b *balance) FindByID(ctx context.Context, id string) (entity.Balance, error) {
	return b.balancePort.FindByID(ctx, id)
}

func (b *balance) FindByAccountAndAsset(ctx context.Context, accountID, asset string) (entity.Balance, error) {
	return b.balancePort.FindByAccountAndAsset(ctx, accountID, asset)
}

func (b *balance) GetAllByAccountID(ctx context.Context, accountID string) ([]dto.BalanceListDTO, error) {
	balances, err := b.balancePort.GetAllByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return entity.ToListDTO(balances), nil
}

func (b *balance) Update(ctx context.Context, id string, request dto.UpdateBalanceRequest) (entity.Balance, error) {
	balanceEntity, err := entity.ToEntityUpdate(request)
	if err != nil {
		return entity.Balance{}, err
	}

	existing, err := b.balancePort.FindByID(ctx, id)
	if err != nil {
		return entity.Balance{}, err
	}
	if existing.ID == "" {
		return entity.Balance{}, ierr.ErrNotFound
	}

	balanceEntity.ID = id
	balanceEntity.AccountID = existing.AccountID // preserve accountID
	balanceEntity.Asset = existing.Asset         // preserve asset

	err = b.balancePort.Update(ctx, *balanceEntity)
	if err != nil {
		return entity.Balance{}, err
	}

	return *balanceEntity, nil
}

func (b *balance) DeleteByID(ctx context.Context, id string) error {
	return b.balancePort.DeleteByID(ctx, id)
}
