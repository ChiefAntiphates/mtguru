package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"

	"mtguru/packages/custom_logger"
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

func init() {
	custom_logger.CreateLogger()
}

func createClient(conf EnvironmentConfig) {

	cfg := weaviate.Config{
		Host:       conf.WEAVIATE_URL,
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: conf.WEAVIATE_API_KEY},
		Headers:    nil,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		slog.Debug(err.Error())
	}

	live, err := client.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		slog.Debug(err.Error())
	}

	slog.Info("%v", live)

}

func main() {

	file, err := os.Open("config.toml")

	if err != nil {
		slog.Error(err.Error())
	}

	defer file.Close()

	configBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	var cfg Environments

	err = toml.Unmarshal(configBytes, &cfg)
	if err != nil {
		slog.Error(err.Error())
		return
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

	createClient(activeConfig)
}
