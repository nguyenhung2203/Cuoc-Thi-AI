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
  connected, aiSpeaking, listening, error, transcript, audioLevel,
  start, stop, commitTurn, sendText,
} = useGeminiLive()

// Mouth opening (px) lip-synced to the AI voice loudness.
const mouthOpen = computed(() => 3 + audioLevel.value * 26)
const mouthWidth = computed(() => 34 - audioLevel.value * 6)
const mouthRadius = computed(() => Math.min(mouthOpen.value / 2, mouthWidth.value / 2))

const speaking = ref(false)     // candidate is holding the talk button
const showEndModal = ref(false)
const showText = ref(false)
const textAnswer = ref('')
const ending = ref(false)
const toast = ref(null)
const transcriptBox = ref(null)

const statusLabel = computed(() => {
  if (error.value) return error.value
  if (!connected.value) return 'Đang kết nối...'
  if (aiSpeaking.value) return 'AI đang nói...'
  if (speaking.value) return 'Đang nghe bạn nói...'
  if (listening.value) return 'Sẵn sàng — nhấn giữ để trả lời'
  return 'Đang chuẩn bị...'
})

// Auto-scroll transcript to the latest line.
watch(transcript, async () => {
  await nextTick()
  if (transcriptBox.value) transcriptBox.value.scrollTop = transcriptBox.value.scrollHeight
}, { deep: true })

onMounted(async () => {
  if (!mockId) {
    toast.value = { type: 'error', message: 'Thiếu mã phiên luyện tập.' }
    return
  }
  try {
    await start({ role, level })
  } catch {
    // error surfaced via `error` ref
  }
})

onBeforeUnmount(() => stop())

// Push-to-talk: hold to speak, release to commit the turn to the AI.
const startTalk = () => {
  if (!connected.value || aiSpeaking.value) return
  speaking.value = true
}
const stopTalk = () => {
  if (!speaking.value) return
  speaking.value = false
  commitTurn()
}

const handleSendText = () => {
  if (!textAnswer.value.trim()) return
  sendText(textAnswer.value.trim())
  textAnswer.value = ''
}

const handleEnd = async () => {
  ending.value = true
  try {
    // Persist the spoken conversation as mock messages, then finalize.
    const turns = transcript.value.filter(t => t.text && t.text.trim())
    await mockService.saveLiveTranscript(mockId, turns).catch(() => {})
    await mockService.endMockInterview(mockId)
  } catch (e) {
    console.error('Lỗi kết thúc phiên', e)
  } finally {
    stop()
    router.push({ path: '/mock-results', query: { mock_id: mockId } })
  }
}
</script>
<template>
  <div class="h-screen flex flex-col bg-gradient-to-b from-slate-900 to-slate-950 font-sans overflow-hidden text-white">
    <!-- Header -->
    <div class="h-16 bg-slate-900/70 backdrop-blur-md border-b border-white/10 flex items-center justify-between px-6 shrink-0">
      <div>
        <h1 class="text-lg font-bold">Luyện tập AI · {{ role }}</h1>
        <p class="text-xs text-slate-400 capitalize">Trình độ: {{ level }}</p>
      </div>
      <button @click="showEndModal = true"
        class="text-rose-300 hover:text-white hover:bg-rose-500/80 font-semibold px-4 py-2 rounded-lg transition-colors flex items-center gap-2">
        <PhoneOff class="w-4 h-4" /> Kết thúc
      </button>
    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <!-- Main -->
    <div class="flex-1 flex flex-col items-center justify-center relative px-4 overflow-hidden">
      <!-- AI Avatar -->
      <div class="relative flex flex-col items-center mb-8">
        <div class="relative">
          <!-- Speaking ripples -->
          <div v-if="aiSpeaking" class="absolute inset-0 rounded-full bg-indigo-500/30 animate-ping" style="animation-duration:1.2s"></div>
          <div v-if="aiSpeaking" class="absolute -inset-4 rounded-full border-2 border-indigo-400/40 animate-pulse"></div>
          <!-- Talking face: eyes + mouth lip-synced to the AI voice -->
          <div class="relative z-10 w-40 h-40 rounded-full bg-gradient-to-tr from-indigo-500 to-purple-500 flex items-center justify-center shadow-[0_0_60px_rgba(99,102,241,0.6)] border-4 border-slate-700 transition-transform duration-150"
               :class="aiSpeaking ? 'scale-105' : 'scale-100'">
            <svg viewBox="0 0 120 120" class="w-32 h-32">
              <!-- Eyes -->
              <ellipse cx="44" cy="50" :rx="7" :ry="aiSpeaking ? 8 : 7" fill="#fff" />
              <ellipse cx="76" cy="50" :rx="7" :ry="aiSpeaking ? 8 : 7" fill="#fff" />
              <circle cx="44" cy="51" r="3.5" fill="#1e1b4b" />
              <circle cx="76" cy="51" r="3.5" fill="#1e1b4b" />
              <!-- Mouth: height follows audioLevel -->
              <rect :x="60 - mouthWidth / 2" :y="78 - mouthOpen / 2"
                    :width="mouthWidth" :height="mouthOpen"
                    :rx="mouthRadius"
                    fill="#1e1b4b" stroke="#fff" stroke-width="2" />
            </svg>
          </div>
        </div>
        <h3 class="mt-6 text-xl font-bold">AI Interviewer</h3>
        <div class="mt-2 flex items-center gap-2 px-4 py-1.5 rounded-full text-sm font-semibold"
             :class="error ? 'bg-rose-500/20 text-rose-300'
               : aiSpeaking ? 'bg-indigo-500/20 text-indigo-200'
               : speaking ? 'bg-emerald-500/20 text-emerald-200'
               : 'bg-white/10 text-slate-300'">
          <Loader2 v-if="!connected && !error" class="w-4 h-4 animate-spin" />
          {{ statusLabel }}
        </div>
      </div>

      <!-- Live transcript -->
      <div ref="transcriptBox" class="w-full max-w-2xl flex-1 max-h-[32vh] overflow-y-auto space-y-3 px-2 mb-4">
        <div v-if="transcript.length === 0" class="text-center text-slate-500 text-sm mt-6">
          Cuộc trò chuyện sẽ hiển thị tại đây khi AI bắt đầu nói...
        </div>
        <div v-for="(t, i) in transcript" :key="i"
             class="flex" :class="t.role === 'ai' ? 'justify-start' : 'justify-end'">
          <div class="max-w-[80%] rounded-2xl px-4 py-2.5 text-[15px] leading-relaxed"
               :class="t.role === 'ai' ? 'bg-white/10 text-slate-100 rounded-tl-sm' : 'bg-blue-600 text-white rounded-tr-sm'">
            {{ t.text }}
          </div>
        </div>
      </div>
    </div>

    <!-- Control bar -->
    <div class="shrink-0 bg-slate-900/80 backdrop-blur-xl border-t border-white/10 p-6">
      <div class="max-w-2xl mx-auto flex flex-col items-center gap-4">
        <!-- Push to talk -->
        <button
          @mousedown="startTalk" @mouseup="stopTalk" @mouseleave="stopTalk"
          @touchstart.prevent="startTalk" @touchend.prevent="stopTalk"
          :disabled="!connected || aiSpeaking"
          class="w-20 h-20 rounded-full flex items-center justify-center transition-all duration-200 shadow-xl disabled:opacity-40 disabled:cursor-not-allowed select-none"
          :class="speaking ? 'bg-rose-500 scale-110 shadow-rose-500/50' : 'bg-gradient-to-r from-blue-600 to-indigo-600 hover:scale-105 shadow-indigo-500/40'">
          <Square v-if="speaking" class="w-7 h-7 fill-current" />
          <Mic v-else class="w-8 h-8" />
        </button>
        <p class="text-xs text-slate-400">Nhấn giữ để nói, thả ra để AI trả lời</p>

        <!-- Text fallback -->
        <button @click="showText = !showText" class="text-xs text-slate-400 hover:text-white flex items-center gap-1">
          <Type class="w-3.5 h-3.5" /> Trả lời bằng văn bản
        </button>
        <div v-if="showText" class="w-full flex gap-2">
          <input v-model="textAnswer" @keyup.enter="handleSendText"
            placeholder="Nhập câu trả lời..."
            class="flex-1 bg-white/10 border border-white/20 rounded-xl px-4 py-2.5 outline-none focus:border-blue-400 text-white placeholder-slate-500" />
          <button @click="handleSendText" :disabled="!textAnswer.trim()"
            class="bg-blue-600 hover:bg-blue-700 disabled:bg-slate-600 px-4 rounded-xl transition-colors">
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
