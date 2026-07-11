<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { candidatePortalService } from '../../services/candidate-portal.service'
import Card from '../../components/common/AppCard.vue'
import Badge from '../../components/common/AppBadge.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import { Briefcase, MapPin, Building2, Banknote, Bookmark, Trash2, Heart, Search, Filter, ExternalLink, ArrowRight } from 'lucide-vue-next'

const router = useRouter()
const savedJobs = ref([])
const loading = ref(true)
const searchQuery = ref('')
const filterLocation = ref('ALL')
const toast = ref(null)

onMounted(async () => {
  try {
    savedJobs.value = await candidatePortalService.getSavedJobs()
  } catch (err) {
    console.error('Lỗi tải danh sách việc làm đã lưu:', err)
  } finally {
    loading.value = false
  }
})

const removeJob = (id, title) => {
  savedJobs.value = savedJobs.value.filter(job => job.id !== id)
  toast.value = { type: 'success', message: `Đã bỏ lưu việc làm "${title || 'chọn'}" khỏi danh sách.` }
}

const clearAllSaved = () => {
  if (confirm('Bạn có chắc chắn muốn xóa toàn bộ danh sách việc làm đã lưu?')) {
    savedJobs.value = []
    toast.value = { type: 'info', message: 'Đã làm trống danh sách việc làm đã lưu.' }
  }
}

const filteredJobs = computed(() => {
  return savedJobs.value.filter(job => {
    const matchesSearch = !searchQuery.value || 
      job.title.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
      job.company_name.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesLocation = filterLocation.value === 'ALL' || 
      (job.location && job.location.toLowerCase().includes(filterLocation.value.toLowerCase()))
      
    return matchesSearch && matchesLocation
  })
})

const getSalaryDisplay = (job) => {
  const min = job.salary_min?.Valid ? job.salary_min.Int64 : null
  const max = job.salary_max?.Valid ? job.salary_max.Int64 : null
  const curr = job.currency?.Valid ? job.currency.String : 'VND'
  if (!min && !max) return 'Thỏa thuận'
  if (min && !max) return `Từ ${min.toLocaleString()} ${curr}`
  if (!min && max) return `Đến ${max.toLocaleString()} ${curr}`
  return `${min.toLocaleString()} - ${max.toLocaleString()} ${curr}`
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const d = new Date(dateString)
  return d.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })
}
</script>

<template>
  <div class="saved-jobs-page space-y-6">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <!-- Khung bọc tiêu đề chuẩn Design System -->
    <div class="header-box animate-rise flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 rounded-2xl bg-[var(--surface)] border border-[var(--border)] shadow-sm gap-4">
      <div>
        <h1 class="text-h1">Việc làm đã lưu</h1>
        <p class="text-secondary mt-1">Danh sách các cơ hội nghề nghiệp bạn quan tâm để chuẩn bị ứng tuyển.</p>
      </div>
      <div class="flex items-center gap-3 shrink-0">
        <Badge class="badge-primary px-3.5 py-1.5 text-sm font-bold">{{ savedJobs.length }} việc làm</Badge>
        <Button variant="primary" @click="router.push('/job-board')" class="flex items-center gap-2">
          <Search :size="16" /> Tìm thêm việc làm
        </Button>
      </div>
    </div>

    <!-- Thanh tìm kiếm và bộ lọc nhanh -->
    <div v-if="savedJobs.length > 0 || searchQuery || filterLocation !== 'ALL'" class="bg-[var(--surface)] p-5 rounded-2xl border border-[var(--border)] shadow-sm space-y-4">
      <div class="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
        <!-- Input search -->
        <div class="relative flex-1">
          <Search :size="18" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-muted)] pointer-events-none" />
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Tìm theo tên công việc, kỹ năng hoặc công ty..." 
            class="w-full pl-10 pr-4 py-2.5 bg-[var(--surface-soft)] border border-[var(--border)] rounded-xl text-sm font-medium text-[var(--text-main)] focus:outline-none focus:border-[var(--primary)] transition-colors"
          />
        </div>

        <!-- Filter theo địa điểm -->
        <div class="flex items-center gap-3 shrink-0">
          <div class="flex items-center gap-2 px-3 py-2 bg-[var(--surface-soft)] border border-[var(--border)] rounded-xl text-sm font-semibold text-[var(--text-main)]">
            <Filter :size="15" class="text-[var(--primary)]" />
            <span>Địa điểm:</span>
            <select v-model="filterLocation" class="bg-transparent border-none font-bold text-[var(--primary)] focus:outline-none cursor-pointer pr-2">
              <option value="ALL">Tất cả khu vực</option>
              <option value="Hà Nội">Hà Nội</option>
              <option value="TP. HCM">TP. HCM</option>
              <option value="Đà Nẵng">Đà Nẵng</option>
              <option value="Remote">Làm việc từ xa (Remote)</option>
            </select>
          </div>

          <Button v-if="savedJobs.length > 0" variant="ghost" size="sm" @click="clearAllSaved" class="text-xs text-[var(--danger)] hover:bg-red-50 border border-transparent hover:border-red-200">
            <Trash2 :size="14" class="mr-1" /> Xóa tất cả
          </Button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center p-12">
      <div class="spinner"></div>
    </div>

    <!-- Empty state chuẩn Design System -->
    <div v-else-if="filteredJobs.length === 0" class="empty-state bg-[var(--surface)] p-12 rounded-2xl border border-dashed border-[var(--border)] text-center max-w-xl mx-auto my-8">
      <div class="w-16 h-16 rounded-full bg-[var(--surface-soft)] text-rose-500 flex items-center justify-center mx-auto mb-4">
        <Heart :size="32" />
      </div>
      <h3 class="text-base font-bold text-[var(--text-main)]">Chưa có việc làm nào phù hợp</h3>
      <p class="text-sm text-[var(--text-secondary)] mt-1.5 leading-relaxed">
        {{ savedJobs.length === 0 ? 'Bạn chưa lưu công việc nào vào danh mục quan tâm. Hãy lướt xem bảng tin tuyển dụng để tìm vị trí ưng ý nhé!' : 'Không tìm thấy việc làm nào khớp với từ khóa tìm kiếm hoặc bộ lọc hiện tại.' }}
      </p>
      <div class="flex items-center justify-center gap-3 mt-6">
        <Button v-if="searchQuery || filterLocation !== 'ALL'" variant="ghost" @click="searchQuery = ''; filterLocation = 'ALL'">Xóa bộ lọc</Button>
        <Button variant="primary" @click="router.push('/job-board')">Khám phá việc làm ngay</Button>
      </div>
    </div>

    <!-- Grid danh sách việc làm đã lưu -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card 
        v-for="job in filteredJobs" 
        :key="job.id" 
        class="flex flex-col justify-between p-6 rounded-2xl border border-[var(--border)] hover:border-[var(--primary)] transition-all duration-300 shadow-sm animate-rise cursor-pointer group"
        @click="router.push(`/careers/${job.company_id}/jobs/${job.id}`)"
      >
        <div class="space-y-4">
          <!-- Top row -->
          <div class="flex justify-between items-start gap-3">
            <h3 class="text-base font-bold text-[var(--text-main)] group-hover:text-[var(--primary)] transition-colors line-clamp-2 leading-snug">
              {{ job.title }}
            </h3>
            <button 
              @click.stop="removeJob(job.id, job.title)" 
              class="p-2 rounded-xl text-rose-500 hover:bg-rose-50 transition-colors shrink-0" 
              title="Bỏ lưu công việc này"
            >
              <Trash2 :size="16" />
            </button>
          </div>

          <!-- Tên công ty -->
          <p class="text-sm font-semibold text-[var(--text-secondary)] flex items-center gap-2">
            <Building2 :size="16" class="text-[var(--primary)]" />
            {{ job.company_name }}
          </p>

          <!-- Tags địa điểm & lương -->
          <div class="flex flex-wrap gap-2 pt-1">
            <Badge class="badge-info px-3 py-1 text-xs font-semibold flex items-center gap-1.5">
              <MapPin :size="13" /> {{ job.location || 'Bất kỳ' }}
            </Badge>
            <Badge class="badge-success px-3 py-1 text-xs font-semibold flex items-center gap-1.5">
              <Banknote :size="13" /> {{ getSalaryDisplay(job) }}
            </Badge>
          </div>
        </div>

        <!-- Footer -->
        <div class="border-t border-[var(--border)] pt-4 mt-6 flex items-center justify-between gap-3">
          <span class="text-xs text-[var(--text-muted)] flex items-center gap-1">
            <Bookmark :size="12" /> Đã lưu: {{ formatDate(job.saved_at) }}
          </span>
          <Button 
            variant="primary" 
            size="sm" 
            @click.stop="router.push(`/careers/${job.company_id}/jobs/${job.id}`)"
            class="text-xs flex items-center gap-1.5 shadow-sm"
          >
            Ứng tuyển <ArrowRight :size="14" />
          </Button>
        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.saved-jobs-page {
  animation: fadeIn 0.4s ease-out;
}
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
