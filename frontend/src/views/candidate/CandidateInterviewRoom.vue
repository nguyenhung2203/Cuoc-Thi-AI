<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from '../../components/common/AppButton.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mic, MicOff, Video, VideoOff, MessageSquare, PhoneOff, CheckCircle, FileText } from 'lucide-vue-next'
import { useLiveKit } from '../../composables/useLiveKit'

const router = useRouter()
const route = useRoute()

const {
  isConnected, error, isMicOn, isCameraOn,
  localVideoEl, remoteVideoEl,
  connectToRoom, toggleMic, toggleCamera, disconnect
} = useLiveKit()

const activeTab = ref('info')

const showLeaveModal = ref(false)
const isLeaving = ref(false)

const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)
const interviewInfo = ref(history.state?.interviewInfo || null)

onMounted(async () => {
  if (history.state?.message) {
    window.history.replaceState({ interviewInfo: history.state.interviewInfo }, document.title)
  }
  
  // Logic kết nối LiveKit
  const token = route.query.token
  if (token) {
    try {
      const livekitUrl = import.meta.env.VITE_LIVEKIT_URL || 'ws://localhost:7880'
      await connectToRoom(livekitUrl, token)
    } catch (err) {
      console.error('Không thể vào phòng', err)
    }
  }
})

onUnmounted(() => {
  disconnect()
})

const handleLeave = () => {
  showLeaveModal.value = true
}

const confirmLeave = () => {
  showLeaveModal.value = false
  isLeaving.value = true
  setTimeout(() => {
    router.push({ path: '/home', state: { message: 'Rời phòng phỏng vấn thành công' } })
  }, 1500)
}
</script>

<template>
  <div class="h-screen flex flex-col bg-slate-900 font-sans text-gray-100 overflow-hidden relative">
    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />
    <Toast v-if="isLeaving" type="info" message="Đang rời phòng phỏng vấn..." :duration="1500" />

    <!-- Top Header -->
    <header class="h-16 bg-slate-800/80 backdrop-blur-md border-b border-slate-700/50 flex items-center justify-between px-6 z-20 shadow-sm shrink-0">
      <div class="flex items-center gap-4">
        <div>
          <h1 class="text-lg font-bold text-white truncate max-w-[300px] md:max-w-[500px]">
            Phỏng vấn: {{ interviewInfo?.job_title || 'Vị trí ứng tuyển' }}
          </h1>
          <p class="text-sm text-slate-400 font-medium">Công ty {{ interviewInfo?.company_name || 'Tuyển dụng' }}</p>
        </div>
      </div>
      <div class="flex items-center gap-4 bg-slate-900/50 px-4 py-1.5 rounded-full border border-slate-700">
        <span class="text-sm font-medium text-emerald-400 flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          <span class="hidden sm:inline">Đường truyền tốt</span>
        </span>
        <span class="w-px h-4 bg-slate-700"></span>
        <span class="text-sm font-mono font-bold text-white tracking-wider">00:15:32</span>
      </div>
    </header>

    <!-- Main Workspace -->
    <div class="flex flex-1 overflow-hidden relative z-10">
      
      <!-- Left: Video Area -->
      <div class="flex-1 flex flex-col p-4 md:p-6 gap-4 relative transition-all duration-300">
        
        <!-- Notification Banner -->
        <div class="bg-indigo-500/10 border border-indigo-500/30 backdrop-blur-md rounded-xl p-3 flex items-center gap-3 text-indigo-200 text-sm shadow-sm shrink-0">
          <div class="w-8 h-8 rounded-lg bg-indigo-500/20 flex items-center justify-center shrink-0">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-400" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" /></svg>
          </div>
          <p><strong class="text-indigo-300">Lời nhắc:</strong> Hãy trả lời tự nhiên. Nhà tuyển dụng sẽ dẫn dắt buổi phỏng vấn.</p>
        </div>

        <!-- Video Grid -->
        <div class="flex-1 relative rounded-2xl overflow-hidden bg-slate-950 shadow-2xl border border-slate-800/60 ring-1 ring-white/5 flex">
          
          <!-- Recruiter Video (Main View) -->
          <div class="absolute inset-0 flex items-center justify-center">
            <!-- Thẻ video thật cho LiveKit (Remote - Recruiter) -->
            <video ref="remoteVideoEl" autoplay playsinline class="w-full h-full object-cover"></video>
            
            <!-- Fallback if no video -->
            <div class="absolute inset-0 flex flex-col items-center justify-center bg-gradient-to-br from-slate-900 to-slate-950 z-0">
              <div class="w-24 h-24 rounded-full bg-slate-800 border-2 border-slate-700 flex items-center justify-center mb-4 text-4xl font-bold text-slate-500 shadow-inner">
                R
              </div>
              <p class="text-slate-400 font-medium animate-pulse">Đang chờ nhà tuyển dụng kết nối camera...</p>
            </div>
            
            <!-- Recruiter Name Badge -->
            <div class="absolute bottom-6 left-6 bg-black/60 backdrop-blur-md text-white px-4 py-2 rounded-xl text-sm font-medium border border-white/10 z-10 shadow-lg flex items-center gap-2">
              <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
              Nhà Tuyển Dụng
            </div>
          </div>
          
          <!-- Candidate Self View (PiP) -->
          <div class="absolute top-6 right-6 w-48 md:w-64 aspect-video bg-slate-800 rounded-xl overflow-hidden border-2 border-slate-700/80 shadow-2xl z-20 group hover:scale-105 transition-transform duration-300 cursor-move">
             <!-- Thẻ video thật cho LiveKit (Local) -->
             <video ref="localVideoEl" autoplay playsinline muted class="w-full h-full object-cover transform -scale-x-100"></video>
             
             <!-- Self Name Badge -->
             <div class="absolute bottom-2 left-2 bg-black/50 backdrop-blur-sm text-white px-2 py-1 rounded-lg text-xs font-medium opacity-0 group-hover:opacity-100 transition-opacity">
               Bạn
             </div>
             <!-- Mic Status icon on self view -->
             <div class="absolute top-2 right-2 w-6 h-6 rounded-full bg-black/50 backdrop-blur-sm flex items-center justify-center">
               <Mic v-if="isMicOn" class="w-3 h-3 text-emerald-400" />
               <MicOff v-else class="w-3 h-3 text-rose-500" />
             </div>
          </div>
        </div>

        <!-- Floating Control Bar -->
        <div class="absolute bottom-8 left-1/2 -translate-x-1/2 bg-slate-800/90 backdrop-blur-xl border border-slate-600/50 p-2 rounded-2xl shadow-2xl flex items-center gap-2 z-30 transition-transform hover:-translate-y-1 duration-300">
          <button @click="toggleMic" class="w-12 h-12 rounded-xl flex items-center justify-center transition-all duration-200" :class="isMicOn ? 'bg-slate-700 hover:bg-slate-600 text-white' : 'bg-rose-500/20 hover:bg-rose-500/30 text-rose-500 border border-rose-500/50'">
            <Mic v-if="isMicOn" class="w-5 h-5" />
            <MicOff v-else class="w-5 h-5" />
          </button>
          
          <button @click="toggleCamera" class="w-12 h-12 rounded-xl flex items-center justify-center transition-all duration-200" :class="isCameraOn ? 'bg-slate-700 hover:bg-slate-600 text-white' : 'bg-rose-500/20 hover:bg-rose-500/30 text-rose-500 border border-rose-500/50'">
            <Video v-if="isCameraOn" class="w-5 h-5" />
            <VideoOff v-else class="w-5 h-5" />
          </button>
          
          <div class="w-px h-8 bg-slate-700 mx-2"></div>
          
          <button @click="handleLeave" class="h-12 px-6 rounded-xl bg-rose-600 hover:bg-rose-500 text-white font-bold shadow-lg shadow-rose-600/20 flex items-center gap-2 transition-colors">
            <PhoneOff class="w-5 h-5" /> <span class="hidden sm:inline">Rời phòng</span>
          </button>
        </div>
      </div>

      <!-- Right: Side Panel (Tabs: JD / Chat) -->
      <div class="w-80 md:w-96 bg-slate-800/50 border-l border-slate-700/50 flex flex-col shrink-0 backdrop-blur-sm relative z-20">
        <!-- Tabs Header -->
        <div class="flex border-b border-slate-700/50 shrink-0">
          <button @click="activeTab = 'info'" class="flex-1 py-4 flex items-center justify-center gap-2 font-medium transition-colors border-b-2" :class="activeTab === 'info' ? 'text-indigo-400 border-indigo-400 bg-indigo-500/5' : 'text-slate-400 border-transparent hover:bg-slate-800 hover:text-slate-200'">
            <FileText class="w-4 h-4" /> JD & Ghi chú
          </button>
          <button @click="activeTab = 'chat'" class="flex-1 py-4 flex items-center justify-center gap-2 font-medium transition-colors border-b-2 relative" :class="activeTab === 'chat' ? 'text-indigo-400 border-indigo-400 bg-indigo-500/5' : 'text-slate-400 border-transparent hover:bg-slate-800 hover:text-slate-200'">
            <MessageSquare class="w-4 h-4" /> Chat
            <span class="absolute top-3 right-8 w-2 h-2 rounded-full bg-rose-500 border-2 border-slate-800"></span>
          </button>
        </div>

        <!-- Tab Content Area -->
        <div class="flex-1 overflow-y-auto custom-scrollbar relative">
          
          <!-- Tab: Info & Notes -->
          <div v-show="activeTab === 'info'" class="absolute inset-0 p-6 flex flex-col">
            <div class="mb-6">
              <h3 class="text-sm font-bold text-slate-300 uppercase tracking-wider mb-3 flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-indigo-400" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M6 2a2 2 0 00-2 2v12a2 2 0 002 2h8a2 2 0 002-2V7.414A2 2 0 0015.414 6L12 2.586A2 2 0 0010.586 2H6zm5 6a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V8z" clip-rule="evenodd" /></svg>
                Yêu cầu công việc
              </h3>
              <ul class="space-y-3">
                <li class="flex items-start gap-3 bg-slate-800/50 p-3 rounded-xl border border-slate-700/50">
                  <span class="w-1.5 h-1.5 rounded-full bg-indigo-400 mt-2 shrink-0"></span>
                  <span class="text-sm text-slate-300 leading-relaxed">Phát triển các tính năng Frontend sử dụng ReactJS.</span>
                </li>
                <li class="flex items-start gap-3 bg-slate-800/50 p-3 rounded-xl border border-slate-700/50">
                  <span class="w-1.5 h-1.5 rounded-full bg-indigo-400 mt-2 shrink-0"></span>
                  <span class="text-sm text-slate-300 leading-relaxed">Tối ưu hóa hiệu suất ứng dụng web.</span>
                </li>
                <li class="flex items-start gap-3 bg-slate-800/50 p-3 rounded-xl border border-slate-700/50">
                  <span class="w-1.5 h-1.5 rounded-full bg-indigo-400 mt-2 shrink-0"></span>
                  <span class="text-sm text-slate-300 leading-relaxed">Phối hợp với UI/UX designer và Backend developer.</span>
                </li>
              </ul>
            </div>
            
            <div class="flex-1 flex flex-col">
              <h3 class="text-sm font-bold text-slate-300 uppercase tracking-wider mb-3 flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-400" viewBox="0 0 20 20" fill="currentColor"><path d="M17.414 2.586a2 2 0 00-2.828 0L7 10.172V13h2.828l7.586-7.586a2 2 0 000-2.828z" /><path fill-rule="evenodd" d="M2 6a2 2 0 012-2h4a1 1 0 010 2H4v10h10v-4a1 1 0 112 0v4a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" clip-rule="evenodd" /></svg>
                Ghi chú của bạn
              </h3>
              <textarea 
                class="flex-1 w-full bg-slate-900/50 border border-slate-700 rounded-xl p-4 text-slate-300 placeholder-slate-600 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none resize-none transition-all" 
                placeholder="Viết nháp câu trả lời hoặc note nhanh tại đây (Nhà tuyển dụng không thấy)..."
              ></textarea>
            </div>
          </div>

          <!-- Tab: Chat -->
          <div v-show="activeTab === 'chat'" class="absolute inset-0 flex flex-col bg-slate-800/30">
            <div class="flex-1 overflow-y-auto p-4 flex flex-col gap-4">
              <!-- Example message -->
              <div class="flex flex-col gap-1 items-start">
                <span class="text-xs text-slate-500 font-medium ml-1">Nhà Tuyển Dụng - 10:05</span>
                <div class="bg-slate-700 text-slate-200 px-4 py-2.5 rounded-2xl rounded-tl-none text-sm max-w-[85%] shadow-sm">
                  Chào bạn, bạn nghe rõ không ạ?
                </div>
              </div>
            </div>
            
            <div class="p-4 bg-slate-800/80 border-t border-slate-700/50 shrink-0">
              <div class="flex items-center gap-2 bg-slate-900/50 border border-slate-700 rounded-xl p-1 focus-within:border-indigo-500 focus-within:ring-1 focus-within:ring-indigo-500 transition-all">
                <input type="text" class="flex-1 bg-transparent border-none text-sm text-white px-3 outline-none placeholder-slate-500" placeholder="Nhập tin nhắn..." />
                <button class="w-8 h-8 flex items-center justify-center bg-indigo-600 text-white rounded-lg hover:bg-indigo-500 transition-colors shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transform rotate-90" viewBox="0 0 20 20" fill="currentColor"><path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z" /></svg>
                </button>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>

    <!-- Leave Modal -->
    <Modal :isOpen="showLeaveModal" @close="showLeaveModal = false" title="Xác nhận rời phòng">
      <div class="p-2">
        <p class="text-gray-700 mb-6 font-medium">Bạn có chắc chắn muốn kết thúc buổi phỏng vấn và rời phòng?</p>
        <div class="flex justify-end gap-3">
          <button @click="showLeaveModal = false" class="px-5 py-2.5 rounded-xl border border-gray-300 text-gray-700 font-bold hover:bg-gray-50 transition-colors">
            Tiếp tục ở lại
          </button>
          <button @click="confirmLeave" class="px-5 py-2.5 rounded-xl bg-rose-600 text-white font-bold hover:bg-rose-500 shadow-md shadow-rose-600/20 transition-colors flex items-center gap-2">
            <PhoneOff class="w-4 h-4" /> Rời khỏi phòng
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
