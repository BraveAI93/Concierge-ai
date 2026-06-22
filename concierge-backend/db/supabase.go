package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	PasswordSalt string    `json:"password_salt"`
	CreatedAt    time.Time `json:"created_at"`
}

type Consent struct {
	ID            string    `json:"id"`
	ProfileID     string    `json:"profile_id"`
	SessionID     string    `json:"session_id"`
	ClientName    string    `json:"client_name"`
	ClientEmail   string    `json:"client_email"`
	FormsAgreed   string    `json:"forms_agreed"`
	Answers       string    `json:"answers"`
	SignatureDate string    `json:"signature_date"`
	CreatedAt     time.Time `json:"created_at"`
}

type Session struct {
	Token     string    `json:"token"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

func supabaseURL() string { return os.Getenv("SUPABASE_URL") }
func supabaseKey() string { return os.Getenv("SUPABASE_KEY") }

func insert(table string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	u := fmt.Sprintf("%s/rest/v1/%s", supabaseURL(), table)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(body))
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

func patch(table, filterCol, filterVal string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	u := fmt.Sprintf("%s/rest/v1/%s?%s=eq.%s", supabaseURL(), table, filterCol, url.QueryEscape(filterVal))
	req, err := http.NewRequest("PATCH", u, bytes.NewBuffer(body))
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
func SaveConsent(con Consent) error            { return insert("consents", con) }
func SaveSession(s Session) error              { return insert("sessions", s) }

func UpdateProfile(p Profile) error {
	return patch("profiles", "slug", p.Slug, p)
}

func GetProfile(slug string) (*Profile, error) {
	u := fmt.Sprintf("%s/rest/v1/profiles?slug=eq.%s&select=*", supabaseURL(), url.QueryEscape(slug))
	req, err := http.NewRequest("GET", u, nil)
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

func GetProfileByEmail(email string) (*Profile, error) {
	u := fmt.Sprintf("%s/rest/v1/profiles?email=eq.%s&select=*", supabaseURL(), url.QueryEscape(email))
	req, err := http.NewRequest("GET", u, nil)
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

func GetSessionSlug(token string) (string, error) {
	u := fmt.Sprintf("%s/rest/v1/sessions?token=eq.%s&select=*", supabaseURL(), url.QueryEscape(token))
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("apikey", supabaseKey())
	req.Header.Set("Authorization", "Bearer "+supabaseKey())
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var sessions []Session
	if err := json.Unmarshal(b, &sessions); err != nil {
		return "", err
	}
	if len(sessions) == 0 {
		return "", fmt.Errorf("session not found")
	}
	return sessions[0].Slug, nil
}

func GetLeadsByProfile(slug string) ([]Lead, error) {
	u := fmt.Sprintf("%s/rest/v1/leads?profile_id=eq.%s&select=*&order=created_at.desc", supabaseURL(), url.QueryEscape(slug))
	req, err := http.NewRequest("GET", u, nil)
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
	var leads []Lead
	if err := json.Unmarshal(b, &leads); err != nil {
		return nil, err
	}
	return leads, nil
}

func GetConversationsByProfile(slug string) ([]Conversation, error) {
	u := fmt.Sprintf("%s/rest/v1/conversations?profile_id=eq.%s&select=*&order=started_at.desc", supabaseURL(), url.QueryEscape(slug))
	req, err := http.NewRequest("GET", u, nil)
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
	var convs []Conversation
	if err := json.Unmarshal(b, &convs); err != nil {
		return nil, err
	}
	return convs, nil
}
