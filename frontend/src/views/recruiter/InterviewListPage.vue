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
  <div class="animate-fade-in space-y-6">
    <Toast v-if="copied" type="success" message="Link phòng phỏng vấn đã được copy!" @close="copied = false" />
    <Toast v-if="routeMessage" type="success" :message="routeMessage" @close="routeMessage = ''" />

    <!-- Page Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Lịch phỏng vấn</h1>
        <p class="text-slate-500 dark:text-slate-400 mt-1 text-sm">Quản lý các buổi phỏng vấn trực tiếp với ứng viên.</p>
      </div>
      <Button @click="router.push('/interviews/new')" class="bg-indigo-600 hover:bg-indigo-700 text-white border-none shadow-md shadow-indigo-500/20">
        <Plus size="16" class="mr-1" /> Tạo lịch phỏng vấn
      </Button>
    </div>

    <!-- Main Content Card -->
    <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 overflow-hidden">
      <!-- Table Section -->
      <div v-if="loading" class="p-12 text-center text-slate-500 dark:text-slate-400 flex flex-col items-center justify-center gap-3">
        <div class="w-8 h-8 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
        <span class="text-sm font-medium">Đang tải lịch phỏng vấn...</span>
      </div>
      
      <div v-else class="w-full overflow-x-auto">
        <Table :columns="columns" :data="interviews" class="w-full text-left text-sm text-slate-600 dark:text-slate-400">
          <template #candidate="{ row }">
            <div class="flex flex-col">
              <div class="font-semibold text-slate-800 dark:text-slate-200">{{ row.candidateName }}</div>
              <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 font-medium">Ứng tuyển: <span class="text-indigo-600 dark:text-indigo-400">{{ row.jobTitle }}</span></div>
            </div>
          </template>
          
          <template #time="{ row }">
            <div class="flex flex-col">
              <div class="font-semibold text-slate-800 dark:text-slate-200">{{ formatTime(row.datetime) }}</div>
              <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{{ formatDate(row.datetime) }}</div>
            </div>
          </template>
          
          <template #status="{ row }">
            <span 
              class="px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider rounded-full border"
              :class="[
                row.status === 'Scheduled' ? 'bg-indigo-50 text-indigo-600 border-indigo-200 dark:bg-indigo-500/10 dark:text-indigo-400 dark:border-indigo-500/20' : 
                'bg-slate-100 text-slate-600 border-slate-200 dark:bg-slate-700 dark:text-slate-300 dark:border-slate-600'
              ]"
            >
              {{ row.status === 'Scheduled' ? 'Đã lên lịch' : row.status }}
            </span>
          </template>
          
          <template #room="{ row }">
            <div class="flex items-center gap-2">
              <button @click.prevent="openEnterRoomModal(row)" class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-indigo-50 hover:bg-indigo-100 dark:bg-indigo-500/10 dark:hover:bg-indigo-500/20 text-indigo-600 dark:text-indigo-400 font-semibold text-xs rounded-lg transition-colors border border-indigo-100 dark:border-indigo-500/20">
                <Video size="14" /> Vào phòng
              </button>
              <button @click="copyLink(row.link)" class="p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-md transition-colors" title="Copy Link">
                <Copy size="14" />
              </button>
            </div>
          </template>
          
          <template #action="{ row }">
            <button @click="router.push(`/interviews/${row.id}`)" class="inline-flex items-center gap-1 px-3 py-1.5 text-slate-500 hover:text-indigo-600 dark:text-slate-400 dark:hover:text-indigo-400 hover:bg-slate-50 dark:hover:bg-slate-800 rounded-lg transition-colors font-medium text-xs">
              <ExternalLink size="14" /> Chi tiết
            </button>
          </template>
        </Table>
      </div>
    </div>

    <!-- Modals -->
    <Modal :isOpen="showEnterRoomModal" @close="showEnterRoomModal = false" title="Vào phòng phỏng vấn">
      <div class="p-1">
        <p class="text-slate-600 dark:text-slate-300 mb-6 leading-relaxed">Bạn có chắc chắn muốn tham gia phòng phỏng vấn này ngay bây giờ không? Camera và Microphone sẽ được kích hoạt.</p>
        <div class="flex justify-end gap-3">
          <Button variant="ghost" @click="showEnterRoomModal = false" class="text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700">Hủy</Button>
          <Button variant="primary" @click="confirmEnterRoom" class="bg-indigo-600 hover:bg-indigo-700 text-white border-none shadow-md shadow-indigo-500/20">Vào phòng</Button>
        </div>
      </div>
    </Modal>
  </div>
</template>
