import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useWebSocket } from '../composables/useWebSocket';

export const useChatStore = defineStore('chat', () => {
  const { sendMessage, on, off } = useWebSocket();

  const messages = ref([]);

  const setupListeners = () => {
    on('chat:message', handleChatMessage);
  };

  const cleanupListeners = () => {
    off('chat:message', handleChatMessage);
  };

  const handleChatMessage = (envelope) => {
    // Check for duplicates
    if (!messages.value.find(m => m.message_id === envelope.payload.message_id)) {
      messages.value.push(envelope.payload);
      // Giữ khoảng 200 tin nhắn gần nhất để tránh tràn RAM
      if (messages.value.length > 200) {
        messages.value.shift();
      }
    }
  };

  const sendChat = (message, visibility = 'room', roomId, interviewId) => {
    // visibility có thể là 'room' hoặc 'recruiter_only'
    sendMessage('chat:send', { message, visibility }, roomId, interviewId);
  };

  const setHistory = (historyArr) => {
    messages.value = historyArr || [];
  };

  const clearChat = () => {
    messages.value = [];
  };

  return {
    messages,
    setupListeners,
    cleanupListeners,
    sendChat,
    setHistory,
    clearChat
  };
});
