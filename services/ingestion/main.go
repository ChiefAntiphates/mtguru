package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
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

func createClient(conf EnvironmentConfig) {

	fmt.Println("Test")

	cfg := weaviate.Config{
		Host:       conf.WEAVIATE_URL,
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: conf.WEAVIATE_API_KEY},
		Headers:    nil,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Println(err)
	}

	live, err := client.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v", live)

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err := os.Open("config.toml")

	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	configBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	var cfg Environments

	err = toml.Unmarshal(configBytes, &cfg)
	if err != nil {
		log.Println("Error unmarshalling TOML:", err)
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
		log.Println("Unknown environment:", env)
		activeConfig = cfg.Localhost // Default to localhost if unknown
	}

	fmt.Println("Config loaded successfully")
	fmt.Println("WEAVIATE_URL:", activeConfig.WEAVIATE_URL)
	fmt.Println("WEAVIATE_API_KEY:", activeConfig.WEAVIATE_API_KEY)
	fmt.Println("OPEN_API_KEY:", activeConfig.OPEN_API_KEY)

	createClient(activeConfig)
}
