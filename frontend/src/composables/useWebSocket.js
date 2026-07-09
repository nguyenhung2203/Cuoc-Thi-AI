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

export function useWebSocket() {
  const connect = (token) => {
    if (ws.value?.readyState === WebSocket.OPEN) return

    // TODO: Determine backend WS URL properly. 
    // Assuming backend runs on port 8080 locally for now, 
    // or proxy via Vite if configured (e.g. wss://domain.com/ws/interview-room)
    const host = window.location.hostname === 'localhost' ? 'localhost:8081' : window.location.host
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = `${protocol}//${host}/ws/interview-room?token=${token}`

    ws.value = new WebSocket(url)

    ws.value.onopen = () => {
      console.log('[WebSocket] Connected')
      isConnected.value = true
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
      
      // Auto-reconnect could be implemented here
      if (event.code !== 1000) { // 1000 is normal closure
        connectionError.value = 'Connection lost. Attempting to reconnect...'
        setTimeout(() => connect(token), 3000)
      }
    }

    ws.value.onerror = (error) => {
      console.error('[WebSocket] Error:', error)
      connectionError.value = 'WebSocket connection error'
    }
  }

  const disconnect = () => {
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
