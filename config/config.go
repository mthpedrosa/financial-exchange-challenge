package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string `mapstructure:"APP_ENV"      validate:"required,oneof=development production"`
	AppName     string `mapstructure:"APP_NAME"     validate:"required"`
	Port        string `mapstructure:"PORT"         validate:"required"`
	LogLevel    string `mapstructure:"LOG_LEVEL"    validate:"required"`
	DatabaseURL string `mapstructure:"DATABASE_URL" validate:"required"`
	RabbitURL   string `mapstructure:"RABBITMQ_URL"   validate:"required"`
	JWTSecret   string `mapstructure:"JWT_SECRET"     validate:"required"`
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, reading environment variables.")
	}

	// load config
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        getEnv("PORT", "8080"),
		AppEnv:      os.Getenv("APP_ENV"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		RabbitURL:   os.Getenv("RABBITMQ_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		AppName:     getEnv("APP_NAME", "Exchange API"),
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		log.Fatalf("Error: Invalid configuration. %v", err)
	}

	log.Println("Configuration loaded successfully.")
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
