package intelligence_runtime

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrConsentNotVerified = errors.New("intelligence runtime: derived-memory consent is not verified")
	ErrActivationCheck    = errors.New("intelligence runtime: activation check failed")
)

// DerivedMemoryConsentVerifier is the only consent boundary introduced by v0.4.
// The parallel P0 Trust/Data Boundary work owns the authoritative consent
// semantics and server-verified implementation. A runtime must provide a
// verified result before it persists conversation-derived memory.
type DerivedMemoryConsentVerifier interface {
	VerifyConversationDerivedMemory(ctx context.Context, binding PersonBinding, source ConversationMessage) error
}

// FailClosedConsentVerifier intentionally rejects all persistence until the
// authoritative P0 verifier is composed by a server-side integration.
type FailClosedConsentVerifier struct{}

func (FailClosedConsentVerifier) VerifyConversationDerivedMemory(context.Context, PersonBinding, ConversationMessage) error {
	return ErrConsentNotVerified
}

// StaticConsentVerifier is a deterministic staging/test adapter. It accepts
// explicit stable subject and source profile pairs only; it is not a production
// consent implementation.
type StaticConsentVerifier struct {
	mu      sync.RWMutex
	allowed map[string]map[string]bool
}

func NewStaticConsentVerifier(bindings []PersonBinding) *StaticConsentVerifier {
	verifier := &StaticConsentVerifier{allowed: make(map[string]map[string]bool)}
	for _, binding := range bindings {
		verifier.Allow(binding.StableSubjectID, binding.SourceProfileID)
	}
	return verifier
}

func (v *StaticConsentVerifier) Allow(stableSubjectID, sourceProfileID string) {
	if v == nil || stableSubjectID == "" || sourceProfileID == "" {
		return
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.allowed[stableSubjectID] == nil {
		v.allowed[stableSubjectID] = make(map[string]bool)
	}
	v.allowed[stableSubjectID][sourceProfileID] = true
}

func (v *StaticConsentVerifier) VerifyConversationDerivedMemory(_ context.Context, binding PersonBinding, source ConversationMessage) error {
	if v == nil || binding.StableSubjectID == "" || source.Conversation.ProfileID == "" {
		return ErrConsentNotVerified
	}
	v.mu.RLock()
	allowed := v.allowed[binding.StableSubjectID][source.Conversation.ProfileID]
	v.mu.RUnlock()
	if !allowed {
		return ErrConsentNotVerified
	}
	return nil
}

// RuntimeActivation is a server-only feature activation and kill-switch port.
// It is intentionally independent from the legacy feature_flags table.
type RuntimeActivation interface {
	Enabled(ctx context.Context) (bool, error)
}

// StaticRuntimeActivation provides deterministic configuration for staging
// tests and local composition. A false value is an immediate kill switch.
type StaticRuntimeActivation struct {
	Allowed bool
	Err     error
}

func (a StaticRuntimeActivation) Enabled(context.Context) (bool, error) {
	if a.Err != nil {
		return false, a.Err
	}
	return a.Allowed, nil
}
