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

const jobs = ref([])
const filters = ref({ job_id: '', status: '', keyword: '' })

const fetchCandidates = async () => {
  loading.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) return
    
    // Clean up empty filters
    const params = {}
    if (filters.value.job_id) params.job_id = filters.value.job_id
    if (filters.value.status) params.status = filters.value.status
    if (filters.value.keyword) params.keyword = filters.value.keyword

    const response = await candidateService.getCandidates(companyId, params)
    candidates.value = response.map(c => ({
      id: c.id,
      name: c.full_name,
      email: c.email,
      appliedJob: c.latest_job?.title || 'Chưa ứng tuyển',
      status: c.status || 'New'
    }))
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
  <div>
    <Toast v-if="routeMessage" type="success" :message="routeMessage" @close="routeMessage = ''" />
    <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" />
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px">
      <div>
        <h1 class="text-h1">Ứng viên</h1>
        <p class="text-helper" style="margin-top: 4px">Theo dõi ứng viên, điểm phù hợp và lịch sử phỏng vấn.</p>
      </div>
      <Button @click="router.push('/candidates/new')"><Plus size="16" /> Thêm ứng viên</Button>
    </div>

    <Card>
      <div style="display: flex; gap: 16px; margin-bottom: 24px">
        <div style="position: relative; flex: 1; max-width: 320px">
          <Search size="16" style="position: absolute; left: 12px; top: 12px; color: var(--text-muted)" />
          <input 
            type="text" 
            placeholder="Tìm tên, email..." 
            class="input-field"
            style="width: 100%; padding-left: 36px"
            v-model="filters.keyword"
            @keyup.enter="fetchCandidates"
          />
        </div>
        <Button variant="secondary" @click="showFilterModal = true">Lọc theo Job</Button>
        <Button variant="secondary" @click="showFilterModal = true">Lọc theo trạng thái</Button>
      </div>

      <div v-if="loading" style="padding: 32px; text-align: center; color: var(--text-muted)">
        Đang tải dữ liệu...
      </div>
      <Table v-else :columns="columns" :data="candidates">
        <template #candidate="{ row }">
          <div>
            <div style="font-weight: 500; color: var(--text-main)">{{ row.name }}</div>
            <div class="text-helper">{{ row.email }}</div>
          </div>
        </template>
        <template #status="{ row }">
          <Badge :type="getStatusType(row.status)">
            {{ row.status }}
          </Badge>
        </template>
        <template #action="{ row }">
          <div style="display: flex; gap: 8px">
            <Button variant="ghost" @click="router.push(`/candidates/${row.id}`)" style="padding: 0 8px">
              <Eye size="16" />
            </Button>
            <Button variant="ghost" @click="router.push(`/candidates/${row.id}`)" style="padding: 0 8px">
              <Edit size="16" />
            </Button>
            <Button variant="ghost" @click="deletingId = row.id; showDeleteModal = true" style="padding: 0 8px; color: var(--danger)">
              <Trash2 size="16" />
            </Button>
          </div>
        </template>
      </Table>
    </Card>

    <Modal :isOpen="showDeleteModal" @close="showDeleteModal = false" title="Xác nhận xóa">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn xóa hồ sơ ứng viên này? Các dữ liệu phỏng vấn liên quan có thể bị ảnh hưởng.</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showDeleteModal = false">Hủy</Button>
        <Button variant="primary" style="background-color: var(--danger); border-color: var(--danger)" @click="confirmDelete">Xóa ứng viên</Button>
      </div>
    </Modal>

    <Modal :isOpen="showFilterModal" @close="showFilterModal = false" title="Lọc ứng viên">
      <div style="display: flex; flex-direction: column; gap: 16px; margin-bottom: 24px">
        <div class="input-group">
          <label class="input-label">Vị trí ứng tuyển (Job)</label>
          <select class="input-field" v-model="filters.job_id">
            <option value="">Tất cả vị trí</option>
            <option v-for="job in jobs" :key="job.id" :value="job.id">{{ job.title }}</option>
          </select>
        </div>
        <div class="input-group">
          <label class="input-label">Trạng thái hồ sơ</label>
          <select class="input-field" v-model="filters.status">
            <option value="">Tất cả trạng thái</option>
            <option value="New">Mới (New)</option>
            <option value="Interviewing">Đang phỏng vấn (Interviewing)</option>
            <option value="Offered">Đã gửi Offer (Offered)</option>
            <option value="Rejected">Từ chối (Rejected)</option>
          </select>
        </div>
      </div>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="clearFilter">Xóa bộ lọc</Button>
        <Button variant="primary" @click="applyFilter">Áp dụng</Button>
      </div>
    </Modal>
  </div>
</template>
