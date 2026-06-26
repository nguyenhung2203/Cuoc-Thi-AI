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
import { mockApi } from '../../utils/mockData'

const router = useRouter()
const jobs = ref([])
const loading = ref(true)
const routeMessage = ref(history.state?.message || '')
const showDeleteModal = ref(false)
const showFilterModal = ref(false)
const deletingId = ref(null)
const localToast = ref(null)

onMounted(async () => {
  const data = await mockApi.jobs.getAll()
  jobs.value = data
  loading.value = false
})

const columns = [
  { header: 'Công việc', key: 'title' },
  { header: 'Trạng thái', key: 'status' },
  { header: 'Ngày tạo', key: 'created' },
  { header: 'Số ứng viên', key: 'applicants' },
  { header: 'Hành động', key: 'action' }
]

const confirmDelete = () => {
  jobs.value = jobs.value.filter(j => j.id !== deletingId.value)
  showDeleteModal.value = false
  localToast.value = { type: 'success', message: 'Đã xóa công việc thành công!' }
}
</script>

<template>
  <div>
    <Toast v-if="routeMessage" type="success" :message="routeMessage" @close="routeMessage = ''" />
    <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" />
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px">
      <div>
        <h1 class="text-h1">Việc làm</h1>
        <p class="text-helper" style="margin-top: 4px">Quản lý các vị trí tuyển dụng và pipeline ứng viên.</p>
      </div>
      <Button @click="router.push('/jobs/new')"><Plus size="16" /> Tạo job mới</Button>
    </div>

    <Card>
      <div style="display: flex; gap: 16px; margin-bottom: 24px">
        <div style="position: relative; flex: 1; max-width: 320px">
          <Search size="16" style="position: absolute; left: 12px; top: 12px; color: var(--text-muted)" />
          <input 
            type="text" 
            placeholder="Tìm kiếm công việc..." 
            class="input-field"
            style="width: 100%; padding-left: 36px"
          />
        </div>
        <Button variant="secondary" @click="showFilterModal = true">Lọc theo trạng thái</Button>
      </div>

      <div v-if="loading" style="padding: 32px; text-align: center; color: var(--text-muted)">
        Đang tải dữ liệu...
      </div>
      <Table v-else :columns="columns" :data="jobs">
        <template #title="{ row }">
          <div style="font-weight: 500; color: var(--text-main); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 250px;" :title="row.title">
            {{ row.title }}
          </div>
        </template>
        <template #status="{ row }">
          <Badge :type="row.status === 'Open' ? 'success' : 'neutral'">
            {{ row.status }}
          </Badge>
        </template>
        <template #action="{ row }">
          <div style="display: flex; gap: 8px">
            <Button variant="ghost" @click="router.push(`/jobs/${row.id}`)" style="padding: 0 8px">
              <Eye size="16" />
            </Button>
            <Button variant="ghost" @click="router.push(`/jobs/${row.id}`)" style="padding: 0 8px">
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
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn xóa công việc này? Các hồ sơ ứng viên liên quan có thể bị ảnh hưởng.</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showDeleteModal = false">Hủy</Button>
        <Button variant="primary" style="background-color: var(--danger); border-color: var(--danger)" @click="confirmDelete">Xóa công việc</Button>
      </div>
    </Modal>

    <Modal :isOpen="showFilterModal" @close="showFilterModal = false" title="Lọc công việc">
      <div style="display: flex; flex-direction: column; gap: 16px; margin-bottom: 24px">
        <div class="input-group">
          <label class="input-label">Trạng thái công việc</label>
          <select class="input-field">
            <option value="">Tất cả trạng thái</option>
            <option value="Open">Đang mở (Open)</option>
            <option value="Closed">Đã đóng (Closed)</option>
          </select>
        </div>
      </div>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showFilterModal = false">Xóa bộ lọc</Button>
        <Button variant="primary" @click="showFilterModal = false; localToast = { type: 'success', message: 'Đã áp dụng bộ lọc!' }">Áp dụng</Button>
      </div>
    </Modal>
  </div>
</template>
