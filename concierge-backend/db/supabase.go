package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Conversation struct {
	ID           string    `json:"id"`
	ProfileID    string    `json:"profile_id"`
	SessionID    string    `json:"session_id"`
	StartedAt    time.Time `json:"started_at"`
	MessageCount int       `json:"message_count"`
	ConsentGiven bool      `json:"consent_given"`
}

type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

func supabaseURL() string {
	return os.Getenv("SUPABASE_URL")
}

func supabaseKey() string {
	return os.Getenv("SUPABASE_KEY")
}

func insert(table string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/rest/v1/%s", supabaseURL(), table)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supabaseKey())
	req.Header.Set("Authorization", "Bearer "+supabaseKey())
	req.Header.Set("Prefer", "return=minimal")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("supabase error: %d", resp.StatusCode)
	}

	return nil
}

func SaveConversation(conv Conversation) error {
	return insert("conversations", conv)
}

func SaveMessage(msg Message) error {
	return insert("messages", msg)
}
