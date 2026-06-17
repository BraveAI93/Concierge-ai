package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/BraveAI93/concierge-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// ── STRUCTS ──
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	ProfileID    string    `json:"profile_id"`
	SessionID    string    `json:"session_id"`
	Messages     []Message `json:"messages"`
	SystemPrompt string    `json:"system_prompt"`
}

type ChatResponse struct {
	Reply          string `json:"reply"`
	ConversationID string `json:"conversation_id"`
}

type AnthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
	Messages  []Message `json:"messages"`
}

type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func main() {
	godotenv.Load()

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "time": time.Now()})
	})

	// Chat endpoint
	r.POST("/chat", handleChat)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Concierge AI backend running on port %s\n", port)
	r.Run(":" + port)
}

func handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if req.SessionID == "" {
		req.SessionID = uuid.New().String()
	}

	convID := uuid.New().String()

	// Save conversation to Supabase
	conv := db.Conversation{
		ID:           convID,
		ProfileID:    req.ProfileID,
		SessionID:    req.SessionID,
		StartedAt:    time.Now(),
		MessageCount: len(req.Messages),
		ConsentGiven: true,
	}
	if err := db.SaveConversation(conv); err != nil {
		fmt.Printf("Warning: could not save conversation: %v\n", err)
	}

	// Save user message
	userMsg := req.Messages[len(req.Messages)-1]
	if err := db.SaveMessage(db.Message{
		ID:             uuid.New().String(),
		ConversationID: convID,
		Role:           "user",
		Content:        userMsg.Content,
		CreatedAt:      time.Now(),
	}); err != nil {
		fmt.Printf("Warning: could not save message: %v\n", err)
	}

	// Call Anthropic
	reply, err := callAnthropic(req.Messages, req.SystemPrompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI service error"})
		return
	}

	// Save assistant reply
	if err := db.SaveMessage(db.Message{
		ID:             uuid.New().String(),
		ConversationID: convID,
		Role:           "assistant",
		Content:        reply,
		CreatedAt:      time.Now(),
	}); err != nil {
		fmt.Printf("Warning: could not save reply: %v\n", err)
	}

	c.JSON(200, ChatResponse{
		Reply:          reply,
		ConversationID: convID,
	})
}

func callAnthropic(messages []Message, systemPrompt string) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	reqBody := AnthropicRequest{
		Model:     "claude-sonnet-4-6",
		MaxTokens: 1024,
		System:    systemPrompt,
		Messages:  messages,
	}

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(respBody, &anthropicResp); err != nil {
		return "", err
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Anthropic")
	}

	return anthropicResp.Content[0].Text, nil
}
