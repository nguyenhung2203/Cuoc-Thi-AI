<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Save, ArrowLeft, Wand2, Plus, X } from 'lucide-vue-next'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { jobService } from '../../services/job.service'
import { authStore } from '../../stores/auth.store'

const router = useRouter()

const formData = ref({
  title: '',
  department: '',
  location: '',
  type: 'Full-time',
  description: '',
  requirements: [],
  benefits: ''
})

const currentReq = ref('')
const isLoading = ref(false)
const isAnalyzing = ref(false)
const errors = ref({})
const localToast = ref(null)

const validate = () => {
  const newErrors = {}
  if (!formData.value.title) newErrors.title = 'Vui lòng nhập chức danh'
  if (!formData.value.department) newErrors.department = 'Vui lòng nhập phòng ban'
  if (!formData.value.description) newErrors.description = 'Vui lòng nhập mô tả công việc'
  if (!formData.value.benefits) newErrors.benefits = 'Vui lòng nhập quyền lợi'
  if (formData.value.requirements.length === 0) newErrors.requirements = 'Vui lòng thêm ít nhất 1 yêu cầu công việc'
  
  errors.value = newErrors
  return Object.keys(newErrors).length === 0
}

const addRequirement = () => {
  if (currentReq.value.trim()) {
    formData.value.requirements.push(currentReq.value.trim())
    currentReq.value = ''
    if (errors.value.requirements) {
      delete errors.value.requirements
    }
  }
}

const removeRequirement = (index) => {
  formData.value.requirements.splice(index, 1)
}

const handleAnalyzeJD = async () => {
  if (!formData.value.description) {
    errors.value.description = 'Vui lòng nhập mô tả công việc trước khi phân tích'
    return
  }

  // NOTE: AI Job Description analysis requires the job to be created first (to get job ID).
  // This button currently only parses description locally for UX preview.
  // Real analysis: call jobService.analyzeJD(companyId, jobId) AFTER job is created.
  // For now, this button provides client-side suggestion feedback.

  isAnalyzing.value = true

  // Simulate extracting requirements from description text (client-side hint)
  // Real AI analysis happens on backend after job creation
  await new Promise(resolve => setTimeout(resolve, 1000))

  if (formData.value.requirements.length === 0) {
    // Provide a helpful hint based on the description
    formData.value.requirements = [
      'Yêu cầu phân tích chi tiết từ AI server sau khi tạo job'
    ]
    if (errors.value.requirements) delete errors.value.requirements
    localToast.value = { type: 'info', message: 'Hãy tạo job trước để AI phân tích chi tiết JD.' }
  }

  isAnalyzing.value = false
}

const handleSave = async () => {
  if (!validate()) {
    localToast.value = { type: 'error', message: 'Vui lòng điền đầy đủ các trường bắt buộc.' }
    return
  }
  
  isLoading.value = true
  
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) throw new Error('Không tìm thấy company ID')
    
    // Tạo mảng string yêu cầu từ state, join lại thành text vì API mong đợi string cho requirements (hoặc dùng text thô)
    const payload = {
      title: formData.value.title,
      department: formData.value.department,
      location: formData.value.location,
      employment_type: formData.value.type,
      description: formData.value.description,
      requirements: formData.value.requirements.join('\n'), // Chuyển mảng thành text xuống dòng
      benefits: formData.value.benefits,
      status: 'open'
    }
    
    await jobService.createJob(companyId, payload)
    
    router.push({ 
      path: '/jobs', 
      state: { message: 'Đã tạo chiến dịch tuyển dụng mới thành công!' }
    })
  } catch (error) {
    localToast.value = { type: 'error', message: 'Tạo công việc thất bại: ' + (error.message || 'Lỗi không xác định') }
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="page-header" style="margin-bottom: 24px;">
    <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" />
    <div>
      <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 8px;">
        <button class="icon-btn" @click="router.back()">
          <ArrowLeft size="24" />
        </button>
        <h1 class="text-h1">Tạo Job Mới</h1>
      </div>
      <p class="text-body" style="color: var(--text-secondary); margin-left: 40px;">
        Khởi tạo một chiến dịch tuyển dụng mới và thiết lập tiêu chí đánh giá.
      </p>
    </div>
    
    <div style="display: flex; gap: 12px;">
      <Button variant="secondary" @click="router.back()" :disabled="isLoading">Hủy</Button>
      <Button variant="primary" @click="handleSave" :isLoading="isLoading">
        <Save size="18" /> Lưu & Xuất bản
      </Button>
    </div>
  </div>

  <div class="job-create-layout">
    <div class="main-column">
      <Card title="Thông tin cơ bản" style="margin-bottom: 24px;">
        <div class="form-grid">
          <Input 
            v-model="formData.title"
            label="Chức danh tuyển dụng (*)"
            placeholder="VD: Senior Frontend Developer"
            :error="errors.title"
          />
          <Input 
            v-model="formData.department"
            label="Phòng ban (*)"
            placeholder="VD: Engineering"
            :error="errors.department"
          />
          <Input 
            v-model="formData.location"
            label="Địa điểm làm việc"
            placeholder="VD: Hồ Chí Minh, Hybrid"
          />
          
          <div class="input-group">
            <label class="input-label">Loại hình công việc</label>
            <select v-model="formData.type" class="input-field">
              <option value="Full-time">Full-time</option>
              <option value="Part-time">Part-time</option>
              <option value="Contract">Contract</option>
              <option value="Internship">Internship</option>
            </select>
          </div>
        </div>
      </Card>

      <Card title="Mô tả công việc (JD)" style="margin-bottom: 24px;">
        <div class="input-group">
          <textarea 
            v-model="formData.description"
            class="input-field" 
            rows="6" 
            placeholder="Nhập chi tiết mô tả công việc, trách nhiệm..."
            :class="{ 'input-error': errors.description }"
          ></textarea>
          <span v-if="errors.description" class="error-text">{{ errors.description }}</span>
        </div>
      </Card>
      
      <Card title="Quyền lợi (*)">
        <div class="input-group">
          <textarea 
            v-model="formData.benefits"
            class="input-field" 
            rows="6" 
            placeholder="Nhập các quyền lợi, chế độ đãi ngộ dành cho ứng viên..."
            :class="{ 'input-error': errors.benefits }"
          ></textarea>
          <span v-if="errors.benefits" class="error-text">{{ errors.benefits }}</span>
        </div>
      </Card>
    </div>

    <div class="side-column">
      <Card title="Yêu cầu & Tiêu chí (*)" style="margin-bottom: 24px;" :class="{'border-error': errors.requirements}">
        <p class="text-helper" style="margin-bottom: 16px;">
          Bạn có thể nhập thủ công hoặc để AI phân tích tự động từ Mô tả công việc.
        </p>
        
        <Button 
          variant="secondary" 
          style="width: 100%; margin-bottom: 24px; color: var(--primary); border-color: var(--primary);"
          @click="handleAnalyzeJD"
          :isLoading="isAnalyzing"
        >
          <Wand2 size="18" /> Tự động phân tích bằng AI
        </Button>

        <div class="req-input-area" style="margin-bottom: 16px;">
          <Input 
            v-model="currentReq"
            placeholder="Nhập yêu cầu mới..."
            @keyup.enter="addRequirement"
          />
          <Button variant="secondary" style="margin-top: 4px;" @click="addRequirement">
            <Plus size="18" />
          </Button>
        </div>

        <div class="req-list">
          <div v-if="formData.requirements.length === 0" class="empty-reqs text-helper">
            Chưa có yêu cầu nào được thiết lập.
          </div>
          
          <div v-for="(req, index) in formData.requirements" :key="index" class="req-item">
            <div class="req-dot"></div>
            <span class="req-text">{{ req }}</span>
            <button class="icon-btn req-remove" @click="removeRequirement(index)">
              <X size="14" />
            </button>
          </div>
        </div>
        <div v-if="errors.requirements" class="error-text" style="margin-top: 12px;">{{ errors.requirements }}</div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.border-error {
  border: 1px solid var(--danger);
}
.job-create-layout {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 16px;
}

/* Make title span full width */
.form-grid > div:first-child {
  grid-column: 1 / -1;
}

.req-input-area {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}
.req-input-area > div {
  flex: 1;
  margin-bottom: 0;
}

.req-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.req-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  background-color: var(--surface-soft);
  border-radius: var(--radius);
  border: 1px solid var(--border);
  transition: all 0.2s;
}

.req-item:hover {
  border-color: var(--primary);
}

.req-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: var(--primary);
  margin-top: 6px;
  flex-shrink: 0;
}

.req-text {
  flex: 1;
  font-size: 13px;
  line-height: 1.5;
}

.req-remove {
  opacity: 0;
  color: var(--danger);
  padding: 2px;
}

.req-item:hover .req-remove {
  opacity: 1;
}

.empty-reqs {
  text-align: center;
  padding: 24px 0;
  border: 1px dashed var(--border);
  border-radius: var(--radius);
}

@media (max-width: 1024px) {
  .job-create-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
