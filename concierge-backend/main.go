package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	// Stripe payments
	r.POST("/stripe/onboard", handleStripeOnboard)
	r.GET("/stripe/status", handleStripeStatus)
	r.POST("/stripe/checkout", handleCreateCheckout)
	r.POST("/stripe/webhook", handleStripeWebhook)
	r.GET("/owner/payments", handleGetPayments)

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

func handleGetProfile(c *gin.Context) {
	slug := c.Param("slug")
	p, err := db.GetProfile(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(200, p)
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
	p, err := db.GetProfile(req.Slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found — complete onboarding first"})
		return
	}
	salt := uuid.New().String()
	hashed := hashPassword(req.Password, salt)
	p.Email = req.Email
	p.PasswordHash = hashed
	p.PasswordSalt = salt
	if err := db.UpdateProfile(*p); err != nil {
		c.JSON(500, gin.H{"error": "Could not save account", "detail": err.Error()})
		return
	}
	token := generateToken()
	if err := db.SaveSession(db.Session{Token: token, Slug: p.Slug, CreatedAt: time.Now()}); err != nil {
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
	publicURL := fmt.Sprintf("https://concierge-ai-gamma.vercel.app/%s", slug)
	dashboardHint := "https://concierge-ai-gamma.vercel.app (tap 'My Profile Login')"
	html := fmt.Sprintf(`
		<h2>Welcome to Concierge AI, %s!</h2>
		<p>Your AI concierge is live. Here's everything you need:</p>
		<p><strong>Your public link (share this with clients):</strong><br>
		<a href="%s">%s</a></p>
		<p><strong>Your dashboard (manage your profile, see leads):</strong><br>
		%s<br>
		Login with this email and the password you chose during setup.</p>
		<p>Save this email — it has everything you need to find your way back.</p>
		<p>— Concierge AI</p>
	`, name, publicURL, publicURL, dashboardHint)
	body := map[string]interface{}{
		"from":    "Concierge AI <" + fromEmail + ">",
		"to":      []string{email},
		"subject": "Your Concierge AI is live — save this email",
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

func handleGetOwnerProfile(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	p, err := db.GetProfile(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found"})
		return
	}
	p.PasswordHash = ""
	p.PasswordSalt = ""
	c.JSON(200, p)
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
