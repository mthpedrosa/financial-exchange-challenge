package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mthpedrosa/financial-exchange-challenge/config"
	accountHandler "github.com/mthpedrosa/financial-exchange-challenge/internal/account/adapters/api"
	accountRepo "github.com/mthpedrosa/financial-exchange-challenge/internal/account/adapters/repository"
	accountApp "github.com/mthpedrosa/financial-exchange-challenge/internal/account/app"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/logger"
)

func main() {
	// load config
	cfg := config.LoadConfig()

	// setup logger based on cfg.LogLevel
	log := logger.New(cfg.LogLevel, cfg.AppName, cfg.AppEnv)
	slog.SetDefault(log)

	// postgres connection
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// rabbitMQ

	// check connection
	if err := db.Ping(context.Background()); err != nil {
		slog.Error("Unable to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to the database")

	// repository
	accountRepository := accountRepo.NewAccountRepository(db)

	// application
	accountApp := accountApp.NewAccountApp(accountRepository)

	// handler
	accountHandler := accountHandler.NewAccountHandler(accountApp)

	// setup server
	server := setupServer(cfg, accountHandler)

	// graceful Shutdown
	go func() {
		if err := server.Start(":" + cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server startup error", "error", err)
			os.Exit(1)
		}
	}()

	// wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // wait for signal

	slog.Warn("Shutting down server...")

	// shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}

	slog.Info("Server shut down gracefully")
}

func setupServer(cfg config.Config, accountHandler accountHandler.Account) *echo.Echo {
	server := echo.New()

	// middlewares
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	v1 := server.Group("/v1")

	accountHandler.RegisterRoutes(v1.Group("/accounts"))

	return server
}
