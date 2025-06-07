package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jinzhu/copier"
)

type Card struct {
	OracleID      string            `json:"oracle_id"`
	Name          string            `json:"name"`
	ReleasedAt    string            `json:"released_at"`
	ManaCost      string            `json:"mana_cost"`
	Cmc           float64           `json:"cmc"`
	TypeLine      string            `json:"type_line"`
	OracleText    *string           `json:"oracle_text,omitempty"`
	Power         *string           `json:"power,omitempty"`
	Toughness     *string           `json:"toughness,omitempty"`
	ImageURIs     map[string]string `json:"image_uris"`
	ScryfallURI   string            `json:"scryfall_uri"`
	Loyalty       *string           `json:"loyalty,omitempty"`
	ColorIdentity *[]string         `json:"color_identity"`
	Keywords      *[]string         `json:"keywords,omitempty"`
	ProducedMana  *[]string         `json:"produced_mana,omitempty"`
	SetName       *string           `json:"set_name"`
	Rarity        *string           `json:"rarity"`
	GuruPrompt    string
	// HandModifier *string `json:"hand_modifier,omitempty"` 		// too niche
	// LifeModifier *string `json:"life_modifier,omitempty"`  		// too niche
	// Colors        []string `json:"colors,omitempty"`				// omitted because color_identity achieves the same
	// Defence      *string `json:"defense,omitempty"` 				// this is just for sieges and might confuse with toughness
}

type GuruPromptFields struct {
	ManaCost      string
	Cmc           float64
	TypeLine      string
	OracleText    string
	Power         string
	Toughness     string
	Defence       string
	Loyalty       string
	HandModifier  string
	LifeModifier  string
	Colors        []string
	ColorIdentity []string
	Keywords      []string
	ProducedMana  []string
	Rarity        string
}

// Map color letters to words
func mapColors(colors []string) []string {
	transformed := make([]string, len(colors))
	for i, color := range colors {
		switch strings.ToUpper(color) {
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
		default:
			transformed[i] = color
			slog.Error("Could not recognise color.", "color", color)
		}
	}
	return transformed
}

// Validate if card should be added to vector db
func isToken(typeLine string) bool {
	return strings.Contains(typeLine, "Token")
}

func generateGuruPrompt(card Card, client *http.Client) (string, error) {
	// Copy required fields out of card ready for Guru prompt generation
	var guruPromptFields GuruPromptFields
	copier.Copy(&guruPromptFields, &card)

	guruPromptFieldsJson, err := json.Marshal(guruPromptFields)
	if err != nil {
		slog.Error("Failed to marshal prompt fields", "error", err.Error())
		return "", err
	}

	slog.Debug("Requesting prompt.", "promptFields", guruPromptFields)

	req, err := http.NewRequest("POST", os.Getenv("CLOUDFLARE_WORKER_URL")+"/prompt", bytes.NewBuffer(guruPromptFieldsJson))
	if err != nil {
		slog.Error("Failed to build request to prompt API", "error", err.Error())
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to call prompt API", "error", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response from prompt API", "error", err.Error())
		return "", err
	}

	var prompt struct {
		Prompt string `json:"response"`
	}

	json.Unmarshal(body, &prompt)
	return prompt.Prompt, nil
}

// Retrieve Guru prompt, process card for vectorisation, and upload to vector db
func processCard(card Card) {

	// Skip non-relevant cards from data
	if isToken(card.TypeLine) {
		return
	}

	client := &http.Client{}

	// Generate and insert guru prompt if applicable
	if !(card.OracleText == nil && card.Keywords == nil) {
		prompt, err := generateGuruPrompt(card, client)
		if err != nil {
			return
		}
		if prompt != "" {
			card.GuruPrompt = prompt
			slog.Info("Prompt recieved: ", "card", card.Name, "prompt", prompt)
		}
	}

	vectorData := extractCardVectorData(card)

	payload := map[string]interface{}{
		"id": card.OracleID,
		"metadata": map[string]interface{}{
			"name":         card.Name,
			"release_date": card.ReleasedAt,
			"rarity":       card.Rarity,
			"set_name":     *card.SetName,
			"colors":       mapColors(*card.ColorIdentity),
			"image_url":    card.ImageURIs["normal"],
			"scryfall_uri": card.ScryfallURI,
		},
		"data": vectorData,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	req, err := http.NewRequest("POST", os.Getenv("CLOUDFLARE_WORKER_URL")+"/insert", bytes.NewBuffer(payloadJson))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Card uploaded to vector index", "card", card.Name, "response", string(body))
}

// Capture the meaningful data that should be vectorised in a natural language format
func extractCardVectorData(card Card) string {
	vectorText := card.TypeLine + "."
	if card.GuruPrompt != "" {
		vectorText += ": " + card.GuruPrompt
	}
	if card.OracleText != nil {
		vectorText += " " + *card.OracleText
	}
	if card.ProducedMana != nil {
		vectorText += " Produces " + strings.Join(*card.ProducedMana, ",") + "mana."
	}
	if card.Keywords != nil && len(*card.Keywords) > 0 {
		vectorText += " " + strings.Join(*card.Keywords, ",")
	}
	return strings.ReplaceAll(strings.TrimSpace(vectorText), "\n", " ")
}

// Handler for loading cards from JSON filepath
func loadCardsFromJson(path string) ([]Card, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var cards []Card
	if err := json.Unmarshal(byteValue, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

// Execute card processing via func process using concurrent goroutines
func processCardsConcurrently(cards []Card, process func(Card)) {
	const workerCount = 30
	jobs := make(chan Card, len(cards))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for card := range jobs { // Receiving from buffered channel jobs blocks while empty
				process(card)
			}
		}()
	}

	// Enqueue cards to jobs channel
	for _, card := range cards {
		jobs <- card
	}

	close(jobs)

	wg.Wait()
}

// Read the card data from the JSON file path and kick off concurrent processing
func populateIndex() {

	cards, err := loadCardsFromJson(os.Getenv("CARD_JSON_PATH"))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	processCardsConcurrently(cards, processCard)
}
