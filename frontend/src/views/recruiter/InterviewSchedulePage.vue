<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import { ArrowLeft, Calendar, Save } from 'lucide-vue-next'
import { jobService } from '../../services/job.service'
import { candidateService } from '../../services/candidate.service'
import { interviewService } from '../../services/interview.service'
import { authStore } from '../../stores/auth.store'
import { isFutureDateTime, requiredTrim, validateForm } from '../../utils/validators.js'

const router = useRouter()
const jobs = ref([])
const candidates = ref([])

const interview = ref({
  jobId: '',
  candidateId: '',
  datetime: ''
})
const saving = ref(false)
const loading = ref(true)
const candidatesLoading = ref(false)
const localToast = ref(null)
const errors = ref({})
const minimumDateTime = computed(() => {
  const date = new Date(Date.now() + 15 * 60_000)
  const offset = date.getTimezoneOffset() * 60_000
  return new Date(date.getTime() - offset).toISOString().slice(0, 16)
})

watch(() => interview.value.jobId, async (jobId) => {
  interview.value.candidateId = ''
  candidates.value = []
  if (!jobId) return

  const companyId = authStore.user?.companies?.[0]?.id
  if (!companyId) return

  candidatesLoading.value = true
  try {
    const data = await candidateService.getCandidates(companyId, {
      job_id: jobId,
      page_size: 100,
    })
    candidates.value = Array.isArray(data) ? data : Array.isArray(data?.data) ? data.data : []
  } catch (error) {
    localToast.value = {
      type: 'error',
      message: 'Không thể tải ứng viên của công việc: ' + (error.message || 'Không xác định'),
    }
  } finally {
    candidatesLoading.value = false
  }
})

onMounted(async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) throw new Error('Không tìm thấy company ID')
    
    const jobsData = await jobService.getJobs(companyId)
    jobs.value = Array.isArray(jobsData) ? jobsData : []
  } catch (error) {
    localToast.value = { type: 'error', message: 'Lỗi tải danh sách: ' + (error.message || 'Không xác định') }
  } finally {
    loading.value = false
  }
})

const handleSave = async (e) => {
  e.preventDefault()
  if (saving.value) return

  const validation = validateForm(interview.value, {
    jobId: [
      (value) => requiredTrim(value, 'Vui lòng chọn công việc.'),
      (value) => jobs.value.some(job => job?.id === value) ? '' : 'Công việc đã chọn không hợp lệ.',
    ],
    candidateId: [
      (value) => requiredTrim(value, 'Vui lòng chọn ứng viên.'),
      (value) => candidates.value.some(candidate => candidate?.id === value) ? '' : 'Ứng viên không thuộc công việc đã chọn.',
    ],
    datetime: [
      (value) => requiredTrim(value, 'Vui lòng chọn ngày giờ phỏng vấn.'),
      (value) => isFutureDateTime(value, 15, 'Lịch phỏng vấn phải cách thời điểm hiện tại ít nhất 15 phút.'),
    ],
  })
  errors.value = validation.errors
  if (!validation.isValid) return

  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    
    // Đảm bảo định dạng chuẩn ISO 8601 UTC
    let date = new Date(interview.value.datetime);
    let isoString = date.toISOString();
    
    const payload = {
      job_id: interview.value.jobId,
      candidate_id: interview.value.candidateId,
      recruiter_id: authStore.user?.id,
      title: 'Phỏng vấn ứng viên',
      scheduled_at: isoString,
      duration_minutes: 60,
      mode: 'real',
      send_invite: true
    }
    
    await interviewService.createInterview(companyId, payload)
    router.push({ path: '/interviews', state: { message: 'Lên lịch phỏng vấn thành công!' } })
  } catch (error) {
    localToast.value = { type: 'error', message: 'Tạo lịch phỏng vấn thất bại: ' + (error.message || 'Lỗi không xác định') }
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div style="max-width: 700px; margin: 0 auto">
    <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" />

    <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 32px">
      <Button variant="ghost" @click="router.push('/interviews')" style="padding: 8px">
        <ArrowLeft size="20" />
      </Button>
      <div>
        <h1 class="text-h1">Lên lịch phỏng vấn</h1>
        <p class="text-helper" style="margin-top: 4px">Chọn ứng viên và thời gian để hệ thống tạo phòng phỏng vấn ảo.</p>
      </div>
    </div>

    <Card>
      <form @submit="handleSave">
        <div style="display: flex; flex-direction: column; gap: 20px">
          <div class="input-group">
            <label class="input-label">Công việc (Vị trí ứng tuyển)</label>
            <select :class="['input-field', { 'input-error': errors.jobId }]" required v-model="interview.jobId">
              <option value="">-- Chọn công việc --</option>
              <option v-for="j in jobs" :key="j.id" :value="j.id">{{ j.title }}</option>
            </select>
            <span v-if="errors.jobId" class="error-text">{{ errors.jobId }}</span>
          </div>

          <div class="input-group">
            <label class="input-label">Ứng viên</label>
            <select
              :class="['input-field', { 'input-error': errors.candidateId }]"
              required
              v-model="interview.candidateId"
              :disabled="!interview.jobId || candidatesLoading || candidates.length === 0"
            >
              <option value="">
                {{ !interview.jobId
                  ? '-- Vui lòng chọn công việc trước --'
                  : candidatesLoading
                    ? '-- Đang tải ứng viên --'
                    : candidates.length === 0
                      ? '-- Chưa có ứng viên ứng tuyển công việc này --'
                      : '-- Chọn ứng viên --' }}
              </option>
              <option v-for="c in candidates" :key="c.id" :value="c.id">{{ c.full_name || c.name }} ({{ c.email }})</option>
            </select>
            <span v-if="errors.candidateId" class="error-text">{{ errors.candidateId }}</span>
          </div>

          <div class="input-group">
            <label class="input-label">Ngày và giờ phỏng vấn</label>
            <input 
              type="datetime-local" 
              :class="['input-field', { 'input-error': errors.datetime }]"
              required
              :min="minimumDateTime"
              v-model="interview.datetime"
            />
            <span v-if="errors.datetime" class="error-text">{{ errors.datetime }}</span>
          </div>
          
          <div style="padding: 16px; background-color: rgba(37, 99, 235, 0.05); border: 1px solid rgba(37, 99, 235, 0.1); border-radius: 8px; display: flex; gap: 12px">
            <div>
              <p class="text-body" style="font-weight: 500; color: var(--primary); margin-bottom: 4px">Thông báo tự động</p>
              <p class="text-helper" style="color: var(--text-secondary)">
                Hệ thống sẽ tự động gửi email chứa link phòng phỏng vấn và lịch trình đến ứng viên ngay sau khi bạn bấm Lưu.
              </p>
            </div>
          </div>

          <div style="display: flex; justify-content: flex-end; margin-top: 16px; gap: 12px">
            <Button type="button" variant="ghost" @click="router.push('/interviews')">Hủy</Button>
            <Button type="submit" :disabled="saving">
              <Calendar size="16" /> {{ saving ? 'Đang tạo lịch...' : 'Lưu và gửi Email mời' }}
            </Button>
          </div>
        </div>
      </form>
    </Card>
  </div>
</template>
