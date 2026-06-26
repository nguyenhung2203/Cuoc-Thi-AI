package events

// VisibilityRule defines who can receive a given server→client event.
type VisibilityRule int

const (
	// VisibleAll — every participant in the room receives the event.
	VisibleAll VisibilityRule = iota
	// VisibleRecruitersOnly — only participants with role "recruiter" receive the event.
	VisibleRecruitersOnly
	// VisibleSenderOnly — only the participant who triggered the event receives it.
	VisibleSenderOnly
	// VisibleConfigDependent — controlled by per-room configuration (e.g. transcript for candidates).
	VisibleConfigDependent
)

// eventVisibility maps every server→client event to its default visibility rule.
// Any event not listed here defaults to VisibleAll.
var eventVisibility = map[string]VisibilityRule{
	// ── Room presence events — visible to everyone ────────────────────────
	EventRoomJoined:         VisibleSenderOnly, // only the joiner gets their own ACK
	EventRoomUserJoined:     VisibleAll,
	EventRoomUserLeft:       VisibleAll,
	EventRoomPresenceUpdate: VisibleAll,

	// ── Interview lifecycle — visible to everyone ─────────────────────────
	EventInterviewStarted:   VisibleAll,
	EventInterviewCompleted: VisibleAll,

	// ── Media — visible to everyone ───────────────────────────────────────
	EventMediaStatusChanged: VisibleAll,

	// ── Chat — depends on message visibility field ────────────────────────
	// Filtering is done at the handler level; the default here is All.
	EventChatMessage: VisibleAll,

	// ── Transcript — config-dependent for candidate ───────────────────────
	EventTranscriptUpdate: VisibleConfigDependent,

	// ── AI events — recruiter only ────────────────────────────────────────
	EventAIThinking:    VisibleRecruitersOnly,
	EventAISuggestion:  VisibleRecruitersOnly,
	EventAIScoreUpdate: VisibleRecruitersOnly,
	EventAIWarning:     VisibleRecruitersOnly,
	EventAIError:       VisibleRecruitersOnly,

	// ── Report & Notes — recruiter only ───────────────────────────────────
	EventReportReady: VisibleRecruitersOnly,
	EventNoteCreate:  VisibleRecruitersOnly,

	// ── Generic error — only the sender ──────────────────────────────────
	EventError: VisibleSenderOnly,
}

// GetVisibility returns the VisibilityRule for a server→client event name.
// Unknown events default to VisibleAll.
func GetVisibility(event string) VisibilityRule {
	if rule, ok := eventVisibility[event]; ok {
		return rule
	}
	return VisibleAll
}

// ShouldSendToParticipant returns true when the given participant is allowed to
// receive the event, based on the centralised visibility rules.
//
// Parameters:
//   - event: the server→client event name (e.g. "ai:suggestion")
//   - participantType: role of the target participant
//   - isSender: true when the target is the same client that triggered the event
//   - transcriptEnabledForCandidate: per-room config flag (only relevant for transcript events)
//   - isMockInterview: in mock mode the candidate sees ai:thinking
func ShouldSendToParticipant(
	event string,
	participantType ParticipantType,
	isSender bool,
	transcriptEnabledForCandidate bool,
	isMockInterview bool,
) bool {
	rule := GetVisibility(event)

	switch rule {
	case VisibleSenderOnly:
		return isSender

	case VisibleRecruitersOnly:
		// Exception: in mock interview mode the candidate can see ai:thinking
		if isMockInterview && event == EventAIThinking {
			return true
		}
		return participantType == ParticipantRecruiter

	case VisibleConfigDependent:
		// transcript:update — recruiter always sees it, candidate only if enabled
		if event == EventTranscriptUpdate {
			if participantType == ParticipantRecruiter {
				return true
			}
			return transcriptEnabledForCandidate
		}
		return true

	default: // VisibleAll
		return true
	}
}

// ShouldSendChatToParticipant handles the special case of chat messages whose
// visibility is embedded in the payload (not just the event name).
func ShouldSendChatToParticipant(
	visibility ChatVisibility,
	participantType ParticipantType,
) bool {
	switch visibility {
	case VisibilityRecruiterOnly:
		return participantType == ParticipantRecruiter
	default: // VisibilityRoom
		return true
	}
}
