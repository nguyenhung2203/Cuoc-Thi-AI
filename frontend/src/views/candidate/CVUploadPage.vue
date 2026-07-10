<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { UploadCloud, FileText, CheckCircle, ArrowLeft, Loader2, Search, Sparkles } from 'lucide-vue-next'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import { candidatePortalService } from '../../services/candidate-portal.service'

const router = useRouter()
const isDragging = ref(false)
const file = ref(null)
const uploadProgress = ref(0)
const uploadStatus = ref('idle') // idle, uploading, parsing, success, error
const parsedData = ref(null)

const handleDragOver = (e) => {
  e.preventDefault()
  isDragging.value = true
}

const handleDragLeave = (e) => {
  e.preventDefault()
  isDragging.value = false
}

const handleDrop = (e) => {
  e.preventDefault()
  isDragging.value = false
  const droppedFiles = e.dataTransfer.files
  if (droppedFiles.length > 0) {
    processFile(droppedFiles[0])
  }
}

const handleFileSelect = (e) => {
  const selectedFiles = e.target.files
  if (selectedFiles.length > 0) {
    processFile(selectedFiles[0])
  }
}

const processFile = async (selectedFile) => {
  if (selectedFile.type !== 'application/pdf') {
    uploadStatus.value = 'error'
    return
  }

  file.value = selectedFile
  uploadStatus.value = 'uploading'
  uploadProgress.value = 0

  try {
    uploadStatus.value = 'parsing'
    // Backend upload CV + parse bằng AI ngay trong 1 request, trả về parsed_data thật.
    const res = await candidatePortalService.uploadCv(selectedFile)
    uploadProgress.value = 100
    if (res && res.cv_url) {
       await loadParsedData(res.parsed_data)
    } else {
       uploadStatus.value = 'error'
    }
  } catch (error) {
    uploadStatus.value = 'error'
  }
}

// Hiển thị dữ liệu CV do AI trích xuất thật. Ưu tiên parsed_data từ response upload,
// bổ sung name/email/phone từ hồ sơ. Không bịa dữ liệu — nếu AI chưa phân tích được thì nói rõ.
const loadParsedData = async (parsedFromUpload) => {
  let parsed = parsedFromUpload || {}
  let profile = {}
  try {
    const res = await candidatePortalService.getProfile()
    profile = res?.data || res || {}
    if (!parsedFromUpload && profile.parsed_data) parsed = profile.parsed_data
  } catch (err) {
    // giữ nguyên parsed từ upload nếu getProfile lỗi
  }
  parsedData.value = {
    name: profile.full_name || '—',
    email: profile.email || '—',
    phone: profile.phone || '—',
    skills: Array.isArray(parsed.skills) ? parsed.skills : [],
    experience: parsed.experience || parsed.summary || 'AI chưa trích xuất được kinh nghiệm từ CV này.'
  }
  uploadStatus.value = 'success'
}

const handleSave = () => {
  // CV đã được lưu ở bước upload; ở đây chỉ điều hướng về trang chủ.
  router.push({
    path: '/home',
    state: { message: 'Đã cập nhật CV thành công!' }
  })
}
</script>

<template>
  <div class="space-y-8 animate-fade-in pb-12 max-w-6xl mx-auto">
    <div class="flex items-center gap-4 mb-8">
      <button class="p-2 bg-white hover:bg-gray-50 text-gray-600 rounded-full shadow-sm transition-colors border border-gray-100" @click="router.back()">
        <ArrowLeft class="w-6 h-6" />
      </button>
      <div>
        <h1 class="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600">Cập nhật CV / Hồ sơ</h1>
        <p class="text-gray-500 mt-1">Upload CV của bạn để AI phân tích kỹ năng và kinh nghiệm.</p>
      </div>
    </div>

    <div class="flex flex-col lg:flex-row gap-8 min-h-[500px]">
      <!-- Upload Zone -->
      <Card class="flex-1 flex flex-col bg-white/80 backdrop-blur-xl border border-white/20 shadow-xl rounded-3xl overflow-hidden relative group">
        <div class="absolute inset-0 bg-gradient-to-br from-blue-50/50 to-indigo-50/50 pointer-events-none"></div>
        
        <div 
          class="flex-1 flex flex-col items-center justify-center p-12 text-center transition-all duration-500 relative z-10 m-6 rounded-2xl border-2 border-dashed"
          :class="[
            isDragging ? 'border-blue-500 bg-blue-50/50 scale-[1.02] shadow-inner' : 'border-gray-200 hover:border-blue-300 hover:bg-gray-50/30',
            file ? 'border-solid border-blue-400 bg-white/60' : 'cursor-pointer'
          ]"
          @dragover="handleDragOver"
          @dragleave="handleDragLeave"
          @drop="handleDrop"
          @click="!file ? $refs.fileInput.click() : null"
        >
          <input 
            type="file" 
            ref="fileInput" 
            accept=".pdf" 
            class="hidden" 
            @change="handleFileSelect"
          />
          
          <!-- State: Idle -->
          <template v-if="uploadStatus === 'idle'">
            <div class="w-24 h-24 mb-6 rounded-full bg-gradient-to-br from-blue-100 to-indigo-100 flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300 relative">
              <div class="absolute inset-0 rounded-full bg-blue-400 opacity-20 animate-ping"></div>
              <UploadCloud class="w-12 h-12 text-blue-600" />
            </div>
            <h3 class="text-2xl font-bold text-gray-800 mb-2">Kéo thả file PDF vào đây</h3>
            <p class="text-gray-500 mb-6 font-medium">Hoặc click để chọn từ thiết bị của bạn (Tối đa 5MB)</p>
            <Button variant="secondary" class="pointer-events-none shadow-md">Chọn File</Button>
          </template>
          
          <!-- State: Uploading -->
          <template v-else-if="uploadStatus === 'uploading'">
            <div class="w-24 h-24 mb-6 rounded-full bg-blue-100 flex items-center justify-center shadow-lg animate-pulse">
              <FileText class="w-12 h-12 text-blue-600" />
            </div>
            <h3 class="text-xl font-bold text-gray-800 mb-6 truncate max-w-full px-4">Đang tải lên: {{ file.name }}</h3>
            
            <div class="w-full max-w-xs bg-gray-100 rounded-full h-2.5 mb-3 overflow-hidden shadow-inner">
              <div class="bg-gradient-to-r from-blue-500 to-indigo-600 h-2.5 rounded-full transition-all duration-200 ease-out relative" :style="{ width: `${uploadProgress}%` }">
                <div class="absolute inset-0 bg-white/20 animate-[shimmer_1s_infinite]"></div>
              </div>
            </div>
            <p class="text-blue-600 font-bold text-sm">{{ uploadProgress }}% hoàn tất</p>
          </template>
          
          <!-- State: Parsing -->
          <template v-else-if="uploadStatus === 'parsing'">
            <div class="w-24 h-24 mb-6 rounded-full bg-indigo-100 flex items-center justify-center shadow-[0_0_30px_rgba(79,70,229,0.3)] relative">
              <div class="absolute inset-0 rounded-full border-4 border-indigo-500 border-t-transparent animate-spin"></div>
              <Search class="w-10 h-10 text-indigo-600 animate-pulse" />
            </div>
            <h3 class="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600 mb-2">AI đang đọc và phân tích...</h3>
            <p class="text-gray-500 font-medium">Hệ thống đang trích xuất kỹ năng và kinh nghiệm của bạn. Vui lòng đợi nhé!</p>
          </template>
          
          <!-- State: Success -->
          <template v-else-if="uploadStatus === 'success'">
            <div class="w-24 h-24 mb-6 rounded-full bg-emerald-100 flex items-center justify-center shadow-lg relative">
              <div class="absolute inset-0 rounded-full bg-emerald-400 opacity-20 animate-[ping_2s_ease-out_1]"></div>
              <CheckCircle class="w-12 h-12 text-emerald-600" />
            </div>
            <h3 class="text-2xl font-bold text-gray-800 mb-2">Phân tích thành công!</h3>
            <div class="flex items-center gap-2 bg-gray-50 px-4 py-2 rounded-lg border border-gray-200 mb-6">
              <FileText class="w-5 h-5 text-gray-400" />
              <span class="text-gray-600 font-medium truncate max-w-[200px]">{{ file.name }}</span>
            </div>
            <Button variant="ghost" @click.stop="uploadStatus = 'idle'; file = null; parsedData = null" class="text-gray-500 hover:text-blue-600">
              Tải lên file khác
            </Button>
          </template>
          
          <!-- State: Error -->
          <template v-else-if="uploadStatus === 'error'">
            <div class="w-24 h-24 mb-6 rounded-full bg-rose-100 flex items-center justify-center shadow-lg">
              <FileText class="w-12 h-12 text-rose-600" />
            </div>
            <h3 class="text-2xl font-bold text-rose-600 mb-2">Lỗi tải lên</h3>
            <p class="text-gray-500 mb-6 font-medium">Chỉ hỗ trợ định dạng PDF. Vui lòng kiểm tra lại file.</p>
            <Button variant="secondary" @click.stop="uploadStatus = 'idle'" class="shadow-sm border-rose-200 hover:border-rose-300 hover:text-rose-600">
              Thử lại
            </Button>
          </template>
        </div>
      </Card>

      <!-- Results Panel -->
      <div v-if="parsedData" class="lg:w-96 animate-[slideInRight_0.5s_ease-out]">
        <Card class="bg-white/95 backdrop-blur-xl border-t-4 border-t-indigo-500 shadow-xl rounded-2xl p-6 h-full flex flex-col">
          <h3 class="text-xl font-bold text-gray-900 mb-2 flex items-center gap-2">
            <Sparkles class="w-5 h-5 text-indigo-500" /> Kết quả AI Trích xuất
          </h3>
          <p class="text-gray-500 text-sm mb-6 pb-4 border-b border-gray-100">
            Vui lòng kiểm tra lại thông tin AI đã nhận diện bên dưới.
          </p>
          
          <div class="space-y-4 flex-1 overflow-y-auto pr-2 custom-scrollbar">
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-1.5">Họ và Tên</label>
              <div class="bg-gray-50/80 border border-gray-100 rounded-lg p-3 text-gray-800 font-medium shadow-sm">{{ parsedData.name }}</div>
            </div>
            
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-1.5">Email</label>
              <div class="bg-gray-50/80 border border-gray-100 rounded-lg p-3 text-gray-800 font-medium shadow-sm">{{ parsedData.email }}</div>
            </div>
            
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-1.5">Số điện thoại</label>
              <div class="bg-gray-50/80 border border-gray-100 rounded-lg p-3 text-gray-800 font-medium shadow-sm">{{ parsedData.phone }}</div>
            </div>
            
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-1.5">Kỹ năng chính</label>
              <div class="flex flex-wrap gap-2 mt-2">
                <span v-for="skill in parsedData.skills" :key="skill"
                  class="px-3 py-1.5 bg-indigo-50 text-indigo-700 border border-indigo-100 rounded-full text-sm font-semibold shadow-sm hover:bg-indigo-100 transition-colors cursor-default">
                  {{ skill }}
                </span>
                <span v-if="!parsedData.skills || parsedData.skills.length === 0" class="text-sm text-gray-400 italic">
                  AI chưa trích xuất được kỹ năng từ CV.
                </span>
              </div>
            </div>
            
            <div>
              <label class="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-1.5">Kinh nghiệm tóm tắt</label>
              <div class="bg-gray-50/80 border border-gray-100 rounded-lg p-3 text-gray-800 font-medium shadow-sm leading-relaxed text-sm">{{ parsedData.experience }}</div>
            </div>
          </div>
          
          <div class="mt-6 pt-4 border-t border-gray-100">
            <button @click="handleSave" class="w-full bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold py-3.5 px-4 rounded-xl shadow-lg hover:shadow-indigo-500/30 transition-all duration-300 transform hover:-translate-y-0.5">
              Xác nhận & Cập nhật Hồ sơ
            </button>
          </div>
        </Card>
      </div>
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

@keyframes slideInRight {
  from { transform: translateX(30px); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}

@keyframes shimmer {
  100% { transform: translateX(100%); }
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 4px;
}
.custom-scrollbar:hover::-webkit-scrollbar-thumb {
  background: #d1d5db;
}
</style>
