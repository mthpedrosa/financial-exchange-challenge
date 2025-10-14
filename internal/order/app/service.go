package app

import (
	"context"
	"errors"
	"math/big"

	accountPort "github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/port"
	balancePort "github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/port"
	instrumentPort "github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/port"
)

type Order interface {
	Create(ctx context.Context, req dto.CreateOrderRequest) (dto.CreateOrderResponse, error)
	FindByID(ctx context.Context, id string) (dto.OrderDTO, error)
	GetAll(ctx context.Context) ([]dto.OrderDTO, error)
	Update(ctx context.Context, req dto.CreateOrderRequest) error
	CancelByID(ctx context.Context, id string) error
	FindByInstrument(ctx context.Context, id string) ([]dto.OrderDTO, error)
}

type orderApp struct {
	orderRepo      port.OrderRepository
	accountRepo    accountPort.AccountRepository
	instrumentRepo instrumentPort.InstrumentRepository
	balanceRepo    balancePort.BalanceRepository
	orderQueue     port.OrderQueue
}

func NewOrderApp(
	orderRepo port.OrderRepository,
	accountRepo accountPort.AccountRepository,
	instrumentRepo instrumentPort.InstrumentRepository,
	balanceRepo balancePort.BalanceRepository,
	orderQueue port.OrderQueue,
) Order {
	return &orderApp{
		orderRepo:      orderRepo,
		accountRepo:    accountRepo,
		instrumentRepo: instrumentRepo,
		balanceRepo:    balanceRepo,
		orderQueue:     orderQueue,
	}
}

func (a *orderApp) Create(ctx context.Context, req dto.CreateOrderRequest) (dto.CreateOrderResponse, error) {
	orderEntity, err := entity.ToEntity(req)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	// check if account exists
	_, err = a.accountRepo.FindByID(ctx, orderEntity.AccountID)
	if err != nil {
		return dto.CreateOrderResponse{}, errors.New("account not found")
	}

	// check if instrument exists
	instrument, err := a.instrumentRepo.FindByID(ctx, orderEntity.InstrumentID)
	if err != nil {
		return dto.CreateOrderResponse{}, errors.New("instrument not found")
	}

	// check if account has sufficient balance for the order
	var asset string
	var requiredAmount *big.Float

	if orderEntity.Type == entity.OrderTypeBuy {
		asset = instrument.QuoteAsset
		requiredAmount = new(big.Float).Mul(orderEntity.Price, orderEntity.Quantity)
	} else {
		asset = instrument.BaseAsset
		requiredAmount = orderEntity.Quantity
	}

	balance, err := a.balanceRepo.FindByAccountAndAsset(ctx, orderEntity.AccountID, asset)
	if err != nil {
		return dto.CreateOrderResponse{}, errors.New("balance not found for required asset")
	}

	if balance.Amount.Cmp(requiredAmount) < 0 {
		return dto.CreateOrderResponse{}, errors.New("insufficient balance")
	}

	id, err := a.orderRepo.Create(ctx, *orderEntity)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	orderEntity.ID = id

	// send order to queue for processing
	if err := a.orderQueue.PublishOrder(ctx, *orderEntity); err != nil {
		return dto.CreateOrderResponse{}, err
	}

	orderEntity.ID = id
	return dto.CreateOrderResponse{ID: id}, nil
}

// FindByID finds an order by its ID.
func (a *orderApp) FindByID(ctx context.Context, id string) (dto.OrderDTO, error) {
	order, err := a.orderRepo.FindByID(ctx, id)
	if err != nil {
		return dto.OrderDTO{}, err
	}
	return order.ToDTO(), nil
}

func (a *orderApp) FindByInstrument(ctx context.Context, id string) ([]dto.OrderDTO, error) {
	orders, err := a.orderRepo.FindByInstrumentID(ctx, id)
	if err != nil {
		return []dto.OrderDTO{}, err
	}
	return entity.ToListDTO(orders), nil
}

// GetAll returns all orders.
func (a *orderApp) GetAll(ctx context.Context) ([]dto.OrderDTO, error) {
	orders, err := a.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return entity.ToListDTO(orders), nil
}

// Update updates an order.
func (a *orderApp) Update(ctx context.Context, req dto.CreateOrderRequest) error {
	orderEntity, err := entity.ToEntity(req)
	if err != nil {
		return err
	}
	return a.orderRepo.Update(ctx, *orderEntity)
}

// DeleteByID deletes an order by its ID.
func (a *orderApp) CancelByID(ctx context.Context, id string) error {
	order, err := a.orderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	order.Status = entity.OrderStatusCancelled
	return a.orderRepo.Update(ctx, order)
}
