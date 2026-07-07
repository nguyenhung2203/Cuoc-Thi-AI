<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../../services/api.service'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Briefcase, MapPin, Clock, Building2, Search, ChevronLeft, ChevronRight } from 'lucide-vue-next'

const router = useRouter()
const jobs = ref([])
const loading = ref(true)
const keyword = ref('')

const currentPage = ref(1)
const totalPages = ref(1)

const unwrap = (val) => {
  if (!val) return ''
  if (typeof val === 'object' && 'String' in val) {
    return val.Valid ? val.String : ''
  }
  return val
}

const fetchJobs = async (page = 1) => {
  loading.value = true
  currentPage.value = page
  try {
    const params = new URLSearchParams()
    params.append('page', page)
    params.append('page_size', 12)
    if (keyword.value) params.append('keyword', keyword.value)
    
    const res = await apiService.getWithMeta(`/public/all-jobs?${params.toString()}`)
    
    // Map to unwrap sql.NullString objects
    jobs.value = (res.data || []).map(j => ({
      ...j,
      company_name: unwrap(j.company_name),
      location: unwrap(j.location),
      employment_type: unwrap(j.employment_type),
      department: unwrap(j.department),
      level: unwrap(j.level),
    }))
    
    if (res.meta) {
      totalPages.value = res.meta.total_pages || 1
    }
  } catch (error) {
    console.error('Failed to load jobs', error)
    jobs.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchJobs(1))

const handleSearch = () => {
  fetchJobs(1)
}

const viewJob = (job) => {
  router.push(`/careers/${job.company_id}/jobs/${job.id}`)
}
</script>

<template>
  <div class="space-y-8 animate-fade-in pb-12 max-w-7xl mx-auto">
    <!-- Header with Gradient -->
    <div class="relative bg-gradient-to-r from-blue-600 to-indigo-700 rounded-3xl p-8 md:p-12 text-white shadow-xl overflow-hidden">
      <!-- Decorative elements -->
      <div class="absolute top-0 right-0 -mr-8 -mt-8 w-64 h-64 rounded-full bg-white opacity-10 blur-3xl"></div>
      <div class="absolute bottom-0 left-0 -ml-8 -mb-8 w-48 h-48 rounded-full bg-blue-400 opacity-20 blur-2xl"></div>
      
      <div class="relative z-10 max-w-2xl">
        <h1 class="text-4xl font-extrabold mb-4 tracking-tight">Khám phá cơ hội nghề nghiệp</h1>
        <p class="text-blue-100 text-lg mb-8">Tìm kiếm hàng ngàn việc làm phù hợp với kỹ năng và định hướng phát triển của bạn.</p>
        
        <!-- Search Box inside Header -->
        <div class="bg-white p-2 rounded-2xl shadow-lg flex flex-col md:flex-row gap-2 max-w-3xl focus-within:ring-4 focus-within:ring-blue-500/30 transition-shadow">
          <div class="relative flex-1 flex items-center">
            <Search class="absolute left-4 text-gray-400 w-5 h-5" />
            <input 
              v-model="keyword"
              type="text"
              placeholder="Nhập chức danh, từ khóa hoặc công ty..."
              class="w-full bg-transparent border-none focus:ring-0 text-gray-800 placeholder-gray-400 py-3 pl-12 pr-4 outline-none font-medium"
              @keydown.enter="handleSearch"
            />
          </div>
          <button @click="handleSearch" class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-8 rounded-xl transition-colors shadow-md hover:shadow-lg whitespace-nowrap">
            Tìm việc ngay
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
      <p class="text-gray-500 font-medium animate-pulse">Đang tìm kiếm việc làm phù hợp...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="jobs.length === 0" class="flex flex-col items-center justify-center py-24 px-4 text-center bg-white rounded-3xl shadow-sm border border-gray-100">
      <div class="w-24 h-24 bg-gray-50 rounded-full flex items-center justify-center mb-6 shadow-inner">
        <Briefcase class="w-12 h-12 text-gray-300" />
      </div>
      <h3 class="text-2xl font-bold text-gray-800 mb-2">Chưa tìm thấy kết quả</h3>
      <p class="text-gray-500 max-w-md">Rất tiếc, chúng tôi không tìm thấy vị trí nào phù hợp với từ khóa của bạn lúc này. Vui lòng thử lại với từ khóa khác.</p>
      <button @click="keyword = ''; handleSearch()" class="mt-6 text-blue-600 font-medium hover:text-blue-800 underline underline-offset-4">Xóa tìm kiếm</button>
    </div>

    <!-- Job List -->
    <div v-else class="space-y-8">
      <div class="flex items-center justify-between">
        <h2 class="text-xl font-bold text-gray-800 flex items-center gap-2">
          <Sparkles class="w-5 h-5 text-amber-500" v-if="!keyword" />
          <span v-if="keyword">Kết quả tìm kiếm cho "<span class="text-blue-600">{{ keyword }}</span>"</span>
          <span v-else>Việc làm đề xuất cho bạn</span>
        </h2>
        <span class="text-sm text-gray-500 font-medium bg-gray-100 px-3 py-1 rounded-full">Trang {{ currentPage }} / {{ totalPages }}</span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- Job Card -->
        <div v-for="job in jobs" :key="job.id" 
          @click="viewJob(job)"
          class="group bg-white rounded-2xl p-6 border border-gray-100 shadow-sm hover:shadow-xl hover:border-blue-200 transition-all duration-300 cursor-pointer flex flex-col h-full transform hover:-translate-y-1">
          
          <div class="flex justify-between items-start mb-4 gap-4">
            <h3 class="text-lg font-bold text-gray-900 group-hover:text-blue-600 transition-colors line-clamp-2 leading-tight flex-1">
              {{ job.title }}
            </h3>
            <span v-if="job.department" class="inline-flex items-center px-2.5 py-1 rounded-md text-xs font-bold bg-blue-50 text-blue-700 whitespace-nowrap">
              {{ job.department }}
            </span>
          </div>
          
          <div class="space-y-2.5 mb-6">
            <div v-if="job.company_name" class="flex items-center text-gray-600 text-sm">
              <Building2 class="w-4 h-4 mr-2.5 text-gray-400" />
              <span class="font-medium truncate">{{ job.company_name }}</span>
            </div>
            <div class="flex items-center text-gray-500 text-sm">
              <MapPin class="w-4 h-4 mr-2.5 text-gray-400" />
              <span class="truncate">{{ job.location || 'Bất kỳ' }}</span>
            </div>
            <div class="flex items-center text-gray-500 text-sm">
              <Clock class="w-4 h-4 mr-2.5 text-gray-400" />
              <span class="truncate">{{ job.employment_type || 'Full-time' }}</span>
            </div>
          </div>
          
          <div class="mt-auto pt-5 border-t border-gray-100">
            <p class="text-gray-600 text-sm line-clamp-2 leading-relaxed">
              {{ job.description }}
            </p>
            <div class="mt-4 flex items-center text-blue-600 font-semibold text-sm opacity-0 group-hover:opacity-100 transition-opacity">
              Xem chi tiết <ChevronRight class="w-4 h-4 ml-1" />
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center items-center gap-4 mt-12 pt-8">
        <button 
          :disabled="currentPage === 1" 
          @click="fetchJobs(currentPage - 1)"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl border border-gray-200 text-gray-700 font-medium hover:bg-gray-50 hover:text-blue-600 disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-gray-700 transition-colors">
          <ChevronLeft class="w-4 h-4" /> Trang trước
        </button>
        
        <div class="flex gap-2">
          <span class="px-4 py-2.5 rounded-xl bg-blue-50 text-blue-700 font-bold border border-blue-100">
            {{ currentPage }}
          </span>
          <span class="px-4 py-2.5 text-gray-400">/</span>
          <span class="px-4 py-2.5 rounded-xl text-gray-600 font-medium">
            {{ totalPages }}
          </span>
        </div>

        <button 
          :disabled="currentPage === totalPages" 
          @click="fetchJobs(currentPage + 1)"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl border border-gray-200 text-gray-700 font-medium hover:bg-gray-50 hover:text-blue-600 disabled:opacity-50 disabled:hover:bg-transparent disabled:hover:text-gray-700 transition-colors">
          Trang sau <ChevronRight class="w-4 h-4" />
        </button>
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
