package intelligence_runtime

import (
	"context"
	"sync"
)

// ServerSessionSubjectLookup must be implemented by a server-authenticated
// session mechanism. It maps an opaque validated session token to an immutable
// stable subject; it must not return a browser-provided profile slug.
type ServerSessionSubjectLookup interface {
	StableSubjectForSession(ctx context.Context, sessionToken string) (string, error)
}

// ServerSessionIdentityAdapter composes server session authentication with a
// canonical binding resolver. It deliberately has no public profile/person
// selector in its API.
type ServerSessionIdentityAdapter struct {
	Sessions ServerSessionSubjectLookup
	Bindings IdentityResolver
}

func (a ServerSessionIdentityAdapter) ResolveSession(ctx context.Context, sessionToken string) (PersonBinding, error) {
	if a.Sessions == nil || a.Bindings == nil || sessionToken == "" {
		return PersonBinding{}, ErrUnknownIdentity
	}
	subject, err := a.Sessions.StableSubjectForSession(ctx, sessionToken)
	if err != nil || subject == "" {
		return PersonBinding{}, ErrUnknownIdentity
	}
	return a.Bindings.Resolve(ctx, AuthenticatedPrincipal{StableSubjectID: subject})
}

// StaticServerSessionSubjectLookup is staging/test-only. It models the output
// of a server-validated session layer, not browser-controlled input.
type StaticServerSessionSubjectLookup struct {
	mu       sync.RWMutex
	subjects map[string]string
}

func NewStaticServerSessionSubjectLookup(values map[string]string) *StaticServerSessionSubjectLookup {
	copyValues := make(map[string]string, len(values))
	for session, subject := range values {
		if session != "" && subject != "" {
			copyValues[session] = subject
		}
	}
	return &StaticServerSessionSubjectLookup{subjects: copyValues}
}

func (l *StaticServerSessionSubjectLookup) StableSubjectForSession(_ context.Context, sessionToken string) (string, error) {
	if l == nil || sessionToken == "" {
		return "", ErrUnknownIdentity
	}
	l.mu.RLock()
	subject, ok := l.subjects[sessionToken]
	l.mu.RUnlock()
	if !ok {
		return "", ErrUnknownIdentity
	}
	return subject, nil
}

// PostgresIdentityResolver reads a durable server-side stable-subject binding
// from the staging repository. It does not read legacy profiles or sessions.
type PostgresIdentityResolver struct {
	Repository *PostgresRuntimeRepository
}

func (r PostgresIdentityResolver) Resolve(ctx context.Context, principal AuthenticatedPrincipal) (PersonBinding, error) {
	if r.Repository == nil || principal.StableSubjectID == "" {
		return PersonBinding{}, ErrUnknownIdentity
	}
	return r.Repository.ResolveBinding(ctx, principal.StableSubjectID)
}

// ExistingLegacySlugSessionLookup is intentionally absent. The legacy session
// model demonstrated by v0.3 maps tokens to slugs, not immutable subjects. A
// production adapter must wait for the reviewed P0 identity migration/backfill
// rather than silently promoting a public slug into a stable subject.
