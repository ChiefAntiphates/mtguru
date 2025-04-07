package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		fmt.Println("Error creating Weaviate client:", err)
	}

	// createIndex and populateIndex only need to be run once to set up the database with schema and populate an index with data
	// populateIndex can take quite a while to run, the default-cards json file has like 100k+ cards
	// createIndex(client)
	// populateIndex(client)

	// example query to get cards with a specific concept
	concept := "spells to draw cards"

	ctx := context.Background()
	response, err := client.GraphQL().Get().
		WithClassName("Card").
		WithFields(
			graphql.Field{Name: "name"},
			graphql.Field{Name: "mana_cost"},
			graphql.Field{Name: "type_line"},
			graphql.Field{Name: "oracle_text"},
			graphql.Field{Name: "power"},
			graphql.Field{Name: "toughness"},
			graphql.Field{Name: "loyalty"},
			graphql.Field{Name: "colors"},
			graphql.Field{Name: "set_name"},
			graphql.Field{Name: "keywords"},
			graphql.Field{Name: "flavor_text"},
			graphql.Field{Name: "rarity"},
		).
		WithNearText(client.GraphQL().NearTextArgBuilder().
			WithConcepts([]string{concept})).
		WithLimit(10).
		Do(ctx)

	if err != nil {
		fmt.Println("Error executing GraphQL query:", err)
		return
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling response to JSON:", err)
		return
	}
	fmt.Println("====")
	fmt.Println(string(jsonResponse))

}

func createIndex(client *weaviate.Client) {
	classObj := &models.Class{
		Class:      "Card",
		Vectorizer: "text2vec-ollama",
		ModuleConfig: map[string]interface{}{
			"text2vec-ollama": map[string]interface{}{
				"apiEndpoint": "http://host.docker.internal:11434",
				"model":       "nomic-embed-text",
			},
			"generative-ollama": map[string]interface{}{
				"apiEndpoint": "http://host.docker.internal:11434",
				"model":       "llama3.2",
			},
		},
		Properties: []*models.Property{
			{
				Name:     "object",
				DataType: []string{"string"},
			},
			{
				Name:     "scryfall_id",
				DataType: []string{"string"},
			},
			{
				Name:     "oracle_id",
				DataType: []string{"string"},
			},
			{
				Name:     "multiverse_ids",
				DataType: []string{"int[]"},
			},
			{
				Name:     "mtgo_id",
				DataType: []string{"int"},
			},
			{
				Name:     "tcgplayer_id",
				DataType: []string{"int"},
			},
			{
				Name:     "name",
				DataType: []string{"string"},
			},
			{
				Name:     "released_at",
				DataType: []string{"string"},
			},
			{
				Name:     "scryfall_uri",
				DataType: []string{"string"},
			},
			// {
			// 	Name:     "image_uris",
			// 	DataType: []string{"object"},
			// },
			{
				Name:     "mana_cost",
				DataType: []string{"string"},
			},
			{
				Name:     "cmc",
				DataType: []string{"number"},
			},
			{
				Name:     "type_line",
				DataType: []string{"string"},
			},
			{
				Name:     "oracle_text",
				DataType: []string{"string"},
			},
			{
				Name:     "power",
				DataType: []string{"string"},
			},
			{
				Name:     "toughness",
				DataType: []string{"string"},
			},
			{
				Name:     "defense",
				DataType: []string{"string"},
			},
			{
				Name:     "loyalty",
				DataType: []string{"string"},
			},
			{
				Name:     "hand_modifier",
				DataType: []string{"string"},
			},
			{
				Name:     "life_modifier",
				DataType: []string{"string"},
			},
			{
				Name:     "colors",
				DataType: []string{"string[]"},
			},
			{
				Name:     "color_identity",
				DataType: []string{"string[]"},
			},
			{
				Name:     "keywords",
				DataType: []string{"string[]"},
			},
			{
				Name:     "produced_mana",
				DataType: []string{"string[]"},
			},
			// {
			// 	Name:     "legalities",
			// 	DataType: []string{"object"},
			// },
			{
				Name:     "games",
				DataType: []string{"string[]"},
			},
			{
				Name:     "reserved",
				DataType: []string{"boolean"},
			},
			{
				Name:     "game_changer",
				DataType: []string{"boolean"},
			},
			{
				Name:     "finishes",
				DataType: []string{"string[]"},
			},
			{
				Name:     "set_id",
				DataType: []string{"string"},
			},
			{
				Name:     "set_name",
				DataType: []string{"string"},
			},
			{
				Name:     "set_type",
				DataType: []string{"string"},
			},
			{
				Name:     "rulings_uri",
				DataType: []string{"string"},
			},
			{
				Name:     "digital",
				DataType: []string{"boolean"},
			},
			{
				Name:     "rarity",
				DataType: []string{"string"},
			},
			{
				Name:     "flavor_text",
				DataType: []string{"string"},
			},
			{
				Name:     "card_back_id",
				DataType: []string{"string"},
			},
			{
				Name:     "artist",
				DataType: []string{"string"},
			},
			{
				Name:     "artist_ids",
				DataType: []string{"string[]"},
			},
			{
				Name:     "border_color",
				DataType: []string{"string"},
			},
			{
				Name:     "booster",
				DataType: []string{"boolean"},
			},
			// {
			// 	Name:     "prices",
			// 	DataType: []string{"object"},
			// },
			// {
			// 	Name:     "related_uris",
			// 	DataType: []string{"object"},
			// },
			// {
			// 	Name:     "purchase_uris",
			// 	DataType: []string{"object"},
			// },
		},
	}

	err := client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		panic(err)
	}

}

func populateIndex(client *weaviate.Client) {
	var cards []Card = parseCardsFromFile()

	// populate index with data
	objects := make([]*models.Object, len(cards))
	for i := range cards {
		objects[i] = &models.Object{
			Class: "Card",
			Properties: map[string]any{
				"object":         cards[i].Object,
				"scryfall_id":    cards[i].ScryfallID,
				"oracle_id":      cards[i].OracleID,
				"multiverse_ids": cards[i].MultiverseIDs,
				"mtgo_id":        cards[i].MtgoID,
				"tcgplayer_id":   cards[i].TcgplayerID,
				"name":           cards[i].Name,
				"released_at":    cards[i].ReleasedAt,
				"scryfall_uri":   cards[i].ScryfallURI,
				// "image_uris":     cards[i].ImageURIs,
				"mana_cost":      cards[i].ManaCost,
				"cmc":            cards[i].Cmc,
				"type_line":      cards[i].TypeLine,
				"oracle_text":    cards[i].OracleText,
				"power":          cards[i].Power,
				"toughness":      cards[i].Toughness,
				"defense":        cards[i].Defence,
				"loyalty":        cards[i].Loyalty,
				"hand_modifier":  cards[i].HandModifier,
				"life_modifier":  cards[i].LifeModifier,
				"colors":         cards[i].Colors,
				"color_identity": cards[i].ColorIdentity,
				"keywords":       cards[i].Keywords,
				"produced_mana":  cards[i].ProducedMana,
				// "legalities":     cards[i].Legalities,
				"games":        cards[i].Games,
				"reserved":     cards[i].Reserved,
				"game_changer": cards[i].GameChanger,
				"finishes":     cards[i].Finishes,
				"set_id":       cards[i].SetID,
				"set_name":     cards[i].SetName,
				"set_type":     cards[i].SetType,
				"rulings_uri":  cards[i].RulingsURI,
				"digital":      cards[i].Digital,
				"rarity":       cards[i].Rarity,
				"flavor_text":  cards[i].FlavorText,
				"card_back_id": cards[i].CardBackID,
				"artist":       cards[i].Artist,
				"artist_ids":   cards[i].ArtistIDs,
				"border_color": cards[i].BorderColor,
				"booster":      cards[i].Booster,
				// "prices":         cards[i].Prices,
				// "related_uris":   cards[i].RelatedURIs,
				// "purchase_uris":  cards[i].PurchaseURIs,
			},
		}
	}

	// batch write items
	batchSize := 100 //
	for i := 0; i < len(objects); i += batchSize {
		end := i + batchSize
		if end > len(objects) {
			end = len(objects)
		}

		fmt.Print("Batching objects from index ", i, " to ", end, "\n")
		batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects[i:end]...).Do(context.Background())
		if err != nil {
			fmt.Println("Batch operation failed:", err.Error())
			return
		}

		for _, res := range batchRes {
			if res.Result.Errors != nil {
				for _, batchError := range res.Result.Errors.Error {
					fmt.Printf("Batch error: %+v\n", batchError)
				}
			}
		}
	}
}

func parseCardsFromFile() []Card {
	jsonFile, err := os.Open("default-cards-20250404213404.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var cards []Card
	json.Unmarshal(byteValue, &cards)

	fmt.Println("Number of cards:", len(cards))

	for i := 0; i < len(cards); i++ {
		fmt.Printf("Card ID: "+cards[i].ScryfallID+" Remaining cards: %d\n", len(cards)-i-1)
	}

	return cards
}

type Card struct {
	Object        string `json:"object"`
	ScryfallID    string `json:"id"`
	OracleID      string `json:"oracle_id"`
	MultiverseIDs []int  `json:"multiverse_ids"`
	MtgoID        int    `json:"mtgo_id"`
	TcgplayerID   int    `json:"tcgplayer_id"`
	Name          string `json:"name"`
	ReleasedAt    string `json:"released_at"`
	ScryfallURI   string `json:"scryfall_uri"`
	// ImageURIs     map[string]string  `json:"image_uris"`
	ManaCost      string   `json:"mana_cost"`
	Cmc           float64  `json:"cmc"`
	TypeLine      string   `json:"type_line"`
	OracleText    string   `json:"oracle_text"`
	Power         string   `json:"power"`
	Toughness     string   `json:"toughness"`
	Defence       string   `json:"defense"`
	Loyalty       string   `json:"loyalty"`
	HandModifier  string   `json:"hand_modifier"`
	LifeModifier  string   `json:"life_modifier"`
	Colors        []string `json:"colors"`
	ColorIdentity []string `json:"color_identity"`
	Keywords      []string `json:"keywords"`
	ProducedMana  []string `json:"produced_mana"`
	// Legalities    map[string]string  `json:"legalities"`
	Games       []string `json:"games"`
	Reserved    bool     `json:"reserved"`
	GameChanger bool     `json:"game_changer"`
	Finishes    []string `json:"finishes"`
	SetID       string   `json:"set_id"`
	SetName     string   `json:"set_name"`
	SetType     string   `json:"set_type"`
	RulingsURI  string   `json:"rulings_uri"`
	Digital     bool     `json:"digital"`
	Rarity      string   `json:"rarity"`
	FlavorText  string   `json:"flavor_text"`
	CardBackID  string   `json:"card_back_id"`
	Artist      string   `json:"artist"`
	ArtistIDs   []string `json:"artist_ids"`
	BorderColor string   `json:"border_color"`
	Booster     bool     `json:"booster"`
	// Prices        map[string]float64 `json:"prices"`
	// RelatedURIs   map[string]string  `json:"related_uris"`
	// PurchaseURIs  map[string]string  `json:"purchase_uris"`
}
