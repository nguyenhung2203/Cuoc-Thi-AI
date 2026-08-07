import { ref, onMounted, onUnmounted, unref } from 'vue'
import { useWebSocket } from './useWebSocket'

export function useRoom(roomId, interviewId) {
  const { isConnected, sendMessage, on, off } = useWebSocket()
  
  const participants = ref([])
  const currentRoom = ref(null)
  const isJoined = ref(false)

  // Event handlers
  const handleRoomJoined = (envelope) => {
    isJoined.value = true
    currentRoom.value = envelope.payload.room
    participants.value = envelope.payload.participants || []
    console.log('[Room] Successfully joined room', envelope.payload)
  }

  const handleUserJoined = (envelope) => {
    const user = envelope.payload.participant
    if (!participants.value.find(p => p.id === user.id)) {
      participants.value.push(user)
    }
    console.log('[Room] User joined', user)
  }

  const handleUserLeft = (envelope) => {
    const userId = envelope.payload.participant_id
    participants.value = participants.value.filter(p => p.id !== userId)
    console.log('[Room] User left', userId)
  }

  const handlePresenceUpdate = (envelope) => {
    const updates = envelope.payload.updates || []
    updates.forEach(update => {
      const idx = participants.value.findIndex(p => p.id === update.participant_id)
      if (idx !== -1) {
        participants.value[idx].status = update.status
      }
    })
  }

  const handleDisconnect = () => {
    isJoined.value = false
    participants.value = []
    currentRoom.value = null
  }

  const joinRoom = () => {
    if (!isConnected.value) {
      console.warn('[Room] Cannot join, WebSocket not connected')
      return
    }
    sendMessage('room:join', {
      client_version: '1.0.0',
      capabilities: ['audio', 'video', 'chat']
    }, unref(roomId), unref(interviewId))
  }

  onMounted(() => {
    on('room:joined', handleRoomJoined)
    on('room:user_joined', handleUserJoined)
    on('room:user_left', handleUserLeft)
    on('room:presence_update', handlePresenceUpdate)
  })

  onUnmounted(() => {
    off('room:joined', handleRoomJoined)
    off('room:user_joined', handleUserJoined)
    off('room:user_left', handleUserLeft)
    off('room:presence_update', handlePresenceUpdate)
    handleDisconnect()
  })

  return {
    isJoined,
    currentRoom,
    participants,
    joinRoom
  }
}
