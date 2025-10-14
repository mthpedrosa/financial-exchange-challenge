package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type account struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) port.AccountRepository {
	return &account{
		db,
	}
}

func (r *account) Create(ctx context.Context, account entity.Account) (string, error) {
	model := ToModel(account)

	query := `INSERT INTO accounts (name, email, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	var id string

	err := r.db.QueryRow(ctx, query, model.Name, model.Email).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *account) Update(ctx context.Context, account entity.Account) error {
	model := ToModel(account)

	query := `UPDATE accounts SET name=$1, email=$2, updated_at=NOW() WHERE id=$3`
	result, err := r.db.Exec(ctx, query, model.Name, model.Email, model.ID)
	if err != nil {
		return err
	}

	// check if any row was affected
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no account found with id: %s", model.ID)
	}

	return nil
}

func (r *account) FindByID(ctx context.Context, id string) (entity.Account, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM accounts WHERE id = $1`

	var acc entity.Account
	err := r.db.QueryRow(ctx, query, id).Scan(
		&acc.ID,
		&acc.Name,
		&acc.Email,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, ierr.ErrNotFound
		}
		return entity.Account{}, err
	}

	return acc, nil
}

func (r *account) GetAccounts(ctx context.Context, filters entity.AccountFilter) ([]entity.Account, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM accounts"
	var args []interface{}
	var conditions []string
	argID := 1

	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argID))
		args = append(args, "%"+filters.Name+"%")
		argID++
	}

	if filters.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", argID))
		args = append(args, filters.Email)
		argID++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []entity.Account
	for rows.Next() {
		var model AccountModel
		if err := rows.Scan(&model.ID, &model.Name, &model.Email, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, ToEntity(model))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *account) DeleteByID(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id=$1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	// check if any row was affected
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no account found with id: %s", id)
	}

	return nil
}
