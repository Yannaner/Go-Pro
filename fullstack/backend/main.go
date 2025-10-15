package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	openai "github.com/sashabaranov/go-openai"
)

// Card represents a playing card
type Card struct {
	Rank string `json:"rank"`
	Suit string `json:"suit"`
}

// Hand represents a poker hand
type Hand struct {
	Cards   []Card  `json:"cards"`
	WinRate float64 `json:"winrate"`
}

// ChatRequest represents a chat message
type ChatRequest struct {
	Message string `json:"message"`
	Hand    *Hand  `json:"hand,omitempty"`
}

// ChatResponse represents AI response
type ChatResponse struct {
	Response string `json:"response"`
}

func main() {
	// Load environment variables
	godotenv.Load()

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/api/analyze", analyzeHandHandler).Methods("POST")
	router.HandleFunc("/api/chat", chatHandler).Methods("POST")

	// CORS
	handler := cors.Default().Handler(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Poker Analyzer API</h1><p>Server is running!</p>")
}

func analyzeHandHandler(w http.ResponseWriter, r *http.Request) {
	var hand Hand
	err := json.NewDecoder(r.Body).Decode(&hand)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// TODO: Implement MediaPipe hand detection
	// TODO: Implement card recognition
	
	// Calculate win rate based on hand strength
	winRate := calculateWinRate(hand.Cards)
	hand.WinRate = winRate

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hand)
}

// calculateWinRate estimates win rate based on cards
// This is a simplified version - you can expand this with actual poker math
func calculateWinRate(cards []Card) float64 {
	if len(cards) == 0 {
		return 0.0
	}

	// Simple heuristic based on high cards
	// Premium hands get higher win rates
	score := 0.0
	
	for _, card := range cards {
		switch card.Rank {
		case "A":
			score += 14
		case "K":
			score += 13
		case "Q":
			score += 12
		case "J":
			score += 11
		case "10":
			score += 10
		default:
			// For numbered cards, try to parse
			score += 5
		}
	}

	// Check for pairs (same rank)
	if len(cards) == 2 && cards[0].Rank == cards[1].Rank {
		score *= 1.5 // Pocket pairs are stronger
	}

	// Check for suited cards
	if len(cards) == 2 && cards[0].Suit == cards[1].Suit {
		score *= 1.2 // Suited hands have better potential
	}

	// Normalize to 0-1 range (max score for AA is about 42)
	winRate := score / 50.0
	if winRate > 0.95 {
		winRate = 0.95
	}
	if winRate < 0.15 {
		winRate = 0.15
	}

	return winRate
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var chatReq ChatRequest
	err := json.NewDecoder(r.Body).Decode(&chatReq)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not configured", http.StatusInternalServerError)
		return
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	// Build the prompt
	var prompt bytes.Buffer
	prompt.WriteString("You are an expert poker coach. ")
	
	if chatReq.Hand != nil && len(chatReq.Hand.Cards) > 0 {
		prompt.WriteString(fmt.Sprintf("The player has the following hand: "))
		for i, card := range chatReq.Hand.Cards {
			if i > 0 {
				prompt.WriteString(", ")
			}
			prompt.WriteString(fmt.Sprintf("%s of %s", card.Rank, card.Suit))
		}
		prompt.WriteString(fmt.Sprintf(" (Win rate: %.1f%%). ", chatReq.Hand.WinRate*100))
	}
	
	prompt.WriteString(fmt.Sprintf("Question: %s", chatReq.Message))

	// Create chat completion request
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini, // or openai.GPT3Dot5Turbo for cheaper option
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert poker coach and strategist. Provide concise, helpful advice about poker hands and strategy.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt.String(),
				},
			},
			MaxTokens:   500,
			Temperature: 0.7,
		},
	)

	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		http.Error(w, "Failed to get AI response", http.StatusInternalServerError)
		return
	}

	// Extract response
	aiResponse := resp.Choices[0].Message.Content
	
	response := ChatResponse{
		Response: aiResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
