<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Table from '../../components/common/AppTable.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import Modal from '../../components/common/AppModal.vue'
import { Plus, Search, Eye, Edit, Trash2 } from 'lucide-vue-next'
import { candidateService } from '../../services/candidate.service'
import { jobService } from '../../services/job.service'
import { authStore } from '../../stores/auth.store'

const router = useRouter()
const candidates = ref([])
const loading = ref(true)
const routeMessage = ref(history.state?.message || '')
const showDeleteModal = ref(false)
const showFilterModal = ref(false)
const deletingId = ref(null)
const localToast = ref(null)
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = ref(1)

const jobs = ref([])
const filters = ref({ job_id: '', status: '', keyword: '' })

const fetchCandidates = async () => {
  loading.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) return

    // Clean up empty filters
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (filters.value.job_id) params.job_id = filters.value.job_id
    if (filters.value.status) params.status = filters.value.status
    if (filters.value.keyword) params.keyword = filters.value.keyword

    const response = await candidateService.getCandidates(companyId, params)

    // Handle both array and paginated response
    if (Array.isArray(response)) {
      candidates.value = response.map(c => ({
        id: c.id,
        name: c.full_name,
        email: c.email,
        appliedJob: c.latest_job?.title || 'Chưa ứng tuyển',
        status: c.status || 'New'
      }))
      totalPages.value = 1
    } else if (response.data && Array.isArray(response.data)) {
      candidates.value = response.data.map(c => ({
        id: c.id,
        name: c.full_name,
        email: c.email,
        appliedJob: c.latest_job?.title || 'Chưa ứng tuyển',
        status: c.status || 'New'
      }))
      totalPages.value = Math.ceil((response.meta?.total || 0) / pageSize.value)
    }
  } catch (error) {
    localToast.value = { type: 'error', message: 'Lỗi tải danh sách ứng viên: ' + (error.message || 'Không xác định') }
    candidates.value = []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) throw new Error('Không tìm thấy company ID')
    
    // Fetch both jobs for dropdown and initial candidates
    const [jobsRes] = await Promise.all([
      jobService.getJobs(companyId),
      fetchCandidates()
    ])
    jobs.value = jobsRes
  } catch (error) {
    console.error(error)
  }
})

const applyFilter = () => {
  showFilterModal.value = false
  fetchCandidates()
  localToast.value = { type: 'success', message: 'Đã áp dụng bộ lọc!' }
}

const clearFilter = () => {
  filters.value.job_id = ''
  filters.value.status = ''
  showFilterModal.value = false
  fetchCandidates()
}

const getStatusType = (status) => {
  if (status === 'New') return 'info'
  if (status === 'Interviewing') return 'warning'
  if (status === 'Offered') return 'success'
  if (status === 'Rejected') return 'danger'
  return 'neutral'
}

const confirmDelete = async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await candidateService.updateCandidate(companyId, deletingId.value, { status: 'deleted' }) // Hoặc gọi API delete nếu có
    candidates.value = candidates.value.filter(c => c.id !== deletingId.value)
    showDeleteModal.value = false
    localToast.value = { type: 'success', message: 'Đã xóa hồ sơ ứng viên thành công!' }
  } catch (error) {
    showDeleteModal.value = false
    localToast.value = { type: 'error', message: 'Lỗi xóa ứng viên: ' + (error.message || 'Không xác định') }
  }
}

const columns = [
  { header: 'Ứng viên', key: 'candidate' },
  { header: 'Vị trí ứng tuyển', key: 'appliedJob' },
  { header: 'Trạng thái', key: 'status' },
  { header: 'Hành động', key: 'action' }
]
</script>

<template>
  <div class="animate-fade-in space-y-6">
    <Toast v-if="routeMessage" type="success" :message="routeMessage" @close="routeMessage = ''" />
    <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" />
    
    <!-- Page Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Ứng viên</h1>
        <p class="text-slate-500 dark:text-slate-400 mt-1 text-sm">Theo dõi ứng viên, điểm phù hợp và lịch sử phỏng vấn.</p>
      </div>
      <Button @click="router.push('/candidates/new')" class="bg-blue-600 hover:bg-blue-700 text-white border-none shadow-md shadow-blue-500/20">
        <Plus size="16" class="mr-1" /> Thêm ứng viên
      </Button>
    </div>

    <!-- Main Content Card -->
    <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 overflow-hidden">
      <!-- Toolbar -->
      <div class="p-6 border-b border-slate-200 dark:border-slate-700 flex flex-col md:flex-row gap-4 bg-slate-50/50 dark:bg-slate-800/50">
        <div class="relative flex-1 max-w-md group">
          <Search size="18" class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-blue-500 transition-colors" />
          <input 
            type="text" 
            placeholder="Tìm tên, email..." 
            class="w-full pl-10 pr-4 py-2 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-slate-700 dark:text-slate-200 placeholder:text-slate-400"
            v-model="filters.keyword"
            @keyup.enter="fetchCandidates"
          />
        </div>
        <div class="flex gap-2">
          <Button variant="secondary" @click="showFilterModal = true" class="bg-white dark:bg-slate-800 border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700">
            Lọc theo Job
          </Button>
          <Button variant="secondary" @click="showFilterModal = true" class="bg-white dark:bg-slate-800 border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700">
            Lọc theo trạng thái
          </Button>
        </div>
      </div>

      <!-- Table Section -->
      <div v-if="loading" class="p-12 text-center text-slate-500 dark:text-slate-400 flex flex-col items-center justify-center gap-3">
        <div class="w-8 h-8 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
        <span class="text-sm font-medium">Đang tải dữ liệu...</span>
      </div>
      
      <div v-else class="w-full overflow-x-auto">
        <Table :columns="columns" :data="candidates" class="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <template #candidate="{ row }">
            <div class="cursor-pointer group flex flex-col" @click="router.push(`/candidates/${row.id}`)">
              <div class="font-semibold text-slate-800 dark:text-slate-200 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                {{ row.name }}
              </div>
              <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{{ row.email }}</div>
            </div>
          </template>
          
          <template #appliedJob="{ row }">
            <span class="text-slate-600 dark:text-slate-300 font-medium">{{ row.appliedJob }}</span>
          </template>
          
          <template #status="{ row }">
            <span 
              class="px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider rounded-full border"
              :class="[
                row.status === 'New' ? 'bg-blue-50 text-blue-600 border-blue-200 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/20' : 
                row.status === 'Interviewing' ? 'bg-amber-50 text-amber-600 border-amber-200 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/20' : 
                row.status === 'Offered' ? 'bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20' : 
                row.status === 'Rejected' ? 'bg-rose-50 text-rose-600 border-rose-200 dark:bg-rose-500/10 dark:text-rose-400 dark:border-rose-500/20' : 
                'bg-slate-100 text-slate-600 border-slate-200 dark:bg-slate-700 dark:text-slate-300 dark:border-slate-600'
              ]"
            >
              {{ row.status }}
            </span>
          </template>
          
          <template #action="{ row }">
            <div class="flex items-center gap-1">
              <button @click="router.push(`/candidates/${row.id}`)" class="p-2 text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-500/10 rounded-lg transition-colors" title="Xem chi tiết">
                <Eye size="18" />
              </button>
              <button @click="router.push(`/candidates/${row.id}`)" class="p-2 text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-500/10 rounded-lg transition-colors" title="Chỉnh sửa">
                <Edit size="18" />
              </button>
              <button @click="deletingId = row.id; showDeleteModal = true" class="p-2 text-slate-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-lg transition-colors" title="Xóa">
                <Trash2 size="18" />
              </button>
            </div>
          </template>
        </Table>
      </div>

      <!-- Pagination -->
      <div v-if="!loading && candidates.length > 0" class="p-6 border-t border-slate-200 dark:border-slate-700 flex items-center justify-between bg-slate-50/50 dark:bg-slate-800/50">
        <div class="text-sm text-slate-600 dark:text-slate-400">
          Trang {{ currentPage }} / {{ totalPages || 1 }}
        </div>
        <div class="flex gap-2">
          <Button
            :disabled="currentPage <= 1"
            variant="secondary"
            @click="currentPage > 1 && (currentPage--, fetchCandidates())"
            class="bg-white dark:bg-slate-800 border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            ← Trước
          </Button>
          <Button
            :disabled="currentPage >= totalPages"
            variant="secondary"
            @click="currentPage < totalPages && (currentPage++, fetchCandidates())"
            class="bg-white dark:bg-slate-800 border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Tiếp →
          </Button>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <Modal :isOpen="showDeleteModal" @close="showDeleteModal = false" title="Xác nhận xóa">
      <div class="p-1">
        <p class="text-slate-600 dark:text-slate-300 mb-6 leading-relaxed">Bạn có chắc chắn muốn xóa hồ sơ ứng viên này? Các dữ liệu phỏng vấn liên quan có thể bị ảnh hưởng.</p>
        <div class="flex justify-end gap-3">
          <Button variant="ghost" @click="showDeleteModal = false" class="text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700">Hủy</Button>
          <Button variant="primary" class="bg-red-600 hover:bg-red-700 border-red-600 text-white" @click="confirmDelete">Xóa ứng viên</Button>
        </div>
      </div>
    </Modal>

    <Modal :isOpen="showFilterModal" @close="showFilterModal = false" title="Lọc ứng viên">
      <div class="flex flex-col gap-5 mb-6 p-1">
        <div class="space-y-2">
          <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Vị trí ứng tuyển (Job)</label>
          <select 
            class="w-full px-4 py-2.5 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-slate-700 dark:text-slate-200"
            v-model="filters.job_id"
          >
            <option value="">Tất cả vị trí</option>
            <option v-for="job in jobs" :key="job.id" :value="job.id">{{ job.title }}</option>
          </select>
        </div>
        <div class="space-y-2">
          <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Trạng thái hồ sơ</label>
          <select 
            class="w-full px-4 py-2.5 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-slate-700 dark:text-slate-200"
            v-model="filters.status"
          >
            <option value="">Tất cả trạng thái</option>
            <option value="New">Mới (New)</option>
            <option value="Interviewing">Đang phỏng vấn (Interviewing)</option>
            <option value="Offered">Đã gửi Offer (Offered)</option>
            <option value="Rejected">Từ chối (Rejected)</option>
          </select>
        </div>
      </div>
      <div class="flex justify-end gap-3 p-1">
        <Button variant="ghost" @click="clearFilter" class="text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700">Xóa bộ lọc</Button>
        <Button variant="primary" @click="applyFilter" class="bg-blue-600 hover:bg-blue-700 text-white border-none shadow-md shadow-blue-500/20">Áp dụng</Button>
      </div>
    </Modal>
  </div>
</template>
