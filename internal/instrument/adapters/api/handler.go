package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/app"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type Instrument interface {
	Create(c echo.Context) error
	FindByID(c echo.Context) error
	GetInstruments(c echo.Context) error
	Update(c echo.Context) error
	DeleteByID(c echo.Context) error
	RegisterRoutes(g *echo.Group)
}

type instrument struct {
	instrumentApp app.Instrument
}

func NewInstrumentHandler(instrumentApp app.Instrument) Instrument {
	return &instrument{
		instrumentApp: instrumentApp,
	}
}

func (h *instrument) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("/:id", h.FindByID)
	g.GET("", h.GetInstruments)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.DeleteByID)
}

// Create godoc
// @Summary      Cria um novo instrumento
// @Description  Cria um novo instrumento financeiro
// @Tags         instruments
// @Accept       json
// @Produce      json
// @Param        instrument  body      dto.CreateInstrumentRequest  true  "Instrument"
// @Success      201    {object}  dto.InstrumentDTO
// @Failure      400    {object}  map[string]string
// @Failure      409    {object}  map[string]string
// @Router       /v1/instruments [post]
func (h *instrument) Create(c echo.Context) error {
	var request dto.CreateInstrumentRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload: "+err.Error())
	}
	if err := request.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "validation failed: "+err.Error())
	}

	createdInstrument, err := h.instrumentApp.Create(c.Request().Context(), request)
	if err != nil {
		if errors.Is(err, ierr.ErrConflict) {
			return c.JSON(http.StatusConflict, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusCreated, createdInstrument)
}

// FindByID godoc
// @Summary      Busca um instrumento por ID
// @Tags         instruments
// @Produce      json
// @Param        id   path      string  true  "Instrument ID"
// @Success      200  {object}  dto.InstrumentDTO
// @Failure      404  {object}  map[string]string
// @Router       /v1/instruments/{id} [get]
func (h *instrument) FindByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "instrument ID cannot be empty")
	}

	instrument, err := h.instrumentApp.FindByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, ierr.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, instrument)
}

// GetInstruments godoc
// @Summary      Lista instrumentos financeiros
// @Tags         instruments
// @Produce      json
// @Param        base_asset   query   string  false  "Base Asset"
// @Param        quote_asset  query   string  false  "Quote Asset"
// @Success      200  {array}   dto.InstrumentListDTO
// @Router       /v1/instruments [get]
func (h *instrument) GetInstruments(c echo.Context) error {
	filter := dto.InstrumentFilter{
		BaseAsset:  c.QueryParam("base_asset"),
		QuoteAsset: c.QueryParam("quote_asset"),
	}

	instruments, err := h.instrumentApp.GetInstruments(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, instruments)
}

// Update godoc
// @Summary      Atualiza um instrumento
// @Tags         instruments
// @Accept       json
// @Produce      json
// @Param        id         path      string                     true  "Instrument ID"
// @Param        instrument body      dto.CreateInstrumentRequest true  "Instrument"
// @Success      200  {object}  dto.InstrumentDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /v1/instruments/{id} [put]
func (h *instrument) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "instrument ID cannot be empty")
	}

	var request dto.CreateInstrumentRequest // Reutilizando o DTO de criação para o update
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	updatedInstrument, err := h.instrumentApp.Update(c.Request().Context(), id, request)
	if err != nil {
		if errors.Is(err, ierr.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		if errors.Is(err, ierr.ErrConflict) {
			return c.JSON(http.StatusConflict, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.JSON(http.StatusOK, updatedInstrument)
}

// DeleteByID godoc
// @Summary      Deleta um instrumento
// @Tags         instruments
// @Produce      json
// @Param        id   path      string  true  "Instrument ID"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /v1/instruments/{id} [delete]
func (h *instrument) DeleteByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "instrument ID cannot be empty")
	}

	if err := h.instrumentApp.DeleteByID(c.Request().Context(), id); err != nil {
		if errors.Is(err, ierr.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return c.NoContent(http.StatusNoContent)
}
