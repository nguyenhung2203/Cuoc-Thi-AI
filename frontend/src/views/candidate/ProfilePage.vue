<script setup>
import { ref, onMounted } from 'vue'
import { authStore } from '../../stores/auth.store'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import { Upload, FileText, CheckCircle, Save, Trash2, Eye, X, Bot, Loader2, UserRound, Target, Code2, Info } from 'lucide-vue-next'
import { candidatePortalService } from '../../services/candidate-portal.service'

const profile = ref({
  name: '',
  email: '',
  phone: '',
  linkedin: '',
  targetRole: '',
  level: 'Middle',
  skills: '',
})

const uploadedCvs = ref([])
const fileInput = ref(null)
const selectedCv = ref(null)

const handleFileUpload = async (e) => {
  const files = e.target.files
  if (files && files.length > 0) {
    for (let i = 0; i < files.length; i++) {
      const file = files[i]
      const cvId = Date.now() + i
      
      // Temporary local URL for preview if needed immediately
      const tempUrl = URL.createObjectURL(file)

      const cvData = {
        id: cvId,
        name: file.name,
        size: (file.size / 1024 / 1024).toFixed(2) + ' MB',
        date: new Date().toLocaleDateString('vi-VN'),
        status: 'analyzing',
        url: tempUrl,
        parsedData: null
      }
      
      if (!uploadedCvs.value.find(cv => cv.name === file.name)) {
        uploadedCvs.value.unshift(cvData)
        
        try {
          // Upload to DB
          const res = await candidatePortalService.uploadCv(file)
          
          const targetCv = uploadedCvs.value.find(cv => cv.id === cvId)
          if (targetCv) {
            targetCv.status = 'done'
            if (res.data && res.data.cv_url) {
              targetCv.url = res.data.cv_url
            }
            // Use real parsed data from CV upload response, if available
            // Otherwise show empty state instead of fake data
            if (res.data && res.data.parsed_data) {
              targetCv.parsedData = res.data.parsed_data
            } else {
              targetCv.parsedData = null
            }
            localStorage.setItem('candidate_cvs', JSON.stringify(uploadedCvs.value))
          }
        } catch (error) {
          console.error("Upload failed", error)
          const targetCv = uploadedCvs.value.find(cv => cv.id === cvId)
          if (targetCv) {
            targetCv.status = 'done'
            toast.value = { type: 'error', message: `Lỗi tải lên ${file.name}` }
          }
        }
      }
    }
  }
}

const deleteCv = (id) => {
  uploadedCvs.value = uploadedCvs.value.filter(cv => cv.id !== id)
  localStorage.setItem('candidate_cvs', JSON.stringify(uploadedCvs.value))
  if (fileInput.value) fileInput.value.value = ''
}

const viewCv = (cv) => {
  selectedCv.value = cv
}

const closeCvModal = () => {
  selectedCv.value = null
}

onMounted(async () => {
  if (authStore.user) {
    profile.value.name = authStore.user.full_name || ''
    profile.value.email = authStore.user.email || ''
  }

  // Load backend profile
  try {
    const res = await candidatePortalService.getProfile()
    if (res.data) {
      if (res.data.full_name) profile.value.name = res.data.full_name
      if (res.data.email) profile.value.email = res.data.email
      
      // Load CV from DB if it exists and not already loaded
      if (res.data.cv_url && res.data.cv_name) {
        const hasCv = uploadedCvs.value.find(cv => cv.name === res.data.cv_name)
        if (!hasCv) {
          uploadedCvs.value.push({
            id: 'db-' + Date.now(),
            name: res.data.cv_name,
            size: 'N/A',
            date: 'Từ hệ thống',
            status: 'done',
            url: res.data.cv_url,
            // Use real parsed CV data if available, otherwise null (empty state)
            parsedData: res.data.parsed_data || null
          })
        }
      }
    }
  } catch (error) {
    console.error("Failed to load profile from DB", error)
  }
  
  const savedProfile = localStorage.getItem('candidate_profile')
  if (savedProfile) {
    const parsed = JSON.parse(savedProfile)
    // Only merge non-core fields since core fields are from DB
    profile.value.targetRole = parsed.targetRole || profile.value.targetRole
    profile.value.level = parsed.level || profile.value.level
    profile.value.phone = parsed.phone || profile.value.phone
    profile.value.linkedin = parsed.linkedin || profile.value.linkedin
    profile.value.skills = parsed.skills || profile.value.skills
  }

  const savedCvs = localStorage.getItem('candidate_cvs')
  if (savedCvs) {
    const parsedCvs = JSON.parse(savedCvs)
    // Merge without duplicating
    parsedCvs.forEach(savedCv => {
      if (!uploadedCvs.value.find(cv => cv.name === savedCv.name)) {
        uploadedCvs.value.push(savedCv)
      }
    })
  }
})

const saving = ref(false)
const toast = ref(null)

const handleSave = async (e) => {
  e.preventDefault()
  saving.value = true
  try {
    await candidatePortalService.updateProfile({
      full_name: profile.value.name,
      avatar_url: profile.value.avatar_url || ''
    })
    localStorage.setItem('candidate_profile', JSON.stringify(profile.value))
    toast.value = { type: 'success', message: 'Hồ sơ cá nhân đã được lưu thành công! Dữ liệu này sẽ được đồng bộ với AI.' }
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi lưu hồ sơ: ' + (error.message || 'Không xác định') }
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="space-y-8 pb-12 max-w-6xl mx-auto">
    <div class="mb-2 animate-rise">
      <h1 class="page-title">Hồ sơ cá nhân & CV</h1>
      <p class="page-subtitle">Cập nhật thông tin để AI có thể đưa ra bài luyện tập chính xác nhất.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">

      <!-- Left Column: Forms -->
      <div class="lg:col-span-2 space-y-8">
        <Card class="card-elevate pf-card animate-rise">
          <div class="pf-card-head">
            <h3 class="pf-card-title">
              <span class="kpi-icon"><UserRound :size="18" /></span>
              Thông tin cơ bản
            </h3>
          </div>

          <div class="p-6">
            <form @submit="handleSave" class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- Custom styling for inputs can be handled by AppInput if it uses standard HTML attributes -->
                <!-- We'll wrap them in div to provide some spacing if needed -->
                <div class="space-y-1">
                  <Input label="Họ và Tên" v-model="profile.name" required />
                </div>
                <div class="space-y-1">
                  <Input label="Email" type="email" v-model="profile.email" disabled />
                </div>
                <div class="space-y-1">
                  <Input label="Số điện thoại" v-model="profile.phone" />
                </div>
                <div class="space-y-1">
                  <Input label="LinkedIn Profile" v-model="profile.linkedin" />
                </div>
              </div>

              <div class="pt-6 mt-6 border-t border-gray-100">
                <h3 class="pf-sub-title">
                  <span class="kpi-icon is-accent"><Target :size="18" /></span>
                  Định hướng nghề nghiệp
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <Input label="Vị trí mục tiêu (Target Role)" v-model="profile.targetRole" />
                  <div class="flex flex-col gap-2">
                    <label class="pf-label">Cấp độ hiện tại</label>
                    <select class="pf-select" v-model="profile.level">
                      <option value="Intern">Intern</option>
                      <option value="Fresher">Fresher</option>
                      <option value="Junior">Junior</option>
                      <option value="Middle">Middle</option>
                      <option value="Senior">Senior</option>
                    </select>
                  </div>
                </div>
              </div>

              <div class="flex justify-end pt-6">
                <button type="submit" :disabled="saving" class="btn btn-primary sheen">
                  <Save :size="18" />
                  {{ saving ? 'Đang lưu...' : 'Lưu hồ sơ' }}
                </button>
              </div>
            </form>
          </div>
        </Card>

        <Card class="card-elevate pf-card animate-rise">
          <div class="pf-card-head">
            <h3 class="pf-card-title">
              <span class="kpi-icon is-accent"><Code2 :size="18" /></span>
              Kỹ năng chuyên môn
            </h3>
          </div>
          <div class="p-6">
            <textarea class="pf-textarea" rows="4" placeholder="Ví dụ: ReactJS, NodeJS, TypeScript..." v-model="profile.skills"></textarea>
            <p class="text-helper mt-2">Phân cách các kỹ năng bằng dấu phẩy (,)</p>
          </div>
        </Card>
      </div>

      <!-- Right Column: CV Upload -->
      <div class="space-y-8">
        <Card class="card-elevate pf-card sticky top-6 animate-rise">
          <div class="pf-card-head">
            <h3 class="pf-card-title">
              <span class="kpi-icon"><FileText :size="18" /></span>
              CV của bạn
            </h3>
          </div>

          <div class="p-6 space-y-4">
            <!-- Hidden File Input -->
            <input type="file" ref="fileInput" @change="handleFileUpload" accept=".pdf,.doc,.docx" multiple class="hidden" style="display: none;" />

            <!-- Upload Area -->
            <div class="pf-dropzone" @click="fileInput.click()">
              <div class="pf-dz-icon"><Upload :size="22" /></div>
              <p class="pf-dz-title">Tải CV lên (Nhiều file)</p>
              <p class="text-helper">PDF, DOCX (Tối đa 5MB/file)</p>
            </div>

            <!-- List of Current CVs -->
            <div v-if="uploadedCvs.length > 0" class="space-y-3 mt-6">
              <h4 class="pf-list-title">Danh sách CV đã tải lên</h4>

              <div v-for="cv in uploadedCvs" :key="cv.id" @click="cv.status === 'done' && viewCv(cv)" class="pf-cv-item group" :class="cv.status === 'done' ? 'is-done' : 'is-loading'">
                <div class="pf-cv-icon" :class="cv.status === 'done' ? '' : 'is-muted'">
                  <FileText :size="20" />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="pf-cv-name truncate">{{ cv.name }}</p>
                  <p v-if="cv.status === 'done'" class="pf-cv-meta is-ok">
                    <CheckCircle :size="13" /> Tải lên: {{ cv.date }} • {{ cv.size }}
                  </p>
                  <p v-else class="pf-cv-meta is-wait">AI đang đọc thông tin...</p>
                </div>

                <div class="flex gap-1">
                  <div v-if="cv.status === 'analyzing'" class="pf-analyzing">
                    <Loader2 :size="15" class="pf-spin" /> Đang phân tích
                  </div>
                  <template v-else>
                    <button type="button" @click.stop="viewCv(cv)" class="pf-icon-btn" title="Xem chi tiết"><Eye :size="18" /></button>
                    <button type="button" @click.stop="deleteCv(cv.id)" class="pf-icon-btn is-danger" title="Xóa CV"><Trash2 :size="18" /></button>
                  </template>
                </div>
              </div>
            </div>

            <div class="ai-block pf-note">
              <Info :size="20" />
              <p>CV của bạn sẽ được AI dùng làm ngữ cảnh (context) để đặt câu hỏi sát với kinh nghiệm thực tế của bạn trong phần Luyện phỏng vấn Mock Interview.</p>
            </div>
          </div>
        </Card>
      </div>

    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <!-- CV Detail Modal -->
    <div v-if="selectedCv" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 md:p-8" @click="closeCvModal">
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-5xl h-[90vh] flex flex-col overflow-hidden animate-fade-in" @click.stop>
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-6 border-b border-slate-100 bg-slate-50 shrink-0">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center text-blue-600">
              <FileText class="w-6 h-6" />
            </div>
            <div>
              <h3 class="text-xl font-bold text-slate-800">{{ selectedCv.name }}</h3>
              <p class="text-sm text-slate-500 mt-1">Tải lên ngày {{ selectedCv.date }} • Dữ liệu được trích xuất bởi AI</p>
            </div>
          </div>
          <button @click="closeCvModal" class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-200 rounded-full transition-colors shrink-0">
            <X class="w-6 h-6" />
          </button>
        </div>
        
        <!-- Modal Body (Mock Viewer + Extracted Data) -->
        <div class="p-6 flex-1 overflow-y-auto bg-slate-100/50 flex flex-col lg:flex-row gap-8">
          
          <!-- Actual PDF Viewer Area -->
          <div class="flex-1 bg-white border border-slate-200 rounded-xl flex flex-col items-center justify-center p-2 text-center min-h-[500px] shadow-sm overflow-hidden relative">
            <object v-if="selectedCv.url" :data="selectedCv.url" type="application/pdf" class="w-full h-full rounded-lg" style="min-height: 550px;">
              <div class="flex flex-col items-center justify-center h-full p-8">
                <FileText class="w-16 h-16 text-slate-300 mb-4" />
                <p class="text-slate-500 font-medium">Trình duyệt của bạn không hỗ trợ xem PDF trực tiếp.</p>
                <a :href="selectedCv.url" target="_blank" class="mt-4 text-blue-600 hover:underline">Tải xuống để xem</a>
              </div>
            </object>
            <div v-else class="flex flex-col items-center justify-center p-8">
              <FileText class="w-16 h-16 text-slate-300 mb-4" />
              <p class="text-slate-500 font-medium">Không thể tải bản xem trước tài liệu</p>
            </div>
          </div>

          <!-- Parsed Data Area -->
          <div class="w-full lg:w-[400px] space-y-6 shrink-0">
            <div>
              <h4 class="text-sm font-bold text-slate-700 uppercase tracking-wide mb-4 flex items-center gap-2">
                <Bot class="w-5 h-5" style="color: var(--accent);" /> Thông tin AI đã trích xuất
              </h4>
              <div class="bg-white p-5 rounded-xl border border-slate-200 shadow-sm space-y-4 text-sm">
                <div class="flex flex-col gap-1 border-b border-slate-100 pb-3">
                  <span class="text-slate-500 font-medium">Chức danh phù hợp:</span>
                  <span class="font-bold text-slate-800 text-base">{{ selectedCv.parsedData?.role || 'Chưa cập nhật' }}</span>
                </div>
                <div class="flex flex-col gap-1 border-b border-slate-100 pb-3">
                  <span class="text-slate-500 font-medium">Đánh giá kinh nghiệm:</span>
                  <span class="font-bold text-slate-800 text-base">{{ selectedCv.parsedData?.level || 'Chưa cập nhật' }}</span>
                </div>
                <div class="flex items-center justify-between pt-1">
                  <span class="text-slate-500 font-medium">Trạng thái phân tích:</span>
                  <span class="font-bold text-emerald-600 flex items-center gap-1.5 bg-emerald-50 px-2.5 py-1 rounded-md">
                    <CheckCircle class="w-4 h-4" /> Hoàn tất
                  </span>
                </div>
              </div>
            </div>

            <div>
              <h4 class="text-sm font-bold text-slate-700 uppercase tracking-wide mb-4">Kỹ năng phát hiện được</h4>
              <div class="bg-white p-5 rounded-xl border border-slate-200 shadow-sm flex flex-wrap gap-2.5">
                <span v-for="skill in (selectedCv.parsedData?.skills || [])" :key="skill" class="modal-skill-chip">
                  {{ skill }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 border-t border-slate-200 bg-white flex justify-end gap-3 shrink-0">
          <button @click="deleteCv(selectedCv.id); closeCvModal()" class="px-5 py-2.5 text-red-600 font-bold hover:bg-red-50 rounded-xl transition-colors">
            Xóa CV này
          </button>
          <button @click="closeCvModal" class="px-8 py-2.5 bg-blue-600 text-white font-bold hover:bg-blue-700 hover:shadow-lg hover:shadow-blue-500/30 rounded-xl transition-all">
            Đóng cửa sổ
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pf-card { padding: 0; overflow: hidden; }
.pf-card-head { padding: 18px 24px; border-bottom: 1px solid var(--border); background: var(--surface-soft); }
.pf-card-title { display: flex; align-items: center; gap: 10px; font-size: 17px; font-weight: 700; color: var(--text-main); }
.pf-sub-title { display: flex; align-items: center; gap: 10px; font-size: 16px; font-weight: 700; color: var(--text-main); margin-bottom: 22px; }
.pf-label { font-size: 13px; font-weight: 600; color: var(--text-secondary); }
.pf-select { width: 100%; height: 44px; padding: 0 14px; border: 1px solid var(--border); border-radius: var(--radius); background: var(--surface); color: var(--text-main); outline: none; transition: all 0.2s ease; font-family: var(--sans); }
.pf-select:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(37,99,235,0.12); }
.pf-textarea { width: 100%; padding: 14px; border: 1px solid var(--border); border-radius: var(--radius); outline: none; resize: none; color: var(--text-main); font-family: var(--sans); font-size: 14px; transition: all 0.2s ease; }
.pf-textarea:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(37,99,235,0.12); }

.pf-dropzone { border: 2px dashed var(--border); border-radius: var(--radius); padding: 24px; text-align: center; cursor: pointer; transition: all 0.3s ease; }
.pf-dropzone:hover { border-color: var(--primary); background: var(--primary-light); }
.pf-dz-icon { display: flex; align-items: center; justify-content: center; width: 48px; height: 48px; margin: 0 auto 12px; border-radius: 50%; background: var(--primary-light); color: var(--primary); transition: transform 0.3s ease; }
.pf-dropzone:hover .pf-dz-icon { transform: translateY(-3px); }
.pf-dz-title { font-weight: 700; color: var(--text-main); margin-bottom: 4px; }
.pf-list-title { font-size: 12px; font-weight: 700; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.04em; }

.pf-cv-item { display: flex; align-items: center; gap: 14px; padding: 14px; border-radius: var(--radius); border: 1px solid var(--border); background: var(--surface); transition: all 0.25s ease; }
.pf-cv-item.is-done { cursor: pointer; }
.pf-cv-item.is-done:hover { box-shadow: var(--shadow-md); border-color: var(--primary-light); }
.pf-cv-item.is-loading { opacity: 0.8; cursor: wait; }
.pf-cv-icon { display: flex; align-items: center; justify-content: center; width: 44px; height: 44px; border-radius: var(--radius); background: var(--primary-light); color: var(--primary); flex-shrink: 0; }
.pf-cv-icon.is-muted { background: var(--surface-soft); color: var(--text-muted); }
.pf-cv-name { font-weight: 700; color: var(--text-main); }
.pf-cv-meta { display: flex; align-items: center; gap: 5px; font-size: 12px; margin-top: 3px; font-weight: 500; }
.pf-cv-meta.is-ok { color: var(--success); }
.pf-cv-meta.is-wait { color: var(--warning); }
.pf-analyzing { display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; border-radius: var(--radius); background: rgba(217,119,6,0.12); color: var(--warning); font-size: 12px; font-weight: 600; }
.pf-icon-btn { padding: 8px; border: none; background: transparent; border-radius: var(--radius); color: var(--primary); cursor: pointer; opacity: 0; transition: all 0.2s ease; }
.pf-cv-item:hover .pf-icon-btn { opacity: 1; }
.pf-icon-btn:hover { background: var(--primary-light); }
.pf-icon-btn.is-danger { color: var(--danger); }
.pf-icon-btn.is-danger:hover { background: rgba(220,38,38,0.1); }
.pf-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.pf-note { display: flex; gap: 10px; padding: 14px 16px; }
.pf-note :deep(svg) { color: var(--accent); flex-shrink: 0; margin-top: 1px; }
.pf-note p { color: var(--text-secondary); font-size: 13px; line-height: 1.55; }

.modal-skill-chip { padding: 6px 12px; border-radius: var(--radius); background: var(--accent-bg); color: var(--accent); border: 1px solid rgba(8,145,178,0.2); font-size: 13px; font-weight: 600; }
</style>
