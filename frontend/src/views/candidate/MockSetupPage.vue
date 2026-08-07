<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { 
  Play, Bot, Info, CheckCircle2, 
  Code, Server, Layers, Cpu, Cloud, Briefcase, Smartphone, ShieldCheck, 
  PlusCircle, Clock, Globe, Zap, Award, Target, ChevronRight,
  Upload, Loader2, Calculator, TrendingUp, Landmark, Handshake, ShoppingBag,
  Megaphone, FileText, Headset, Settings, UserPlus, Building2, Scale, Truck,
  ShoppingCart, Package, HardHat, DraftingCompass, Home, GraduationCap,
  HeartPulse, Palette, Hotel, Plane, Factory, Newspaper, Leaf, ClipboardCheck, Users,
  Search,
} from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { langStore } from '../../stores/lang.store'
import { isOneOf, maxLength, minLength, normalizeText, requiredTrim, validateForm } from '../../utils/validators.js'
import { getRoleCatalog } from '../../utils/mockRoleSuggestions.js'

const router = useRouter()
const loading = ref(false)
const profileLoading = ref(true)
const toast = ref(null)
const cvFileId = ref('')
const cvName = ref('')
const roleQuery = ref('')
const errors = ref({})

const setup = ref({
  jobRole: '',
  customRoleName: '',
  customJd: '',
  showCustomInput: false,
  level: 'fresher',
  type: 'behavior',
  language: 'vi',
  duration: 'standard',
  style: 'professional',
  useCurrentCv: true,
})

const iconMap = {
  Code, Server, Layers, Cpu, Cloud, Briefcase, Smartphone, ShieldCheck,
  Calculator, TrendingUp, Landmark, Handshake, ShoppingBag, Megaphone,
  FileText, Award, Headset, Settings, UserPlus, Building2, Scale, Truck,
  ShoppingCart, Package, HardHat, DraftingCompass, Home, GraduationCap,
  HeartPulse, Palette, Hotel, Plane, Factory, Newspaper, Leaf, ClipboardCheck,
  Users, Globe,
}

const catalogRoles = computed(() => getRoleCatalog().map((role) => ({
  ...role,
  icon: iconMap[role.icon] || Briefcase,
})))

const filteredRoles = computed(() => {
  const q = roleQuery.value.trim().toLowerCase()
  if (!q) return catalogRoles.value
  return catalogRoles.value.filter((role) =>
    role.name.toLowerCase().includes(q) || (role.desc || '').toLowerCase().includes(q)
  )
})

const hasCv = computed(() => !!cvFileId.value)
const canUseStudio = computed(() => !profileLoading.value)

const levels = [
  { id: 'fresher', name: 'Fresher / Intern', badge: 'Entry Level', desc: 'Mới bắt đầu / thực tập — kiến thức nền tảng' },
  { id: 'junior', name: 'Junior', badge: '1 - 2 Năm', desc: 'Làm được việc cơ bản, cần hướng dẫn' },
  { id: 'middle', name: 'Middle', badge: '3 - 5 Năm', desc: 'Tự chủ công việc, xử lý tình huống thực tế' },
  { id: 'senior', name: 'Senior', badge: '5+ Năm', desc: 'Chuyên sâu nghiệp vụ, mentoring đồng nghiệp' },
  { id: 'lead', name: 'Lead / Manager', badge: 'Lead', desc: 'Điều phối nhóm, định hướng công việc' }
]

const interviewTypes = [
  { id: 'tech', name: 'Phỏng vấn Chuyên môn', desc: 'Kiến thức nghiệp vụ, kỹ năng cứng theo vị trí ứng tuyển' },
  { id: 'behavior', name: 'Phỏng vấn Hành vi (STAR)', desc: 'Tình huống thực tế, làm việc nhóm, xử lý áp lực' },
  { id: 'hr', name: 'Phỏng vấn Nhân sự (HR)', desc: 'Định hướng nghề nghiệp, văn hóa, lương thưởng' },
  { id: 'system', name: 'Case / Tình huống nghiệp vụ', desc: 'Phân tích case, đưa phương án giải quyết' },
  { id: 'algo', name: 'Tư duy & Giải quyết vấn đề', desc: 'Logic, ưu tiên, lập kế hoạch (phù hợp nhiều ngành)' }
]

const languages = [
  { id: 'vi', name: 'Tiếng Việt', desc: 'Phỏng vấn chuẩn mực bằng Tiếng Việt' },
  { id: 'en', name: 'Tiếng Anh', desc: 'Luyện phản xạ Tiếng Anh chuyên ngành' }
]

const durations = [
  { id: 'quick', name: 'Phỏng vấn Nhanh', badge: '3 câu hỏi · ~15 phút', desc: 'Kiểm tra phản xạ & ôn luyện nhanh' },
  { id: 'standard', name: 'Phỏng vấn Tiêu chuẩn', badge: '5 câu hỏi · ~30 phút', desc: 'Mô phỏng phỏng vấn thực tế' },
  { id: 'deep', name: 'Phỏng vấn Chuyên sâu', badge: '8 câu hỏi · ~45 phút', desc: 'Đào sâu kiến trúc & thử thách áp lực' }
]

const styles = [
  { id: 'friendly', name: 'Thân thiện & Gợi mở', desc: 'Tạo tâm lý thoải mái, gợi ý nhẹ nhàng' },
  { id: 'professional', name: 'Chuyên nghiệp & Chuẩn mực', desc: 'Bám sát tiêu chuẩn đánh giá' },
  { id: 'challenging', name: 'Khó tính & Hỏi xoáy', desc: 'Xoáy sâu vào các vấn đề hóc búa' }
]

onMounted(async () => {
  profileLoading.value = true
  try {
    const res = await candidatePortalService.getProfile()
    const profile = res?.data || res || {}
    cvFileId.value = profile.cv_file_id || ''
    cvName.value = profile.cv_name || ''
    setup.value.useCurrentCv = !!cvFileId.value
  } catch (err) {
    cvFileId.value = ''
    setup.value.useCurrentCv = false
  } finally {
    profileLoading.value = false
  }
})

const goUploadCv = () => {
  router.push({ path: '/profile', query: { tab: 'profile' } })
}

const selectRole = (roleId) => {
  setup.value.jobRole = roleId
  setup.value.showCustomInput = false
  setup.value.customRoleName = ''
}

const selectCustomRole = () => {
  setup.value.showCustomInput = true
  setup.value.jobRole = 'custom'
}

const onCustomRoleInput = () => {
  setup.value.showCustomInput = true
  setup.value.jobRole = 'custom'
}

const getSelectedRoleName = () => {
  if (setup.value.showCustomInput && setup.value.customRoleName.trim()) {
    return setup.value.customRoleName.trim()
  }
  const r = catalogRoles.value.find(item => item.id === setup.value.jobRole)
  return r ? r.name : ''
}

const getSelectedLevelName = () => {
  const l = levels.find(item => item.id === setup.value.level)
  return l ? l.name : 'Middle'
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

  loading.value = true
  try {
    const finalRole = getSelectedRoleName()
    const payload = {
      target_role: finalRole,
      target_level: setup.value.level,
    }
    if (cvFileId.value && setup.value.useCurrentCv) {
      payload.cv_file_id = cvFileId.value
    }
    const session = await mockService.createMockInterview(payload)

    await mockService.startMockInterview(session.id)

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
    toast.value = {
      type: 'error',
      message: error?.message || 'Không thể khởi động phỏng vấn. Vui lòng kiểm tra kết nối máy chủ hoặc thử lại sau.'
    }
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
        <Button variant="primary" class="btn-nowrap" :disabled="loading || !canUseStudio" @click="handleStart">
          <Play :size="16" class="shrink-0" />
          <span>{{ loading ? langStore.t('setup', 'starting') : langStore.t('setup', 'startNow') }}</span>
        </Button>
      </div>
    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <!-- Loading -->
    <Card v-if="profileLoading" class="cv-gate-card animate-rise">
      <div class="cv-gate-inner">
        <Loader2 class="pf-spin text-primary" :size="28" />
        <p class="cv-gate-desc">Đang tải hồ sơ...</p>
      </div>
    </Card>

    <!-- Main 2-Column Grid -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
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
              <p class="section-subtitle">
                Chọn vị trí bạn muốn luyện — đa ngành (Kinh tế, Kinh doanh, Marketing, IT…).
              </p>
            </div>
          </div>

          <div class="role-search">
            <Search :size="16" class="role-search-icon" />
            <input
              v-model="roleQuery"
              type="search"
              class="role-search-input"
              placeholder="Tìm vị trí (ví dụ: kế toán, marketing, sales…)"
            />
          </div>

          <div class="role-grid-wrap">
            <div class="role-grid">
              <button 
                v-for="role in filteredRoles" 
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
            </div>
            <p v-if="!filteredRoles.length" class="role-empty">Không tìm thấy vị trí phù hợp. Hãy nhập vị trí tùy chỉnh bên dưới.</p>
          </div>

          <!-- Fixed custom job form (always below 3-row grid, like search) -->
          <div class="custom-job-bar" :class="{ active: setup.showCustomInput }">
            <div class="custom-job-bar-icon">
              <PlusCircle :size="18" />
            </div>
            <input
              v-model="setup.customRoleName"
              type="text"
              class="custom-job-input"
              placeholder="Nhập vị trí / chức danh tùy chỉnh (ví dụ: Chuyên viên Kinh tế…"
              @focus="selectCustomRole"
              @input="onCustomRoleInput"
            />
          </div>

          <!-- Custom Role & JD Expandable Box -->
          <div v-if="setup.showCustomInput" class="custom-jd-box animate-fadeIn">
            <div class="field">
              <div class="flex items-center justify-between">
                <label class="field-label font-bold text-main">Mô tả công việc / Yêu cầu (JD) — tùy chọn</label>
                <span class="text-xs text-primary font-semibold">⚡ AI xoáy sâu skill trong JD</span>
              </div>
              <textarea 
                v-model="setup.customJd" 
                rows="4" 
                placeholder="Dán JD từ TopCV, LinkedIn hoặc ghi nhanh yêu cầu công việc..." 
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

          <!-- CV optional nudge -->
          <div class="cv-integrate-bar" :class="{ active: hasCv && setup.useCurrentCv }">
            <label v-if="hasCv" class="cv-toggle-label">
              <input
                v-model="setup.useCurrentCv"
                type="checkbox"
                class="cv-checkbox"
              />
              <div>
                <div class="font-bold text-main text-sm">
                  Dùng CV trong hồ sơ
                  <span class="text-xs text-[var(--success)] font-semibold ml-1">Khuyến nghị</span>
                </div>
                <p class="cv-desc">
                  {{ cvName || 'CV đã tải lên' }} — AI hỏi sát kinh nghiệm/kỹ năng trong CV hơn.
                </p>
              </div>
            </label>
            <div v-else class="cv-toggle-label">
              <Upload :size="20" class="text-[var(--primary)] shrink-0" />
              <div>
                <div class="font-bold text-main text-sm">
                  Chưa có CV — vẫn luyện được
                  <span class="text-xs text-secondary font-semibold ml-1">Không bắt buộc</span>
                </div>
                <p class="cv-desc">
                  Câu hỏi sẽ theo vị trí bạn chọn. Tải CV trên Hồ sơ để AI cá nhân hóa tốt hơn.
                  <button type="button" class="cv-upload-link" @click="goUploadCv">Tải CV ngay</button>
                </p>
              </div>
            </div>
          </div>

          <!-- Start Action Button -->
          <div class="mt-8 flex flex-col sm:flex-row items-center justify-between gap-4 pt-6 border-t border-[var(--border)]">
            <div class="flex items-center gap-2 text-xs text-secondary font-medium">
              <ShieldCheck :size="16" class="text-[var(--success)] shrink-0" />
              <span>{{ langStore.t('setup', 'realtimeScoring') }}</span>
            </div>

            <Button 
              variant="primary" 
              class="w-full sm:w-auto btn-start-studio"
              :disabled="loading || !canUseStudio"
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
.role-search {
  position: relative;
  margin-bottom: 12px;
}
.role-search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary, #64748B);
  pointer-events: none;
}
.role-search-input {
  width: 100%;
  border: 1px solid var(--border, #E2E8F0);
  border-radius: 10px;
  padding: 10px 12px 10px 36px;
  font-size: 14px;
  background: var(--surface, #fff);
  color: var(--text-main, #0f172a);
}
.role-search-input:focus {
  outline: none;
  border-color: var(--primary, #2563EB);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
}
.role-grid-wrap {
  /* Viewport cố định ~3 hàng; danh sách cuộn bên trong */
  max-height: calc(3 * 86px + 2 * 12px + 20px);
  overflow-y: auto;
  border: 1px solid var(--border, #E2E8F0);
  border-radius: 12px;
  padding: 10px;
  background: var(--surface-soft, #F8FAFC);
}
@media (max-width: 639px) {
  .role-grid-wrap {
    max-height: calc(3 * 86px + 2 * 12px + 20px);
  }
}
.role-empty {
  margin: 12px 4px 4px;
  font-size: 13px;
  color: var(--text-secondary, #64748B);
}
.role-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 12px;
}
@media (min-width: 640px) {
  .role-grid { grid-template-columns: repeat(2, 1fr); }
}
.custom-job-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 12px;
  border: 1px dashed #CBD5E1;
  border-radius: 10px;
  padding: 10px 12px;
  background: var(--surface, #fff);
  transition: border-color 0.15s, box-shadow 0.15s;
}
.custom-job-bar.active {
  border-style: solid;
  border-color: var(--accent, #10B981);
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.12);
}
.custom-job-bar-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent, #10B981);
  background: rgba(16, 185, 129, 0.08);
  flex-shrink: 0;
}
.custom-job-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: 14px;
  color: var(--text-main, #0f172a);
  min-width: 0;
}
.custom-job-input::placeholder {
  color: var(--text-secondary, #94A3B8);
}
.role-card {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 14px 16px;
  min-height: 86px;
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
}
.cv-upload-link {
  display: inline;
  margin-left: 4px;
  border: none;
  background: none;
  padding: 0;
  color: var(--primary, #2563EB);
  font-weight: 600;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 2px;
}
.cv-upload-link:hover {
  color: var(--primary-hover, #1D4ED8);
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

.cv-gate-card {
  padding: 40px 24px;
  text-align: center;
}
.cv-gate-inner {
  max-width: 480px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}
.cv-gate-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-light, #dbeafe);
  color: var(--primary, #2563eb);
}
.cv-gate-title {
  font-size: 20px;
  font-weight: 750;
  color: var(--text-main);
  margin: 0;
}
.cv-gate-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.55;
  margin: 0 0 8px;
}
.skill-hint-bar {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 16px;
  padding: 10px 12px;
  border-radius: var(--radius, 10px);
  background: rgba(37, 99, 235, 0.08);
  color: var(--primary, #2563eb);
  font-size: 12px;
  font-weight: 600;
  line-height: 1.45;
}
.skill-hint-bar.is-warn {
  background: rgba(245, 158, 11, 0.12);
  color: #b45309;
}
.role-badge {
  display: inline-block;
  margin-left: 6px;
  padding: 1px 7px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 700;
  vertical-align: middle;
  background: rgba(16, 185, 129, 0.15);
  color: #047857;
}
.role-badge.is-related {
  background: rgba(100, 116, 139, 0.14);
  color: #475569;
}
.pf-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
