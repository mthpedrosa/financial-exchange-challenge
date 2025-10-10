package config_test

import (
	"testing"

	"github.com/mthpedrosa/financial-exchange-challenge/config"
)

func TestLoadConfig_Success(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("PORT", "3000")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	t.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	t.Setenv("JWT_SECRET", "test-secret")

	cfg := config.LoadConfig()

	// check if the values are as expected
	if cfg.Port != "3000" {
		t.Errorf("esperado Port '3000', mas recebi '%s'", cfg.Port)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("esperado AppEnv 'development', mas recebi '%s'", cfg.AppEnv)
	}
	if cfg.DatabaseURL != "postgres://test:test@localhost:5432/testdb" {
		t.Errorf("DatabaseURL não corresponde ao esperado")
	}
}
