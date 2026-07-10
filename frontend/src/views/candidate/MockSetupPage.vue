<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { Play, FileText, Sparkles } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'

const router = useRouter()
const loading = ref(false)
const toast = ref(null)

const setup = ref({
  jobRole: 'frontend',
  level: 'junior',
  type: 'tech',
  style: 'friendly',
  useCurrentCv: false
})

const currentCvId = ref(null)
const currentCvName = ref('')

import { onMounted } from 'vue'
import { candidatePortalService } from '../../services/candidate-portal.service'

onMounted(async () => {
  try {
    const profile = await candidatePortalService.getProfile()
    if (profile && profile.cv_file_id) {
      if (typeof profile.cv_file_id === 'object' && profile.cv_file_id.Valid) {
        currentCvId.value = profile.cv_file_id.String
      } else if (typeof profile.cv_file_id === 'string') {
        currentCvId.value = profile.cv_file_id
      }
      currentCvName.value = profile.cv_name || ''
    }
  } catch (error) {
    console.error('Không thể lấy thông tin CV:', error)
  }
})

const handleStart = async (e) => {
  e.preventDefault()
  loading.value = true
  try {
    const payload = {
      target_role: setup.value.jobRole,
      target_level: setup.value.level
    }
    
    if (setup.value.useCurrentCv && currentCvId.value) {
      payload.cv_file_id = currentCvId.value
    } else if (setup.value.useCurrentCv && !currentCvId.value) {
      toast.value = { type: 'error', message: 'Bạn chưa có CV nào trong hồ sơ!' }
      loading.value = false
      return
    }
    
    const session = await mockService.createMockInterview(payload)
    // Bước 2: Bắt đầu session → nhận câu hỏi đầu tiên — API_SPEC §11.2
    await mockService.startMockInterview(session.id)
    // Chuyển vào phòng phỏng vấn mock
    router.push({ path: '/mock-room', query: { mock_id: session.id } })
  } catch (error) {
    console.error(error)
    toast.value = { type: 'error', message: 'Không thể khởi động phỏng vấn. Backend đang được kết nối.' }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="space-y-8 animate-fade-in pb-12 max-w-5xl mx-auto">
    <Toast 
      v-if="toast" 
      :type="toast.type" 
      :message="toast.message" 
      @close="toast = null" 
    />
    <div class="mb-8">
      <h1 class="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600">Luyện phỏng vấn cùng AI</h1>
      <p class="text-gray-500 mt-2 text-lg">AI sẽ đóng vai người phỏng vấn thật, đặt câu hỏi theo vị trí ứng tuyển và đưa feedback chi tiết.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Form Setup -->
      <div class="lg:col-span-2 space-y-6">
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl p-8 relative overflow-hidden">
          <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-blue-100 to-transparent rounded-bl-full opacity-50"></div>
          
          <h2 class="text-2xl font-bold text-gray-800 mb-8 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            Cấu hình buổi phỏng vấn
          </h2>
          
          <form @submit="handleStart" class="space-y-6 relative z-10">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-gray-700">Vị trí ứng tuyển (Target Role)</label>
                <select class="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none transition-all bg-gray-50 focus:bg-white text-gray-800 font-medium" required v-model="setup.jobRole">
                  <option value="frontend">Frontend Developer</option>
                  <option value="backend">Backend Developer</option>
                  <option value="fullstack">Fullstack Developer</option>
                  <option value="pm">Product Manager</option>
                </select>
              </div>

              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-gray-700">Cấp độ (Level)</label>
                <select class="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none transition-all bg-gray-50 focus:bg-white text-gray-800 font-medium" required v-model="setup.level">
                  <option value="fresher">Fresher</option>
                  <option value="junior">Junior</option>
                  <option value="middle">Middle</option>
                  <option value="senior">Senior</option>
                </select>
              </div>

              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-gray-700">Loại phỏng vấn</label>
                <select class="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 outline-none transition-all bg-gray-50 focus:bg-white text-gray-800 font-medium" required v-model="setup.type">
                  <option value="tech">Phỏng vấn Kỹ thuật (Technical)</option>
                  <option value="behavior">Phỏng vấn Hành vi (Behavioral)</option>
                  <option value="hr">Phỏng vấn Nhân sự (HR)</option>
                </select>
              </div>

              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-gray-700">Phong cách AI (Interviewer Style)</label>
                <select class="w-full px-4 py-3 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 outline-none transition-all bg-gray-50 focus:bg-white text-gray-800 font-medium" required v-model="setup.style">
                  <option value="friendly">Thân thiện, gợi mở</option>
                  <option value="professional">Chuyên nghiệp, tiêu chuẩn</option>
                  <option value="challenging">Khó tính, hay hỏi xoáy</option>
                </select>
              </div>
            </div>

            <div class="mt-8 pt-8 border-t border-gray-100 flex flex-col md:flex-row justify-between items-center gap-6">
              <div class="flex flex-col gap-1 w-full md:w-auto">
                <label class="flex items-center gap-3 cursor-pointer group bg-blue-50 hover:bg-blue-100 px-4 py-3 rounded-xl transition-colors border border-blue-100">
                  <input type="checkbox" v-model="setup.useCurrentCv" class="w-5 h-5 text-blue-600 rounded border-gray-300 focus:ring-blue-500" />
                  <div class="flex items-center gap-2 text-blue-800 font-medium">
                    <FileText class="w-5 h-5 text-blue-500" /> Sử dụng CV hiện tại trong Hồ sơ
                  </div>
                </label>
                <p v-if="setup.useCurrentCv && currentCvName" class="text-sm text-gray-500 italic pl-4 border-l-2 border-blue-300 ml-2 mt-1">
                  Đang dùng: <span class="font-medium text-gray-700">{{ currentCvName }}</span>
                </p>
              </div>
              
              <button type="submit" :disabled="loading" class="w-full md:w-auto bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold py-3.5 px-8 rounded-xl shadow-lg hover:shadow-indigo-500/30 transition-all duration-300 transform hover:-translate-y-0.5 flex items-center justify-center gap-2 disabled:opacity-50 disabled:transform-none">
                <Play class="w-5 h-5 fill-current" /> 
                {{ loading ? 'Đang khởi tạo phòng...' : 'Bắt đầu ngay' }}
              </button>
            </div>
          </form>
        </Card>
      </div>

      <!-- Right Column -->
      <div class="lg:col-span-1 space-y-6">
        <!-- AI Interviewer Card -->
        <Card class="bg-gradient-to-br from-slate-800 to-slate-900 border-0 rounded-2xl overflow-hidden relative shadow-2xl">
          <!-- Animated glowing background -->
          <div class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-indigo-400 via-transparent to-transparent animate-pulse" style="animation-duration: 3s;"></div>
          
          <div class="p-8 text-center relative z-10">
            <div class="w-24 h-24 mx-auto bg-gradient-to-tr from-indigo-500 to-purple-500 rounded-full flex items-center justify-center shadow-[0_0_40px_rgba(99,102,241,0.5)] mb-6 border-4 border-slate-700">
              <Sparkles class="w-12 h-12 text-white animate-pulse" />
            </div>
            <h3 class="text-2xl font-bold text-white mb-2">AI Interviewer</h3>
            <p class="text-indigo-200 font-medium mb-6">Sẵn sàng hỗ trợ bạn</p>
            
            <div class="inline-flex items-center gap-2 bg-slate-800 px-4 py-2 rounded-full border border-slate-700">
              <span class="relative flex h-3 w-3">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
              </span>
              <span class="text-sm font-semibold text-emerald-400 uppercase tracking-wider">Trực tuyến</span>
            </div>
          </div>
        </Card>
        
        <!-- Checklist -->
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl p-6">
          <h3 class="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            Lưu ý trước khi bắt đầu
          </h3>
          <ul class="space-y-4">
            <li class="flex items-start gap-3 text-gray-600 font-medium text-sm">
              <div class="mt-0.5 bg-emerald-100 p-1 rounded-full text-emerald-600">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
              </div>
              Chuẩn bị micro và ngồi trong không gian yên tĩnh.
            </li>
            <li class="flex items-start gap-3 text-gray-600 font-medium text-sm">
              <div class="mt-0.5 bg-emerald-100 p-1 rounded-full text-emerald-600">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
              </div>
              Sẽ có khoảng 3-5 câu hỏi tùy thuộc vào chức danh.
            </li>
            <li class="flex items-start gap-3 text-gray-600 font-medium text-sm">
              <div class="mt-0.5 bg-emerald-100 p-1 rounded-full text-emerald-600">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
              </div>
              Bạn có thể trả lời bằng Giọng nói (khuyên dùng) hoặc Văn bản.
            </li>
          </ul>
        </Card>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in {
  animation: fade-in 0.5s ease-out forwards;
}
</style>
