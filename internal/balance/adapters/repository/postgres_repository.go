package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type balanceRepository struct {
	db *pgxpool.Pool
}

func NewBalanceRepository(db *pgxpool.Pool) port.BalanceRepository {
	return &balanceRepository{db: db}
}

// Create inserts a new balance. It checks for duplicates based on account_id and asset.
func (r *balanceRepository) Create(ctx context.Context, balance entity.Balance) (string, error) {
	// Check for duplicate
	queryCheck := `SELECT id FROM balances WHERE account_id = $1 AND asset = $2`
	var existingID string
	err := r.db.QueryRow(ctx, queryCheck, balance.AccountID, balance.Asset).Scan(&existingID)
	if err == nil {
		return "", ierr.ErrConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	m := ToModel(balance)
	query := `INSERT INTO balances (account_id, asset, amount, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	var id string
	err = r.db.QueryRow(ctx, query, m.AccountID, m.Asset, m.Amount).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// FindByID returns a balance by its ID.
func (r *balanceRepository) FindByID(ctx context.Context, id string) (entity.Balance, error) {
	query := `SELECT id, account_id, asset, amount, created_at, updated_at FROM balances WHERE id = $1`
	var m BalanceModel

	var amountStr string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.AccountID,
		&m.Asset,
		&amountStr,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Balance{}, ierr.ErrNotFound
		}
		return entity.Balance{}, err
	}
	m.Amount, _ = new(big.Float).SetString(amountStr)
	return m.ToEntity(), nil
}

// FindByAccountAndAsset returns a balance for a given account and asset.
func (r *balanceRepository) FindByAccountAndAsset(ctx context.Context, accountID, asset string) (entity.Balance, error) {
	query := `SELECT id, account_id, asset, amount, created_at, updated_at FROM balances WHERE account_id = $1 AND asset = $2`
	var m BalanceModel
	var amountStr string
	err := r.db.QueryRow(ctx, query, accountID, asset).Scan(
		&m.ID,
		&m.AccountID,
		&m.Asset,
		&amountStr,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Balance{}, ierr.ErrNotFound
		}
		return entity.Balance{}, err
	}
	m.Amount, _ = new(big.Float).SetString(amountStr)
	return m.ToEntity(), nil
}

// Update modifies an existing balance.
func (r *balanceRepository) Update(ctx context.Context, balance entity.Balance) error {
	m := ToModel(balance)
	query := `UPDATE balances SET amount = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(ctx, query, m.Amount, m.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no balance found with id: %s", m.ID)
	}
	return nil
}

// DeleteByID removes a balance by its ID.
func (r *balanceRepository) DeleteByID(ctx context.Context, id string) error {
	query := `DELETE FROM balances WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no balance found with id: %s", id)
	}
	return nil
}

// GetAllByAccountID returns all balances for a given account ID.
func (r *balanceRepository) GetAllByAccountID(ctx context.Context, accountID string) ([]entity.Balance, error) {
	query := `SELECT id, account_id, asset, amount, created_at, updated_at FROM balances WHERE account_id = $1`
	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []entity.Balance
	for rows.Next() {
		var m BalanceModel
		var amountStr string
		if err := rows.Scan(&m.ID, &m.AccountID, &m.Asset, &amountStr, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		m.Amount, _ = new(big.Float).SetString(amountStr)
		balances = append(balances, m.ToEntity())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return balances, nil
}
