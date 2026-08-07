<script setup>
import { ref, onMounted } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Table from '../../components/common/AppTable.vue'
import Toast from '../../components/common/AppToast.vue'
import { authStore } from '../../stores/auth.store'
import { auditService } from '../../services/notification.service'
import { FileText, RefreshCw, Search } from 'lucide-vue-next'

const logs = ref([])
const loading = ref(true)
const toast = ref(null)
const page = ref(1)
const resourceType = ref('')

const columns = [
  { header: 'Hành động', key: 'action' },
  { header: 'Loại tài nguyên', key: 'resource_type' },
  { header: 'ID Tài nguyên', key: 'resource_id' },
  { header: 'Người thực hiện', key: 'actor' },
  { header: 'Thời gian', key: 'created_at' }
]

const loadLogs = async () => {
  loading.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) return
    const params = { page: page.value, page_size: 20 }
    if (resourceType.value) params.resource_type = resourceType.value
    
    const data = await auditService.getAuditLogs(companyId, params)
    // Map data
    logs.value = (Array.isArray(data) ? data : (data.data || [])).map(log => ({
      id: log.id,
      action: log.action || '—',
      resource_type: log.resource_type || '—',
      resource_id: log.resource_id || '—',
      actor: log.actor_user_id || 'System',
      created_at: log.created_at ? new Date(log.created_at).toLocaleString('vi-VN') : '—'
    }))
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi tải nhật ký hệ thống.' }
  } finally {
    loading.value = false
  }
}

onMounted(loadLogs)
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px">
      <div>
        <h1 class="text-h1">Nhật ký hệ thống</h1>
        <p class="text-helper" style="margin-top: 4px">Xem lịch sử các thao tác thay đổi dữ liệu trong hệ thống.</p>
      </div>
      <Button variant="secondary" @click="loadLogs">
        <RefreshCw size="16" style="margin-right: 8px" :class="{ 'animate-spin': loading }" /> 
        Làm mới
      </Button>
    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <Card>
      <div style="display: flex; gap: 16px; margin-bottom: 24px">
        <select class="input-field" v-model="resourceType" @change="loadLogs" style="width: 240px">
          <option value="">Tất cả tài nguyên</option>
          <option value="job">Công việc (Job)</option>
          <option value="candidate">Ứng viên (Candidate)</option>
          <option value="interview">Phỏng vấn (Interview)</option>
          <option value="rubric">Tiêu chí (Rubric)</option>
          <option value="question_bank">Kho câu hỏi</option>
        </select>
      </div>

      <div v-if="loading" style="padding: 40px 0; text-align: center; color: var(--text-muted)">
        <RefreshCw class="animate-spin" size="24" style="margin: 0 auto 16px" />
        <p>Đang tải dữ liệu...</p>
      </div>
      
      <Table v-else :columns="columns" :data="logs">
        <template #action="{ row }">
          <span style="font-weight: 500; color: var(--primary)">{{ row.action }}</span>
        </template>
        <template #resource_type="{ row }">
          <span class="badge badge-neutral" style="text-transform: capitalize">{{ row.resource_type }}</span>
        </template>
        <template #actor="{ row }">
          <span style="color: var(--text-secondary)">{{ row.actor }}</span>
        </template>
      </Table>
      
      <div v-if="!loading && logs.length === 0" style="padding: 40px 0; text-align: center; color: var(--text-muted)">
        <FileText size="48" style="margin: 0 auto 16px; opacity: 0.5" />
        <p>Không có dữ liệu nhật ký.</p>
      </div>
    </Card>
  </div>
</template>
