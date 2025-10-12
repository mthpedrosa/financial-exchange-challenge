package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/port"
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
	// Supondo que ToModel converte a entidade para um modelo de DB
	// e que o modelo tem os mesmos campos Name e Email.
	model := ToModel(account)

	query := `INSERT INTO accounts (name, email, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	var id string

	// Usando os campos do 'model' (ou do 'account' se não houver transformação)
	err := r.db.QueryRow(ctx, query, model.Name, model.Email).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *account) FindByID(ctx context.Context, email string) (entity.Account, error) {
	var model AccountModel

	query := `SELECT id, name, email, created_at, updated_at FROM accounts WHERE id=$1`
	err := r.db.QueryRow(ctx, query, email).Scan(&model.ID, &model.Name, &model.Email, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return entity.Account{}, err
	}

	return ToEntity(model), nil
}

func (r *account) GetAccounts(ctx context.Context, filters entity.AccountFilter) ([]entity.Account, error) {
	query := `SELECT * FROM accounts`
	rows, err := r.db.Query(ctx, query)
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
