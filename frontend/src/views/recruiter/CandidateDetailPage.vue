<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import { mockApi } from '../../utils/mockData'
import { ArrowLeft, Save, Upload, FileText, Sparkles, CheckCircle } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const fileInputRef = ref(null)
const id = route.params.id
const isNew = id === 'new'

const candidate = ref({ name: '', email: '', appliedJob: '', status: 'New' })
const loading = ref(!isNew)
const saving = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const parsedData = ref(null)

onMounted(async () => {
  if (!isNew) {
    const data = await mockApi.candidates.getById(id)
    if (data) candidate.value = data
    loading.value = false
  }
})

const handleSave = async (e) => {
  e.preventDefault()
  saving.value = true
  if (isNew) {
    await mockApi.candidates.create(candidate.value)
  } else {
    await new Promise(r => setTimeout(r, 600))
  }
  saving.value = false
  router.push({ path: '/candidates', state: { message: 'Lưu thông tin ứng viên thành công!' } })
}

const handleFileUpload = (e) => {
  const file = e.target.files[0]
  if (!file) return

  uploading.value = true
  uploadProgress.value = 0

  const interval = setInterval(() => {
    uploadProgress.value += 25
    if (uploadProgress.value >= 100) {
      clearInterval(interval)
      uploading.value = false
      candidate.value.cv = file.name
      
      setTimeout(() => {
        parsedData.value = {
          skills: ['React', 'JavaScript', 'Node.js', 'TypeScript', 'CSS'],
          experience: '3 years Frontend',
          education: 'BS Computer Science'
        }
      }, 1000)
    }
  }, 400)
}
</script>

<template>
  <div v-if="loading">Đang tải...</div>
  <div v-else>
    <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 32px">
      <Button variant="ghost" @click="router.push('/candidates')" style="padding: 8px">
        <ArrowLeft size="20" />
      </Button>
      <div>
        <div style="display: flex; align-items: center; gap: 12px">
          <h1 class="text-h1">{{ isNew ? 'Thêm Ứng Viên Mới' : candidate.name }}</h1>
          <Badge v-if="!isNew" type="info">{{ candidate.status }}</Badge>
        </div>
        <p class="text-helper" style="margin-top: 4px">Candidates > {{ isNew ? 'New' : candidate.name }}</p>
      </div>
    </div>

    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px">
      <Card title="Thông tin cá nhân">
        <form @submit="handleSave">
          <Input 
            label="Họ và Tên" 
            v-model="candidate.name" 
            required 
          />
          <Input 
            label="Email" 
            type="email"
            v-model="candidate.email" 
            required 
          />
          <Input 
            label="Vị trí ứng tuyển" 
            v-model="candidate.appliedJob" 
            required 
          />
          
          <div class="input-group">
            <label class="input-label">Trạng thái</label>
            <select class="input-field" v-model="candidate.status">
              <option value="New">Mới (New)</option>
              <option value="Interviewing">Đang phỏng vấn (Interviewing)</option>
              <option value="Offered">Đã gửi Offer (Offered)</option>
              <option value="Rejected">Từ chối (Rejected)</option>
            </select>
          </div>

          <div style="display: flex; justify-content: flex-end; margin-top: 24px; gap: 12px">
            <Button type="button" variant="ghost" @click="router.push('/candidates')">Hủy</Button>
            <Button type="submit" :disabled="saving">
              <Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu ứng viên' }}
            </Button>
          </div>
        </form>
      </Card>

      <div style="display: flex; flex-direction: column; gap: 24px">
        <Card title="Hồ sơ (CV & Resume)">
          <div v-if="candidate.cv" style="display: flex; align-items: center; gap: 16px; padding: 16px; background-color: var(--surface-soft); border-radius: var(--radius)">
            <FileText size="32" color="var(--primary)" />
            <div style="flex: 1">
              <p class="text-body" style="font-weight: 500">{{ candidate.cv }}</p>
              <p class="text-helper" style="color: var(--success); display: flex; align-items: center; gap: 4px; margin-top: 4px">
                <CheckCircle size="14" /> Tải lên thành công
              </p>
            </div>
            <Button variant="secondary" @click="fileInputRef?.click()">Thay đổi</Button>
          </div>
          <div v-else 
            style="border: 2px dashed var(--border); border-radius: var(--radius-lg); padding: 40px 24px; text-align: center; cursor: pointer; background-color: var(--surface-soft)"
            @click="fileInputRef?.click()"
          >
            <div v-if="uploading">
              <p class="text-body" style="margin-bottom: 12px; font-weight: 500">Đang tải lên... {{ uploadProgress }}%</p>
              <div style="height: 6px; background: var(--border); border-radius: 3px; overflow: hidden">
                <div :style="`height: 100%; background: var(--primary); width: ${uploadProgress}%; transition: width 0.3s`"></div>
              </div>
            </div>
            <div v-else>
              <Upload size="32" color="var(--text-muted)" style="margin: 0 auto 12px" />
              <p class="text-body" style="font-weight: 500; margin-bottom: 4px">Click để tải CV lên</p>
              <p class="text-helper">Hỗ trợ PDF, DOCX tối đa 5MB</p>
            </div>
          </div>
          <input 
            type="file" 
            ref="fileInputRef" 
            style="display: none" 
            accept=".pdf,.doc,.docx" 
            @change="handleFileUpload"
          />
        </Card>

        <Card v-if="parsedData" title="AI Bóc tách dữ liệu CV" style="border-top: 4px solid var(--accent)">
          <div style="display: flex; flex-direction: column; gap: 16px">
            <div style="display: flex; align-items: center; gap: 8px; color: var(--accent)">
              <Sparkles size="16" /> <span style="font-size: 14px; font-weight: 500">Hoàn tất phân tích tự động</span>
            </div>
            <div>
              <h4 class="text-helper" style="text-transform: uppercase; margin-bottom: 8px; letter-spacing: 0.05em">Kỹ năng nổi bật</h4>
              <div style="display: flex; gap: 8px; flex-wrap: wrap">
                <span v-for="s in parsedData.skills" :key="s" style="background-color: var(--surface-soft); padding: 4px 10px; border-radius: 6px; font-size: 13px; border: 1px solid var(--border)">
                  {{ s }}
                </span>
              </div>
            </div>
            <div>
              <h4 class="text-helper" style="text-transform: uppercase; margin-bottom: 4px; letter-spacing: 0.05em">Kinh nghiệm</h4>
              <p class="text-body">{{ parsedData.experience }}</p>
            </div>
            <div>
              <h4 class="text-helper" style="text-transform: uppercase; margin-bottom: 4px; letter-spacing: 0.05em">Học vấn</h4>
              <p class="text-body">{{ parsedData.education }}</p>
            </div>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>
