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
	godotenv.Load() // loading .env from cmd/CoCMetaAnalyser
	cocToken := os.Getenv("COC_API")
	if cocToken == "" {
		return nil, fmt.Errorf("Please provide COC_API in .env")
	}
	return &Config{
		CocToken: cocToken,
	}, nil
}
