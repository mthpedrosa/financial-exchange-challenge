package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) port.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order entity.Order) (string, error) {
	query := `INSERT INTO orders (account_id, instrument_id, type, status, price, quantity, remaining_quantity, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id`
	var id string
	err := r.db.QueryRow(ctx, query,
		order.AccountID,
		order.InstrumentID,
		string(order.Type),
		string(order.Status),
		order.Price.Text('f', 10),
		order.Quantity.Text('f', 18),
		order.RemainingQuantity.Text('f', 18),
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (entity.Order, error) {
	query := `SELECT id, account_id, instrument_id, type, status, price, quantity, remaining_quantity, created_at, updated_at FROM orders WHERE id = $1`
	var o entity.Order
	var priceStr, quantityStr, remainingStr string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID,
		&o.AccountID,
		&o.InstrumentID,
		&o.Type,
		&o.Status,
		&priceStr,
		&quantityStr,
		&remainingStr,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, ierr.ErrNotFound
		}
		return entity.Order{}, err
	}
	o.Price, _ = new(big.Float).SetString(priceStr)
	o.Quantity, _ = new(big.Float).SetString(quantityStr)
	o.RemainingQuantity, _ = new(big.Float).SetString(remainingStr)
	return o, nil
}

func (r *orderRepository) GetAll(ctx context.Context) ([]entity.Order, error) {
	query := `SELECT id, account_id, instrument_id, type, status, price, quantity, remaining_quantity, created_at, updated_at FROM orders`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		var priceStr, quantityStr, remainingStr string
		if err := rows.Scan(
			&o.ID,
			&o.AccountID,
			&o.InstrumentID,
			&o.Type,
			&o.Status,
			&priceStr,
			&quantityStr,
			&remainingStr,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		o.Price, _ = new(big.Float).SetString(priceStr)
		o.Quantity, _ = new(big.Float).SetString(quantityStr)
		o.RemainingQuantity, _ = new(big.Float).SetString(remainingStr)
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Update(ctx context.Context, order entity.Order) error {
	query := `UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`
	result, err := r.db.Exec(ctx, query, string(order.Status), order.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no order found with id: %s", order.ID)
	}
	return nil
}

func (r *orderRepository) FindByInstrumentID(ctx context.Context, id string) ([]entity.Order, error) {
	query := `SELECT * FROM orders WHERE instrument_id = $1`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		var priceStr, quantityStr, remainingStr string
		if err := rows.Scan(
			&o.ID,
			&o.AccountID,
			&o.InstrumentID,
			&o.Type,
			&o.Status,
			&priceStr,
			&quantityStr,
			&remainingStr,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		o.Price, _ = new(big.Float).SetString(priceStr)
		o.Quantity, _ = new(big.Float).SetString(quantityStr)
		o.RemainingQuantity, _ = new(big.Float).SetString(remainingStr)
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
