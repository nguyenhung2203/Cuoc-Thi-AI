import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useWebSocket } from '../composables/useWebSocket';

export const useChatStore = defineStore('chat', () => {
  const { sendMessage, on, off } = useWebSocket();

  const messages = ref(JSON.parse(sessionStorage.getItem('interview_chat_history') || '[]'));

  const saveToStorage = () => {
    try {
      sessionStorage.setItem('interview_chat_history', JSON.stringify(messages.value));
    } catch (e) { /* ignore */ }
  };

  const setupListeners = () => {
    on('chat:message', handleChatMessage);
  };

  const cleanupListeners = () => {
    off('chat:message', handleChatMessage);
  };

  const handleChatMessage = (envelope) => {
    const payload = envelope.payload;
    if (!payload) return;
    // Check for exact ID match or recent local optimistic match with exact same text
    const existingIndex = messages.value.findIndex(m => 
      m.message_id === payload.message_id || 
      (m.sender_type === 'local' && m.message === payload.message && (Date.now() - new Date(m.created_at || Date.now()).getTime() < 10000))
    );
    if (existingIndex !== -1) {
      // Update existing optimistic message with server fields
      messages.value[existingIndex] = { ...messages.value[existingIndex], ...payload };
    } else {
      messages.value.push(payload);
      if (messages.value.length > 200) {
        messages.value.shift();
      }
    }
    saveToStorage();
  };

  const sendChat = (message, visibility = 'room', roomId, interviewId, senderName = 'Bạn') => {
    if (!message || typeof message !== 'string' || !message.trim()) return;
    const msgId = 'msg-' + Math.random().toString(36).substring(2, 9) + '-' + Date.now();
    // Optimistic update để hiển thị ngay trên UI
    if (!messages.value.find(m => m.message_id === msgId)) {
      messages.value.push({
        message_id: msgId,
        sender_name: senderName,
        sender_type: 'local',
        message: message.trim(),
        visibility: visibility || 'room',
        created_at: new Date().toISOString()
      });
      if (messages.value.length > 200) {
        messages.value.shift();
      }
      saveToStorage();
    }
    sendMessage('chat:send', { message: message.trim(), visibility, message_id: msgId }, roomId, interviewId);
  };

  const setHistory = (historyArr) => {
    messages.value = historyArr || [];
    saveToStorage();
  };

  const clearChat = () => {
    messages.value = [];
    sessionStorage.removeItem('interview_chat_history');
  };

  return {
    messages,
    setupListeners,
    cleanupListeners,
    sendChat,
    setHistory,
    clearChat
  };
}, {
  persist: {
    paths: ['messages']
  }
});
