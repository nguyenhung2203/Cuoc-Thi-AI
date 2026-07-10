import { ref } from 'vue'

/**
 * useGeminiLive — drives a spoken, real-time mock interview through the
 * backend Gemini Live proxy at ws(s)://<gateway>/ws/mock-live.
 *
 * Responsibilities:
 *  - open the WebSocket (role/level as query params)
 *  - capture mic audio as PCM16 @16kHz via an AudioWorklet and stream it up
 *  - play PCM16 @24kHz audio chunks the AI sends back, gaplessly
 *  - surface live transcripts (candidate + AI) and speaking state
 *
 * The Gemini API key lives on the server; the browser only talks to the proxy.
 */
export function useGeminiLive() {
  const connected = ref(false)
  const aiSpeaking = ref(false)
  const listening = ref(false)
  const error = ref(null)
  // transcript entries: { role: 'ai'|'candidate', text }
  const transcript = ref([])
  // 0..1 loudness of the AI voice, drives avatar lip-sync (free, no 3rd party).
  const audioLevel = ref(0)

  let ws = null
  let micCtx = null
  let playCtx = null
  let micNode = null
  let playNode = null
  let micStream = null
  let analyser = null
  let levelRAF = null

  const wsBase = () => {
    let base = import.meta.env.VITE_WS_URL
    if (!base) {
      const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const host = window.location.hostname === 'localhost' ? 'localhost:8081' : window.location.host
      base = `${proto}//${host}`
    }
    return base.replace(/\/$/, '')
  }

  // Append transcript text, merging consecutive same-role fragments into one
  // bubble so streaming partials read naturally.
  const pushTranscript = (role, text) => {
    if (!text) return
    const last = transcript.value[transcript.value.length - 1]
    if (last && last.role === role && !last.done) {
      last.text += text
    } else {
      transcript.value.push({ role, text, done: false })
    }
  }

  const markTurnDone = () => {
    const last = transcript.value[transcript.value.length - 1]
    if (last) last.done = true
  }

  async function start({ role, level } = {}) {
    error.value = null
    transcript.value = []
    try {
      await setupAudio()
    } catch (e) {
      error.value = 'Không truy cập được micro. Hãy cấp quyền micro và thử lại.'
      throw e
    }

    const q = new URLSearchParams()
    if (role) q.set('role', role)
    if (level) q.set('level', level)
    ws = new WebSocket(`${wsBase()}/ws/mock-live?${q.toString()}`)
    ws.binaryType = 'arraybuffer'

    ws.onopen = () => { connected.value = true }
    ws.onerror = () => { error.value = 'Mất kết nối tới AI phỏng vấn.' }
    ws.onclose = () => { connected.value = false; listening.value = false }
    ws.onmessage = onServerMessage
  }

  function onServerMessage(ev) {
    let msg
    try { msg = JSON.parse(ev.data) } catch { return }
    switch (msg.type) {
      case 'ready':
        listening.value = true
        break
      case 'audio':
        playChunk(msg.data)
        break
      case 'transcript':
        pushTranscript(msg.role, msg.text)
        break
      case 'turn_complete':
        markTurnDone()
        break
      case 'interrupted':
        if (playNode) playNode.port.postMessage('flush')
        aiSpeaking.value = false
        break
      case 'error':
        error.value = msg.message || 'Lỗi AI'
        break
    }
  }

  async function setupAudio() {
    // --- Playback context @24kHz (AI output sample rate) ---
    playCtx = new (window.AudioContext || window.webkitAudioContext)({ sampleRate: 24000 })
    await playCtx.audioWorklet.addModule('/pcm-player-processor.js')
    playNode = new AudioWorkletNode(playCtx, 'pcm-player-processor')
    playNode.port.onmessage = (e) => {
      aiSpeaking.value = e.data === 'speaking'
      if (!aiSpeaking.value) audioLevel.value = 0
    }
    // Analyser taps the AI voice to drive avatar lip-sync (RMS loudness 0..1).
    analyser = playCtx.createAnalyser()
    analyser.fftSize = 256
    playNode.connect(analyser)
    analyser.connect(playCtx.destination)
    startLevelLoop()

    // --- Capture context (native rate; worklet resamples to 16kHz) ---
    micStream = await navigator.mediaDevices.getUserMedia({
      audio: { channelCount: 1, echoCancellation: true, noiseSuppression: true },
    })
    micCtx = new (window.AudioContext || window.webkitAudioContext)()
    await micCtx.audioWorklet.addModule('/mic-capture-processor.js')
    const src = micCtx.createMediaStreamSource(micStream)
    micNode = new AudioWorkletNode(micCtx, 'mic-capture-processor')
    micNode.port.onmessage = (e) => {
      // e.data is a PCM16 ArrayBuffer @16kHz — forward as base64.
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'audio', data: arrayBufferToBase64(e.data) }))
      }
    }
    src.connect(micNode)
    // Keep the graph alive without echoing mic to speakers.
    const sink = micCtx.createGain()
    sink.gain.value = 0
    micNode.connect(sink)
    sink.connect(micCtx.destination)
  }

  function playChunk(b64) {
    if (!playNode) return
    const buf = base64ToArrayBuffer(b64)
    playNode.port.postMessage(buf, [buf])
  }

  // Continuously measure AI-voice loudness so the avatar mouth can lip-sync.
  function startLevelLoop() {
    const data = new Uint8Array(analyser.frequencyBinCount)
    const tick = () => {
      if (!analyser) return
      analyser.getByteTimeDomainData(data)
      let sum = 0
      for (let i = 0; i < data.length; i++) {
        const v = (data[i] - 128) / 128
        sum += v * v
      }
      const rms = Math.sqrt(sum / data.length)
      // Smooth and boost for a lively mouth without clipping.
      const target = aiSpeaking.value ? Math.min(1, rms * 3.2) : 0
      audioLevel.value = audioLevel.value * 0.6 + target * 0.4
      levelRAF = requestAnimationFrame(tick)
    }
    levelRAF = requestAnimationFrame(tick)
  }

  // Tell the AI the candidate finished speaking (commit the turn).
  function commitTurn() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'end' }))
    }
  }

  // Optional typed fallback.
  function sendText(text) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'text', text }))
      pushTranscript('candidate', text)
      markTurnDone()
    }
  }

  function stop() {
    try { levelRAF && cancelAnimationFrame(levelRAF) } catch {}
    try { ws && ws.close() } catch {}
    try { micStream && micStream.getTracks().forEach(t => t.stop()) } catch {}
    try { micNode && micNode.disconnect() } catch {}
    try { analyser && analyser.disconnect() } catch {}
    try { playNode && playNode.disconnect() } catch {}
    try { micCtx && micCtx.close() } catch {}
    try { playCtx && playCtx.close() } catch {}
    ws = micCtx = playCtx = micNode = playNode = micStream = analyser = levelRAF = null
    connected.value = false
    listening.value = false
    aiSpeaking.value = false
    audioLevel.value = 0
  }

  return {
    connected, aiSpeaking, listening, error, transcript, audioLevel,
    start, stop, commitTurn, sendText,
  }
}

function arrayBufferToBase64(buf) {
  let binary = ''
  const bytes = new Uint8Array(buf)
  const chunk = 0x8000
  for (let i = 0; i < bytes.length; i += chunk) {
    binary += String.fromCharCode.apply(null, bytes.subarray(i, i + chunk))
  }
  return btoa(binary)
}

function base64ToArrayBuffer(b64) {
  const binary = atob(b64)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i)
  return bytes.buffer
}
