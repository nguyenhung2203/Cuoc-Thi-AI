<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import WelcomeAlert from '../../components/common/WelcomeAlert.vue'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { authStore } from '../../stores/auth.store'

const router = useRouter()
const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)

const stats = ref({
  upcoming_interviews: 0,
  completed_mock_tests: 0,
  average_mock_score: 0,
  profile_completeness: 0
})
const upcomingInterviews = ref([])
const loading = ref(true)

onMounted(async () => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }

  if (!authStore.isAuthenticated && !localStorage.getItem('access_token')) {
    stats.value = {
      upcoming_interviews: 0,
      completed_mock_tests: 0,
      average_mock_score: 0,
      profile_completeness: 0
    }
    upcomingInterviews.value = []
    loading.value = false
    return
  }

  try {
    const [statsData, interviewsData] = await Promise.all([
      candidatePortalService.getDashboardStats(),
      candidatePortalService.getInterviews()
    ])
    
    stats.value = statsData
    // Filter only future interviews or recently active ones
    upcomingInterviews.value = interviewsData.filter(i => i.status === 'scheduled' || i.status === 'active').slice(0, 3)
  } catch (err) {
    console.error('Lỗi tải dữ liệu dashboard:', err)
  } finally {
    loading.value = false
  }
})

const formatDate = (dateString) => {
  if (!dateString) return ''
  const d = new Date(dateString)
  return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' }) + ', ' + d.toLocaleDateString('vi-VN')
}
</script>

<template>
  <div class="space-y-6 animate-fade-in pb-10">
    <WelcomeAlert 
      v-if="entryToast" 
      role="candidate"
      title="Thành công!"
      :message="entryToast.message" 
      @close="entryToast = null" 
    />

    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div>
        <h1 class="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600">
          Chào mừng trở lại, {{ authStore.user?.full_name || 'Ứng viên' }} 👋
        </h1>
        <p class="text-gray-500 mt-2 text-base">
          Theo dõi lịch phỏng vấn, luyện tập với AI và cải thiện kỹ năng trả lời của bạn.
        </p>
      </div>
    </div>
    
    <!-- Loading State -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
      <p class="text-gray-500 font-medium animate-pulse">Đang tải dữ liệu dashboard...</p>
    </div>
    
    <template v-else>
      <!-- Stats Overview Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Card 1 -->
        <Card class="relative overflow-hidden group bg-white/80 backdrop-blur-xl border border-white/20 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 rounded-2xl">
          <div class="absolute inset-0 bg-gradient-to-br from-blue-50 to-transparent opacity-50"></div>
          <div class="relative p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="p-3 bg-blue-100/50 text-blue-600 rounded-xl group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
              </div>
            </div>
            <h3 class="text-gray-500 text-sm font-medium">Lịch sắp tới</h3>
            <p class="text-4xl font-extrabold text-gray-900 mt-1">{{ stats.upcoming_interviews }}</p>
          </div>
        </Card>

        <!-- Card 2 -->
        <Card class="relative overflow-hidden group bg-white/80 backdrop-blur-xl border border-white/20 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 rounded-2xl">
          <div class="absolute inset-0 bg-gradient-to-br from-indigo-50 to-transparent opacity-50"></div>
          <div class="relative p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="p-3 bg-indigo-100/50 text-indigo-600 rounded-xl group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              </div>
            </div>
            <h3 class="text-gray-500 text-sm font-medium">Luyện tập AI đã xong</h3>
            <p class="text-4xl font-extrabold text-gray-900 mt-1">{{ stats.completed_mock_tests }}</p>
          </div>
        </Card>

        <!-- Card 3 -->
        <Card class="relative overflow-hidden group bg-white/80 backdrop-blur-xl border border-white/20 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 rounded-2xl">
          <div class="absolute inset-0 bg-gradient-to-br from-amber-50 to-transparent opacity-50"></div>
          <div class="relative p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="p-3 bg-amber-100/50 text-amber-600 rounded-xl group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" /></svg>
              </div>
            </div>
            <h3 class="text-gray-500 text-sm font-medium">Điểm AI trung bình</h3>
            <p class="text-4xl font-extrabold text-gray-900 mt-1">{{ stats.average_mock_score.toFixed(1) }}</p>
          </div>
        </Card>

        <!-- Card 4 -->
        <Card class="relative overflow-hidden group bg-white/80 backdrop-blur-xl border border-white/20 shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 rounded-2xl">
          <div class="absolute inset-0 bg-gradient-to-br from-emerald-50 to-transparent opacity-50"></div>
          <div class="relative p-6">
            <div class="flex items-center justify-between mb-4">
              <div class="p-3 bg-emerald-100/50 text-emerald-600 rounded-xl group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
              </div>
            </div>
            <h3 class="text-gray-500 text-sm font-medium">Mức độ hoàn thiện CV</h3>
            <div class="flex items-baseline gap-2 mt-1">
              <p class="text-4xl font-extrabold text-gray-900">{{ stats.profile_completeness }}%</p>
            </div>
            <!-- Mini progress bar -->
            <div class="w-full bg-gray-200 rounded-full h-1.5 mt-4 overflow-hidden">
              <div class="bg-gradient-to-r from-emerald-400 to-emerald-500 h-1.5 rounded-full transition-all duration-1000 ease-out" :style="`width: ${stats.profile_completeness}%`"></div>
            </div>
          </div>
        </Card>
      </div>
      
      <div class="grid grid-cols-1 lg:grid-cols-3 mt-8 gap-8">
        <!-- Lịch phỏng vấn sắp tới -->
        <div class="lg:col-span-2 space-y-8">
          <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl p-6">
            <div class="flex justify-between items-center mb-6">
              <h3 class="text-xl font-bold text-gray-800 flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                Lịch phỏng vấn sắp tới
              </h3>
              <button class="text-blue-600 hover:text-blue-800 font-semibold text-sm transition-colors" @click="router.push('/my-interviews')">
                Xem tất cả &rarr;
              </button>
            </div>
            
            <div v-if="upcomingInterviews.length === 0" class="py-12 flex flex-col items-center justify-center border-2 border-dashed border-gray-200 rounded-xl bg-gray-50/50">
              <div class="w-16 h-16 bg-blue-50 text-blue-300 rounded-full flex items-center justify-center mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" /></svg>
              </div>
              <p class="text-gray-500 font-medium">Bạn chưa có lịch phỏng vấn nào sắp tới.</p>
              <Button variant="outline" class="mt-4" @click="router.push('/job-board')">Tìm việc ngay</Button>
            </div>
            
            <div class="space-y-4">
              <div v-for="iv in upcomingInterviews" :key="iv.id" 
                class="group border border-gray-100 hover:border-blue-200 rounded-xl p-5 transition-all duration-300 hover:shadow-md bg-white flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div>
                  <h4 class="font-bold text-gray-900 text-lg group-hover:text-blue-600 transition-colors">{{ iv.job_title || iv.title }}</h4>
                  <p class="text-gray-500 text-sm font-medium mt-1 flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" /></svg>
                    {{ iv.company_name || 'Công ty ẩn danh' }}
                  </p>
                  <div class="flex flex-wrap gap-2 mt-3">
                    <span class="inline-flex items-center px-2.5 py-1 rounded-md text-xs font-semibold bg-blue-50 text-blue-700">
                      {{ formatDate(iv.scheduled_at) }}
                    </span>
                    <span :class="['inline-flex items-center px-2.5 py-1 rounded-md text-xs font-semibold', iv.mode === 'real' ? 'bg-rose-50 text-rose-700' : 'bg-gray-100 text-gray-700']">
                      {{ iv.mode === 'real' ? 'Phỏng vấn thật' : 'Phỏng vấn thử' }}
                    </span>
                  </div>
                </div>
                <div>
                  <Button v-if="iv.mode === 'real'" variant="primary" class="w-full sm:w-auto shadow-md hover:shadow-lg shadow-blue-500/30" @click="iv.join_link ? router.push(iv.join_link) : null">
                    Tham gia ngay
                  </Button>
                </div>
              </div>
            </div>
          </Card>
          
          <!-- Banner: AI Practice -->
          <div class="relative overflow-hidden rounded-2xl bg-gradient-to-r from-sky-400 via-blue-500 to-indigo-600 p-8 text-white shadow-xl">
            <div class="absolute top-0 right-0 -mt-16 -mr-16 text-white/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-64 w-64" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/></svg>
            </div>
            <div class="relative z-10">
              <h3 class="text-2xl font-bold mb-2">Sẵn sàng vượt qua mọi câu hỏi phỏng vấn?</h3>
              <p class="text-blue-100 mb-6 max-w-lg text-lg">Trải nghiệm phỏng vấn 1-kèm-1 với AI Interviewer của chúng tôi. Luyện tập không giới hạn, nhận phản hồi ngay lập tức.</p>
              <button @click="router.push('/mock-setup')" class="bg-white text-blue-600 hover:bg-blue-50 font-bold py-3 px-6 rounded-full shadow-lg transition-transform transform hover:scale-105 flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd" /></svg>
                Bắt đầu luyện tập
              </button>
            </div>
          </div>
        </div>
        
        <!-- Sidebar -->
        <div class="space-y-8">
          <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl p-6">
            <h3 class="text-lg font-bold text-gray-800 mb-5 flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.95 11.95 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" /></svg>
              Hành trang ứng viên
            </h3>
            <div class="space-y-4">
              <div class="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer" @click="router.push('/my-cv')">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                  </div>
                  <div>
                    <span class="block text-sm font-bold text-gray-800">Tải lên CV</span>
                    <span class="block text-xs text-gray-500">Bắt buộc để AI phân tích</span>
                  </div>
                </div>
                <span v-if="stats.profile_completeness > 50" class="text-emerald-500"><svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg></span>
                <span v-else class="text-amber-500"><svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" /></svg></span>
              </div>
              
              <div class="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer" @click="router.push('/profile')">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                  </div>
                  <div>
                    <span class="block text-sm font-bold text-gray-800">Thêm Kỹ năng</span>
                    <span class="block text-xs text-gray-500">Giúp nhà tuyển dụng tìm thấy bạn</span>
                  </div>
                </div>
                <span class="text-amber-500"><svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" /></svg></span>
              </div>
            </div>
          </Card>
          
          <Card class="relative overflow-hidden bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl p-6 border-t-4 border-t-indigo-500">
            <h3 class="text-lg font-bold text-gray-800 mb-3 flex items-center gap-2">
              <span class="animate-pulse">✨</span> AI Career Coach
            </h3>
            <div class="bg-indigo-50 rounded-xl p-4 border border-indigo-100 relative">
              <div class="absolute -left-2 -top-2 w-6 h-6 bg-indigo-500 rounded-full flex items-center justify-center shadow-md">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 text-white" viewBox="0 0 20 20" fill="currentColor"><path d="M11 3a1 1 0 10-2 0v1a1 1 0 102 0V3zM15.657 5.757a1 1 0 00-1.414-1.414l-.707.707a1 1 0 001.414 1.414l.707-.707zM18 10a1 1 0 01-1 1h-1a1 1 0 110-2h1a1 1 0 011 1zM5.05 6.464A1 1 0 106.464 5.05l-.707-.707a1 1 0 00-1.414 1.414l.707.707zM5 10a1 1 0 01-1 1H3a1 1 0 110-2h1a1 1 0 011 1zM8 16v-1h4v1a2 2 0 11-4 0zM12 14c.015-.34.208-.646.477-.859a4 4 0 10-4.954 0c.27.213.462.519.476.859h4.002z" /></svg>
              </div>
              <p class="text-sm text-gray-700 leading-relaxed pl-2">
                Dựa trên kết quả phỏng vấn gần đây, tốc độ nói của bạn rất tốt, tuy nhiên bạn nên luyện tập thêm cách trả lời rành mạch các câu hỏi về <strong class="text-indigo-600 font-bold">Kỹ năng chuyên môn sâu</strong>.
              </p>
            </div>
            <button class="w-full mt-4 py-2.5 px-4 bg-white border-2 border-indigo-100 text-indigo-600 font-semibold rounded-xl hover:bg-indigo-50 transition-colors shadow-sm" @click="router.push('/mock-setup')">
              Luyện chủ đề này
            </button>
          </Card>
        </div>
      </div>
    </template>
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
