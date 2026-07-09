import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useWebSocket } from '../composables/useWebSocket';

export const useRoomStore = defineStore('room', () => {
  const { connect, disconnect, sendMessage, on, off, isConnected, connectionError } = useWebSocket();

  // State
  const roomId = ref(null);
  const interviewId = ref(null);
  const status = ref('waiting'); // waiting, active, paused, completed, cancelled, expired
  const participants = ref([]);
  const myParticipantId = ref(null);
  const startedAt = ref(null); // ISO timestamp khi phỏng vấn bắt đầu (dùng tính timer)
  const endedAt = ref(null);
  const reportStatus = ref(null); // 'generating', 'ready', 'failed'
  const cancelReason = ref(null);

  // Heartbeat interval reference (để cleanup)
  let heartbeatInterval = null;

  // =====================================================================
  // Actions
  // =====================================================================

  const connectRoom = (token, room_id, interview_id) => {
    roomId.value = room_id;
    interviewId.value = interview_id;
    connect(token);
    
    // Đăng ký nhận sự kiện sau khi kết nối
    setupListeners();
    
    // [FIX #1] Gửi heartbeat định kỳ mỗi 20 giây (spec yêu cầu 15-30s)
    startHeartbeat();
  };

  const disconnectRoom = () => {
    stopHeartbeat();
    cleanupListeners();
    disconnect();
    // Reset state
    roomId.value = null;
    interviewId.value = null;
    participants.value = [];
    myParticipantId.value = null;
    status.value = 'waiting';
    startedAt.value = null;
    endedAt.value = null;
    reportStatus.value = null;
    cancelReason.value = null;
  };

  // =====================================================================
  // [FIX #1] Heartbeat — gửi room:heartbeat mỗi 20s để server giữ online
  // =====================================================================

  const startHeartbeat = () => {
    stopHeartbeat(); // clear cũ nếu có
    heartbeatInterval = setInterval(() => {
      if (isConnected.value && roomId.value) {
        sendMessage('room:heartbeat', { connection_state: 'online' }, roomId.value, interviewId.value);
      }
    }, 20000); // 20 giây
  };

  const stopHeartbeat = () => {
    if (heartbeatInterval) {
      clearInterval(heartbeatInterval);
      heartbeatInterval = null;
    }
  };

  // =====================================================================
  // Listener Setup / Cleanup
  // =====================================================================

  const setupListeners = () => {
    on('room:joined', handleRoomJoined);
    on('room:user_joined', handleUserJoined);
    on('room:user_left', handleUserLeft);
    on('room:presence_update', handlePresenceUpdate);
    on('interview:started', handleInterviewStarted);
    on('interview:completed', handleInterviewCompleted);
    on('media:status_changed', handleMediaStatusChanged);
    // [FIX #2] Thêm listeners cho expired, paused, resumed, cancelled
    on('room:expired', handleRoomExpired);
    on('interview:paused', handleInterviewPaused);
    on('interview:resumed', handleInterviewResumed);
    on('interview:cancelled', handleInterviewCancelled);
    // Generic error listener
    on('error', handleError);
    // Report ready
    on('report:ready', handleReportReady);
  };

  const cleanupListeners = () => {
    off('room:joined', handleRoomJoined);
    off('room:user_joined', handleUserJoined);
    off('room:user_left', handleUserLeft);
    off('room:presence_update', handlePresenceUpdate);
    off('interview:started', handleInterviewStarted);
    off('interview:completed', handleInterviewCompleted);
    off('media:status_changed', handleMediaStatusChanged);
    // [FIX #2] Cleanup thêm
    off('room:expired', handleRoomExpired);
    off('interview:paused', handleInterviewPaused);
    off('interview:resumed', handleInterviewResumed);
    off('interview:cancelled', handleInterviewCancelled);
    off('error', handleError);
    off('report:ready', handleReportReady);
  };

  // =====================================================================
  // Event Handlers
  // =====================================================================

  const handleRoomJoined = (envelope) => {
    myParticipantId.value = envelope.payload.participant_id;
    status.value = envelope.payload.room_status;
    participants.value = envelope.payload.participants || [];
  };

  const handleUserJoined = (envelope) => {
    const existingIndex = participants.value.findIndex(p => p.participant_id === envelope.payload.participant_id);
    if (existingIndex >= 0) {
      participants.value[existingIndex] = { ...participants.value[existingIndex], ...envelope.payload };
    } else {
      participants.value.push({
        participant_id: envelope.payload.participant_id,
        display_name: envelope.payload.display_name,
        participant_type: envelope.payload.participant_type,
        connection_state: 'online',
        media_status: { mic_enabled: false, camera_enabled: false, screen_sharing: false }
      });
    }
  };

  const handleUserLeft = (envelope) => {
    const p = participants.value.find(p => p.participant_id === envelope.payload.participant_id);
    if (p) p.connection_state = 'offline';
  };

  const handlePresenceUpdate = (envelope) => {
    participants.value = envelope.payload.participants || [];
  };

  const handleInterviewStarted = (envelope) => {
    status.value = envelope.payload.status;
    startedAt.value = envelope.payload.started_at || new Date().toISOString();
  };

  const handleInterviewCompleted = (envelope) => {
    status.value = envelope.payload.status;
    endedAt.value = envelope.payload.ended_at || new Date().toISOString();
    reportStatus.value = envelope.payload.report_status || null;
  };

  const handleMediaStatusChanged = (envelope) => {
    const p = participants.value.find(p => p.participant_id === envelope.payload.participant_id);
    if (p) {
      p.media_status = envelope.payload;
    }
  };

  // [FIX #2] — Xử lý room:expired
  const handleRoomExpired = (envelope) => {
    status.value = 'expired';
    console.warn('[Room] Room đã hết hạn:', envelope.payload);
    // Component sẽ watch status === 'expired' để redirect user
  };

  // [FIX #3] — Xử lý interview:paused
  const handleInterviewPaused = (envelope) => {
    status.value = envelope.payload.status || 'paused';
  };

  // [FIX #3] — Xử lý interview:resumed
  const handleInterviewResumed = (envelope) => {
    status.value = envelope.payload.status || 'active';
  };

  // [FIX #3] — Xử lý interview:cancelled
  const handleInterviewCancelled = (envelope) => {
    status.value = envelope.payload.status || 'cancelled';
    cancelReason.value = envelope.payload.reason || null;
  };

  // Generic error handler
  const handleError = (envelope) => {
    console.error('[Room] Server error:', envelope.payload);
  };

  // Report ready handler
  const handleReportReady = (envelope) => {
    reportStatus.value = 'ready';
    console.log('[Room] Report ready:', envelope.payload);
  };

  // =====================================================================
  // Commands to Server
  // =====================================================================

  const startInterview = () => {
    sendMessage('interview:start', { consent_recording: true, consent_ai: true }, roomId.value, interviewId.value);
  };

  const endInterview = () => {
    sendMessage('interview:end', { generate_report: true }, roomId.value, interviewId.value);
  };

  // [FIX #3] — Thêm actions gửi pause, resume, cancel
  const pauseInterview = () => {
    sendMessage('interview:pause', {}, roomId.value, interviewId.value);
  };

  const resumeInterview = () => {
    sendMessage('interview:resume', {}, roomId.value, interviewId.value);
  };

  const cancelInterview = (reason = '') => {
    sendMessage('interview:cancel', { reason }, roomId.value, interviewId.value);
  };
  
  const updateMediaStatus = (mic_enabled, camera_enabled, screen_sharing = false) => {
    sendMessage('media:status', { mic_enabled, camera_enabled, screen_sharing }, roomId.value, interviewId.value);
  };

  const createNote = (content, tags = []) => {
    sendMessage('note:create', { content, tags }, roomId.value, interviewId.value);
  };

  const markQuestionAsked = (questionId, askedAtMs) => {
    sendMessage('question:mark_asked', { question_id: questionId, asked_at_ms: askedAtMs }, roomId.value, interviewId.value);
  };

  return {
    roomId,
    interviewId,
    status,
    participants,
    myParticipantId,
    startedAt,
    endedAt,
    reportStatus,
    cancelReason,
    isConnected,
    connectionError,
    connectRoom,
    disconnectRoom,
    startInterview,
    endInterview,
    pauseInterview,
    resumeInterview,
    cancelInterview,
    updateMediaStatus,
    createNote,
    markQuestionAsked
  };
});
