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

// ownerFixture is a profile with every private field populated with a
// distinctive sentinel, so a leak of any one of them is visible in the body.
func ownerFixture(slug string) db.Profile {
	return db.Profile{
		ID:                 "owner-id-" + slug,
		Slug:               slug,
		Name:               "Alice",
		Business:           "Alice Therapies",
		Profession:         "Sports Therapist",
		Location:           "London",
		SystemPrompt:       "You are Alice's concierge.",
		ProfileData:        `{"tag":"Move better"}`,
		Accent:             "#c9a96e",
		Active:             true,
		Email:              "alice@example.com",
		PasswordHash:       "owner-hash-must-never-ship",
		PasswordSalt:       "owner-salt-must-never-ship",
		GoogleRefreshToken: "owner-google-refresh-must-never-ship",
		DigestFrequency:    "weekly",
		DigestLastSent:     time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
		StripeAccountID:    "acct_must_never_ship",
		CreatedAt:          time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
	}
}

// ownerForbiddenKeys must never appear in either owner profile response.
var ownerForbiddenKeys = []string{
	"id",
	"password_hash",
	"password_salt",
	"google_refresh_token",
	"stripe_account_id",
	"digest_last_sent",
	"created_at",
}

// ownerSecretValues must never appear anywhere in either response body,
// whatever key they might have been serialized under.
func ownerSecretValues(p db.Profile) []string {
	return []string{
		p.ID,
		p.PasswordHash,
		p.PasswordSalt,
		p.GoogleRefreshToken,
		p.StripeAccountID,
	}
}

// installOwnerSeams swaps the owner profile seams and restores them on cleanup.
func installOwnerSeams(t *testing.T, authSlug string, profiles ...db.Profile) {
	t.Helper()
	origAuth, origGet, origByEmail := ownerAuth, ownerGetProfile, ownerGetProfilesByEmail

	bySlug := map[string]db.Profile{}
	for _, p := range profiles {
		bySlug[p.Slug] = p
	}

	ownerAuth = func(*gin.Context) (string, bool) { return authSlug, true }
	ownerGetProfile = func(slug string) (*db.Profile, error) {
		p, ok := bySlug[slug]
		if !ok {
			return nil, errNotFound
		}
		cp := p
		return &cp, nil
	}
	ownerGetProfilesByEmail = func(email string) ([]db.Profile, error) {
		var out []db.Profile
		for _, p := range profiles {
			if p.Email == email {
				out = append(out, p)
			}
		}
		return out, nil
	}

	t.Cleanup(func() {
		ownerAuth, ownerGetProfile, ownerGetProfilesByEmail = origAuth, origGet, origByEmail
	})
}

// callOwner routes a request through a bare gin engine, with no Authorization
// header unless one is supplied.
func callOwner(t *testing.T, method, path string, handler gin.HandlerFunc, authHeader string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Handle(method, path, handler)

	req := httptest.NewRequest(method, path, nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ─── GET /owner/profile (singular) ─────────────────────────────────────────

func TestOwnerProfileReturnsAllowlistedFields(t *testing.T) {
	p := ownerFixture("alice")
	installOwnerSeams(t, "alice", p)

	w := callOwner(t, http.MethodGet, "/owner/profile", handleGetOwnerProfile, "Bearer stub")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	want := map[string]any{
		"slug":             p.Slug,
		"name":             p.Name,
		"business":         p.Business,
		"profession":       p.Profession,
		"location":         p.Location,
		"system_prompt":    p.SystemPrompt,
		"profile_data":     p.ProfileData,
		"accent":           p.Accent,
		"active":           p.Active,
		"email":            p.Email,
		"digest_frequency": p.DigestFrequency,
	}
	for key, wantVal := range want {
		gotVal, ok := got[key]
		if !ok {
			t.Errorf("expected allowlisted field %q to be present", key)
			continue
		}
		if gotVal != wantVal {
			t.Errorf("field %q: got %v, want %v", key, gotVal, wantVal)
		}
	}

	// The allowlist is exact: any field beyond the eleven above is a leak.
	if len(got) != len(want) {
		t.Errorf("expected exactly %d fields, got %d: %v", len(want), len(got), got)
	}
}

func TestOwnerProfileOmitsPrivateFields(t *testing.T) {
	p := ownerFixture("alice")
	installOwnerSeams(t, "alice", p)

	w := callOwner(t, http.MethodGet, "/owner/profile", handleGetOwnerProfile, "Bearer stub")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}

	var got map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	for _, key := range ownerForbiddenKeys {
		if _, present := got[key]; present {
			t.Errorf("private field %q must not appear in GET /owner/profile", key)
		}
	}

	body := w.Body.String()
	for _, val := range ownerSecretValues(p) {
		if strings.Contains(body, val) {
			t.Errorf("secret value %q leaked into GET /owner/profile body: %s", val, body)
		}
	}
}

func TestOwnerProfileNotFoundUnchanged(t *testing.T) {
	installOwnerSeams(t, "ghost")

	w := callOwner(t, http.MethodGet, "/owner/profile", handleGetOwnerProfile, "Bearer stub")
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

// ─── GET /owner/profiles (plural) ──────────────────────────────────────────

func TestOwnerProfilesReturnsAllowlistedFields(t *testing.T) {
	alice := ownerFixture("alice")
	second := ownerFixture("alice-studio")
	second.Name = "Alice Studio"
	installOwnerSeams(t, "alice", alice, second)

	w := callOwner(t, http.MethodGet, "/owner/profiles", handleGetOwnerProfiles, "Bearer stub")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}

	var resp struct {
		Profiles []map[string]any `json:"profiles"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(resp.Profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d: %v", len(resp.Profiles), resp.Profiles)
	}

	want := []map[string]any{
		{"slug": alice.Slug, "name": alice.Name, "business": alice.Business, "profession": alice.Profession},
		{"slug": second.Slug, "name": second.Name, "business": second.Business, "profession": second.Profession},
	}
	for i, entry := range resp.Profiles {
		for key, wantVal := range want[i] {
			gotVal, ok := entry[key]
			if !ok {
				t.Errorf("profile %d: expected allowlisted field %q to be present", i, key)
				continue
			}
			if gotVal != wantVal {
				t.Errorf("profile %d field %q: got %v, want %v", i, key, gotVal, wantVal)
			}
		}
		// The allowlist is exact: the switcher needs these four and nothing else.
		if len(entry) != len(want[i]) {
			t.Errorf("profile %d: expected exactly %d fields, got %d: %v", i, len(want[i]), len(entry), entry)
		}
	}
}

func TestOwnerProfilesOmitsPrivateFields(t *testing.T) {
	alice := ownerFixture("alice")
	installOwnerSeams(t, "alice", alice)

	w := callOwner(t, http.MethodGet, "/owner/profiles", handleGetOwnerProfiles, "Bearer stub")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}

	var resp struct {
		Profiles []map[string]json.RawMessage `json:"profiles"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(resp.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(resp.Profiles))
	}
	// email is allowed on the singular endpoint but has no consumer on the
	// switcher, so it must not appear here either.
	for _, key := range append(ownerForbiddenKeys, "email", "digest_frequency") {
		if _, present := resp.Profiles[0][key]; present {
			t.Errorf("private field %q must not appear in GET /owner/profiles", key)
		}
	}

	body := w.Body.String()
	for _, val := range append(ownerSecretValues(alice), alice.Email) {
		if strings.Contains(body, val) {
			t.Errorf("secret value %q leaked into GET /owner/profiles body: %s", val, body)
		}
	}
}

// ─── auth behaviour is unchanged ───────────────────────────────────────────

// No seam is installed here, so these exercise the real authenticateToken.
// An absent Authorization header short-circuits before any DB call, so this
// needs no live Supabase.

func TestOwnerProfileUnauthenticatedStill401(t *testing.T) {
	w := callOwner(t, http.MethodGet, "/owner/profile", handleGetOwnerProfile, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (body: %s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Unauthorized") {
		t.Errorf("expected unchanged Unauthorized body, got %q", w.Body.String())
	}
}

func TestOwnerProfilesUnauthenticatedStill401(t *testing.T) {
	w := callOwner(t, http.MethodGet, "/owner/profiles", handleGetOwnerProfiles, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (body: %s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Unauthorized") {
		t.Errorf("expected unchanged Unauthorized body, got %q", w.Body.String())
	}
}

// A rejected session must also 401 rather than fall through to data.
func TestOwnerProfilesRejectedSessionStill401(t *testing.T) {
	orig := ownerAuth
	ownerAuth = func(*gin.Context) (string, bool) { return "", false }
	t.Cleanup(func() { ownerAuth = orig })

	w := callOwner(t, http.MethodGet, "/owner/profiles", handleGetOwnerProfiles, "Bearer invalid")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (body: %s)", w.Code, w.Body.String())
	}
}
