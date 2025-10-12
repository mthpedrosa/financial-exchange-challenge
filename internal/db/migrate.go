package db

import (
	"embed"
	"errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(databaseURL string) {
	slog.Info("Checking for database migrations...")

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		slog.Error("Failed to create migration source from embedded files", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, databaseURL)
	if err != nil {
		slog.Error("Failed to create migrate instance", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("Database migrations checked successfully")
}
