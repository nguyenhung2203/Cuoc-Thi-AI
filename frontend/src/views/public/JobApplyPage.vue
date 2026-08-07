<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { publicService } from '../../services/public.service'
import { validateFile } from '../../utils/validators.js'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { authStore } from '../../stores/auth.store'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { Briefcase, MapPin, Clock, ArrowLeft, UploadCloud, Sparkles, Target, AlertCircle, BookOpen, Building2, Globe, Users, AlignLeft, CheckSquare, Gift, DollarSign, ExternalLink, Award } from 'lucide-vue-next'
import { loadSavedJobIds, toggleSavedJob } from '../../utils/savedJobs'
import { formatSalaryTrieu, formatExperience } from '../../utils/formatters'

const route = useRoute()
const router = useRouter()
const companyId = route.params.company_id
const jobId = route.params.job_id
const job = ref(null)
const company = ref(null)
const loading = ref(true)
const submitting = ref(false)
const cvError = ref('')
const toast = ref(null)

const salaryLabel = computed(() => formatSalaryTrieu(job.value?.salary_min, job.value?.salary_max))
const experienceLabel = computed(() => formatExperience(job.value?.level))

const cvFile = ref(null)
const cvPreviewUrl = ref(null)
const isDragging = ref(false)
const fileInput = ref(null)

const isLoggedIn = computed(() => authStore.isAuthenticated)
const user = computed(() => authStore.user)

// AI CV↔JD match (real, from backend). null = not loaded, {has_cv:false} = no CV yet.
const match = ref(null)
const matchLoading = ref(false)

const loadMatch = async () => {
  if (!isLoggedIn.value) return
  matchLoading.value = true
  try {
    const res = await candidatePortalService.getJobMatch(jobId)
    match.value = res.data || res
  } catch (e) {
    console.error('Failed to load AI match', e)
    match.value = null
  } finally {
    matchLoading.value = false
  }
}

const unwrap = (val) => {
  if (!val) return ''
  if (typeof val === 'object') {
    if ('String' in val) return val.Valid ? val.String : ''
    if ('Int64' in val) return val.Valid ? val.Int64 : ''
    if ('Float64' in val) return val.Valid ? val.Float64 : ''
    if ('Bool' in val) return val.Valid ? val.Bool : false
  }
  return val
}

onMounted(async () => {
  try {
    // If auth store exists but not initialized, let's init it
    if (authStore.isAuthenticated && !authStore.user) {
      await authStore.init()
    }

    const [res, compRes] = await Promise.all([
      publicService.getJobDetails(companyId, jobId),
      publicService.getCompanyDetails(companyId)
    ])
    
    const rawJob = res.data || res
    company.value = compRes.data || compRes
    job.value = {
      ...rawJob,
      location: unwrap(rawJob.location),
      employment_type: unwrap(rawJob.employment_type),
      department: unwrap(rawJob.department),
      requirements: unwrap(rawJob.requirements),
      benefits: unwrap(rawJob.benefits),
      salary_min: unwrap(rawJob.salary_min),
      salary_max: unwrap(rawJob.salary_max),
      currency: unwrap(rawJob.currency),
      level: unwrap(rawJob.level),
      salary_type: unwrap(rawJob.salary_type)
    }
  } catch (error) {
    console.error('Failed to load job', error)
    toast.value = { type: 'error', message: 'Không tìm thấy thông tin công việc.' }
  } finally {
    loading.value = false
  }
  refreshSavedIds()
  loadMatch()
})

const validateAndSetFile = (file) => {
  if (!file) return
  const error = validateFile(file, {
    required: true,
    maxBytes: 5 * 1024 * 1024,
    mimeTypes: ['application/pdf', 'image/png', 'image/jpeg'],
    extensions: ['pdf', 'png', 'jpg', 'jpeg'],
    typeMessage: 'Vui lòng chọn file PDF, PNG hoặc JPG.',
    sizeMessage: 'CV không được vượt quá 5MB.',
  })
  cvError.value = error
  if (error) {
    cvFile.value = null
    if (cvPreviewUrl.value) URL.revokeObjectURL(cvPreviewUrl.value)
    cvPreviewUrl.value = null
    toast.value = { type: 'error', message: error }
    return
  }

  cvFile.value = file
  if (cvPreviewUrl.value) URL.revokeObjectURL(cvPreviewUrl.value)
  cvPreviewUrl.value = URL.createObjectURL(file)
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
  if (submitting.value) return
  if (!isLoggedIn.value) {
    requireLogin()
    return
  }

  if (job.value?.status && job.value.status !== 'open') {
    toast.value = { type: 'error', message: 'Tin tuyển dụng này đã đóng hoặc tạm dừng, không thể ứng tuyển.' }
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
      router.push('/my-applications')
    }, 1500)
  } catch (error) {
    const msg = String(error?.message || error?.code || '').toLowerCase()
    if (msg.includes('already applied') || msg.includes('conflict') || error?.code === 'CONFLICT') {
      toast.value = { type: 'warning', message: 'Bạn đã ứng tuyển vị trí này rồi.' }
    } else if (msg.includes('not open') || msg.includes('closed')) {
      toast.value = { type: 'error', message: 'Tin tuyển dụng không còn nhận hồ sơ.' }
    } else {
      toast.value = { type: 'error', message: error.message || 'Có lỗi xảy ra khi nộp đơn.' }
    }
  } finally {
    submitting.value = false
  }
}

const savedJobIds = ref([])
const refreshSavedIds = () => {
  savedJobIds.value = loadSavedJobIds()
}

const toggleSaveJob = () => {
  if (!job.value) return
  const snapshot = {
    ...job.value,
    company_id: companyId,
    company_name: company.value?.name || job.value.company_name || 'Chưa cập nhật',
  }
  const { list, saved } = toggleSavedJob(snapshot)
  savedJobIds.value = list.map((j) => j.id)
  toast.value = saved
    ? { type: 'success', message: 'Đã lưu tin tuyển dụng trên thiết bị này.' }
    : { type: 'info', message: 'Đã bỏ lưu tin tuyển dụng.' }
}

const showBeginnerGuide = ref(false)

const beginnerSkills = computed(() => {
  if (!job.value) return []
  const text = ((job.value.requirements || '') + ' ' + (job.value.description || '')).toLowerCase()
  const skillMap = [
    { key: 'javascript', label: 'JavaScript / ES6+', doc: 'Nền tảng lập trình web, xử lý logic phía client' },
    { key: 'vue', label: 'Vue.js Framework', doc: 'Framework xây dựng giao diện ứng dụng web reactive' },
    { key: 'react', label: 'React.js', doc: 'Thư viện UI phổ biến xây dựng component-based' },
    { key: 'node', label: 'Node.js / Express', doc: 'Lập trình Backend bằng JavaScript' },
    { key: 'go', label: 'Golang (Go)', doc: 'Ngôn ngữ backend hiệu năng cao, microservices' },
    { key: 'python', label: 'Python / FastAPI / Django', doc: 'Lập trình backend, phân tích dữ liệu & AI' },
    { key: 'sql', label: 'SQL / PostgreSQL / MySQL', doc: 'Truy vấn và thiết kế cơ sở dữ liệu quan hệ' },
    { key: 'docker', label: 'Docker & Containerization', doc: 'Đóng gói và triển khai ứng dụng' },
    { key: 'git', label: 'Git & GitHub/GitLab', doc: 'Quản lý mã nguồn & làm việc nhóm' },
    { key: 'api', label: 'RESTful API / JSON', doc: 'Giao tiếp dữ liệu giữa Frontend và Backend' },
    { key: 'css', label: 'CSS3 / Tailwind / Responsive', doc: 'Thiết kế giao diện đẹp và tương thích mobile' },
  ]
  const matched = skillMap.filter(item => text.includes(item.key))
  if (matched.length > 0) return matched
  // Fallback defaults if no specific keyword found
  return [
    { label: 'Kỹ năng chuyên môn ngành ' + (job.value.department || 'IT'), doc: 'Tìm hiểu các khái niệm cơ bản và quy trình làm việc thực tế' },
    { label: 'Tư duy giải quyết vấn đề (Problem Solving)', doc: 'Cách phân tích bài toán và tìm giải pháp tối ưu' },
    { label: 'Sử dụng Git & Quy trình làm việc nhóm', doc: 'Quản lý phiên bản mã nguồn và phối hợp làm việc' }
  ]
})

const scrollToApply = () => {
  const applySection = document.getElementById('apply-section')
  if (applySection) {
    applySection.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<template>
  <div class="job-apply-page">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <Button variant="ghost" class="btn-back" @click="router.push('/job-board')">
      <ArrowLeft size="16" style="margin-right: 8px;" /> Quay lại danh sách việc làm
    </Button>

    <div v-if="loading" style="text-align: center; padding: 40px;">
      <div class="spinner"></div>
    </div>

    <div v-else-if="job" class="job-layout-grid">
      <!-- Left: Job Details -->
      <div class="left-col">
        
        <!-- Job Header (TopCV Style) -->
        <Card class="top-card animate-rise">
          <h1 class="job-title-top">{{ job.title }}</h1>
          
          <div class="job-quick-stats">
            <div class="stat-item">
              <div class="stat-icon text-primary"><DollarSign size="20" /></div>
              <div>
                <div class="stat-label">Mức lương</div>
                <div class="stat-value text-salary">{{ salaryLabel }}</div>
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-icon text-primary"><MapPin size="20" /></div>
              <div>
                <div class="stat-label">Địa điểm</div>
                <div class="stat-value">{{ job.location || 'Bất kỳ' }}</div>
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-icon text-primary"><Briefcase size="20" /></div>
              <div>
                <div class="stat-label">Kinh nghiệm</div>
                <div class="stat-value">{{ experienceLabel }}</div>
              </div>
            </div>
          </div>
          
          <div class="action-buttons">
            <Button variant="primary" class="btn-apply-big" @click="scrollToApply">
              <UploadCloud size="18" style="margin-right: 8px" /> Ứng tuyển ngay
            </Button>
            <Button 
              variant="outline" 
              class="btn-save" 
              :class="{ 'is-saved': savedJobIds.includes(job.id) }" 
              @click="toggleSaveJob"
            >
              {{ savedJobIds.includes(job.id) ? 'Đã lưu' : 'Lưu tin' }}
            </Button>
          </div>
        </Card>

        <!-- AI Career Coach Panel -->
        <div v-if="isLoggedIn" class="ai-panel animate-rise" style="animation-delay: 0.1s;">
          <div class="ai-bg-icon">
            <Sparkles :size="80" />
          </div>
          <div class="ai-panel-header">
            <div style="display: flex; align-items: center; gap: 8px;">
              <Sparkles :size="24" /> AI Career Coach
            </div>
            <div v-if="match && match.has_cv" class="ai-score-ring" :style="`--score: ${Math.round(match.fit_score)};`">
              <div class="ai-score-inner">
                <span class="ai-score-val">{{ Math.round(match.fit_score) }}%</span>
                <span class="ai-score-label">Phù hợp</span>
              </div>
            </div>
          </div>

          <!-- Loading state -->
          <div v-if="matchLoading" class="ai-section">
            <div class="ai-summary-box">
              <p>Đang phân tích mức độ phù hợp giữa CV của bạn và công việc...</p>
            </div>
          </div>

          <!-- No CV: honest CTA instead of fake numbers -->
          <div v-else-if="match && !match.has_cv" class="ai-section">
            <div class="ai-summary-box">
              <p><strong>Chưa có dữ liệu phù hợp.</strong> Hãy tải lên CV của bạn để AI phân tích mức độ phù hợp với công việc này.</p>
            </div>
          </div>

          <!-- Real AI match result -->
          <template v-else-if="match && match.has_cv">
            <div class="ai-section">
              <h4 class="ai-section-title">
                <Target :size="18" class="text-info" /> Nhận định của AI
              </h4>
              <div class="ai-summary-box">
                <p>{{ match.summary }}</p>
                <p v-if="match.recommendation"><strong>Gợi ý:</strong> {{ match.recommendation }}</p>
              </div>
            </div>

            <div class="ai-grid">
              <!-- Matched skills -->
              <div class="ai-card">
                <h4 class="ai-section-title">
                  <CheckSquare :size="16" class="text-success" /> Kỹ năng phù hợp
                </h4>
                <ul class="ai-list">
                  <li v-for="(sk, i) in match.matched_skills" :key="'m'+i" class="ai-list-item text-success">
                    <span class="ai-icon">✓</span>
                    <span>{{ sk }}</span>
                  </li>
                  <li v-if="!match.matched_skills || match.matched_skills.length === 0" class="ai-list-item text-muted">
                    <span>Chưa xác định được kỹ năng khớp.</span>
                  </li>
                </ul>
              </div>

              <!-- Missing skills -->
              <div class="ai-card">
                <h4 class="ai-section-title">
                  <AlertCircle :size="16" class="text-warning" /> Kỹ năng cần bổ sung
                </h4>
                <ul class="ai-list">
                  <li v-for="(sk, i) in match.missing_skills" :key="'x'+i" class="ai-list-item text-warning">
                    <span class="ai-icon">!</span>
                    <span>{{ sk }}</span>
                  </li>
                  <li v-if="!match.missing_skills || match.missing_skills.length === 0" class="ai-list-item text-success">
                    <span class="ai-icon">✓</span>
                    <span>Không thiếu kỹ năng nổi bật.</span>
                  </li>
                </ul>
              </div>
            </div>
          </template>
        </div>

        <!-- JD Content (TopCV Style) -->
        <Card class="jd-content-card animate-rise" style="animation-delay: 0.15s;">
          <div class="flex items-center justify-between flex-wrap gap-3 mb-6 pb-4 border-b border-[var(--border)]">
            <h2 class="section-heading mb-0" style="margin-bottom: 0">Chi tiết tin tuyển dụng</h2>
            
            <!-- AI Guide Button for Beginners -->
            <button 
              @click="showBeginnerGuide = !showBeginnerGuide" 
              class="flex items-center gap-2 text-sm font-bold px-4 py-2 rounded-full border transition-all cursor-pointer shadow-sm"
              :class="showBeginnerGuide ? 'bg-primary text-white border-primary' : 'bg-[var(--accent-bg)] text-[var(--accent)] border-[var(--accent)] hover:opacity-90'">
              <Sparkles :size="16" />
              <span>{{ showBeginnerGuide ? 'Ẩn Hướng dẫn AI' : '🤖 Người mới chưa hiểu JD? Bấm AI Giải thích & Gợi ý học' }}</span>
            </button>
          </div>

          <!-- AI Beginner Explanation & Roadmap Panel -->
          <div v-if="showBeginnerGuide" class="p-5 mb-6 rounded-xl border border-[var(--accent)]/30 bg-[var(--surface-soft)] space-y-4 animate-fadeIn">
            <div class="flex items-center gap-2 font-bold text-[var(--accent)] text-base">
              <BookOpen :size="20" /> Trợ lý AI: Tóm tắt công việc & Hướng dẫn học tập cho người mới
            </div>
            
            <div class="text-sm leading-relaxed text-[var(--text-secondary)]">
              <p class="mb-2"><strong>💡 Bạn là người mới bắt đầu?</strong> Đừng lo lắng nếu mô tả công việc chứa nhiều thuật ngữ chuyên môn. AI đã phân tích JD này và đúc kết các kỹ năng trọng tâm bạn cần chuẩn bị:</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-3 pt-2">
              <div v-for="(item, idx) in beginnerSkills" :key="idx" class="p-3 bg-[var(--surface)] border border-[var(--border)] rounded-lg flex flex-col justify-between">
                <div>
                  <div class="font-bold text-sm text-[var(--text-main)] flex items-center gap-2">
                    <span class="w-6 h-6 rounded-full bg-[var(--primary-light)] text-[var(--primary)] text-xs flex items-center justify-center font-extrabold">{{ idx + 1 }}</span>
                    {{ item.label }}
                  </div>
                  <div class="text-xs text-[var(--text-secondary)] mt-1.5 leading-snug">{{ item.doc }}</div>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between flex-wrap gap-3 pt-3 border-t border-[var(--border)] mt-2">
              <span class="text-xs text-[var(--text-muted)]">🎯 Hãy tự tin luyện phỏng vấn thử với AI để kiểm tra mức độ sẵn sàng!</span>
              <button @click="router.push('/mock-setup')" class="text-xs font-bold text-[var(--primary)] hover:underline flex items-center gap-1">
                Luyện phỏng vấn thử ngay <Sparkles :size="14" />
              </button>
            </div>
          </div>
          
          <div class="jd-block">
            <h3 class="jd-block-title">Mô tả công việc</h3>
            <div class="jd-text">{{ job.description }}</div>
          </div>

          <div class="jd-block" v-if="job.requirements">
            <h3 class="jd-block-title">Yêu cầu ứng viên</h3>
            <div class="jd-text">{{ job.requirements }}</div>
          </div>

          <div class="jd-block" v-if="job.benefits">
            <h3 class="jd-block-title">Quyền lợi</h3>
            <div class="jd-text">{{ job.benefits }}</div>
          </div>
          
          <div class="jd-block">
            <h3 class="jd-block-title">Địa điểm làm việc</h3>
            <div class="jd-text">- {{ job.location || 'Bất kỳ' }}</div>
          </div>
          
          <div class="mt-8 pt-6 border-t border-[var(--border)] text-center">
            <p class="text-sm text-secondary mb-4">Hạn nộp hồ sơ: Kể từ ngày đăng tin đến khi tuyển đủ</p>
            <Button variant="primary" @click="scrollToApply">
              Ứng tuyển ngay
            </Button>
          </div>
        </Card>
      </div>

      <!-- Right: Sidebar -->
      <div class="right-col">
        <!-- Company Sidebar Card -->
        <Card class="sidebar-card company-card animate-rise" v-if="company">
          <div class="company-head">
            <div class="company-logo-sm">
              <img :src="unwrap(company.logo_url) || '/images/logo.png'" alt="Company Logo" />
            </div>
            <h3 class="company-name-sm">{{ company.name }}</h3>
          </div>
          
          <div class="company-info-list">
            <div class="info-row">
              <Users size="16" class="text-secondary shrink-0"/> 
              <span class="info-val"><strong>Quy mô:</strong> {{ unwrap(company.size) || '50-100 nhân viên' }}</span>
            </div>
            <div class="info-row">
              <Building2 size="16" class="text-secondary shrink-0"/> 
              <span class="info-val"><strong>Lĩnh vực:</strong> {{ unwrap(company.industry) || 'IT / Công nghệ' }}</span>
            </div>
            <div class="info-row">
              <MapPin size="16" class="text-secondary shrink-0"/> 
              <span class="info-val"><strong>Địa điểm:</strong> Trụ sở chính</span>
            </div>
          </div>
          
          <div class="company-link">
            <a :href="`/careers/${company.id}`" class="text-primary font-medium text-sm flex items-center justify-center gap-1 hover:underline">
              Xem trang công ty <ExternalLink size="14"/>
            </a>
          </div>
        </Card>

        <!-- General Info Sidebar -->
        <Card class="sidebar-card animate-rise" style="animation-delay: 0.1s;">
          <h3 class="sidebar-heading">Thông tin chung</h3>
          <div class="general-info-grid">
            <div class="gen-info-item">
              <div class="gen-icon"><Award size="18" /></div>
              <div>
                <div class="gen-label">Kinh nghiệm</div>
                <div class="gen-val">{{ experienceLabel }}</div>
              </div>
            </div>
            <div class="gen-info-item">
              <div class="gen-icon"><Briefcase size="18" /></div>
              <div>
                <div class="gen-label">Hình thức</div>
                <div class="gen-val">{{ job.employment_type || 'Full-time' }}</div>
              </div>
            </div>
            <div class="gen-info-item">
              <div class="gen-icon"><Target size="18" /></div>
              <div>
                <div class="gen-label">Phòng ban</div>
                <div class="gen-val">{{ job.department || 'General' }}</div>
              </div>
            </div>
          </div>
        </Card>

        <!-- Application Form -->
        <Card class="sidebar-card sticky-apply animate-rise" style="animation-delay: 0.2s;" id="apply-section">
          <h2 class="apply-title">Nộp đơn ứng tuyển</h2>

          <div v-if="!isLoggedIn" class="apply-login-prompt">
            <p>Bạn cần đăng nhập bằng tài khoản Ứng viên để nộp đơn.</p>
            <Button variant="primary" style="width: 100%" @click="requireLogin">Đăng nhập / Đăng ký</Button>
          </div>

          <div v-else>
            <div class="candidate-info-box">
              <div class="candidate-label">Thông tin ứng viên:</div>
              <div class="candidate-name">{{ user?.full_name }}</div>
              <div class="candidate-email">{{ user?.email }}</div>
            </div>

            <div class="upload-section">
              <label class="upload-label">Tải lên CV (PDF, PNG, JPG)</label>
              
              <div 
                class="upload-zone"
                :class="{ 'dragging': isDragging }"
                @dragover.prevent="isDragging = true"
                @dragleave.prevent="isDragging = false"
                @drop.prevent="handleDrop"
                @click="!cvFile && fileInput.click()"
              >
                <UploadCloud v-if="!cvFile" size="32" class="text-primary mb-3" />
                
                <div v-if="cvFile" class="file-preview-box">
                  <div class="file-name">{{ cvFile.name }}</div>
                  
                  <div v-if="cvFile.type.startsWith('image/')" class="img-preview">
                    <img :src="cvPreviewUrl" />
                  </div>
                  <div v-else-if="cvFile.type === 'application/pdf'" class="pdf-preview">
                    <iframe :src="cvPreviewUrl"></iframe>
                  </div>
                  
                  <Button variant="secondary" size="sm" class="mt-3" @click.stop="cvFile = null; cvPreviewUrl = null">
                    Xóa / Chọn file khác
                  </Button>
                </div>
                <div v-else>
                  <span class="font-medium text-primary">Bấm để tải lên</span> hoặc kéo thả file vào đây
                  <div class="text-xs text-muted mt-1">Hỗ trợ PDF, PNG, JPG (Max 10MB)</div>
                </div>
              </div>
              <input type="file" ref="fileInput" accept=".pdf,.png,.jpg,.jpeg" style="display: none" @change="handleFileUpload" />
          <p v-if="cvError" class="error-text mt-2">{{ cvError }}</p>
            </div>

            <Button variant="primary" style="width: 100%" :loading="submitting" @click="submitApplication">
              Nộp đơn ngay
            </Button>
            
            <div class="mock-cta">
              <h4 class="mock-cta-title">Chưa tự tin ứng tuyển?</h4>
              <p class="mock-cta-desc">Hãy để AI đóng vai nhà tuyển dụng và hỏi thử bạn dựa trên chính JD này. Sẵn sàng chưa?</p>
              <Button variant="outline" class="mock-cta-btn" @click="router.push('/mock-setup')">
                <Sparkles size="16" /> Luyện phỏng vấn thử
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.job-apply-page {
  padding-bottom: 64px;
}
.btn-back {
  margin-bottom: 24px;
  color: var(--text-secondary);
}

/* Layout Grid */
.job-layout-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
  align-items: start;
}
@media (min-width: 992px) {
  .job-layout-grid {
    grid-template-columns: 2fr 1fr;
  }
}
.left-col {
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.right-col {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Top Job Header Card */
.top-card {
  padding: 32px;
}
.job-title-top {
  font-size: 26px;
  font-weight: 700;
  color: var(--text-main);
  margin: 0 0 24px 0;
  line-height: 1.3;
}
.job-quick-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  margin-bottom: 32px;
}
.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
}
.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--primary-light, #DBEAFE);
  display: flex;
  align-items: center;
  justify-content: center;
}
.stat-label {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 2px;
}
.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-main);
}
.stat-value.text-salary {
  color: #e11d48;
}
.action-buttons {
  display: flex;
  gap: 16px;
}
.btn-apply-big {
  flex: 1;
  max-width: 300px;
  height: 48px;
  font-size: 15px;
  font-weight: 600;
}
.btn-save {
  height: 48px;
  padding: 0 24px;
  font-weight: 600;
}
.btn-save.is-saved {
  border-color: var(--danger);
  color: var(--danger);
  background-color: rgba(239, 68, 68, 0.05);
}

/* JD Details Section */
.jd-content-card {
  padding: 32px;
}
.section-heading {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-main);
  margin: 0 0 24px 0;
  padding-left: 14px;
  border-left: 4px solid var(--primary);
}
.jd-block {
  margin-bottom: 28px;
}
.jd-block:last-child {
  margin-bottom: 0;
}
.jd-block-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 12px;
}
.jd-text {
  font-size: 15px;
  color: var(--text-secondary);
  line-height: 1.8;
  white-space: pre-wrap;
}

/* Sidebar Cards */
.sidebar-card {
  padding: 24px;
}
.sidebar-heading {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-main);
  margin: 0 0 16px 0;
}

/* Company Sidebar Info */
.company-head {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}
.company-logo-sm {
  width: 64px;
  height: 64px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  padding: 6px;
  background: var(--surface);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.company-logo-sm img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.company-name-sm {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-main);
  margin: 0;
  line-height: 1.4;
}
.company-info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.info-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
}
.info-val strong {
  color: var(--text-main);
  font-weight: 600;
}
.company-link {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px dashed var(--border);
  text-align: center;
}

/* General Info */
.general-info-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.gen-info-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}
.gen-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--surface-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary);
  flex-shrink: 0;
}
.gen-label {
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 2px;
}
.gen-val {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-main);
}

/* Apply Section */
.sticky-apply {
  position: sticky;
  top: 24px;
}
.apply-title {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 20px 0;
  color: var(--text-main);
}
.apply-login-prompt {
  text-align: center;
  padding: 24px 0;
}
.apply-login-prompt p {
  margin-bottom: 16px;
  color: var(--text-secondary);
  font-size: 14px;
}
.candidate-info-box {
  background: var(--surface-soft);
  padding: 16px;
  border-radius: var(--radius);
  margin-bottom: 20px;
}
.candidate-label {
  font-weight: 600;
  font-size: 13px;
  margin-bottom: 4px;
  color: var(--text-secondary);
}
.candidate-name {
  color: var(--text-main);
  font-weight: 600;
  font-size: 15px;
}
.candidate-email {
  color: var(--text-muted);
  font-size: 13px;
}
.upload-section {
  margin-bottom: 24px;
}
.upload-label {
  display: block;
  font-weight: 600;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-main);
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
.file-preview-box {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.file-name {
  font-weight: 600;
  color: var(--primary);
  margin-bottom: 12px;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.img-preview {
  width: 100%;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}
.img-preview img {
  max-width: 100%;
  max-height: 400px;
  object-fit: contain;
  display: block;
  margin: 0 auto;
}
.pdf-preview {
  width: 100%;
  height: 400px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  overflow: hidden;
}
.pdf-preview iframe {
  width: 100%;
  height: 100%;
  border: none;
}

/* AI Panel (Copied from previous design but tweaked to match) */
.ai-panel { background: var(--accent-bg); border: 1px solid rgba(8,145,178,0.2); border-radius: var(--radius-lg); padding: 28px; position: relative; overflow: hidden; }
.ai-bg-icon { position: absolute; top: 0; right: 0; padding: 16px; opacity: 0.1; pointer-events: none; color: var(--accent); }
.ai-panel-header { display: flex; align-items: center; gap: 8px; justify-content: space-between; margin-bottom: 24px; color: var(--accent); font-weight: 700; font-size: 18px; position: relative; z-index: 10; }
.ai-section { margin-bottom: 24px; position: relative; z-index: 10; }
.ai-section-title { display: flex; align-items: center; gap: 8px; font-weight: 700; color: var(--text-main); font-size: 15px; margin-bottom: 12px; }
.text-info { color: var(--info); }
.text-warning { color: var(--warning); }
.text-success { color: var(--success); }
.text-danger { color: var(--danger); }
.ai-summary-box { font-size: 14px; color: var(--text-main); line-height: 1.7; background: var(--surface); padding: 16px; border-radius: var(--radius); border: 1px solid var(--border); box-shadow: var(--shadow-sm); }
.ai-grid { display: grid; grid-template-columns: 1fr; gap: 20px; position: relative; z-index: 10; }
@media (min-width: 1280px) { .ai-grid { grid-template-columns: repeat(2, 1fr); } }
.ai-card { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); padding: 20px; box-shadow: var(--shadow-sm); }
.ai-list { display: flex; flex-direction: column; gap: 12px; font-size: 14px; }
.ai-list-item { display: flex; align-items: flex-start; gap: 8px; }
.ai-icon { font-weight: 700; margin-top: 2px; }
.ai-upskill-list { display: flex; flex-direction: column; gap: 8px; }
.ai-accordion { border-bottom: 1px solid var(--border); padding-bottom: 8px; margin-bottom: 8px; }
.ai-accordion:last-child { border-bottom: none; padding-bottom: 0; margin-bottom: 0; }
.ai-accordion-title { font-weight: 600; color: var(--text-main); font-size: 14px; cursor: pointer; list-style: none; display: flex; justify-content: space-between; align-items: center; }
.ai-accordion-title::-webkit-details-marker { display: none; }
.ai-accordion-title::after { content: '+'; font-weight: 400; color: var(--text-muted); font-size: 16px; }
details[open] .ai-accordion-title::after { content: '-'; }
.ai-accordion-content { color: var(--text-secondary); font-size: 13px; margin-top: 8px; line-height: 1.6; }
.ai-score-ring { width: 56px; height: 56px; border-radius: 50%; background: conic-gradient(var(--success) calc(var(--score) * 1%), var(--border) 0); display: flex; align-items: center; justify-content: center; position: relative; }
.ai-score-inner { width: 48px; height: 48px; background: var(--surface); border-radius: 50%; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.ai-score-val { font-size: 14px; font-weight: 700; color: var(--success); line-height: 1; }
.ai-score-label { font-size: 9px; color: var(--text-secondary); margin-top: 2px; }

/* Mock CTA Styles */
.mock-cta { margin-top: 32px; padding-top: 24px; border-top: 1px solid var(--border); }
.mock-cta-title { font-weight: 600; font-size: 15px; margin-bottom: 8px; color: var(--text-main); }
.mock-cta-desc { font-size: 13px; color: var(--text-secondary); margin-bottom: 16px; line-height: 1.6; }
.mock-cta-btn { width: 100%; border-color: var(--accent); color: var(--accent); display: flex; align-items: center; justify-content: center; gap: 8px; transition: all 0.2s; font-weight: 600; }
.mock-cta-btn:hover { background: var(--accent-bg); }
</style>
