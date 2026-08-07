package realtime

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"

	"backend/internal/models"
)

// mock_live.go implements a WebSocket proxy that bridges a candidate's browser
// to the Gemini Live API for a spoken, real-time mock interview.
//
// Flow:
//   browser mic (PCM16 @16kHz) --> WS /ws/mock-live --> Gemini Live
//   Gemini Live (PCM16 @24kHz + transcript) --> WS --> browser (play + show)
//
// The Gemini API key never leaves the server. The browser only talks to us.

// mockLiveModel is a Live-capable Gemini model (native audio dialog).
// Overridable via GEMINI_LIVE_MODEL. This native-audio model serves the
// Live API (bidiGenerateContent) which the older 2.0-flash-live alias no
// longer does on current API versions.
const mockLiveModel = "gemini-2.5-flash-native-audio-latest"

// mockLiveVoice is the prebuilt voice used for the AI interviewer.
const mockLiveVoice = "Puck"

// Client -> server frame from the browser.
// type "audio": base64 PCM16 mono @16kHz in `data`.
// type "end":   candidate finished speaking (commit turn).
// type "text":  optional typed answer fallback.
type mockLiveClientMsg struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Text string `json:"text,omitempty"`
}

// Server -> browser frame.
// type "audio":       base64 PCM16 mono @24kHz chunk to play.
// type "transcript":  { role, text } partial/final transcript.
// type "turn_complete": AI finished its turn.
// type "interrupted": AI generation was interrupted (flush playback).
// type "ready":       Gemini session established.
// type "error":       { message }.
type mockLiveServerMsg struct {
	Type    string `json:"type"`
	Data    string `json:"data,omitempty"`
	Role    string `json:"role,omitempty"`
	Text    string `json:"text,omitempty"`
	Message string `json:"message,omitempty"`
}

// handleMockLive upgrades the connection and proxies it to Gemini Live.
func (s *Server) handleMockLive(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	model := os.Getenv("GEMINI_LIVE_MODEL")
	if s.aiSettings != nil {
		if cfg, err := s.aiSettings.GetRuntimeConfig(r.Context()); err == nil {
			if len(cfg.APIKeys) > 0 {
				apiKey = cfg.APIKeys[0]
			}
			if cfg.LiveModel != "" {
				model = cfg.LiveModel
			}
		}
	}
	if apiKey == "" || apiKey == "your-gemini-api-key" {
		http.Error(w, "AI voice service not configured (missing GEMINI_API_KEY)", http.StatusServiceUnavailable)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[mock-live] upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// System instruction personalises the interviewer to the requested role.
	role := r.URL.Query().Get("role")
	if role == "" {
		role = "Software Developer"
	}
	level := r.URL.Query().Get("level")
	if level == "" {
		level = "middle"
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if model == "" {
		model = mockLiveModel
	}

	session, err := newGeminiLiveSession(ctx, apiKey, model, role, level)
	if err != nil {
		log.Printf("[mock-live] gemini connect failed: %v", err)
		writeMockLiveJSON(conn, mockLiveServerMsg{Type: "error", Message: "Không kết nối được AI phỏng vấn"})
		return
	}
	defer session.Close()

	start := time.Now()
	var usageMu sync.Mutex
	var promptTokens, responseTokens, totalTokens int32
	outcome := "success"
	var outcomeMu sync.Mutex
	setOutcome := func(status string) {
		outcomeMu.Lock()
		if outcome == "success" || status == "failed" {
			outcome = status
		}
		outcomeMu.Unlock()
	}
	usageSink := func(usage *genai.UsageMetadata) {
		if usage == nil {
			return
		}
		usageMu.Lock()
		promptTokens = usage.PromptTokenCount
		responseTokens = usage.ResponseTokenCount
		totalTokens = usage.TotalTokenCount
		usageMu.Unlock()
	}
	defer func() {
		if s.aiLogSvc == nil {
			return
		}
		usageMu.Lock()
		in, out, total := promptTokens, responseTokens, totalTokens
		usageMu.Unlock()
		if total == 0 {
			total = in + out
		}
		outcomeMu.Lock()
		status := outcome
		outcomeMu.Unlock()
		cost := float64(in)*0.30/1_000_000 + float64(out)*1.25/1_000_000
		s.aiLogSvc.LogAsync(&models.AIRequestLog{
			CompanyID: "", Provider: "gemini", Model: sql.NullString{String: model, Valid: true},
			Operation: sql.NullString{String: "mock_interview_live", Valid: true},
			InputJSON: models.JSONB(`{"modality":"audio"}`), LatencyMs: sql.NullInt32{Int32: int32(time.Since(start).Milliseconds()), Valid: true},
			TokensIn: sql.NullInt32{Int32: in, Valid: in > 0}, TokensOut: sql.NullInt32{Int32: out, Valid: out > 0}, TotalTokens: sql.NullInt32{Int32: total, Valid: total > 0},
			Cost: sql.NullFloat64{Float64: cost, Valid: total > 0}, Status: status, CreatedAt: time.Now(),
		})
	}()

	writeMockLiveJSON(conn, mockLiveServerMsg{Type: "ready"})

	// A single writer goroutine owns the browser socket to avoid concurrent writes.
	var writeMu sync.Mutex
	send := func(m mockLiveServerMsg) {
		writeMu.Lock()
		defer writeMu.Unlock()
		_ = conn.WriteJSON(m)
	}

	// We will wait for the client to send a 'start_interview' message to trigger the greeting.

	// Gemini -> browser pump.
	go pumpGeminiToBrowser(ctx, cancel, session, send, usageSink, setOutcome)

	// browser -> Gemini pump (blocks until the socket closes).
	pumpBrowserToGemini(ctx, cancel, conn, session, send, setOutcome)
}

// newGeminiLiveSession opens a Live session configured for audio dialog in Vietnamese.
func newGeminiLiveSession(ctx context.Context, apiKey, model, role, level string) (*genai.Session, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	sysText := "Bạn là một nhà tuyển dụng chuyên nghiệp, thân thiện, đang phỏng vấn thử cho vị trí " +
		role + " trình độ " + level + ". " +
		"Hãy nói bằng tiếng Việt tự nhiên. Mỗi lượt chỉ hỏi một câu hỏi rõ ràng, " +
		"lắng nghe câu trả lời, phản hồi ngắn gọn rồi hỏi câu tiếp theo. " +
		"Đào sâu khi cần và giữ không khí khích lệ. Không đánh giá các yếu tố nhạy cảm."

	cfg := &genai.LiveConnectConfig{
		ResponseModalities: []genai.Modality{genai.ModalityAudio},
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: sysText}},
		},
		SpeechConfig: &genai.SpeechConfig{
			LanguageCode: "vi-VN",
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{VoiceName: mockLiveVoice},
			},
		},
		// Ask Gemini to also return text transcripts of both sides.
		InputAudioTranscription:  &genai.AudioTranscriptionConfig{},
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},
	}

	return client.Live.Connect(ctx, model, cfg)
}

// pumpBrowserToGemini reads client frames and forwards audio/text to Gemini.
func pumpBrowserToGemini(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, session *genai.Session, send func(mockLiveServerMsg), setOutcome func(string)) {
	defer cancel()
	conn.SetReadLimit(2 << 20) // 2 MiB per frame is ample for audio chunks.
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, raw, err := conn.ReadMessage()
		if err != nil {
			setOutcome("partial")
			return // browser closed
		}

		var msg mockLiveClientMsg
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "audio":
			pcm, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil || len(pcm) == 0 {
				continue
			}
			if err := session.SendRealtimeInput(genai.LiveRealtimeInput{
				Audio: &genai.Blob{Data: pcm, MIMEType: "audio/pcm;rate=16000"},
			}); err != nil {
				setOutcome("failed")
				send(mockLiveServerMsg{Type: "error", Message: "Không gửi được âm thanh tới AI"})
			}
		case "end":
			// Signal end of the user's audio stream so Gemini responds.
			_ = session.SendRealtimeInput(genai.LiveRealtimeInput{AudioStreamEnd: true})
		case "text":
			if msg.Text != "" {
				_ = session.SendClientContent(genai.LiveClientContentInput{
					Turns: []*genai.Content{{
						Role:  "user",
						Parts: []*genai.Part{{Text: msg.Text}},
					}},
					TurnComplete: genai.Ptr(true),
				})
			}
		case "start_interview":
			_ = session.SendClientContent(genai.LiveClientContentInput{
				Turns: []*genai.Content{{
					Role:  "user",
					Parts: []*genai.Part{{Text: "Xin chào, tôi đã sẵn sàng. Hãy bắt đầu buổi phỏng vấn và đặt câu hỏi đầu tiên."}},
				}},
				TurnComplete: genai.Ptr(true),
			})
		case "close":
			return
		}
	}
}

// pumpGeminiToBrowser reads server messages from Gemini and forwards them.
func pumpGeminiToBrowser(ctx context.Context, cancel context.CancelFunc, session *genai.Session, send func(mockLiveServerMsg), usageSink func(*genai.UsageMetadata), setOutcome func(string)) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		resp, err := session.Receive()
		if err != nil {
			// Surface unexpected Gemini disconnects instead of silently ending the
			// browser socket after only a partial transcription has arrived.
			if ctx.Err() == nil {
				setOutcome("failed")
				log.Printf("[mock-live] gemini receive failed: %v", err)
				send(mockLiveServerMsg{Type: "error", Message: "Kết nối AI bị gián đoạn. Vui lòng bắt đầu lại phiên luyện tập."})
			}
			return
		}
		if resp == nil {
			continue
		}
		if usageSink != nil && resp.UsageMetadata != nil {
			usageSink(resp.UsageMetadata)
		}
		if resp.ServerContent == nil {
			continue
		}
		sc := resp.ServerContent

		if sc.Interrupted {
			send(mockLiveServerMsg{Type: "interrupted"})
		}

		// Input (candidate) transcript.
		if sc.InputTranscription != nil && sc.InputTranscription.Text != "" {
			send(mockLiveServerMsg{Type: "transcript", Role: "candidate", Text: sc.InputTranscription.Text})
		}
		// Output (AI) transcript.
		if sc.OutputTranscription != nil && sc.OutputTranscription.Text != "" {
			send(mockLiveServerMsg{Type: "transcript", Role: "ai", Text: sc.OutputTranscription.Text})
		}

		// Model audio + any text parts.
		if sc.ModelTurn != nil {
			for _, part := range sc.ModelTurn.Parts {
				if part == nil {
					continue
				}
				if part.InlineData != nil && len(part.InlineData.Data) > 0 {
					send(mockLiveServerMsg{
						Type: "audio",
						Data: base64.StdEncoding.EncodeToString(part.InlineData.Data),
					})
				}
			}
		}

		if sc.TurnComplete {
			send(mockLiveServerMsg{Type: "turn_complete"})
		}
	}
}

func writeMockLiveJSON(conn *websocket.Conn, m mockLiveServerMsg) {
	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_ = conn.WriteJSON(m)
	_ = conn.SetWriteDeadline(time.Time{})
}
