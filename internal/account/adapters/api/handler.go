package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/app"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
)

type Account interface {
	Create(context echo.Context) error
	GetAccounts(context echo.Context) error
	FindByID(context echo.Context) error
	DeleteByID(context echo.Context) error
	Update(context echo.Context) error
	RegisterRoutes(g *echo.Group)
}

type account struct {
	accountApp app.Account
}

func NewAccountHandler(accountApp app.Account) Account {
	return &account{
		accountApp: accountApp,
	}
}

func (h *account) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("/:id", h.FindByID)
	g.GET("", h.GetAccounts)
	g.DELETE("/:id", h.DeleteByID)
	g.PUT("/:id", h.Update)
}

func (h *account) Create(ctx echo.Context) error {
	var request dto.CreateAccountRequest

	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload: "+err.Error())
	}

	if err := request.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "validation failed: "+err.Error())
	}

	account, err := h.accountApp.Create(ctx.Request().Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, ierr.ErrConflict):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		case errors.Is(err, ierr.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			slog.Error("error creating account", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
		}
	}

	return ctx.JSON(http.StatusCreated, account)
}

func (h *account) FindByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "account ID cannot be empty")
	}

	account, err := h.accountApp.FindByID(ctx.Request().Context(), id)
	if err != nil {
		slog.Error("error finding account by ID", "error", err)
		if errors.Is(err, ierr.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, ierr.ErrNotFound) // retun 404 for "not found"
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}

	return ctx.JSON(http.StatusOK, account)
}

func (h *account) GetAccounts(ctx echo.Context) error {
	email := ctx.QueryParam("email")
	name := ctx.QueryParam("name")

	filter := dto.AccountFilter{
		Email: email,
		Name:  name,
	}

	accounts, err := h.accountApp.GetAccounts(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, accounts)
}

func (h *account) DeleteByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id cannot be empty")
	}

	if err := h.accountApp.DeleteByID(ctx.Request().Context(), id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *account) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id cannot be empty")
	}

	var request dto.UpdateAccountRequest
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := request.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updatedAccount, err := h.accountApp.Update(ctx.Request().Context(), id, request)
	if err != nil {
		switch {
		case errors.Is(err, ierr.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			slog.Error("error updating account", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
		}
	}

	return ctx.JSON(http.StatusOK, updatedAccount)
}
