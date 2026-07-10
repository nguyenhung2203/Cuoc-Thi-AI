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
    <div class="flex items-center gap-4 mb-8 animate-rise">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="20" />
      </button>
      <div>
        <h1 class="page-title">Cập nhật CV / Hồ sơ</h1>
        <p class="page-subtitle">Upload CV của bạn để AI phân tích kỹ năng và kinh nghiệm.</p>
      </div>
    </div>

    <div class="flex flex-col lg:flex-row gap-8 min-h-[500px]">
      <!-- Upload Zone -->
      <Card class="flex-1 flex flex-col card-elevate overflow-hidden relative group animate-rise">
        <div
          class="dropzone"
          :class="[
            isDragging ? 'is-dragging' : '',
            file ? 'has-file' : 'cursor-pointer'
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
            <div class="dz-icon">
              <span class="dz-icon-ping"></span>
              <UploadCloud :size="44" />
            </div>
            <h3 class="dz-title">Kéo thả file PDF vào đây</h3>
            <p class="dz-sub">Hoặc click để chọn từ thiết bị của bạn (Tối đa 5MB)</p>
            <Button variant="secondary" class="pointer-events-none">Chọn File</Button>
          </template>
          
          <!-- State: Uploading -->
          <template v-else-if="uploadStatus === 'uploading'">
            <div class="dz-icon is-soft"><FileText :size="44" /></div>
            <h3 class="dz-title-sm">Đang tải lên: {{ file.name }}</h3>
            <div class="dz-progress">
              <div class="dz-progress-bar" :style="{ width: `${uploadProgress}%` }"></div>
            </div>
            <p class="dz-percent">{{ uploadProgress }}% hoàn tất</p>
          </template>

          <!-- State: Parsing -->
          <template v-else-if="uploadStatus === 'parsing'">
            <div class="dz-icon is-accent-icon">
              <span class="dz-spin-ring"></span>
              <Search :size="40" />
            </div>
            <h3 class="dz-title" style="color: var(--accent);">AI đang đọc và phân tích...</h3>
            <p class="dz-sub">Hệ thống đang trích xuất kỹ năng và kinh nghiệm của bạn. Vui lòng đợi nhé!</p>
          </template>

          <!-- State: Success -->
          <template v-else-if="uploadStatus === 'success'">
            <div class="dz-icon is-success-icon">
              <span class="dz-icon-ping is-success-ping"></span>
              <CheckCircle :size="44" />
            </div>
            <h3 class="dz-title">Phân tích thành công!</h3>
            <div class="dz-file-chip">
              <FileText :size="18" />
              <span class="truncate">{{ file.name }}</span>
            </div>
            <Button variant="ghost" @click.stop="uploadStatus = 'idle'; file = null; parsedData = null">
              Tải lên file khác
            </Button>
          </template>

          <!-- State: Error -->
          <template v-else-if="uploadStatus === 'error'">
            <div class="dz-icon is-error-icon"><FileText :size="44" /></div>
            <h3 class="dz-title" style="color: var(--danger);">Lỗi tải lên</h3>
            <p class="dz-sub">Chỉ hỗ trợ định dạng PDF. Vui lòng kiểm tra lại file.</p>
            <Button variant="secondary" @click.stop="uploadStatus = 'idle'">Thử lại</Button>
          </template>
        </div>
      </Card>

      <!-- Results Panel -->
      <div v-if="parsedData" class="lg:w-96 result-panel">
        <Card class="card-elevate result-card">
          <h3 class="section-heading mb-2">
            <span class="kpi-icon is-accent"><Sparkles :size="18" /></span>
            Kết quả AI Trích xuất
          </h3>
          <p class="result-hint">Vui lòng kiểm tra lại thông tin AI đã nhận diện bên dưới.</p>

          <div class="space-y-4 flex-1 overflow-y-auto pr-1 custom-scrollbar">
            <div>
              <label class="result-label">Họ và Tên</label>
              <div class="result-value">{{ parsedData.name }}</div>
            </div>
            <div>
              <label class="result-label">Email</label>
              <div class="result-value">{{ parsedData.email }}</div>
            </div>
            <div>
              <label class="result-label">Số điện thoại</label>
              <div class="result-value">{{ parsedData.phone }}</div>
            </div>
            <div>
              <label class="result-label">Kỹ năng chính</label>
              <div class="flex flex-wrap gap-2 mt-1">
                <span v-for="skill in parsedData.skills" :key="skill" class="skill-chip">{{ skill }}</span>
                <span v-if="!parsedData.skills || parsedData.skills.length === 0" class="result-empty">
                  AI chưa trích xuất được kỹ năng từ CV.
                </span>
              </div>
            </div>
            <div>
              <label class="result-label">Kinh nghiệm tóm tắt</label>
              <div class="result-value result-value-text">{{ parsedData.experience }}</div>
            </div>
          </div>

          <div class="result-foot">
            <button @click="handleSave" class="btn btn-primary sheen w-full">
              Xác nhận & Cập nhật Hồ sơ
            </button>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.back-btn { display: inline-flex; align-items: center; justify-content: center; width: 42px; height: 42px; border-radius: 50%; background: var(--surface); border: 1px solid var(--border); color: var(--text-secondary); cursor: pointer; box-shadow: var(--shadow-sm); transition: all 0.2s ease; flex-shrink: 0; }
.back-btn:hover { color: var(--primary); border-color: var(--primary-light); }

.dropzone { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; text-align: center; padding: 48px; margin: 22px; border: 2px dashed var(--border); border-radius: var(--radius-lg); transition: all 0.35s ease; }
.dropzone.cursor-pointer { cursor: pointer; }
.dropzone.cursor-pointer:hover { border-color: var(--primary); background: var(--primary-light); }
.dropzone.is-dragging { border-color: var(--primary); background: var(--primary-light); transform: scale(1.01); }
.dropzone.has-file { border-style: solid; border-color: var(--primary); }

.dz-icon { position: relative; width: 88px; height: 88px; display: flex; align-items: center; justify-content: center; border-radius: 50%; background: var(--primary-light); color: var(--primary); margin-bottom: 22px; box-shadow: var(--shadow-md); }
.dz-icon.is-soft { animation: softPulse 1.4s ease-in-out infinite; }
.dz-icon.is-accent-icon { background: var(--accent-bg); color: var(--accent); }
.dz-icon.is-success-icon { background: rgba(22,163,74,0.12); color: var(--success); }
.dz-icon.is-error-icon { background: rgba(220,38,38,0.12); color: var(--danger); }
.dz-icon-ping { position: absolute; inset: 0; border-radius: 50%; background: var(--primary); opacity: 0.18; animation: dzping 1.8s ease-out infinite; }
.dz-icon-ping.is-success-ping { background: var(--success); }
.dz-spin-ring { position: absolute; inset: 0; border-radius: 50%; border: 3px solid var(--accent); border-top-color: transparent; animation: spin 0.9s linear infinite; }
@keyframes dzping { 0% { transform: scale(1); opacity: 0.18; } 100% { transform: scale(1.35); opacity: 0; } }
@keyframes softPulse { 0%,100% { transform: scale(1); } 50% { transform: scale(1.06); } }
@keyframes spin { to { transform: rotate(360deg); } }

.dz-title { font-size: 20px; font-weight: 700; color: var(--text-main); margin-bottom: 8px; }
.dz-title-sm { font-size: 16px; font-weight: 600; color: var(--text-main); margin-bottom: 20px; max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; padding: 0 12px; }
.dz-sub { color: var(--text-secondary); margin-bottom: 22px; }
.dz-progress { width: 100%; max-width: 20rem; height: 8px; border-radius: 999px; background: var(--surface-soft); overflow: hidden; margin-bottom: 10px; }
.dz-progress-bar { height: 100%; border-radius: 999px; background: var(--gradient-brand); transition: width 0.2s ease; }
.dz-percent { color: var(--primary); font-weight: 700; font-size: 13px; }
.dz-file-chip { display: inline-flex; align-items: center; gap: 8px; padding: 8px 16px; border-radius: var(--radius); background: var(--surface-soft); border: 1px solid var(--border); color: var(--text-secondary); font-weight: 500; margin-bottom: 22px; max-width: 260px; }

.result-panel { animation: riseIn 0.5s cubic-bezier(0.2,0.8,0.2,1) both; }
.result-card { display: flex; flex-direction: column; height: 100%; padding: 24px; border-top: 3px solid var(--accent); }
.result-hint { color: var(--text-secondary); font-size: 14px; margin: 8px 0 18px; padding-bottom: 16px; border-bottom: 1px solid var(--border); }
.result-label { display: block; font-size: 11px; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 6px; }
.result-value { background: var(--surface-soft); border: 1px solid var(--border); border-radius: var(--radius); padding: 12px; color: var(--text-main); font-weight: 500; }
.result-value-text { font-size: 14px; line-height: 1.55; font-weight: 400; }
.result-empty { font-size: 14px; color: var(--text-muted); font-style: italic; }
.skill-chip { padding: 6px 12px; border-radius: var(--radius-full); background: var(--accent-bg); color: var(--accent); border: 1px solid rgba(8,145,178,0.2); font-size: 13px; font-weight: 600; }
.result-foot { margin-top: 22px; padding-top: 18px; border-top: 1px solid var(--border); }

.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: var(--border); border-radius: 4px; }

@media (prefers-reduced-motion: reduce) {
  .dz-icon.is-soft, .dz-icon-ping, .dz-spin-ring { animation: none !important; }
  .result-panel { animation: none !important; }
}
</style>
