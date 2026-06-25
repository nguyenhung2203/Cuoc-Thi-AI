<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import { ArrowLeft, Calendar, Mail, AlertCircle } from 'lucide-vue-next'
import { mockApi } from '../../utils/mockData'

const router = useRouter()
const jobs = ref([])
const candidates = ref([])

const interview = ref({
  jobTitle: '',
  candidateName: '',
  datetime: ''
})
const saving = ref(false)

onMounted(async () => {
  const jobsData = await mockApi.jobs.getAll()
  jobs.value = jobsData
  const candidatesData = await mockApi.candidates.getAll()
  candidates.value = candidatesData
})

const handleSave = async (e) => {
  e.preventDefault()
  if (!interview.value.jobTitle || !interview.value.candidateName || !interview.value.datetime) return
  
  saving.value = true
  await mockApi.interviews.create(interview.value)
  saving.value = false
  router.push({ path: '/interviews', state: { message: 'Lên lịch phỏng vấn thành công!' } })
}
</script>

<template>
  <div style="max-width: 700px; margin: 0 auto">
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
            <select class="input-field" required v-model="interview.jobTitle">
              <option value="">-- Chọn công việc --</option>
              <option v-for="j in jobs" :key="j.id" :value="j.title">{{ j.title }}</option>
            </select>
          </div>

          <div class="input-group">
            <label class="input-label">Ứng viên</label>
            <select class="input-field" required v-model="interview.candidateName">
              <option value="">-- Chọn ứng viên --</option>
              <option v-for="c in candidates" :key="c.id" :value="c.name">{{ c.name }} ({{ c.email }})</option>
            </select>
          </div>

          <div class="input-group">
            <label class="input-label">Ngày và giờ phỏng vấn</label>
            <input 
              type="datetime-local" 
              class="input-field" 
              required
              v-model="interview.datetime"
            />
          </div>
          
          <div style="padding: 16px; background-color: rgba(37, 99, 235, 0.05); border: 1px solid rgba(37, 99, 235, 0.1); border-radius: 8px; display: flex; gap: 12px">
            <AlertCircle size="20" color="var(--primary)" style="flex-shrink: 0" />
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
