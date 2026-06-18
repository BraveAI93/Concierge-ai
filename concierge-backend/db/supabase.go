package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type Lead struct {
	ID        string    `json:"id"`
	ProfileID string    `json:"profile_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Score     string    `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

type Profile struct {
	ID           string    `json:"id"`
	Slug         string    `json:"slug"`
	Name         string    `json:"name"`
	Business     string    `json:"business"`
	Profession   string    `json:"profession"`
	Location     string    `json:"location"`
	SystemPrompt string    `json:"system_prompt"`
	ProfileData  string    `json:"profile_data"`
	Accent       string    `json:"accent"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

func supabaseURL() string { return os.Getenv("SUPABASE_URL") }
func supabaseKey() string { return os.Getenv("SUPABASE_KEY") }

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
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func SaveConversation(conv Conversation) error { return insert("conversations", conv) }
func SaveMessage(msg Message) error            { return insert("messages", msg) }
func SaveLead(lead Lead) error                 { return insert("leads", lead) }
func SaveProfile(p Profile) error              { return insert("profiles", p) }

func GetProfile(slug string) (*Profile, error) {
	url := fmt.Sprintf("%s/rest/v1/profiles?slug=eq.%s&select=*", supabaseURL(), slug)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", supabaseKey())
	req.Header.Set("Authorization", "Bearer "+supabaseKey())
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var profiles []Profile
	if err := json.Unmarshal(b, &profiles); err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &profiles[0], nil
}
