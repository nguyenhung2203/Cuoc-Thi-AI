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
        if (data.cv_file?.id) {
          const res = await candidateService.getCVUrl(data.cv_file.id, companyId)
          if (res && res.url) {
            cvPreviewUrl.value = res.url
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

const handleParseCV = async () => {
  if (!candidate.value.cv_file_id || isNew.value) {
    localToast.value = { type: 'error', message: 'Vui lòng lưu ứng viên với CV trước' }
    return
  }
  isAiParsing.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await candidateService.parseCV(companyId, id)

    const updated = await candidateService.getCandidate(companyId, id)
    if (updated.ai_cv_summary || updated.parsed_cv_json) {
      parsedData.value = {
        skills: updated.skills || [],
        experience: updated.experience || updated.ai_cv_summary || '',
        education: updated.education || ''
      }
    }
    localToast.value = { type: 'success', message: 'Đã phân tích CV thành công!' }
  } catch (error) {
    localToast.value = { type: 'error', message: 'Lỗi phân tích CV: ' + (error.message || '') }
  } finally {
    isAiParsing.value = false
  }
}
</script>

<template>
  <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" class="fixed top-5 right-5 z-[9999]" />
  <div v-if="loading" class="flex flex-col items-center justify-center min-h-[400px] text-slate-500 dark:text-slate-400">
    <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
    <span class="font-medium">Đang tải chi tiết ứng viên...</span>
  </div>
  <div v-else class="animate-fade-in space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center gap-4 bg-white dark:bg-slate-800 p-6 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700">
      <button @click="router.push('/candidates')" class="p-2 hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-500 dark:text-slate-400 rounded-lg transition-colors border border-transparent hover:border-slate-200 dark:hover:border-slate-600 shrink-0">
        <ArrowLeft size="20" />
      </button>
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">{{ isNew ? 'Thêm Ứng Viên Mới' : candidate.name }}</h1>
          <span v-if="!isNew" 
            class="px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider rounded-full border"
            :class="[
              candidate.status.toLowerCase() === 'new' ? 'bg-blue-50 text-blue-600 border-blue-200 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/20' : 
              candidate.status.toLowerCase() === 'interviewing' ? 'bg-amber-50 text-amber-600 border-amber-200 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/20' : 
              candidate.status.toLowerCase() === 'offered' ? 'bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20' : 
              candidate.status.toLowerCase() === 'rejected' ? 'bg-rose-50 text-rose-600 border-rose-200 dark:bg-rose-500/10 dark:text-rose-400 dark:border-rose-500/20' : 
              'bg-slate-100 text-slate-600 border-slate-200 dark:bg-slate-700 dark:text-slate-300 dark:border-slate-600'
            ]"
          >
            {{ candidate.status }}
          </span>
        </div>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">Candidates > {{ isNew ? 'New' : candidate.name }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 relative">
      <!-- Loading Overlay -->
      <div v-if="uploading || isAiParsing" class="absolute inset-0 z-50 bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm rounded-2xl flex flex-col items-center justify-center border border-indigo-100 dark:border-indigo-900 shadow-2xl">
        <div class="w-16 h-16 mb-6 rounded-full bg-indigo-100 dark:bg-indigo-900/50 flex items-center justify-center shadow-lg relative animate-bounce">
           <Sparkles size="32" class="text-indigo-600 dark:text-indigo-400" />
        </div>
        <h3 class="text-2xl font-bold text-slate-800 dark:text-slate-100 mb-2">
          {{ uploading ? 'Đang tải CV lên hệ thống...' : 'AI đang đọc và phân tích CV...' }}
        </h3>
        <p class="text-slate-500 dark:text-slate-400 font-medium max-w-sm text-center">
          Vui lòng đợi trong giây lát. Hệ thống đang trích xuất thông tin kỹ năng và kinh nghiệm từ CV của ứng viên.
        </p>
        <div v-if="uploading" class="w-64 mt-6">
          <div class="h-2 bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
            <div class="h-full bg-indigo-600 transition-all duration-300" :style="`width: ${uploadProgress}%`"></div>
          </div>
          <p class="text-center text-sm font-bold text-indigo-600 mt-2">{{ uploadProgress }}%</p>
        </div>
      </div>

      <!-- Form Section -->
      <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6 self-start">
        <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-6">Thông tin cá nhân</h2>
        <form @submit="handleSave" class="space-y-5">
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Họ và Tên</label>
            <input 
              v-model="candidate.name" 
              required 
              minlength="2"
              class="w-full px-4 py-2.5 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all text-slate-700 dark:text-slate-200 placeholder:text-slate-400"
            />
          </div>
          
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Email</label>
            <input 
              type="email"
              v-model="candidate.email" 
              required 
              class="w-full px-4 py-2.5 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all text-slate-700 dark:text-slate-200 placeholder:text-slate-400"
            />
          </div>
          
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Vị trí ứng tuyển <span class="text-red-500">*</span></label>
            <select 
              v-model="candidate.job_id" 
              required 
              :disabled="!isNew"
              class="w-full px-4 py-2.5 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all text-slate-700 dark:text-slate-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <option value="" disabled>-- Chọn vị trí ứng tuyển --</option>
              <option v-for="job in jobs" :key="job.id" :value="job.id">{{ job.title }}</option>
            </select>
            <p v-if="!isNew" class="text-xs text-slate-500 dark:text-slate-400 mt-1">Không thể thay đổi vị trí của ứng viên đã tạo</p>
          </div>
          
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Trạng thái</label>
            <select 
              v-model="candidate.status"
              class="w-full px-4 py-2.5 bg-white dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 transition-all text-slate-700 dark:text-slate-200"
            >
              <option value="new">Mới (New)</option>
              <option value="interviewing">Đang phỏng vấn (Interviewing)</option>
              <option value="offered">Đã gửi Offer (Offered)</option>
              <option value="rejected">Từ chối (Rejected)</option>
            </select>
          </div>

          <div class="flex justify-end pt-4 gap-3 border-t border-slate-100 dark:border-slate-700">
            <Button type="button" variant="ghost" @click="router.push('/candidates')" class="text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700">Hủy</Button>
            <Button type="submit" :disabled="saving" class="bg-indigo-600 hover:bg-indigo-700 text-white border-none shadow-md shadow-indigo-500/20">
              <Save size="16" class="mr-2" v-if="!saving" /> 
              <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2"></div>
              {{ saving ? 'Đang lưu...' : 'Lưu ứng viên' }}
            </Button>
          </div>
        </form>
      </div>

      <div class="space-y-6">
        <!-- CV Card -->
        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-6">Hồ sơ (CV & Resume)</h2>
          <div v-if="candidate.cv || selectedFile" class="border border-slate-200 dark:border-slate-700 rounded-xl p-4 bg-slate-50 dark:bg-slate-900">
            <div class="flex items-center gap-4 mb-3">
              <FileText size="32" class="text-indigo-500" />
              <div class="flex-1">
                <p v-if="selectedFile" class="font-medium text-slate-800 dark:text-slate-200">{{ selectedFile.name }}</p>
                <p v-else class="font-medium text-indigo-600 dark:text-indigo-400 cursor-pointer hover:underline" @click="handleViewCV" title="Nhấn để xem CV">{{ candidate.cv }}</p>
                <p class="text-emerald-600 dark:text-emerald-400 text-sm flex items-center gap-1 mt-1 font-medium">
                  <CheckCircle size="14" /> {{ selectedFile ? 'Sẵn sàng tải lên' : 'Đã lưu trên hệ thống' }}
                </p>
              </div>
              <Button type="button" variant="secondary" @click="fileInputRef?.click()" :disabled="uploading || isAiParsing" class="bg-white dark:bg-slate-800 border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-700">
                Thay đổi
              </Button>
            </div>
            
            <div v-if="cvPreviewUrl" class="w-full border border-slate-200 dark:border-slate-700 rounded-lg overflow-hidden mt-4 bg-slate-100 dark:bg-slate-900 flex items-center justify-center min-h-[400px]">
              <img v-if="selectedFile && selectedFile.type.startsWith('image/')" :src="cvPreviewUrl" class="max-w-full max-h-[400px] object-contain mx-auto block" />
              <iframe v-else-if="selectedFile && selectedFile.type === 'application/pdf'" :src="cvPreviewUrl" width="100%" height="400px" class="border-none block"></iframe>
              <iframe v-else-if="!selectedFile && cvPreviewUrl" :src="cvPreviewUrl" width="100%" height="400px" class="border-none block"></iframe>
            </div>

            <div v-if="candidate.cv_file_id && !selectedFile" class="flex gap-2 mt-4">
              <Button type="button" variant="secondary" @click="handleViewCV" class="flex-1 bg-white dark:bg-slate-800 border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-700">
                Xem CV
              </Button>
              <Button type="button" :disabled="uploading || isAiParsing" @click="handleParseCV" class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white border-none">
                <Sparkles size="16" class="mr-2" v-if="!isAiParsing" />
                {{ isAiParsing ? 'Đang phân tích...' : 'Phân tích CV' }}
              </Button>
            </div>

          </div>
          <div v-else 
            class="border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-xl p-10 text-center cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors group"
            @click="fileInputRef?.click()"
          >
            <Upload size="36" class="text-slate-400 group-hover:text-indigo-500 mx-auto mb-3 transition-colors" />
            <p class="font-semibold text-slate-700 dark:text-slate-300 mb-1">Click để tải CV lên</p>
            <p class="text-sm text-slate-500 dark:text-slate-400">Hỗ trợ PDF, DOCX, PNG, JPG (tối đa 5MB)</p>
          </div>
          <input 
            type="file" 
            ref="fileInputRef" 
            class="hidden" 
            accept=".pdf,.doc,.docx,.png,.jpg,.jpeg" 
            @change="handleFileUpload"
          />
        </div>

        <!-- AI Parsing Card -->
        <div v-if="parsedData" class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6 border-t-4 border-t-indigo-500 relative overflow-hidden">
          <div class="absolute -right-6 -top-6 text-indigo-500/10 pointer-events-none">
            <Sparkles size="100" />
          </div>
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-5 relative z-10">AI Bóc tách dữ liệu CV</h2>
          <div class="flex flex-col gap-6 relative z-10">
            <div class="flex items-center gap-2 text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-500/10 p-3 rounded-lg border border-indigo-100 dark:border-indigo-500/20">
              <Sparkles size="18" /> <span class="text-sm font-semibold">Hoàn tất phân tích tự động</span>
            </div>
            
            <div>
              <h4 class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-3">Kỹ năng nổi bật</h4>
              <div class="flex flex-wrap gap-2">
                <span v-for="s in parsedData.skills" :key="s" class="px-3 py-1 rounded-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300 text-xs font-semibold">
                  {{ s }}
                </span>
              </div>
            </div>
            
            <div>
              <h4 class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2">Kinh nghiệm</h4>
              <p class="text-sm text-slate-700 dark:text-slate-300 leading-relaxed">{{ parsedData.experience }}</p>
            </div>
            
            <div>
              <h4 class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2">Học vấn</h4>
              <p class="text-sm text-slate-700 dark:text-slate-300 leading-relaxed">{{ parsedData.education }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
