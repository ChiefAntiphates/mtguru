package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"mtguru/packages/custom_logger"
	"net/http"
	"os"
)

func init() {
	custom_logger.CreateLogger()
}

type Card struct {
	OracleID      string   `json:"oracle_id,omitempty"`
	Name          string   `json:"name,omitempty"`
	ReleasedAt    string   `json:"released_at,omitempty"`
	ManaCost      string   `json:"mana_cost,omitempty"`
	Cmc           float64  `json:"cmc,omitempty"`
	TypeLine      string   `json:"type_line,omitempty"`
	OracleText    string   `json:"oracle_text,omitempty"`
	Power         string   `json:"power,omitempty"`
	Toughness     string   `json:"toughness,omitempty"`
	Defence       string   `json:"defense,omitempty"`
	Loyalty       string   `json:"loyalty,omitempty"`
	HandModifier  string   `json:"hand_modifier,omitempty"`
	LifeModifier  string   `json:"life_modifier,omitempty"`
	Colors        []string `json:"colors,omitempty"`
	ColorIdentity []string `json:"color_identity,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
	ProducedMana  []string `json:"produced_mana,omitempty"`
	SetID         string   `json:"set_id,omitempty"`
	SetName       string   `json:"set_name,omitempty"`
	Rarity        string   `json:"rarity,omitempty"`
	FlavorText    string   `json:"flavor_text,omitempty"`
	Artist        string   `json:"artist,omitempty"`
	ArtistIDs     []string `json:"artist_ids,omitempty"`
	BorderColor   string   `json:"border_color,omitempty"`
	GuruPrompt    string
}

type PromptFields struct {
	ManaCost      string   `json:"mana_cost,omitempty"`
	Cmc           float64  `json:"cmc,omitempty"`
	TypeLine      string   `json:"type_line,omitempty"`
	OracleText    string   `json:"oracle_text,omitempty"`
	Power         string   `json:"power,omitempty"`
	Toughness     string   `json:"toughness,omitempty"`
	Defence       string   `json:"defense,omitempty"`
	Loyalty       string   `json:"loyalty,omitempty"`
	HandModifier  string   `json:"hand_modifier,omitempty"`
	LifeModifier  string   `json:"life_modifier,omitempty"`
	Colors        []string `json:"colors,omitempty"`
	ColorIdentity []string `json:"color_identity,omitempty"`
	Keywords      []string `json:"keywords,omitempty"`
	ProducedMana  []string `json:"produced_mana,omitempty"`
	Rarity        string   `json:"rarity,omitempty"`
}

func (c Card) String() string {
	return c.Name
}

func parseCardsFromFile() []Card {

	// jsonFile, err := os.Open("data/oracle-cards-20250505090231.json")
	jsonFile, err := os.Open("data/example.json")

	if err != nil {
		slog.Error(err.Error())
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var cards []Card
	var promptFields []PromptFields
	json.Unmarshal(byteValue, &cards)
	json.Unmarshal(byteValue, &promptFields)

	for index, value := range promptFields {
		var promptJson, err2 = json.Marshal(value)
		if err2 != nil {
			slog.Error("ah oh dear")
		}

		// Make request to create Guru Prompt
		req, err := http.NewRequest("POST", os.Getenv("CLOUDFLARE_WORKER_URL")+"/prompt", bytes.NewBuffer(promptJson))
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error(err.Error())
		}

		var prompt struct {
			Prompt string `json:"response"`
		}
		json.Unmarshal(body, &prompt)

		cards[index].GuruPrompt = prompt.Prompt
	}

	// Should we consider replacing the colour letters with words

	var uniqueCards []Card
	var duplicateCards []Card
	oracleIDMap := make(map[string]bool)

	for _, card := range cards {
		if oracleIDMap[card.OracleID] {
			duplicateCards = append(duplicateCards, card)
		} else {
			oracleIDMap[card.OracleID] = true
			uniqueCards = append(uniqueCards, card)
		}
	}
	cards = uniqueCards

	// slog.Info("Number of total cards", "count", len(cards))
	slog.Info("Number of unique cards:", "count", len(uniqueCards))
	slog.Info("Number of duplicate cards:", "count", len(duplicateCards))

	// for i := 0; i < len(cards); i++ {
	// 	fmt.Printf("Card ID: "+cards[i].ScryfallID+" Remaining cards: %d\n", len(cards)-i-1)
	// }

	return cards
}

func populateIndex() {
	var cards []Card = parseCardsFromFile()
	slog.Info("Number of total cards", "count", len(cards))

	for _, card := range cards {
		var cardJson, err = json.Marshal(card)
		if err != nil {
			slog.Error(err.Error())
		}
		slog.Info(string(cardJson))
		// TODO: POST payload to vectorise and store DB
	}
}
