import { ref } from 'vue'

// Event emitter for pub/sub pattern
class EventEmitter {
  constructor() {
    this.listeners = new Map()
  }

  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set())
    }
    this.listeners.get(event).add(callback)
  }

  off(event, callback) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).delete(callback)
    }
  }

  emit(event, data) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => callback(data))
    }
  }
}

const ws = ref(null)
const isConnected = ref(false)
const connectionError = ref(null)
const emitter = new EventEmitter()
let reconnectTimer = null
let reconnectAttempts = 0
let shouldReconnect = false
let lastToken = ''

export function useWebSocket() {
  const connect = (token) => {
    if (ws.value?.readyState === WebSocket.OPEN || ws.value?.readyState === WebSocket.CONNECTING) return
    shouldReconnect = true
    lastToken = token

    // Resolve the realtime gateway base from env (VITE_WS_URL), e.g.
    // ws://localhost:8081. Falls back to deriving from the current origin so
    // it still works behind a reverse proxy without extra config.
    let base = import.meta.env.VITE_WS_URL
    if (!base) {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      base = `${protocol}//${window.location.host}`
    }
    base = base.replace(/\/+$/, '').replace(/\/ws\/interview-room$/, '')
    const url = `${base}/ws/interview-room?token=${encodeURIComponent(token)}`

    ws.value = new WebSocket(url)

    ws.value.onopen = () => {
      console.log('[WebSocket] Connected')
      isConnected.value = true
      reconnectAttempts = 0
      connectionError.value = null
    }

    ws.value.onmessage = (event) => {
      try {
        const envelope = JSON.parse(event.data)
        // Emit the specific event
        emitter.emit(envelope.event, envelope)
        
        // Also emit a catch-all event for debugging
        emitter.emit('*', envelope)
      } catch (e) {
        console.error('[WebSocket] Failed to parse message:', e)
      }
    }

    ws.value.onclose = (event) => {
      console.log('[WebSocket] Disconnected', event.code, event.reason)
      isConnected.value = false
      ws.value = null
      
      if (event.code !== 1000 && shouldReconnect && reconnectAttempts < 8) {
        connectionError.value = 'Connection lost. Attempting to reconnect...'
        const delay = Math.min(30000, 1000 * (2 ** reconnectAttempts))
        reconnectAttempts += 1
        clearTimeout(reconnectTimer)
        reconnectTimer = setTimeout(() => connect(lastToken), delay)
      }
    }

    ws.value.onerror = (error) => {
      console.error('[WebSocket] Error:', error)
      connectionError.value = 'WebSocket connection error'
    }
  }

  const disconnect = () => {
    shouldReconnect = false
    clearTimeout(reconnectTimer)
    reconnectTimer = null
    reconnectAttempts = 0
    if (ws.value) {
      ws.value.close(1000, 'User navigating away')
      ws.value = null
      isConnected.value = false
    }
  }

  const sendMessage = (event, payload = {}, room_id = '', interview_id = '') => {
    if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
      console.warn('[WebSocket] Cannot send message, not connected')
      return false
    }

    const request_id = Math.random().toString(36).substring(2, 15)
    
    const envelope = {
      event,
      request_id,
      room_id,
      interview_id,
      payload
    }

    ws.value.send(JSON.stringify(envelope))
    return true
  }

  return {
    isConnected,
    connectionError,
    connect,
    disconnect,
    sendMessage,
    on: emitter.on.bind(emitter),
    off: emitter.off.bind(emitter)
  }
}
