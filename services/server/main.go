package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mtguru/packages/config"
	"mtguru/packages/custom_logger"
	"net/http"

	"github.com/rs/cors"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type MTGuruSearchRequestFilters struct {
	SetType string `json:"set_type"`
	Color   string `json:"colors"`
	Rarity  string `json:"rarity"`
}

type MTGuruSearchRequest struct {
	Query   string                     `json:"query"`
	Filters MTGuruSearchRequestFilters `json:"filters"`
	// Filters map[string]string `json:"filters"`
}

var activeConfig config.EnvironmentConfig
var client *weaviate.Client

func init() {
	// init is called before main, so we can set up our logger and client here
	custom_logger.CreateLogger()
	activeConfig = config.CreateConfig()
	client = createClient(activeConfig)
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

func searchDatabase(search_string string, search_filters MTGuruSearchRequestFilters) *models.GraphQLResponse {

	// search_string := "make my units fly"

	slog.Info(fmt.Sprintf("SetType: %v", search_filters.SetType))
	slog.Info(fmt.Sprintf("Color: %v", search_filters.Color))
	slog.Info(fmt.Sprintf("Rarity: %v", search_filters.Rarity))

	ctx := context.Background()

	where := filters.Where().
		WithOperator(filters.And).
		WithOperands([]*filters.WhereBuilder{
			filters.Where().
				WithPath([]string{"set_type"}).
				WithOperator(filters.NotEqual).
				WithValueString("token"),
			filters.Where().
				WithPath([]string{"set_type"}).
				WithOperator(filters.NotEqual).
				WithValueString("memorabilia"),
		})

	// if search_filters.SetType != "" {
	// 	where = where.WithOperands([]*filters.WhereBuilder{
	// 		filters.Where().
	// 			WithPath([]string{"set_type"}).
	// 			WithOperator(filters.Equal).
	// 			WithValueString(search_filters.SetType),
	// 	})
	// }
	slog.Debug(string(search_filters.Color[0]))
	// if search_filters.Color != "" {
	// 	string_color := string(search_filters.Color[0])
	// 	where = where.WithOperands([]*filters.WhereBuilder{
	// 		filters.Where().
	// 			WithPath([]string{"colors"}).
	// 			WithOperator(filters.ContainsAny).
	// 			WithValueString(strings.ToUpper(string_color)),
	// 	})
	// }

	// if search_filters.Rarity != "" {
	// 	where = where.WithOperands([]*filters.WhereBuilder{
	// 		filters.Where().
	// 			WithPath([]string{"rarity"}).
	// 			WithOperator(filters.Equal).
	// 			WithValueString(search_filters.Rarity),
	// 	})
	// }

	response, err := client.GraphQL().Get().
		WithClassName("Mtguru").
		// WithFields is used to specify the fields you want to retrieve from the cards matched in the json resposne
		WithFields(
			graphql.Field{Name: "name"},
			// graphql.Field{Name: "mana_cost"},
			// graphql.Field{Name: "type_line"},
			graphql.Field{Name: "oracle_text"},
			// graphql.Field{Name: "power"},
			// graphql.Field{Name: "toughness"},
			// graphql.Field{Name: "loyalty"},
			graphql.Field{Name: "colors"},
			graphql.Field{Name: "set_name"},
			// graphql.Field{Name: "keywords"},
			// graphql.Field{Name: "flavor_text"},
			// graphql.Field{Name: "rarity"},
			graphql.Field{Name: "set_type"},
			graphql.Field{Name: "scryfall_uri"},
			graphql.Field{Name: "image_uris", Fields: []graphql.Field{
				{Name: "normal"},
				{Name: "large"},
			}},
			graphql.Field{Name: "_additional", Fields: []graphql.Field{
				{Name: "distance"}}},
		).
		WithNearText(client.GraphQL().NearTextArgBuilder().
			WithConcepts([]string{search_string})).
		WithLimit(29).
		WithWhere(where).
		Do(ctx)

	if err != nil {
		slog.Debug(err.Error())
	}

	slog.Info("Prompt:", "prompt", search_string)
	slog.Debug("Response:", "matches", response)

	return response
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

	results := searchDatabase(requestBody.Query, requestBody.Filters)

	responseJSON, err := json.Marshal(results)
	if err != nil {
		slog.Debug("Error marshalling response", "error", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	// searchDatabase(requestBody.Query)
}

func initHandler() http.Handler {
	mux := http.NewServeMux()

	// mux.HandleFunc("GET /api/health", alive)
	mux.HandleFunc("POST /api/search", searchHandler)
	return cors.Default().Handler(mux)

}

func main() {

	handler := initHandler()
	slog.Info("Starting server on port 8888...")
	http.ListenAndServe(":8888", handler)

}
