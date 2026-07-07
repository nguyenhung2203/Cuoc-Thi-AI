<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { candidatePortalService } from '../../services/candidate-portal.service'
import Badge from '../../components/common/AppBadge.vue'
import Button from '../../components/common/AppButton.vue'

const router = useRouter()
const interviews = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await candidatePortalService.getInterviews()
    interviews.value = data
  } catch (err) {
    console.error('Lỗi tải danh sách phỏng vấn:', err)
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
  <div class="space-y-8 animate-fade-in pb-12 max-w-6xl mx-auto px-4 md:px-0">
    <div class="mb-8 border-b border-gray-100 pb-6">
      <h1 class="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600 tracking-tight">Phỏng vấn của tôi</h1>
      <p class="text-gray-500 mt-2 text-lg font-medium">Quản lý các lịch phỏng vấn sắp tới và lịch sử phỏng vấn.</p>
    </div>
    
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="w-12 h-12 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
      <p class="text-gray-500 font-medium animate-pulse">Đang tải danh sách phỏng vấn...</p>
    </div>
    
    <div v-else-if="interviews.length === 0" class="bg-slate-50 border-2 border-dashed border-slate-200 rounded-3xl p-12 text-center flex flex-col items-center justify-center relative overflow-hidden group">
      <div class="absolute inset-0 bg-gradient-to-br from-blue-50 to-indigo-50 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
      
      <div class="w-20 h-20 bg-indigo-100 text-indigo-600 rounded-2xl flex items-center justify-center mb-6 relative z-10 shadow-sm group-hover:scale-110 transition-transform duration-500">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/></svg>
      </div>
      <h2 class="text-2xl font-bold text-slate-800 mb-3 relative z-10">Chưa có lịch phỏng vấn</h2>
      <p class="text-slate-500 max-w-md mx-auto text-lg leading-relaxed relative z-10">
        Bạn hiện chưa có lịch phỏng vấn nào sắp tới. Khi nhà tuyển dụng gửi lời mời, lịch sẽ xuất hiện tại đây.
      </p>
    </div>
    
    <div v-else class="space-y-4">
      <div v-for="iv in interviews" :key="iv.id" class="bg-white border border-gray-100 rounded-2xl p-6 flex flex-col md:flex-row justify-between items-start md:items-center gap-6 hover:shadow-xl hover:border-indigo-100 transition-all duration-300 group relative overflow-hidden">
        
        <!-- Hover highlight effect -->
        <div class="absolute left-0 top-0 bottom-0 w-1 bg-gradient-to-b from-blue-500 to-indigo-600 opacity-0 group-hover:opacity-100 transition-opacity"></div>
        
        <div class="flex-1 min-w-0">
          <h3 class="text-xl font-bold text-gray-800 mb-1 truncate group-hover:text-indigo-600 transition-colors">{{ iv.job_title || iv.title }}</h3>
          <p class="text-gray-500 font-medium mb-4 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M4 4a2 2 0 012-2h8a2 2 0 012 2v12a1 1 0 110 2h-3a1 1 0 01-1-1v-2a1 1 0 00-1-1H9a1 1 0 00-1 1v2a1 1 0 01-1 1H4a1 1 0 110-2V4zm3 1h2v2H7V5zm2 4H7v2h2V9zm2-4h2v2h-2V5zm2 4h-2v2h2V9z" clip-rule="evenodd" /></svg>
            {{ iv.company_name || 'Công ty ẩn danh' }}
          </p>
          
          <div class="flex flex-wrap gap-2 items-center">
            <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-blue-50 text-blue-700 text-sm font-bold border border-blue-100">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" /></svg>
              {{ formatDate(iv.scheduled_at) }}
            </span>
            
            <span class="inline-flex items-center px-3 py-1 rounded-lg text-sm font-bold border" :class="iv.mode === 'real' ? 'bg-rose-50 text-rose-700 border-rose-100' : 'bg-slate-100 text-slate-700 border-slate-200'">
              {{ iv.mode === 'real' ? 'Phỏng vấn thật' : 'Phỏng vấn thử' }}
            </span>
            
            <span v-if="iv.status === 'completed'" class="inline-flex items-center px-3 py-1 rounded-lg bg-emerald-50 text-emerald-700 text-sm font-bold border border-emerald-100">
              Đã hoàn thành
            </span>
            <span v-else-if="iv.status === 'cancelled'" class="inline-flex items-center px-3 py-1 rounded-lg bg-amber-50 text-amber-700 text-sm font-bold border border-amber-100">
              Đã hủy
            </span>
          </div>
        </div>
        
        <div class="flex gap-3 w-full md:w-auto shrink-0">
          <button v-if="iv.mode === 'real' && iv.status !== 'completed' && iv.status !== 'cancelled'" @click="iv.join_link ? router.push(iv.join_link) : null" class="w-full md:w-auto px-6 py-2.5 bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-bold rounded-xl shadow-md hover:shadow-lg hover:shadow-indigo-500/30 transform hover:-translate-y-0.5 transition-all flex items-center justify-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z" /></svg>
            Tham gia
          </button>
        </div>
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
  animation: fade-in 0.4s ease-out forwards;
}
</style>
