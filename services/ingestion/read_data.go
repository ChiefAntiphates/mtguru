package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mtguru/packages/custom_logger"
	"net/http"
	"os"
	"strings"
)

func init() {
	custom_logger.CreateLogger()
}

// Map colour letters to words
func mapColors(colors []string) []string {
	transformed := make([]string, len(colors))
	for i, color := range colors {
		switch color {
		case "B":
			transformed[i] = "Black"
		case "U":
			transformed[i] = "Blue"
		case "G":
			transformed[i] = "Green"
		case "R":
			transformed[i] = "Red"
		case "W":
			transformed[i] = "White"
		}
	}

	return transformed
}

type Card struct {
	OracleID   string  `json:"oracle_id"`
	Name       string  `json:"name"`
	ReleasedAt string  `json:"released_at"`
	ManaCost   string  `json:"mana_cost"`
	Cmc        float64 `json:"cmc"`
	TypeLine   string  `json:"type_line"`
	OracleText string  `json:"oracle_text"`
	Power      *string `json:"power,omitempty"`
	Toughness  *string `json:"toughness,omitempty"`
	// Defence      *string `json:"defense,omitempty"` // this is just for sieges and might confuse with toughness
	Loyalty *string `json:"loyalty,omitempty"`
	// HandModifier *string `json:"hand_modifier,omitempty"` //very niche
	// LifeModifier *string `json:"life_modifier,omitempty"`  //very niche
	// Colors        []string `json:"colors,omitempty"`
	ColorIdentity *[]string `json:"color_identity"`
	Keywords      *[]string `json:"keywords,omitempty"`
	ProducedMana  *[]string `json:"produced_mana,omitempty"`
	// SetID         string   `json:"set_id,omitempty"`
	SetName *string `json:"set_name"`
	Rarity  *string `json:"rarity"`
	// FlavorText    string   `json:"flavor_text,omitempty"`
	// Artist      string   `json:"artist,omitempty"`
	// ArtistIDs   []string `json:"artist_ids,omitempty"`
	// BorderColor string   `json:"border_color,omitempty"`
	GuruPrompt string
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

	// The double unmarshal here doesn't feel very performant
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

func extractCardVectorData(card Card) string {
	// vectorData := make(map[string]interface{})
	vectorText := card.TypeLine + ": " + card.GuruPrompt + " " + card.OracleText

	// vectorData["mana_cost"] = card.ManaCost
	// vectorData["cmc"] = card.Cmc

	// if card.Power != nil {
	// 	vectorData["power"] = *card.Power
	// }
	// if card.Toughness != nil {
	// 	vectorData["toughness"] = *card.Toughness
	// }
	// if card.Loyalty != nil {
	// 	vectorData["loyalty"] = *card.Loyalty
	// }

	if card.ProducedMana != nil {
		vectorText += " Produces " + strings.Join(*card.ProducedMana, ",") + "mana."
	}
	if card.Keywords != nil && len(*card.Keywords) > 0 {
		vectorText += " " + strings.Join(*card.Keywords, ",")
	}
	return strings.ReplaceAll(strings.TrimSpace(vectorText), "\n", " ")
}

func populateIndex() {
	var cards []Card = parseCardsFromFile()
	slog.Info("Number of total cards", "count", len(cards))

	for _, card := range cards {
		// TODO: POST payload to vectorise and store DB

		vectorData := extractCardVectorData(card)
		fmt.Println(vectorData)

		payload := map[string]interface{}{
			"id": card.OracleID,
			"metadata": map[string]interface{}{
				"name":         card.Name,
				"release_date": card.ReleasedAt,
				"rarity":       card.Rarity,
				"set_name":     *card.SetName,
				"colors":       mapColors(*card.ColorIdentity),
			},
			"data": vectorData,
		}

		var payloadJson, err = json.Marshal(payload)
		if err != nil {
			slog.Error(err.Error())
		}

		slog.Info(string(payloadJson))

		// Make request to create Guru Insert
		req, err := http.NewRequest("POST", os.Getenv("CLOUDFLARE_WORKER_URL")+"/insert", bytes.NewBuffer(payloadJson))
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
		slog.Info(string(body))
	}

}
