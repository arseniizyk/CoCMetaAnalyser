package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arseniizyk/CoCMetaAnalyser/internal/pkg/cocdev"
	"github.com/joho/godotenv"
)

type Config struct {
	CocToken string
	Cleanup  func()
}

func New() (*Config, error) {
	// Try loading .env from current directory or parent directories if available
	_ = godotenv.Load(".env", "../../.env", "../.env")

	cocToken := os.Getenv("COC_API")
	if cocToken != "" {
		return &Config{
			CocToken: cocToken,
		}, nil
	}

	// If COC_API is not provided, check if COC_EMAIL and COC_PASSWORD are set for automated login
	email := os.Getenv("COC_EMAIL")
	password := os.Getenv("COC_PASSWORD")

	if email != "" && password != "" {
		log.Println("COC_API not provided. Attempting automated login via Supercell Developer Portal...")
		devClient, err := cocdev.New(email, password)
		if err != nil {
			return nil, fmt.Errorf("initializing Supercell developer client: %w", err)
		}

		if err := devClient.Login(); err != nil {
			return nil, fmt.Errorf("Supercell portal login failed: %w", err)
		}

		ip, err := devClient.GetPublicIP()
		if err != nil {
			return nil, fmt.Errorf("detecting runner IP: %w", err)
		}

		// Clean up any stale temporary keys to stay within 10-key limit
		_ = devClient.CleanupOldKeys("coc-meta-temp-")

		keyName := fmt.Sprintf("coc-meta-temp-%d", time.Now().Unix())
		key, err := devClient.CreateKey(keyName, ip)
		if err != nil {
			return nil, fmt.Errorf("generating temporary key for IP %s: %w", ip, err)
		}

		log.Printf("Successfully generated temporary key for runner IP %s", ip)

		cleanup := func() {
			log.Println("Cleaning up temporary Supercell API key...")
			if err := devClient.RevokeKey(key.ID); err != nil {
				log.Printf("Warning: failed to revoke temporary key: %v", err)
			} else {
				log.Println("Temporary API key successfully revoked.")
			}
		}

		return &Config{
			CocToken: key.Key,
			Cleanup:  cleanup,
		}, nil
	}

	return nil, fmt.Errorf("please provide COC_API (or COC_EMAIL and COC_PASSWORD) in .env or environment variables")
}
