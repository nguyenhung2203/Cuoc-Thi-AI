<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter, useRoute, onBeforeRouteLeave } from 'vue-router'
import Button from '../../components/common/AppButton.vue'
import Card from '../../components/common/AppCard.vue'
import Badge from '../../components/common/AppBadge.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mic, MicOff, Video, VideoOff, MonitorUp, MessageSquare, PhoneOff, Sparkles, CheckCircle, AlertTriangle, FileText, ChevronRight, ChevronLeft, ChevronsLeft, ChevronsRight, Send, LogOut, Play, Pause } from 'lucide-vue-next'
import { useLiveKit } from '../../composables/useLiveKit'
import { useSpeechToText } from '../../composables/useSpeechToText'
import { roomService } from '../../services/room.service'
import { authStore } from '../../stores/auth.store'
import { rubricService } from '../../services/rubric.service'

// Import Stores
import { useRoomStore } from '../../stores/room.store'
import { useChatStore } from '../../stores/chat.store'
import { useTranscriptStore } from '../../stores/transcript.store'
import { useAiStore } from '../../stores/ai.store'

const router = useRouter()
const route = useRoute()

// Stores
const roomStore = useRoomStore()
const chatStore = useChatStore()
const transcriptStore = useTranscriptStore()
const aiStore = useAiStore()
const rubricCriteria = ref([])

const {
  isConnected: isLiveKitConnected, error: liveKitError, isMicOn, isCameraOn, isScreenSharing,
  localVideoEl, remoteVideoEl,
  connectToRoom, toggleMic, toggleCamera, toggleScreenShare, disconnect: liveKitDisconnect,
  clearError
} = useLiveKit()

// Live transcription of THIS recruiter's mic (free Web Speech API). Pushes
// transcript:partial/final up the realtime WS; the gateway attributes the
// speaker, broadcasts, persists, and triggers AI suggestion/scoring.
const { supported: sttSupported, listening: sttListening, start: startSTT, stop: stopSTT } = useSpeechToText()

const syncSTT = () => {
  const shouldRun = roomStore.status === 'active' && isMicOn.value && sttSupported.value
  if (shouldRun && !sttListening.value) {
    startSTT({
      roomId: roomStore.roomId,
      interviewId: roomStore.interviewId,
      speakerType: 'recruiter',
      speakerName: authStore.user?.full_name || 'Nhà tuyển dụng',
    })
  } else if (!shouldRun && sttListening.value) {
    stopSTT()
  }
}

watch(() => [roomStore.status, isMicOn.value], syncSTT)

const activeTab = ref(sessionStorage.getItem('recruiter_active_tab') || 'assistant')
watch(activeTab, (val) => { sessionStorage.setItem('recruiter_active_tab', val); })

const isPanelExpanded = ref(sessionStorage.getItem('recruiter_panel_expanded') !== 'false')
watch(isPanelExpanded, (val) => { sessionStorage.setItem('recruiter_panel_expanded', val); })

const chatInput = ref('')
const chatCooldown = ref(false)

const handleSendMessage = () => {
  const message = chatInput.value.trim()
  if (!message || chatCooldown.value) return
  if (!roomStore.roomId || !roomStore.interviewId || !roomStore.isConnected) {
    activeToast.value = { type: 'error', message: 'Không thể gửi tin nhắn khi phòng chưa kết nối.' }
    return
  }
  if (message.length > 2000) {
    activeToast.value = { type: 'error', message: 'Tin nhắn không được vượt quá 2.000 ký tự.' }
    return
  }
  chatStore.sendChat(message, 'room', roomStore.roomId, roomStore.interviewId, 'Nhà tuyển dụng (Bạn)')
  chatInput.value = ''
  chatCooldown.value = true
  setTimeout(() => { chatCooldown.value = false }, 500)
}

const requestAIScore = async () => {
  const criterionIds = rubricCriteria.value.map(item => item.id).filter(Boolean)
  if (!criterionIds.length) {
    activeToast.value = { type: 'warning', message: 'Công việc chưa có tiêu chí rubric hợp lệ.' }
    return
  }
  await aiStore.requestScoreUpdate(criterionIds, 'latest_answer', roomStore.roomId, roomStore.interviewId)
}

const isEnding = ref(false)
const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)
const activeToast = ref(null)
const interviewId = route.params.interviewId || history.state?.interviewId || route.params.id || null

watch(() => aiStore.aiError, (err) => {
  if (err) {
    activeToast.value = {
      type: err.severity === 'critical' ? 'error' : 'warning',
      message: err.message || 'Lỗi hệ thống AI.'
    }
  }
})

watch(() => roomStore.reportStatus, (status) => {
  if (status === 'ready') {
    activeToast.value = {
      type: 'success',
      message: 'Báo cáo phỏng vấn đã được chuẩn bị thành công!'
    }
  } else if (status === 'failed') {
    activeToast.value = {
      type: 'error',
      message: 'Gặp lỗi khi tạo báo cáo phỏng vấn.'
    }
  }
})

onMounted(async () => {
  if (history.state?.message) {
    window.history.replaceState({ interviewId: history.state.interviewId }, document.title)
  }

  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (companyId) {
      const data = await rubricService.getRubrics(companyId)
      const rubrics = Array.isArray(data) ? data : (data?.items || data?.data || [])
      rubricCriteria.value = rubrics.flatMap(r => r.rubric_criteria || r.criteria || [])
    }
  } catch (error) {
    console.error('Không tải được rubric cho chấm điểm AI', error)
  }

  chatStore.setupListeners()
  transcriptStore.setupListeners()
  aiStore.setupListeners()
  
  if (interviewId) {
    try {
      const companyId = authStore.user?.companies?.[0]?.id
      if (companyId) {
        const response = await roomService.getRoomToken(companyId, interviewId)
        const token = response.token || response.livekit_token || response.room_access_token
        
        if (token) {
          // Connect WebSocket Realtime (Store)
          const roomId = response.room_id || `room-${interviewId}`
          roomStore.connectRoom(token, roomId, interviewId)

          // Gửi room:join sau khi WebSocket mở
          const tryJoinAsRecruiter = (retries = 10) => {
            if (roomStore.isConnected) {
              roomStore.sendRoomJoin('recruiter')
            } else if (retries > 0) {
              setTimeout(() => tryJoinAsRecruiter(retries - 1), 300)
            }
          }
          setTimeout(() => tryJoinAsRecruiter(), 500)

          // Connect LiveKit (hoặc native getUserMedia nếu không có server)
          const livekitUrl = import.meta.env.VITE_LIVEKIT_URL || `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}`
          await connectToRoom(livekitUrl, token)
        }
      }
    } catch (err) {
      console.error('Không thể lấy room token', err)
    }
  }
})

const showLeaveWarningModal = ref(false)
const pendingRouteResolve = ref(null)

const handleBeforeUnload = (e) => {
  if (!isEnding.value && roomStore.status !== 'completed' && roomStore.status !== 'expired' && roomStore.status !== 'cancelled') {
    e.preventDefault()
    e.returnValue = ''
  }
}

onMounted(() => {
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onUnmounted(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
  stopSTT()
  liveKitDisconnect()
  roomStore.disconnectRoom()
  chatStore.cleanupListeners()
  transcriptStore.cleanupListeners()
  aiStore.cleanupListeners()
})

onBeforeRouteLeave((to, from) => {
  if (isEnding.value || roomStore.status === 'completed' || roomStore.status === 'expired' || roomStore.status === 'cancelled') {
    return true
  }
  return new Promise((resolve) => {
    pendingRouteResolve.value = resolve
    showLeaveWarningModal.value = true
  })
})

const handleConfirmLeave = () => {
  showLeaveWarningModal.value = false
  if (pendingRouteResolve.value) {
    pendingRouteResolve.value(true)
    pendingRouteResolve.value = null
  }
}

const handleCancelLeave = () => {
  showLeaveWarningModal.value = false
  if (pendingRouteResolve.value) {
    pendingRouteResolve.value(false)
    pendingRouteResolve.value = null
  }
}

const handleEndCall = () => {
  showEndModal.value = true
}

const confirmEndCall = () => {
  showEndModal.value = false
  isEnding.value = true
  roomStore.endInterview() // Gửi event end

  setTimeout(() => {
    const target = interviewId ? `/interviews/${interviewId}/report` : '/reports'
    router.push({ path: target, state: { message: 'Đã lưu kết quả phỏng vấn thành công' } })
  }, 1500)
}

const handleStartCall = () => {
  roomStore.startInterview()
}

const handlePauseCall = () => {
  roomStore.pauseInterview()
}

const handleResumeCall = () => {
  roomStore.resumeInterview()
}

// Toggle chính cho nút Bắt đầu/Tạm dừng:
// - Nếu đang diễn ra (active) -> tạm dừng
// - Nếu đang tạm dừng (paused) -> tiếp tục
// - Ngược lại (waiting) -> bắt đầu
const handleToggleCallState = () => {
  if (roomStore.status === 'active') {
    handlePauseCall()
  } else if (roomStore.status === 'paused') {
    handleResumeCall()
  } else {
    handleStartCall()
  }
}

const handleCancelCall = () => {
  roomStore.cancelInterview('Recruiter hủy buổi phỏng vấn')
}

// [FIX] Watch status expired/cancelled → redirect user
watch(() => roomStore.status, (newStatus) => {
  if (newStatus === 'expired') {
    setTimeout(() => {
      router.push({ path: '/interviews', state: { message: 'Phòng phỏng vấn đã hết hạn.' } })
    }, 2000)
  }
  if (newStatus === 'cancelled') {
    setTimeout(() => {
      router.push({ path: '/interviews', state: { message: 'Buổi phỏng vấn đã bị hủy.' } })
    }, 2000)
  }
})

// [FIX] Live elapsed timer from startedAt
const elapsedSeconds = ref(0)
let timerInterval = null

watch(() => roomStore.startedAt, (val) => {
  if (val) {
    if (timerInterval) clearInterval(timerInterval)
    timerInterval = setInterval(() => {
      elapsedSeconds.value = Math.floor((Date.now() - new Date(val).getTime()) / 1000)
    }, 1000)
  }
})

const elapsedFormatted = computed(() => {
  const h = Math.floor(elapsedSeconds.value / 3600)
  const m = Math.floor((elapsedSeconds.value % 3600) / 60)
  const s = elapsedSeconds.value % 60
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
})

onUnmounted(() => { if (timerInterval) clearInterval(timerInterval) })

const recruiterNoteContent = ref(localStorage.getItem('recruiter_personal_note_' + (route.params.interviewId || 'default')) || '')
watch(recruiterNoteContent, (val) => { localStorage.setItem('recruiter_personal_note_' + (route.params.interviewId || 'default'), val || ''); })
const saveRecruiterNote = () => {
  if (!recruiterNoteContent.value.trim()) return
  roomStore.createNote(recruiterNoteContent.value.trim(), ['manual'])
  recruiterNoteContent.value = ''
  activeToast.value = {
    type: 'success',
    message: 'Đã lưu ghi chú thành công!'
  }
}

const askedQuestionIds = ref(new Set())
const handleMarkQuestionAsked = (questionId) => {
  askedQuestionIds.value.add(questionId)
  roomStore.markQuestionAsked(questionId, elapsedSeconds.value * 1000)
}

// Helpers for UI
const candidateParticipant = computed(() => {
  return roomStore.participants.find(p => p.participant_type === 'candidate') || null
})

const formatTime = (ms) => {
  const date = new Date(ms)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

// Auto-hide control bar after 3.5s of mouse inactivity
const isControlBarVisible = ref(true)
let controlBarTimer = null

const resetControlBarTimer = () => {
  isControlBarVisible.value = true
  if (controlBarTimer) clearTimeout(controlBarTimer)
  controlBarTimer = setTimeout(() => {
    isControlBarVisible.value = false
  }, 3500)
}

const keepControlBarVisible = () => {
  isControlBarVisible.value = true
  if (controlBarTimer) clearTimeout(controlBarTimer)
}

onMounted(() => {
  resetControlBarTimer()
})
onUnmounted(() => {
  if (controlBarTimer) clearTimeout(controlBarTimer)
})

</script>

<template>
  <div style="height: 100vh; display: flex; flex-direction: column; background-color: var(--background)">
    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />
    <Toast v-if="activeToast" :type="activeToast.type" :message="activeToast.message" @close="activeToast = null" />
    <Toast v-if="isEnding" type="success" message="Đã kết thúc phỏng vấn. Đang lưu kết quả..." :duration="1500" />
    <!-- Cảnh báo quyền/thiết bị media: duration=0 để bám lại tới khi user xử lý -->
    <Toast v-if="liveKitError" type="error" :message="liveKitError" :duration="0" @close="clearError()" />
    
    <!-- Room Expired/Cancelled Banner -->
    <div v-if="roomStore.status === 'expired'" style="background-color: var(--warning); color: white; padding: 12px 24px; text-align: center; font-size: 14px; font-weight: 600;">
      <AlertTriangle size="16" style="vertical-align: text-bottom; margin-right: 4px;" />
      Phòng phỏng vấn đã hết hạn. Đang chuyển hướng...
    </div>
    <div v-if="roomStore.status === 'cancelled'" style="background-color: var(--danger); color: white; padding: 12px 24px; text-align: center; font-size: 14px; font-weight: 600;">
      <AlertTriangle size="16" style="vertical-align: text-bottom; margin-right: 4px;" />
      Buổi phỏng vấn đã bị hủy. {{ roomStore.cancelReason ? `Lý do: ${roomStore.cancelReason}` : '' }}
    </div>
    
    <!-- Paused Banner -->
    <div v-if="roomStore.status === 'paused'" style="background-color: var(--warning); color: #1a1a2e; padding: 8px 24px; text-align: center; font-size: 14px; display: flex; align-items: center; justify-content: center; gap: 12px;">
      <span>⏸️ Buổi phỏng vấn đang tạm dừng.</span>
      <Button variant="primary" style="height: 28px; font-size: 12px; padding: 0 12px;" @click="handleResumeCall">Tiếp tục</Button>
    </div>

    <!-- AI Error/Warning Banner -->
    <div v-if="aiStore.aiError && aiStore.aiError.severity === 'critical'" style="background-color: var(--danger); color: white; padding: 8px 24px; text-align: center; font-size: 14px;">
      <AlertTriangle size="16" style="vertical-align: text-bottom; margin-right: 4px;" />
      {{ aiStore.aiError.message || 'AI hiện không khả dụng.' }}
    </div>

    <!-- Header -->
    <div style="height: 48px; background-color: var(--surface); border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 16px">
      <div style="display: flex; align-items: center; gap: 16px">
        <div>
          <h1 class="text-h2" style="font-size: 16px; line-height: 1.2;">Phòng phỏng vấn <span v-if="candidateParticipant">- {{ candidateParticipant.display_name }}</span></h1>
          <div style="display: flex; gap: 12px; margin-top: 1px">
            <span v-if="roomStore.status === 'active'" class="text-helper" style="color: var(--danger); display: flex; align-items: center; gap: 4px">
              <div style="width: 8px; height: 8px; border-radius: 50%; background-color: var(--danger); animation: pulse 1.5s infinite;"></div> Đang ghi âm & Transcript
            </span>
            <span v-else-if="roomStore.status === 'waiting'" class="text-helper" style="color: var(--warning)">Đang chờ bắt đầu...</span>
            <span v-else-if="roomStore.status === 'paused'" class="text-helper" style="color: var(--warning)">⏸️ Tạm dừng</span>
            <span v-else-if="roomStore.status === 'completed'" class="text-helper" style="color: var(--success)">✅ Đã kết thúc</span>
            <span v-else class="text-helper">{{ roomStore.status }}</span>
            
            <!-- Live Timer -->
            <span v-if="roomStore.startedAt" class="text-helper" style="color: var(--primary); font-family: monospace; font-weight: 600;">{{ elapsedFormatted }}</span>
          </div>
        </div>
      </div>
      <div style="display: flex; align-items: center; gap: 12px">
        <Badge :type="roomStore.status === 'active' ? 'info' : 'secondary'">
          <Sparkles size="12" style="margin-right: 4px" /> {{ roomStore.status === 'active' ? 'AI Active' : 'AI Inactive' }}
        </Badge>
        <Button variant="secondary" style="padding: 6px 12px; font-size: 13px;" @click="router.push('/interviews')" title="Rời phòng tạm thời, buổi phỏng vấn vẫn được lưu">
          <LogOut size="16" style="margin-right: 6px" /> Rời tạm thời
        </Button>
      </div>
    </div>

    <!-- Main Content -->
    <div style="display: flex; flex: 1; overflow: hidden">
      
      <!-- Left: Video Area -->
      <div style="flex: 1; display: flex; flex-direction: column; padding: 6px; position: relative; overflow: hidden;">
        <!-- Video Grid -->
        <div style="flex: 1; display: flex; position: relative; width: 100%; height: 100%;">
          <!-- Candidate Video (Main Full Area) -->
          <div 
            @mousemove="resetControlBarTimer"
            @click="resetControlBarTimer"
            style="flex: 1; background-color: #0F172A; border-radius: var(--radius-lg); display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden; box-shadow: var(--shadow-lg); border: 1px solid rgba(255,255,255,0.08);"
          >
            <!-- Thẻ video thật cho LiveKit -->
            <video ref="remoteVideoEl" autoplay playsinline style="width: 100%; height: 100%; object-fit: cover; position: absolute; top: 0; left: 0;"></video>
            
            <div v-if="!candidateParticipant || candidateParticipant.connection_state === 'offline'" style="text-align: center; color: white; position: relative; z-index: 1;">
              <div style="width: 88px; height: 88px; border-radius: 50%; background: linear-gradient(135deg, rgba(255,255,255,0.15), rgba(255,255,255,0.05)); border: 1px solid rgba(255,255,255,0.2); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; font-size: 28px; font-weight: bold; box-shadow: 0 10px 25px -5px rgba(0,0,0,0.5);">N</div>
              <div style="font-size: 16px; font-weight: 500; color: #E2E8F0;">Đang chờ ứng viên tham gia...</div>
              <div style="font-size: 13px; color: #94A3B8; margin-top: 4px;">Video của ứng viên sẽ hiển thị ở trung tâm màn hình này</div>
            </div>
            
            <div v-if="candidateParticipant" style="position: absolute; top: 16px; left: 16px; background-color: rgba(15,23,42,0.75); backdrop-filter: blur(8px); color: white; padding: 6px 14px; border-radius: 9999px; font-size: 13px; font-weight: 500; z-index: 5; border: 1px solid rgba(255,255,255,0.15); display: flex; align-items: center; gap: 6px;">
              <div style="width: 8px; height: 8px; border-radius: 50%; background-color: #22C55E;"></div>
              {{ candidateParticipant.display_name }} (Ứng viên)
            </div>

            <!-- Floating Control Bar (Center Bottom of Video) -->
            <div 
              @mouseenter="keepControlBarVisible"
              @mouseleave="resetControlBarTimer"
              :style="{ 
                position: 'absolute', bottom: '24px', left: '50%', 
                transform: isControlBarVisible ? 'translateX(-50%) translateY(0)' : 'translateX(-50%) translateY(24px)', 
                opacity: isControlBarVisible ? 1 : 0, 
                pointerEvents: isControlBarVisible ? 'auto' : 'none', 
                transition: 'all 0.35s cubic-bezier(0.4, 0, 0.2, 1)', 
                backgroundColor: 'rgba(15, 23, 42, 0.85)', backdropFilter: 'blur(16px)', 
                border: '1px solid rgba(255, 255, 255, 0.15)', padding: '8px 16px', borderRadius: '9999px', 
                boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.6)', display: 'flex', alignItems: 'center', gap: '12px', zIndex: 20 
              }"
            >
              <button :title="isMicOn ? 'Tắt micro' : 'Bật micro'" @click="() => { toggleMic(); roomStore.updateMediaStatus(isMicOn, isCameraOn); }" :style="{ width: '44px', height: '44px', borderRadius: '50%', border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', transition: 'all 0.2s', backgroundColor: isMicOn ? 'rgba(51, 65, 85, 0.8)' : '#EF4444', color: 'white' }">
                <Mic v-if="isMicOn" size="20" />
                <MicOff v-else size="20" />
              </button>
              
              <button :title="isCameraOn ? 'Tắt camera' : 'Bật camera'" @click="() => { toggleCamera(); roomStore.updateMediaStatus(isMicOn, isCameraOn); }" :style="{ width: '44px', height: '44px', borderRadius: '50%', border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', transition: 'all 0.2s', backgroundColor: isCameraOn ? 'rgba(51, 65, 85, 0.8)' : '#EF4444', color: 'white' }">
                <Video v-if="isCameraOn" size="20" />
                <VideoOff v-else size="20" />
              </button>
              
              <button :title="isScreenSharing ? 'Dừng chia sẻ màn hình' : 'Chia sẻ màn hình'" @click="toggleScreenShare" :style="{ width: '44px', height: '44px', borderRadius: '50%', border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', transition: 'all 0.2s', backgroundColor: isScreenSharing ? '#2563EB' : 'rgba(51, 65, 85, 0.8)', color: 'white' }">
                <MonitorUp size="20" />
              </button>
              
              <button title="Mở khung Chat" @click="activeTab = 'chat'; isPanelExpanded = true" :style="{ width: '44px', height: '44px', borderRadius: '50%', border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', transition: 'all 0.2s', backgroundColor: activeTab === 'chat' && isPanelExpanded ? '#3B82F6' : 'rgba(51, 65, 85, 0.8)', color: 'white' }">
                <MessageSquare size="20" />
              </button>
              
              <div style="width: 1px; height: 28px; background-color: rgba(255, 255, 255, 0.2); margin: 0 4px;"></div>
              
              <!-- Nút Bắt đầu / Tạm dừng (Chỉ hiện nếu CHƯA Hoàn thành hoặc Hủy) -->
              <button v-if="roomStore.status !== 'completed' && roomStore.status !== 'cancelled' && roomStore.status !== 'expired'" :title="roomStore.status === 'active' ? 'Tạm dừng buổi phỏng vấn' : 'Bắt đầu phỏng vấn'" @click="handleToggleCallState" :style="{ backgroundColor: roomStore.status === 'active' ? 'rgba(51, 65, 85, 0.8)' : '#22C55E', color: 'white', border: 'none', height: '44px', padding: '0 20px', borderRadius: '9999px', fontWeight: 600, fontSize: '14px', display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', transition: 'all 0.2s', boxShadow: roomStore.status === 'active' ? 'none' : '0 4px 12px rgba(34, 197, 94, 0.4)' }">
                <Pause v-if="roomStore.status === 'active'" size="16" fill="white" />
                <Play v-else size="16" fill="white" />
                <span style="white-space: nowrap;">{{ roomStore.status === 'active' ? 'Tạm dừng' : 'Bắt đầu phỏng vấn' }}</span>
              </button>
              
              <!-- Nút Tiếp tục (Chỉ hiện khi đang Tạm dừng) -->
              <button v-if="roomStore.status === 'paused'" title="Tiếp tục phỏng vấn" @click="handleResumeCall" style="background-color: #3B82F6; color: white; border: none; height: 44px; padding: 0 20px; border-radius: 9999px; font-weight: 600; font-size: 14px; display: flex; align-items: center; gap: 8px; cursor: pointer; transition: background-color 0.2s; box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);">
                <Play size="16" fill="white" /> <span style="white-space: nowrap;">Tiếp tục</span>
              </button>
              
              <!-- Nút Kết thúc (Chỉ hiện khi đã Bắt đầu hoặc Đang tạm dừng) -->
              <button v-if="roomStore.status === 'active' || roomStore.status === 'paused'" title="Kết thúc cuộc gọi và xuất báo cáo" @click="handleEndCall" style="background-color: #EF4444; color: white; border: none; height: 44px; padding: 0 20px; border-radius: 9999px; font-weight: 600; font-size: 14px; display: flex; align-items: center; gap: 8px; cursor: pointer; transition: background-color 0.2s; box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);">
                <PhoneOff size="18" /> <span style="white-space: nowrap;">Kết thúc</span>
              </button>
            </div>
          </div>
          
          <!-- Recruiter Video (PiP top right) -->
          <div style="width: 220px; height: 140px; background-color: #1E293B; border-radius: var(--radius-md); display: flex; align-items: center; justify-content: center; position: absolute; top: 16px; right: 16px; box-shadow: 0 10px 25px -5px rgba(0,0,0,0.5); border: 2px solid rgba(255,255,255,0.15); overflow: hidden; z-index: 10; transition: all 0.3s ease;">
             <video ref="localVideoEl" autoplay playsinline muted style="width: 100%; height: 100%; object-fit: cover; position: absolute; top: 0; left: 0; transform: scaleX(-1);"></video>
             <div v-if="!isCameraOn" style="color: #94A3B8; font-size: 12px; font-weight: 500; position: relative; z-index: 1;">Camera của bạn đang tắt</div>
             <div style="position: absolute; bottom: 6px; left: 8px; background-color: rgba(0,0,0,0.6); color: white; padding: 2px 8px; border-radius: 4px; font-size: 11px; z-index: 2;">Bạn (Nhà tuyển dụng)</div>
          </div>
        </div>
      </div>

      <!-- Right: AI Panel -->
      <div v-if="aiStore.aiError?.severity !== 'critical'" class="bg-white border-l border-slate-200 flex flex-col shrink-0 shadow-sm relative z-20 transition-all duration-300" :style="{ width: isPanelExpanded ? '460px' : '48px' }">
        
        <!-- Collapsed Sidebar Icon Mode -->
        <div v-if="!isPanelExpanded" class="flex flex-col items-center py-3 gap-2.5 h-full bg-slate-50 border-l border-slate-200">
          <button @click="isPanelExpanded = true" class="w-9 h-9 rounded-xl bg-white hover:bg-slate-100 text-slate-600 border border-slate-200 flex items-center justify-center transition-colors shadow-sm" title="Mở rộng bảng công cụ">
            <ChevronsLeft class="w-4 h-4" />
          </button>
          <div class="w-6 h-px bg-slate-200 my-0.5"></div>
          
          <button @click="() => { activeTab = 'assistant'; isPanelExpanded = true; }" class="w-9 h-9 rounded-xl flex items-center justify-center transition-colors relative" :class="activeTab === 'assistant' ? 'bg-blue-50 text-blue-600 border border-blue-200 shadow-sm font-bold' : 'text-slate-500 hover:bg-slate-100 hover:text-slate-800'" title="AI Assistant">
            <Sparkles class="w-4 h-4" />
          </button>
          <button @click="() => { activeTab = 'rubric'; isPanelExpanded = true; }" class="w-9 h-9 rounded-xl flex items-center justify-center transition-colors relative" :class="activeTab === 'rubric' ? 'bg-blue-50 text-blue-600 border border-blue-200 shadow-sm font-bold' : 'text-slate-500 hover:bg-slate-100 hover:text-slate-800'" title="Tiêu chí">
            <CheckCircle class="w-4 h-4" />
          </button>
          <button @click="() => { activeTab = 'transcript'; isPanelExpanded = true; }" class="w-9 h-9 rounded-xl flex items-center justify-center transition-colors relative" :class="activeTab === 'transcript' ? 'bg-blue-50 text-blue-600 border border-blue-200 shadow-sm font-bold' : 'text-slate-500 hover:bg-slate-100 hover:text-slate-800'" title="Transcript">
            <MessageSquare class="w-4 h-4" />
          </button>
          <button @click="() => { activeTab = 'chat'; isPanelExpanded = true; }" class="w-9 h-9 rounded-xl flex items-center justify-center transition-colors relative" :class="activeTab === 'chat' ? 'bg-blue-50 text-blue-600 border border-blue-200 shadow-sm font-bold' : 'text-slate-500 hover:bg-slate-100 hover:text-slate-800'" title="Chat">
            <MessageSquare class="w-4 h-4" />
            <span v-if="chatStore.messages.length > 0" class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-rose-500 border border-white"></span>
          </button>
          <button @click="() => { activeTab = 'notes'; isPanelExpanded = true; }" class="w-9 h-9 rounded-xl flex items-center justify-center transition-colors relative" :class="activeTab === 'notes' ? 'bg-blue-50 text-blue-600 border border-blue-200 shadow-sm font-bold' : 'text-slate-500 hover:bg-slate-100 hover:text-slate-800'" title="Ghi chú">
            <FileText class="w-4 h-4" />
          </button>
        </div>

        <!-- Expanded Sidebar Mode -->
        <template v-else>
          <!-- Tabs Header -->
          <div class="flex border-b border-slate-200 shrink-0 bg-slate-50 items-center justify-between pr-2">
            <div class="flex flex-1 overflow-x-auto custom-scrollbar">
              <button @click="activeTab = 'assistant'" class="py-3 px-2 flex items-center justify-center gap-1.5 font-medium text-xs transition-colors border-b-2 whitespace-nowrap" :class="activeTab === 'assistant' ? 'text-blue-600 border-blue-600 bg-blue-50/60 font-semibold' : 'text-slate-500 border-transparent hover:bg-slate-100 hover:text-slate-800'" title="AI Assistant">
                <Sparkles class="w-4 h-4 shrink-0" /> AI Assistant
              </button>
              <button @click="activeTab = 'rubric'" class="py-3 px-2 flex items-center justify-center gap-1.5 font-medium text-xs transition-colors border-b-2 whitespace-nowrap" :class="activeTab === 'rubric' ? 'text-blue-600 border-blue-600 bg-blue-50/60 font-semibold' : 'text-slate-500 border-transparent hover:bg-slate-100 hover:text-slate-800'" title="Tiêu chí">
                <CheckCircle class="w-4 h-4 shrink-0" /> Tiêu chí
              </button>
              <button @click="activeTab = 'transcript'" class="py-3 px-2 flex items-center justify-center gap-1.5 font-medium text-xs transition-colors border-b-2 whitespace-nowrap" :class="activeTab === 'transcript' ? 'text-blue-600 border-blue-600 bg-blue-50/60 font-semibold' : 'text-slate-500 border-transparent hover:bg-slate-100 hover:text-slate-800'" title="Transcript">
                <MessageSquare class="w-4 h-4 shrink-0" /> Transcript
              </button>
              <button @click="activeTab = 'chat'" class="py-3 px-2 flex items-center justify-center gap-1.5 font-medium text-xs transition-colors border-b-2 relative whitespace-nowrap" :class="activeTab === 'chat' ? 'text-blue-600 border-blue-600 bg-blue-50/60 font-semibold' : 'text-slate-500 border-transparent hover:bg-slate-100 hover:text-slate-800'" title="Chat">
                <MessageSquare class="w-4 h-4 shrink-0" /> Chat
                <span v-if="chatStore.messages.length > 0" class="absolute top-2 right-1.5 w-2 h-2 rounded-full bg-rose-500 border border-white"></span>
              </button>
              <button @click="activeTab = 'notes'" class="py-3 px-2 flex items-center justify-center gap-1.5 font-medium text-xs transition-colors border-b-2 whitespace-nowrap" :class="activeTab === 'notes' ? 'text-blue-600 border-blue-600 bg-blue-50/60 font-semibold' : 'text-slate-500 border-transparent hover:bg-slate-100 hover:text-slate-800'" title="Ghi chú">
                <FileText class="w-4 h-4 shrink-0" /> Ghi chú
              </button>
            </div>
            <button @click="isPanelExpanded = false" class="w-8 h-8 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-200/60 flex items-center justify-center transition-colors ml-1 shrink-0" title="Thu gọn bảng công cụ">
              <ChevronsRight class="w-4 h-4" />
            </button>
          </div>

          <!-- Tab Content -->
          <div class="flex-1 overflow-y-auto overflow-x-hidden p-5 bg-white flex flex-col relative custom-scrollbar">
          
          <!-- AI Error Degraded Indicator -->
          <div v-if="aiStore.aiError && aiStore.aiError.severity === 'degraded'" style="margin-bottom: 16px; font-size: 12px; color: var(--warning);">
            <AlertTriangle size="12" style="vertical-align: middle;"/> {{ aiStore.aiError.message }}
          </div>

          <div v-if="activeTab === 'assistant'" style="display: flex; flex-direction: column; gap: 24px">
            
            <div v-if="aiStore.isThinking" style="text-align: center; padding: 20px; color: var(--text-muted)">
              <div class="spinner" style="margin: 0 auto 10px;"></div>
              {{ aiStore.aiThinkingMessage }}
            </div>

            <div v-else>
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px">
                <h3 class="text-body" style="font-weight: 600;">Gợi ý câu hỏi tiếp theo</h3>
                <Button variant="secondary" style="height: 28px; font-size: 12px; padding: 0 10px;" @click="aiStore.requestSuggestion('general', null, roomStore.roomId, roomStore.interviewId)">Tạo mới</Button>
              </div>
              <div v-if="aiStore.suggestions.length === 0" class="text-helper text-center text-muted" style="margin-top: 40px">
                Chưa có gợi ý nào từ AI.
              </div>
              <div style="display: flex; flex-direction: column; gap: 12px">
                <div v-for="sugg in aiStore.suggestions" :key="sugg.suggestion_id" :style="{ padding: '12px', border: '1px solid var(--border)', borderRadius: 'var(--radius)', backgroundColor: 'var(--surface)', opacity: askedQuestionIds.has(sugg.suggestion_id) ? 0.6 : 1 }">
                  <div style="font-size: 11px; color: var(--text-muted); margin-bottom: 4px; display: flex; justify-content: space-between;">
                    <span>Target: {{ sugg.target_skill || 'Chung' }}</span>
                    <span>Độ tin cậy: {{ Math.round((sugg.confidence||0)*100) }}%</span>
                  </div>
                  <p class="text-body" style="margin-bottom: 12px">{{ sugg.content }}</p>
                  <div v-if="sugg.reason" style="font-size: 12px; background: rgba(8,145,178,0.1); padding: 8px; border-radius: 4px; margin-bottom: 12px; color: var(--accent);">
                    Lý do: {{ sugg.reason }}
                  </div>
                  <div style="display: flex; justify-content: flex-end;">
                    <Button 
                      v-if="!askedQuestionIds.has(sugg.suggestion_id)" 
                      variant="secondary" 
                      style="height: 24px; font-size: 11px; padding: 0 8px;" 
                      @click="handleMarkQuestionAsked(sugg.suggestion_id)"
                    >
                      Đánh dấu đã hỏi
                    </Button>
                    <span v-else style="font-size: 11px; color: var(--success); font-weight: 600;">✓ Đã hỏi</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'rubric'" style="display: flex; flex-direction: column; gap: 16px">
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <h3 class="text-body" style="font-weight: 600;">Đánh giá theo rubric</h3>
              <Button variant="secondary" style="height: 28px; font-size: 12px; padding: 0 10px;" @click="requestAIScore">Chấm điểm bằng AI</Button>
            </div>
            <div v-if="aiStore.scores.length === 0" class="text-helper text-center text-muted">
              Chưa có điểm đánh giá nào.
            </div>
            <div v-for="score in aiStore.scores" :key="score.criterion_name">
              <div style="display: flex; justify-content: space-between; margin-bottom: 8px">
                <span class="text-body" style="font-weight: 500">{{ score.criterion_name }}</span>
                <span class="text-body" :style="{ color: score.status === 'scored' ? 'var(--success)' : 'var(--warning)' }">
                  {{ score.status === 'scored' ? score.score : '--' }} / {{ score.max_score }}
                </span>
              </div>
              <div style="height: 6px; background-color: var(--surface-soft); border-radius: 3px">
                <div v-if="score.status === 'scored'" style="height: 100%; border-radius: 3px; background-color: var(--success);" :style="{ width: (score.score / score.max_score * 100) + '%' }"></div>
              </div>
              <p v-if="score.evidence" class="text-helper" style="margin-top: 4px; color: var(--text-muted)">Bằng chứng: {{ score.evidence }}</p>
              <p v-if="score.status === 'insufficient_evidence'" class="text-helper" style="margin-top: 4px; color: var(--warning)">Chưa đủ bằng chứng để chấm điểm.</p>
            </div>
          </div>

          <div v-if="activeTab === 'transcript'" style="display: flex; flex-direction: column; gap: 16px">
             <div v-if="transcriptStore.transcripts.length === 0 && Object.keys(transcriptStore.activePartials).length === 0" class="text-helper text-center text-muted">
               Chưa có hội thoại nào được ghi nhận.
             </div>
             
             <!-- Final Transcripts -->
             <div v-for="item in transcriptStore.transcripts" :key="item.transcript_id" :style="{ padding: '12px', backgroundColor: 'var(--surface-soft)', borderRadius: 'var(--radius)', borderLeft: item.confidence < 0.5 ? '4px solid var(--warning)' : 'none' }">
              <div style="display: flex; justify-content: space-between; margin-bottom: 4px">
                <span class="text-helper" style="font-weight: 600; color: var(--primary)">
                  {{ item.speaker_name || item.speaker_type }}
                  <span v-if="item.confidence < 0.5" style="color: var(--warning); font-size: 11px; margin-left: 8px;" title="Độ tin cậy nhận diện thấp">(⚠️ Nhận diện thấp)</span>
                </span>
                <span class="text-helper" style="font-size: 11px">{{ formatTime(item.start_time_ms) }}</span>
              </div>
              <p class="text-body" style="color: var(--text-main);">
                {{ item.content }}
              </p>
            </div>
            
            <!-- Partial Transcripts -->
            <div v-for="item in transcriptStore.activePartials" :key="item.transcript_id" style="padding: 12px; background-color: rgba(241, 245, 249, 0.5); border-radius: var(--radius); opacity: 0.7;">
              <div style="display: flex; justify-content: space-between; margin-bottom: 4px">
                <span class="text-helper" style="font-weight: 600; color: var(--primary)">
                  {{ item.speaker_name || item.speaker_type }} (đang nói...)
                </span>
              </div>
              <p class="text-body" style="color: var(--text-main);">
                {{ item.content }}
              </p>
            </div>
          </div>
          
          <div v-if="activeTab === 'chat'" style="display: flex; flex-direction: column; height: 100%;">
            <div style="flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 12px; margin-bottom: 16px; padding-right: 4px;">
              <div v-if="chatStore.messages.length === 0" class="text-helper text-center text-muted" style="margin-top: 20px;">
                Chưa có tin nhắn nào.
              </div>
              <div v-for="msg in chatStore.messages" :key="msg.message_id" style="display: flex; flex-direction: column;" :style="{ alignItems: (msg.sender_type === 'recruiter' || msg.sender_type === 'local') ? 'flex-end' : 'flex-start' }">
                <span style="font-size: 11px; color: var(--text-muted); margin-bottom: 2px;">{{ (msg.sender_type === 'recruiter' || msg.sender_type === 'local') ? 'Bạn (Nhà tuyển dụng)' : (msg.sender_name || 'Ứng viên') }}</span>
                <div style="padding: 10px 14px; border-radius: 12px; max-width: 85%; font-size: 13px;" :style="{ backgroundColor: (msg.sender_type === 'recruiter' || msg.sender_type === 'local') ? 'var(--primary)' : 'var(--surface-soft)', color: (msg.sender_type === 'recruiter' || msg.sender_type === 'local') ? '#ffffff' : 'var(--text-main)' }">
                  {{ msg.message }}
                </div>
              </div>
            </div>
            <div style="display: flex; align-items: center; gap: 8px; background-color: var(--surface); border: 1px solid var(--border); border-radius: 24px; padding: 4px 6px 4px 16px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); transition: border-color 0.2s; flex-shrink: 0;">
               <input type="text" maxlength="2000" v-model="chatInput" placeholder="Nhập tin nhắn..." style="flex: 1; border: none; background: transparent; outline: none; font-size: 13.5px; color: var(--text-main); padding: 6px 0;" @keyup.enter="handleSendMessage">
               <button @click="handleSendMessage" style="width: 34px; height: 34px; border-radius: 50%; background-color: var(--primary); color: white; border: none; display: flex; align-items: center; justify-content: center; cursor: pointer; transition: transform 0.15s, opacity 0.2s; flex-shrink: 0; box-shadow: 0 2px 6px rgba(37, 99, 235, 0.3);" title="Gửi tin nhắn">
                 <Send size="16" />
               </button>
            </div>
          </div>

          <div v-if="activeTab === 'notes'" style="display: flex; flex-direction: column; gap: 16px; height: 100%;">
            <div style="flex: 1; display: flex; flex-direction: column; gap: 8px;">
              <h3 class="text-body" style="font-weight: 600;">Ghi chú phỏng vấn</h3>
              <textarea
                v-model="recruiterNoteContent"
                class="app-input"
                style="flex: 1; min-height: 150px; resize: none; padding: 12px; font-family: inherit; font-size: 14px; border-radius: var(--radius); border: 1px solid var(--border); background-color: var(--surface-soft);"
                placeholder="Nhập ghi chú nội bộ của nhà tuyển dụng tại đây (Candidate không nhìn thấy)..."
              ></textarea>
            </div>
            <div style="display: flex; justify-content: flex-end;">
              <Button variant="primary" @click="saveRecruiterNote">Lưu ghi chú</Button>
            </div>
          </div>

          </div>
        </template>
      </div>
    </div>

    <Modal :isOpen="showEndModal" @close="showEndModal = false" title="Kết thúc phỏng vấn">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn kết thúc buổi phỏng vấn này?</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showEndModal = false">Huỷ</Button>
        <Button variant="primary" @click="confirmEndCall" style="background-color: var(--danger); color: white; border-color: var(--danger)">OK</Button>
      </div>
    </Modal>

    <Modal :isOpen="showLeaveWarningModal" @close="handleCancelLeave" title="Xác nhận rời phòng phỏng vấn">
      <p class="text-body" style="margin-bottom: 24px">Bạn đang trong phòng phỏng vấn. Bạn có chắc chắn muốn rời khỏi phòng ngay bây giờ? (Bạn vẫn có thể quay lại sau từ danh sách phỏng vấn)</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="handleCancelLeave">Ở lại phòng</Button>
        <Button variant="primary" @click="handleConfirmLeave" style="background-color: var(--warning); color: #1a1a2e; font-weight: 600;">Rời đi</Button>
      </div>
    </Modal>
  </div>
</template>
<style scoped>
.tabs-container::-webkit-scrollbar {
  display: none;
}
.tabs-container {
  scrollbar-width: none;
  -ms-overflow-style: none;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
  font-weight: 500;
  border-bottom: 2px solid transparent;
  white-space: nowrap;
  font-size: 13px;
}
.nav-item:hover {
  color: var(--text-main);
}
.nav-item.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
}
.spinner {
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-left-color: var(--primary);
  border-radius: 50%;
  width: 24px;
  height: 24px;
  animation: spin 1s linear infinite;
}
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
</style>
