<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Button from '../../components/common/AppButton.vue'
import Card from '../../components/common/AppCard.vue'
import Badge from '../../components/common/AppBadge.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mic, MicOff, Video, VideoOff, MonitorUp, MessageSquare, PhoneOff, Sparkles, CheckCircle, AlertTriangle, FileText, ChevronRight, ChevronLeft } from 'lucide-vue-next'
import { useLiveKit } from '../../composables/useLiveKit'
import { roomService } from '../../services/room.service'
import { authStore } from '../../stores/auth.store'

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

const {
  isConnected: isLiveKitConnected, error: liveKitError, isMicOn, isCameraOn,
  localVideoEl, remoteVideoEl,
  connectToRoom, toggleMic, toggleCamera, disconnect: liveKitDisconnect
} = useLiveKit()

const activeTab = ref('assistant')
const isPanelExpanded = ref(true)

const showEndModal = ref(false)
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

  // Khởi tạo Listeners cho các Stores
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
          const livekitUrl = import.meta.env.VITE_LIVEKIT_URL || 'ws://localhost:7880'
          await connectToRoom(livekitUrl, token)
        }
      }
    } catch (err) {
      console.error('Không thể lấy room token', err)
    }
  }
})

onUnmounted(() => {
  liveKitDisconnect()
  roomStore.disconnectRoom()
  chatStore.cleanupListeners()
  transcriptStore.cleanupListeners()
  aiStore.cleanupListeners()
})

const handleEndCall = () => {
  showEndModal.value = true
}

const confirmEndCall = () => {
  showEndModal.value = false
  isEnding.value = true
  roomStore.endInterview() // Gửi event end
  
  setTimeout(() => {
    router.push({ path: '/reports', state: { message: 'Đã lưu kết quả phỏng vấn thành công' } })
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

const recruiterNoteContent = ref('')
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

</script>

<template>
  <div style="height: 100vh; display: flex; flex-direction: column; background-color: var(--background)">
    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />
    <Toast v-if="activeToast" :type="activeToast.type" :message="activeToast.message" @close="activeToast = null" />
    <Toast v-if="isEnding" type="success" message="Đã kết thúc phỏng vấn. Đang lưu kết quả..." :duration="1500" />
    
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
    <div style="height: 64px; background-color: var(--surface); border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 24px">
      <div style="display: flex; align-items: center; gap: 16px">
        <div>
          <h1 class="text-h2">Phòng phỏng vấn <span v-if="candidateParticipant">- {{ candidateParticipant.display_name }}</span></h1>
          <div style="display: flex; gap: 12px; margin-top: 4px">
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
        <Button v-if="roomStore.status === 'waiting'" variant="primary" @click="handleStartCall">Bắt đầu Phỏng vấn</Button>
        <Button v-if="roomStore.status === 'active'" variant="secondary" @click="handlePauseCall">Tạm dừng</Button>
        <Button variant="secondary" @click="router.push('/interviews')">Rời tạm thời</Button>
      </div>
    </div>

    <!-- Main Content -->
    <div style="display: flex; flex: 1; overflow: hidden">
      
      <!-- Left: Video Area -->
      <div style="flex: 1; display: flex; flex-direction: column; padding: 12px; position: relative; overflow: hidden;">
        <!-- Video Grid -->
        <div style="flex: 1; display: flex; position: relative; width: 100%; height: 100%;">
          <!-- Candidate Video (Main Full Area) -->
          <div style="flex: 1; background-color: #0F172A; border-radius: var(--radius-lg); display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden; box-shadow: var(--shadow-lg); border: 1px solid rgba(255,255,255,0.08);">
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
            <div style="position: absolute; bottom: 24px; left: 50%; transform: translateX(-50%); background-color: rgba(15, 23, 42, 0.85); backdrop-filter: blur(16px); border: 1px solid rgba(255, 255, 255, 0.15); padding: 8px 16px; border-radius: 9999px; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.6); display: flex; align-items: center; gap: 12px; z-index: 20;">
              <button :title="isMicOn ? 'Tắt micro' : 'Bật micro'" @click="() => { toggleMic(); roomStore.updateMediaStatus(isMicOn, isCameraOn); }" :style="{ width: '44px', height: '44px', borderRadius: '50%', border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', transition: 'all 0.2s', backgroundColor: isMicOn ? 'rgba(51, 65, 85, 0.8)' : '#EF4444', color: 'white' }">
                <Mic v-if="isMicOn" size="20" />
                <MicOff v-else size="20" />
              </button>
              
              <button :title="isCameraOn ? 'Tắt camera' : 'Bật camera'" @click="() => { toggleCamera(); roomStore.updateMediaStatus(isMicOn, isCameraOn); }" :style="{ width: '44px', height: '44px', borderRadius: '50%', border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', transition: 'all 0.2s', backgroundColor: isCameraOn ? 'rgba(51, 65, 85, 0.8)' : '#EF4444', color: 'white' }">
                <Video v-if="isCameraOn" size="20" />
                <VideoOff v-else size="20" />
              </button>
              
              <button title="Chia sẻ màn hình" :style="{ width: '44px', height: '44px', borderRadius: '50%', border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', transition: 'all 0.2s', backgroundColor: 'rgba(51, 65, 85, 0.8)', color: 'white' }">
                <MonitorUp size="20" />
              </button>
              
              <button title="Mở khung Chat" @click="activeTab = 'chat'; isPanelExpanded = true" :style="{ width: '44px', height: '44px', borderRadius: '50%', border: 'none', display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', transition: 'all 0.2s', backgroundColor: activeTab === 'chat' && isPanelExpanded ? '#3B82F6' : 'rgba(51, 65, 85, 0.8)', color: 'white' }">
                <MessageSquare size="20" />
              </button>
              
              <div style="width: 1px; height: 28px; background-color: rgba(255, 255, 255, 0.2); margin: 0 4px;"></div>
              
              <button title="Kết thúc cuộc gọi" @click="handleEndCall" style="background-color: #EF4444; color: white; border: none; height: 44px; padding: 0 20px; border-radius: 9999px; font-weight: 600; font-size: 14px; display: flex; align-items: center; gap: 8px; cursor: pointer; transition: background-color 0.2s; box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);">
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
      <div v-if="aiStore.aiError?.severity !== 'critical'" :style="{ width: isPanelExpanded ? '420px' : '64px', minWidth: isPanelExpanded ? '420px' : '64px', maxWidth: isPanelExpanded ? '420px' : '64px', flexShrink: 0, backgroundColor: 'var(--surface)', borderLeft: '1px solid var(--border)', display: 'flex', flexDirection: 'column', transition: 'width 0.3s ease, min-width 0.3s ease, max-width 0.3s ease' }">
        
        <!-- Toggle Button & Header -->
        <div :style="{ display: 'flex', alignItems: 'center', justifyContent: isPanelExpanded ? 'space-between' : 'center', padding: '12px', borderBottom: '1px solid var(--border)', flexShrink: 0 }">
          <span v-if="isPanelExpanded" style="font-weight: 600; font-size: 14px; color: var(--text-main)">Công cụ Hỗ trợ</span>
          <Button variant="ghost" style="padding: 4px; height: auto;" @click="isPanelExpanded = !isPanelExpanded">
            <ChevronRight v-if="isPanelExpanded" size="20" />
            <ChevronLeft v-else size="20" />
          </Button>
        </div>

        <!-- Tabs -->
        <div class="tabs-container" :style="{ display: 'flex', flexDirection: isPanelExpanded ? 'row' : 'column', borderBottom: isPanelExpanded ? '1px solid var(--border)' : 'none', overflowX: isPanelExpanded ? 'auto' : 'hidden', overflowY: 'hidden', flexShrink: 0 }">
          <div :class="['nav-item', activeTab === 'assistant' ? 'active' : '']" :style="{ flex: isPanelExpanded ? '0 0 auto' : 'unset', justifyContent: 'center', cursor: 'pointer', padding: isPanelExpanded ? '16px 12px' : '16px 0' }" @click="activeTab = 'assistant'; isPanelExpanded = true" title="AI Assistant">
            <Sparkles size="16" /> <span v-if="isPanelExpanded">AI Assistant</span>
          </div>
          <div :class="['nav-item', activeTab === 'rubric' ? 'active' : '']" :style="{ flex: isPanelExpanded ? '0 0 auto' : 'unset', justifyContent: 'center', cursor: 'pointer', padding: isPanelExpanded ? '16px 12px' : '16px 0' }" @click="activeTab = 'rubric'; isPanelExpanded = true" title="Tiêu chí">
            <CheckCircle size="16" /> <span v-if="isPanelExpanded">Tiêu chí</span>
          </div>
          <div :class="['nav-item', activeTab === 'transcript' ? 'active' : '']" :style="{ flex: isPanelExpanded ? '0 0 auto' : 'unset', justifyContent: 'center', cursor: 'pointer', padding: isPanelExpanded ? '16px 12px' : '16px 0' }" @click="activeTab = 'transcript'; isPanelExpanded = true" title="Transcript">
            <MessageSquare size="16" /> <span v-if="isPanelExpanded">Transcript</span>
          </div>
          <div :class="['nav-item', activeTab === 'chat' ? 'active' : '']" :style="{ flex: isPanelExpanded ? '0 0 auto' : 'unset', justifyContent: 'center', cursor: 'pointer', padding: isPanelExpanded ? '16px 12px' : '16px 0' }" @click="activeTab = 'chat'; isPanelExpanded = true" title="Chat">
             <MessageSquare size="16" /> <span v-if="isPanelExpanded">Chat</span>
          </div>
          <div :class="['nav-item', activeTab === 'notes' ? 'active' : '']" :style="{ flex: isPanelExpanded ? '0 0 auto' : 'unset', justifyContent: 'center', cursor: 'pointer', padding: isPanelExpanded ? '16px 12px' : '16px 0' }" @click="activeTab = 'notes'; isPanelExpanded = true" title="Ghi chú">
            <FileText size="16" /> <span v-if="isPanelExpanded">Ghi chú</span>
          </div>
        </div>

        <!-- Tab Content -->
        <div v-show="isPanelExpanded" style="flex: 1; overflow-y: auto; overflow-x: hidden; padding: 24px;">
          
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
            <div style="flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 12px; margin-bottom: 16px;">
              <div v-if="chatStore.messages.length === 0" class="text-helper text-center text-muted">
                Không có tin nhắn.
              </div>
              <div v-for="msg in chatStore.messages" :key="msg.message_id" style="padding: 8px 12px; border-radius: var(--radius); background-color: var(--surface-soft);">
                <div style="font-size: 11px; font-weight: bold; color: var(--primary); margin-bottom: 2px;">{{ msg.sender_name }} <span v-if="msg.visibility === 'recruiter_only'" style="color:var(--danger)">(Internal)</span></div>
                <div class="text-body">{{ msg.message }}</div>
              </div>
            </div>
            <div style="display: flex; gap: 8px;">
               <input type="text" id="chatInput" placeholder="Nhập tin nhắn..." class="app-input" style="flex: 1;" @keyup.enter="(e) => { chatStore.sendChat(e.target.value, 'room', roomStore.roomId, roomStore.interviewId); e.target.value=''; }">
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
      </div>
    </div>

    <Modal :isOpen="showEndModal" @close="showEndModal = false" title="Kết thúc phỏng vấn">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn kết thúc buổi phỏng vấn này?</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showEndModal = false">Huỷ</Button>
        <Button variant="primary" @click="confirmEndCall" style="background-color: var(--danger); color: white; border-color: var(--danger)">OK</Button>
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
