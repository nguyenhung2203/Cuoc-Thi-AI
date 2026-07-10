<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { ArrowLeft, Video, Copy, Calendar, Clock, User, Briefcase, Mail, FileText, MessageSquare } from 'lucide-vue-next'
import { interviewService } from '../../services/interview.service'
import { jobService } from '../../services/job.service'
import { candidateService } from '../../services/candidate.service'
import { transcriptService } from '../../services/transcript.service'
import { authStore } from '../../stores/auth.store'

const route = useRoute()
const router = useRouter()
const id = route.params.id

const interview = ref(null)
const loading = ref(true)
const copied = ref(false)
const toast = ref(null)
const transcripts = ref([])
const loadingTranscripts = ref(false)

onMounted(async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (companyId) {
      const [data, jobs, candidates] = await Promise.all([
        interviewService.getInterview(companyId, id),
        jobService.getJobs(companyId),
        candidateService.getCandidates(companyId)
      ])
      
      const jobMap = jobs.reduce((acc, j) => { acc[j.id] = j; return acc; }, {})
      const candidateMap = candidates.reduce((acc, c) => { acc[c.id] = c; return acc; }, {})
      
      const candidate = data.candidate || candidateMap[data.candidate_id] || {}
      const job = data.job || jobMap[data.job_id?.String || data.job_id] || {}
      const dt = typeof data.scheduled_at === 'object' && data.scheduled_at !== null ? data.scheduled_at.Time : data.scheduled_at

      interview.value = {
        ...data,
        candidateName: candidate.full_name || candidate.name || 'Không rõ ứng viên',
        jobTitle: job.title || 'Không rõ vị trí',
        datetime: dt,
        status: data.status === 'scheduled' ? 'Scheduled' : data.status,
        link: data.invite_url || ''
      }
    }
  } catch (error) {
    toast.value = { type: 'error', message: 'Không thể tải thông tin phỏng vấn' }
  } finally {
    loading.value = false
  }
})

const copyLink = () => {
  if (interview.value && interview.value.link) {
    navigator.clipboard.writeText(interview.value.link)
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
  }
}

const handleSaveNotes = () => {
  toast.value = { type: 'success', message: 'Ghi chú đã được lưu thành công!' }
}

const dateObj = computed(() => {
  if (!interview.value || !interview.value.datetime) return null;
  const d = new Date(interview.value.datetime);
  return isNaN(d.getTime()) ? null : d;
})

const loadTranscripts = async () => {
  loadingTranscripts.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (companyId && interview.value) {
      const data = await transcriptService.getTranscripts(companyId, id)
      transcripts.value = (Array.isArray(data) ? data : (data?.data || [])).map(t => ({
        id: t.id,
        speaker: t.speaker_name || t.speaker_type || 'Unknown',
        speakerType: t.speaker_type || 'unknown',
        content: t.edited_content || t.content || '',
        time: t.created_at ? new Date(t.created_at).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit', second: '2-digit' }) : ''
      }))
    }
  } catch (error) {
    console.error('Failed to load transcripts:', error)
  } finally {
    loadingTranscripts.value = false
  }
}

// Load transcripts after interview data is ready
onMounted(async () => {
  // Wait for interview loading to finish, then load transcripts
  const checkAndLoad = setInterval(() => {
    if (!loading.value && interview.value) {
      clearInterval(checkAndLoad)
      loadTranscripts()
    }
  }, 200)
  // Timeout after 10s
  setTimeout(() => clearInterval(checkAndLoad), 10000)
})
</script>

<template>
  <div style="max-width: 1000px; margin: 0 auto; padding-bottom: 40px">
    <Toast v-if="copied" type="success" message="Link phòng phỏng vấn đã được copy!" @close="copied = false" />
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <div v-if="loading" style="padding: 32px; text-align: center; color: var(--text-muted)">Đang tải thông tin phỏng vấn...</div>
    <div v-else-if="!interview" style="padding: 32px; text-align: center; color: var(--warning)">Không tìm thấy thông tin phỏng vấn.</div>
    <div v-else>
      <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 32px">
        <Button variant="ghost" @click="router.push('/interviews')" style="padding: 8px">
          <ArrowLeft size="20" />
        </Button>
        <div>
          <div style="display: flex; align-items: center; gap: 12px">
            <h1 class="text-h1">Chi tiết buổi phỏng vấn</h1>
            <Badge :type="interview.status === 'Scheduled' ? 'info' : 'neutral'">
              {{ interview.status === 'Scheduled' ? 'Đã lên lịch' : interview.status }}
            </Badge>
          </div>
          <p class="text-helper" style="margin-top: 4px">Lịch phỏng vấn > {{ interview.candidateName }}</p>
        </div>
      </div>

      <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px">
        <div style="display: flex; flex-direction: column; gap: 24px">
          <Card title="Thông tin chung">
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px">
              <div>
                <p class="text-helper" style="margin-bottom: 4px; display: flex; align-items: center; gap: 6px"><User size="14" /> Ứng viên</p>
                <p class="text-body" style="font-weight: 500">{{ interview.candidateName }}</p>
              </div>
              <div>
                <p class="text-helper" style="margin-bottom: 4px; display: flex; align-items: center; gap: 6px"><Briefcase size="14" /> Vị trí ứng tuyển</p>
                <p class="text-body" style="font-weight: 500">{{ interview.jobTitle }}</p>
              </div>
              <div>
                <p class="text-helper" style="margin-bottom: 4px; display: flex; align-items: center; gap: 6px"><Calendar size="14" /> Ngày phỏng vấn</p>
                <p class="text-body" style="font-weight: 500">{{ dateObj.toLocaleDateString('vi-VN') }}</p>
              </div>
              <div>
                <p class="text-helper" style="margin-bottom: 4px; display: flex; align-items: center; gap: 6px"><Clock size="14" /> Thời gian</p>
                <p class="text-body" style="font-weight: 500">{{ dateObj.toLocaleTimeString('vi-VN', {hour: '2-digit', minute:'2-digit'}) }}</p>
              </div>
            </div>
            
            <div style="margin-top: 24px; padding-top: 24px; border-top: 1px solid var(--border)">
              <p class="text-helper" style="margin-bottom: 12px; font-weight: 500; color: var(--text-main)">Hành động</p>
              <div style="display: flex; gap: 12px">
                <Button v-if="interview.status !== 'Completed'" @click="router.push({ path: '/recruiter-room', state: { message: 'Vào phòng phỏng vấn thành công!', interviewId: interview.id } })">
                  <Video size="16" /> Vào phòng phỏng vấn
                </Button>
                <Button v-if="interview.status === 'Completed'" @click="router.push(`/interviews/${interview.id}/report`)" variant="primary" style="background-color: var(--accent); border-color: var(--accent); color: white">
                  <FileText size="16" /> Xem Báo cáo AI
                </Button>
                <Button variant="secondary" @click="copyLink">
                  <Copy size="16" /> Copy Link Invite
                </Button>
                <Button variant="ghost">
                  <Mail size="16" /> Gửi email nhắc nhở
                </Button>
              </div>
            </div>
          </Card>
          
          <Card title="Cấu hình AI Assistant" style="border-top: 4px solid var(--accent)">
            <ul class="text-body" style="padding-left: 20px; display: flex; flex-direction: column; gap: 12px">
              <li><strong>Phân tích realtime:</strong> Bật</li>
              <li><strong>Tạo transcript:</strong> Bật</li>
              <li><strong>Gợi ý câu hỏi (Rubric-based):</strong> Bật</li>
              <li><strong>Tự động chấm điểm:</strong> Bật</li>
            </ul>
          </Card>
        </div>

        <div>
          <Card title="Ghi chú nội bộ">
            <textarea 
              class="input-field" 
              rows="6"
              placeholder="Nhập ghi chú hoặc nhắc nhở trước buổi phỏng vấn (Chỉ recruiter xem được)..."
              style="width: 100%; margin-bottom: 16px"
            ></textarea>
            <Button variant="secondary" style="width: 100%" @click="handleSaveNotes">Lưu ghi chú</Button>
          </Card>
        </div>
      </div>

      <!-- Transcript Section -->
      <div style="margin-top: 24px">
        <Card>
          <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 20px">
            <MessageSquare size="20" style="color: var(--accent)" />
            <h3 class="text-h3" style="margin: 0">Nội dung cuộc trò chuyện (Transcript)</h3>
            <button @click="loadTranscripts" class="ml-auto text-xs text-indigo-600 hover:underline font-medium">Làm mới</button>
          </div>

          <div v-if="loadingTranscripts" style="padding: 24px 0; text-align: center; color: var(--text-muted)">
            <div class="w-6 h-6 border-3 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mx-auto mb-3"></div>
            Đang tải transcript...
          </div>

          <div v-else-if="transcripts.length === 0" style="padding: 32px 0; text-align: center; color: var(--text-muted)">
            <MessageSquare size="40" style="margin: 0 auto 12px; opacity: 0.4" />
            <p>Chưa có transcript nào được ghi lại cho buổi phỏng vấn này.</p>
          </div>

          <div v-else class="space-y-3" style="max-height: 500px; overflow-y: auto; padding-right: 8px">
            <div v-for="t in transcripts" :key="t.id"
              class="flex gap-3 p-3 rounded-lg transition-colors"
              :class="t.speakerType === 'interviewer' ? 'bg-indigo-50 dark:bg-indigo-500/10' : 'bg-slate-50 dark:bg-slate-800'"
            >
              <div class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold shrink-0"
                :class="t.speakerType === 'interviewer' ? 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300' : 'bg-slate-200 text-slate-700 dark:bg-slate-700 dark:text-slate-300'"
              >
                {{ t.speaker.charAt(0).toUpperCase() }}
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-xs font-semibold" :class="t.speakerType === 'interviewer' ? 'text-indigo-700 dark:text-indigo-300' : 'text-slate-700 dark:text-slate-300'">{{ t.speaker }}</span>
                  <span class="text-[11px] text-slate-400">{{ t.time }}</span>
                </div>
                <p class="text-sm text-slate-700 dark:text-slate-300 leading-relaxed">{{ t.content }}</p>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>
