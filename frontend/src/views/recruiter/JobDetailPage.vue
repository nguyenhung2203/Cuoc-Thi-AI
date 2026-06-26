<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import { mockApi } from '../../utils/mockData'
import { Sparkles, ArrowLeft, Save, AlertCircle, CheckCircle } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const id = route.params.id
const isNew = id === 'new'

const job = ref({ title: '', status: 'Open', description: '' })
const loading = ref(!isNew)
const saving = ref(false)
const aiAnalyzing = ref(false)
const rubric = ref(null)

onMounted(async () => {
  if (!isNew) {
    const data = await mockApi.jobs.getById(id)
    if (data) job.value = data
    loading.value = false
  }
})

const handleSave = async (e) => {
  e.preventDefault()
  saving.value = true
  if (isNew) {
    await mockApi.jobs.create({ ...job.value, created: new Date().toISOString().split('T')[0] })
  } else {
    await new Promise(r => setTimeout(r, 600))
  }
  saving.value = false
  router.push({ path: '/jobs', state: { message: 'Lưu thông tin công việc thành công!' } })
}

const handleAiAnalyze = () => {
  aiAnalyzing.value = true
  setTimeout(() => {
    aiAnalyzing.value = false
    rubric.value = [
      { criterion: 'Technical Skills', weight: '40%' },
      { criterion: 'Communication', weight: '30%' },
      { criterion: 'Problem Solving', weight: '30%' }
    ]
  }, 1500)
}
</script>

<template>
  <div v-if="loading">Đang tải...</div>
  <div v-else>
    <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 32px">
      <Button variant="ghost" @click="router.push('/jobs')" style="padding: 8px">
        <ArrowLeft size="20" />
      </Button>
      <div>
        <div style="display: flex; align-items: center; gap: 12px">
          <h1 class="text-h1">{{ isNew ? 'Tạo Job Mới' : job.title }}</h1>
          <Badge v-if="!isNew" :type="job.status === 'Open' ? 'success' : 'neutral'">{{ job.status }}</Badge>
        </div>
        <p class="text-helper" style="margin-top: 4px">Jobs > {{ isNew ? 'New' : job.title }}</p>
      </div>
    </div>

    <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px">
      <div style="display: flex; flex-direction: column; gap: 24px">
        <Card title="Thông tin chung">
          <form @submit="handleSave">
            <Input 
              label="Tiêu đề công việc" 
              v-model="job.title" 
              required 
            />
            
            <div class="input-group">
              <label class="input-label">Trạng thái</label>
              <select class="input-field" v-model="job.status">
                <option value="Open">Đang mở (Open)</option>
                <option value="Closed">Đã đóng (Closed)</option>
              </select>
            </div>

            <div class="input-group">
              <label class="input-label">Mô tả công việc (JD)</label>
              <textarea 
                class="input-field" 
                rows="10"
                v-model="job.description"
                placeholder="Nhập yêu cầu công việc..."
              ></textarea>
            </div>

            <div style="display: flex; justify-content: flex-end; margin-top: 24px; gap: 12px">
              <Button type="button" variant="ghost" @click="router.push('/jobs')">Hủy</Button>
              <Button type="submit" :disabled="saving">
                <Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu thông tin' }}
              </Button>
            </div>
          </form>
        </Card>
        
        <Card v-if="!isNew" title="Danh sách ứng viên">
          <div style="text-align: center; color: var(--text-muted); padding: 32px 0">
            <p class="text-body">Chưa có ứng viên nào nộp đơn.</p>
          </div>
        </Card>
      </div>

      <div style="display: flex; flex-direction: column; gap: 24px">
        <Card title="AI Gợi ý tối ưu JD" style="border-top: 4px solid var(--accent)">
          <p class="text-body" style="color: var(--text-muted); margin-bottom: 16px">
            Hệ thống AI sẽ tự động đọc JD và phân tích ra bộ tiêu chí chấm điểm (Rubric) phù hợp nhất.
          </p>
          <Button 
            variant="secondary" 
            style="width: 100%; border-color: var(--accent); color: var(--accent)" 
            @click="handleAiAnalyze"
            :disabled="aiAnalyzing || !job.description"
          >
            <Sparkles size="16" /> 
            {{ aiAnalyzing ? 'Đang phân tích...' : 'Bóc tách Rubric bằng AI' }}
          </Button>

          <div v-if="rubric" style="margin-top: 24px; border-top: 1px solid var(--border); padding-top: 16px">
            <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px">
              <CheckCircle size="16" color="var(--success)" />
              <h4 class="text-h2" style="font-size: 14px">Bộ tiêu chí đề xuất</h4>
            </div>
            <ul class="text-body" style="padding-left: 24px; display: flex; flex-direction: column; gap: 8px">
              <li v-for="(r, i) in rubric" :key="i">
                {{ r.criterion }} - <span style="color: var(--primary); font-weight: 600">{{ r.weight }}</span>
              </li>
            </ul>
          </div>
          <div v-if="!job.description && !rubric" style="margin-top: 16px; display: flex; align-items: center; gap: 8px; color: var(--warning); font-size: 13px">
            <AlertCircle size="14" /> Cần nhập JD để AI có thể phân tích.
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>
