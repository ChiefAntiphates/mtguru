package main

import (
	"context"
	"log/slog"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"

	"mtguru/packages/config"
	"mtguru/packages/custom_logger"

	"github.com/joho/godotenv"
)

var activeConfig config.EnvironmentConfig
var client *weaviate.Client

func init() {
	custom_logger.CreateLogger()
	activeConfig = config.CreateConfig()
	// client = createClient(activeConfig)
}

func createClient(conf config.EnvironmentConfig) *weaviate.Client {

	cfg := weaviate.Config{
		Host:   conf.WEAVIATE_URL,
		Scheme: "http",
		// AuthConfig: auth.ApiKey{Value: conf.WEAVIATE_API_KEY},
		Headers: map[string]string{
			"X-OpenAI-Api-Key": conf.OPEN_API_KEY,
		},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		slog.Debug(err.Error())
	}

	live, err := client.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		slog.Debug(err.Error())
	}

	slog.Debug("%v", "is_remote_server_up", live)
	return client

}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}
	// createIndex(client)
	populateIndex()
	// searchDatabase(client)
	// updateCollection(client)
}
