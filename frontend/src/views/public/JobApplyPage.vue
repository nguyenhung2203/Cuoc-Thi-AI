<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { publicService } from '../../services/public.service'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { authStore } from '../../stores/auth.store'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { Briefcase, MapPin, Clock, ArrowLeft, UploadCloud } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const companyId = route.params.company_id
const jobId = route.params.job_id
const job = ref(null)
const loading = ref(true)
const submitting = ref(false)
const toast = ref(null)

const cvFile = ref(null)
const cvPreviewUrl = ref(null)
const isDragging = ref(false)
const fileInput = ref(null)
const hasApplied = ref(false)

const isLoggedIn = computed(() => authStore.isAuthenticated)
const user = computed(() => authStore.user)

const unwrap = (val) => {
  if (!val) return ''
  if (typeof val === 'object' && 'String' in val) {
    return val.Valid ? val.String : ''
  }
  return val
}

onMounted(async () => {
  try {
    // If auth store exists but not initialized, let's init it
    if (authStore.isAuthenticated && !authStore.user) {
      await authStore.init()
    }

    if (authStore.isAuthenticated) {
      try {
        const checkRes = await candidatePortalService.checkApplied(jobId)
        hasApplied.value = checkRes.data?.has_applied || checkRes.has_applied || false
      } catch (err) {
        console.error('Failed to check if applied', err)
      }
    }

    const res = await publicService.getJobDetails(companyId, jobId)
    const rawJob = res.data || res
    job.value = {
      ...rawJob,
      location: unwrap(rawJob.location),
      employment_type: unwrap(rawJob.employment_type),
      department: unwrap(rawJob.department),
      requirements: unwrap(rawJob.requirements),
      benefits: unwrap(rawJob.benefits)
    }
  } catch (error) {
    console.error('Failed to load job', error)
    toast.value = { type: 'error', message: 'Không tìm thấy thông tin công việc.' }
  } finally {
    loading.value = false
  }
})

const validateAndSetFile = (file) => {
  if (!file) return
  const validTypes = ['application/pdf', 'image/png', 'image/jpeg', 'image/jpg']
  if (validTypes.includes(file.type)) {
    cvFile.value = file
    if (cvPreviewUrl.value) URL.revokeObjectURL(cvPreviewUrl.value)
    cvPreviewUrl.value = URL.createObjectURL(file)
  } else {
    toast.value = { type: 'error', message: 'Vui lòng chọn file PDF, PNG hoặc JPG.' }
  }
}

const handleFileUpload = (e) => {
  validateAndSetFile(e.target.files[0])
}

const handleDrop = (e) => {
  isDragging.value = false
  validateAndSetFile(e.dataTransfer.files[0])
}

const requireLogin = () => {
  // Pass current url to redirect back after login
  router.push({ path: '/login', query: { redirect: route.fullPath } })
}

const submitApplication = async () => {
  if (!isLoggedIn.value) {
    requireLogin()
    return
  }

  if (!cvFile.value) {
    toast.value = { type: 'error', message: 'Vui lòng tải lên CV của bạn.' }
    return
  }

  submitting.value = true
  const formData = new FormData()
  formData.append('cv_file', cvFile.value)

  try {
    await publicService.applyForJob(jobId, formData)
    toast.value = { type: 'success', message: 'Ứng tuyển thành công! Nhà tuyển dụng sẽ sớm liên hệ với bạn.' }
    // Clean up
    cvFile.value = null
    if (cvPreviewUrl.value) URL.revokeObjectURL(cvPreviewUrl.value)
    cvPreviewUrl.value = null
    setTimeout(() => {
      router.push(`/careers/${companyId}`)
    }, 2000)
  } catch (error) {
    toast.value = { type: 'error', message: error.message || 'Có lỗi xảy ra khi nộp đơn.' }
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="job-apply-page">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <Button variant="ghost" style="margin-bottom: 24px; color: var(--text-secondary);" @click="router.push(`/careers/${companyId}`)">
      <ArrowLeft size="16" style="margin-right: 8px;" /> Quay lại danh sách
    </Button>

    <div v-if="loading" style="text-align: center; padding: 40px;">
      <div class="spinner"></div>
    </div>

    <div v-else-if="job" style="display: grid; grid-template-columns: 2fr 1fr; gap: 32px;">
      <!-- Left: Job Details -->
      <div>
        <Card style="padding: 32px; margin-bottom: 24px;">
          <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 24px;">
            <div>
              <h1 style="font-size: 28px; margin: 0 0 12px 0; color: var(--primary);">{{ job.title }}</h1>
              <div style="display: flex; gap: 16px; color: var(--text-secondary); font-size: 15px;">
                <span style="display: flex; align-items: center; gap: 4px;">
                  <MapPin size="16" /> {{ job.location || 'Bất kỳ' }}
                </span>
                <span style="display: flex; align-items: center; gap: 4px;">
                  <Clock size="16" /> {{ job.employment_type || 'Full-time' }}
                </span>
              </div>
            </div>
            <Badge variant="primary">{{ job.department || 'General' }}</Badge>
          </div>

          <div class="job-section">
            <h3>Mô tả công việc</h3>
            <div style="white-space: pre-wrap; line-height: 1.6; color: var(--text-main);">{{ job.description }}</div>
          </div>

          <div class="job-section">
            <h3>Yêu cầu ứng viên</h3>
            <div style="white-space: pre-wrap; line-height: 1.6; color: var(--text-main);">{{ job.requirements }}</div>
          </div>

          <div class="job-section" v-if="job.benefits">
            <h3>Quyền lợi</h3>
            <div style="white-space: pre-wrap; line-height: 1.6; color: var(--text-main);">{{ job.benefits }}</div>
          </div>
        </Card>
      </div>

      <!-- Right: Application Form -->
      <div>
        <Card style="padding: 24px; position: sticky; top: 24px;">
          <h2 style="margin: 0 0 24px 0; font-size: 20px;">Nộp đơn ứng tuyển</h2>

          <div v-if="hasApplied" style="text-align: center; padding: 24px 0;">
            <div style="margin-bottom: 16px; color: var(--success); font-weight: 500;">Bạn đã ứng tuyển vị trí này rồi.</div>
            <Button variant="outline" style="width: 100%" @click="router.push('/home')">Về trang quản lý ứng viên</Button>
          </div>

          <div v-else-if="!isLoggedIn" style="text-align: center; padding: 24px 0;">
            <div style="margin-bottom: 16px; color: var(--text-secondary);">Bạn cần đăng nhập bằng tài khoản Ứng viên để nộp đơn.</div>
            <Button variant="primary" style="width: 100%" @click="requireLogin">Đăng nhập / Đăng ký</Button>
          </div>

          <div v-else>
            <div style="margin-bottom: 16px;">
              <div style="font-weight: 600; font-size: 14px; margin-bottom: 4px;">Thông tin ứng viên:</div>
              <div style="color: var(--text-main);">{{ user?.full_name }}</div>
              <div style="color: var(--text-secondary); font-size: 13px;">{{ user?.email }}</div>
            </div>

            <div style="margin-bottom: 24px;">
              <label style="display: block; font-weight: 500; margin-bottom: 8px; font-size: 14px;">Tải lên CV (PDF, PNG, JPG)</label>
              
              <div 
                class="upload-zone"
                :class="{ 'dragging': isDragging }"
                @dragover.prevent="isDragging = true"
                @dragleave.prevent="isDragging = false"
                @drop.prevent="handleDrop"
                @click="!cvFile && fileInput.click()"
              >
                <UploadCloud v-if="!cvFile" size="32" color="var(--primary)" style="margin-bottom: 12px;" />
                
                <div v-if="cvFile" style="width: 100%; display: flex; flex-direction: column; align-items: center;">
                  <div style="font-weight: 500; color: var(--primary); margin-bottom: 12px; max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ cvFile.name }}</div>
                  
                  <div v-if="cvFile.type.startsWith('image/')" style="width: 100%; border: 1px solid var(--border); border-radius: var(--radius-md); overflow: hidden;">
                    <img :src="cvPreviewUrl" style="max-width: 100%; max-height: 400px; object-fit: contain; display: block; margin: 0 auto;" />
                  </div>
                  <div v-else-if="cvFile.type === 'application/pdf'" style="width: 100%; height: 400px; border: 1px solid var(--border); border-radius: var(--radius-md); overflow: hidden;">
                    <iframe :src="cvPreviewUrl" width="100%" height="100%" style="border: none;"></iframe>
                  </div>
                  
                  <Button variant="secondary" size="sm" style="margin-top: 12px;" @click.stop="cvFile = null; cvPreviewUrl = null">
                    Xóa / Chọn file khác
                  </Button>
                </div>
                <div v-else>
                  <span style="color: var(--primary); font-weight: 500;">Bấm để tải lên</span> hoặc kéo thả file vào đây
                  <div style="font-size: 12px; color: var(--text-muted); margin-top: 4px;">Hỗ trợ PDF, PNG, JPG (Max 10MB)</div>
                </div>
              </div>
              <input type="file" ref="fileInput" accept=".pdf,.png,.jpg,.jpeg" style="display: none" @change="handleFileUpload" />
            </div>

            <Button variant="primary" style="width: 100%" :loading="submitting" @click="submitApplication">
              Nộp đơn ngay
            </Button>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.job-section {
  margin-top: 32px;
}
.job-section h3 {
  font-size: 18px;
  color: var(--text-main);
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
}
.upload-zone {
  border: 2px dashed var(--border);
  border-radius: var(--radius-md);
  padding: 32px 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: var(--background);
}
.upload-zone:hover, .upload-zone.dragging {
  border-color: var(--primary);
  background-color: rgba(37, 99, 235, 0.05);
}
</style>
