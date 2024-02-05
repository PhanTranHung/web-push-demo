package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type EnvConfigs struct {
	VapidContact    string
	VapidPublicKey  string
	VapidPrivateKey string
}

func LoadEnvConfig() ConfigLoader {
	return func(config *Configs) error {

		// Load secret when running from IDE
		wd, _ := os.Getwd() // working space dir (IDE)
		envSecret := filepath.Join(wd, ".env")

		err := godotenv.Load(envSecret)
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		vapidContact := os.Getenv("VAPID_CONTACT")
		if vapidContact == "" {
			return fmt.Errorf("failed to load VAPID_CONTACT in env file")
		}

		vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
		if vapidPublicKey == "" {
			return fmt.Errorf("failed to load VAPID_PUBLIC_KEY in env file")
		}

		vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")
		if vapidPrivateKey == "" {
			return fmt.Errorf("failed to load VAPID_PRIVATE_KEY in env file")
		}

		config.env.VapidContact = vapidContact
		config.env.VapidPublicKey = vapidPublicKey
		config.env.VapidPrivateKey = vapidPrivateKey

		return nil

	}
}
