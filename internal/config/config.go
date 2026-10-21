package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CocToken string
}

func New() (*Config, error) {
	// Try loading .env from current directory or parent directories if available
	_ = godotenv.Load(".env", "../../.env", "../.env")
	cocToken := os.Getenv("COC_API")
	if cocToken == "" {
		return nil, fmt.Errorf("please provide COC_API in .env or environment variable")
	}
	return &Config{
		CocToken: cocToken,
	}, nil
}
