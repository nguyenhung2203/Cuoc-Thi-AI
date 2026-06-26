package livekit

import (
	"context"
	"log"
	"sync"
	"time"
)

// AudioHookService simulates capturing audio from a LiveKit room and sending it to the AI STT Orchestrator.
type AudioHookService struct {
	mu           sync.Mutex
	activeHooks  map[string]context.CancelFunc // roomID -> CancelFunc
}

// NewAudioHookService creates a new AudioHookService.
func NewAudioHookService() *AudioHookService {
	return &AudioHookService{
		activeHooks: make(map[string]context.CancelFunc),
	}
}

// StartHook starts the audio hook for the given room.
// In a real implementation, this would connect to the LiveKit server using lksdk, 
// subscribe to audio tracks, and send the Opus frames to the Python Orchestrator.
func (s *AudioHookService) StartHook(roomID, interviewID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prevent duplicate hooks for the same room
	if _, exists := s.activeHooks[roomID]; exists {
		log.Printf("[audio_hook] Hook already running for room=%s", roomID)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.activeHooks[roomID] = cancel

	log.Printf("[audio_hook] Starting audio stream hook for room=%s, interview=%s", roomID, interviewID)

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("[audio_hook] Stopped audio stream hook for room=%s", roomID)
				return
			case <-ticker.C:
				// Simulated: sending an audio chunk to AI Orchestrator
				log.Printf("[audio_hook] Forwarding audio chunk from room=%s -> Python AI Orchestrator (STT)", roomID)
			}
		}
	}()
}

// StopHook stops the audio hook for the given room.
func (s *AudioHookService) StopHook(roomID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cancel, exists := s.activeHooks[roomID]; exists {
		cancel()
		delete(s.activeHooks, roomID)
		log.Printf("[audio_hook] Requested stop for audio stream hook in room=%s", roomID)
	}
}
