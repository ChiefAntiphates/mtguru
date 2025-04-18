package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/rs/cors"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"

	"mtguru/packages/custom_logger"
)

var activeConfig EnvironmentConfig
var client *weaviate.Client

func init() {
	custom_logger.CreateLogger()
	activeConfig = CreateConfig()
	client = createClient(activeConfig)
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

func searchDatabase(search_string string) {

	// search_string := "make my units fly"

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
			WithConcepts([]string{search_string})).
		WithLimit(3).
		Do(ctx)

	if err != nil {
		slog.Debug(err.Error())
	}

	slog.Info("Prompt:", "prompt", search_string)
	slog.Info("Response:", "matches", response)

}

func alive(w http.ResponseWriter, r *http.Request) {
	slog.Info("I'm alive!")
}

type MTGuruSearchRequest struct {
	Query   string            `json:"query"`
	Filters map[string]string `json:"filters"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody MTGuruSearchRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		slog.Debug("Error decoding request body", "error", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	slog.Info("Received search request:", "query", requestBody.Query, "filters", requestBody.Filters)
	searchDatabase(requestBody.Query)
}

func initHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", alive)
	mux.HandleFunc("POST /api/search", searchHandler)
	return cors.Default().Handler(mux)

}

func main() {
	// createIndex(client)
	// populateIndex(client)
	// searchDatabase(client)
	handler := initHandler()
	slog.Info("Starting server on port 8080...")
	http.ListenAndServe(":8080", handler)

}
