package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BraveAI93/concierge-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

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
	Score          string `json:"score"`
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
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Owner-Key")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"service": "Concierge AI Backend", "status": "ok"}) })
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "time": time.Now()}) })
	r.POST("/chat", handleChat)
	r.GET("/analytics", handleAnalytics)
	r.POST("/lead", handleLead)
	r.POST("/alert", handleCreateAlert)
	r.GET("/owner/notifications", handleGetNotifications)
	r.POST("/profile", handleSaveProfile)
	r.GET("/profile/:slug", handleGetProfile)
	r.GET("/check-slug/:slug", handleCheckSlug)
	r.GET("/debug-env", handleDebugEnv)
	r.POST("/consent", handleSaveConsent)
	r.POST("/auth/signup", handleSignup)
	r.POST("/auth/login", handleLogin)
	r.GET("/owner/profile", handleGetOwnerProfile)
	r.PUT("/owner/profile", handleUpdateOwnerProfile)
	r.GET("/owner/leads", handleGetOwnerLeads)
	r.GET("/owner/conversations", handleGetOwnerConversations)
	r.GET("/admin/stats", handleAdminStats)

	// A3: booking requests
	r.POST("/booking-request", handleCreateBookingRequest)
	r.GET("/owner/bookings", handleGetOwnerBookings)
	r.PATCH("/owner/bookings/:id", handleUpdateBooking)

	// A3: notes
	r.POST("/owner/notes", handleCreateNote)
	r.GET("/owner/notes", handleGetNotes)
	r.PATCH("/owner/notes/:id", handleUpdateNote)
	r.DELETE("/owner/notes/:id", handleDeleteNote)

	// A3: news
	r.GET("/owner/news", handleGetNews)
	r.POST("/owner/news", handleSaveNews)

	// A3: digest
	r.GET("/owner/digest", handleGetDigest)

	// A3: booking prefs
	r.GET("/owner/booking-prefs", handleGetBookingPrefs)
	r.POST("/owner/booking-prefs", handleSaveBookingPrefs)

	// Multi-profile support
	r.GET("/owner/profiles", handleGetOwnerProfiles)
	r.POST("/owner/profiles", handleAddOwnerProfile)
	r.POST("/owner/switch-profile", handleSwitchProfile)

	// Stripe payments
	r.POST("/stripe/onboard", handleStripeOnboard)
	r.GET("/stripe/status", handleStripeStatus)
	r.POST("/stripe/checkout", handleCreateCheckout)
	r.POST("/stripe/webhook", handleStripeWebhook)
	r.GET("/owner/payments", handleGetPayments)

	// Media & AI
	r.POST("/media/upload", handleMediaUpload)
	r.POST("/ai/import-services", handleImportServices)

	// Client form pages (public)
	r.GET("/forms/:slug/:formType", handleGetFormInfo)
	r.POST("/forms/:slug/:formType", handleSubmitForm)
	r.GET("/owner/form-submissions", handleGetFormSubmissions)
	r.GET("/flags", handleGetFeatureFlags)
	r.POST("/pa/websearch", handlePAWebSearch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Concierge AI backend running on port %s\n", port)
	r.Run(":" + port)
}

func handleDebugEnv(c *gin.Context) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	c.JSON(200, gin.H{
		"supabase_url_set":    url != "",
		"supabase_url_prefix": safePrefix(url, 20),
		"supabase_key_set":    key != "",
		"supabase_key_len":    len(key),
	})
}

func safePrefix(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}

func scoreLead(messages []Message) string {
	score := 0
	hotWords := []string{"book", "booking", "available", "price", "cost", "how much", "when can", "tomorrow", "this week", "pay", "deposit", "prenot", "disponibil", "prezzo", "quanto", "costo", "budget", "interested", "ready"}
	warmWords := []string{"how", "what", "tell me", "interested", "looking for", "need", "want", "could", "would", "info", "more"}
	allText := ""
	for _, m := range messages {
		if m.Role == "user" {
			allText += strings.ToLower(m.Content) + " "
		}
	}
	for _, w := range hotWords {
		if strings.Contains(allText, w) {
			score += 2
		}
	}
	for _, w := range warmWords {
		if strings.Contains(allText, w) {
			score++
		}
	}
	userMsgs := 0
	for _, m := range messages {
		if m.Role == "user" {
			userMsgs++
		}
	}
	if userMsgs >= 4 {
		score += 2
	}
	if score >= 5 {
		return "hot"
	} else if score >= 2 {
		return "warm"
	}
	return "cold"
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
	conv := db.Conversation{ID: convID, ProfileID: req.ProfileID, SessionID: req.SessionID, StartedAt: time.Now(), MessageCount: len(req.Messages), ConsentGiven: true}
	if err := db.SaveConversation(conv); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}
	userMsg := req.Messages[len(req.Messages)-1]
	if err := db.SaveMessage(db.Message{ID: uuid.New().String(), ConversationID: convID, Role: "user", Content: userMsg.Content, CreatedAt: time.Now()}); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}
	reply, err := callAnthropic(req.Messages, req.SystemPrompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI service error"})
		return
	}
	if err := db.SaveMessage(db.Message{ID: uuid.New().String(), ConversationID: convID, Role: "assistant", Content: reply, CreatedAt: time.Now()}); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}
	score := scoreLead(req.Messages)
	c.JSON(200, ChatResponse{Reply: reply, ConversationID: convID, Score: score})
}

func handleLead(c *gin.Context) {
	var req struct {
		ProfileID string `json:"profile_id"`
		SessionID string `json:"session_id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Score     string `json:"score"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	lead := db.Lead{ID: uuid.New().String(), ProfileID: req.ProfileID, Name: req.Name, Email: req.Email, Score: req.Score, CreatedAt: time.Now()}
	if err := db.SaveLead(lead); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}
	if req.Score == "hot" {
		go sendHotLeadEmail(req.Name, req.Email, req.ProfileID)
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func sendHotLeadEmail(name, email, profileID string) {
	apiKey := os.Getenv("RESEND_API_KEY")
	ownerEmail := os.Getenv("OWNER_EMAIL")
	if apiKey == "" || ownerEmail == "" {
		return
	}
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}
	body := map[string]interface{}{
		"from":    "Concierge AI <" + fromEmail + ">",
		"to":      []string{ownerEmail},
		"subject": "Hot lead from your Concierge!",
		"html":    fmt.Sprintf("<h2>New hot lead!</h2><p><b>Name:</b> %s</p><p><b>Email:</b> %s</p><p><b>Profile:</b> %s</p>", name, email, profileID),
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Warning: Resend hot-lead email failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		rb, _ := io.ReadAll(resp.Body)
		fmt.Printf("Warning: Resend returned %d for hot-lead to %s: %s\n", resp.StatusCode, ownerEmail, string(rb))
	}
}

// ─── OWNER NOTIFICATION EVENTS ─────────────────────────
// V1: durable notification/event record, stored via the existing notes
// table (note_type="alert") since the repo has no migration mechanism for
// a dedicated table. See docs/architecture/BRAIN_SPINE_READINESS_AUDIT.md.

type AlertPayload struct {
	AlertType      string `json:"alert_type"`
	Topic          string `json:"topic"`
	Excerpt        string `json:"excerpt"`
	ConversationID string `json:"conversation_id,omitempty"`
	SessionID      string `json:"session_id,omitempty"`
	EmailStatus    string `json:"email_status"`          // "sent" | "failed" | "disabled_missing_env"
	EmailError     string `json:"email_error,omitempty"` // safe, redacted — never includes secrets
}

// Delivery statuses returned by sendAlertEmail. "disabled_missing_env" is not
// a failure — it means RESEND_API_KEY/OWNER_EMAIL are not configured, so no
// send was attempted.
const (
	EmailStatusSent          = "sent"
	EmailStatusFailed        = "failed"
	EmailStatusDisabledNoEnv = "disabled_missing_env"
)

func handleCreateAlert(c *gin.Context) {
	var req struct {
		ProfileID      string `json:"profile_id"`
		SessionID      string `json:"session_id"`
		ConversationID string `json:"conversation_id"`
		Topic          string `json:"topic"`
		Excerpt        string `json:"excerpt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ProfileID == "" || req.Topic == "" {
		c.JSON(400, gin.H{"error": "profile_id and topic required"})
		return
	}
	excerpt := req.Excerpt
	if len(excerpt) > 300 {
		excerpt = excerpt[:300]
	}

	// Attempt email delivery synchronously (not fire-and-forget) so the
	// caller and the durable record both get the real outcome, not a guess.
	emailStatus, emailErr := sendAlertEmail(req.Topic, excerpt, req.ProfileID)
	emailEnabled := emailStatus != EmailStatusDisabledNoEnv

	payload := AlertPayload{
		AlertType:      "sensitive_topic",
		Topic:          req.Topic,
		Excerpt:        excerpt,
		ConversationID: req.ConversationID,
		SessionID:      req.SessionID,
		EmailStatus:    emailStatus,
		EmailError:     emailErr,
	}
	content, err := json.Marshal(payload)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not encode alert"})
		return
	}
	n := db.Note{
		ID:        uuid.New().String(),
		ProfileID: req.ProfileID,
		ClientID:  req.SessionID,
		Content:   string(content),
		NoteType:  "alert",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.SaveNote(n); err != nil {
		resp := gin.H{
			"error":                    "Could not create notification record",
			"detail":                   err.Error(),
			"dashboard_record_created": false,
			"email_enabled":            emailEnabled,
			"email_status":             emailStatus,
		}
		if emailErr != "" {
			resp["email_error"] = emailErr
		}
		c.JSON(500, resp)
		return
	}
	logAuditEvent(req.ProfileID, "alert_created", "chat_alert_endpoint", "note", n.ID, emailStatus, map[string]interface{}{
		"topic": req.Topic,
	})
	resp := gin.H{
		"status":                   "ok",
		"id":                       n.ID,
		"dashboard_record_created": true,
		"email_enabled":            emailEnabled,
		"email_status":             emailStatus,
	}
	if emailErr != "" {
		resp["email_error"] = emailErr
	}
	c.JSON(200, resp)
}

func handleGetNotifications(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	alerts, err := db.GetAlertsByProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"notifications": []interface{}{}})
		return
	}
	c.JSON(200, gin.H{"notifications": alerts})
}

// sendAlertEmail attempts real, synchronous delivery via Resend and returns
// an explicit status plus a safe (no secrets, no raw response body) error
// message. It never claims "sent" unless Resend accepted the request.
func sendAlertEmail(topic, excerpt, profileID string) (status string, safeErr string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Warning: sendAlertEmail panic: %v\n", r)
			status = EmailStatusFailed
			safeErr = "internal error while sending email"
		}
	}()
	apiKey := os.Getenv("RESEND_API_KEY")
	ownerEmail := os.Getenv("OWNER_EMAIL")
	if apiKey == "" || ownerEmail == "" {
		return EmailStatusDisabledNoEnv, ""
	}
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}
	body := map[string]interface{}{
		"from":    "Concierge AI <" + fromEmail + ">",
		"to":      []string{ownerEmail},
		"subject": "Sensitive topic flagged in your Concierge chat",
		"html":    fmt.Sprintf("<h2>A conversation was flagged</h2><p><b>Topic:</b> %s</p><p><b>Excerpt:</b> %s</p><p><b>Profile:</b> %s</p>", topic, excerpt, profileID),
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("Warning: sendAlertEmail request build failed: %v\n", err)
		return EmailStatusFailed, "could not build email request"
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Warning: Resend alert email failed: %v\n", err)
		return EmailStatusFailed, "could not reach email provider"
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		rb, _ := io.ReadAll(resp.Body)
		fmt.Printf("Warning: Resend returned %d for alert email to %s: %s\n", resp.StatusCode, ownerEmail, string(rb))
		return EmailStatusFailed, fmt.Sprintf("email provider returned status %d", resp.StatusCode)
	}
	return EmailStatusSent, ""
}

func handleSaveProfile(c *gin.Context) {
	var p db.Profile
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	if p.Slug == "" {
		c.JSON(400, gin.H{"error": "Slug required"})
		return
	}
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	if err := db.SaveProfile(p); err != nil {
		c.JSON(500, gin.H{"error": "Could not save profile", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "slug": p.Slug})
}

// publicProfile is the allowlisted public view of a db.Profile. /profile/:slug is
// unauthenticated, and db.Profile carries credentials (password hash/salt),
// connected-account secrets (Google refresh token, Stripe account) and internal
// metadata — none of which may ever reach a public response. Fields are added
// here only when a public page genuinely needs them.
type publicProfile struct {
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Business     string `json:"business"`
	Profession   string `json:"profession"`
	Location     string `json:"location"`
	SystemPrompt string `json:"system_prompt"`
	ProfileData  string `json:"profile_data"`
	Accent       string `json:"accent"`
	Active       bool   `json:"active"`
}

func newPublicProfile(p db.Profile) publicProfile {
	return publicProfile{
		Slug:         p.Slug,
		Name:         p.Name,
		Business:     p.Business,
		Profession:   p.Profession,
		Location:     p.Location,
		SystemPrompt: p.SystemPrompt,
		ProfileData:  p.ProfileData,
		Accent:       p.Accent,
		Active:       p.Active,
	}
}

// Indirection so the public profile route can be tested without a live Supabase backend.
var publicGetProfile = db.GetProfile

func handleGetProfile(c *gin.Context) {
	slug := c.Param("slug")
	p, err := publicGetProfile(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(200, newPublicProfile(*p))
}

func handleCheckSlug(c *gin.Context) {
	slug := c.Param("slug")
	reserved := map[string]bool{"bruno": true, "marco": true, "nour": true, "sofia": true, "alex": true, "custom": true, "admin": true, "api": true, "privacy": true, "index": true, "dashboard": true}
	if reserved[slug] {
		c.JSON(200, gin.H{"available": false, "reason": "reserved"})
		return
	}
	_, err := db.GetProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"available": true})
		return
	}
	c.JSON(200, gin.H{"available": false, "reason": "taken"})
}

func handleAnalytics(c *gin.Context) {
	ownerKey := c.GetHeader("X-Owner-Key")
	validKey := os.Getenv("OWNER_KEY")
	if validKey == "" {
		validKey = "demo123"
	}
	if ownerKey != validKey {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(200, gin.H{"conversations": queryCount("conversations"), "messages": queryCount("messages"), "leads": queryCount("leads"), "bookingIntents": 0, "blocked": 0, "byProfile": map[string]int{}, "topTopics": map[string]int{}})
}

func queryCount(table string) int {
	url := fmt.Sprintf("%s/rest/v1/%s?select=count", os.Getenv("SUPABASE_URL"), table)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("apikey", os.Getenv("SUPABASE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_KEY"))
	req.Header.Set("Prefer", "count=exact")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	contentRange := resp.Header.Get("Content-Range")
	if contentRange != "" {
		parts := strings.Split(contentRange, "/")
		if len(parts) == 2 {
			count := 0
			fmt.Sscanf(parts[1], "%d", &count)
			return count
		}
	}
	return 0
}

func callAnthropic(messages []Message, systemPrompt string) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	reqBody := AnthropicRequest{Model: "claude-sonnet-4-6", MaxTokens: 1024, System: systemPrompt, Messages: messages}
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
		return "", fmt.Errorf("empty response")
	}
	return anthropicResp.Content[0].Text, nil
}

func handleSaveConsent(c *gin.Context) {
	var req struct {
		ProfileID     string   `json:"profile_id"`
		SessionID     string   `json:"session_id"`
		ClientName    string   `json:"client_name"`
		ClientEmail   string   `json:"client_email"`
		FormsAgreed   []string `json:"forms_agreed"`
		Answers       string   `json:"answers"`
		SignatureDate string   `json:"signature_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	formsJSON, _ := json.Marshal(req.FormsAgreed)
	consent := db.Consent{
		ID:            uuid.New().String(),
		ProfileID:     req.ProfileID,
		SessionID:     req.SessionID,
		ClientName:    req.ClientName,
		ClientEmail:   req.ClientEmail,
		FormsAgreed:   string(formsJSON),
		Answers:       req.Answers,
		SignatureDate: req.SignatureDate,
		CreatedAt:     time.Now(),
	}
	if err := db.SaveConsent(consent); err != nil {
		c.JSON(500, gin.H{"error": "Could not save consent", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "id": consent.ID})
}

// ─── AUTH ──────────────────────────────────────────────

func hashPassword(password, salt string) string {
	h := sha256.Sum256([]byte(password + salt + os.Getenv("AUTH_SALT_SECRET")))
	return hex.EncodeToString(h[:])
}

func generateToken() string {
	return uuid.New().String() + uuid.New().String()
}

// Indirection so signup can be tested without a live Supabase backend.
var (
	signupGetProfile    = db.GetProfile
	signupUpdateProfile = db.UpdateProfile
	signupSaveSession   = db.SaveSession
)

func handleSignup(c *gin.Context) {
	var req struct {
		Slug     string `json:"slug"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if req.Slug == "" || req.Email == "" || len(req.Password) < 6 {
		c.JSON(400, gin.H{"error": "Slug, email and password (6+ chars) required"})
		return
	}
	p, err := signupGetProfile(req.Slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found — complete onboarding first"})
		return
	}
	// Credentials already registered: never overwrite them from an unauthenticated
	// signup request. Reject before touching email, hash or salt.
	if strings.TrimSpace(p.PasswordHash) != "" {
		c.JSON(409, gin.H{"error": "Account already exists — please log in"})
		return
	}
	salt := uuid.New().String()
	hashed := hashPassword(req.Password, salt)
	p.Email = req.Email
	p.PasswordHash = hashed
	p.PasswordSalt = salt
	if err := signupUpdateProfile(*p); err != nil {
		c.JSON(500, gin.H{"error": "Could not save account", "detail": err.Error()})
		return
	}
	token := generateToken()
	if err := signupSaveSession(db.Session{Token: token, Slug: p.Slug, CreatedAt: time.Now()}); err != nil {
		c.JSON(500, gin.H{"error": "Could not create session"})
		return
	}
	go sendWelcomeEmail(p.Name, req.Email, p.Slug)
	c.JSON(200, gin.H{"status": "ok", "token": token, "slug": p.Slug})
}

func sendWelcomeEmail(name, email, slug string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Warning: sendWelcomeEmail panic: %v\n", r)
		}
	}()
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" || email == "" {
		return
	}
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}
	publicURL := fmt.Sprintf("https://bravebybruno.com/%s", slug)
	dashboardURL := "https://bravebybruno.com"
	firstName := name
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1.0"/>
<title>Your Concierge is live</title>
</head>
<body style="margin:0;padding:0;background-color:#0a0a0a;font-family:'Helvetica Neue',Helvetica,Arial,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#0a0a0a;min-height:100vh;">
<tr><td align="center" style="padding:48px 16px 40px;">
<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width:560px;">

  <!-- Header -->
  <tr><td align="center" style="padding-bottom:36px;">
    <div style="font-size:13px;letter-spacing:0.26em;text-transform:uppercase;color:#c9a96e;font-weight:600;">✦ THE CONCIERGE</div>
  </td></tr>

  <!-- Card -->
  <tr><td style="background:#111010;border:1px solid rgba(201,169,110,0.18);border-radius:24px;padding:44px 40px 40px;box-shadow:0 0 80px rgba(201,169,110,0.05),0 32px 80px rgba(0,0,0,0.7);">

    <!-- H1 -->
    <h1 style="margin:0 0 10px;font-size:32px;font-weight:700;color:#ffffff;letter-spacing:-0.01em;line-height:1.15;">Welcome, %s!</h1>
    <p style="margin:0 0 36px;font-size:15px;color:rgba(255,255,255,0.38);font-weight:300;line-height:1.6;">Your AI concierge is live and ready for clients.</p>

    <!-- Gold divider -->
    <div style="height:1px;background:linear-gradient(90deg,transparent,rgba(201,169,110,0.3),transparent);margin-bottom:36px;"></div>

    <!-- Primary CTA: Share with clients -->
    <div style="margin-bottom:14px;">
      <a href="%s" style="display:block;background:linear-gradient(135deg,#c9a96e,#8a5a10);color:#0a0a0a;text-decoration:none;padding:15px 28px;border-radius:50px;font-size:14px;font-weight:700;letter-spacing:0.06em;text-align:center;text-transform:uppercase;">Share with clients →</a>
    </div>
    <div style="text-align:center;margin-bottom:36px;">
      <span style="font-size:12px;color:rgba(201,169,110,0.45);word-break:break-all;">%s</span>
    </div>

    <!-- Secondary CTA: Dashboard -->
    <div style="margin-bottom:40px;text-align:center;">
      <a href="%s" style="display:inline-block;background:transparent;border:1px solid rgba(201,169,110,0.3);color:#c9a96e;text-decoration:none;padding:12px 28px;border-radius:50px;font-size:13px;font-weight:500;letter-spacing:0.05em;">Open my dashboard →</a>
    </div>

    <!-- Gold divider -->
    <div style="height:1px;background:linear-gradient(90deg,transparent,rgba(201,169,110,0.18),transparent);margin-bottom:36px;"></div>

    <!-- Next steps -->
    <div style="font-size:10px;font-weight:600;letter-spacing:0.16em;text-transform:uppercase;color:rgba(201,169,110,0.5);margin-bottom:20px;">Your next steps</div>
    <table cellpadding="0" cellspacing="0" border="0" width="100%%">
      <tr><td style="padding:14px 0;border-bottom:1px solid rgba(255,255,255,0.05);">
        <table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
          <td width="36" style="vertical-align:top;padding-top:1px;">
            <div style="width:24px;height:24px;border-radius:50%%;background:rgba(201,169,110,0.12);border:1px solid rgba(201,169,110,0.25);text-align:center;line-height:24px;font-size:10px;font-weight:700;color:#c9a96e;">1</div>
          </td>
          <td>
            <div style="font-size:14px;color:#ffffff;font-weight:500;margin-bottom:2px;">Share your link</div>
            <div style="font-size:12px;color:rgba(255,255,255,0.35);font-weight:300;">Send it to clients, add it to your Instagram bio, put it everywhere.</div>
          </td>
        </tr></table>
      </td></tr>
      <tr><td style="padding:14px 0;border-bottom:1px solid rgba(255,255,255,0.05);">
        <table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
          <td width="36" style="vertical-align:top;padding-top:1px;">
            <div style="width:24px;height:24px;border-radius:50%%;background:rgba(201,169,110,0.12);border:1px solid rgba(201,169,110,0.25);text-align:center;line-height:24px;font-size:10px;font-weight:700;color:#c9a96e;">2</div>
          </td>
          <td>
            <div style="font-size:14px;color:#ffffff;font-weight:500;margin-bottom:2px;">Upload your photos</div>
            <div style="font-size:12px;color:rgba(255,255,255,0.35);font-weight:300;">Add images and videos to your profile — Edit Profile → Upload media.</div>
          </td>
        </tr></table>
      </td></tr>
      <tr><td style="padding:14px 0;">
        <table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
          <td width="36" style="vertical-align:top;padding-top:1px;">
            <div style="width:24px;height:24px;border-radius:50%%;background:rgba(201,169,110,0.12);border:1px solid rgba(201,169,110,0.25);text-align:center;line-height:24px;font-size:10px;font-weight:700;color:#c9a96e;">3</div>
          </td>
          <td>
            <div style="font-size:14px;color:#ffffff;font-weight:500;margin-bottom:2px;">Connect your calendar</div>
            <div style="font-size:12px;color:rgba(255,255,255,0.35);font-weight:300;">Link Google Calendar so your AI can manage bookings automatically.</div>
          </td>
        </tr></table>
      </td></tr>
    </table>

  </td></tr>

  <!-- Instagram footer -->
  <tr><td align="center" style="padding-top:32px;padding-bottom:10px;">
    <a href="https://www.instagram.com/bravebybruno" style="text-decoration:none;display:inline-flex;align-items:center;gap:6px;">
      <span style="font-size:12px;color:rgba(201,169,110,0.5);letter-spacing:0.1em;">@bravebybruno</span>
    </a>
  </td></tr>

  <!-- Powered by -->
  <tr><td align="center" style="padding-bottom:40px;">
    <div style="font-size:11px;letter-spacing:0.14em;text-transform:uppercase;color:#c9a96e;font-weight:500;">Powered by Brave by Bruno</div>
  </td></tr>

</table>
</td></tr>
</table>
</body>
</html>`, firstName, publicURL, publicURL, dashboardURL)
	body := map[string]interface{}{
		"from":    "Brave by Bruno <" + fromEmail + ">",
		"to":      []string{email},
		"subject": "✦ Your Concierge is live, " + firstName,
		"html":    html,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Warning: Resend request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		rb, _ := io.ReadAll(resp.Body)
		fmt.Printf("Warning: Resend returned %d for %s: %s\n", resp.StatusCode, email, string(rb))
	}
}

func handleLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	p, err := db.GetProfileByEmail(req.Email)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}
	hashed := hashPassword(req.Password, p.PasswordSalt)
	if hashed != p.PasswordHash || p.PasswordHash == "" {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}
	token := generateToken()
	if err := db.SaveSession(db.Session{Token: token, Slug: p.Slug, CreatedAt: time.Now()}); err != nil {
		c.JSON(500, gin.H{"error": "Could not create session"})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "token": token, "slug": p.Slug})
}

func authenticateToken(c *gin.Context) (string, bool) {
	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return "", false
	}
	slug, err := db.GetSessionSlug(token)
	if err != nil {
		return "", false
	}
	return slug, true
}

// The owner endpoints below return allowlisted views of db.Profile rather than
// blanking individual fields. Blanking is a denylist: it fails open, so every
// field added to db.Profile in future would be published automatically. These
// responses are authenticated and scoped to the caller's own account, but they
// still land in browser-reachable JSON, so credentials, OAuth tokens, payment
// identifiers and internal metadata are excluded entirely.

// ownerProfileDetail is the allowlisted view returned by GET /owner/profile.
// email and digest_frequency are included because the dashboard genuinely
// renders them (Account section, and the digest-frequency selector).
type ownerProfileDetail struct {
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	Business        string `json:"business"`
	Profession      string `json:"profession"`
	Location        string `json:"location"`
	SystemPrompt    string `json:"system_prompt"`
	ProfileData     string `json:"profile_data"`
	Accent          string `json:"accent"`
	Active          bool   `json:"active"`
	Email           string `json:"email"`
	DigestFrequency string `json:"digest_frequency"`
}

func newOwnerProfileDetail(p db.Profile) ownerProfileDetail {
	return ownerProfileDetail{
		Slug:            p.Slug,
		Name:            p.Name,
		Business:        p.Business,
		Profession:      p.Profession,
		Location:        p.Location,
		SystemPrompt:    p.SystemPrompt,
		ProfileData:     p.ProfileData,
		Accent:          p.Accent,
		Active:          p.Active,
		Email:           p.Email,
		DigestFrequency: p.DigestFrequency,
	}
}

// ownerProfileSummary is the allowlisted per-profile view returned by
// GET /owner/profiles, which only backs the dashboard's profile switcher.
type ownerProfileSummary struct {
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Business   string `json:"business"`
	Profession string `json:"profession"`
}

func newOwnerProfileSummary(p db.Profile) ownerProfileSummary {
	return ownerProfileSummary{
		Slug:       p.Slug,
		Name:       p.Name,
		Business:   p.Business,
		Profession: p.Profession,
	}
}

// Indirection so the two owner profile read endpoints can be tested without a
// live Supabase backend. authenticateToken itself is unchanged.
var (
	ownerAuth               = authenticateToken
	ownerGetProfile         = db.GetProfile
	ownerGetProfilesByEmail = db.GetProfilesByEmail
)

func handleGetOwnerProfile(c *gin.Context) {
	slug, ok := ownerAuth(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	p, err := ownerGetProfile(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(200, newOwnerProfileDetail(*p))
}

func handleUpdateOwnerProfile(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	var updates db.Profile
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	existing, err := db.GetProfile(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}
	updates.ID = existing.ID
	updates.Slug = existing.Slug
	updates.Email = existing.Email
	updates.PasswordHash = existing.PasswordHash
	updates.PasswordSalt = existing.PasswordSalt
	updates.CreatedAt = existing.CreatedAt
	if err := db.UpdateProfile(updates); err != nil {
		c.JSON(500, gin.H{"error": "Could not update profile", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func handleGetOwnerLeads(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	leads, err := db.GetLeadsByProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"leads": []interface{}{}})
		return
	}
	c.JSON(200, gin.H{"leads": leads})
}

func handleGetOwnerConversations(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	convs, err := db.GetConversationsByProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"conversations": []interface{}{}, "count": 0})
		return
	}
	c.JSON(200, gin.H{"conversations": convs, "count": len(convs)})
}

// ─── ADMIN (Bruno only) ────────────────────────────────

func handleAdminStats(c *gin.Context) {
	adminKey := c.GetHeader("X-Admin-Key")
	validKey := os.Getenv("ADMIN_KEY")
	if validKey == "" || adminKey != validKey {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	totalUsers := queryCount("profiles")
	totalConversations := queryCount("conversations")
	totalMessages := queryCount("messages")
	totalLeads := queryCount("leads")
	estCostGBP := float64(totalMessages) * 0.003
	c.JSON(200, gin.H{
		"total_users":         totalUsers,
		"total_conversations": totalConversations,
		"total_messages":      totalMessages,
		"total_leads":         totalLeads,
		"estimated_cost_gbp":  fmt.Sprintf("%.2f", estCostGBP),
	})
}

// ─── A3: BOOKING REQUESTS ──────────────────────────────

func handleCreateBookingRequest(c *gin.Context) {
	var req struct {
		ProfileID   string `json:"profile_id"`
		ClientName  string `json:"client_name"`
		ClientEmail string `json:"client_email"`
		SessionID   string `json:"session_id"`
		ServiceName string `json:"service_name"`
		PrimarySlot string `json:"primary_slot"`
		BackupSlot  string `json:"backup_slot"`
		Message     string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	br := db.BookingRequest{
		ID:          uuid.New().String(),
		ProfileID:   req.ProfileID,
		ClientName:  req.ClientName,
		ClientEmail: req.ClientEmail,
		SessionID:   req.SessionID,
		ServiceName: req.ServiceName,
		PrimarySlot: req.PrimarySlot,
		BackupSlot:  req.BackupSlot,
		Message:     req.Message,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := db.SaveBookingRequest(br); err != nil {
		c.JSON(500, gin.H{"error": "Could not save booking request", "detail": err.Error()})
		return
	}
	logAuditEvent(req.ProfileID, "booking_request_created", "booking_endpoint", "booking_request", br.ID, "ok", map[string]interface{}{
		"service_name": req.ServiceName,
	})
	c.JSON(200, gin.H{"status": "ok", "id": br.ID})
}

func handleGetOwnerBookings(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	brs, err := db.GetBookingRequestsByProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"bookings": []interface{}{}})
		return
	}
	c.JSON(200, gin.H{"bookings": brs})
}

func handleUpdateBooking(c *gin.Context) {
	_, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	var req struct {
		Status      string `json:"status"`
		OwnerReply  string `json:"owner_reply"`
		CounterSlot string `json:"counter_slot"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	data := map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	}
	if req.OwnerReply != "" {
		data["owner_reply"] = req.OwnerReply
	}
	if req.CounterSlot != "" {
		data["counter_slot"] = req.CounterSlot
	}
	if err := db.UpdateBookingRequest(id, data); err != nil {
		c.JSON(500, gin.H{"error": "Could not update booking", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

// ─── A3: NOTES ─────────────────────────────────────────

func handleCreateNote(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	var req struct {
		ClientID string `json:"client_id"`
		Content  string `json:"content"`
		NoteType string `json:"note_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	noteType := req.NoteType
	if noteType == "" {
		noteType = "personal"
	}
	n := db.Note{
		ID:        uuid.New().String(),
		ProfileID: slug,
		ClientID:  req.ClientID,
		Content:   req.Content,
		NoteType:  noteType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.SaveNote(n); err != nil {
		c.JSON(500, gin.H{"error": "Could not save note", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "id": n.ID})
}

func handleGetNotes(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	notes, err := db.GetNotesByProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"notes": []interface{}{}})
		return
	}
	c.JSON(200, gin.H{"notes": notes})
}

func handleUpdateNote(c *gin.Context) {
	_, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if err := db.UpdateNote(id, req.Content); err != nil {
		c.JSON(500, gin.H{"error": "Could not update note", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func handleDeleteNote(c *gin.Context) {
	_, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	if err := db.DeleteNote(id); err != nil {
		c.JSON(500, gin.H{"error": "Could not delete note", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

// ─── A3: NEWS — GET returns 200 with empty payload when no rows exist ──────────────────────────────────────────

func handleGetNews(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	news, err := db.GetLatestNews(slug)
	if err != nil {
		c.JSON(200, gin.H{"items": nil, "date": ""})
		return
	}
	c.JSON(200, news)
}

func handleSaveNews(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Warning: handleSaveNews panic: %v\n", r)
		}
	}()
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	var req struct {
		Date  string `json:"date"`
		Items string `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}
	dn := db.DailyNews{
		ID:        uuid.New().String(),
		ProfileID: slug,
		Date:      req.Date,
		Items:     req.Items,
		CreatedAt: time.Now(),
	}
	if err := db.SaveDailyNews(dn); err != nil {
		c.JSON(500, gin.H{"error": "Could not save news", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

// ─── A3: DIGEST ────────────────────────────────────────

func handleGetDigest(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	leads, _ := db.GetLeadsByProfile(slug)
	convs, _ := db.GetConversationsByProfile(slug)
	news, _ := db.GetLatestNews(slug)
	notes, _ := db.GetNotesByProfile(slug)

	hotLeads := 0
	for _, l := range leads {
		if l.Score == "hot" {
			hotLeads++
		}
	}

	digest := gin.H{
		"total_leads":         len(leads),
		"hot_leads":           hotLeads,
		"total_conversations": len(convs),
		"total_notes":         len(notes),
		"latest_news_date":    "",
	}
	if news != nil {
		digest["latest_news_date"] = news.Date
	}
	c.JSON(200, digest)
}

// ─── A3: BOOKING PREFS ─────────────────────────────────

func handleGetBookingPrefs(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	prefs, err := db.GetBookingPrefsByProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"prefs": []interface{}{}})
		return
	}
	c.JSON(200, gin.H{"prefs": prefs})
}

func handleSaveBookingPrefs(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	var req db.BookingPrefs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	req.ID = uuid.New().String()
	req.ProfileID = slug
	req.UpdatedAt = time.Now()
	if err := db.SaveBookingPrefs(req); err != nil {
		c.JSON(500, gin.H{"error": "Could not save booking prefs", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

// ─── MULTI-PROFILE ─────────────────────────────────

func handleGetOwnerProfiles(c *gin.Context) {
	slug, ok := ownerAuth(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	owner, err := ownerGetProfile(slug)
	if err != nil || owner.Email == "" {
		c.JSON(200, gin.H{"profiles": []interface{}{}})
		return
	}
	profiles, err := ownerGetProfilesByEmail(owner.Email)
	if err != nil {
		c.JSON(200, gin.H{"profiles": []interface{}{}})
		return
	}
	summaries := make([]ownerProfileSummary, 0, len(profiles))
	for _, p := range profiles {
		summaries = append(summaries, newOwnerProfileSummary(p))
	}
	c.JSON(200, gin.H{"profiles": summaries})
}

func handleAddOwnerProfile(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	owner, err := db.GetProfile(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Owner profile not found"})
		return
	}
	var p db.Profile
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}
	if p.Slug == "" {
		c.JSON(400, gin.H{"error": "Slug required"})
		return
	}
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.Email = owner.Email
	p.PasswordHash = owner.PasswordHash
	p.PasswordSalt = owner.PasswordSalt
	if err := db.SaveProfile(p); err != nil {
		c.JSON(500, gin.H{"error": "Could not save profile", "detail": err.Error()})
		return
	}
	token := generateToken()
	if err := db.SaveSession(db.Session{Token: token, Slug: p.Slug, CreatedAt: time.Now()}); err != nil {
		c.JSON(500, gin.H{"error": "Could not create session"})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "slug": p.Slug, "token": token})
}

func handleSwitchProfile(c *gin.Context) {
	currentSlug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	owner, err := db.GetProfile(currentSlug)
	if err != nil || owner.Email == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	var req struct {
		Slug string `json:"slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Slug == "" {
		c.JSON(400, gin.H{"error": "Slug required"})
		return
	}
	target, err := db.GetProfile(req.Slug)
	if err != nil || target.Email != owner.Email {
		c.JSON(403, gin.H{"error": "Profile not linked to your account"})
		return
	}
	token := generateToken()
	if err := db.SaveSession(db.Session{Token: token, Slug: target.Slug, CreatedAt: time.Now()}); err != nil {
		c.JSON(500, gin.H{"error": "Could not create session"})
		return
	}
	c.JSON(200, gin.H{"token": token, "slug": target.Slug})
}

// ─── MEDIA UPLOAD ───────────────────────────────

func handleMediaUpload(c *gin.Context) {
	// Optional auth — if present, use slug path; otherwise use temp UUID path
	pathPrefix := "temp/" + uuid.New().String()
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if s, err := db.GetSessionSlug(token); err == nil && s != "" {
			pathPrefix = s
		}
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	if header.Size > 50*1024*1024 {
		c.JSON(400, gin.H{"error": "File too large (max 50MB)"})
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not read file"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".bin"
	}
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = header.Header.Get("Content-Type")
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	storagePath := fmt.Sprintf("%s/%s%s", pathPrefix, uuid.New().String(), ext)
	supabaseBase := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	uploadURL := fmt.Sprintf("%s/storage/v1/object/media/%s", supabaseBase, storagePath)

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileData))
	if err != nil {
		c.JSON(500, gin.H{"error": "Upload request failed"})
		return
	}
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Upload failed: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Storage error %d: %s", resp.StatusCode, string(b))})
		return
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/media/%s", supabaseBase, storagePath)
	c.JSON(200, gin.H{"url": publicURL})
}

// ─── AI SERVICE IMPORT FROM SCREENSHOT ─────────

func callAnthropicVision(imageBase64, mimeType, prompt string) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	body := map[string]interface{}{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 2048,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "image",
						"source": map[string]interface{}{
							"type":       "base64",
							"media_type": mimeType,
							"data":       imageBase64,
						},
					},
					{
						"type": "text",
						"text": prompt,
					},
				},
			},
		},
	}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(respBody, &anthropicResp); err != nil {
		return "", fmt.Errorf("parse error: %s", string(respBody))
	}
	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return anthropicResp.Content[0].Text, nil
}

func handleImportServices(c *gin.Context) {
	_, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		ImageBase64 string `json:"image_base64"`
		MimeType    string `json:"mime_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ImageBase64 == "" {
		c.JSON(400, gin.H{"error": "image_base64 required"})
		return
	}
	if req.MimeType == "" {
		req.MimeType = "image/jpeg"
	}

	// Validate it's actually base64
	if _, err := base64.StdEncoding.DecodeString(req.ImageBase64[:min(64, len(req.ImageBase64))]); err != nil {
		c.JSON(400, gin.H{"error": "Invalid base64 image data"})
		return
	}

	text, err := callAnthropicVision(
		req.ImageBase64,
		req.MimeType,
		"Extract all services from this image. Return a JSON array where each element has these fields: name (string), duration_minutes (number or null), price (number or null), currency (string, default \"£\"), description (string). Return ONLY valid JSON array, no other text, no markdown.",
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI analysis failed: " + err.Error()})
		return
	}

	// Strip markdown code fences if present
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var services []map[string]interface{}
	if err := json.Unmarshal([]byte(text), &services); err != nil {
		c.JSON(500, gin.H{"error": "Could not parse AI response", "raw": text})
		return
	}
	c.JSON(200, gin.H{"services": services})
}

// ─── CLIENT FORM PAGES ─────────────────────────────

func handleGetFormInfo(c *gin.Context) {
	slug := c.Param("slug")
	p, err := db.GetProfile(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(200, gin.H{
		"name":       p.Name,
		"profession": p.Profession,
		"accent":     p.Accent,
		"slug":       p.Slug,
	})
}

func handleSubmitForm(c *gin.Context) {
	slug := c.Param("slug")
	formType := c.Param("formType")

	var req struct {
		ClientName  string `json:"client_name"`
		ClientEmail string `json:"client_email"`
		Responses   string `json:"responses"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ClientName == "" {
		c.JSON(400, gin.H{"error": "client_name required"})
		return
	}

	sub := db.FormSubmission{
		ID:          uuid.New().String(),
		ProfileSlug: slug,
		FormType:    formType,
		ClientName:  req.ClientName,
		ClientEmail: req.ClientEmail,
		Responses:   req.Responses,
		SubmittedAt: time.Now(),
	}
	if err := db.SaveFormSubmission(sub); err != nil {
		c.JSON(500, gin.H{"error": "Could not save submission", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "id": sub.ID})
}

func handleGetFormSubmissions(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	subs, err := db.GetFormSubmissionsByProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"submissions": []interface{}{}})
		return
	}
	c.JSON(200, gin.H{"submissions": subs})
}

func handleGetFeatureFlags(c *gin.Context) {
	flags, err := db.GetFeatureFlags()
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not load feature flags", "detail": err.Error()})
		return
	}
	c.JSON(200, gin.H{"flags": flags})
}

// ─── OWNER PA WEB SEARCH (Sprint 01 Day 6) ──────────────────
// Owner-only. Never used by the public /[slug] concierge chat, which
// has its own separate prompt-building path (buildPrompt in the
// frontend) and never calls this endpoint.
//
// Gate order is deliberate: (1) feature flag must be ACTIVE_PRIVATE,
// (2) valid owner session token, (3) provider call. If either of the
// first two fails, the provider is never touched and no key is spent.

type WebSearchResult struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Source    string `json:"source,omitempty"`
	Snippet   string `json:"snippet,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

func featureFlagState(name string) string {
	flags, err := db.GetFeatureFlags()
	if err != nil {
		return "UNKNOWN"
	}
	for _, f := range flags {
		if f.Name == name {
			return f.State
		}
	}
	return "UNKNOWN"
}

func handlePAWebSearch(c *gin.Context) {
	if featureFlagState("pa_web_search") != "ACTIVE_PRIVATE" {
		resp := gin.H{
			"connected": false,
			"reason":    "pa_web_search is not ACTIVE_PRIVATE yet",
			"results":   []WebSearchResult{},
		}
		if !braveSearchConfigured() {
			resp["setup_needed"] = "BRAVE_SEARCH_API_KEY"
		}
		c.JSON(200, resp)
		return
	}

	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Query     string `json:"query"`
		ProfileID string `json:"profile_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Query == "" {
		c.JSON(400, gin.H{"error": "query required"})
		return
	}
	profileID := req.ProfileID
	if profileID == "" {
		profileID = slug
	}

	start := time.Now()
	logAuditEvent(profileID, "pa_web_search_requested", "pa_web_search_endpoint", "web_search", slug, "ok", map[string]interface{}{
		"provider": "brave_search",
	})

	results, err := braveWebSearch(req.Query)
	latencyMs := time.Since(start).Milliseconds()

	if err != nil {
		logAuditEvent(profileID, "pa_web_search_failed", "pa_web_search_endpoint", "web_search", slug, "failed", map[string]interface{}{
			"provider":   "brave_search",
			"latency_ms": latencyMs,
			"error_code": webSearchErrorCode(err),
		})
		c.JSON(200, gin.H{
			"connected": false,
			"reason":    "search provider error",
			"results":   []WebSearchResult{},
		})
		return
	}

	logAuditEvent(profileID, "pa_web_search_completed", "pa_web_search_endpoint", "web_search", slug, "ok", map[string]interface{}{
		"provider":     "brave_search",
		"result_count": len(results),
		"latency_ms":   latencyMs,
	})

	c.JSON(200, gin.H{"connected": true, "results": results})
}

// braveSearchConfigured reports whether a provider key is present, without
// making a network call. DORMANT_BUILT posture: the adapter below never
// fabricates a result from model knowledge when the key is absent.
func braveSearchConfigured() bool {
	return os.Getenv("BRAVE_SEARCH_API_KEY") != ""
}

func braveWebSearch(query string) ([]WebSearchResult, error) {
	apiKey := os.Getenv("BRAVE_SEARCH_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("BRAVE_SEARCH_API_KEY not configured")
	}
	baseURL := os.Getenv("BRAVE_SEARCH_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.search.brave.com/res/v1/web/search"
	}
	req, err := http.NewRequest("GET", baseURL+"?q="+url.QueryEscape(query), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Subscription-Token", apiKey)
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("brave search request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("brave search returned status %d", resp.StatusCode)
	}
	var parsed struct {
		Web struct {
			Results []struct {
				Title       string `json:"title"`
				URL         string `json:"url"`
				Description string `json:"description"`
				Age         string `json:"age"`
			} `json:"results"`
		} `json:"web"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("brave search returned an unparseable response")
	}
	results := make([]WebSearchResult, 0, min(len(parsed.Web.Results), 5))
	for _, r := range parsed.Web.Results {
		if len(results) >= 5 {
			break
		}
		results = append(results, WebSearchResult{
			Title: r.Title, URL: r.URL, Source: "Brave Search", Snippet: r.Description, Timestamp: r.Age,
		})
	}
	return results, nil
}

// webSearchErrorCode classifies an error into a safe, low-risk code for
// audit_events — never the raw error string, which can embed the query
// (Go's http client errors often include the full request URL).
func webSearchErrorCode(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "not configured"):
		return "missing_api_key"
	case strings.Contains(msg, "unparseable"):
		return "invalid_response"
	case strings.Contains(msg, "returned status"):
		return "provider_http_error"
	default:
		return "request_error"
	}
}

// logAuditEvent is best-effort and must never block or fail the caller's
// primary action — a failed audit write is logged server-side, not surfaced.
// metadata must only contain low-risk structural fields (topic, service_name,
// etc.) — never raw messages, contact details, form contents, or payment data.
func logAuditEvent(profileID, eventType, source, entityType, entityID, status string, metadata map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Warning: logAuditEvent panic: %v\n", r)
		}
	}()
	metaBytes, err := json.Marshal(metadata)
	if err != nil {
		metaBytes = []byte("{}")
	}
	ev := db.AuditEvent{
		ID:           uuid.New().String(),
		ProfileID:    profileID,
		EventType:    eventType,
		Source:       source,
		EntityType:   entityType,
		EntityID:     entityID,
		Status:       status,
		MetadataJSON: metaBytes,
		CreatedAt:    time.Now(),
	}
	if err := db.SaveAuditEvent(ev); err != nil {
		fmt.Printf("Warning: audit event %s (%s) failed to record: %v\n", eventType, entityID, err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
