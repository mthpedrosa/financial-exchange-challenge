package app

import (
	"context"
	"errors"

	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/entity"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/port"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type Instrument interface {
	Create(ctx context.Context, request dto.CreateInstrumentRequest) (string, error)
	FindByID(ctx context.Context, id string) (*entity.Instrument, error)
	GetInstruments(ctx context.Context, filter dto.InstrumentFilter) ([]*entity.Instrument, error)
	Update(ctx context.Context, id string, request dto.CreateInstrumentRequest) (*entity.Instrument, error)
	DeleteByID(ctx context.Context, id string) error
}

type instrument struct {
	instrumentPort port.InstrumentRepository
}

func NewInstrumentApp(instrumentPort port.InstrumentRepository) Instrument {
	return &instrument{
		instrumentPort: instrumentPort,
	}
}

func (i *instrument) Create(ctx context.Context, request dto.CreateInstrumentRequest) (string, error) {
	instrumentEntity, err := entity.ToEntity(request)
	if err != nil {
		return "", ierr.ErrInvalidInput
	}

	// check duplicate
	_, err = i.instrumentPort.FindByAssets(ctx, instrumentEntity.BaseAsset, instrumentEntity.QuoteAsset)
	if err == nil {
		return "", ierr.ErrConflict
	}

	if !errors.Is(err, ierr.ErrNotFound) {
		return "", err
	}

	// create
	createdInstrumentID, createErr := i.instrumentPort.Create(ctx, instrumentEntity)
	if createErr != nil {
		return "", createErr
	}
	return createdInstrumentID, nil
}

func (i *instrument) FindByID(ctx context.Context, id string) (*entity.Instrument, error) {
	return i.instrumentPort.FindByID(ctx, id)
}

func (i *instrument) GetInstruments(ctx context.Context, filter dto.InstrumentFilter) ([]*entity.Instrument, error) {
	instruments, err := i.instrumentPort.FindAll(ctx, entity.ToEntityFilter(filter))
	if err != nil {
		return nil, err
	}
	return instruments, nil
}

func (i *instrument) DeleteByID(ctx context.Context, id string) error {
	return i.instrumentPort.DeleteByID(ctx, id)
}

func (i *instrument) Update(ctx context.Context, id string, request dto.CreateInstrumentRequest) (*entity.Instrument, error) {
	// check duplicate
	instrumentToUpdate, err := i.instrumentPort.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// update fields
	instrumentToUpdate.BaseAsset = request.BaseAsset
	instrumentToUpdate.QuoteAsset = request.QuoteAsset

	if err := i.instrumentPort.Update(ctx, instrumentToUpdate); err != nil {
		return nil, err
	}

	return instrumentToUpdate, nil
}
