<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter, useRoute, onBeforeRouteLeave } from 'vue-router'
import Button from '../../components/common/AppButton.vue'
import Modal from '../../components/common/AppModal.vue'
// Toast trước đây bị dùng trong template mà KHÔNG import — cả entryToast lẫn
// isLeaving toast chưa từng render được.
import Toast from '../../components/common/AppToast.vue'
import { Mic, MicOff, Video, VideoOff, MonitorUp, MessageSquare, PhoneOff, CheckCircle, FileText, ChevronsRight, ChevronsLeft, X, Send } from 'lucide-vue-next'
import { useLiveKit } from '../../composables/useLiveKit'
import { useSpeechToText } from '../../composables/useSpeechToText'

// Import Stores
import { useRoomStore } from '../../stores/room.store'
import { useChatStore } from '../../stores/chat.store'

const router = useRouter()
const route = useRoute()

// Stores
const roomStore = useRoomStore()
const chatStore = useChatStore()

const {
  isConnected: isLiveKitConnected, error: liveKitError, isMicOn, isCameraOn, isScreenSharing,
  localVideoEl, remoteVideoEl,
  connectToRoom, toggleMic, toggleCamera, toggleScreenShare, disconnect: liveKitDisconnect,
  clearError
} = useLiveKit()

// Live transcription of THIS candidate's mic (free Web Speech API). Pushes
// transcript:partial/final up the realtime WS; the gateway attributes the
// speaker, broadcasts, persists, and triggers AI suggestion/scoring.
const { supported: sttSupported, listening: sttListening, start: startSTT, stop: stopSTT } = useSpeechToText()

const syncSTT = () => {
  const shouldRun = roomStore.status === 'active' && isMicOn.value && sttSupported.value
  if (shouldRun && !sttListening.value) {
    startSTT({
      roomId: roomStore.roomId,
      interviewId: roomStore.interviewId,
      speakerType: 'candidate',
      speakerName: 'Ứng viên',
    })
  } else if (!shouldRun && sttListening.value) {
    stopSTT()
  }
}

watch(() => [roomStore.status, isMicOn.value], syncSTT)

const activeTab = ref(sessionStorage.getItem('candidate_active_tab') || 'info')
watch(activeTab, (val) => {
  sessionStorage.setItem('candidate_active_tab', val);
})

const showLeaveModal = ref(false)
const isLeaving = ref(false)

const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)
const interviewInfo = ref(history.state?.interviewInfo || null)
// Parse JWT to extract room_id and interview_id
const parseJwt = (token) => {
  try {
    return JSON.parse(atob(token.split('.')[1]));
  } catch (e) {
    return {};
  }
}

const token = route.query.token || ''
const tokenClaims = parseJwt(token)
const roomId = tokenClaims.room_id || (tokenClaims.video && tokenClaims.video.room) || ''
const interviewId = tokenClaims.interview_id || ''

const chatInput = ref('')
const chatError = ref('')
const chatCooldown = ref(false)

const handleSendMessage = () => {
  const message = chatInput.value.trim()
  if (!message || chatCooldown.value) return
  if (!roomId || !interviewId || !roomStore.isConnected) {
    chatError.value = 'Không thể gửi tin nhắn khi phòng chưa kết nối.'
    return
  }
  if (message.length > 2000) {
    chatError.value = 'Tin nhắn không được vượt quá 2.000 ký tự.'
    return
  }
  chatError.value = ''
  chatStore.sendChat(message, 'room', roomId, interviewId, 'Ứng viên (Bạn)')
  chatInput.value = ''
  chatCooldown.value = true
  setTimeout(() => { chatCooldown.value = false }, 500)
}

onMounted(async () => {
  if (history.state?.message) {
    window.history.replaceState({ interviewInfo: history.state.interviewInfo }, document.title)
  }
  
  // Khởi tạo Listeners cho Store
  chatStore.setupListeners()

  // Logic kết nối LiveKit & WebSocket
  if (token) {
    try {
      // Connect WebSocket Realtime
      roomStore.connectRoom(token, roomId, interviewId)

      // Gửi room:join (candidate) sau khi WebSocket mở
      const tryJoinAsCandidate = (retries = 10) => {
        if (roomStore.isConnected) {
          roomStore.sendRoomJoin('candidate')
        } else if (retries > 0) {
          setTimeout(() => tryJoinAsCandidate(retries - 1), 300)
        }
      }
      setTimeout(() => tryJoinAsCandidate(), 500)

      // Connect LiveKit (hoặc native getUserMedia nếu không có server)
      const livekitUrl = import.meta.env.VITE_LIVEKIT_URL || `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}`
      await connectToRoom(livekitUrl, token)
    } catch (err) {
      console.error('Không thể vào phòng', err)
    }
  }
  window.addEventListener('beforeunload', handleCandidateBeforeUnload)
})

const showNavWarningModal = ref(false)
const pendingCandidateRouteResolve = ref(null)

const handleCandidateBeforeUnload = (e) => {
  if (!isLeaving.value && roomStore.status !== 'completed' && roomStore.status !== 'expired' && roomStore.status !== 'cancelled') {
    e.preventDefault()
    e.returnValue = ''
  }
}

onBeforeRouteLeave((to, from) => {
  if (isLeaving.value || roomStore.status === 'completed' || roomStore.status === 'expired' || roomStore.status === 'cancelled') {
    return true
  }
  return new Promise((resolve) => {
    pendingCandidateRouteResolve.value = resolve
    showNavWarningModal.value = true
  })
})

const handleCandidateConfirmLeave = () => {
  showNavWarningModal.value = false
  if (pendingCandidateRouteResolve.value) {
    pendingCandidateRouteResolve.value(true)
    pendingCandidateRouteResolve.value = null
  }
}

const handleCandidateCancelLeave = () => {
  showNavWarningModal.value = false
  if (pendingCandidateRouteResolve.value) {
    pendingCandidateRouteResolve.value(false)
    pendingCandidateRouteResolve.value = null
  }
}

onUnmounted(() => {
  window.removeEventListener('beforeunload', handleCandidateBeforeUnload)
  stopSTT()
  liveKitDisconnect()
  roomStore.disconnectRoom()
  chatStore.cleanupListeners()
})

const handleLeave = () => {
  showLeaveModal.value = true
}

const confirmLeave = () => {
  showLeaveModal.value = false
  isLeaving.value = true
  // Candidate không có quyền gọi interview:end, chỉ disconnect WebSocket
  roomStore.disconnectRoom()
  liveKitDisconnect()
  setTimeout(() => {
    router.push({ path: '/home', state: { message: 'Rời phòng phỏng vấn thành công' } })
  }, 1500)
}

// [FIX] Watch status expired/cancelled/completed → redirect candidate
watch(() => roomStore.status, (newStatus) => {
  if (newStatus === 'expired') {
    setTimeout(() => {
      router.push({ path: '/home', state: { message: 'Phòng phỏng vấn đã hết hạn.' } })
    }, 2000)
  }
  if (newStatus === 'cancelled') {
    setTimeout(() => {
      router.push({ path: '/home', state: { message: 'Buổi phỏng vấn đã bị hủy.' } })
    }, 2000)
  }
  if (newStatus === 'completed') {
    setTimeout(() => {
      router.push({ path: '/home', state: { message: 'Buổi phỏng vấn đã kết thúc. Cảm ơn bạn!' } })
    }, 3000)
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

// Helpers for UI
const recruiterParticipant = computed(() => {
  return roomStore.participants.find(p => p.participant_type === 'recruiter') || null
})

const isPanelExpanded = ref(sessionStorage.getItem('candidate_panel_expanded') !== 'false')
watch(isPanelExpanded, (val) => {
  sessionStorage.setItem('candidate_panel_expanded', val);
})

const candidatePersonalNote = ref(localStorage.getItem('candidate_personal_note_' + (interviewId || 'default')) || '')
watch(candidatePersonalNote, (val) => {
  localStorage.setItem('candidate_personal_note_' + (interviewId || 'default'), val || '');
})

const showReminder = ref(true)

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
  <div class="h-screen flex flex-col bg-slate-50 font-sans text-slate-800 overflow-hidden relative">
    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />
    <Toast v-if="isLeaving" type="info" message="Đang rời phòng phỏng vấn..." :duration="1500" />
    <!-- Cảnh báo quyền/thiết bị media: duration=0 để bám lại tới khi user xử lý -->
    <Toast v-if="liveKitError" type="error" :message="liveKitError" :duration="0" @close="clearError()" />

    <!-- Top Header -->
    <header class="h-12 bg-white border-b border-slate-200 flex items-center justify-between px-4 z-20 shadow-sm shrink-0">
      <div class="flex items-center gap-4">
        <div>
          <h1 class="text-base font-bold text-slate-800 truncate max-w-[300px] md:max-w-[500px]">
            Phỏng vấn: {{ interviewInfo?.job_title || 'Vị trí ứng tuyển' }}
          </h1>
        </div>
      </div>
      <div class="flex items-center gap-3 bg-slate-100 px-3 py-1 rounded-full border border-slate-200 shadow-inner">
        <span class="text-xs font-medium text-[var(--success)] flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full bg-[var(--success)]" :class="{ 'animate-pulse': roomStore.isConnected }"></span>
          <span class="hidden sm:inline">{{ roomStore.isConnected ? 'Đường truyền tốt' : 'Mất kết nối' }}</span>
        </span>
        <span class="w-px h-3.5 bg-slate-300"></span>
        <span class="text-xs font-mono font-bold text-slate-700 tracking-wider flex items-center gap-1.5">
          <span v-if="roomStore.status === 'active'" class="text-rose-600 flex items-center gap-1.5">
            <span class="w-2 h-2 rounded-full bg-rose-500 animate-pulse"></span>
            Đang Live
          </span>
          <span v-else-if="roomStore.status === 'waiting'" class="text-amber-600">Đang chờ</span>
          <span v-else-if="roomStore.status === 'paused'" class="text-amber-600">⏸️ Tạm dừng</span>
          <span v-else-if="roomStore.status === 'completed'" class="text-[var(--success)]">✅ Đã kết thúc</span>
          <span v-else>{{ roomStore.status }}</span>
          
          <span v-if="roomStore.startedAt" class="text-[var(--primary)] font-mono ml-1.5">{{ elapsedFormatted }}</span>
        </span>
      </div>
    </header>

    <!-- Main Workspace -->
    <div class="flex flex-1 overflow-hidden relative z-10 bg-slate-900">
      
      <!-- Left: Video Area -->
      <div class="flex-1 flex flex-col p-0.5 md:p-1 gap-1 relative transition-all duration-300 min-w-0">
        
        <!-- Video Grid (Fills almost 100% height) -->
        <div 
          @mousemove="resetControlBarTimer"
          @click="resetControlBarTimer"
          class="flex-1 relative rounded-xl overflow-hidden bg-slate-100 shadow-2xl border border-slate-200 ring-1 ring-white/5 flex items-center justify-center"
        >
          
          <!-- Floating Notification Banner inside video top-left (Dismissible) -->
          <div v-if="showReminder" class="absolute top-4 left-4 z-10 max-w-lg bg-cyan-50 border border-cyan-200 rounded-xl px-4 py-2.5 flex items-center gap-3 text-cyan-800 text-xs md:text-sm shadow-xl transition-all">
            <div class="w-7 h-7 rounded-lg bg-cyan-100 flex items-center justify-center shrink-0">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-cyan-600" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" /></svg>
            </div>
            <p class="flex-1"><strong class="text-cyan-700 font-semibold">Lời nhắc:</strong> Hãy trả lời tự nhiên. Nhà tuyển dụng sẽ dẫn dắt buổi phỏng vấn.</p>
            <button @click="showReminder = false" class="text-slate-400 hover:text-white p-1 transition-colors">
              <X class="w-4 h-4" />
            </button>
          </div>

          <!-- Recruiter Video (Main View) -->
          <div class="absolute inset-0 flex items-center justify-center">
            <!-- Thẻ video thật cho LiveKit (Remote - Recruiter) -->
            <video ref="remoteVideoEl" autoplay playsinline class="w-full h-full object-cover"></video>
            
            <!-- Fallback if no video -->
            <div v-if="!recruiterParticipant || recruiterParticipant.connection_state === 'offline'" class="absolute inset-0 flex flex-col items-center justify-center bg-slate-100 z-0">
              <div class="w-24 h-24 rounded-full bg-slate-200 border-2 border-slate-300 flex items-center justify-center mb-4 text-4xl font-bold text-slate-500 shadow-inner">
                R
              </div>
              <p class="text-slate-500 font-medium animate-pulse">Đang chờ nhà tuyển dụng kết nối...</p>
            </div>
            
            <!-- Recruiter Name Badge -->
            <div v-if="recruiterParticipant" class="absolute bottom-6 left-6 bg-black/60 backdrop-blur-md text-white px-4 py-2 rounded-xl text-sm font-medium border border-white/10 z-10 shadow-lg flex items-center gap-2">
              <div class="w-2 h-2 rounded-full bg-[var(--success)]"></div>
              {{ recruiterParticipant.display_name }} (Nhà Tuyển Dụng)
            </div>
          </div>
          
          <!-- Candidate Self View (PiP) -->
          <div class="absolute top-6 right-6 w-48 md:w-60 aspect-video bg-slate-200 rounded-xl overflow-hidden border-2 border-slate-300 shadow-lg z-20 group hover:scale-105 transition-all duration-300 cursor-move">
             <!-- Thẻ video thật cho LiveKit (Local) -->
             <video ref="localVideoEl" autoplay playsinline muted class="w-full h-full object-cover transform -scale-x-100"></video>
             
             <!-- Self Name Badge -->
             <div class="absolute bottom-2 left-2 bg-white/80 text-slate-800 px-2 py-1 rounded-lg text-xs font-medium opacity-0 group-hover:opacity-100 transition-opacity">
               Bạn
             </div>
             <!-- Mic Status icon on self view -->
             <div class="absolute top-2 right-2 w-6 h-6 rounded-full bg-white/80 flex items-center justify-center">
               <Mic v-if="isMicOn" class="w-3 h-3 text-[var(--success)]" />
               <MicOff v-else class="w-3 h-3 text-rose-500" />
             </div>
          </div>

          <div 
            @mouseenter="keepControlBarVisible"
            @mouseleave="resetControlBarTimer"
            :style="{
              opacity: isControlBarVisible ? 1 : 0,
              transform: isControlBarVisible ? 'translateX(-50%) translateY(0)' : 'translateX(-50%) translateY(24px)',
              pointerEvents: isControlBarVisible ? 'auto' : 'none',
              transition: 'all 0.35s cubic-bezier(0.4, 0, 0.2, 1)'
            }"
            class="absolute bottom-6 left-1/2 bg-white border border-slate-200 px-4 py-2 rounded-full shadow-lg flex items-center gap-3 z-30"
          >
            <button :title="isMicOn ? 'Tắt micro' : 'Bật micro'" @click="() => { toggleMic(); roomStore.updateMediaStatus(isMicOn, isCameraOn); }" class="w-11 h-11 rounded-full flex items-center justify-center transition-all duration-200 shadow-md" :class="isMicOn ? 'bg-slate-100 hover:bg-slate-200 text-slate-700' : 'bg-rose-500 hover:bg-rose-600 text-white'">
              <Mic v-if="isMicOn" class="w-5 h-5" />
              <MicOff v-else class="w-5 h-5" />
            </button>
            
            <button :title="isCameraOn ? 'Tắt camera' : 'Bật camera'" @click="() => { toggleCamera(); roomStore.updateMediaStatus(isMicOn, isCameraOn); }" class="w-11 h-11 rounded-full flex items-center justify-center transition-all duration-200 shadow-md" :class="isCameraOn ? 'bg-slate-100 hover:bg-slate-200 text-slate-700' : 'bg-rose-500 hover:bg-rose-600 text-white'">
              <Video v-if="isCameraOn" class="w-5 h-5" />
              <VideoOff v-else class="w-5 h-5" />
            </button>

            <button :title="isScreenSharing ? 'Dừng chia sẻ màn hình' : 'Chia sẻ màn hình'" @click="toggleScreenShare" class="w-11 h-11 rounded-full flex items-center justify-center transition-all duration-200 shadow-md" :class="isScreenSharing ? 'bg-[var(--primary)] text-white' : 'bg-slate-100 hover:bg-slate-200 text-slate-700'">
              <MonitorUp class="w-5 h-5" />
            </button>

            <button title="Mở khung Chat" @click="() => { activeTab = 'chat'; isPanelExpanded = true; }" class="w-11 h-11 rounded-full flex items-center justify-center transition-all duration-200 shadow-md" :class="activeTab === 'chat' && isPanelExpanded ? 'bg-[var(--primary)] text-white' : 'bg-slate-100 hover:bg-slate-200 text-slate-700'">
              <MessageSquare class="w-5 h-5" />
            </button>
            
            <div class="w-px h-7 bg-slate-300 mx-1"></div>
            
            <button @click="handleLeave" class="h-11 px-6 rounded-full bg-rose-600 hover:bg-rose-500 text-white font-bold shadow-lg shadow-rose-600/30 flex items-center gap-2 transition-all hover:scale-105 active:scale-95 text-sm">
              <PhoneOff class="w-4 h-4" /> <span class="hidden sm:inline">Rời phòng</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Right: Side Panel (Collapsible: Width 340px when expanded, 48px when collapsed) -->
      <div class="bg-white border-l border-slate-200 flex flex-col shrink-0 shadow-sm relative z-20 transition-all duration-300" :style="{ width: isPanelExpanded ? '340px' : '48px' }">
        
        <!-- Collapsed Sidebar Icon Mode -->
        <div v-if="!isPanelExpanded" class="flex flex-col items-center py-3 gap-2.5 h-full bg-slate-50 border-l border-slate-200">
          <button @click="isPanelExpanded = true" class="w-9 h-9 rounded-xl bg-white hover:bg-slate-100 text-slate-600 border border-slate-200 flex items-center justify-center transition-colors shadow-sm" title="Mở rộng bảng công cụ">
            <ChevronsLeft class="w-4 h-4" />
          </button>
          <div class="w-6 h-px bg-slate-200 my-0.5"></div>
          <button @click="() => { activeTab = 'info'; isPanelExpanded = true; }" class="w-9 h-9 rounded-xl flex items-center justify-center transition-colors relative" :class="activeTab === 'info' ? 'bg-[var(--primary-light)] text-[var(--primary)] border border-[var(--primary-light)] shadow-sm font-bold' : 'text-slate-500 hover:bg-slate-100 hover:text-slate-800'" title="JD & Ghi chú">
            <FileText class="w-4 h-4" />
          </button>
          <button @click="() => { activeTab = 'chat'; isPanelExpanded = true; }" class="w-9 h-9 rounded-xl flex items-center justify-center transition-colors relative" :class="activeTab === 'chat' ? 'bg-[var(--primary-light)] text-[var(--primary)] border border-[var(--primary-light)] shadow-sm font-bold' : 'text-slate-500 hover:bg-slate-100 hover:text-slate-800'" title="Chat">
            <MessageSquare class="w-4 h-4" />
            <span v-if="chatStore.messages.length > 0" class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-rose-500 border border-white"></span>
          </button>
        </div>

        <!-- Expanded Sidebar Mode -->
        <template v-else>
          <!-- Tabs Header -->
          <div class="flex border-b border-slate-200 shrink-0 bg-slate-50 items-center justify-between pr-3">
            <div class="flex flex-1">
              <button @click="activeTab = 'info'" class="flex-1 py-3.5 flex items-center justify-center gap-2 font-medium text-sm transition-colors border-b-2" :class="activeTab === 'info' ? 'text-[var(--primary)] border-[var(--primary)] bg-[var(--primary-light)]/60 font-semibold' : 'text-slate-500 border-transparent hover:bg-slate-100 hover:text-slate-800'">
                <FileText class="w-4 h-4" /> JD & Ghi chú
              </button>
              <button @click="activeTab = 'chat'" class="flex-1 py-3.5 flex items-center justify-center gap-2 font-medium text-sm transition-colors border-b-2 relative" :class="activeTab === 'chat' ? 'text-[var(--primary)] border-[var(--primary)] bg-[var(--primary-light)]/60 font-semibold' : 'text-slate-500 border-transparent hover:bg-slate-100 hover:text-slate-800'">
                <MessageSquare class="w-4 h-4" /> Chat
                <span v-if="chatStore.messages.length > 0" class="absolute top-2.5 right-6 w-2 h-2 rounded-full bg-rose-500 border-2 border-white"></span>
              </button>
            </div>
            <button @click="isPanelExpanded = false" class="w-8 h-8 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-200/60 flex items-center justify-center transition-colors ml-1 shrink-0" title="Thu gọn bảng công cụ">
              <ChevronsRight class="w-4 h-4" />
            </button>
          </div>

          <!-- Tab Content Area -->
          <div class="flex-1 overflow-y-auto custom-scrollbar relative flex flex-col bg-white">
            
            <!-- Tab: Info & Notes -->
            <div v-show="activeTab === 'info'" class="flex-1 p-5 flex flex-col gap-6">
              <div>
                <h3 class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-3 flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-[var(--primary)]" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M6 2a2 2 0 00-2 2v12a2 2 0 002 2h8a2 2 0 002-2V7.414A2 2 0 0015.414 6L12 2.586A2 2 0 0010.586 2H6zm5 6a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V8z" clip-rule="evenodd" /></svg>
                  Yêu cầu công việc
                </h3>
                <div class="bg-slate-50 p-3.5 rounded-xl border border-slate-200 text-slate-700">
                  <p class="text-sm leading-relaxed">{{ interviewInfo?.job_title ? `Vị trí ứng tuyển: ${interviewInfo.job_title}` : 'Phát triển và hoàn thiện các tính năng theo yêu cầu của nhà tuyển dụng.' }}</p>
                </div>
              </div>
              
              <div class="flex-1 flex flex-col">
                <h3 class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-3 flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-[var(--success)]" viewBox="0 0 20 20" fill="currentColor"><path d="M17.414 2.586a2 2 0 00-2.828 0L7 10.172V13h2.828l7.586-7.586a2 2 0 000-2.828z" /><path fill-rule="evenodd" d="M2 6a2 2 0 012-2h4a1 1 0 010 2H4v10h10v-4a1 1 0 112 0v4a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" clip-rule="evenodd" /></svg>
                  Ghi chú cá nhân của bạn
                </h3>
                <textarea 
                  v-model="candidatePersonalNote"
                  class="flex-1 w-full bg-white border border-slate-300 rounded-xl p-3.5 text-sm text-slate-800 placeholder-slate-400 focus:ring-2 focus:ring-[var(--accent)]/20 focus:border-[var(--accent)] outline-none resize-none transition-all custom-scrollbar shadow-inner" 
                  placeholder="Viết nháp ý chính hoặc ghi chú nhanh tại đây (Chỉ một mình bạn nhìn thấy)..."
                ></textarea>
              </div>
            </div>

            <!-- Tab: Chat -->
            <div v-show="activeTab === 'chat'" class="flex-1 flex flex-col h-full bg-slate-50/50">
              <div class="flex-1 overflow-y-auto p-4 flex flex-col gap-3 custom-scrollbar">
                <div v-if="chatStore.messages.length === 0" class="flex flex-col items-center justify-center h-full text-center text-slate-400 text-sm py-10 gap-2">
                  <MessageSquare class="w-8 h-8 text-slate-300" />
                  <span>Chưa có tin nhắn nào.<br/>Hãy gửi tin nhắn đầu tiên để trò chuyện!</span>
                </div>

                <div v-for="msg in chatStore.messages" :key="msg.message_id" class="flex flex-col gap-1" :class="(msg.sender_type === 'candidate' || msg.sender_type === 'local') ? 'items-end self-end' : 'items-start self-start'">
                  <span class="text-[11px] text-slate-500 font-medium px-1">
                    {{ (msg.sender_type === 'candidate' || msg.sender_type === 'local') ? 'Bạn (Ứng viên)' : (msg.sender_name || 'Nhà tuyển dụng') }}
                  </span>
                  <div 
                    class="px-3.5 py-2 rounded-2xl text-sm max-w-[85%] shadow-sm leading-relaxed"
                    :class="(msg.sender_type === 'candidate' || msg.sender_type === 'local') ? 'bg-[var(--primary)] text-white rounded-tr-sm shadow-md' : 'bg-white text-slate-800 border border-slate-200 rounded-tl-sm shadow-sm'"
                  >
                    {{ msg.message }}
                  </div>
                </div>
              </div>
              
              <div class="p-3 bg-white border-t border-slate-200 shrink-0">
                <div class="flex items-center gap-2 bg-slate-100 border border-slate-300 rounded-full py-1 pl-4 pr-1.5 focus-within:border-[var(--primary)] focus-within:ring-2 focus-within:ring-[var(--accent)]/20 transition-all shadow-inner">
                  <input 
                    type="text" 
                    v-model="chatInput" 
                    @keyup.enter="handleSendMessage"
                    class="flex-1 bg-transparent border-none text-sm text-slate-800 outline-none placeholder-slate-400 py-1.5" 
                    maxlength="2000"
                    placeholder="Nhập tin nhắn..."
                  />
                  <button 
                    @click="handleSendMessage"
                    class="w-8 h-8 rounded-full bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white flex items-center justify-center transition-all hover:scale-105 active:scale-95 shrink-0 shadow-md"
                    title="Gửi tin nhắn"
                  >
                    <Send class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- Leave Modal -->
    <Modal :isOpen="showLeaveModal" @close="showLeaveModal = false" title="Xác nhận rời phòng">
      <div class="p-2">
        <p class="text-gray-700 mb-6 font-medium">Bạn có chắc chắn muốn kết thúc buổi phỏng vấn và rời phòng?</p>
        <div class="flex justify-end gap-3">
          <button @click="showLeaveModal = false" class="px-5 py-2.5 rounded-full border border-gray-300 text-gray-700 font-bold hover:bg-gray-50 transition-colors">
            Tiếp tục ở lại
          </button>
          <button @click="confirmLeave" class="px-5 py-2.5 rounded-full bg-rose-600 text-white font-bold hover:bg-rose-500 shadow-md shadow-rose-600/20 transition-colors flex items-center gap-2">
            <PhoneOff class="w-4 h-4" /> Rời khỏi phòng
          </button>
        </div>
      </div>
    </Modal>

    <Modal :isOpen="showNavWarningModal" @close="handleCandidateCancelLeave" title="Cảnh báo rời phòng phỏng vấn">
      <div class="p-2">
        <p class="text-gray-700 mb-6 font-medium">Bạn đang trong phòng phỏng vấn. Bạn có chắc chắn muốn rời đi hoặc chuyển sang trang khác không?</p>
        <div class="flex justify-end gap-3">
          <button @click="handleCandidateCancelLeave" class="px-5 py-2.5 rounded-full border border-gray-300 text-gray-700 font-bold hover:bg-gray-50 transition-colors">
            Ở lại phòng
          </button>
          <button @click="handleCandidateConfirmLeave" class="px-5 py-2.5 rounded-full bg-amber-500 text-slate-900 font-bold hover:bg-amber-400 shadow-md transition-colors flex items-center gap-2">
            Rời khỏi phòng
          </button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
/* Scrollbar customizing for the panel */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #475569;
  border-radius: 20px;
}
.custom-scrollbar:hover::-webkit-scrollbar-thumb {
  background-color: #64748b;
}
</style>
