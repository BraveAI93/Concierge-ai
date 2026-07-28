package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/BraveAI93/concierge-backend/db"
	"github.com/gin-gonic/gin"
)

// fullProfile is a profile with every private field populated, so a leak of any
// one of them is visible in the response body.
func fullProfile() db.Profile {
	return db.Profile{
		ID:                 "11111111-2222-3333-4444-555555555555",
		Slug:               "alice",
		Name:               "Alice",
		Business:           "Alice Therapies",
		Profession:         "Sports Therapist",
		Location:           "London",
		SystemPrompt:       "You are Alice's concierge.",
		ProfileData:        `{"tag":"Move better"}`,
		Accent:             "#c9a96e",
		Active:             true,
		Email:              "alice@example.com",
		PasswordHash:       "hash-should-never-be-public",
		PasswordSalt:       "salt-should-never-be-public",
		GoogleRefreshToken: "google-refresh-should-never-be-public",
		DigestFrequency:    "weekly",
		DigestLastSent:     time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
		StripeAccountID:    "acct_should_never_be_public",
		CreatedAt:          time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
	}
}

// installPublicGetProfile swaps the public profile seam and restores it on cleanup.
func installPublicGetProfile(t *testing.T, fn func(string) (*db.Profile, error)) {
	t.Helper()
	orig := publicGetProfile
	publicGetProfile = fn
	t.Cleanup(func() { publicGetProfile = orig })
}

func getProfile(t *testing.T, slug string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/profile/:slug", handleGetProfile)

	req := httptest.NewRequest(http.MethodGet, "/profile/"+slug, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGetProfileReturnsPublicFields(t *testing.T) {
	p := fullProfile()
	installPublicGetProfile(t, func(slug string) (*db.Profile, error) {
		if slug != "alice" {
			return nil, errNotFound
		}
		cp := p
		return &cp, nil
	})

	w := getProfile(t, "alice")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	want := map[string]any{
		"slug":          p.Slug,
		"name":          p.Name,
		"business":      p.Business,
		"profession":    p.Profession,
		"location":      p.Location,
		"system_prompt": p.SystemPrompt,
		"profile_data":  p.ProfileData,
		"accent":        p.Accent,
		"active":        p.Active,
	}
	for key, wantVal := range want {
		gotVal, ok := got[key]
		if !ok {
			t.Errorf("expected public field %q to be present", key)
			continue
		}
		if gotVal != wantVal {
			t.Errorf("field %q: got %v, want %v", key, gotVal, wantVal)
		}
	}

	// The allowlist is exact: no field beyond the nine above may be serialized.
	if len(got) != len(want) {
		t.Errorf("expected exactly %d public fields, got %d: %v", len(want), len(got), got)
	}
}

func TestGetProfileOmitsPrivateFields(t *testing.T) {
	p := fullProfile()
	installPublicGetProfile(t, func(string) (*db.Profile, error) {
		cp := p
		return &cp, nil
	})

	w := getProfile(t, "alice")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}

	var got map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	forbiddenKeys := []string{
		"id",
		"email",
		"password_hash",
		"password_salt",
		"google_refresh_token",
		"stripe_account_id",
		"digest_frequency",
		"digest_last_sent",
		"created_at",
	}
	for _, key := range forbiddenKeys {
		if _, present := got[key]; present {
			t.Errorf("private field %q must not appear in the public response", key)
		}
	}

	// Belt and braces: the secret values themselves must not appear anywhere in
	// the body, whatever key they might have been serialized under.
	body := w.Body.String()
	forbiddenValues := []string{
		p.ID,
		p.Email,
		p.PasswordHash,
		p.PasswordSalt,
		p.GoogleRefreshToken,
		p.StripeAccountID,
		p.DigestFrequency,
	}
	for _, val := range forbiddenValues {
		if strings.Contains(body, val) {
			t.Errorf("private value %q leaked into the public response body: %s", val, body)
		}
	}
}

func TestGetProfileNotFoundUnchanged(t *testing.T) {
	installPublicGetProfile(t, func(string) (*db.Profile, error) {
		return nil, errNotFound
	})

	w := getProfile(t, "does-not-exist")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d (body: %s)", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["error"] != "Profile not found" {
		t.Errorf("expected unchanged 404 body, got %q", w.Body.String())
	}
}
