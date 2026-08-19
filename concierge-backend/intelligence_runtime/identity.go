package intelligence_runtime

import (
	"context"
	"sync"
)

// StaticIdentityResolver is a deterministic local implementation of the
// stable-subject mapping contract. It models a future server-only binding store
// and is used by runtime tests; it never accepts or resolves public slugs.
type StaticIdentityResolver struct {
	mu       sync.RWMutex
	bindings map[string]PersonBinding
}

func NewStaticIdentityResolver(bindings []PersonBinding) *StaticIdentityResolver {
	resolver := &StaticIdentityResolver{bindings: make(map[string]PersonBinding, len(bindings))}
	for _, binding := range bindings {
		if binding.StableSubjectID == "" || binding.SourceProfileID == "" || binding.Person.ID == "" || binding.World.PersonID != binding.Person.ID {
			continue
		}
		resolver.bindings[binding.StableSubjectID] = cloneBinding(binding)
	}
	return resolver
}

func (r *StaticIdentityResolver) Resolve(_ context.Context, principal AuthenticatedPrincipal) (PersonBinding, error) {
	if r == nil || principal.StableSubjectID == "" {
		return PersonBinding{}, ErrUnknownIdentity
	}
	r.mu.RLock()
	binding, ok := r.bindings[principal.StableSubjectID]
	r.mu.RUnlock()
	if !ok {
		return PersonBinding{}, ErrUnknownIdentity
	}
	return cloneBinding(binding), nil
}

func cloneBinding(binding PersonBinding) PersonBinding {
	binding.AllowedSourceProfileIDs = append([]string(nil), binding.AllowedSourceProfileIDs...)
	binding.World.EntityIDs = append([]string(nil), binding.World.EntityIDs...)
	binding.World.ContextIDs = append([]string(nil), binding.World.ContextIDs...)
	return binding
}
