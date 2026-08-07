package realtime

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TranscriptRecord represents a final transcript row to be batch-inserted into interview_transcripts.
type TranscriptRecord struct {
	ID            string
	InterviewID   string
	ParticipantID string
	SpeakerType   string
	SpeakerName   string
	Content       string
	Language      string
	StartTimeMs   int64
	EndTimeMs     int64
	Confidence    float64
	Source        string
	IsFinal       bool
	CreatedAt     time.Time
}

// TranscriptBatchSaver accumulates final transcript records in memory and flushes
// them in batches (e.g. every 5 seconds or when capacity threshold is reached) to avoid
// blocking WebSocket and STT pipelines.
type TranscriptBatchSaver struct {
	mu             sync.Mutex
	db             *sqlx.DB
	records        []TranscriptRecord
	batchSize      int
	flushInterval  time.Duration
	stopCh         chan struct{}
	wg             sync.WaitGroup
	flushedRecords []TranscriptRecord
}

// NewTranscriptBatchSaver creates and starts a new async batch saver.
// db may be nil in tests; when nil, records are buffered in memory only
// (flushedRecords) and no SQL is executed.
func NewTranscriptBatchSaver(db *sqlx.DB, batchSize int, flushInterval time.Duration) *TranscriptBatchSaver {
	if batchSize <= 0 {
		batchSize = 100
	}
	if flushInterval <= 0 {
		flushInterval = 5 * time.Second
	}
	s := &TranscriptBatchSaver{
		db:             db,
		records:        make([]TranscriptRecord, 0, batchSize),
		batchSize:      batchSize,
		flushInterval:  flushInterval,
		stopCh:         make(chan struct{}),
		flushedRecords: make([]TranscriptRecord, 0),
	}
	s.wg.Add(1)
	go s.worker()
	return s
}

// Push adds a transcript record to the buffer queue.
// Returns true if queued, or false if ignored (e.g. IsFinal is false).
func (s *TranscriptBatchSaver) Push(record TranscriptRecord) bool {
	if !record.IsFinal {
		return false // strictly enforce storing only final transcripts
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now().UTC()
	}
	if record.Language == "" {
		record.Language = "vi"
	}
	if record.Source == "" {
		record.Source = "audio"
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, record)
	if len(s.records) >= s.batchSize {
		s.flushLocked()
	}
	return true
}

// Flush immediately writes all accumulated records to storage.
func (s *TranscriptBatchSaver) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flushLocked()
}

func (s *TranscriptBatchSaver) flushLocked() {
	if len(s.records) == 0 {
		return
	}

	batch := s.records
	s.records = make([]TranscriptRecord, 0, s.batchSize)

	// Keep an in-memory copy for the chat-history read path and tests.
	s.flushedRecords = append(s.flushedRecords, batch...)

	// No DB (tests / dev without Postgres): memory-only, nothing to persist.
	if s.db == nil {
		return
	}

	// Build a single multi-row INSERT. participant_id is a nullable UUID FK to
	// interview_participants; the realtime layer doesn't create those rows, so
	// we insert NULL rather than a bogus id that would violate the FK.
	var (
		placeholders []string
		args         []interface{}
	)
	col := 1
	for _, r := range batch {
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			col, col+1, col+2, col+3, col+4, col+5, col+6, col+7, col+8, col+9, col+10, col+11,
		))
		var participantID interface{}
		if _, err := uuid.Parse(r.ParticipantID); err == nil {
			participantID = r.ParticipantID
		} else {
			participantID = nil
		}
		args = append(args,
			r.ID, r.InterviewID, participantID, r.SpeakerType, r.SpeakerName,
			r.Content, r.Language, r.StartTimeMs, r.EndTimeMs, r.Confidence,
			r.Source, r.IsFinal,
		)
		col += 12
	}

	q := `INSERT INTO interview_transcripts
		(id, interview_id, participant_id, speaker_type, speaker_name, content, language, start_time_ms, end_time_ms, confidence, source, is_final)
		VALUES ` + strings.Join(placeholders, ", ") + `
		ON CONFLICT (id) DO NOTHING`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, q, args...); err != nil {
		log.Printf("[transcript-saver] batch insert failed (%d rows): %v", len(batch), err)
	}
}

func (s *TranscriptBatchSaver) worker() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Flush()
		case <-s.stopCh:
			s.Flush()
			return
		}
	}
}

// Close stops the background worker and performs a final flush.
func (s *TranscriptBatchSaver) Close() {
	select {
	case <-s.stopCh:
		return // already closed
	default:
		close(s.stopCh)
		s.wg.Wait()
	}
}

// GetFlushedByInterview returns all flushed records for a given interview ID, sorted by timing.
func (s *TranscriptBatchSaver) GetFlushedByInterview(interviewID string) []TranscriptRecord {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []TranscriptRecord
	for _, r := range s.flushedRecords {
		if r.InterviewID == interviewID {
			out = append(out, r)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].StartTimeMs == out[j].StartTimeMs {
			return out[i].CreatedAt.Before(out[j].CreatedAt)
		}
		return out[i].StartTimeMs < out[j].StartTimeMs
	})
	return out
}
