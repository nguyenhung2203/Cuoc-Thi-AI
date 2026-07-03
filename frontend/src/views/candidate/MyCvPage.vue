<script setup>
import { ref, onMounted } from 'vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { Upload, FileText, CheckCircle2 } from 'lucide-vue-next'
import { candidatePortalService } from '../../services/candidate-portal.service'

const fileInput = ref(null)
const uploading = ref(false)
const toast = ref(null)
const cvUrl = ref(null)
const cvPreviewUrl = ref(null)
const isImage = ref(false)
const cvName = ref('')

onMounted(async () => {
  try {
    const profile = await candidatePortalService.getProfile()
    if (profile.cv_url) {
      cvUrl.value = profile.cv_url
      cvName.value = profile.cv_name || 'My_CV_Current.pdf'
      const lowerUrl = cvUrl.value.toLowerCase()
      if (lowerUrl.endsWith('.png') || lowerUrl.endsWith('.jpg') || lowerUrl.endsWith('.jpeg')) {
        isImage.value = true
      }
    }
  } catch (err) {
    console.error('Lỗi tải profile:', err)
  }
})

const handleFileUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const validTypes = ['application/pdf', 'image/png', 'image/jpeg', 'image/jpg']
  if (!validTypes.includes(file.type)) {
    toast.value = { type: 'error', message: 'Vui lòng chọn file PDF, PNG hoặc JPG' }
    return
  }
  
  if (file.size > 5 * 1024 * 1024) {
    toast.value = { type: 'error', message: 'Kích thước file không được vượt quá 5MB' }
    return
  }

  isImage.value = file.type.startsWith('image/')
  if (cvPreviewUrl.value) URL.revokeObjectURL(cvPreviewUrl.value)
  cvPreviewUrl.value = URL.createObjectURL(file)

  uploading.value = true
  toast.value = null

  try {
    await candidatePortalService.uploadCv(file)
    toast.value = { type: 'success', message: 'Tải lên CV thành công!' }
    // Fetch profile again to update the cv_url
    const profile = await candidatePortalService.getProfile()
    cvUrl.value = profile.cv_url || 'dummy-url-for-demo' // Backend is stubbed for now
    cvName.value = profile.cv_name || file.name
    // Remove the revokeObjectURL here so the user can continue seeing what they uploaded
  } catch (err) {
    toast.value = { type: 'error', message: err.message || 'Lỗi tải lên file' }
  } finally {
    uploading.value = false
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }
}
</script>

<template>
  <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
  
  <div class="page-header">
    <div>
      <h1 class="text-h1">CV của tôi</h1>
      <p class="text-helper" style="margin-top: 4px">Tải lên và quản lý hồ sơ xin việc (CV/Resume) của bạn.</p>
    </div>
  </div>

  <input 
    type="file" 
    ref="fileInput" 
    accept=".pdf,.png,.jpg,.jpeg" 
    style="display: none" 
    @change="handleFileUpload"
  />
  
  <div v-if="!cvUrl" style="padding: 48px; text-align: center; border: 1px dashed var(--border); border-radius: var(--radius-lg); background-color: var(--surface-soft)">
    <div style="margin-bottom: 16px; display: inline-flex; justify-content: center; align-items: center; width: 64px; height: 64px; border-radius: 50%; background-color: rgba(37, 99, 235, 0.1); color: var(--primary)">
      <FileText size="32" />
    </div>
    <h2 class="text-h2" style="margin-bottom: 8px">Chưa có CV nào</h2>
    <p class="text-body" style="color: var(--text-secondary); max-width: 400px; margin: 0 auto 24px">
      Tải lên CV của bạn (định dạng PDF, PNG, JPG) để các nhà tuyển dụng dễ dàng đánh giá kỹ năng của bạn hơn.
    </p>
    <Button @click="fileInput.click()" :disabled="uploading">
      <Upload size="16" style="margin-right: 8px" /> 
      {{ uploading ? 'Đang tải lên...' : 'Tải lên CV mới' }}
    </Button>
  </div>
  
  <div v-else style="padding: 32px; border: 1px solid var(--border); border-radius: var(--radius-lg); background-color: var(--surface)">
    <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 24px">
      <div style="width: 48px; height: 48px; border-radius: 8px; background-color: rgba(34, 197, 94, 0.1); color: var(--success); display: flex; align-items: center; justify-content: center">
        <CheckCircle2 size="24" />
      </div>
      <div>
        <h3 class="text-h3">CV của bạn đã sẵn sàng</h3>
        <p class="text-body" style="color: var(--text-secondary)">Các nhà tuyển dụng có thể xem CV này khi bạn ứng tuyển.</p>
      </div>
    </div>
    
    <div style="padding: 16px; border: 1px dashed var(--border); border-radius: 8px; background-color: var(--surface-soft); display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px">
      <div style="display: flex; align-items: center; gap: 12px">
        <FileText size="20" color="var(--primary)" />
        <span style="font-weight: 500">{{ cvName }}</span>
      </div>
      <div style="display: flex; gap: 8px">
        <Button variant="outline" @click="fileInput.click()" :disabled="uploading">
          {{ uploading ? 'Đang cập nhật...' : 'Cập nhật CV' }}
        </Button>
      </div>
    </div>
    
    <div v-if="cvPreviewUrl || cvUrl" style="width: 100%; border: 1px solid var(--border); border-radius: var(--radius-md); overflow: hidden;">
      <img v-if="isImage" :src="cvPreviewUrl || cvUrl" style="max-width: 100%; max-height: 500px; object-fit: contain; display: block; margin: 0 auto;" />
      <iframe v-else :src="cvPreviewUrl || cvUrl" width="100%" height="500px" style="border: none; display: block;"></iframe>
    </div>
  </div>
</template>
