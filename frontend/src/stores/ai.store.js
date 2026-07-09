import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useWebSocket } from '../composables/useWebSocket';

export const useAiStore = defineStore('ai', () => {
  const { sendMessage, on, off } = useWebSocket();

  const suggestions = ref([]);
  const scores = ref([]);
  const isThinking = ref(null); // null or string task name (e.g. 'suggest_follow_up')
  const aiError = ref(null);
  const aiWarning = ref(null);
  const aiThinkingMessage = ref('');
  const aiThinkingTimer = ref(null);
  const aiHardTimeoutTimer = ref(null);

  const setupListeners = () => {
    on('ai:thinking', handleAiThinking);
    on('ai:suggestion', handleAiSuggestion);
    on('ai:score_update', handleAiScoreUpdate);
    on('ai:error', handleAiError);
    on('ai:warning', handleAiWarning);
  };

  const cleanupListeners = () => {
    off('ai:thinking', handleAiThinking);
    off('ai:suggestion', handleAiSuggestion);
    off('ai:score_update', handleAiScoreUpdate);
    off('ai:error', handleAiError);
    off('ai:warning', handleAiWarning);
    clearAllTimers();
  };

  const handleAiThinking = (envelope) => {
    const payload = envelope.payload;
    isThinking.value = payload.task;
    aiThinkingMessage.value = payload.message || 'AI đang phân tích...';
    aiError.value = null; // Clear any previous error

    // Setup timeouts based on REALTIME_EVENTS.md rules
    clearAllTimers();
    
    // Soft timeout
    const expectedDuration = payload.expected_duration_ms || 3000;
    const softTimeout = expectedDuration * 2;
    aiThinkingTimer.value = setTimeout(() => {
      aiThinkingMessage.value = 'AI đang xử lý, chờ thêm chút...';
    }, softTimeout);

    // Hard timeout
    let hardTimeout = 15000;
    if (payload.task === 'score_answer') hardTimeout = 20000;
    if (payload.task === 'generate_summary') hardTimeout = 30000;

    aiHardTimeoutTimer.value = setTimeout(() => {
      isThinking.value = null;
      aiError.value = {
        severity: 'degraded',
        message: 'Kết quả AI bị trễ, vui lòng thử lại.'
      };
    }, hardTimeout);
  };

  const handleAiSuggestion = (envelope) => {
    clearAllTimers();
    isThinking.value = null;
    const existingIndex = suggestions.value.findIndex(s => s.suggestion_id === envelope.payload.suggestion_id);
    if (existingIndex >= 0) {
      suggestions.value[existingIndex] = envelope.payload;
    } else {
      suggestions.value.push(envelope.payload);
    }
  };

  const handleAiScoreUpdate = (envelope) => {
    clearAllTimers();
    isThinking.value = null;
    // Cập nhật điểm, payload.scores là mảng các tiêu chí
    const newScores = envelope.payload.scores || [];
    newScores.forEach(newScore => {
      const existingIndex = scores.value.findIndex(s => s.criterion_name === newScore.criterion_name);
      if (existingIndex >= 0) {
        scores.value[existingIndex] = newScore;
      } else {
        scores.value.push(newScore);
      }
    });
  };

  const handleAiError = (envelope) => {
    clearAllTimers();
    isThinking.value = null;
    aiError.value = envelope.payload;
    
    // Nếu recoverable = true thì frontend tự retry (có thể implement logic retry ở component hoặc đây)
    // if (envelope.payload.recoverable) { ... }
  };

  const handleAiWarning = (envelope) => {
    aiWarning.value = envelope.payload;
  };

  const requestSuggestion = (focus, last_transcript_id, roomId, interviewId) => {
    sendMessage('ai:request_suggestion', { focus, last_transcript_id }, roomId, interviewId);
  };
  
  const requestScoreUpdate = (criterion_ids, scope = 'latest_answer', roomId, interviewId) => {
    sendMessage('ai:request_score_update', { criterion_ids, scope }, roomId, interviewId);
  };

  const clearAllTimers = () => {
    if (aiThinkingTimer.value) clearTimeout(aiThinkingTimer.value);
    if (aiHardTimeoutTimer.value) clearTimeout(aiHardTimeoutTimer.value);
    aiThinkingTimer.value = null;
    aiHardTimeoutTimer.value = null;
  };

  const clearData = () => {
    suggestions.value = [];
    scores.value = [];
    isThinking.value = null;
    aiError.value = null;
    aiWarning.value = null;
    clearAllTimers();
  };

  return {
    suggestions,
    scores,
    isThinking,
    aiThinkingMessage,
    aiError,
    aiWarning,
    setupListeners,
    cleanupListeners,
    requestSuggestion,
    requestScoreUpdate,
    clearData
  };
});
