<script setup>
import { ref, onMounted } from 'vue'
import { authStore } from '../../stores/auth.store'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import { Upload, FileText, CheckCircle, Save, Trash2, Eye, X, Bot, Loader2 } from 'lucide-vue-next'
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
  <div class="space-y-8 animate-fade-in pb-12 max-w-6xl mx-auto">
    <div class="mb-8">
      <h1 class="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-indigo-600 to-purple-600">Hồ sơ cá nhân & CV</h1>
      <p class="text-gray-500 mt-2 text-lg">Cập nhật thông tin để AI có thể đưa ra bài luyện tập chính xác nhất.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      
      <!-- Left Column: Forms -->
      <div class="lg:col-span-2 space-y-8">
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl overflow-hidden">
          <div class="bg-gradient-to-r from-indigo-50 to-purple-50 p-6 border-b border-indigo-100">
            <h3 class="text-xl font-bold text-indigo-900 flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
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
                <h3 class="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                  Định hướng nghề nghiệp
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <Input label="Vị trí mục tiêu (Target Role)" v-model="profile.targetRole" />
                  <div class="flex flex-col gap-2">
                    <label class="text-sm font-semibold text-gray-700">Cấp độ hiện tại</label>
                    <select class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 outline-none transition-all bg-white text-gray-800" v-model="profile.level">
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
                <button type="submit" :disabled="saving" class="bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-bold py-3 px-8 rounded-xl shadow-lg hover:shadow-indigo-500/30 transition-all duration-300 transform hover:-translate-y-0.5 disabled:opacity-50 disabled:transform-none flex items-center gap-2">
                  <Save class="w-5 h-5" /> 
                  {{ saving ? 'Đang lưu...' : 'Lưu hồ sơ' }}
                </button>
              </div>
            </form>
          </div>
        </Card>

        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl overflow-hidden">
          <div class="bg-gradient-to-r from-emerald-50 to-teal-50 p-6 border-b border-emerald-100">
            <h3 class="text-xl font-bold text-emerald-900 flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" /></svg>
              Kỹ năng chuyên môn
            </h3>
          </div>
          <div class="p-6">
            <textarea class="w-full p-4 border border-gray-200 rounded-xl focus:border-emerald-500 focus:ring-4 focus:ring-emerald-500/20 outline-none transition-all text-gray-800 font-medium resize-none" rows="4" placeholder="Ví dụ: ReactJS, NodeJS, TypeScript..." v-model="profile.skills"></textarea>
            <p class="text-sm text-gray-500 mt-2">Phân cách các kỹ năng bằng dấu phẩy (,)</p>
          </div>
        </Card>
      </div>

      <!-- Right Column: CV Upload -->
      <div class="space-y-8">
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl overflow-hidden sticky top-6">
          <div class="bg-gradient-to-r from-blue-50 to-cyan-50 p-6 border-b border-blue-100">
            <h3 class="text-xl font-bold text-blue-900 flex items-center gap-2">
              <FileText class="w-6 h-6 text-blue-600" />
              CV của bạn
            </h3>
          </div>
          
          <div class="p-6 space-y-4">
            <!-- Hidden File Input -->
            <input type="file" ref="fileInput" @change="handleFileUpload" accept=".pdf,.doc,.docx" multiple class="hidden" style="display: none;" />

            <!-- Upload Area -->
            <div class="group border-2 border-dashed border-gray-200 hover:border-blue-400 bg-gray-50 hover:bg-blue-50/50 rounded-xl p-6 text-center cursor-pointer transition-all duration-300 relative overflow-hidden" @click="fileInput.click()">
              <div class="absolute inset-0 bg-gradient-to-br from-blue-400/0 to-indigo-400/0 group-hover:from-blue-400/10 group-hover:to-indigo-400/10 transition-all duration-500"></div>
              <div class="w-12 h-12 bg-white rounded-full shadow-sm flex items-center justify-center mx-auto mb-3 group-hover:-translate-y-2 group-hover:shadow-md transition-all duration-300">
                <Upload class="w-6 h-6 text-blue-500" />
              </div>
              <p class="text-gray-900 font-bold mb-1 group-hover:text-blue-700 transition-colors">Tải CV lên (Nhiều file)</p>
              <p class="text-gray-500 text-sm font-medium">PDF, DOCX (Tối đa 5MB/file)</p>
            </div>

            <!-- List of Current CVs -->
            <div v-if="uploadedCvs.length > 0" class="space-y-3 mt-6">
              <h4 class="text-sm font-bold text-gray-700 uppercase tracking-wide">Danh sách CV đã tải lên</h4>
              
              <div v-for="cv in uploadedCvs" :key="cv.id" @click="cv.status === 'done' && viewCv(cv)" class="bg-blue-50/40 p-4 rounded-xl border border-blue-100 flex items-center gap-4 transition-all group" :class="cv.status === 'done' ? 'hover:shadow-md hover:bg-blue-50/80 cursor-pointer' : 'opacity-80 cursor-wait'">
                <div class="w-12 h-12 bg-white rounded-lg shadow-sm flex items-center justify-center flex-shrink-0" :class="cv.status === 'done' ? 'group-hover:scale-105 transition-transform' : ''">
                  <FileText class="w-6 h-6" :class="cv.status === 'done' ? 'text-blue-600' : 'text-slate-400'" />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="font-bold text-gray-800 truncate">{{ cv.name }}</p>
                  <p v-if="cv.status === 'done'" class="text-xs text-emerald-600 flex items-center gap-1 mt-1 font-medium">
                    <CheckCircle class="w-3 h-3" /> Tải lên: {{ cv.date }} • {{ cv.size }}
                  </p>
                  <p v-else class="text-xs text-amber-600 flex items-center gap-1 mt-1 font-medium">
                    AI đang đọc thông tin...
                  </p>
                </div>
                
                <div class="flex gap-2">
                  <div v-if="cv.status === 'analyzing'" class="px-3 py-1.5 bg-amber-100/50 text-amber-700 rounded-lg font-medium text-xs flex items-center gap-2">
                    <Loader2 class="w-4 h-4 animate-spin" /> Đang phân tích
                  </div>
                  <template v-else>
                    <button type="button" @click.stop="viewCv(cv)" class="p-2 text-blue-600 hover:bg-blue-100 rounded-lg transition-colors opacity-0 group-hover:opacity-100" title="Xem chi tiết">
                      <Eye class="w-5 h-5" />
                    </button>
                    <button type="button" @click.stop="deleteCv(cv.id)" class="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors opacity-0 group-hover:opacity-100" title="Xóa CV">
                      <Trash2 class="w-5 h-5" />
                    </button>
                  </template>
                </div>
              </div>
            </div>
            
            <div class="bg-amber-50 border border-amber-100 rounded-xl p-4 flex gap-3">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-amber-500 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              <p class="text-amber-800 text-sm leading-relaxed font-medium">
                CV của bạn sẽ được AI dùng làm ngữ cảnh (context) để đặt câu hỏi sát với kinh nghiệm thực tế của bạn trong phần Luyện phỏng vấn Mock Interview.
              </p>
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
                <Bot class="w-5 h-5 text-indigo-500" /> Thông tin AI đã trích xuất
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
                <span v-for="skill in (selectedCv.parsedData?.skills || [])" :key="skill" class="px-3 py-1.5 bg-indigo-50 text-indigo-700 rounded-lg text-sm font-semibold border border-indigo-100">
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
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in {
  animation: fade-in 0.5s ease-out forwards;
}
</style>
