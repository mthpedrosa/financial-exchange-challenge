package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type instrument struct {
	db *pgxpool.Pool
}

func NewInstrumentRepository(db *pgxpool.Pool) port.InstrumentRepository {
	return &instrument{
		db,
	}
}

func (r *instrument) Create(ctx context.Context, instrument *entity.Instrument) (string, error) {
	model := ToModel(instrument)
	fmt.Print(instrument, "chequei até aqui repositoru ")

	query := `INSERT INTO instruments (base_asset, quote_asset, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id`
	var id string

	err := r.db.QueryRow(ctx, query, model.BaseAsset, model.QuoteAsset).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *instrument) Update(ctx context.Context, instrument *entity.Instrument) error {
	model := ToModel(instrument)

	query := `UPDATE instruments SET base_asset=$1, quote_asset=$2, updated_at=NOW() WHERE id=$3`
	result, err := r.db.Exec(ctx, query, model.BaseAsset, model.QuoteAsset, model.ID)
	if err != nil {
		return err
	}

	// check if any row was affected
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no account found with id: %d", model.ID)
	}

	return nil
}

func (r *instrument) FindByID(ctx context.Context, id string) (*entity.Instrument, error) {
	query := `SELECT id, base_asset, quote_asset, created_at, updated_at FROM instruments WHERE id = $1`

	var acc entity.Instrument
	err := r.db.QueryRow(ctx, query, id).Scan(
		&acc.ID,
		&acc.BaseAsset,
		&acc.QuoteAsset,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &entity.Instrument{}, ierr.ErrNotFound
		}
		return &entity.Instrument{}, err
	}

	return &acc, nil
}

func (r *instrument) FindAll(ctx context.Context, filter *entity.InstrumentFilter) ([]*entity.Instrument, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT id, base_asset, quote_asset, created_at, updated_at FROM instruments WHERE 1=1")

	args := []interface{}{}
	argId := 1

	// dynamic filters
	if filter.BaseAsset != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND base_asset = $%d", argId))
		args = append(args, filter.BaseAsset)
		argId++
	}
	if filter.QuoteAsset != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND quote_asset = $%d", argId))
		args = append(args, filter.QuoteAsset)
		argId++
	}

	query := queryBuilder.String()
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instruments []*entity.Instrument
	for rows.Next() {
		var model InstrumentModel
		if err := rows.Scan(
			&model.ID,
			&model.BaseAsset,
			&model.QuoteAsset,
			&model.CreatedAt,
			&model.UpdatedAt,
		); err != nil {
			return nil, err
		}
		instruments = append(instruments, ToEntity(&model))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return instruments, nil
}

func (r *instrument) DeleteByID(ctx context.Context, id string) error {
	query := `DELETE FROM instruments WHERE id=$1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	// check if any row was affected
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no account found with id: %d", id)
	}

	return nil
}

func (r *instrument) FindByAssets(ctx context.Context, baseAsset, quoteAsset string) (*entity.Instrument, error) {
	query := `
        SELECT id, base_asset, quote_asset, created_at, updated_at 
        FROM instruments 
        WHERE base_asset = $1 AND quote_asset = $2`

	var model InstrumentModel
	err := r.db.QueryRow(ctx, query, baseAsset, quoteAsset).Scan(
		&model.ID,
		&model.BaseAsset,
		&model.QuoteAsset,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ierr.ErrNotFound
		}
		return nil, err
	}

	return ToEntity(&model), nil
}
