package intelligence_runtime

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/BraveAI93/concierge-backend/kernel"
)

// ConservativeConversationAdapter supports only source text that explicitly
// asks to schedule or reschedule a meeting/appointment. It intentionally does
// not infer preferences, commitments, or calendar availability.
type ConservativeConversationAdapter struct{}

func (ConservativeConversationAdapter) Map(binding PersonBinding, source ConversationMessage, now time.Time) (IngestionBundle, error) {
	conversation := source.Conversation
	message := source.Message
	if binding.Person.ID == "" || binding.SourceProfileID == "" || conversation.ID == "" || message.ID == "" || message.ConversationID != conversation.ID || conversation.ProfileID != binding.SourceProfileID || message.CreatedAt.IsZero() || now.IsZero() {
		return IngestionBundle{}, ErrUnsupportedSource
	}
	if !isClientRole(message.Role) || !isSchedulingRequest(message.Content) {
		return IngestionBundle{}, ErrUnsupportedSource
	}

	baseID := "legacy-message:" + message.ID
	temporal := kernel.TemporalState{
		EventAt:     message.CreatedAt,
		RecordedAt:  now,
		EffectiveAt: now,
		AttentionAt: now,
	}
	if err := temporal.Validate(); err != nil {
		return IngestionBundle{}, err
	}
	sourceRef := fmt.Sprintf("concierge://conversations/%s/messages/%s", conversation.ID, message.ID)
	sourceRecord := SourceRecord{
		ID:             baseID,
		PersonID:       binding.Person.ID,
		ProfileID:      conversation.ProfileID,
		ConversationID: conversation.ID,
		SessionID:      conversation.SessionID,
		MessageID:      message.ID,
		MessageRole:    message.Role,
		Content:        message.Content,
		ConversationAt: conversation.StartedAt,
		MessageAt:      message.CreatedAt,
		StoredAt:       now,
	}
	claim := kernel.Claim{
		ID:        "claim:" + baseID,
		PersonID:  binding.Person.ID,
		MemoryID:  "memory:" + baseID,
		Statement: "The source message contains an unresolved scheduling request.",
		Predicate: "requests_scheduling_follow_up",
		Object:    "source_message:" + message.ID,
		Temporal:  temporal,
		Freshness: kernel.FreshnessState{LastValidatedAt: now, Status: kernel.FreshnessFresh},
		CreatedAt: now,
	}
	evidence := kernel.Evidence{
		ID:        "evidence:" + baseID,
		PersonID:  binding.Person.ID,
		ClaimID:   claim.ID,
		Stance:    kernel.EvidenceSupports,
		Summary:   "Verbatim source message: " + message.Content,
		Quality:   1,
		Relevance: 1,
		Authority: 0.60,
		Provenance: kernel.Provenance{
			SourceType: "concierge_conversation_message",
			SourceRef:  sourceRef,
			Actor:      message.Role,
			CapturedAt: now,
			Checksum:   message.ID,
		},
		Temporal:  temporal,
		CreatedAt: now,
	}
	claim.EvidenceIDs = []string{evidence.ID}
	memory := kernel.Memory{
		ID:         claim.MemoryID,
		PersonID:   binding.Person.ID,
		Kind:       kernel.MemoryEpisodic,
		Summary:    "Unresolved scheduling request preserved from source message " + message.ID,
		ClaimIDs:   []string{claim.ID},
		ContextIDs: []string{"conversation:" + conversation.ID},
		Temporal:   temporal,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	event := kernel.Event{
		ID:          "event:" + baseID,
		PersonID:    binding.Person.ID,
		Kind:        "legacy_conversation_message",
		Summary:     "Scheduling request received in source message " + message.ID,
		ContextIDs:  []string{"conversation:" + conversation.ID},
		Temporal:    temporal,
		EvidenceIDs: []string{evidence.ID},
		CreatedAt:   now,
	}
	intent := &kernel.PendingIntent{
		ID:         "intent:" + baseID,
		PersonID:   binding.Person.ID,
		Summary:    "Review and prepare a scheduling response supported by source message " + message.ID,
		State:      kernel.IntentCaptured,
		MemoryID:   memory.ID,
		ContextIDs: []string{"conversation:" + conversation.ID},
		Temporal:   temporal,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	openLoop := &kernel.OpenLoop{
		ID:              "loop:" + baseID,
		PersonID:        binding.Person.ID,
		PendingIntentID: intent.ID,
		Label:           "Resolve scheduling request from source message " + message.ID,
		Attention:       temporal,
		InteractionGap: kernel.InteractionGapState{
			PersonID:          binding.Person.ID,
			ContextID:         "conversation:" + conversation.ID,
			LastInteractionAt: message.CreatedAt,
			ObservedAt:        now,
		},
		ContextIDs: []string{"conversation:" + conversation.ID},
		AttentionNeed: kernel.EffortAttention{
			EstimatedEffort:    15 * time.Minute,
			EstimatedAttention: 10 * time.Minute,
			InterruptionCost:   0.10,
			ContextSwitchCost:  0.10,
		},
		CreatedAt: now,
	}
	return IngestionBundle{
		IdempotencyKey: sourceRecord.ID,
		Source:         sourceRecord,
		Event:          event,
		Evidence:       evidence,
		Memory:         memory,
		Claim:          claim,
		Intent:         intent,
		OpenLoop:       openLoop,
		Deadline:       explicitDeadline(message.Content),
	}, nil
}

func isClientRole(role string) bool {
	role = strings.ToLower(strings.TrimSpace(role))
	return role == "user" || role == "client"
}

func isSchedulingRequest(content string) bool {
	text := strings.ToLower(content)
	hasSchedulingVerb := strings.Contains(text, "schedule") || strings.Contains(text, "reschedule") || strings.Contains(text, "appointment")
	hasRequestSignal := strings.Contains(text, "please") || strings.Contains(text, "can you") || strings.Contains(text, "could you") || strings.Contains(text, "need")
	return hasSchedulingVerb && hasRequestSignal
}

var explicitDeadlinePattern = regexp.MustCompile(`(?i)deadline:\s*([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z)`)

// explicitDeadline accepts only a literal RFC3339 source substring. A natural
// language phrase is not converted into a deadline because that would invent a
// semantic time not directly supported by the stored message.
func explicitDeadline(content string) kernel.ActionWindow {
	match := explicitDeadlinePattern.FindStringSubmatch(content)
	if len(match) != 2 {
		return kernel.ActionWindow{}
	}
	at, err := time.Parse(time.RFC3339, match[1])
	if err != nil {
		return kernel.ActionWindow{}
	}
	return kernel.ActionWindow{LatestSafeActionAt: at, EstimatedEffort: 15 * time.Minute}
}
