package main

import (
	"context"
	"log/slog"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"

	"mtguru/packages/custom_logger"
)

var activeConfig EnvironmentConfig

func init() {
	custom_logger.CreateLogger()
	activeConfig = CreateConfig()
}

func createClient(conf EnvironmentConfig) *weaviate.Client {

	cfg := weaviate.Config{
		Host:       conf.WEAVIATE_URL,
		Scheme:     "https",
		AuthConfig: auth.ApiKey{Value: conf.WEAVIATE_API_KEY},
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

func searchDatabase(client *weaviate.Client) {

	ctx := context.Background()
	response, err := client.GraphQL().Get().
		WithClassName("Mtguru").
		// WithFields is used to specify the fields you want to retrieve from the cards matched
		WithFields(
			graphql.Field{Name: "name"},
			// graphql.Field{Name: "mana_cost"},
			// graphql.Field{Name: "type_line"},
			graphql.Field{Name: "oracle_text"},
			// graphql.Field{Name: "power"},
			// graphql.Field{Name: "toughness"},
			// graphql.Field{Name: "loyalty"},
			// graphql.Field{Name: "colors"},
			graphql.Field{Name: "set_name"},
			// graphql.Field{Name: "keywords"},
			// graphql.Field{Name: "flavor_text"},
			// graphql.Field{Name: "rarity"},
			graphql.Field{Name: "_additional", Fields: []graphql.Field{
				{Name: "distance"}}},
		).
		WithNearText(client.GraphQL().NearTextArgBuilder().
			WithConcepts([]string{"make dragons stronger"})).
		WithLimit(3).
		Do(ctx)

	if err != nil {
		slog.Debug(err.Error())
	}

	slog.Info("Response:", "matches", response)

}

func main() {
	client := createClient(activeConfig)
	// createIndex(client)
	// populateIndex(client)
	searchDatabase(client)
}
