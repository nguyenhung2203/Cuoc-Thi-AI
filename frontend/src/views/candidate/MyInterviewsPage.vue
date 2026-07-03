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
  <div class="page-header">
    <div>
      <h1 class="text-h1">Phỏng vấn của tôi</h1>
      <p class="text-helper" style="margin-top: 4px">Quản lý các lịch phỏng vấn sắp tới và lịch sử phỏng vấn.</p>
    </div>
  </div>
  
  <div v-if="loading" style="padding: 48px; text-align: center; color: var(--text-muted)">
    Đang tải dữ liệu...
  </div>
  
  <div v-else-if="interviews.length === 0" style="padding: 48px; text-align: center; border: 1px dashed var(--border); border-radius: var(--radius-lg); background-color: var(--surface-soft)">
    <div style="margin-bottom: 16px; display: inline-flex; justify-content: center; align-items: center; width: 64px; height: 64px; border-radius: 50%; background-color: rgba(37, 99, 235, 0.1); color: var(--primary)">
      <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/></svg>
    </div>
    <h2 class="text-h2" style="margin-bottom: 8px">Chưa có lịch phỏng vấn</h2>
    <p class="text-body" style="color: var(--text-secondary); max-width: 400px; margin: 0 auto">
      Bạn hiện chưa có lịch phỏng vấn nào sắp tới. Khi nhà tuyển dụng gửi lời mời, lịch sẽ xuất hiện tại đây.
    </p>
  </div>
  
  <div v-else style="margin-top: 24px; display: flex; flex-direction: column; gap: 16px">
    <div v-for="iv in interviews" :key="iv.id" style="background-color: var(--surface); border: 1px solid var(--border); border-radius: 8px; padding: 20px; display: flex; justify-content: space-between; align-items: center; box-shadow: 0 1px 2px rgba(0,0,0,0.05)">
      <div>
        <div style="font-size: 18px; font-weight: 600; margin-bottom: 4px">{{ iv.job_title || iv.title }}</div>
        <div class="text-body" style="color: var(--text-secondary); font-weight: 500">{{ iv.company_name || 'Công ty ẩn danh' }}</div>
        <div style="display: flex; gap: 12px; margin-top: 12px">
          <Badge type="primary">{{ formatDate(iv.scheduled_at) }}</Badge>
          <Badge :type="iv.mode === 'real' ? 'danger' : 'default'">{{ iv.mode === 'real' ? 'Phỏng vấn thật' : 'Phỏng vấn thử' }}</Badge>
          <Badge v-if="iv.status === 'completed'" type="success">Đã hoàn thành</Badge>
          <Badge v-else-if="iv.status === 'cancelled'" type="warning">Đã hủy</Badge>
        </div>
      </div>
      <div style="display: flex; gap: 8px">
        <Button v-if="iv.mode === 'real' && iv.status !== 'completed' && iv.status !== 'cancelled'" variant="primary" @click="iv.join_link ? router.push(iv.join_link) : null">Tham gia</Button>
      </div>
    </div>
  </div>
</template>
