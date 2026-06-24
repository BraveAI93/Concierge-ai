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
	ID                  string    `json:"id"`
	Slug                string    `json:"slug"`
	Name                string    `json:"name"`
	Business            string    `json:"business"`
	Profession          string    `json:"profession"`
	Location            string    `json:"location"`
	SystemPrompt        string    `json:"system_prompt"`
	ProfileData         string    `json:"profile_data"`
	Accent              string    `json:"accent"`
	Active              bool      `json:"active"`
	Email               string    `json:"email"`
	PasswordHash        string    `json:"password_hash"`
	PasswordSalt        string    `json:"password_salt"`
	GoogleRefreshToken  string    `json:"google_refresh_token"`
	DigestFrequency     string    `json:"digest_frequency"` // "none","biweekly","weekly","bimonthly","monthly"
	DigestLastSent      time.Time `json:"digest_last_sent"`
	StripeAccountID     string    `json:"stripe_account_id"`
	CreatedAt           time.Time `json:"created_at"`
}

// BookingPayment — tracks deposit payments via Stripe
type BookingPayment struct {
	ID            string    `json:"id"`
	ProfileID     string    `json:"profile_id"`
	BookingID     string    `json:"booking_id"`
	ClientEmail   string    `json:"client_email"`
	ClientName    string    `json:"client_name"`
	ServiceName   string    `json:"service_name"`
	AmountPence   int       `json:"amount_pence"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"` // pending|paid|refunded
	StripeSession string    `json:"stripe_session_id"`
	CreatedAt     time.Time `json:"created_at"`
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

// BookingRequest — propose/confirm flow
type BookingRequest struct {
	ID            string    `json:"id"`
	ProfileID     string    `json:"profile_id"`
	ClientName    string    `json:"client_name"`
	ClientEmail   string    `json:"client_email"`
	SessionID     string    `json:"session_id"`
	ServiceName   string    `json:"service_name"`
	PrimarySlot   string    `json:"primary_slot"`
	BackupSlot    string    `json:"backup_slot"`
	Message       string    `json:"message"`
	Status        string    `json:"status"` // "pending","accepted","declined","counter"
	OwnerReply    string    `json:"owner_reply"`
	CounterSlot   string    `json:"counter_slot"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Note — owner notes on clients or general
type Note struct {
	ID        string    `json:"id"`
	ProfileID string    `json:"profile_id"`
	ClientID  string    `json:"client_id"` // email or session_id, empty for personal notes
	Content   string    `json:"content"`
	NoteType  string    `json:"note_type"` // "personal" | "client"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DailyNews — AI-generated news per profile
type DailyNews struct {
	ID        string    `json:"id"`
	ProfileID string    `json:"profile_id"`
	Date      string    `json:"date"` // YYYY-MM-DD
	Items     string    `json:"items"` // JSON array of {title, summary, relevance}
	CreatedAt time.Time `json:"created_at"`
}

// BookingPrefs — per-service booking preferences
type BookingPrefs struct {
	ID             string    `json:"id"`
	ProfileID      string    `json:"profile_id"`
	ServiceName    string    `json:"service_name"`
	SlotDuration   int       `json:"slot_duration"`   // minutes
	BufferMinutes  int       `json:"buffer_minutes"`
	PreferredHours string    `json:"preferred_hours"` // e.g. "09:00-17:00"
	MaxPerDay      int       `json:"max_per_day"`
	UpdatedAt      time.Time `json:"updated_at"`
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

func upsert(table string, data interface{}) error {
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
	req.Header.Set("Prefer", "resolution=merge-duplicates,return=minimal")
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

func getMany(u string) ([]byte, error) {
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
	return io.ReadAll(resp.Body)
}

// ── EXISTING ──────────────────────────────────────────────

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
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var profiles []Profile
	if err := json.Unmarshal(b, &profiles); err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &profiles[0], nil
}

func GetAllActiveProfiles() ([]Profile, error) {
	u := fmt.Sprintf("%s/rest/v1/profiles?active=eq.true&select=*", supabaseURL())
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var profiles []Profile
	if err := json.Unmarshal(b, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func GetProfileByEmail(email string) (*Profile, error) {
	u := fmt.Sprintf("%s/rest/v1/profiles?email=eq.%s&select=*", supabaseURL(), url.QueryEscape(email))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
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
	b, err := getMany(u)
	if err != nil {
		return "", err
	}
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
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var leads []Lead
	if err := json.Unmarshal(b, &leads); err != nil {
		return nil, err
	}
	return leads, nil
}

func GetConversationsByProfile(slug string) ([]Conversation, error) {
	u := fmt.Sprintf("%s/rest/v1/conversations?profile_id=eq.%s&select=*&order=started_at.desc", supabaseURL(), url.QueryEscape(slug))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var convs []Conversation
	if err := json.Unmarshal(b, &convs); err != nil {
		return nil, err
	}
	return convs, nil
}

// ── NEW A3 ────────────────────────────────────────────────

func SaveBookingRequest(br BookingRequest) error { return insert("booking_requests", br) }

func UpdateBookingRequest(id string, data map[string]interface{}) error {
	body, _ := json.Marshal(data)
	u := fmt.Sprintf("%s/rest/v1/booking_requests?id=eq.%s", supabaseURL(), url.QueryEscape(id))
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

func GetBookingRequestsByProfile(slug string) ([]BookingRequest, error) {
	u := fmt.Sprintf("%s/rest/v1/booking_requests?profile_id=eq.%s&select=*&order=created_at.desc", supabaseURL(), url.QueryEscape(slug))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var brs []BookingRequest
	if err := json.Unmarshal(b, &brs); err != nil {
		return nil, err
	}
	return brs, nil
}

func GetBookingRequest(id string) (*BookingRequest, error) {
	u := fmt.Sprintf("%s/rest/v1/booking_requests?id=eq.%s&select=*", supabaseURL(), url.QueryEscape(id))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var brs []BookingRequest
	if err := json.Unmarshal(b, &brs); err != nil {
		return nil, err
	}
	if len(brs) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &brs[0], nil
}

func SaveNote(n Note) error { return insert("notes", n) }

func UpdateNote(id string, content string) error {
	return patch("notes", "id", id, map[string]interface{}{"content": content, "updated_at": time.Now()})
}

func DeleteNote(id string) error {
	u := fmt.Sprintf("%s/rest/v1/notes?id=eq.%s", supabaseURL(), url.QueryEscape(id))
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("apikey", supabaseKey())
	req.Header.Set("Authorization", "Bearer "+supabaseKey())
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func GetNotesByProfile(slug string) ([]Note, error) {
	u := fmt.Sprintf("%s/rest/v1/notes?profile_id=eq.%s&select=*&order=updated_at.desc", supabaseURL(), url.QueryEscape(slug))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var notes []Note
	if err := json.Unmarshal(b, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func GetClientNotesByProfile(slug string) ([]Note, error) {
	u := fmt.Sprintf("%s/rest/v1/notes?profile_id=eq.%s&note_type=eq.client&select=*&order=updated_at.desc", supabaseURL(), url.QueryEscape(slug))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var notes []Note
	if err := json.Unmarshal(b, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func GetClientNote(slug, clientID string) (*Note, error) {
	u := fmt.Sprintf("%s/rest/v1/notes?profile_id=eq.%s&client_id=eq.%s&note_type=eq.client&select=*", supabaseURL(), url.QueryEscape(slug), url.QueryEscape(clientID))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var notes []Note
	if err := json.Unmarshal(b, &notes); err != nil {
		return nil, err
	}
	if len(notes) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &notes[0], nil
}

func SaveDailyNews(dn DailyNews) error { return upsert("daily_news", dn) }

func GetLatestNews(slug string) (*DailyNews, error) {
	u := fmt.Sprintf("%s/rest/v1/daily_news?profile_id=eq.%s&select=*&order=date.desc&limit=1", supabaseURL(), url.QueryEscape(slug))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var news []DailyNews
	if err := json.Unmarshal(b, &news); err != nil {
		return nil, err
	}
	if len(news) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &news[0], nil
}

func SaveBookingPrefs(bp BookingPrefs) error { return upsert("booking_prefs", bp) }

func GetBookingPrefsByProfile(slug string) ([]BookingPrefs, error) {
	u := fmt.Sprintf("%s/rest/v1/booking_prefs?profile_id=eq.%s&select=*", supabaseURL(), url.QueryEscape(slug))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var prefs []BookingPrefs
	if err := json.Unmarshal(b, &prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}

// ── STRIPE PAYMENTS ───────────────────────────────────────────

func SaveBookingPayment(bp BookingPayment) error { return insert("booking_payments", bp) }

func UpdateBookingPaymentBySession(sessionID, status string) error {
	return patch("booking_payments", "stripe_session_id", sessionID, map[string]interface{}{"status": status})
}

func GetPaymentsByProfile(slug string) ([]BookingPayment, error) {
	u := fmt.Sprintf("%s/rest/v1/booking_payments?profile_id=eq.%s&select=*&order=created_at.desc", supabaseURL(), url.QueryEscape(slug))
	b, err := getMany(u)
	if err != nil {
		return nil, err
	}
	var payments []BookingPayment
	if err := json.Unmarshal(b, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}
