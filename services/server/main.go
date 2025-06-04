package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"mtguru/packages/config"
	"mtguru/packages/custom_logger"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/rs/cors"
)

type MTGuruSearchRequestFilters struct {
	SetType string `json:"set_type"`
	Color   string `json:"colors"`
	Rarity  string `json:"rarity"`
}

type MTGuruSearchRequest struct {
	Query   string                     `json:"query"`
	Count   int                        `json:"count,omitempty"`
	Filters MTGuruSearchRequestFilters `json:"filters"`
	// Filters map[string]string `json:"filters"`
}

type MTGuruSearchResponse struct {
	Count   int         `json:"count"`
	Matches []CardMatch `json:"matches"`
}

type CardMatch struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
	// Values   []interface{} `json:"values"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Colors      []string `json:"colors"`
	ImageURL    string   `json:"image_url"`
	ScryfallURI string   `json:"scryfall_uri"`
	Name        string   `json:"name"`
	Rarity      string   `json:"rarity"`
	ReleaseDate string   `json:"release_date"`
	SetName     string   `json:"set_name"`
}

var activeConfig config.EnvironmentConfig

func init() {
	// init is called before main, so we can set up our logger and client here
	custom_logger.CreateLogger()
	activeConfig = config.CreateConfig()
	// client = createClient(activeConfig)
}

// func searchDatabase(search_string string, search_filters MTGuruSearchRequestFilters) *models.GraphQLResponse {

// 	// search_string := "make my units fly"

// 	slog.Info(fmt.Sprintf("SetType: %v", search_filters.SetType))
// 	slog.Info(fmt.Sprintf("Color: %v", search_filters.Color))
// 	slog.Info(fmt.Sprintf("Rarity: %v", search_filters.Rarity))

// 	ctx := context.Background()

// 	where := filters.Where().
// 		WithOperator(filters.And).
// 		WithOperands([]*filters.WhereBuilder{
// 			filters.Where().
// 				WithPath([]string{"set_type"}).
// 				WithOperator(filters.NotEqual).
// 				WithValueString("token"),
// 			filters.Where().
// 				WithPath([]string{"set_type"}).
// 				WithOperator(filters.NotEqual).
// 				WithValueString("memorabilia"),
// 		})

// 	// if search_filters.SetType != "" {
// 	// 	where = where.WithOperands([]*filters.WhereBuilder{
// 	// 		filters.Where().
// 	// 			WithPath([]string{"set_type"}).
// 	// 			WithOperator(filters.Equal).
// 	// 			WithValueString(search_filters.SetType),
// 	// 	})
// 	// }
// 	slog.Debug(string(search_filters.Color[0]))
// 	// if search_filters.Color != "" {
// 	// 	string_color := string(search_filters.Color[0])
// 	// 	where = where.WithOperands([]*filters.WhereBuilder{
// 	// 		filters.Where().
// 	// 			WithPath([]string{"colors"}).
// 	// 			WithOperator(filters.ContainsAny).
// 	// 			WithValueString(strings.ToUpper(string_color)),
// 	// 	})
// 	// }

// 	// if search_filters.Rarity != "" {
// 	// 	where = where.WithOperands([]*filters.WhereBuilder{
// 	// 		filters.Where().
// 	// 			WithPath([]string{"rarity"}).
// 	// 			WithOperator(filters.Equal).
// 	// 			WithValueString(search_filters.Rarity),
// 	// 	})
// 	// }

// 	response, err := client.GraphQL().Get().
// 		WithClassName("Mtguru").
// 		// WithFields is used to specify the fields you want to retrieve from the cards matched in the json resposne
// 		WithFields(
// 			graphql.Field{Name: "name"},
// 			// graphql.Field{Name: "mana_cost"},
// 			// graphql.Field{Name: "type_line"},
// 			graphql.Field{Name: "oracle_text"},
// 			// graphql.Field{Name: "power"},
// 			// graphql.Field{Name: "toughness"},
// 			// graphql.Field{Name: "loyalty"},
// 			graphql.Field{Name: "colors"},
// 			graphql.Field{Name: "set_name"},
// 			// graphql.Field{Name: "keywords"},
// 			// graphql.Field{Name: "flavor_text"},
// 			// graphql.Field{Name: "rarity"},
// 			graphql.Field{Name: "set_type"},
// 			graphql.Field{Name: "scryfall_uri"},
// 			graphql.Field{Name: "image_uris", Fields: []graphql.Field{
// 				{Name: "normal"},
// 				{Name: "large"},
// 			}},
// 			graphql.Field{Name: "_additional", Fields: []graphql.Field{
// 				{Name: "distance"}}},
// 		).
// 		WithNearText(client.GraphQL().NearTextArgBuilder().
// 			WithConcepts([]string{search_string})).
// 		WithLimit(29).
// 		WithWhere(where).
// 		Do(ctx)

// 	if err != nil {
// 		slog.Debug(err.Error())
// 	}

// 	slog.Info("Prompt:", "prompt", search_string)
// 	slog.Debug("Response:", "matches", response)

// 	return response
// }

func searchHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody MTGuruSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		slog.Debug("Error decoding request body", "error", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// any further checks on request body here

	requestBody.Count = 16

	slog.Info("Received search request:", "query", requestBody.Query, "filters", requestBody.Filters)

	var payloadJson, err = json.Marshal(requestBody)
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info(string(payloadJson))

	// Make request to create Guru Insert
	req, err := http.NewRequest("GET", os.Getenv("CLOUDFLARE_WORKER_URL")+"/search", bytes.NewBuffer(payloadJson))
	if err != nil {
		slog.Error(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error())
	}

	// Close response body
	defer resp.Body.Close()

	var response MTGuruSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error(err.Error())
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		slog.Debug("Error marshalling response", "error", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "https://mtguru.com")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

func initHandler() http.Handler {
	mux := http.NewServeMux()

	// mux.HandleFunc("GET /api/health", alive)
	mux.HandleFunc("/api/search", searchHandler)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return cors.Default().Handler(mux)

}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.Info("Received request:", "request", request)

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Hello from Lambda!",
	}

	return response, nil
}

func main() {

	handler := initHandler()
	// slog.Info("Starting server on port 8888...")
	lambda.Start(httpadapter.NewV2(handler).ProxyWithContext)
	// http.ListenAndServe(":8888", handler)
	// lambda.Start(handler)
}
