<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { 
  Play, FileText, Settings2, Bot, Info, CheckCircle2, Sparkles, 
  Code, Server, Layers, Cpu, Cloud, Briefcase, Smartphone, ShieldCheck, 
  PlusCircle, Clock, Globe, Zap, Award, Target, MessageSquare, AlertCircle, ChevronRight
} from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { langStore } from '../../stores/lang.store'
import { isOneOf, maxLength, minLength, normalizeText, requiredTrim, validateForm } from '../../utils/validators.js'

const router = useRouter()
const loading = ref(false)
const toast = ref(null)
const cvFileId = ref('')
const errors = ref({})

// Interactive setup options
const setup = ref({
  jobRole: 'frontend',
  customRoleName: '',
  customJd: '',
  showCustomInput: false,
  level: 'middle',
  type: 'tech',
  language: 'vi',
  duration: 'standard',
  style: 'professional',
  useCurrentCv: true
})

// Quick-select preset job roles with rich icons
const presetRoles = [
  { id: 'frontend', name: 'Frontend Developer', desc: 'React, Vue, Architecture', icon: Code },
  { id: 'backend', name: 'Backend Developer', desc: 'Go, Java, Node, API', icon: Server },
  { id: 'fullstack', name: 'Fullstack Engineer', desc: 'End-to-End System', icon: Layers },
  { id: 'ai_ml', name: 'AI / ML Engineer', desc: 'LLM, RAG, PyTorch', icon: Cpu },
  { id: 'devops', name: 'DevOps & Cloud', desc: 'AWS, K8s, CI/CD', icon: Cloud },
  { id: 'pm', name: 'Product Manager', desc: 'Product Strategy', icon: Briefcase },
  { id: 'mobile', name: 'Mobile Developer', desc: 'iOS, Android, React Native', icon: Smartphone },
  { id: 'qa', name: 'QA & Automation', desc: 'Test Automation', icon: ShieldCheck }
]

// Seniority levels
const levels = [
  { id: 'fresher', name: 'Fresher / Intern', badge: 'Entry Level', desc: 'Kiến thức nền tảng CS, tư duy thuật toán & cơ bản' },
  { id: 'junior', name: 'Junior Engineer', badge: '1 - 2 Năm', desc: 'Thực hành coding, bug fixing & làm việc nhóm' },
  { id: 'middle', name: 'Middle Engineer', badge: '3 - 5 Năm', desc: 'Độc lập giải quyết vấn đề, trade-off & clean code' },
  { id: 'senior', name: 'Senior Engineer', badge: '5+ Năm', desc: 'Kiến trúc hệ thống, scale, tối ưu & leadership' },
  { id: 'lead', name: 'Tech Lead / Principal', badge: 'Expert / Lead', desc: 'Định hướng công nghệ, quản lý đội ngũ & high-concurrency' }
]

// Interview focus types
const interviewTypes = [
  { id: 'tech', name: 'Phỏng vấn Kỹ thuật (Technical & Core)', desc: 'Chuyên sâu ngôn ngữ, framework, best practices & kiến thức nền tảng' },
  { id: 'system', name: 'Thiết kế Hệ thống (System Design)', desc: 'Vẽ kiến trúc, microservices, cacher, load balancing & scalability' },
  { id: 'algo', name: 'Thuật toán & Problem Solving', desc: 'Cấu trúc dữ liệu, tối ưu độ phức tạp Big-O & tư duy giải thuật' },
  { id: 'behavior', name: 'Phỏng vấn Hành vi (STAR Behavioral)', desc: 'Khử xung đột, làm việc nhóm, quản lý áp lực & xử lý tình huống thực tế' },
  { id: 'hr', name: 'Phỏng vấn Nhân sự (HR & Culture Fit)', desc: 'Định hướng sự nghiệp, lương thưởng, độ phù hợp văn hóa doanh nghiệp' }
]

// Languages
const languages = [
  { id: 'vi', name: 'Tiếng Việt', desc: 'Phỏng vấn chuẩn mực bằng Tiếng Việt' },
  { id: 'en', name: 'Tiếng Anh', desc: 'Luyện phản xạ Tiếng Anh chuyên ngành IT' }
]

// Durations
const durations = [
  { id: 'quick', name: 'Phỏng vấn Nhanh', badge: '3 câu hỏi · ~15 phút', desc: 'Kiểm tra phản xạ & ôn luyện nhanh' },
  { id: 'standard', name: 'Phỏng vấn Tiêu chuẩn', badge: '5 câu hỏi · ~30 phút', desc: 'Mô phỏng phỏng vấn thực tế' },
  { id: 'deep', name: 'Phỏng vấn Chuyên sâu', badge: '8 câu hỏi · ~45 phút', desc: 'Đào sâu kiến trúc & thử thách áp lực' }
]

// AI Personas
const styles = [
  { id: 'friendly', name: 'Thân thiện & Gợi mở', desc: 'Tạo tâm lý thoải mái, gợi ý nhẹ nhàng' },
  { id: 'professional', name: 'Chuyên nghiệp & Chuẩn mực', desc: 'Bám sát tiêu chuẩn đánh giá' },
  { id: 'challenging', name: 'Khó tính & Hỏi xoáy', desc: 'Xoáy sâu vào các vấn đề hóc búa' }
]

onMounted(async () => {
  try {
    const res = await candidatePortalService.getProfile()
    const profile = res?.data || res || {}
    cvFileId.value = profile.cv_file_id || ''
    if (!cvFileId.value) {
      setup.value.useCurrentCv = false
    }
  } catch (err) {
    cvFileId.value = ''
    setup.value.useCurrentCv = false
  }
})

const selectRole = (roleId) => {
  setup.value.jobRole = roleId
  setup.value.showCustomInput = false
}

const selectCustomRole = () => {
  setup.value.showCustomInput = true
  setup.value.jobRole = 'custom'
}

const getSelectedRoleName = () => {
  if (setup.value.showCustomInput && setup.value.customRoleName.trim()) {
    return setup.value.customRoleName.trim()
  }
  const r = presetRoles.find(item => item.id === setup.value.jobRole)
  return r ? r.name : 'Senior Software Engineer'
}

const getSelectedLevelName = () => {
  const l = levels.find(item => item.id === setup.value.level)
  return l ? l.name : 'Middle Engineer'
}

const getSelectedDurationBadge = () => {
  const d = durations.find(item => item.id === setup.value.duration)
  return d ? d.badge : '5 câu hỏi · ~30 phút'
}

const handleStart = async (e) => {
  e.preventDefault()
  if (loading.value) return
  const values = { ...setup.value, customRoleName: normalizeText(setup.value.customRoleName), customJd: normalizeText(setup.value.customJd) }
  const validation = validateForm(values, {
    jobRole: [(value) => requiredTrim(value, 'Vui lòng chọn vị trí phỏng vấn.')],
    level: [(value) => isOneOf(value, levels.map(item => item.id), 'Cấp độ không hợp lệ.')],
    type: [(value) => isOneOf(value, interviewTypes.map(item => item.id), 'Loại phỏng vấn không hợp lệ.')],
    language: [(value) => isOneOf(value, languages.map(item => item.id), 'Ngôn ngữ không hợp lệ.')],
    duration: [(value) => isOneOf(value, durations.map(item => item.id), 'Thời lượng không hợp lệ.')],
    style: [(value) => isOneOf(value, styles.map(item => item.id), 'Phong cách AI không hợp lệ.')],
    customRoleName: [
      (value) => setup.value.showCustomInput ? requiredTrim(value, 'Vui lòng nhập tên vị trí.') : '',
      (value) => setup.value.showCustomInput ? minLength(value, 2, 'Tên vị trí phải có ít nhất 2 ký tự.') : '',
      (value) => maxLength(value, 100, 'Tên vị trí không được vượt quá 100 ký tự.'),
    ],
    customJd: [(value) => maxLength(value, 20000, 'JD không được vượt quá 20.000 ký tự.')],
  })
  errors.value = validation.errors
  if (!validation.isValid) {
    toast.value = { type: 'error', message: Object.values(validation.errors)[0] }
    return
  }
  if (setup.value.useCurrentCv && !cvFileId.value) {
    toast.value = { type: 'warning', message: 'Bạn chưa có CV trong hồ sơ. Vui lòng tải CV lên trước hoặc bỏ chọn tùy chọn sử dụng CV.' }
    return
  }
  if (setup.value.showCustomInput && !setup.value.customRoleName.trim()) {
    toast.value = { type: 'warning', message: 'Vui lòng nhập tên vị trí hoặc chức danh bạn muốn phỏng vấn.' }
    return
  }

  loading.value = true
  try {
    const finalRole = getSelectedRoleName()
    
    // Bước 1: Tạo session mới trong Backend Go
    const session = await mockService.createMockInterview({
      target_role: finalRole,
      target_level: setup.value.level,
      cv_file_id: setup.value.useCurrentCv ? cvFileId.value : undefined
    })

    // Bước 2: Khởi động session
    await mockService.startMockInterview(session.id)

    // Bước 3: Chuyển vào phòng phỏng vấn mock
    router.push({
      path: '/mock-room',
      query: {
        mock_id: session.id,
        role: finalRole,
        level: setup.value.level,
        type: setup.value.type,
        style: setup.value.style,
        language: setup.value.language,
        duration: setup.value.duration,
        jd: setup.value.showCustomInput ? setup.value.customJd : ''
      },
    })
  } catch (error) {
    console.error(error)
    toast.value = { type: 'error', message: 'Không thể khởi động phỏng vấn. Vui lòng kiểm tra kết nối máy chủ.' }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="setup-container max-w-6xl mx-auto pb-16 space-y-8">
    <!-- Top Header Banner (Exact style as PracticeHistoryPage) -->
    <div class="header-box animate-rise">
      <div class="header-info">
        <h1 class="page-title">{{ langStore.t('setup', 'title') }}</h1>
        <p class="page-subtitle">
          {{ langStore.t('setup', 'subtitle') }}
        </p>
      </div>

      <div class="header-cta">
        <Button variant="primary" class="btn-nowrap" :disabled="loading" @click="handleStart">
          <Play :size="16" class="shrink-0" />
          <span>{{ loading ? langStore.t('setup', 'starting') : langStore.t('setup', 'startNow') }}</span>
        </Button>
      </div>
    </div>

    <!-- Main 2-Column Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
      <!-- Left / Main Form Column (2 Cols) -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Section 1: Job Role Selection -->
        <Card class="studio-section animate-rise">
          <div class="section-top">
            <div class="section-icon-box primary-bg">
              <Target :size="18" style="color: var(--primary)" />
            </div>
            <div>
              <h2 class="section-title">{{ langStore.t('setup', 'step1Title') }}</h2>
              <p class="section-subtitle">{{ langStore.t('setup', 'step1Desc') }}</p>
            </div>
          </div>

          <div class="role-grid">
            <button 
              v-for="role in presetRoles" 
              :key="role.id"
              type="button"
              @click="selectRole(role.id)"
              class="role-card"
              :class="{ active: !setup.showCustomInput && setup.jobRole === role.id }"
            >
              <div class="role-icon">
                <component :is="role.icon" :size="20" />
              </div>
              <div class="role-content">
                <div class="role-name">{{ role.name }}</div>
                <div class="role-desc">{{ role.desc }}</div>
              </div>
              <div class="radio-indicator"></div>
            </button>

            <!-- Custom Role / JD Button -->
            <button 
              type="button"
              @click="selectCustomRole"
              class="role-card custom-role-card"
              :class="{ active: setup.showCustomInput }"
            >
              <div class="role-icon custom-icon">
                <PlusCircle :size="20" />
              </div>
              <div class="role-content">
                <div class="role-name">Nhập Vị trí khác / Dán JD tùy chỉnh</div>
                <div class="role-desc">Tùy chỉnh 100% câu hỏi theo mô tả công việc cụ thể bạn đang ứng tuyển</div>
              </div>
              <div class="radio-indicator"></div>
            </button>
          </div>

          <!-- Custom Role & JD Expandable Box -->
          <div v-if="setup.showCustomInput" class="custom-jd-box animate-fadeIn">
            <div class="field">
              <label class="field-label font-bold text-main">Tên Chức danh / Vị trí cụ thể <span class="text-rose-500">*</span></label>
              <input 
                v-model="setup.customRoleName" 
                type="text" 
                placeholder="Ví dụ: Senior Data Engineer, Cloud Solution Architect, Golang Backend Developer..." 
                class="app-input"
              />
            </div>

            <div class="field mt-4">
              <div class="flex items-center justify-between">
                <label class="field-label font-bold text-main">Mô tả công việc / Yêu cầu kỹ thuật (Job Description - JD)</label>
                <span class="text-xs text-primary font-semibold">⚡ AI sẽ xoáy sâu vào các skill trong JD này</span>
              </div>
              <textarea 
                v-model="setup.customJd" 
                rows="4" 
                placeholder="Dán nội dung JD từ TopCV, LinkedIn hoặc ghi chú nhanh các công nghệ bắt buộc (ví dụ: yêu cầu thành thạo Kafka, Redis clustering, kinh nghiệm làm việc với hệ thống 1 triệu CCU...)" 
                class="app-textarea"
              ></textarea>
            </div>
          </div>
        </Card>

        <!-- Section 2: Seniority Level -->
        <Card class="studio-section animate-rise">
          <div class="section-top">
            <div class="section-icon-box accent-bg">
              <Award :size="18" style="color: var(--accent)" />
            </div>
            <div>
              <h2 class="section-title">{{ langStore.t('setup', 'step2Title') }}</h2>
              <p class="section-subtitle">{{ langStore.t('setup', 'step2Desc') }}</p>
            </div>
          </div>

          <div class="level-grid">
            <button 
              v-for="lvl in levels" 
              :key="lvl.id"
              type="button"
              @click="setup.level = lvl.id"
              class="level-card"
              :class="{ active: setup.level === lvl.id }"
            >
              <div class="level-header">
                <span class="level-name">{{ lvl.name }}</span>
                <span class="level-pill">{{ lvl.badge }}</span>
              </div>
              <p class="level-desc">{{ lvl.desc }}</p>
            </button>
          </div>
        </Card>

        <!-- Section 3: Interview Type & Focus -->
        <Card class="studio-section animate-rise">
          <div class="section-top">
            <div class="section-icon-box info-bg">
              <Layers :size="18" style="color: var(--primary)" />
            </div>
            <div>
              <h2 class="section-title">{{ langStore.t('setup', 'step3Title') }}</h2>
              <p class="section-subtitle">{{ langStore.t('setup', 'step3Desc') }}</p>
            </div>
          </div>

          <div class="type-grid">
            <button 
              v-for="t in interviewTypes" 
              :key="t.id"
              type="button"
              @click="setup.type = t.id"
              class="type-card"
              :class="{ active: setup.type === t.id }"
            >
              <div class="type-head">
                <span class="type-name">{{ t.name }}</span>
              </div>
              <p class="type-desc">{{ t.desc }}</p>
            </button>
          </div>
        </Card>

        <!-- Section 4: Duration & Language -->
        <Card class="studio-section animate-rise">
          <div class="section-top">
            <div class="section-icon-box warning-bg">
              <Clock :size="18" style="color: var(--warning)" />
            </div>
            <div>
              <h2 class="section-title">{{ langStore.t('setup', 'step4Title') }}</h2>
              <p class="section-subtitle">{{ langStore.t('setup', 'step4Desc') }}</p>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Duration -->
            <div class="space-y-3">
              <label class="font-bold text-sm text-main flex items-center gap-2">
                <Clock :size="15" class="text-primary" />
                <span>Thời lượng & Số câu hỏi:</span>
              </label>
              <div class="space-y-2.5">
                <button 
                  v-for="dur in durations" 
                  :key="dur.id"
                  type="button"
                  @click="setup.duration = dur.id"
                  class="option-pill"
                  :class="{ active: setup.duration === dur.id }"
                >
                  <div class="pill-head">
                    <span class="font-bold text-main">{{ dur.name }}</span>
                    <span class="pill-badge">{{ dur.badge }}</span>
                  </div>
                  <p class="pill-desc">{{ dur.desc }}</p>
                </button>
              </div>
            </div>

            <!-- Language -->
            <div class="space-y-3">
              <label class="font-bold text-sm text-main flex items-center gap-2">
                <Globe :size="15" class="text-primary" />
                <span>Ngôn ngữ thoại & câu hỏi AI:</span>
              </label>
              <div class="space-y-2.5">
                <button 
                  v-for="lang in languages" 
                  :key="lang.id"
                  type="button"
                  @click="setup.language = lang.id"
                  class="option-pill"
                  :class="{ active: setup.language === lang.id }"
                >
                  <div class="font-bold text-main">{{ lang.name }}</div>
                  <p class="pill-desc">{{ lang.desc }}</p>
                </button>
              </div>
            </div>
          </div>
        </Card>

        <!-- Section 5: AI Persona & CV Integration -->
        <Card class="studio-section animate-rise">
          <div class="section-top">
            <div class="section-icon-box primary-bg">
              <Bot :size="18" style="color: var(--primary)" />
            </div>
            <div>
              <h2 class="section-title">{{ langStore.t('setup', 'step5Title') }}</h2>
              <p class="section-subtitle">{{ langStore.t('setup', 'step5Desc') }}</p>
            </div>
          </div>

          <!-- AI Persona Grid -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
            <button 
              v-for="st in styles" 
              :key="st.id"
              type="button"
              @click="setup.style = st.id"
              class="persona-card"
              :class="{ active: setup.style === st.id }"
            >
              <div class="font-bold text-main mb-1">{{ st.name }}</div>
              <p class="persona-desc">{{ st.desc }}</p>
            </button>
          </div>

          <!-- CV Toggle Bar -->
          <div class="cv-integrate-bar" :class="{ active: setup.useCurrentCv }">
            <label class="cv-toggle-label">
              <input type="checkbox" v-model="setup.useCurrentCv" :disabled="!cvFileId" class="cv-checkbox" />
              <FileText :size="20" class="text-primary shrink-0" />
              <div>
                <div class="font-bold text-main text-sm">
                  Đồng bộ & hỏi dựa theo CV thực tế trong hồ sơ ứng viên
                  <span v-if="!cvFileId" class="text-xs text-rose-500 font-normal ml-1">(Hồ sơ chưa tải CV lên)</span>
                  <span v-else class="text-xs text-[var(--success)] font-semibold ml-1">✓ Đã sẵn sàng</span>
                </div>
                <p class="cv-desc">AI sẽ đọc kinh nghiệm, dự án (Projects) và công nghệ ghi trong CV của bạn để đặt câu hỏi xác thực</p>
              </div>
            </label>
          </div>

          <!-- Start Action Button -->
          <!-- Start Action Button -->
          <div class="mt-8 flex flex-col sm:flex-row items-center justify-between gap-4 pt-6 border-t border-[var(--border)]">
            <div class="flex items-center gap-2 text-xs text-secondary font-medium">
              <ShieldCheck :size="16" class="text-[var(--success)] shrink-0" />
              <span>{{ langStore.t('setup', 'realtimeScoring') }}</span>
            </div>

            <Button 
              variant="primary" 
              class="w-full sm:w-auto btn-start-studio"
              :disabled="loading"
              @click="handleStart"
            >
              <Play :size="18" class="shrink-0" />
              <span>{{ loading ? langStore.t('setup', 'starting') : langStore.t('setup', 'btnStartStudio') }}</span>
              <ChevronRight :size="18" class="shrink-0 ml-1 opacity-80" />
            </Button>
          </div>
        </Card>
      </div>

      <!-- Right Column: Live AI Studio Preview & Checklist -->
      <div class="lg:col-span-1 space-y-6 sticky top-6">
        <!-- Live Configuration Tower Card -->
        <div class="ai-studio-tile animate-rise">
          <div class="studio-status-tag">
            <span class="dot-live"><span class="dot-live-ping"></span></span>
            <span>STUDIO ONLINE</span>
          </div>

          <div class="studio-avatar">
            <Bot :size="48" />
            <span class="avatar-ring"></span>
          </div>

          <h3 class="studio-bot-name">AI Interviewer Pro</h3>
          <p class="studio-bot-role">Được phát triển trên nền tảng Gemini Real-time API</p>

          <div class="summary-box">
            <div class="summary-row">
              <span class="summary-label">{{ langStore.t('setup', 'selectedRole') }}:</span>
              <span class="summary-val">{{ getSelectedRoleName() }}</span>
            </div>
            <div class="summary-row">
              <span class="summary-label">{{ langStore.t('setup', 'selectedLevel') }}:</span>
              <span class="summary-val">{{ getSelectedLevelName() }}</span>
            </div>
            <div class="summary-row">
              <span class="summary-label">{{ langStore.t('setup', 'selectedTime') }}:</span>
              <span class="summary-val">{{ getSelectedDurationBadge() }}</span>
            </div>
            <div class="summary-row border-0 pt-1">
              <span class="summary-label">{{ langStore.t('setup', 'selectedPersona') }}:</span>
              <span class="summary-val capitalize">{{ setup.style }}</span>
            </div>
          </div>
        </div>

        <!-- Rubric Enterprise & Preparation Checklist (Collapsible) -->
        <div class="space-y-4">
          <details class="studio-accordion bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm animate-rise" open>
            <summary class="flex items-center gap-3 p-4 cursor-pointer font-bold text-main text-sm hover:bg-slate-50">
              <div class="kpi-icon-box warning-bg shrink-0"><Zap :size="18" style="color: var(--warning)" /></div>
              Chuẩn Đánh giá Rubric Enterprise
            </summary>
            <div class="p-4 pt-0 border-t border-slate-100 bg-slate-50/50">
              <p class="text-xs text-secondary leading-relaxed mb-3 mt-3">AI tự động chấm điểm đa chiều theo 6 trục:</p>
              <div class="flex flex-wrap gap-2">
                <span class="badge badge-info">Chuyên môn</span>
                <span class="badge badge-success">Cấu trúc STAR</span>
                <span class="badge badge-primary">Trade-offs</span>
                <span class="badge badge-warning">Giao tiếp</span>
                <span class="badge badge-danger">Áp lực & Phản xạ</span>
                <span class="badge badge-neutral">Thuật ngữ IT</span>
              </div>
            </div>
          </details>

          <details class="studio-accordion bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm animate-rise">
            <summary class="flex items-center gap-3 p-4 cursor-pointer font-bold text-main text-sm hover:bg-slate-50">
              <div class="kpi-icon-box info-bg shrink-0"><Info :size="18" style="color: var(--primary)" /></div>
              Hướng dẫn Thực chiến
            </summary>
            <div class="p-4 pt-0 border-t border-slate-100 bg-slate-50/50">
              <ul class="studio-tips mt-3 space-y-2 text-sm text-secondary">
                <li class="flex items-start gap-2"><CheckCircle2 :size="16" class="text-[var(--success)] shrink-0 mt-0.5" /><span>Micro rõ ràng, không gian tĩnh.</span></li>
                <li class="flex items-start gap-2"><CheckCircle2 :size="16" class="text-[var(--success)] shrink-0 mt-0.5" /><span>Nhấn giữ phím hoặc Gõ chữ đều được.</span></li>
                <li class="flex items-start gap-2"><CheckCircle2 :size="16" class="text-[var(--success)] shrink-0 mt-0.5" /><span>Nhấn Kết thúc sớm để nhận Report.</span></li>
              </ul>
            </div>
          </details>
        </div>
      </div>
    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
  </div>
</template>

<style scoped>
/* Header Banner — Exact match with PracticeHistoryPage (`header-box`) */
.header-box {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  padding: 24px 28px;
  border-radius: var(--radius-lg, 16px);
  background: var(--surface, #FFFFFF);
  border: 1px solid var(--border, #E2E8F0);
  box-shadow: var(--shadow-sm);
}
@media (min-width: 640px) {
  .header-box {
    flex-direction: row;
    align-items: center;
  }
}
.header-info {
  flex: 1;
}
.brand-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: var(--radius-full, 9999px);
  background: var(--primary-light, #DBEAFE);
  color: var(--primary, #2563EB);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  margin-bottom: 12px;
}
.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
  margin: 0 0 8px 0;
}
.page-subtitle {
  font-size: 14px;
  color: var(--text-secondary, #475569);
  max-width: 700px;
  line-height: 1.6;
  margin: 0;
}
.header-cta {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
.btn-nowrap {
  white-space: nowrap;
}

/* Section Cards */
.studio-section {
  padding: 28px;
  border-radius: 16px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface, #FFFFFF);
  box-shadow: 0 4px 20px rgba(0,0,0,0.03);
}
.section-top {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 24px;
}
.section-icon-box {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.primary-bg { background: rgba(37,99,235,0.1); }
.accent-bg { background: rgba(16,185,129,0.1); }
.info-bg { background: rgba(59,130,246,0.1); }
.warning-bg { background: rgba(245,158,11,0.1); }

.section-title {
  font-size: 17px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
  margin: 0 0 4px 0;
}
.section-subtitle {
  font-size: 13px;
  color: var(--text-secondary, #64748B);
  margin: 0;
}

/* Role Grid */
.role-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 12px;
}
@media (min-width: 640px) {
  .role-grid { grid-template-columns: repeat(2, 1fr); }
}
.role-card {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 16px;
  border-radius: 12px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface-soft, #F8FAFC);
  text-align: left;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  position: relative;
}
.role-card:hover {
  border-color: var(--primary, #2563EB);
  background: var(--surface, #FFFFFF);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(37,99,235,0.06);
}
.role-card.active {
  border-color: var(--primary, #2563EB);
  background: rgba(37,99,235,0.06);
  box-shadow: 0 0 0 2px var(--primary, #2563EB);
}
.role-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: var(--surface, #FFFFFF);
  border: 1px solid var(--border, #E2E8F0);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary, #2563EB);
  flex-shrink: 0;
}
.custom-icon {
  color: var(--accent, #10B981);
}
.role-content {
  flex: 1;
}
.role-name {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
  margin-bottom: 2px;
}
.role-desc {
  font-size: 12px;
  color: var(--text-secondary, #64748B);
  line-height: 1.4;
}
.custom-role-card {
  grid-column: 1 / -1;
  border-style: dashed;
  border-color: #CBD5E1;
}
.custom-role-card.active {
  border-style: solid;
}

/* Custom JD Box */
.custom-jd-box {
  margin-top: 20px;
  padding: 20px;
  border-radius: 12px;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
}
.app-input {
  width: 100%;
  height: 44px;
  padding: 0 14px;
  border-radius: 8px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface, #FFFFFF);
  font-size: 14px;
  color: var(--text-main, #0F172A);
  outline: none;
  margin-top: 6px;
  transition: all 0.2s ease;
}
.app-input:focus {
  border-color: var(--primary, #2563EB);
  box-shadow: 0 0 0 3px rgba(37,99,235,0.1);
}
.app-textarea {
  width: 100%;
  padding: 12px 14px;
  border-radius: 8px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface, #FFFFFF);
  font-size: 13px;
  color: var(--text-main, #0F172A);
  outline: none;
  margin-top: 6px;
  transition: all 0.2s ease;
  resize: vertical;
}
.app-textarea:focus {
  border-color: var(--primary, #2563EB);
  box-shadow: 0 0 0 3px rgba(37,99,235,0.1);
}

/* Level Grid */
.level-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 12px;
}
@media (min-width: 640px) {
  .level-grid { grid-template-columns: repeat(2, 1fr); }
}
.level-card {
  padding: 16px;
  border-radius: 12px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface-soft, #F8FAFC);
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
}
.level-card:hover {
  border-color: var(--primary, #2563EB);
  background: var(--surface, #FFFFFF);
}
.level-card.active {
  border-color: var(--primary, #2563EB);
  background: rgba(37,99,235,0.06);
  box-shadow: 0 0 0 2px var(--primary, #2563EB);
}
.level-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}
.level-name {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
}
.level-pill {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 999px;
  background: rgba(37,99,235,0.1);
  color: var(--primary, #2563EB);
}
.level-desc {
  font-size: 12px;
  color: var(--text-secondary, #64748B);
  margin: 0;
  line-height: 1.45;
}

/* Type Grid */
.type-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 12px;
}
@media (min-width: 640px) {
  .type-grid { grid-template-columns: repeat(2, 1fr); }
}
.type-card {
  padding: 16px;
  border-radius: 12px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface-soft, #F8FAFC);
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
}
.type-card:hover {
  border-color: var(--primary, #2563EB);
  background: var(--surface, #FFFFFF);
}
.type-card.active {
  border-color: var(--primary, #2563EB);
  background: rgba(37,99,235,0.06);
  box-shadow: 0 0 0 2px var(--primary, #2563EB);
}
.type-name {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
}
.type-desc {
  font-size: 12px;
  color: var(--text-secondary, #64748B);
  margin: 4px 0 0 0;
  line-height: 1.45;
  display: none;
}
.type-card.active .type-desc, .type-card:hover .type-desc {
  display: block;
}

/* Option Pill */
.option-pill {
  width: 100%;
  padding: 14px;
  border-radius: 10px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface-soft, #F8FAFC);
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
}
.option-pill:hover {
  border-color: var(--primary, #2563EB);
  background: var(--surface, #FFFFFF);
}
.option-pill.active {
  border-color: var(--primary, #2563EB);
  background: rgba(37,99,235,0.06);
  box-shadow: 0 0 0 2px var(--primary, #2563EB);
}
.pill-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.pill-badge {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 999px;
  background: rgba(245,158,11,0.12);
  color: #D97706;
}
.pill-desc {
  font-size: 12px;
  color: var(--text-secondary, #64748B);
  margin-top: 4px;
  line-height: 1.45;
}

/* Persona Card */
.persona-card {
  padding: 14px;
  border-radius: 10px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface-soft, #F8FAFC);
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
}
.persona-card:hover {
  border-color: var(--primary, #2563EB);
  background: var(--surface, #FFFFFF);
}
.persona-card.active {
  border-color: var(--primary, #2563EB);
  background: rgba(37,99,235,0.06);
  box-shadow: 0 0 0 2px var(--primary, #2563EB);
}
.persona-desc {
  font-size: 12px;
  color: var(--text-secondary, #64748B);
  margin-top: 4px;
  line-height: 1.45;
}

/* CV Toggle Bar */
.cv-integrate-bar {
  padding: 16px;
  border-radius: 12px;
  background: rgba(37,99,235,0.05);
  border: 1px solid rgba(37,99,235,0.2);
  transition: all 0.2s ease;
}
.cv-integrate-bar:hover {
  background: var(--surface, #FFFFFF);
  border-color: var(--primary, #2563EB);
}
.cv-integrate-bar.active {
  background: rgba(37,99,235,0.08);
  border-color: var(--primary, #2563EB);
}
.cv-toggle-label {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  cursor: pointer;
}
.cv-checkbox {
  width: 18px;
  height: 18px;
  accent-color: var(--primary, #2563EB);
  margin-top: 2px;
  cursor: pointer;
}
.cv-desc {
  font-size: 12px;
  color: var(--text-secondary, #64748B);
  margin-top: 4px;
  line-height: 1.45;
  display: none;
}
.cv-integrate-bar.active .cv-desc, .cv-integrate-bar:hover .cv-desc {
  display: block;
}

/* Start Button */
.btn-start-studio {
  height: 48px;
  padding: 0 32px;
  font-size: 14px;
  font-weight: 600;
  border-radius: var(--radius, 12px);
  box-shadow: var(--shadow-md);
}

/* AI Studio Tile Right Column — Cohesive with Navy Brand Identity & Design System */
.ai-studio-tile {
  position: relative;
  overflow: hidden;
  padding: 28px 24px;
  border-radius: var(--radius-lg, 16px);
  background: var(--surface, #FFFFFF);
  color: var(--text-main, #0F172A);
  text-align: center;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border, #E2E8F0);
  border-top: 4px solid var(--primary, #1E3A8A);
}
.studio-status-tag {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border-radius: var(--radius-full, 9999px);
  background: var(--primary-light, #DBEAFE);
  border: 1px solid rgba(30, 58, 138, 0.15);
  color: var(--primary, #1E3A8A);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.05em;
  margin-bottom: 20px;
  box-shadow: var(--shadow-sm);
}
.dot-live {
  position: relative;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--success, #10B981);
}
.dot-live-ping {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  background: var(--success, #10B981);
  animation: aiPulse 1.8s ease-out infinite;
}
@keyframes aiPulse {
  0% { transform: scale(1); opacity: 0.8; }
  100% { transform: scale(1.6); opacity: 0; }
}
.studio-avatar {
  position: relative;
  width: 80px;
  height: 80px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--primary-light, #DBEAFE);
  color: var(--primary, #1E3A8A);
  border: 2px solid rgba(30, 58, 138, 0.2);
  box-shadow: var(--shadow-sm);
}
.avatar-ring {
  position: absolute;
  inset: -6px;
  border-radius: 50%;
  border: 2px solid rgba(30, 58, 138, 0.2);
  animation: aiPulse 2.4s ease-out infinite;
}
.studio-bot-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
  margin: 0 0 4px 0;
}
.studio-bot-role {
  font-size: 13px;
  font-weight: 600;
  color: var(--primary, #1E3A8A);
  margin: 0 0 20px 0;
}
.summary-box {
  background: var(--surface-soft, #F8FAFC);
  border-radius: var(--radius, 12px);
  padding: 14px 16px;
  border: 1px solid var(--border, #E2E8F0);
  text-align: left;
  box-shadow: var(--shadow-sm);
  margin-bottom: 8px;
}
.summary-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px dashed rgba(255,255,255,0.2);
  font-size: 13px;
}
.summary-row:last-child {
  border-bottom: none;
  padding-bottom: 2px;
}
.summary-label {
  color: var(--text-secondary, #475569);
  font-weight: 500;
}
.summary-val {
  color: var(--text-main, #0F172A);
  font-weight: 600;
  text-align: right;
}

/* Studio Cards */
.studio-card {
  padding: 20px;
  border-radius: var(--radius-lg, 16px);
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface, #FFFFFF);
  box-shadow: var(--shadow-sm);
}
.rubric-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin: 0;
  padding: 0;
  list-style: none;
}
.rubric-list li {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-main, #0F172A);
}
.rubric-list li span {
  color: var(--primary, #2563EB);
  font-weight: 700;
}
.studio-tips {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 0;
  padding: 0;
  list-style: none;
}
.studio-tips li {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-size: 13px;
  color: var(--text-secondary, #475569);
  line-height: 1.5;
}
</style>
