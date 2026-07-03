<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { candidateService } from '../../services/candidate.service'
import { jobService } from '../../services/job.service'
import { authStore } from '../../stores/auth.store'
import { ArrowLeft, Save, Upload, FileText, Sparkles, CheckCircle } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const fileInputRef = ref(null)
const id = route.params.id
const isNew = ref(id === 'new')

const jobs = ref([])
const candidate = ref({ name: '', email: '', job_id: '', status: 'New' })
const loading = ref(!isNew.value)
const saving = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const parsedData = ref(null)
const localToast = ref(null)
const selectedFile = ref(null)
const cvPreviewUrl = ref(null)
const isAiParsing = ref(false)

onMounted(async () => {
  if (!isNew.value) {
    try {
      const companyId = authStore.user?.companies?.[0]?.id
      if (companyId) {
        const [data, jobsData] = await Promise.all([
          candidateService.getCandidate(companyId, id),
          jobService.getJobs(companyId)
        ])
        jobs.value = jobsData
        candidate.value = {
          ...data,
          name: data.full_name || data.name,
          job_id: data.latest_job?.id || '',
          status: data.status || 'new',
          cv: data.cv_file?.original_name || '',
          cv_file_id: data.cv_file?.id || ''
        }
        if (data.ai_cv_summary || data.parsed_cv_json) {
          parsedData.value = {
            skills: data.skills || [],
            experience: data.experience || data.ai_cv_summary || '',
            education: data.education || ''
          }
        }
      }
    } catch (error) {
      localToast.value = { type: 'error', message: 'Không thể tải chi tiết ứng viên' }
    } finally {
      loading.value = false
    }
  } else {
    // isNew
    try {
      const companyId = authStore.user?.companies?.[0]?.id
      if (companyId) {
        jobs.value = await jobService.getJobs(companyId)
      }
    } catch (err) {
      console.error(err)
    } finally {
      loading.value = false
    }
  }
})

const handleSave = async (e) => {
  e.preventDefault()
  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    
    let fileId = null
    if (selectedFile.value) {
      uploading.value = true
      uploadProgress.value = 50
      // Dùng endpoint CV upload, truyền 'new' làm candidate_id vì backend không check candidate_id lúc upload
      const fileResponse = await candidateService.uploadCV(companyId, id, selectedFile.value)
      fileId = fileResponse.id
      uploadProgress.value = 100
      uploading.value = false
    }

    let finalCandidateId = id
    if (isNew.value) {
      const createPayload = {
        full_name: candidate.value.name,
        email: candidate.value.email,
        phone: candidate.value.phone || '',
        job_id: candidate.value.job_id
      }
      if (fileId) createPayload.cv_file_id = fileId
      
      const newCandidate = await candidateService.createCandidate(companyId, createPayload)
      finalCandidateId = newCandidate.id
    } else {
      const updatePayload = {
        full_name: candidate.value.name,
        email: candidate.value.email,
        phone: candidate.value.phone || '',
        status: candidate.value.status
      }
      if (fileId) updatePayload.cv_file_id = fileId
      
      await candidateService.updateCandidate(companyId, id, updatePayload)
    }

    if (fileId) {
      isAiParsing.value = true
      await candidateService.parseCV(companyId, finalCandidateId)
      isAiParsing.value = false
      
      const updated = await candidateService.getCandidate(companyId, finalCandidateId)
      if (updated.ai_cv_summary) {
        parsedData.value = {
          skills: updated.skills || [],
          experience: updated.experience || updated.ai_cv_summary || '',
          education: updated.education || ''
        }
      }
    }
    
    localToast.value = { type: 'success', message: 'Lưu thông tin ứng viên thành công!' }
    setTimeout(() => {
      router.push('/candidates')
    }, 1500)
  } catch (error) {
    let msg = error?.message || 'Lưu ứng viên thất bại'
    if (msg === 'a candidate with this email already exists') {
      msg = 'Ứng viên với Email này đã tồn tại trong hệ thống!'
    }
    localToast.value = { type: 'error', message: msg }
  } finally {
    saving.value = false
    uploading.value = false
    isAiParsing.value = false
    if (cvPreviewUrl.value) {
      URL.revokeObjectURL(cvPreviewUrl.value)
      cvPreviewUrl.value = null
    }
    selectedFile.value = null
  }
}

const validateAndSetFile = (file) => {
  if (!file) return
  const validTypes = ['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 'image/png', 'image/jpeg', 'image/jpg']
  
  if (validTypes.includes(file.type)) {
    selectedFile.value = file
    if (cvPreviewUrl.value) URL.revokeObjectURL(cvPreviewUrl.value)
    
    if (file.type === 'application/pdf' || file.type.startsWith('image/')) {
      cvPreviewUrl.value = URL.createObjectURL(file)
    } else {
      cvPreviewUrl.value = null
    }
  } else {
    localToast.value = { type: 'error', message: 'Vui lòng chọn file PDF, Word hoặc Ảnh (PNG/JPG).' }
  }
}

const handleFileUpload = (e) => {
  validateAndSetFile(e.target.files[0])
}

const handleViewCV = async () => {
  if (!candidate.value.cv_file_id) return
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const res = await candidateService.getCVUrl(candidate.value.cv_file_id, companyId)
    if (res && res.url) {
      window.open(res.url, '_blank')
    } else {
      localToast.value = { type: 'error', message: 'Không lấy được đường dẫn CV' }
    }
  } catch (error) {
    localToast.value = { type: 'error', message: 'Lỗi khi xem CV: ' + (error.message || '') }
  }
}
</script>

<template>
  <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" style="position: fixed; top: 20px; right: 20px; z-index: 9999;" />
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
            minlength="2"
          />
          <Input 
            label="Email" 
            type="email"
            v-model="candidate.email" 
            required 
          />
          
          <div class="input-group">
            <label class="input-label">Vị trí ứng tuyển <span style="color:var(--danger)">*</span></label>
            <select class="input-field" v-model="candidate.job_id" required :disabled="!isNew">
              <option value="" disabled>-- Chọn vị trí ứng tuyển --</option>
              <option v-for="job in jobs" :key="job.id" :value="job.id">{{ job.title }}</option>
            </select>
            <p v-if="!isNew" class="text-helper" style="margin-top: 4px; font-size: 12px">Không thể thay đổi vị trí của ứng viên đã tạo</p>
          </div>
          
          <div class="input-group">
            <label class="input-label">Trạng thái</label>
            <select class="input-field" v-model="candidate.status">
              <option value="new">Mới (New)</option>
              <option value="interviewing">Đang phỏng vấn (Interviewing)</option>
              <option value="offered">Đã gửi Offer (Offered)</option>
              <option value="rejected">Từ chối (Rejected)</option>
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
          <div v-if="candidate.cv || selectedFile" style="border: 1px solid var(--border); border-radius: var(--radius-lg); padding: 16px; background-color: var(--surface-soft)">
            <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 12px">
              <FileText size="32" color="var(--primary)" />
              <div style="flex: 1">
                <p v-if="selectedFile" class="text-body" style="font-weight: 500">{{ selectedFile.name }}</p>
                <p v-else class="text-body" style="font-weight: 500; color: var(--primary); cursor: pointer; text-decoration: underline" @click="handleViewCV" title="Nhấn để xem CV">{{ candidate.cv }}</p>
                <p class="text-helper" style="color: var(--success); display: flex; align-items: center; gap: 4px; margin-top: 4px">
                  <CheckCircle size="14" /> {{ selectedFile ? 'Sẵn sàng tải lên' : 'Đã lưu trên hệ thống' }}
                </p>
              </div>
              <Button type="button" variant="secondary" @click="fileInputRef?.click()" :disabled="uploading || isAiParsing">Thay đổi</Button>
            </div>
            
            <div v-if="cvPreviewUrl" style="width: 100%; border: 1px solid var(--border); border-radius: var(--radius-md); overflow: hidden; margin-top: 16px; margin-bottom: 12px;">
              <img v-if="selectedFile && selectedFile.type.startsWith('image/')" :src="cvPreviewUrl" style="max-width: 100%; max-height: 400px; object-fit: contain; display: block; margin: 0 auto;" />
              <iframe v-else-if="selectedFile && selectedFile.type === 'application/pdf'" :src="cvPreviewUrl" width="100%" height="400px" style="border: none; display: block;"></iframe>
            </div>
            <div v-if="uploading" style="margin-top: 12px">
              <p class="text-body" style="margin-bottom: 8px; font-weight: 500; font-size: 13px">Đang tải lên... {{ uploadProgress }}%</p>
              <div style="height: 4px; background: var(--border); border-radius: 2px; overflow: hidden">
                <div :style="`height: 100%; background: var(--primary); width: ${uploadProgress}%; transition: width 0.3s`"></div>
              </div>
            </div>
            <div v-else-if="isAiParsing" style="margin-top: 16px; text-align: center; padding: 12px">
              <Sparkles size="24" color="var(--accent)" style="margin: 0 auto 8px; animation: pulse 1.5s infinite" />
              <p class="text-body" style="font-weight: 500; font-size: 14px">AI đang phân tích CV...</p>
            </div>
          </div>
          <div v-else 
            style="border: 2px dashed var(--border); border-radius: var(--radius-lg); padding: 40px 24px; text-align: center; cursor: pointer; background-color: var(--surface-soft)"
            @click="fileInputRef?.click()"
          >
            <div>
              <Upload size="32" color="var(--text-muted)" style="margin: 0 auto 12px" />
              <p class="text-body" style="font-weight: 500; margin-bottom: 4px">Click để tải CV lên</p>
              <p class="text-helper">Hỗ trợ PDF, DOCX, PNG, JPG (tối đa 5MB)</p>
            </div>
          </div>
          <input 
            type="file" 
            ref="fileInputRef" 
            style="display: none" 
            accept=".pdf,.doc,.docx,.png,.jpg,.jpeg" 
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
