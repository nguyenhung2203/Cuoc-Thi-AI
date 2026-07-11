<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { AlertTriangle, Video, Mic, ShieldCheck } from 'lucide-vue-next'
import { interviewService } from '../../services/interview.service'

const router = useRouter()
const route = useRoute()
const agreed = ref(false)
const loading = ref(true)
const errorMsg = ref('')
const interviewInfo = ref(null)
const inviteToken = ref(route.query.token || '')

onMounted(async () => {
  if (!inviteToken.value) {
    errorMsg.value = 'Link mời phỏng vấn không hợp lệ hoặc bị thiếu token.'
    loading.value = false
    return
  }

  try {
    const data = await interviewService.joinByToken(inviteToken.value)
    interviewInfo.value = data
  } catch (error) {
    errorMsg.value = error.message || 'Không thể xác thực link mời phỏng vấn. Link có thể đã hết hạn.'
  } finally {
    loading.value = false
  }
})

const handleJoin = () => {
  if (!agreed.value) return
  // Pass token and details to candidate room
  router.push({ 
    path: '/candidate-room', 
    query: { token: interviewInfo.value.room_access_token },
    state: { message: 'Vào phòng phỏng vấn thành công!', interviewInfo: interviewInfo.value } 
  })
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-6 relative overflow-hidden font-sans">
    <!-- Relaxing background animations -->
    <div class="absolute inset-0 z-0 overflow-hidden pointer-events-none">
      <!-- Clean background, no blobs as per design system -->
    </div>

    <div class="max-w-2xl w-full relative z-10 animate-fade-in-up">
      <div class="text-center mb-10">
        <div class="w-20 h-20 bg-[var(--primary)] text-white rounded-2xl flex items-center justify-center mx-auto mb-5 shadow-lg transform hover:scale-105 transition-transform">
          <ShieldCheck class="w-10 h-10" />
        </div>
        <h1 class="text-h1 font-extrabold tracking-tight">Chuẩn bị vào phòng phỏng vấn</h1>
        
        <div v-if="loading" class="mt-4 flex flex-col items-center justify-center gap-3">
          <div class="w-8 h-8 border-4 border-blue-200 border-t-[var(--primary)] rounded-full animate-spin"></div>
          <p class="text-gray-500 font-medium animate-pulse">Đang xác thực thông tin phòng...</p>
        </div>
        <div v-else-if="errorMsg" class="mt-4 text-rose-600 font-medium bg-rose-50 inline-block px-4 py-2 rounded-lg border border-rose-100">
          {{ errorMsg }}
        </div>
        <div v-else class="mt-4 inline-flex items-center gap-2 bg-white px-5 py-2 rounded-full shadow-sm border border-gray-200">
          <span class="w-2 h-2 rounded-full bg-[var(--success)] animate-pulse"></span>
          <p class="text-gray-700 font-bold">{{ interviewInfo?.job_title }} <span class="text-gray-400 mx-1">|</span> {{ interviewInfo?.company_name }}</p>
        </div>
      </div>

      <Card v-if="!loading && !errorMsg" class="shadow-lg rounded-2xl p-8 md:p-10">
        <!-- AI & Privacy Notice -->
        <div class="bg-amber-50/80 border border-amber-200/60 rounded-2xl p-6 mb-8 flex flex-col md:flex-row gap-5 hover:bg-amber-50 transition-colors">
          <div class="w-12 h-12 rounded-full bg-amber-100 flex items-center justify-center shrink-0">
            <AlertTriangle class="w-6 h-6 text-amber-600" />
          </div>
          <div>
            <h3 class="text-lg font-bold text-amber-800 mb-2">Thông báo về Quyền riêng tư & AI</h3>
            <p class="text-amber-700/90 text-sm leading-relaxed mb-3 font-medium">
              Buổi phỏng vấn này sẽ có sự hỗ trợ của Trí tuệ nhân tạo (AI). Nhằm mục đích đánh giá công bằng và cung cấp báo cáo chi tiết cho Nhà tuyển dụng, toàn bộ diễn biến bao gồm:
            </p>
            <ul class="space-y-2">
              <li class="flex items-start gap-2 text-amber-800/80 text-sm font-medium">
                <span class="w-1.5 h-1.5 rounded-full bg-amber-400 mt-1.5 shrink-0"></span>
                Dữ liệu âm thanh (Giọng nói) sẽ được ghi lại và chuyển ngữ thành văn bản.
              </li>
              <li class="flex items-start gap-2 text-amber-800/80 text-sm font-medium">
                <span class="w-1.5 h-1.5 rounded-full bg-amber-400 mt-1.5 shrink-0"></span>
                Dữ liệu văn bản sẽ được AI phân tích tự động.
              </li>
              <li class="flex items-start gap-2 text-amber-800/80 text-sm font-medium">
                <span class="w-1.5 h-1.5 rounded-full bg-amber-400 mt-1.5 shrink-0"></span>
                Hình ảnh từ Camera sẽ được truyền trực tiếp nhưng KHÔNG lưu trữ.
              </li>
            </ul>
          </div>
        </div>

        <!-- System Checks -->
        <div class="flex items-center gap-4 p-5 bg-slate-50 border border-slate-200 rounded-2xl mb-8">
          <div class="flex gap-3 text-[var(--primary)] bg-white p-2 rounded-xl shadow-sm border border-[var(--primary-light)]">
            <div class="relative">
              <Video class="w-6 h-6" />
              <div class="absolute -top-1 -right-1 w-2.5 h-2.5 bg-[var(--success)] border-2 border-white rounded-full"></div>
            </div>
            <div class="w-px h-6 bg-[var(--primary-light)]"></div>
            <div class="relative">
              <Mic class="w-6 h-6" />
              <div class="absolute -top-1 -right-1 w-2.5 h-2.5 bg-[var(--success)] border-2 border-white rounded-full animate-pulse"></div>
            </div>
          </div>
          <p class="text-gray-600 text-sm font-medium leading-relaxed flex-1">
            Hệ thống sẽ yêu cầu quyền truy cập Camera và Micro ở bước tiếp theo để tiến hành phỏng vấn.
          </p>
        </div>

        <!-- Agreement -->
        <label class="flex items-start gap-4 cursor-pointer mb-10 group">
          <div class="relative flex items-center justify-center shrink-0 mt-0.5">
            <input 
              type="checkbox" 
              v-model="agreed"
              class="peer appearance-none w-6 h-6 border-2 border-gray-300 rounded-lg checked:bg-[var(--primary)] checked:border-[var(--primary)] transition-colors cursor-pointer focus:ring-4 focus:ring-[var(--accent)]/20 outline-none" 
            />
            <svg class="absolute w-4 h-4 text-white opacity-0 peer-checked:opacity-100 transition-opacity pointer-events-none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </div>
          <span class="text-gray-700 font-medium leading-relaxed group-hover:text-gray-900 transition-colors">
            Tôi đã đọc, hiểu rõ và đồng ý với việc sử dụng hệ thống AI phân tích và ghi âm trong buổi phỏng vấn này.
          </span>
        </label>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row gap-4">
          <button @click="router.push('/home')" class="flex-1 py-3.5 px-6 bg-white border-2 border-gray-200 text-gray-700 font-bold rounded-xl hover:bg-gray-50 hover:border-gray-300 transition-colors shadow-sm focus:ring-4 focus:ring-gray-100 outline-none">
            Từ chối & Quay lại
          </button>
          <button @click="handleJoin" :disabled="!agreed" class="flex-1 py-3.5 px-6 bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white font-bold rounded-xl shadow-md transition-all duration-300 transform hover:-translate-y-0.5 disabled:opacity-50 disabled:transform-none disabled:shadow-none outline-none">
            Tham gia phỏng vấn
          </button>
        </div>
      </Card>
      
      <Card v-if="errorMsg" class="shadow-lg rounded-2xl p-10 text-center animate-fade-in-up">
        <div class="w-16 h-16 bg-rose-100 text-rose-500 rounded-full flex items-center justify-center mx-auto mb-6">
          <AlertTriangle class="w-8 h-8" />
        </div>
        <p class="text-gray-800 font-bold text-lg mb-8">{{ errorMsg }}</p>
        <button @click="router.push('/home')" class="py-3 px-8 bg-gray-900 text-white font-bold rounded-xl hover:bg-gray-800 transition-colors shadow-lg shadow-gray-900/20">
          Quay về trang chủ
        </button>
      </Card>
    </div>
  </div>
</template>

<style scoped>
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-up {
  animation: fade-in-up 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}
.animate-blob {
  animation: blob 10s infinite;
}
.animation-delay-2000 {
  animation-delay: 2s;
}
.animation-delay-4000 {
  animation-delay: 4s;
}
</style>
