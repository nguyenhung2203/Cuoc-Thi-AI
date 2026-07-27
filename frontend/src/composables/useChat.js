import { ref, onMounted, onUnmounted } from 'vue'
import { useWebSocket } from './useWebSocket'

export function useChat(roomId, interviewId) {
  const { isConnected, sendMessage, on, off } = useWebSocket()
  const messages = ref([])

  const handleIncomingMessage = (envelope) => {
    const msg = envelope.payload
    messages.value.push({
      id: msg.id || Date.now().toString(),
      sender_id: msg.sender_id,
      sender_name: msg.sender_name,
      sender_role: msg.sender_role,
      content: msg.content,
      timestamp: msg.timestamp || new Date().toISOString(),
      visibility: msg.visibility || 'room'
    })
  }

  const sendChatMessage = (content, visibility = 'room') => {
    if (!isConnected.value) {
      console.warn('[Chat] Cannot send message, WebSocket not connected')
      return false
    }

    // Optimistic update could be added here if we had a temporary ID, 
    // but typically we wait for the server broadcast to maintain order.
    // However, if the server doesn't broadcast back to the sender, we should add it.
    // For now we assume the server broadcasts to everyone including sender.

    return sendMessage('chat:send', {
      content,
      visibility
    }, roomId, interviewId)
  }

  onMounted(() => {
    on('chat:message', handleIncomingMessage)
  })

  onUnmounted(() => {
    off('chat:message', handleIncomingMessage)
  })

  return {
    messages,
    sendChatMessage
  }
}
