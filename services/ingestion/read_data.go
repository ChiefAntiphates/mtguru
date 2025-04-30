package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mtguru/packages/custom_logger"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func init() {
	custom_logger.CreateLogger()
}

type Card struct {
	Object        string            `json:"object"`
	ScryfallID    string            `json:"id"`
	OracleID      string            `json:"oracle_id"`
	MultiverseIDs []int             `json:"multiverse_ids"`
	MtgoID        int               `json:"mtgo_id"`
	TcgplayerID   int               `json:"tcgplayer_id"`
	Name          string            `json:"name"`
	ReleasedAt    string            `json:"released_at"`
	ScryfallURI   string            `json:"scryfall_uri"`
	ImageURIs     map[string]string `json:"image_uris"`
	ManaCost      string            `json:"mana_cost"`
	Cmc           float64           `json:"cmc"`
	TypeLine      string            `json:"type_line"`
	OracleText    string            `json:"oracle_text"`
	Power         string            `json:"power"`
	Toughness     string            `json:"toughness"`
	Defence       string            `json:"defense"`
	Loyalty       string            `json:"loyalty"`
	HandModifier  string            `json:"hand_modifier"`
	LifeModifier  string            `json:"life_modifier"`
	Colors        []string          `json:"colors"`
	ColorIdentity []string          `json:"color_identity"`
	Keywords      []string          `json:"keywords"`
	ProducedMana  []string          `json:"produced_mana"`
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

func parseCardsFromFile() []Card {

	jsonFile, err := os.Open("data/oracle-cards-20250429210412.json")
	if err != nil {
		slog.Error(err.Error())
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var cards []Card
	json.Unmarshal(byteValue, &cards)

	// var uniqueCards []Card
	// var duplicateCards []Card
	// oracleIDMap := make(map[string]bool)

	// for _, card := range cards {
	// 	if oracleIDMap[card.OracleID] {
	// 		duplicateCards = append(duplicateCards, card)
	// 	} else {
	// 		oracleIDMap[card.OracleID] = true
	// 		uniqueCards = append(uniqueCards, card)
	// 	}
	// }
	// cards = uniqueCards

	slog.Info("Number of total cards", "count", len(cards))
	// slog.Info("Number of unique cards:", "count", len(uniqueCards))
	// slog.Info("Number of duplicate cards:", "count", len(duplicateCards))

	// for i := 0; i < len(cards); i++ {
	// 	fmt.Printf("Card ID: "+cards[i].ScryfallID+" Remaining cards: %d\n", len(cards)-i-1)
	// }

	return cards
}

func updateCollection(client *weaviate.Client) {
	// var cards []Card = parseCardsFromFile()

	vabatchSize := 20
	className := "Mtguru"
	classProperties := []string{"title"}

	getBatchWithCursor := func(client *weaviate.Client,
		className string, classProperties []string, batchSize int, cursor string) (*models.GraphQLResponse, error) {
		fields := []graphql.Field{}
		for _, prop := range classProperties {
			fields = append(fields, graphql.Field{Name: prop})
		}
		fields = append(fields, graphql.Field{Name: "_additional { id vector }"})

		get := client.GraphQL().Get().
			WithClassName(className).
			// Optionally retrieve the vector embedding by adding `vector` to the _additional fields
			WithFields(fields...).
			WithLimit(batchSize)

		if cursor != "" {
			return get.WithAfter(cursor).Do(context.Background())
		}
		return get.Do(context.Background())
	}

	response, err := getBatchWithCursor(client, className, classProperties, vabatchSize, "")
	if err != nil {
		slog.Error("Error fetching batch with cursor", "error", err)
		return
	}
	slog.Info("Batch fetched successfully", "response", response)
}

func populateIndex(client *weaviate.Client) {
	var cards []Card = parseCardsFromFile()

	// populate index with data
	objects := make([]*models.Object, len(cards))
	for i := range cards {
		objects[i] = &models.Object{
			Class: "mtguru",
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
				"image_uris":     cards[i].ImageURIs,
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

		slog.Info(fmt.Sprintf("Batching objects from index %d to %d\n", i, end))
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

func createIndex(client *weaviate.Client) {
	// define the collection
	classObj := &models.Class{
		Class:      "mtguru",
		Vectorizer: "text2vec-openai",
		ModuleConfig: map[string]interface{}{
			"text2vec-openai": map[string]interface{}{
				"sourceProperties": []string{"title"},
				"model":            "text-embedding-3-large",
				"dimensions":       1024, // Optional (e.g. 1024, 512, 256)
			},
			"generative-cohere": map[string]interface{}{},
			"vectorizePropertyName": map[string]bool{
				"scryfall_id":    false,
				"oracle_id":      false,
				"multiverse_ids": false,
				"mtgo_id":        false,
				"tcgplayer_id":   false,
				"scryfall_uri":   false,
				"image_uris":     false,
				"games":          false,
				"reserved":       false,
				"game_changer":   false,
				"finishes":       false,
				"set_id":         false,
				"rulings_uri":    false,
				"digital":        false,
				"card_back_id":   false,
				"artist_ids":     false,
				"border_color":   false,
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
			{
				Name:     "image_uris",
				DataType: []string{"object"},
				NestedProperties: []*models.NestedProperty{
					{
						Name:     "small",
						DataType: []string{"text"},
					},
					{
						Name:     "normal",
						DataType: []string{"text"},
					},
					{
						Name:     "large",
						DataType: []string{"text"},
					},
					{
						Name:     "png",
						DataType: []string{"text"},
					},
					{
						Name:     "art_crop",
						DataType: []string{"text"},
					},
					{
						Name:     "border_crop",
						DataType: []string{"text"},
					},
				},
			},
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

	slog.Info("Creating collection 'mtguru'...")

	err := client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		panic(err)
	}

	slog.Info("Collection 'mtguru' created")

}
