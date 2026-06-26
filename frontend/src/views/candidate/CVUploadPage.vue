<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { UploadCloud, FileText, CheckCircle, ArrowLeft, Loader2, Search } from 'lucide-vue-next'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'

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
  
  // Fake upload progress
  const interval = setInterval(() => {
    uploadProgress.value += 10
    if (uploadProgress.value >= 100) {
      clearInterval(interval)
      startParsing()
    }
  }, 200)
}

const startParsing = async () => {
  uploadStatus.value = 'parsing'
  
  // Fake AI parsing delay
  await new Promise(resolve => setTimeout(resolve, 2500))
  
  parsedData.value = {
    name: 'Nguyễn Văn A',
    email: 'nguyenvana@example.com',
    phone: '0901234567',
    skills: ['Vue.js', 'React', 'JavaScript', 'TypeScript', 'Node.js'],
    experience: '3 năm kinh nghiệm lập trình Frontend'
  }
  
  uploadStatus.value = 'success'
}

const handleSave = () => {
  router.push({ 
    path: '/home', 
    state: { message: 'Đã cập nhật CV thành công!' }
  })
}
</script>

<template>
  <div class="page-header">
    <div style="display: flex; align-items: center; gap: 16px;">
      <button class="icon-btn" @click="router.back()">
        <ArrowLeft size="24" />
      </button>
      <h1 class="text-h1">Cập nhật CV / Hồ sơ</h1>
    </div>
  </div>

  <div class="upload-layout">
    <Card style="flex: 1; display: flex; flex-direction: column;">
      <div 
        class="dropzone" 
        :class="{ 'is-dragging': isDragging, 'has-file': file }"
        @dragover="handleDragOver"
        @dragleave="handleDragLeave"
        @drop="handleDrop"
        @click="!file ? $refs.fileInput.click() : null"
      >
        <input 
          type="file" 
          ref="fileInput" 
          accept=".pdf" 
          style="display: none" 
          @change="handleFileSelect"
        />
        
        <template v-if="uploadStatus === 'idle'">
          <div class="upload-icon-wrapper">
            <UploadCloud size="48" color="var(--primary)" />
          </div>
          <h3>Kéo thả file PDF vào đây</h3>
          <p class="text-helper">Hoặc click để chọn từ thiết bị của bạn (Tối đa 5MB)</p>
          <Button variant="secondary" style="margin-top: 16px; pointer-events: none;">Chọn File</Button>
        </template>
        
        <template v-else-if="uploadStatus === 'uploading'">
          <div class="upload-icon-wrapper pulsing">
            <FileText size="48" color="var(--primary)" />
          </div>
          <h3>Đang tải lên: {{ file.name }}</h3>
          <div class="progress-container">
            <div class="progress-bar" :style="{ width: `${uploadProgress}%` }"></div>
          </div>
          <p class="text-helper">{{ uploadProgress }}% hoàn tất</p>
        </template>
        
        <template v-else-if="uploadStatus === 'parsing'">
          <div class="upload-icon-wrapper spinning">
            <Search size="48" color="var(--accent)" />
          </div>
          <h3>AI đang đọc và phân tích CV của bạn...</h3>
          <p class="text-helper">Việc này có thể mất vài giây. Vui lòng đợi nhé!</p>
        </template>
        
        <template v-else-if="uploadStatus === 'success'">
          <div class="upload-icon-wrapper bg-success-light">
            <CheckCircle size="48" color="var(--success)" />
          </div>
          <h3>Tải lên và Phân tích thành công!</h3>
          <p class="text-helper">{{ file.name }}</p>
          <Button variant="ghost" style="margin-top: 16px;" @click.stop="uploadStatus = 'idle'; file = null; parsedData = null">
            Tải lên file khác
          </Button>
        </template>
        
        <template v-else-if="uploadStatus === 'error'">
          <div class="upload-icon-wrapper bg-danger-light">
            <FileText size="48" color="var(--danger)" />
          </div>
          <h3 style="color: var(--danger)">Lỗi tải lên</h3>
          <p class="text-helper">Chỉ hỗ trợ định dạng PDF. Vui lòng thử lại.</p>
          <Button variant="secondary" style="margin-top: 16px;" @click.stop="uploadStatus = 'idle'">
            Thử lại
          </Button>
        </template>
      </div>
    </Card>

    <div class="side-panel" v-if="parsedData">
      <Card title="Kết quả AI Trích xuất">
        <p class="text-helper" style="margin-bottom: 20px;">
          AI của chúng tôi đã đọc CV và nhận diện được các thông tin sau. Vui lòng kiểm tra lại.
        </p>
        
        <div class="parsed-field">
          <label>Họ và Tên</label>
          <div class="field-value">{{ parsedData.name }}</div>
        </div>
        
        <div class="parsed-field">
          <label>Email</label>
          <div class="field-value">{{ parsedData.email }}</div>
        </div>
        
        <div class="parsed-field">
          <label>Số điện thoại</label>
          <div class="field-value">{{ parsedData.phone }}</div>
        </div>
        
        <div class="parsed-field">
          <label>Kỹ năng chính</label>
          <div class="skills-wrapper">
            <span v-for="skill in parsedData.skills" :key="skill" class="skill-tag">
              {{ skill }}
            </span>
          </div>
        </div>
        
        <div class="parsed-field">
          <label>Kinh nghiệm tóm tắt</label>
          <div class="field-value">{{ parsedData.experience }}</div>
        </div>
        
        <div style="margin-top: 32px;">
          <Button variant="primary" style="width: 100%;" @click="handleSave">
            Xác nhận & Cập nhật Hồ sơ
          </Button>
        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.upload-layout {
  display: flex;
  gap: 24px;
  min-height: 500px;
}

.dropzone {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 2px dashed var(--border);
  border-radius: var(--radius-lg);
  background-color: var(--surface-soft);
  transition: all 0.3s ease;
  cursor: pointer;
  padding: 40px 20px;
  text-align: center;
}

.dropzone:hover:not(.has-file), .dropzone.is-dragging {
  border-color: var(--primary);
  background-color: rgba(37, 99, 235, 0.05);
}

.dropzone.has-file {
  cursor: default;
  border-style: solid;
}

.upload-icon-wrapper {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  background-color: white;
  box-shadow: var(--shadow-md);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  transition: all 0.3s;
}

.bg-success-light { background-color: rgba(22, 163, 74, 0.1); }
.bg-danger-light { background-color: rgba(220, 38, 38, 0.1); }

.dropzone h3 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-main);
  margin-bottom: 8px;
}

.progress-container {
  width: 100%;
  max-width: 300px;
  height: 8px;
  background-color: var(--border);
  border-radius: 4px;
  margin: 16px 0 8px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background-color: var(--primary);
  transition: width 0.2s ease;
}

.pulsing {
  animation: pulse 2s infinite;
}

.spinning {
  animation: pulse 1.5s infinite alternate;
}

@keyframes pulse {
  0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(37, 99, 235, 0.4); }
  70% { transform: scale(1); box-shadow: 0 0 0 15px rgba(37, 99, 235, 0); }
  100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(37, 99, 235, 0); }
}

.side-panel {
  width: 400px;
  animation: slideInRight 0.4s ease-out;
}

.parsed-field {
  margin-bottom: 16px;
}

.parsed-field label {
  display: block;
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.field-value {
  font-size: 14px;
  color: var(--text-main);
  font-weight: 500;
  padding: 8px 12px;
  background-color: var(--surface-soft);
  border-radius: var(--radius);
  border: 1px solid var(--border);
}

.skills-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px 0;
}

.skill-tag {
  padding: 4px 12px;
  background-color: rgba(37, 99, 235, 0.1);
  color: var(--primary);
  border-radius: 99px;
  font-size: 13px;
  font-weight: 500;
}

@keyframes slideInRight {
  from { transform: translateX(20px); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}

@media (max-width: 1024px) {
  .upload-layout {
    flex-direction: column;
  }
  .side-panel {
    width: 100%;
  }
}
</style>
