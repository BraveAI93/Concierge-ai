package main

// ══════════════════════════════════════════════════════════════
// STRIPE INTEGRATION — The Concierge by Brave by Bruno
// Fee structure:
//   - Stripe processing: 1.4% + 20p (UK cards) — charged to professional
//   - Platform fee: 0% month 1, then 1% from month 2
//   - Client NEVER sees extra fees
// ══════════════════════════════════════════════════════════════

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BraveAI93/concierge-backend/db"
	"github.com/gin-gonic/gin"
)

// ── Stripe helpers ────────────────────────────────────────────

func stripePost(endpoint string, params url.Values) ([]byte, error) {
	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY not set")
	}
	req, err := http.NewRequest("POST", "https://api.stripe.com/v1/"+endpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("stripe error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func stripeGet(endpoint string) ([]byte, error) {
	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY not set")
	}
	req, err := http.NewRequest("GET", "https://api.stripe.com/v1/"+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(apiKey, "")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// ── Platform fee calculator ────────────────────────────────────
// Month 1: 0%, Month 2+: 1%

func calculatePlatformFee(amountPence int, profileCreatedAt time.Time) int {
	if time.Since(profileCreatedAt).Hours() < 720 {
		return 0
	}
	return amountPence / 100
}

// ── Email helper ──────────────────────────────────────────────

func sendResendEmail(to, subject, html string) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" || to == "" {
		return
	}
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}
	body := map[string]interface{}{
		"from":    "The Concierge <" + fromEmail + ">",
		"to":      []string{to},
		"subject": subject,
		"html":    html,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Warning: sendResendEmail failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

// ── Stripe webhook signature verification ─────────────────────
// Header format: t=timestamp,v1=sig1,v1=sig2

func verifyStripeSignature(payload []byte, sigHeader, secret string) bool {
	var timestamp string
	var sigs []string
	for _, part := range strings.Split(sigHeader, ",") {
		if strings.HasPrefix(part, "t=") {
			timestamp = strings.TrimPrefix(part, "t=")
		} else if strings.HasPrefix(part, "v1=") {
			sigs = append(sigs, strings.TrimPrefix(part, "v1="))
		}
	}
	if timestamp == "" || len(sigs) == 0 {
		return false
	}
	signed := timestamp + "." + string(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signed))
	expected := hex.EncodeToString(mac.Sum(nil))
	for _, sig := range sigs {
		if hmac.Equal([]byte(sig), []byte(expected)) {
			return true
		}
	}
	return false
}

// ── Stripe Connect onboarding ─────────────────────────────────

func handleStripeOnboard(c *gin.Context) {
	if os.Getenv("STRIPE_SECRET_KEY") == "" {
		c.JSON(200, gin.H{"connected": false, "reason": "not_configured"})
		return
	}
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

	returnURL := "https://concierge-ai-gamma.vercel.app?stripe=connected"
	refreshURL := "https://concierge-ai-gamma.vercel.app?stripe=refresh"

	// Already has account — generate new onboarding link
	if p.StripeAccountID != "" {
		body, err := stripePost("account_links", url.Values{
			"account":     {p.StripeAccountID},
			"refresh_url": {refreshURL},
			"return_url":  {returnURL},
			"type":        {"account_onboarding"},
		})
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not create Stripe link"})
			return
		}
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		c.JSON(200, gin.H{"url": result["url"], "existing": true})
		return
	}

	// Create new Express account
	params := url.Values{}
	params.Set("type", "express")
	params.Set("country", "GB")
	params.Set("email", p.Email)
	params.Set("capabilities[card_payments][requested]", "true")
	params.Set("capabilities[transfers][requested]", "true")
	params.Set("business_profile[name]", p.Business)
	params.Set("metadata[slug]", slug)

	body, err := stripePost("accounts", params)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create Stripe account: " + err.Error()})
		return
	}
	var account map[string]interface{}
	json.Unmarshal(body, &account)
	accountID, _ := account["id"].(string)
	if accountID == "" {
		c.JSON(500, gin.H{"error": "Stripe returned no account ID"})
		return
	}

	p.StripeAccountID = accountID
	if err := db.UpdateProfile(*p); err != nil {
		fmt.Printf("Warning: could not save stripe_account_id for %s: %v\n", slug, err)
	}

	linkBody, err := stripePost("account_links", url.Values{
		"account":     {accountID},
		"refresh_url": {refreshURL},
		"return_url":  {returnURL},
		"type":        {"account_onboarding"},
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create onboarding link"})
		return
	}
	var link map[string]interface{}
	json.Unmarshal(linkBody, &link)
	c.JSON(200, gin.H{"url": link["url"], "account_id": accountID})
}

// ── Stripe Connect account status ─────────────────────────────

func handleStripeStatus(c *gin.Context) {
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
	if p.StripeAccountID == "" {
		c.JSON(200, gin.H{"connected": false, "charges_enabled": false})
		return
	}
	body, err := stripeGet("accounts/" + p.StripeAccountID)
	if err != nil {
		c.JSON(200, gin.H{"connected": false, "error": err.Error()})
		return
	}
	var account map[string]interface{}
	json.Unmarshal(body, &account)
	c.JSON(200, gin.H{
		"connected":       true,
		"charges_enabled": account["charges_enabled"],
		"payouts_enabled": account["payouts_enabled"],
		"account_id":      p.StripeAccountID,
	})
}

// ── Create checkout session for deposit ───────────────────────

func handleCreateCheckout(c *gin.Context) {
	var req struct {
		ProfileSlug string `json:"profile_slug"`
		BookingID   string `json:"booking_id"`
		ClientEmail string `json:"client_email"`
		ClientName  string `json:"client_name"`
		ServiceName string `json:"service_name"`
		AmountPence int    `json:"amount_pence"`
		Currency    string `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if req.Currency == "" {
		req.Currency = "gbp"
	}
	p, err := db.GetProfile(req.ProfileSlug)
	if err != nil || p.StripeAccountID == "" {
		c.JSON(400, gin.H{"error": "Professional not set up for payments"})
		return
	}

	platformFee := calculatePlatformFee(req.AmountPence, p.CreatedAt)
	baseURL := "https://concierge-ai-gamma.vercel.app/" + req.ProfileSlug

	params := url.Values{}
	params.Set("mode", "payment")
	params.Set("customer_email", req.ClientEmail)
	params.Set("line_items[0][price_data][currency]", req.Currency)
	params.Set("line_items[0][price_data][unit_amount]", strconv.Itoa(req.AmountPence))
	params.Set("line_items[0][price_data][product_data][name]", req.ServiceName+" — Deposit")
	params.Set("line_items[0][price_data][product_data][description]", "Booking deposit for "+req.ServiceName+" with "+p.Name)
	params.Set("line_items[0][quantity]", "1")
	params.Set("payment_intent_data[transfer_data][destination]", p.StripeAccountID)
	if platformFee > 0 {
		params.Set("payment_intent_data[application_fee_amount]", strconv.Itoa(platformFee))
	}
	params.Set("success_url", baseURL+"?payment=success&booking="+req.BookingID)
	params.Set("cancel_url", baseURL+"?payment=cancelled")
	params.Set("metadata[booking_id]", req.BookingID)
	params.Set("metadata[profile_slug]", req.ProfileSlug)
	params.Set("metadata[client_name]", req.ClientName)

	body, err := stripePost("checkout/sessions", params)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create checkout: " + err.Error()})
		return
	}
	var session map[string]interface{}
	json.Unmarshal(body, &session)
	sessionID, _ := session["id"].(string)
	sessionURL, _ := session["url"].(string)

	payment := db.BookingPayment{
		ID:            fmt.Sprintf("pay_%d", time.Now().UnixMilli()),
		ProfileID:     req.ProfileSlug,
		BookingID:     req.BookingID,
		ClientEmail:   req.ClientEmail,
		ClientName:    req.ClientName,
		ServiceName:   req.ServiceName,
		AmountPence:   req.AmountPence,
		Currency:      req.Currency,
		Status:        "pending",
		StripeSession: sessionID,
		CreatedAt:     time.Now(),
	}
	db.SaveBookingPayment(payment)

	c.JSON(200, gin.H{
		"checkout_url": sessionURL,
		"session_id":   sessionID,
		"amount_pence": req.AmountPence,
		"platform_fee": platformFee,
	})
}

// ── Stripe webhook ─────────────────────────────────────────────

func handleStripeWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not read body"})
		return
	}

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	signature := c.GetHeader("Stripe-Signature")
	if webhookSecret != "" && signature != "" {
		if !verifyStripeSignature(body, signature, webhookSecret) {
			c.JSON(400, gin.H{"error": "Invalid signature"})
			return
		}
	}

	var event map[string]interface{}
	json.Unmarshal(body, &event)
	eventType, _ := event["type"].(string)

	switch eventType {
	case "checkout.session.completed":
		data, _ := event["data"].(map[string]interface{})
		obj, _ := data["object"].(map[string]interface{})
		sessionID, _ := obj["id"].(string)
		meta, _ := obj["metadata"].(map[string]interface{})
		bookingID, _ := meta["booking_id"].(string)
		profileSlug, _ := meta["profile_slug"].(string)
		clientName, _ := meta["client_name"].(string)

		db.UpdateBookingRequest(bookingID, map[string]interface{}{
			"status":     "deposit_paid",
			"updated_at": time.Now(),
		})
		db.UpdateBookingPaymentBySession(sessionID, "paid")

		go func() {
			p, err := db.GetProfile(profileSlug)
			if err != nil || p.Email == "" {
				return
			}
			html := fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto;padding:24px;background:#0c0a08;color:#e8dcc8">
				<div style="color:#c9a96e;font-size:11px;letter-spacing:0.15em;text-transform:uppercase;margin-bottom:16px">✦ The Concierge</div>
				<h2 style="color:#e8dcc8;margin-bottom:8px">Deposit received</h2>
				<p style="color:rgba(232,220,200,0.6);margin-bottom:16px">%s has paid their booking deposit.</p>
				<a href="https://concierge-ai-gamma.vercel.app" style="display:inline-block;background:linear-gradient(135deg,#c9a96e,#7a4f0e);color:#0c0a08;padding:12px 24px;border-radius:8px;text-decoration:none;font-weight:600;font-size:13px">View in Dashboard →</a>
			</div>`, clientName)
			sendResendEmail(p.Email, "💰 Deposit received from "+clientName, html)
		}()
	}

	c.JSON(200, gin.H{"received": true})
}

// ── Get payments for owner ─────────────────────────────────────

func handleGetPayments(c *gin.Context) {
	slug, ok := authenticateToken(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	payments, err := db.GetPaymentsByProfile(slug)
	if err != nil {
		c.JSON(200, gin.H{"payments": []interface{}{}, "total_pence": 0, "total_display": "£0.00"})
		return
	}
	total := 0
	for _, p := range payments {
		if p.Status == "paid" {
			total += p.AmountPence
		}
	}
	c.JSON(200, gin.H{
		"payments":      payments,
		"total_pence":   total,
		"total_display": fmt.Sprintf("£%.2f", float64(total)/100),
	})
}
