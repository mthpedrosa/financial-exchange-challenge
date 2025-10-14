package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/app"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/balance/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type Balance interface {
	Create(context echo.Context) error
	FindByID(context echo.Context) error
	FindByAccountAndAsset(context echo.Context) error
	GetAllByAccountID(context echo.Context) error
	Update(context echo.Context) error
	DeleteByID(context echo.Context) error
	RegisterRoutes(g *echo.Group)
}

type balance struct {
	balanceApp app.Balance
}

func NewBalanceHandler(balanceApp app.Balance) Balance {
	return &balance{
		balanceApp: balanceApp,
	}
}

func (h *balance) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("/:id", h.FindByID)
	g.GET("/account/:account_id", h.GetAllByAccountID)
	g.GET("/account/:account_id/asset/:asset", h.FindByAccountAndAsset)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.DeleteByID)
}

// Create godoc
// @Summary      Cria um novo balance
// @Description  Cria um novo balance para uma conta e asset
// @Tags         balances
// @Accept       json
// @Produce      json
// @Param        balance  body      dto.CreateBalanceRequest  true  "Balance"
// @Success      201    {object}  dto.BalanceListDTO
// @Failure      400    {object}  map[string]string
// @Failure      409    {object}  map[string]string
// @Router       /v1/balances [post]
func (h *balance) Create(ctx echo.Context) error {
	var request dto.CreateBalanceRequest

	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload: "+err.Error())
	}

	if err := request.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "validation failed: "+err.Error())
	}

	resp, err := h.balanceApp.Create(ctx.Request().Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, ierr.ErrConflict):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			slog.Error("error creating balance", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return ctx.JSON(http.StatusCreated, resp)
}

// FindByID godoc
// @Summary      Busca um balance por ID
// @Tags         balances
// @Produce      json
// @Param        id   path      string  true  "Balance ID"
// @Success      200  {object}  dto.BalanceListDTO
// @Failure      404  {object}  map[string]string
// @Router       /v1/balances/{id} [get]
func (h *balance) FindByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "balance ID cannot be empty")
	}

	balance, err := h.balanceApp.FindByID(ctx.Request().Context(), id)
	if err != nil {
		slog.Error("error finding balance by ID", "error", err)
		if errors.Is(err, ierr.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, ierr.ErrNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}

	return ctx.JSON(http.StatusOK, balance)
}

// FindByAccountAndAsset godoc
// @Summary      Busca um balance por account e asset
// @Tags         balances
// @Produce      json
// @Param        account_id   path   string  true  "Account ID"
// @Param        asset        path   string  true  "Asset"
// @Success      200  {object}  dto.BalanceListDTO
// @Failure      404  {object}  map[string]string
// @Router       /v1/balances/account/{account_id}/asset/{asset} [get]
func (h *balance) FindByAccountAndAsset(ctx echo.Context) error {
	accountID := ctx.Param("account_id")
	asset := ctx.Param("asset")
	if accountID == "" || asset == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "account_id and asset are required")
	}

	balance, err := h.balanceApp.FindByAccountAndAsset(ctx.Request().Context(), accountID, asset)
	if err != nil {
		slog.Error("error finding balance by account and asset", "error", err)
		if errors.Is(err, ierr.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, ierr.ErrNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}

	return ctx.JSON(http.StatusOK, balance)
}

// GetAllByAccountID godoc
// @Summary      Lista todos os balances de uma conta
// @Tags         balances
// @Produce      json
// @Param        account_id   path   string  true  "Account ID"
// @Success      200  {array}  dto.BalanceListDTO
// @Router       /v1/balances/account/{account_id} [get]
func (h *balance) GetAllByAccountID(ctx echo.Context) error {
	accountID := ctx.Param("account_id")
	if accountID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "account_id is required")
	}

	balances, err := h.balanceApp.GetAllByAccountID(ctx.Request().Context(), accountID)
	if err != nil {
		slog.Error("error getting balances by account_id", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}

	return ctx.JSON(http.StatusOK, balances)
}

// Update godoc
// @Summary      Atualiza um balance
// @Tags         balances
// @Accept       json
// @Produce      json
// @Param        id      path      string                  true  "Balance ID"
// @Param        balance body      dto.UpdateBalanceRequest true  "Balance"
// @Success      200  {object}  dto.BalanceListDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /v1/balances/{id} [put]
func (h *balance) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id cannot be empty")
	}

	var request dto.UpdateBalanceRequest
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := request.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updatedBalance, err := h.balanceApp.Update(ctx.Request().Context(), id, request)
	if err != nil {
		switch {
		case errors.Is(err, ierr.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			slog.Error("error updating balance", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
		}
	}

	return ctx.JSON(http.StatusOK, updatedBalance)
}

// DeleteByID godoc
// @Summary      Deleta um balance
// @Tags         balances
// @Produce      json
// @Param        id   path      string  true  "Balance ID"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /v1/balances/{id} [delete]
func (h *balance) DeleteByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id cannot be empty")
	}

	if err := h.balanceApp.DeleteByID(ctx.Request().Context(), id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
