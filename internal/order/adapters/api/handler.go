package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/app"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/order/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type Order interface {
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
	FindByID(ctx echo.Context) error
	GetOrders(ctx echo.Context) error
	CancelByID(ctx echo.Context) error
	FindByInstrument(ctx echo.Context) error
	RegisterRoutes(g *echo.Group)
}

type order struct {
	orderApp app.Order
}

func NewOrderHandler(orderApp app.Order) Order {
	return &order{
		orderApp: orderApp,
	}
}

func (h *order) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("/:id", h.FindByID)
	g.GET("", h.GetOrders)
	g.PUT("/:id", h.Update)
	g.POST("/:id/cancel", h.CancelByID)
	g.GET("/instrument/:instrument_id", h.FindByInstrument)
}

// Create godoc
// @Summary      Cria uma nova ordem
// @Description  Cria uma nova ordem e envia para a fila
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body      dto.CreateOrderRequest  true  "Order"
// @Success      201    {object}  dto.OrderDTO
// @Failure      400    {object}  map[string]string
// @Failure      409    {object}  map[string]string
// @Router       /v1/orders [post]
func (h *order) Create(ctx echo.Context) error {
	var request dto.CreateOrderRequest

	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload: "+err.Error())
	}

	if err := request.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "validation failed: "+err.Error())
	}

	order, err := h.orderApp.Create(ctx.Request().Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, ierr.ErrConflict):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.Is(err, ierr.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			slog.Error("error creating order", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
		}
	}

	return ctx.JSON(http.StatusCreated, order)
}

// Update godoc
// @Summary      Atualiza uma ordem
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id     path      string                 true  "Order ID"
// @Param        order  body      dto.CreateOrderRequest true  "Order"
// @Success      204    "No Content"
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /v1/orders/{id} [put]
func (h *order) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order ID cannot be empty")
	}

	var request dto.CreateOrderRequest
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := request.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Aqui você pode adaptar para aceitar um DTO específico de update, se desejar.
	if err := h.orderApp.Update(ctx.Request().Context(), request); err != nil {
		switch {
		case errors.Is(err, ierr.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			slog.Error("error updating order", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
		}
	}

	return ctx.NoContent(http.StatusNoContent)
}

// FindByID godoc
// @Summary      Busca uma ordem por ID
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  dto.OrderDTO
// @Failure      404  {object}  map[string]string
// @Router       /v1/orders/{id} [get]
func (h *order) FindByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order ID cannot be empty")
	}

	order, err := h.orderApp.FindByID(ctx.Request().Context(), id)
	if err != nil {
		slog.Error("error finding order by ID", "error", err)
		if errors.Is(err, ierr.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, ierr.ErrNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}

	return ctx.JSON(http.StatusOK, order)
}

// GetOrders godoc
// @Summary      Lista todas as ordens
// @Tags         orders
// @Produce      json
// @Success      200  {array}   dto.OrderDTO
// @Router       /v1/orders [get]
func (h *order) GetOrders(ctx echo.Context) error {
	orders, err := h.orderApp.GetAll(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}
	return ctx.JSON(http.StatusOK, orders)
}

// CancelByID godoc
// @Summary      Cancela uma ordem
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /v1/orders/{id}/cancel [post]
func (h *order) CancelByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "order ID cannot be empty")
	}

	if err := h.orderApp.CancelByID(ctx.Request().Context(), id); err != nil {
		switch {
		case errors.Is(err, ierr.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			slog.Error("error cancelling order", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
		}
	}

	return ctx.NoContent(http.StatusNoContent)
}

// FindByInstrument godoc
// @Summary      Busca as orders por um Intrument
// @Tags         orders
// @Produce      json
// @Param        instrument_id   path   string  true  "Instrument ID"
// @Param        asset        path   string  true  "Asset"
// @Success      200  {array}   dto.OrderDTO
// @Failure      404  {message}  map[string]string "instrumentID required"
// @Router       /v1/orders/instrument/{instrumentID} [get]
func (h *order) FindByInstrument(ctx echo.Context) error {
	instrumentID := ctx.Param("instrument_id")
	if instrumentID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "instrumentID are required")
	}

	balance, err := h.orderApp.FindByInstrument(ctx.Request().Context(), instrumentID)
	if err != nil {
		slog.Error("error finding orders by instrument", "error", err)
		if errors.Is(err, ierr.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, ierr.ErrNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}

	return ctx.JSON(http.StatusOK, balance)
}
