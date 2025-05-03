package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type EnvironmentConfig struct {
	WEAVIATE_URL     string `toml:"WEAVIATE_URL"`
	WEAVIATE_API_KEY string `toml:"WEAVIATE_API_KEY"`
	OPEN_API_KEY     string `toml:"OPEN_API_KEY"`
}

type Environments struct {
	Localhost EnvironmentConfig `toml:"localhost"`
	Prod      EnvironmentConfig `toml:"prod"`
}

func CreateConfig() EnvironmentConfig {

	file, err := os.Open("config.toml")

	if err != nil {
		slog.Error(err.Error())
	}

	defer file.Close()

	configBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error(err.Error())
	}

	var cfg Environments

	err = toml.Unmarshal(configBytes, &cfg)
	if err != nil {
		slog.Error(err.Error())
	}

	fmt.Printf("Config struct: %+v\n", cfg)
	env := os.Getenv("mtguru_env") // e.g., "development", "staging", "production"

	var activeConfig EnvironmentConfig

	switch env {
	case "localhost":
		activeConfig = cfg.Localhost
	case "prod":
		activeConfig = cfg.Prod
	default:
		slog.Info("Unknown environment:", env, ". Defaulting to localhost")
		activeConfig = cfg.Localhost // Default to localhost if unknown
	}

	slog.Info("Config loaded successfully")
	slog.Info("WEAVIATE_URL", "weaviate_url", activeConfig.WEAVIATE_URL)
	slog.Info("WEAVIATE_API_KEY:", "weaviate_api_key", activeConfig.WEAVIATE_API_KEY)
	slog.Info("OPEN_API_KEY:", "open_api_key", activeConfig.OPEN_API_KEY)

	return activeConfig
}
