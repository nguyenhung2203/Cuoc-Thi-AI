import { ref } from 'vue'
import { useWebSocket } from './useWebSocket'

/**
 * useSpeechToText — live transcription for the REAL interview room (2 humans).
 *
 * Strategy (free tier, provider-swappable):
 *   Each browser transcribes ITS OWN microphone via the Web Speech API and
 *   pushes `transcript:partial` / `transcript:final` events up the existing
 *   realtime WebSocket. Because every client only ever transcribes the local
 *   speaker, speaker attribution is implicit — no server-side diarization
 *   needed. The gateway already resolves the speaker, broadcasts
 *   `transcript:update`, persists finals, and triggers AI suggestion/scoring.
 *
 * Swapping to a paid provider later (Google Cloud STT / Deepgram) means feeding
 * the SAME `transcript:partial/final` events (or the `/internal/rooms/{id}/
 * transcript` webhook) from a server-side ASR that subscribes to LiveKit
 * tracks. The gateway pipeline downstream does not change.
 *
 * Web Speech API availability: Chrome/Edge (webkitSpeechRecognition). Firefox
 * has no support; `supported` reports false so the UI can degrade gracefully.
 */
export function useSpeechToText() {
  const { sendMessage } = useWebSocket()

  const listening = ref(false)
  const supported = ref(
    typeof window !== 'undefined' &&
      (window.SpeechRecognition || window.webkitSpeechRecognition) != null,
  )
  const error = ref(null)

  let recognition = null
  let roomId = ''
  let interviewId = ''
  let speakerType = 'candidate'
  let speakerName = ''
  // We want the STT to run for the whole interview. Web Speech stops itself
  // periodically (after silence / ~60s), so we auto-restart while active.
  let active = false
  // start_time_ms per in-flight result index, so a partial and its later final
  // share the same value and the gateway dedupes them into one transcript row.
  const startByIndex = {}

  function buildRecognition() {
    const Ctor = window.SpeechRecognition || window.webkitSpeechRecognition
    const rec = new Ctor()
    rec.lang = 'vi-VN'
    rec.continuous = true
    rec.interimResults = true
    rec.maxAlternatives = 1

    rec.onresult = (event) => {
      for (let i = event.resultIndex; i < event.results.length; i++) {
        const result = event.results[i]
        const alt = result[0]
        if (!alt) continue
        const content = (alt.transcript || '').trim()
        if (!content) continue

        if (startByIndex[i] == null) startByIndex[i] = Date.now()
        const startMs = startByIndex[i]
        const confidence = typeof alt.confidence === 'number' && alt.confidence > 0 ? alt.confidence : 0.9

        const payload = {
          speaker_type: speakerType,
          speaker_name: speakerName,
          content,
          start_time_ms: startMs,
          end_time_ms: Date.now(),
          confidence,
        }

        if (result.isFinal) {
          sendMessage('transcript:final', payload, roomId, interviewId)
          delete startByIndex[i]
        } else {
          sendMessage('transcript:partial', payload, roomId, interviewId)
        }
      }
    }

    rec.onerror = (e) => {
      // 'no-speech' and 'aborted' are benign and self-heal via onend restart.
      if (e.error && e.error !== 'no-speech' && e.error !== 'aborted') {
        error.value = e.error
        console.warn('[STT] recognition error:', e.error)
      }
      if (e.error === 'not-allowed' || e.error === 'service-not-allowed') {
        // Permission denied — do not fight the browser with restarts.
        active = false
        listening.value = false
      }
    }

    rec.onend = () => {
      // Auto-restart to keep transcribing for the whole session.
      if (active) {
        try {
          rec.start()
        } catch {
          // start() throws if called too soon; retry shortly.
          setTimeout(() => {
            if (active) {
              try { rec.start() } catch {}
            }
          }, 300)
        }
      } else {
        listening.value = false
      }
    }

    return rec
  }

  /**
   * start begins local transcription and wires events to the given room.
   * @param {object} opts
   * @param {string} opts.roomId
   * @param {string} opts.interviewId
   * @param {'candidate'|'recruiter'} opts.speakerType — who this browser's mic belongs to
   * @param {string} [opts.speakerName]
   */
  function start({ roomId: rid, interviewId: iid, speakerType: stype, speakerName: sname } = {}) {
    if (!supported.value) {
      error.value = 'unsupported'
      return false
    }
    roomId = rid || ''
    interviewId = iid || ''
    speakerType = stype || 'candidate'
    speakerName = sname || ''
    error.value = null

    if (!recognition) recognition = buildRecognition()
    active = true
    try {
      recognition.start()
      listening.value = true
    } catch (e) {
      // Already started — treat as active.
      listening.value = true
    }
    return true
  }

  function stop() {
    active = false
    for (const k in startByIndex) delete startByIndex[k]
    if (recognition) {
      try { recognition.stop() } catch {}
    }
    listening.value = false
  }

  return { listening, supported, error, start, stop }
}
