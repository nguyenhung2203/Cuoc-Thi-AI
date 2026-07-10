import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useWebSocket } from '../composables/useWebSocket';
import { useRoomStore } from './room.store';

export const useChatStore = defineStore('chat', () => {
  const { sendMessage, on, off } = useWebSocket();

  // Lưu trữ tin nhắn theo roomId
  const messagesByRoom = ref(JSON.parse(sessionStorage.getItem('interview_chat_history_map') || '{}'));

  const saveToStorage = () => {
    try {
      sessionStorage.setItem('interview_chat_history_map', JSON.stringify(messagesByRoom.value));
    } catch (e) { /* ignore */ }
  };

  // Tính toán roomId hiện tại từ roomStore
  const currentRoomId = computed(() => {
    const roomStore = useRoomStore();
    return roomStore.roomId;
  });

  // Lấy danh sách tin nhắn của phòng hiện tại
  const messages = computed(() => {
    if (!currentRoomId.value) return [];
    return messagesByRoom.value[currentRoomId.value] || [];
  });

  const setupListeners = () => {
    on('chat:message', handleChatMessage);
  };

  const cleanupListeners = () => {
    off('chat:message', handleChatMessage);
  };

  const handleChatMessage = (envelope) => {
    const payload = envelope.payload;
    if (!payload) return;
    
    const targetRoomId = payload.room_id || currentRoomId.value;
    if (!targetRoomId) return;

    if (!messagesByRoom.value[targetRoomId]) {
      messagesByRoom.value[targetRoomId] = [];
    }

    const roomMsgs = messagesByRoom.value[targetRoomId];

    // Check for exact ID match or recent local optimistic match with exact same text
    const existingIndex = roomMsgs.findIndex(m => 
      m.message_id === payload.message_id || 
      (m.sender_type === 'local' && m.message === payload.message && (Date.now() - new Date(m.created_at || Date.now()).getTime() < 10000))
    );
    
    if (existingIndex !== -1) {
      // Update existing optimistic message with server fields
      roomMsgs[existingIndex] = { ...roomMsgs[existingIndex], ...payload };
    } else {
      roomMsgs.push(payload);
      if (roomMsgs.length > 200) {
        roomMsgs.shift();
      }
    }
    saveToStorage();
  };

  const sendChat = (message, visibility = 'room', roomId, interviewId, senderName = 'Bạn') => {
    if (!message || typeof message !== 'string' || !message.trim()) return;
    const msgId = 'msg-' + Math.random().toString(36).substring(2, 9) + '-' + Date.now();
    
    const targetRoomId = roomId || currentRoomId.value;
    if (!targetRoomId) return;

    if (!messagesByRoom.value[targetRoomId]) {
      messagesByRoom.value[targetRoomId] = [];
    }
    const roomMsgs = messagesByRoom.value[targetRoomId];

    // Optimistic update để hiển thị ngay trên UI
    if (!roomMsgs.find(m => m.message_id === msgId)) {
      roomMsgs.push({
        message_id: msgId,
        sender_name: senderName,
        sender_type: 'local',
        message: message.trim(),
        visibility: visibility || 'room',
        created_at: new Date().toISOString(),
        room_id: targetRoomId
      });
      if (roomMsgs.length > 200) {
        roomMsgs.shift();
      }
      saveToStorage();
    }
    sendMessage('chat:send', { message: message.trim(), visibility, message_id: msgId }, targetRoomId, interviewId);
  };

  const setHistory = (historyArr, roomId) => {
    const targetRoomId = roomId || currentRoomId.value;
    if (!targetRoomId) return;
    messagesByRoom.value[targetRoomId] = historyArr || [];
    saveToStorage();
  };

  const clearChat = (roomId) => {
    const targetRoomId = roomId || currentRoomId.value;
    if (targetRoomId) {
      messagesByRoom.value[targetRoomId] = [];
      saveToStorage();
    }
  };

  // Dọn dẹp storage cũ nếu tồn tại
  try {
    sessionStorage.removeItem('interview_chat_history');
  } catch (e) { /* ignore */ }

  return {
    messagesByRoom, // Cần return để pinia persist hoạt động đúng
    messages,
    setupListeners,
    cleanupListeners,
    sendChat,
    setHistory,
    clearChat
  };
}, {
  persist: {
    paths: ['messagesByRoom']
  }
});
