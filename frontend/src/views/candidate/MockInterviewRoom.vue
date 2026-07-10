<script setup>
import { ref, onMounted, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '../../components/common/AppButton.vue'
import Card from '../../components/common/AppCard.vue'
import Badge from '../../components/common/AppBadge.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mic, ArrowRight, Play, Square, RefreshCw, Send, CheckCircle } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'

const router = useRouter()
const route = useRoute()
const mockId = route.query.mock_id

const recording = ref(false)
const questionIndex = ref(1)
const analyzing = ref(false)
const feedback = ref(null)
const textAnswer = ref('')
const showEndModal = ref(false)
const messages = ref([])
const chatContainer = ref(null)

const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)
const toast = ref(null)

// Lấy câu hỏi hiện tại (message AI cuối chưa có câu trả lời)
const currentQuestion = computed(() => {
  const aiMsgs = messages.value.filter(m => m.sender_type === 'ai' && !m.is_feedback)
  return aiMsgs.length > 0 ? aiMsgs[aiMsgs.length - 1] : null
})

const loadMessages = async (retryCount = 0) => {
  if (!mockId) return
  try {
    const msgs = await mockService.getMessages(mockId)
    messages.value = Array.isArray(msgs) ? msgs : []
    
    // Nếu messages trống (có thể do backend đang gen câu hỏi đầu tiên), thử lại sau 1s (tối đa 3 lần)
    if (messages.value.length === 0 && retryCount < 3) {
      setTimeout(() => loadMessages(retryCount + 1), 1000)
      return
    }

    // Parse feedback từ message AI có score_json
    const feedbackMsg = [...messages.value].reverse().find(m => m.sender_type === 'ai' && m.score_json)
    if (feedbackMsg) {
      try {
        const parsed = typeof feedbackMsg.score_json === 'string'
          ? JSON.parse(feedbackMsg.score_json)
          : feedbackMsg.score_json
        feedback.value = {
          score: `${parsed.score ?? '?'}/10`,
          message: feedbackMsg.content || 'AI đã đánh giá câu trả lời của bạn.',
          improvement: parsed.improvements ? parsed.improvements.join(', ') : ''
        }
      } catch { feedback.value = null }
    } else {
      feedback.value = null
    }
    scrollToBottom()
  } catch (err) {
    console.error('Lỗi tải tin nhắn', err)
  }
}

onMounted(() => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }
  loadMessages()
})

const scrollToBottom = () => {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })
}

const handleRecord = () => {
  // Mock recording logic for now since voice parsing is not requested yet
  if (!recording.value) {
    recording.value = true
    feedback.value = null
  } else {
    recording.value = false
    handleSendText() // Fallback to text send for now
  }
}

const handleSendText = async () => {
  if (!textAnswer.value && !recording.value) return
  
  const answerText = textAnswer.value || 'Câu trả lời ghi âm (chức năng voice chưa tích hợp)'
  
  // Immediately add user's message to UI and clear input
  messages.value.push({
    sender_type: 'user',
    content: answerText
  })
  textAnswer.value = ''
  scrollToBottom()
  
  analyzing.value = true
  feedback.value = null
  
  // Lấy question_id từ câu hỏi hiện tại
  const questionId = currentQuestion.value?.id || currentQuestion.value?.question_id || null

  try {
    const result = await mockService.submitAnswer(mockId, {
      question_id: questionId,
      answer_text: answerText
    })
    
    questionIndex.value++

    // Reload messages để có câu hỏi tiếp theo + feedback
    await loadMessages()

    // Parse feedback từ result.ai_message (nếu backend parse scoreJSON ra)
    if (result?.ai_message?.score_json) {
      const fb = typeof result.ai_message.score_json === 'string' 
        ? JSON.parse(result.ai_message.score_json) 
        : result.ai_message.score_json
      feedback.value = {
        score: `${fb.score ?? '?'}/10`,
        message: Array.isArray(fb.strengths) ? fb.strengths.join('. ') : 'AI đã đánh giá.',
        improvement: Array.isArray(fb.improvements) ? fb.improvements.join('. ') : ''
      }
    }
  } catch (error) {
    console.error('Lỗi gửi câu trả lời', error)
  } finally {
    analyzing.value = false
  }
}

const handleNext = async () => {
  if (!mockId) {
    router.push({ path: '/mock-results' })
    return
  }

  try {
    await mockService.endMockInterview(mockId)
  } catch (error) {
    console.error('Lỗi kết thúc mock interview:', error)
  } finally {
    router.push({ path: `/mock-results/${mockId}` })
  }
}

// STT: Speech Recognition
let recognition = null
if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
  const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition
  recognition = new SpeechRecognition()
  recognition.continuous = true
  recognition.interimResults = true
  recognition.lang = 'vi-VN'

  recognition.onresult = (event) => {
    let interimTranscript = ''
    let finalTranscript = ''
    for (let i = event.resultIndex; i < event.results.length; ++i) {
      if (event.results[i].isFinal) {
        finalTranscript += event.results[i][0].transcript
      } else {
        interimTranscript += event.results[i][0].transcript
      }
    }
    // Update the input field with ongoing transcription
    textAnswer.value = finalTranscript || interimTranscript
  }
  
  recognition.onerror = (event) => {
    console.error('Speech recognition error', event.error)
    recording.value = false
    toast.value = { type: 'error', message: 'Lỗi ghi âm: ' + event.error }
  }
  
  recognition.onend = () => {
    if (recording.value) {
      // Auto restart if it stopped but we are still in "recording" state
      recognition.start()
    }
  }
}

const toggleRecord = () => {
  if (!recognition) {
    toast.value = { type: 'error', message: 'Trình duyệt không hỗ trợ Web Speech API' }
    return
  }
  if (!recording.value) {
    recording.value = true
    feedback.value = null
    textAnswer.value = ''
    recognition.start()
  } else {
    recording.value = false
    recognition.stop()
    // Sau khi tắt thì tự động send luôn hoặc bắt user bấm Gửi?
    // Sẽ không tự động send, để user review text
  }
}

</script>

<template>
  <div class="h-screen flex flex-col bg-slate-50 font-sans overflow-hidden">
    <!-- Header -->
    <div class="h-16 bg-white/80 backdrop-blur-md border-b border-gray-200 flex items-center justify-between px-6 shadow-sm z-20">
      <div class="flex items-center gap-4">
        <div>
          <h1 class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600">Luyện tập AI: Frontend Developer</h1>
          <p class="text-sm font-medium text-gray-500 mt-0.5">Câu hỏi {{ questionIndex }} / 5</p>
        </div>
      </div>
      <button @click="showEndModal = true" class="text-rose-500 hover:text-rose-600 font-semibold px-4 py-2 hover:bg-rose-50 rounded-lg transition-colors flex items-center gap-2">
        <Square class="w-4 h-4" /> Kết thúc sớm
      </button>
    </div>

    <Modal :isOpen="showEndModal" @close="showEndModal = false" title="Kết thúc sớm">
      <p class="text-gray-600 mb-6">Bạn có chắc chắn muốn kết thúc bài thi sớm? Kết quả sẽ được tính trên những câu bạn đã trả lời.</p>
      <div class="flex justify-end gap-3">
        <button class="px-4 py-2 text-gray-600 font-medium hover:bg-gray-100 rounded-lg transition-colors" @click="showEndModal = false">Huỷ</button>
        <button class="px-4 py-2 bg-rose-500 hover:bg-rose-600 text-white font-medium rounded-lg shadow-sm transition-colors" @click="handleNext">Xác nhận</button>
      </div>
    </Modal>

    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />

    <!-- Main Content -->
    <div class="flex flex-1 overflow-hidden relative">
      
      <!-- Left: AI Interviewer -->
      <div class="w-[400px] flex flex-col bg-white border-r border-gray-200 shadow-[4px_0_24px_rgba(0,0,0,0.02)] z-10 relative">
        <!-- AI Avatar Area -->
        <div class="p-8 flex flex-col items-center justify-center border-b border-gray-100 relative overflow-hidden bg-gradient-to-b from-slate-900 to-slate-800 shrink-0">
          <!-- Animated glowing background -->
          <div class="absolute inset-0 opacity-30 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-indigo-500 via-transparent to-transparent animate-pulse" style="animation-duration: 3s;"></div>
          
          <div class="relative z-10 w-28 h-28 rounded-full bg-gradient-to-tr from-indigo-500 to-purple-500 flex items-center justify-center text-white mb-6 shadow-[0_0_40px_rgba(99,102,241,0.5)] border-4 border-slate-700">
            <span class="text-3xl font-black tracking-wider">AI</span>
            <!-- Speaking ripples -->
            <div v-if="!analyzing" class="absolute inset-0 rounded-full border-2 border-indigo-400 opacity-0 animate-[ping_2s_cubic-bezier(0,0,0.2,1)_infinite]"></div>
            <div v-if="!analyzing" class="absolute inset-0 rounded-full border-2 border-purple-400 opacity-0 animate-[ping_2.5s_cubic-bezier(0,0,0.2,1)_infinite]" style="animation-delay: 0.5s;"></div>
          </div>
          
          <h3 class="text-xl font-bold text-white mb-4 relative z-10">AI Interviewer</h3>
          
          <div class="relative z-10 flex items-center gap-2 px-5 py-2 rounded-full border shadow-inner backdrop-blur-md transition-all duration-300"
               :class="analyzing ? 'bg-indigo-500/20 border-indigo-400/30' : 'bg-emerald-500/20 border-emerald-400/30'">
            <template v-if="analyzing">
              <RefreshCw class="w-4 h-4 text-indigo-300 animate-spin" /> 
              <span class="text-sm font-semibold text-indigo-200">Đang phân tích...</span>
            </template>
            <template v-else>
              <div class="relative flex h-2.5 w-2.5 mr-1">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
              </div>
              <span class="text-sm font-semibold text-emerald-200">Đang lắng nghe</span>
            </template>
          </div>
        </div>
      </div>

      <!-- Right: Candidate Answer & Interaction -->
      <div class="flex-1 flex flex-col relative bg-slate-50/30">
        
        <div ref="chatContainer" class="flex-1 overflow-y-auto p-8 pb-32 scroll-smooth">
          <div class="max-w-3xl mx-auto space-y-8">
            
            <!-- Chat History -->
            <div class="space-y-6 mb-8">
              <div v-if="messages.length === 0" class="flex flex-col items-center justify-center text-center opacity-50 mt-20 animate-pulse">
                <Mic class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <p class="text-gray-500 font-medium">Hệ thống đang chuẩn bị câu hỏi đầu tiên...</p>
              </div>
              <div v-for="(msg, idx) in messages" :key="idx" class="animate-fade-in-up">
                <div v-if="msg.sender_type === 'ai'" class="flex gap-3">
                  <div class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center shrink-0 border border-indigo-200 mt-1 shadow-sm">
                    <span class="text-xs font-bold text-indigo-700">AI</span>
                  </div>
                  <div class="bg-white border border-gray-100 shadow-sm rounded-2xl rounded-tl-sm p-4 text-gray-800 text-[15px] leading-relaxed relative">
                    {{ msg.content }}
                    <!-- Indicate if this is the latest AI question without feedback -->
                    <span v-if="idx === messages.length - 1 && !analyzing && !msg.is_feedback" class="absolute -right-2 -bottom-2 flex h-3 w-3">
                      <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                      <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
                    </span>
                  </div>
                </div>
                <div v-else class="flex gap-3 justify-end">
                  <div class="bg-blue-600 text-white shadow-sm rounded-2xl rounded-tr-sm p-4 text-[15px] leading-relaxed max-w-[85%]">
                    {{ msg.content }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Feedback Area -->
            <div v-if="feedback" class="bg-white rounded-2xl shadow-xl border-t-4 border-t-emerald-500 overflow-hidden animate-fade-in-down mt-6">
              <div class="p-6">
                <div class="flex justify-between items-center mb-4">
                  <h3 class="text-lg font-bold text-gray-800 flex items-center gap-2">
                    <CheckCircle class="w-6 h-6 text-emerald-500" /> Phản hồi từ AI
                  </h3>
                  <div class="bg-emerald-100 text-emerald-700 font-bold px-3 py-1 rounded-lg text-sm">
                    Điểm: {{ feedback.score }}
                  </div>
                </div>
                <p class="text-gray-700 mb-6 leading-relaxed">{{ feedback.message }}</p>
                <div class="bg-amber-50 rounded-xl p-4 border-l-4 border-amber-400">
                  <p class="text-sm font-bold text-amber-800 mb-1 flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path d="M11 3a1 1 0 10-2 0v1a1 1 0 102 0V3zM15.657 5.757a1 1 0 00-1.414-1.414l-.707.707a1 1 0 001.414 1.414l.707-.707zM18 10a1 1 0 01-1 1h-1a1 1 0 110-2h1a1 1 0 011 1zM5.05 6.464A1 1 0 106.464 5.05l-.707-.707a1 1 0 00-1.414 1.414l.707.707zM5 10a1 1 0 01-1 1H3a1 1 0 110-2h1a1 1 0 011 1zM8 16v-1h4v1a2 2 0 11-4 0zM12 14c.015-.34.208-.646.477-.859a4 4 0 10-4.954 0c.27.213.462.519.476.859h4.002z" /></svg>
                    Gợi ý cải thiện:
                  </p>
                  <p class="text-amber-700 text-sm leading-relaxed">{{ feedback.improvement || 'Nên cung cấp thêm các ví dụ thực tế cụ thể.' }}</p>
                </div>
                
                <div class="flex justify-end mt-6">
                  <button @click="handleNext" class="bg-emerald-600 hover:bg-emerald-700 text-white font-bold py-2.5 px-6 rounded-xl shadow-lg hover:shadow-emerald-500/30 transition-all duration-300 transform hover:-translate-y-0.5 flex items-center gap-2">
                    {{ questionIndex < 5 ? 'Câu hỏi tiếp theo' : 'Xem báo cáo tổng hợp' }} <ArrowRight class="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Empty space filler when no feedback -->
            <div v-else class="h-full flex items-center justify-center text-center opacity-50 mt-20">
              <div>
                <Mic class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <p class="text-gray-500 font-medium">Bắt đầu ghi âm hoặc nhập câu trả lời của bạn</p>
              </div>
            </div>

          </div>
        </div>

        <!-- Floating Action Bar for Controls -->
        <div class="absolute bottom-0 left-0 right-0 bg-white/90 backdrop-blur-xl border-t border-gray-200 p-6 shadow-[0_-10px_40px_rgba(0,0,0,0.05)] z-20">
          <div class="max-w-4xl mx-auto flex items-end gap-6">
            
            <!-- Voice Record Button -->
            <div class="relative shrink-0 flex flex-col items-center group">
              <button @click="toggleRecord" 
                 class="w-16 h-16 rounded-full flex items-center justify-center transition-all duration-300 shadow-xl z-10"
                 :class="recording ? 'bg-rose-500 hover:bg-rose-600 shadow-rose-500/40 animate-pulse' : 'bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 shadow-indigo-500/30 hover:scale-105'">
                <Square v-if="recording" class="w-6 h-6 text-white fill-current" />
                <Mic v-else class="w-7 h-7 text-white" />
              </button>
              
              <!-- Voice ripples -->
              <div v-if="recording" class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full rounded-full border-2 border-rose-400 animate-ping" style="animation-duration: 1.5s;"></div>
              
              <div class="absolute -top-10 whitespace-nowrap text-xs font-bold text-gray-500 px-3 py-1 bg-white border border-gray-200 rounded-lg shadow-sm opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none">
                {{ recording ? 'Dừng ghi âm' : 'Trả lời bằng giọng nói' }}
              </div>
            </div>

            <!-- Text Input -->
            <div class="flex-1 relative bg-gray-50 rounded-2xl border border-gray-200 focus-within:border-blue-500 focus-within:ring-4 focus-within:ring-blue-500/10 focus-within:bg-white transition-all overflow-hidden flex shadow-inner">
              <textarea 
                class="w-full bg-transparent p-4 outline-none resize-none min-h-[60px] max-h-[160px] text-gray-800 placeholder-gray-400" 
                placeholder="Hoặc nhập câu trả lời bằng văn bản tại đây..."
                rows="2"
                v-model="textAnswer"
                :disabled="recording || analyzing"
              ></textarea>
              <div class="absolute bottom-3 right-3">
                <button 
                  @click="handleSendText"
                  :disabled="!textAnswer || recording || analyzing"
                  class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white p-2 rounded-xl transition-colors shadow-sm"
                  title="Gửi câu trả lời"
                >
                  <Send class="w-5 h-5" />
                </button>
              </div>
            </div>

          </div>
          
          <!-- Recording status indicator -->
          <div v-if="recording" class="max-w-4xl mx-auto mt-4 pl-[88px] flex items-center gap-3">
            <span class="relative flex h-3 w-3">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-rose-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-3 w-3 bg-rose-500"></span>
            </span>
            <span class="text-sm font-bold text-rose-600">Đang thu âm... (Tối đa 3 phút)</span>
            
            <!-- Fake Waveform -->
            <div class="flex gap-1 ml-4 h-6 items-center">
              <div v-for="(h, i) in [2, 5, 3, 7, 4, 8, 3, 6, 2, 5]" :key="i" 
                   class="w-1 bg-rose-400 rounded-full" 
                   :style="{ height: `${h * 10}%`, animation: `pulse-height ${0.3 + (i%4)*0.1}s ease-in-out infinite alternate` }"></div>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-up {
  animation: fade-in-up 0.4s ease-out forwards;
}

@keyframes fade-in-down {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-down {
  animation: fade-in-down 0.4s ease-out forwards;
}

@keyframes pulse-height {
  from { transform: scaleY(0.3); }
  to { transform: scaleY(1); }
}
</style>
