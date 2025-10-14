package port

import (
	"context"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order entity.Order) (string, error)
	FindByID(ctx context.Context, id string) (entity.Order, error)
	GetAll(ctx context.Context) ([]entity.Order, error)
	Update(ctx context.Context, order entity.Order) error
}

type OrderQueue interface {
	PublishOrder(ctx context.Context, order entity.Order) error
}
