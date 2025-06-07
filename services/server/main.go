package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"mtguru/packages/config"
	"mtguru/packages/custom_logger"
	"net/http"
	"os"

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
}

type MTGuruSearchResponse struct {
	Count   int         `json:"count"`
	Matches []CardMatch `json:"matches"`
}

type CardMatch struct {
	ID       string   `json:"id"`
	Score    float64  `json:"score"`
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

func init() {
	// init is called before main, so we set up our logger and config here
	custom_logger.CreateLogger()
	config.CreateConfig()
}

// /api/search
func searchHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody MTGuruSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		slog.Debug("Error decoding request body", "error", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	requestBody.Count = 16 // Hard coded number of results wanted back, should make paginated or adjustable moving forwards

	slog.Info("Received search request:", "query", requestBody.Query, "filters", requestBody.Filters)

	var payloadJson, err = json.Marshal(requestBody)
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Debug("Making search request to Cloudflare:", "request", string(payloadJson))

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

	var searchResults MTGuruSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		slog.Error(err.Error())
	}

	searchResultsJson, err := json.Marshal(searchResults)
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
	w.Write(searchResultsJson)

}

func initHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/search", searchHandler)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return cors.Default().Handler(mux)

}

func main() {
	handler := initHandler()
	slog.Info("Starting server on lambda...")
	lambda.Start(httpadapter.NewV2(handler).ProxyWithContext)
}
