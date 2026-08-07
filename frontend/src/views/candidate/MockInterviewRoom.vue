<script setup>
import { ref, onMounted, onBeforeUnmount, computed, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mic, Square, PhoneOff, Loader2, Type, Send } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'
import { useGeminiLive } from '../../composables/useGeminiLive'

const router = useRouter()
const route = useRoute()
const mockId = route.query.mock_id
const role = route.query.role || 'Software Developer'
const level = route.query.level || 'middle'

const {
  connected, aiSpeaking, listening, error, transcript, audioLevel, isRecording,
  start, stop, commitTurn, sendText,
} = useGeminiLive()

// Mouth opening (px) lip-synced to the AI voice loudness.
const mouthOpen = computed(() => 3 + audioLevel.value * 26)
const mouthWidth = computed(() => 34 - audioLevel.value * 6)
const mouthRadius = computed(() => Math.min(mouthOpen.value / 2, mouthWidth.value / 2))

const showEndModal = ref(false)
const showText = ref(false)
const textAnswer = ref('')
const ending = ref(false)
const toast = ref(null)
const transcriptBox = ref(null)

const visibleTranscript = computed(() => {
  return transcript.value.filter(turn => turn.text && turn.text.trim())
})

const statusLabel = computed(() => {
  if (error.value) return error.value
  if (!connected.value) return 'Đang kết nối...'
  if (aiSpeaking.value) return 'AI đang nói...'
  if (isRecording.value) return 'Đang nghe bạn nói...'
  if (listening.value) return 'Sẵn sàng — nhấn giữ để trả lời'
  return 'Đang chuẩn bị...'
})

// Auto-scroll transcript to the latest line.
watch(transcript, async () => {
  await nextTick()
  if (transcriptBox.value) transcriptBox.value.scrollTop = transcriptBox.value.scrollHeight
}, { deep: true })

const connectFailed = ref(false)

onMounted(async () => {
  if (!mockId) {
    toast.value = { type: 'error', message: 'Thiếu mã phiên luyện tập.' }
    connectFailed.value = true
    return
  }
  try {
    await start({ role, level })
  } catch (e) {
    connectFailed.value = true
    toast.value = { type: 'error', message: e?.message || 'Không kết nối được AI. Bạn có thể thoát và thử lại.' }
  }
})

onBeforeUnmount(() => stop())

// Push-to-talk: hold to speak, release to commit the turn to the AI.
const startTalk = () => {
  if (!connected.value || aiSpeaking.value) return
  isRecording.value = true
}
const stopTalk = () => {
  if (!isRecording.value) return
  isRecording.value = false
  commitTurn()
}

const handleSendText = () => {
  const answer = textAnswer.value.trim()
  if (!answer) {
    toast.value = { type: 'error', message: 'Vui lòng nhập câu trả lời.' }
    return
  }
  if (!connected.value) {
    toast.value = { type: 'error', message: 'Không thể gửi khi AI chưa kết nối.' }
    return
  }
  if (aiSpeaking.value) {
    toast.value = { type: 'warning', message: 'Vui lòng chờ AI nói xong trước khi trả lời.' }
    return
  }
  if (answer.length > 5000) {
    toast.value = { type: 'error', message: 'Câu trả lời không được vượt quá 5.000 ký tự.' }
    return
  }
  sendText(answer)
  textAnswer.value = ''
}

const handleEnd = async () => {
  ending.value = true
  try {
    // Persist the spoken conversation as mock messages, then finalize.
    const turns = visibleTranscript.value.slice(-200).map(turn => ({
      role: turn.role,
      text: turn.text.trim().slice(0, 5000),
    }))
    await mockService.saveLiveTranscript(mockId, turns).catch(() => {})
    await mockService.endMockInterview(mockId).catch(() => {})
  } catch (e) {
    console.error('Lỗi kết thúc phiên', e)
    toast.value = { type: 'warning', message: 'Phiên đã dừng; một phần dữ liệu có thể chưa được lưu.' }
  } finally {
    stop()
    ending.value = false
    router.push({ path: '/mock-results', query: { mock_id: mockId } })
  }
}

const leaveWithoutEnd = () => {
  stop()
  router.push('/mock-results')
}
</script>
<template>
  <div class="h-screen flex flex-col bg-white font-sans overflow-hidden text-slate-900">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
    <!-- Header -->
    <div class="h-16 bg-white border-b border-slate-200 flex items-center justify-between px-6 shrink-0">
      <div>
        <h1 class="text-lg font-bold">Luyện tập AI · {{ role }}</h1>
        <p class="text-xs text-slate-500 capitalize">Trình độ: {{ level }}</p>
      </div>
      <button @click="(connectFailed || error) ? leaveWithoutEnd() : (showEndModal = true)"
        class="text-rose-500 hover:text-white hover:bg-rose-500 font-semibold px-4 py-2 rounded-lg transition-colors flex items-center gap-2">
        <PhoneOff class="w-4 h-4" /> {{ (connectFailed || error) ? 'Thoát' : 'Kết thúc' }}
      </button>
    </div>

    <div v-if="connectFailed || error" class="mx-6 mt-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700 flex items-center justify-between gap-3">
      <span>{{ error || 'Không kết nối được AI luyện tập. Bạn có thể thoát và thử lại sau.' }}</span>
      <button class="font-semibold underline shrink-0" @click="leaveWithoutEnd">Về lịch sử</button>
    </div>

    <!-- Main -->
    <div class="flex-1 min-h-0 grid grid-cols-1 lg:grid-cols-[minmax(280px,360px)_minmax(0,1fr)] gap-4 bg-slate-50 p-4 lg:p-6 overflow-hidden">
      <!-- AI interviewer -->
      <section class="rounded-2xl border border-slate-200 bg-white shadow-sm flex flex-col items-center justify-center p-6 min-h-0">
        <div class="relative flex flex-col items-center">
          <div class="relative">
            <div v-if="aiSpeaking" class="absolute inset-0 rounded-full bg-cyan-500/30 animate-ping" style="animation-duration:1.2s"></div>
            <div v-if="aiSpeaking" class="absolute -inset-4 rounded-full border-2 border-cyan-400/40 animate-pulse"></div>
            <div class="relative z-10 w-36 h-36 sm:w-40 sm:h-40 rounded-full bg-cyan-600 flex items-center justify-center shadow-lg border-4 border-cyan-100 transition-transform duration-150"
                 :class="aiSpeaking ? 'scale-105' : 'scale-100'">
              <svg viewBox="0 0 120 120" class="w-28 h-28 sm:w-32 sm:h-32">
                <ellipse cx="44" cy="50" :rx="7" :ry="aiSpeaking ? 8 : 7" fill="#fff" />
                <ellipse cx="76" cy="50" :rx="7" :ry="aiSpeaking ? 8 : 7" fill="#fff" />
                <circle cx="44" cy="51" r="3.5" fill="#164e63" />
                <circle cx="76" cy="51" r="3.5" fill="#164e63" />
                <rect :x="60 - mouthWidth / 2" :y="78 - mouthOpen / 2"
                      :width="mouthWidth" :height="mouthOpen" :rx="mouthRadius"
                      fill="#164e63" stroke="#fff" stroke-width="2" />
              </svg>
            </div>
          </div>
          <h3 class="mt-6 text-xl font-bold">AI Interviewer</h3>
          <div class="mt-2 flex items-center gap-2 px-4 py-1.5 rounded-full text-sm font-semibold text-center"
               :class="error ? 'bg-rose-50 text-rose-600'
                 : aiSpeaking ? 'bg-cyan-50 text-cyan-700'
                 : isRecording ? 'bg-emerald-50 text-emerald-700'
                 : 'bg-slate-100 text-slate-600'">
            <Loader2 v-if="!connected && !error" class="w-4 h-4 animate-spin" />
            {{ statusLabel }}
          </div>
          <p class="mt-6 text-center text-sm leading-relaxed text-slate-500">
            Hãy trả lời tự nhiên. AI sẽ đặt câu hỏi và phản hồi theo từng lượt.
          </p>
        </div>
      </section>

      <!-- Live transcript -->
      <section class="rounded-2xl border border-slate-200 bg-white shadow-sm flex flex-col min-h-0 overflow-hidden">
        <div class="shrink-0 border-b border-slate-200 px-5 py-4">
          <h2 class="font-bold text-slate-900">Cuộc trò chuyện</h2>
          <p class="mt-1 text-xs text-slate-500">Nội dung trao đổi với AI sẽ hiển thị tại đây</p>
        </div>
        <div ref="transcriptBox" class="flex-1 min-h-0 overflow-y-auto space-y-4 p-5">
          <div v-if="visibleTranscript.length === 0" class="h-full flex items-center justify-center text-center text-slate-500 text-sm">
            Cuộc trò chuyện sẽ hiển thị tại đây khi AI bắt đầu nói...
          </div>
          <div v-for="(t, i) in visibleTranscript" :key="i"
               class="flex" :class="t.role === 'ai' ? 'justify-start' : 'justify-end'">
            <div class="max-w-[92%] lg:max-w-[82%] rounded-2xl px-4 py-3 text-[15px] leading-relaxed shadow-sm whitespace-pre-wrap break-words"
                 :class="t.role === 'ai' ? 'bg-slate-100 text-slate-800 rounded-tl-sm' : 'bg-[var(--primary)] text-white rounded-tr-sm'">
              <div class="mb-1 text-[11px] font-semibold uppercase tracking-wide opacity-60">
                {{ t.role === 'ai' ? 'AI Interviewer' : 'Bạn' }}
              </div>
              {{ t.text }}
            </div>
          </div>
        </div>
      </section>
    </div>

    <!-- Control bar -->
    <div class="shrink-0 bg-white border-t border-slate-200 p-6">
      <div class="max-w-2xl mx-auto flex flex-col items-center gap-4">
        <!-- Push to talk -->
        <button
          @mousedown="startTalk" @mouseup="stopTalk" @mouseleave="stopTalk"
          @touchstart.prevent="startTalk" @touchend.prevent="stopTalk"
          :disabled="!connected || aiSpeaking"
          class="w-20 h-20 rounded-full flex items-center justify-center transition-all duration-200 shadow-md disabled:opacity-40 disabled:cursor-not-allowed select-none text-white"
          :class="isRecording ? 'bg-rose-500 scale-110 shadow-rose-500/30' : 'bg-[var(--primary)] hover:bg-[var(--primary-hover)] hover:scale-105 shadow-md'">
          <Square v-if="isRecording" class="w-7 h-7 fill-current" />
          <Mic v-else class="w-8 h-8" />
        </button>
        <p class="text-xs text-slate-500">Nhấn giữ để nói, thả ra để AI trả lời</p>

        <!-- Text fallback -->
        <button @click="showText = !showText" class="text-xs text-slate-500 hover:text-slate-700 flex items-center gap-1">
          <Type class="w-3.5 h-3.5" /> Trả lời bằng văn bản
        </button>
        <div v-if="showText" class="w-full flex gap-2">
          <input v-model="textAnswer" maxlength="5000" @keyup.enter="handleSendText"
            placeholder="Nhập câu trả lời..."
            class="flex-1 bg-white border border-slate-300 rounded-xl px-4 py-2.5 outline-none focus:border-[var(--accent)] text-slate-900 placeholder-slate-400" />
          <button @click="handleSendText" :disabled="!textAnswer.trim()"
            class="bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white disabled:bg-slate-300 px-4 rounded-xl transition-colors">
            <Send class="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>

    <!-- End modal -->
    <Modal :isOpen="showEndModal" @close="showEndModal = false" title="Kết thúc buổi luyện tập">
      <p class="text-gray-600 mb-6">Kết thúc và xem báo cáo tổng hợp? Cuộc trò chuyện sẽ được lưu lại.</p>
      <div class="flex justify-end gap-3">
        <button class="px-4 py-2 text-gray-600 font-medium hover:bg-gray-100 rounded-lg" @click="showEndModal = false">Tiếp tục luyện</button>
        <button class="px-4 py-2 bg-rose-500 hover:bg-rose-600 text-white font-medium rounded-lg flex items-center gap-2 disabled:opacity-60"
          :disabled="ending" @click="handleEnd">
          <Loader2 v-if="ending" class="w-4 h-4 animate-spin" /> Kết thúc & xem báo cáo
        </button>
      </div>
    </Modal>
  </div>
</template>
