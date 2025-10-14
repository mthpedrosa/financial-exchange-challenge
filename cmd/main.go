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

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mthpedrosa/financial-exchange-challenge/config"
	accountHandler "github.com/mthpedrosa/financial-exchange-challenge/internal/account/adapters/api"
	accountRepo "github.com/mthpedrosa/financial-exchange-challenge/internal/account/adapters/repository"
	accountApp "github.com/mthpedrosa/financial-exchange-challenge/internal/account/app"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/db"
	instrumentHandler "github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/adapters/api"
	instrumentRepo "github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/adapters/repository"
	instrumentApp "github.com/mthpedrosa/financial-exchange-challenge/internal/instrument/app"
	echoSwagger "github.com/swaggo/echo-swagger"

	balanceHandler "github.com/mthpedrosa/financial-exchange-challenge/internal/balance/adapters/api"
	balanceRepo "github.com/mthpedrosa/financial-exchange-challenge/internal/balance/adapters/repository"
	balanceApp "github.com/mthpedrosa/financial-exchange-challenge/internal/balance/app"

	_ "github.com/mthpedrosa/financial-exchange-challenge/docs"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/logger"
	orderHandler "github.com/mthpedrosa/financial-exchange-challenge/internal/order/adapters/api"
	orderRepo "github.com/mthpedrosa/financial-exchange-challenge/internal/order/adapters/repository"
	orderApp "github.com/mthpedrosa/financial-exchange-challenge/internal/order/app"
	amqp "github.com/rabbitmq/amqp091-go"
)

// @title           Financial Exchange Challenge API
// @version         1.0
// @description     Esta é a documentação da API para o desafio.
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1
func main() {
	// load config
	cfg := config.LoadConfig()

	// setup logger based on cfg.LogLevel
	log := logger.New(cfg.LogLevel, cfg.AppName, cfg.AppEnv)
	slog.SetDefault(log)

	// migrations
	db.RunMigrations(cfg.DatabaseURL)

	// postgres connection
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// rabbitMQ
	rabbitConn, err := amqp.Dial(cfg.RabbitURL)
	if err != nil {
		slog.Error("Unable to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		slog.Error("Unable to open RabbitMQ channel", "error", err)
		os.Exit(1)
	}
	defer rabbitChannel.Close()

	queue, err := rabbitChannel.QueueDeclare(
		"orders", // queue name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		slog.Error("Unable to declare RabbitMQ queue", "error", err)
		os.Exit(1)
	}

	// check connection
	if err := db.Ping(context.Background()); err != nil {
		slog.Error("Unable to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to the database")

	// repository
	accountRepository := accountRepo.NewAccountRepository(db)
	instrumentRepository := instrumentRepo.NewInstrumentRepository(db)
	balanceRepository := balanceRepo.NewBalanceRepository(db)
	orderQueueRepository := orderRepo.NewOrderQueueRepository(rabbitChannel, queue.Name)
	orderRepository := orderRepo.NewOrderRepository(db)

	// application
	accountApp := accountApp.NewAccountApp(accountRepository)
	instrumentApp := instrumentApp.NewInstrumentApp(instrumentRepository)
	balanceApp := balanceApp.NewBalanceApp(balanceRepository, accountRepository)
	orderApp := orderApp.NewOrderApp(
		orderRepository,
		accountRepository,
		instrumentRepository,
		balanceRepository,
		orderQueueRepository,
	)

	// handler
	accountHandler := accountHandler.NewAccountHandler(accountApp)
	instrumentHandler := instrumentHandler.NewInstrumentHandler(instrumentApp)
	balanceHandler := balanceHandler.NewBalanceHandler(balanceApp)
	orderHandler := orderHandler.NewOrderHandler(orderApp)

	// setup server
	server := setupServer(cfg, accountHandler, instrumentHandler, balanceHandler, orderHandler)

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

func setupServer(cfg config.Config, accountHandler accountHandler.Account, instrumentHandler instrumentHandler.Instrument, balanceHandler balanceHandler.Balance, orderHandler orderHandler.Order) *echo.Echo {
	server := echo.New()

	// middlewares
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// swagger
	server.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := server.Group("/v1")

	accountHandler.RegisterRoutes(v1.Group("/accounts"))
	instrumentHandler.RegisterRoutes(v1.Group("/instruments"))
	balanceHandler.RegisterRoutes(v1.Group("/balances"))
	orderHandler.RegisterRoutes(v1.Group("/orders"))

	return server
}
