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
  <div class="space-y-8 animate-fade-in pb-12 max-w-6xl mx-auto">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
    
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600">CV của tôi</h1>
        <p class="text-gray-500 mt-2 text-lg">Tải lên và quản lý hồ sơ xin việc (CV/Resume) của bạn.</p>
      </div>
    </div>

    <input 
      type="file" 
      ref="fileInput" 
      accept=".pdf,.png,.jpg,.jpeg" 
      class="hidden" 
      @change="handleFileUpload"
    />

    <!-- Empty State -->
    <Card v-if="!cvUrl && !cvPreviewUrl" class="bg-white/80 backdrop-blur-xl border border-white/20 shadow-xl rounded-3xl overflow-hidden p-12 text-center relative group">
      <div class="absolute inset-0 bg-gradient-to-br from-blue-50/50 to-indigo-50/50 pointer-events-none"></div>
      <div class="relative z-10 flex flex-col items-center justify-center py-12">
        <div class="w-32 h-32 bg-blue-50 rounded-full flex items-center justify-center mb-8 shadow-inner group-hover:scale-110 transition-transform duration-500">
          <FileText class="w-16 h-16 text-blue-400" />
        </div>
        <h2 class="text-2xl font-bold text-gray-800 mb-3">Chưa có CV nào</h2>
        <p class="text-gray-500 mb-8 max-w-md mx-auto text-lg">
          Tải lên CV của bạn (định dạng PDF, PNG, JPG) để các nhà tuyển dụng dễ dàng đánh giá kỹ năng của bạn hơn.
        </p>
        
        <button @click="fileInput.click()" :disabled="uploading" class="bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold py-3.5 px-8 rounded-xl shadow-lg hover:shadow-indigo-500/30 transition-all duration-300 transform hover:-translate-y-0.5 text-lg flex items-center gap-2 disabled:opacity-50 disabled:transform-none">
          <Upload class="w-6 h-6" /> 
          {{ uploading ? 'Đang tải lên...' : 'Tải lên CV mới' }}
        </button>
      </div>
    </Card>

    <!-- Uploaded State -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <!-- CV Info Sidebar -->
      <div class="lg:col-span-1 space-y-6">
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl p-6 relative overflow-hidden">
          <div class="absolute top-0 right-0 w-24 h-24 bg-gradient-to-br from-blue-100 to-transparent rounded-bl-full opacity-50"></div>
          
          <div class="flex items-center gap-4 mb-6">
            <div class="w-12 h-12 bg-emerald-50 rounded-xl flex items-center justify-center text-emerald-600 shadow-sm">
              <CheckCircle2 class="w-6 h-6" />
            </div>
            <div>
              <p class="font-bold text-gray-900 leading-tight">CV của bạn <br> đã sẵn sàng</p>
            </div>
          </div>
          
          <div class="bg-blue-50 border border-blue-100 rounded-xl p-4 mb-6">
            <div class="flex items-center gap-2 text-blue-700 font-medium mb-1 truncate">
              <FileText class="w-4 h-4 flex-shrink-0" />
              <span class="truncate" :title="cvName">{{ cvName }}</span>
            </div>
          </div>
          
          <button @click="fileInput.click()" :disabled="uploading" class="w-full py-2.5 px-4 bg-white border-2 border-gray-200 text-gray-700 font-bold rounded-xl hover:bg-gray-50 hover:border-gray-300 transition-colors shadow-sm flex items-center justify-center gap-2 disabled:opacity-50">
            <Upload class="w-4 h-4" /> 
            {{ uploading ? 'Đang cập nhật...' : 'Cập nhật CV mới' }}
          </button>
        </Card>
      </div>
      
      <!-- CV Preview Area -->
      <Card class="lg:col-span-3 bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl overflow-hidden p-2">
        <div class="bg-gray-100 rounded-xl overflow-hidden h-[800px] relative border border-gray-200 shadow-inner flex items-center justify-center">
          <img v-if="isImage" :src="cvPreviewUrl || cvUrl" class="max-w-full max-h-full object-contain drop-shadow-lg" />
          <iframe v-else :src="cvPreviewUrl || cvUrl" width="100%" height="100%" class="border-none w-full h-full bg-white"></iframe>
        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in {
  animation: fade-in 0.5s ease-out forwards;
}
</style>
