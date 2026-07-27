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
import { maxLength, normalizeText, isOneOf } from '../../utils/validators.js'
import { JOB_STATUSES } from '../../utils/constants.js'
import { jobService } from '../../services/job.service'
import { authStore } from '../../stores/auth.store'

const router = useRouter()
const jobs = ref([])
const loading = ref(true)
const filters = ref({ keyword: '', status: '' })
const routeMessage = ref(history.state?.message || '')
const showDeleteModal = ref(false)
const showFilterModal = ref(false)
const deletingId = ref(null)
const localToast = ref(null)

const isValidId = (value) => typeof value === 'string' || typeof value === 'number'
const validateFilters = () => {
  filters.value.keyword = normalizeText(filters.value.keyword).slice(0, 255)
  if (filters.value.status && isOneOf(filters.value.status, JOB_STATUSES)) {
    filters.value.status = ''
    localToast.value = { type: 'error', message: 'Bộ lọc trạng thái không hợp lệ.' }
    return false
  }
  return true
}

const fetchJobs = async () => {
  if (!validateFilters()) return
  loading.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) {
      throw new Error('Không tìm thấy company ID')
    }
    
    const params = {}
    if (filters.value.keyword) params.keyword = filters.value.keyword
    if (filters.value.status) params.status = filters.value.status

    const response = await jobService.getJobs(companyId, params)
    const safeJobs = Array.isArray(response) ? response : []
    jobs.value = safeJobs.map(j => ({
      id: j?.id,
      title: j?.title || 'Chưa có tiêu đề',
      status: j?.status || 'draft',
      created: j?.created_at ? new Date(j.created_at).toLocaleDateString('vi-VN') : 'Chưa cập nhật',
      applicants: Number(j?.candidate_count || 0)
    }))
  } catch (error) {
    localToast.value = { type: 'error', message: 'Lỗi tải dữ liệu: ' + (error.message || 'Không xác định') }
    jobs.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchJobs()
})

const applyFilter = () => {
  showFilterModal.value = false
  fetchJobs()
  localToast.value = { type: 'success', message: 'Đã áp dụng bộ lọc!' }
}

const clearFilter = () => {
  filters.value.status = ''
  showFilterModal.value = false
  fetchJobs()
}

const columns = [
  { header: 'Công việc', key: 'title' },
  { header: 'Trạng thái', key: 'status' },
  { header: 'Ngày tạo', key: 'created' },
  { header: 'Số ứng viên', key: 'applicants' },
  { header: 'Hành động', key: 'action' }
]

const confirmDelete = async () => {
  if (!isValidId(deletingId.value)) {
    localToast.value = { type: 'error', message: 'Không xác định được công việc cần xóa.' }
    return
  }
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await jobService.deleteJob(companyId, deletingId.value)
    jobs.value = jobs.value.filter(j => j.id !== deletingId.value)
    showDeleteModal.value = false
    localToast.value = { type: 'success', message: 'Đã xóa công việc thành công!' }
  } catch (error) {
    showDeleteModal.value = false
    localToast.value = { type: 'error', message: 'Lỗi xóa công việc: ' + (error.message || 'Không xác định') }
  }
}
</script>

<template>
  <div class="animate-fade-in space-y-6">
    <Toast v-if="routeMessage" type="success" :message="routeMessage" @close="routeMessage = ''" />
    <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" />
    
    <!-- Page Header -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Việc làm</h1>
        <p class="text-slate-500 dark:text-slate-400 mt-1 text-sm">Quản lý các vị trí tuyển dụng và pipeline ứng viên.</p>
      </div>
      <Button @click="router.push('/jobs/new')" class="bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white border-none shadow-md">
        <Plus size="16" class="mr-1" /> Tạo job mới
      </Button>
    </Card>

    <!-- Main Content Card -->
    <Card class="rounded-2xl shadow-sm overflow-hidden">
      <!-- Toolbar -->
      <div class="p-6 border-b border-slate-200 dark:border-slate-700 flex flex-col sm:flex-row gap-4 bg-slate-50/50 dark:bg-slate-800/50">
        <div class="relative flex-1 max-w-md group">
          <Search size="18" class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-[var(--primary)] transition-colors" />
          <input 
            type="text" 
            placeholder="Tìm kiếm công việc..." 
            class="w-full pl-10 pr-4 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[var(--primary)]/50 focus:border-[var(--primary)] transition-all text-[var(--text-primary)] placeholder:text-[var(--text-secondary)]"
            v-model="filters.keyword"
            @keyup.enter="fetchJobs"
          />
        </div>
        <Button variant="secondary" @click="showFilterModal = true" class="bg-[var(--surface)] border-[var(--border)] text-[var(--text-primary)] hover:bg-[var(--surface-hover)]">
          Lọc theo trạng thái
        </Button>
      </div>

      <!-- Table Section -->
      <div v-if="loading" class="p-12 text-center text-slate-500 dark:text-slate-400 flex flex-col items-center justify-center gap-3">
        <div class="w-8 h-8 border-4 border-[var(--primary-light)] border-t-[var(--primary)] rounded-full animate-spin"></div>
        <span class="text-sm font-medium">Đang tải dữ liệu...</span>
      </div>
      
      <div v-else class="w-full overflow-x-auto">
        <Table :columns="columns" :data="jobs" class="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <template #title="{ row }">
            <div class="font-semibold text-slate-800 dark:text-slate-200 whitespace-nowrap overflow-hidden text-ellipsis max-w-[250px] cursor-pointer hover:text-[var(--primary)] dark:hover:text-[var(--primary-light)] transition-colors" :title="row.title" @click="router.push(`/jobs/${row.id}`)">
              {{ row.title }}
            </div>
          </template>
          <template #status="{ row }">
            <span 
              class="px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider rounded-full border"
              :class="[
                row.status.toLowerCase() === 'open' ? 'bg-[var(--success-bg)] text-[var(--success)] border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20' : 
                row.status.toLowerCase() === 'closed' ? 'bg-rose-50 text-rose-600 border-rose-200 dark:bg-rose-500/10 dark:text-rose-400 dark:border-rose-500/20' : 
                'bg-slate-100 text-slate-600 border-slate-200 dark:bg-slate-700 dark:text-slate-300 dark:border-slate-600'
              ]"
            >
              {{ row.status }}
            </span>
          </template>
          <template #created="{ row }">
            <div class="text-slate-500 dark:text-slate-400 font-medium">{{ row.created }}</div>
          </template>
          <template #applicants="{ row }">
            <div class="inline-flex items-center justify-center px-2.5 py-1 rounded-lg bg-[var(--primary-light)] dark:bg-blue-500/10 text-[var(--primary)] dark:text-blue-400 font-semibold border border-[var(--primary-light)] dark:border-blue-500/20">
              {{ row.applicants }}
            </div>
          </template>
          <template #action="{ row }">
            <div class="flex items-center gap-1">
              <button @click="router.push(`/jobs/${row.id}`)" class="p-2 text-slate-400 hover:text-[var(--primary)] dark:hover:text-blue-400 hover:bg-[var(--primary-light)] dark:hover:bg-blue-500/10 rounded-lg transition-colors" title="Xem chi tiết">
                <Eye size="18" />
              </button>
              <button @click="router.push(`/jobs/${row.id}`)" class="p-2 text-slate-400 hover:text-[var(--primary)] dark:hover:text-blue-400 hover:bg-[var(--primary-light)] dark:hover:bg-blue-500/10 rounded-lg transition-colors" title="Chỉnh sửa">
                <Edit size="18" />
              </button>
              <button @click="deletingId = row.id; showDeleteModal = true" class="p-2 text-slate-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-lg transition-colors" title="Xóa">
                <Trash2 size="18" />
              </button>
            </div>
          </template>
        </Table>
      </div>
    </Card>

    <!-- Modals -->
    <Modal :isOpen="showDeleteModal" @close="showDeleteModal = false" title="Xác nhận xóa">
      <div class="p-1">
        <p class="text-slate-600 dark:text-slate-300 mb-6 leading-relaxed">Bạn có chắc chắn muốn xóa công việc này? Các hồ sơ ứng viên liên quan có thể bị ảnh hưởng.</p>
        <div class="flex justify-end gap-3">
          <Button variant="ghost" @click="showDeleteModal = false" class="text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700">Hủy</Button>
          <Button variant="primary" class="bg-red-600 hover:bg-red-700 border-red-600 text-white" @click="confirmDelete">Xóa công việc</Button>
        </div>
      </div>
    </Modal>

    <Modal :isOpen="showFilterModal" @close="showFilterModal = false" title="Lọc công việc">
      <div class="flex flex-col gap-5 mb-6 p-1">
        <div class="space-y-2">
          <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Trạng thái công việc</label>
          <select 
            class="w-full px-4 py-2.5 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-[var(--text-primary)]"
            v-model="filters.status"
          >
            <option value="">Tất cả trạng thái</option>
            <option value="open">Đang mở (Open)</option>
            <option value="closed">Đã đóng (Closed)</option>
            <option value="draft">Bản nháp (Draft)</option>
            <option value="paused">Tạm dừng (Paused)</option>
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
