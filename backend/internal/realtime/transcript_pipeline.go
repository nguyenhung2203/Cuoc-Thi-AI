package realtime

import (
	"log"
	"time"

	"backend/internal/realtime/events"
)

// transcriptJob represents a single transcript update queued for broadcast.
type transcriptJob struct {
	roomID  string
	payload events.TranscriptUpdatePayload
}

// TranscriptPipeline manages an asynchronous buffer queue for broadcasting
// STT transcript updates from the AI Orchestrator into rooms without blocking HTTP workers.
type TranscriptPipeline struct {
	roomManager *RoomManager
	inboundCh   chan transcriptJob
}

// NewTranscriptPipeline creates a pipeline with a specified buffer capacity and starts its worker.
func NewTranscriptPipeline(rm *RoomManager, bufferSize int) *TranscriptPipeline {
	tp := &TranscriptPipeline{
		roomManager: rm,
		inboundCh:   make(chan transcriptJob, bufferSize),
	}
	go tp.worker()
	return tp
}

// Push enqueues a transcript update payload. Returns false if the buffer is full.
func (tp *TranscriptPipeline) Push(roomID string, payload events.TranscriptUpdatePayload) bool {
	if payload.CreatedAt.IsZero() {
		payload.CreatedAt = time.Now().UTC()
	}

	select {
	case tp.inboundCh <- transcriptJob{roomID: roomID, payload: payload}:
		return true
	default:
		log.Printf("[transcript-pipeline] buffer full (cap=%d), dropping transcript update for room=%s", cap(tp.inboundCh), roomID)
		return false
	}
}

// worker drains the inbound channel and broadcasts each transcript job to the appropriate room.
func (tp *TranscriptPipeline) worker() {
	for job := range tp.inboundCh {
		room := tp.roomManager.Get(job.roomID)
		if room == nil {
			continue
		}

		env, err := events.NewEnvelope(events.EventTranscriptUpdate, "", room.ID, room.InterviewID, job.payload)
		if err != nil {
			log.Printf("[transcript-pipeline] failed to create envelope: %v", err)
			continue
		}

		rawJSON, err := env.ToJSON()
		if err != nil {
			continue
		}

		// BroadcastWithVisibility applies role visibility rules:
		// Recruiter always sees it; Candidate only if Room.TranscriptEnabledForCandidate is true.
		room.BroadcastWithVisibility(events.EventTranscriptUpdate, "", rawJSON)
	}
}
