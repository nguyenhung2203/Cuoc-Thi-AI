import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useWebSocket } from '../composables/useWebSocket';

export const useTranscriptStore = defineStore('transcript', () => {
  const { on, off } = useWebSocket();

  const transcripts = ref([]);
  const activePartials = ref({});

  const setupListeners = () => {
    on('transcript:update', handleTranscriptUpdate);
  };

  const cleanupListeners = () => {
    off('transcript:update', handleTranscriptUpdate);
  };

  const handleTranscriptUpdate = (envelope) => {
    const payload = envelope.payload;
    const { transcript_id, is_final } = payload;

    if (is_final) {
      // Remove from partials if it exists
      if (activePartials.value[transcript_id]) {
        delete activePartials.value[transcript_id];
      }
      
      // Check if it already exists in final transcripts
      const existingIndex = transcripts.value.findIndex(t => t.transcript_id === transcript_id);
      if (existingIndex >= 0) {
        transcripts.value[existingIndex] = payload;
      } else {
        transcripts.value.push(payload);
        // Sort by start time just in case they arrive out of order
        transcripts.value.sort((a, b) => a.start_time_ms - b.start_time_ms);
      }
    } else {
      // Store in active partials
      activePartials.value[transcript_id] = payload;
    }
  };

  const setHistory = (historyArr) => {
    transcripts.value = historyArr || [];
    // Sort
    transcripts.value.sort((a, b) => a.start_time_ms - b.start_time_ms);
  };

  const clearTranscripts = () => {
    transcripts.value = [];
    activePartials.value = {};
  };

  return {
    transcripts,
    activePartials,
    setupListeners,
    cleanupListeners,
    setHistory,
    clearTranscripts
  };
}, {
  persist: {
    paths: ['transcripts', 'activePartials']
  }
});
