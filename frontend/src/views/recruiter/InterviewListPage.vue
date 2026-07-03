<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Table from '../../components/common/AppTable.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import Modal from '../../components/common/AppModal.vue'
import { Plus, Video, Copy, ExternalLink } from 'lucide-vue-next'
import { interviewService } from '../../services/interview.service'
import { jobService } from '../../services/job.service'
import { candidateService } from '../../services/candidate.service'
import { authStore } from '../../stores/auth.store'

const router = useRouter()
const interviews = ref([])
const loading = ref(true)
const copied = ref(false)
const routeMessage = ref(history.state?.message || '')
const showEnterRoomModal = ref(false)
const selectedInterview = ref(null)

const openEnterRoomModal = (row) => {
  selectedInterview.value = row
  showEnterRoomModal.value = true
}

const confirmEnterRoom = () => {
  showEnterRoomModal.value = false
  // Lưu interviewId vào state để trang Room biết vào phòng nào
  router.push({ path: '/recruiter-room', state: { message: 'Vào phòng phỏng vấn thành công!', interviewId: selectedInterview.value.id } })
}

onMounted(async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) throw new Error('Không tìm thấy company ID')
    
    const [response, jobs, candidates] = await Promise.all([
      interviewService.getInterviews(companyId),
      jobService.getJobs(companyId),
      candidateService.getCandidates(companyId)
    ])
    
    const jobMap = jobs.reduce((acc, j) => { acc[j.id] = j; return acc; }, {})
    const candidateMap = candidates.reduce((acc, c) => { acc[c.id] = c; return acc; }, {})

    interviews.value = response.map(i => {
      const candidate = i.candidate || candidateMap[i.candidate_id] || {}
      const job = i.job || jobMap[i.job_id?.String || i.job_id] || {}
      const dt = typeof i.scheduled_at === 'object' && i.scheduled_at !== null ? i.scheduled_at.Time : i.scheduled_at
      
      return {
        id: i.id,
        candidateName: candidate.full_name || candidate.name || 'Không rõ ứng viên',
        jobTitle: job.title || 'Không rõ vị trí',
        datetime: dt,
        status: i.status === 'scheduled' ? 'Scheduled' : i.status,
        link: i.invite_url || '' // Mock field, actual is in room object
      }
    })
  } catch (error) {
    console.error('Lỗi tải danh sách phỏng vấn:', error)
    // Could set a toast error here
  } finally {
    loading.value = false
  }
})

const copyLink = (link) => {
  navigator.clipboard.writeText(link)
  copied.value = true
  setTimeout(() => copied.value = false, 2000)
}

const columns = [
  { header: 'Vị trí & Ứng viên', key: 'candidate' },
  { header: 'Thời gian', key: 'time' },
  { header: 'Trạng thái', key: 'status' },
  { header: 'Phòng phỏng vấn', key: 'room' },
  { header: 'Hành động', key: 'action' }
]

const formatTime = (isoString) => {
  const date = new Date(isoString)
  return date.toLocaleTimeString('vi-VN', {hour: '2-digit', minute:'2-digit'})
}
const formatDate = (isoString) => {
  const date = new Date(isoString)
  return date.toLocaleDateString('vi-VN')
}
</script>

<template>
  <div>
    <Toast v-if="copied" type="success" message="Link phòng phỏng vấn đã được copy!" @close="copied = false" />
    <Toast v-if="routeMessage" type="success" :message="routeMessage" @close="routeMessage = ''" />

    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px">
      <div>
        <h1 class="text-h1">Lịch phỏng vấn</h1>
        <p class="text-helper" style="margin-top: 4px">Quản lý các buổi phỏng vấn trực tiếp với ứng viên.</p>
      </div>
      <Button @click="router.push('/interviews/new')"><Plus size="16" /> Tạo lịch phỏng vấn</Button>
    </div>

    <Card>
      <div v-if="loading" style="padding: 32px; text-align: center; color: var(--text-muted)">
        Đang tải lịch phỏng vấn...
      </div>
      <Table v-else :columns="columns" :data="interviews">
        <template #candidate="{ row }">
          <div>
            <div style="font-weight: 500; color: var(--text-main)">{{ row.candidateName }}</div>
            <div class="text-helper">Ứng tuyển: {{ row.jobTitle }}</div>
          </div>
        </template>
        <template #time="{ row }">
          <div>
            <div style="font-weight: 500; color: var(--text-main)">{{ formatTime(row.datetime) }}</div>
            <div class="text-helper">{{ formatDate(row.datetime) }}</div>
          </div>
        </template>
        <template #status="{ row }">
          <Badge :type="row.status === 'Scheduled' ? 'info' : 'neutral'">
            {{ row.status === 'Scheduled' ? 'Đã lên lịch' : row.status }}
          </Badge>
        </template>
        <template #room="{ row }">
          <div style="display: flex; align-items: center; gap: 8px">
            <a href="#" @click.prevent="openEnterRoomModal(row)" style="display: inline-flex; align-items: center; gap: 4px; color: var(--primary); font-weight: 500; font-size: 13px">
              <Video size="14" /> Vào phòng
            </a>
            <Button variant="ghost" @click="copyLink(row.link)" style="padding: 4px" title="Copy Link">
              <Copy size="14" />
            </Button>
          </div>
        </template>
        <template #action="{ row }">
          <Button variant="ghost" style="padding: 4px 8px; font-size: 13px" @click="router.push(`/interviews/${row.id}`)">
            <ExternalLink size="14" style="margin-right: 4px" /> Chi tiết
          </Button>
        </template>
      </Table>
    </Card>

    <Modal :isOpen="showEnterRoomModal" @close="showEnterRoomModal = false" title="Vào phòng phỏng vấn">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn tham gia phòng phỏng vấn này ngay bây giờ không? Camera và Microphone sẽ được kích hoạt.</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showEnterRoomModal = false">Hủy</Button>
        <Button variant="primary" @click="confirmEnterRoom">Vào phòng</Button>
      </div>
    </Modal>
  </div>
</template>
