package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BraveAI93/concierge-backend/db"
	"github.com/gin-gonic/gin"
)

// fakeProfileStore is an in-memory stand-in for the profile table, wired in
// through the signup* function variables in main.go.
type fakeProfileStore struct {
	profiles map[string]db.Profile
	sessions []db.Session
}

func newFakeStore(seed db.Profile) *fakeProfileStore {
	return &fakeProfileStore{profiles: map[string]db.Profile{seed.Slug: seed}}
}

// install swaps the signup seams for the fake store and returns a restore func.
func (f *fakeProfileStore) install(t *testing.T) {
	t.Helper()
	origGet, origUpdate, origSave := signupGetProfile, signupUpdateProfile, signupSaveSession
	signupGetProfile = func(slug string) (*db.Profile, error) {
		p, ok := f.profiles[slug]
		if !ok {
			return nil, errNotFound
		}
		cp := p
		return &cp, nil
	}
	signupUpdateProfile = func(p db.Profile) error {
		f.profiles[p.Slug] = p
		return nil
	}
	signupSaveSession = func(s db.Session) error {
		f.sessions = append(f.sessions, s)
		return nil
	}
	t.Cleanup(func() {
		signupGetProfile, signupUpdateProfile, signupSaveSession = origGet, origUpdate, origSave
	})
}

type notFoundError struct{}

func (notFoundError) Error() string { return "profile not found" }

var errNotFound = notFoundError{}

func postSignup(t *testing.T, slug, email, password string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/auth/signup", handleSignup)

	body, err := json.Marshal(map[string]string{"slug": slug, "email": email, "password": password})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSignupFirstTimeSucceeds(t *testing.T) {
	store := newFakeStore(db.Profile{Slug: "alice", Name: "Alice"})
	store.install(t)

	w := postSignup(t, "alice", "alice@example.com", "hunter2")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}

	saved := store.profiles["alice"]
	if saved.Email != "alice@example.com" {
		t.Errorf("expected email to be set, got %q", saved.Email)
	}
	if saved.PasswordSalt == "" {
		t.Error("expected password salt to be set")
	}
	if want := hashPassword("hunter2", saved.PasswordSalt); saved.PasswordHash != want {
		t.Errorf("stored hash %q does not match hash of submitted password", saved.PasswordHash)
	}
	if len(store.sessions) != 1 || store.sessions[0].Slug != "alice" {
		t.Errorf("expected one session for alice, got %+v", store.sessions)
	}
}

func TestSignupSecondAttemptReturns409(t *testing.T) {
	salt := "existing-salt"
	store := newFakeStore(db.Profile{
		Slug:         "alice",
		Email:        "alice@example.com",
		PasswordHash: hashPassword("original-password", salt),
		PasswordSalt: salt,
	})
	store.install(t)

	w := postSignup(t, "alice", "attacker@evil.example", "attacker-password")
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d (body: %s)", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !strings.Contains(resp["error"], "Account already exists") {
		t.Errorf("expected neutral 'Account already exists' error, got %q", resp["error"])
	}
	if len(store.sessions) != 0 {
		t.Errorf("expected no session to be created, got %+v", store.sessions)
	}
}

func TestSignupRejectionLeavesCredentialsUnchanged(t *testing.T) {
	salt := "existing-salt"
	original := db.Profile{
		Slug:         "alice",
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: hashPassword("original-password", salt),
		PasswordSalt: salt,
	}
	store := newFakeStore(original)
	store.install(t)

	if w := postSignup(t, "alice", "attacker@evil.example", "attacker-password"); w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d (body: %s)", w.Code, w.Body.String())
	}

	after := store.profiles["alice"]
	if after != original {
		t.Fatalf("profile was modified by a rejected signup:\n before: %+v\n after:  %+v", original, after)
	}
	// The original password must still verify against the stored hash.
	if hashPassword("original-password", after.PasswordSalt) != after.PasswordHash {
		t.Error("original password no longer matches stored credentials")
	}
}
