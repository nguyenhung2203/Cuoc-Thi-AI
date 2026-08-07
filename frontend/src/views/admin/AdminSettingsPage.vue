<script setup>
import { ref, reactive, onMounted } from 'vue'
import Card from '../../components/common/AppCard.vue'
import { apiService } from '../../services/api.service'
import { usePlatformStore } from '../../stores/platform.store'
import { Settings, Sliders, Bot, Shield, Save, RefreshCw, Plus, Edit3, CheckCircle2, AlertCircle, Trash2, Cpu, Key, FileText, Bell } from 'lucide-vue-next'
import { isEmail, isOneOf, isUrl, maxLength, normalizeText, requiredTrim } from '../../utils/validators.js'

const platformStore = usePlatformStore()

const activeTab = ref('general')
const loading = ref(false)
const saving = ref(false)
const successMessage = ref('')
const errorMessage = ref('')

// General & Security Settings State
const settings = reactive({
  system_name: platformStore.systemName || 'ViệcLàm AI Recruitment',
  brand_name: platformStore.brandName || 'ViệcLàm',
  brand_badge: platformStore.brandBadge || 'AI',
  brand_slogan: platformStore.brandSlogan || 'Nền tảng Phỏng vấn & Tuyển dụng Thông minh',
  brand_logo_url: platformStore.brandLogoUrl || '',
  support_email: platformStore.supportEmail || 'support@vieclam.ai',
  maintenance_mode: false,
  max_upload_size_mb: 20,
  default_passing_score: 70,
  default_ai_model: 'gemini-3.1-flash-lite',
  default_ai_voice_model: 'gemini-2.5-flash-native-audio-latest',
  jwt_token_expiry_hours: 24,
  admin_2fa_required: false,
  notification_ttl_days: 30,
  notification_max_per_user: 200,
  enable_email_notifications: true,
  notify_on_new_applicant: true,
  notify_on_report_ready: true,
  notify_on_interview_cancelled: true
})

// AI Prompt Templates State
const promptTemplates = ref([])
const loadingPrompts = ref(false)
const aiSettings = reactive({
  text_model: 'gemini-3.1-flash-lite',
  voice_model: 'gemini-2.5-flash-native-audio-latest',
  api_key: '',
  api_key_configured: false,
  api_key_masked: ''
})
const removingApiKey = ref(false)
const showPromptModal = ref(false)
const editingPrompt = reactive({
  name: '',
  model: 'gemini-3.1-flash-lite',
  content: ''
})

const fetchSettings = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await apiService.get('/admin/settings')
    const data = Array.isArray(res) ? res[0] : (res.data || res || {})
    if (data && typeof data === 'object') {
      if (data.system_name !== undefined) {
        settings.system_name = (data.system_name && data.system_name.includes('WeMake')) 
          ? 'ViệcLàm AI Platform' 
          : data.system_name
      }
      if (data.maintenance_mode !== undefined) settings.maintenance_mode = data.maintenance_mode
      if (data.max_upload_size_mb !== undefined) settings.max_upload_size_mb = data.max_upload_size_mb
      if (data.default_passing_score !== undefined) settings.default_passing_score = data.default_passing_score
      if (data.default_ai_model !== undefined) settings.default_ai_model = data.default_ai_model
      if (data.default_ai_voice_model !== undefined) settings.default_ai_voice_model = data.default_ai_voice_model
      if (data.jwt_token_expiry_hours !== undefined) settings.jwt_token_expiry_hours = data.jwt_token_expiry_hours
      if (data.admin_2fa_required !== undefined) settings.admin_2fa_required = data.admin_2fa_required
      if (data.notification_ttl_days !== undefined) settings.notification_ttl_days = data.notification_ttl_days
      if (data.notification_max_per_user !== undefined) settings.notification_max_per_user = data.notification_max_per_user
      if (data.enable_email_notifications !== undefined) settings.enable_email_notifications = data.enable_email_notifications
      if (data.notify_on_new_applicant !== undefined) settings.notify_on_new_applicant = data.notify_on_new_applicant
      if (data.notify_on_report_ready !== undefined) settings.notify_on_report_ready = data.notify_on_report_ready
      if (data.notify_on_interview_cancelled !== undefined) settings.notify_on_interview_cancelled = data.notify_on_interview_cancelled
      if (data.brand_name !== undefined) settings.brand_name = data.brand_name
      if (data.brand_badge !== undefined) settings.brand_badge = data.brand_badge
      if (data.brand_slogan !== undefined) settings.brand_slogan = data.brand_slogan
      if (data.brand_logo_url !== undefined) settings.brand_logo_url = data.brand_logo_url
      if (data.support_email !== undefined) settings.support_email = data.support_email
    }
  } catch (err) {
    console.error('Failed to load settings:', err)
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  if (saving.value) return
  settings.system_name = normalizeText(settings.system_name)
  settings.support_email = normalizeText(settings.support_email)
  settings.brand_logo_url = normalizeText(settings.brand_logo_url)
  let validationError = requiredTrim(settings.system_name, 'Vui lòng nhập tên hệ thống.')
    || maxLength(settings.system_name, 255, 'Tên hệ thống không được vượt quá 255 ký tự.')
    || requiredTrim(settings.support_email, 'Vui lòng nhập email hỗ trợ.')
    || isEmail(settings.support_email)
  if (!validationError && settings.brand_logo_url) validationError = isUrl(settings.brand_logo_url)
  const numericRules = [
    ['Dung lượng upload', settings.max_upload_size_mb, 1, 1024],
    ['Điểm đạt', settings.default_passing_score, 0, 100],
    ['Thời hạn JWT', settings.jwt_token_expiry_hours, 1, 720],
    ['Thời gian lưu thông báo', settings.notification_ttl_days, 1, 3650],
    ['Số thông báo tối đa', settings.notification_max_per_user, 1, 10000],
  ]
  for (const [label, value, min, max] of numericRules) {
    const number = Number(value)
    if (!validationError && (!Number.isFinite(number) || number < min || number > max)) validationError = `${label} phải nằm trong khoảng ${min} đến ${max}.`
  }
  if (validationError) {
    errorMessage.value = validationError
    return
  }
  saving.value = true
  successMessage.value = ''
  errorMessage.value = ''
  try {
    // Cập nhật cấu hình thương hiệu động trên toàn bộ ứng dụng
    platformStore.updateConfig(settings)
    await apiService.put('/admin/settings', { ...settings })
    successMessage.value = 'Cập nhật cấu hình hệ thống & thương hiệu thành công!'
    setTimeout(() => { successMessage.value = '' }, 4000)
  } catch (err) {
    console.error('Failed to save settings:', err)
    errorMessage.value = 'Lỗi khi lưu cấu hình: ' + (err.response?.data?.message || err.message)
  } finally {
    saving.value = false
  }
}

const fetchAISettings = async () => {
  try {
    const res = await apiService.get('/admin/ai-settings')
    const data = res.data || res || {}
    aiSettings.text_model = data.text_model || 'gemini-3.1-flash-lite'
    aiSettings.voice_model = data.voice_model || 'gemini-2.5-flash-native-audio-latest'
    aiSettings.api_key_configured = Boolean(data.api_key_configured)
    aiSettings.api_key_masked = data.api_key_masked || ''
    aiSettings.api_key = ''
  } catch (err) {
    console.error('Failed to load AI settings:', err)
  }
}

const saveAISettings = async () => {
  if (saving.value) return
  const validationError = isOneOf(aiSettings.text_model, ['gemini-3.1-flash-lite', 'gemini-2.5-flash-lite', 'gemini-2.5-flash'], 'Mô hình AI text không hợp lệ.')
    || isOneOf(aiSettings.voice_model, ['gemini-2.5-flash-native-audio-latest'], 'Mô hình AI voice không hợp lệ.')
  if (validationError) { errorMessage.value = validationError; return }
  saving.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const payload = { text_model: aiSettings.text_model, voice_model: aiSettings.voice_model }
    if (aiSettings.api_key.trim()) payload.api_key = aiSettings.api_key.trim()
    const res = await apiService.put('/admin/ai-settings', payload)
    const data = res.data || res || {}
    aiSettings.api_key = ''
    aiSettings.api_key_configured = Boolean(data.api_key_configured)
    aiSettings.api_key_masked = data.api_key_masked || ''
    settings.default_ai_model = aiSettings.text_model
    settings.default_ai_voice_model = aiSettings.voice_model
    successMessage.value = 'Đã lưu cấu hình Gemini an toàn.'
  } catch (err) {
    errorMessage.value = 'Lỗi khi lưu cấu hình AI: ' + (err.response?.data?.message || err.message)
  } finally { saving.value = false }
}

const removeApiKey = async () => {
  if (!window.confirm('Xóa Gemini API key đã lưu? Các chức năng AI sẽ ngừng hoạt động nếu không có key từ môi trường.')) return
  removingApiKey.value = true
  try {
    const res = await apiService.put('/admin/ai-settings', { text_model: aiSettings.text_model, voice_model: aiSettings.voice_model, clear_api_key: true })
    const data = res.data || res || {}
    aiSettings.api_key_configured = Boolean(data.api_key_configured)
    aiSettings.api_key_masked = data.api_key_masked || ''
    aiSettings.api_key = ''
    successMessage.value = 'Đã xóa Gemini API key.'
  } catch (err) {
    errorMessage.value = 'Không thể xóa API key: ' + (err.response?.data?.message || err.message)
  } finally { removingApiKey.value = false }
}

const fetchPromptTemplates = async () => {
  loadingPrompts.value = true
  try {
    const res = await apiService.get('/admin/ai-prompts')
    const data = res && typeof res === 'object' ? (res.data ?? res) : []
    promptTemplates.value = Array.isArray(data) ? data : []
  } catch (err) {
    promptTemplates.value = []
    console.error('Failed to load AI prompts:', err)
  } finally {
    loadingPrompts.value = false
  }
}

const openNewPromptModal = () => {
  editingPrompt.name = 'RUBRIC_EVALUATION_V2'
  editingPrompt.model = aiSettings.text_model || 'gemini-3.1-flash-lite'
  editingPrompt.content = 'Bạn là chuyên gia nhân sự AI. Hãy đánh giá câu trả lời sau dựa trên tiêu chí...'
  showPromptModal.value = true
}

const openEditPromptModal = (tmpl) => {
  editingPrompt.name = tmpl.name
  editingPrompt.model = tmpl.model || aiSettings.text_model || 'gemini-3.1-flash-lite'
  editingPrompt.content = tmpl.content || ''
  showPromptModal.value = true
}

const savePromptTemplate = async () => {
  if (saving.value) return
  editingPrompt.name = normalizeText(editingPrompt.name)
  editingPrompt.content = normalizeText(editingPrompt.content)
  const promptError = requiredTrim(editingPrompt.name, 'Vui lòng nhập tên prompt.')
    || maxLength(editingPrompt.name, 255, 'Tên prompt không được vượt quá 255 ký tự.')
    || requiredTrim(editingPrompt.content, 'Vui lòng nhập nội dung prompt.')
    || maxLength(editingPrompt.content, 50000, 'Nội dung prompt không được vượt quá 50.000 ký tự.')
    || isOneOf(editingPrompt.model, ['gemini-3.1-flash-lite', 'gemini-2.5-flash-lite', 'gemini-2.5-flash'], 'Mô hình AI không hợp lệ.')
  if (promptError) {
    errorMessage.value = promptError
    setTimeout(() => { errorMessage.value = '' }, 3000)
    return
  }
  saving.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await apiService.post('/admin/ai-prompts', {
      name: editingPrompt.name,
      model: editingPrompt.model,
      content: editingPrompt.content
    })
    showPromptModal.value = false
    await fetchPromptTemplates()
    successMessage.value = `Đã lưu phiên bản mới của Prompt "${editingPrompt.name}" thành công!`
    setTimeout(() => { successMessage.value = '' }, 4000)
  } catch (err) {
    errorMessage.value = 'Lỗi lưu Prompt: ' + (err.response?.data?.message || err.message)
    setTimeout(() => { errorMessage.value = '' }, 5000)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchSettings()
  fetchAISettings()
  fetchPromptTemplates()
})
</script>

<template>
  <div class="space-y-6 animate-fade-in pb-12">
    <!-- Page Header -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm border border-[var(--border)]">
      <div>
        <h1 class="text-h1 flex items-center gap-2.5">
          <Settings size="26" class="text-[var(--primary)]" />
          Cài đặt & Cấu hình Hệ thống
        </h1>
        <p class="text-[var(--text-secondary)] text-sm mt-1">
          Quản lý toàn bộ thông số kỹ thuật, mô hình AI, Prompt Templates và chính sách bảo mật hệ sinh thái.
        </p>
      </div>
      <button 
        @click="fetchSettings(); fetchPromptTemplates()" 
        :disabled="loading"
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-[var(--surface-soft)] hover:bg-[var(--border)] text-[var(--text-main)] rounded-xl font-semibold text-sm transition-colors shrink-0 disabled:opacity-50"
      >
        <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới cấu hình
      </button>
    </Card>

    <!-- Alert Notifications -->
    <div v-if="successMessage" class="p-4 bg-[var(--success)]/10 border border-[var(--success)]/20 rounded-2xl text-[var(--success)] font-semibold text-sm flex items-center gap-2">
      <CheckCircle2 size="18" class="shrink-0 text-[var(--success)]" />
      <span>{{ successMessage }}</span>
    </div>
    <div v-if="errorMessage" class="p-4 bg-[var(--danger)]/10 border border-[var(--danger)]/20 rounded-2xl text-[var(--danger)] font-semibold text-sm flex items-center gap-2">
      <AlertCircle size="18" class="shrink-0 text-[var(--danger)]" />
      <span>{{ errorMessage }}</span>
    </div>

    <!-- Tabs Navigation -->
    <div class="flex items-center gap-1.5 border-b border-[var(--border)] pb-2 flex-wrap">
      <button 
        @click="activeTab = 'general'"
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg font-semibold text-xs transition-all"
        :class="activeTab === 'general' ? 'bg-[var(--primary)] text-white shadow-xs' : 'text-[var(--text-secondary)] hover:bg-[var(--surface-soft)]'"
      >
        <Sliders size="14" />
        <span>Cấu hình chung</span>
      </button>

      <button 
        @click="activeTab = 'ai'"
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg font-semibold text-xs transition-all"
        :class="activeTab === 'ai' ? 'bg-[var(--primary)] text-white shadow-xs' : 'text-[var(--text-secondary)] hover:bg-[var(--surface-soft)]'"
      >
        <Bot size="14" />
        <span>AI & Prompts</span>
        <span class="px-1.5 py-0.5 rounded-full text-[10px] font-bold" :class="activeTab === 'ai' ? 'bg-white/20 text-white' : 'bg-[var(--primary-light)] text-[var(--primary)]'">
          {{ promptTemplates.length }}
        </span>
      </button>

      <button 
        @click="activeTab = 'security'"
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg font-semibold text-xs transition-all"
        :class="activeTab === 'security' ? 'bg-[var(--primary)] text-white shadow-xs' : 'text-[var(--text-secondary)] hover:bg-[var(--surface-soft)]'"
      >
        <Shield size="14" />
        <span>Bảo mật</span>
      </button>

      <button 
        @click="activeTab = 'notifications'"
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg font-semibold text-xs transition-all"
        :class="activeTab === 'notifications' ? 'bg-[var(--primary)] text-white shadow-xs' : 'text-[var(--text-secondary)] hover:bg-[var(--surface-soft)]'"
      >
        <Bell size="14" />
        <span>Thông báo</span>
      </button>
    </div>

    <!-- Tab 1: General Settings -->
    <Card v-if="activeTab === 'general'" class="rounded-2xl border border-[var(--border)] shadow-sm p-6 space-y-4">
      <div class="border-b border-[var(--border)] pb-3">
        <h3 class="text-base font-bold text-[var(--text-main)] flex items-center gap-2">
          <Sliders class="text-[var(--primary)]" size="18" />
          Thông số hoạt động nền tảng
        </h3>
        <p class="text-xs text-[var(--text-secondary)] mt-1">
          Các thiết lập sẽ áp dụng tức thì cho toàn bộ các dịch vụ tuyển dụng và phỏng vấn AI.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- System Name -->
        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Tên nền tảng hệ thống
          </label>
          <input 
            type="text" 
            v-model="settings.system_name"
            placeholder="Ví dụ: ViệcLàm AI Platform"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
          <span class="text-[10px] text-[var(--text-secondary)] mt-1 block">Dùng trong báo cáo, log server & email hệ thống</span>
        </div>

        <!-- Brand Name -->
        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Tên thương hiệu hiển thị
          </label>
          <input 
            type="text" 
            v-model="settings.brand_name"
            placeholder="Ví dụ: ViệcLàm, Talent"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
          <span class="text-[10px] text-[var(--primary)] mt-1 block font-semibold">⭐ Hiển thị ở Navbar, Footer, Sidebar</span>
        </div>

        <!-- Brand Badge -->
        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Badge Logo
          </label>
          <input 
            type="text" 
            v-model="settings.brand_badge"
            placeholder="Ví dụ: AI, PRO, VN"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
        </div>

        <!-- Support Email -->
        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Email hỗ trợ hệ thống
          </label>
          <input 
            type="email" 
            v-model="settings.support_email"
            placeholder="support@vieclam.ai"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
        </div>

        <!-- Brand Slogan -->
        <div class="md:col-span-2">
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Slogan & Tiêu đề phụ Website
          </label>
          <input 
            type="text" 
            v-model="settings.brand_slogan"
            placeholder="Nền tảng Phỏng vấn & Tuyển dụng Thông minh"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
        </div>

        <!-- Custom Logo URL -->
        <div class="md:col-span-2">
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1 flex items-center justify-between">
            <span>Custom Logo URL (PNG, JPG, SVG)</span>
            <span class="text-[10px] font-normal text-[var(--accent)]">Để trống dùng Icon AI mặc định</span>
          </label>
          <input 
            type="text" 
            v-model="settings.brand_logo_url"
            placeholder="/images/logo.png hoặc https://example.com/logo.svg"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
        </div>

        <!-- Max File Size & Default Score in a 2-col grid -->
        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Dung lượng file tối đa (MB)
          </label>
          <input 
            type="number" 
            v-model.number="settings.max_upload_size_mb"
            min="1" max="100"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
        </div>

        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Điểm đạt tối thiểu / 100
          </label>
          <input 
            type="number" 
            v-model.number="settings.default_passing_score"
            min="10" max="100"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
        </div>
      </div>

      <!-- Maintenance Mode Toggle -->
      <div class="pt-3 border-t border-[var(--border)] flex items-center justify-between">
        <div>
          <h4 class="font-bold text-[var(--text-main)] text-sm">Chế độ bảo trì (Maintenance Mode)</h4>
          <p class="text-xs text-[var(--text-secondary)] mt-0.5">
            Khi bật, toàn bộ ứng viên và nhà tuyển dụng sẽ nhận thông báo nâng cấp. Chỉ Admin mới có quyền truy cập.
          </p>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" v-model="settings.maintenance_mode" class="sr-only peer">
          <div class="w-14 h-7 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all dark:border-slate-600 peer-checked:bg-[var(--highlight)]"></div>
        </label>
      </div>

      <!-- Save Button -->
      <div class="pt-4 border-t border-[var(--border)] flex justify-end">
        <button 
          @click="saveSettings" 
          :disabled="saving"
          class="px-5 py-2 bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white font-bold text-sm rounded-xl shadow-xs transition-all disabled:opacity-50 flex items-center gap-2"
        >
          <Save size="16" />
          {{ saving ? 'Đang lưu...' : 'Lưu Cấu hình' }}
        </button>
      </div>
    </Card>

    <!-- Tab 2: AI & Prompt Templates -->
    <div v-if="activeTab === 'ai'" class="space-y-6">
      <!-- AI Model Configuration Card -->
      <Card class="p-8 rounded-2xl border border-[var(--border)] shadow-sm space-y-6">
        <div class="border-b border-[var(--border)] pb-4">
          <h3 class="text-lg font-bold text-[var(--text-main)] flex items-center gap-2">
            <Bot class="text-[var(--primary)]" size="20" />
            Cấu hình mô hình AI Phỏng vấn mặc định
          </h3>
          <p class="text-xs text-[var(--text-secondary)] mt-1">
            Thiết lập mô hình AI phù hợp cho hai chế độ: Phỏng vấn bằng chữ (Text) và Phỏng vấn qua giọng nói (Live Voice).
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="md:col-span-2 p-4 rounded-xl border border-[var(--border)] bg-[var(--surface-soft)]">
            <label class="block text-xs font-bold text-[var(--text-main)] uppercase tracking-wider mb-2 flex items-center gap-1.5">
              <Key size="14" class="text-[var(--primary)]" /> Gemini API key
            </label>
            <input v-model="aiSettings.api_key" type="password" autocomplete="new-password" placeholder="Để trống để giữ key hiện tại" class="w-full px-4 py-3 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)]" />
            <p class="text-[11px] text-[var(--text-secondary)] mt-2" v-if="aiSettings.api_key_configured">Đã cấu hình: {{ aiSettings.api_key_masked }}. Key không được hiển thị lại.</p>
            <p class="text-[11px] text-[var(--text-secondary)] mt-2" v-else>Key được lưu phía máy chủ và không gửi lại trình duyệt.</p>
            <button v-if="aiSettings.api_key_configured" @click="removeApiKey" :disabled="removingApiKey" type="button" class="mt-3 px-3 py-1.5 text-xs font-bold text-red-600 border border-red-200 rounded-lg hover:bg-red-50 disabled:opacity-50">{{ removingApiKey ? 'Đang xóa...' : 'Xóa API key' }}</button>
          </div>
          <!-- Text AI Model -->
          <div>
            <label class="block text-xs font-bold text-[var(--text-main)] uppercase tracking-wider mb-2 flex items-center gap-1.5">
              <FileText size="14" class="text-[var(--primary)]" /> Mô hình AI Phỏng vấn dạng Text
            </label>
            <select 
              v-model="aiSettings.text_model"
              class="w-full px-4 py-3 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm font-semibold text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
            >
              <option value="gemini-3.1-flash-lite">Gemini 3.1 Flash Lite (Khuyên dùng - quota 500 lượt/ngày)</option>
              <option value="gemini-2.5-flash-lite">Gemini 2.5 Flash Lite (Dự phòng tiết kiệm)</option>
              <option value="gemini-2.5-flash">Gemini 2.5 Flash (Chất lượng cao hơn, quota thấp)</option>
            </select>
            <span class="text-[11px] text-[var(--text-secondary)] mt-1.5 block">Sử dụng cho phỏng vấn trực tiếp bằng chữ & chat.</span>
          </div>

          <!-- Voice/Live AI Model -->
          <div>
            <label class="block text-xs font-bold text-[var(--text-main)] uppercase tracking-wider mb-2 flex items-center gap-1.5">
              <Bot size="14" class="text-[var(--primary)]" /> Mô hình AI Phỏng vấn dạng Voice / Live
            </label>
            <select 
              v-model="aiSettings.voice_model"
              class="w-full px-4 py-3 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm font-semibold text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
            >
              <option value="gemini-2.5-flash-native-audio-latest">Gemini 2.5 Flash Native Audio (Khuyên dùng - Live Voice)</option>
            </select>
            <span class="text-[11px] text-[var(--text-secondary)] mt-1.5 block">Sử dụng cho phòng phỏng vấn thử giọng nói real-time (Gemini Live WebSocket).</span>
          </div>
        </div>

        <div class="flex justify-end pt-2">
          <button 
            @click="saveAISettings"
            :disabled="saving"
            class="px-5 py-2.5 bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white font-bold text-sm rounded-xl shadow-xs transition-all disabled:opacity-50 flex items-center gap-2"
          >
            <Save size="16" /> Lưu cấu hình Mô hình AI
          </button>
        </div>
      </Card>

      <!-- Action Bar -->
      <Card class="p-6 rounded-2xl border border-[var(--border)] shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h3 class="text-lg font-bold text-[var(--text-main)] flex items-center gap-2">
            <Cpu class="text-[var(--primary)] animate-pulse" size="20" />
            Thư viện Prompt Templates AI
          </h3>
          <p class="text-xs text-[var(--text-secondary)] mt-1">
            Mỗi lần cập nhật nội dung prompt sẽ tự động tạo một phiên bản mới (Version history) trong cơ sở dữ liệu.
          </p>
        </div>
        <button 
          @click="openNewPromptModal"
          class="px-5 py-2.5 bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white font-bold text-sm rounded-xl shadow-xs transition-all flex items-center gap-2 shrink-0"
        >
          <Plus size="18" /> Thêm Prompt Template mới
        </button>
      </Card>

      <!-- Prompts Grid -->
      <div v-if="loadingPrompts" class="py-12 text-center text-[var(--text-secondary)]">
        <RefreshCw size="24" class="animate-spin mx-auto mb-2 text-[var(--primary)]" />
        Đang tải danh sách Prompt Templates...
      </div>

      <Card v-else-if="promptTemplates.length === 0" class="p-12 rounded-2xl text-center border border-[var(--border)]">
        <Bot size="36" class="mx-auto mb-3 text-[var(--text-secondary)]" />
        <h4 class="font-bold text-[var(--text-main)]">Chưa có Prompt Template nào được tùy biến</h4>
        <p class="text-xs text-[var(--text-secondary)] mt-1 mb-4">Hệ thống đang sử dụng Prompt AI mặc định được nạp từ mã nguồn.</p>
        <button @click="openNewPromptModal" class="px-4 py-2 bg-[var(--primary)] text-white font-bold text-xs rounded-xl">Tạo Prompt mẫu ngay</button>
      </Card>

      <div v-else class="grid grid-cols-1 gap-4">
        <div 
          v-for="tmpl in promptTemplates" 
          :key="tmpl.id"
          class="bg-[var(--surface)] p-6 rounded-2xl border border-[var(--border)] shadow-sm hover:border-[var(--primary)]/30 transition-all flex flex-col md:flex-row md:items-center justify-between gap-4"
        >
          <div class="space-y-2 flex-1">
            <div class="flex items-center gap-3">
              <span class="px-3 py-1 rounded-lg bg-[var(--primary-light)] text-[var(--primary)] font-extrabold font-mono text-xs">
                v{{ tmpl.version || 1 }}
              </span>
              <h4 class="text-base font-extrabold text-[var(--text-main)]">{{ tmpl.name }}</h4>
              <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-[var(--success)]/10 text-[var(--success)]">
                Active
              </span>
              <span class="text-xs text-[var(--text-secondary)] font-mono">Model: {{ tmpl.model || 'gemini-3.1-flash-lite' }}</span>
            </div>
            <p class="text-xs text-[var(--text-secondary)] font-mono line-clamp-2 bg-[var(--background)] p-3 rounded-xl border border-[var(--border)]">
              {{ tmpl.content }}
            </p>
          </div>
          
          <div class="flex items-center gap-2 shrink-0">
            <button 
              @click="openEditPromptModal(tmpl)"
              class="px-4 py-2 bg-[var(--surface-soft)] hover:bg-[var(--border)] text-[var(--text-main)] font-bold text-xs rounded-xl transition-colors flex items-center gap-1.5"
            >
              <Edit3 size="14" /> Chỉnh sửa / Tạo bản v{{ (tmpl.version || 1) + 1 }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tab 3: Security & Sessions -->
    <Card v-if="activeTab === 'security'" class="rounded-2xl border border-[var(--border)] shadow-sm p-6 space-y-4">
      <div class="border-b border-[var(--border)] pb-3">
        <h3 class="text-base font-bold text-[var(--text-main)] flex items-center gap-2">
          <Shield class="text-[var(--primary)]" size="18" />
          Bảo mật phiên làm việc (Sessions &amp; Auth)
        </h3>
        <p class="text-xs text-[var(--text-secondary)] mt-1">
          Thiết lập tiêu chuẩn an toàn truy cập cho các tài khoản Quản trị và Nhà tuyển dụng.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- JWT Expiry -->
        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1 flex items-center gap-1.5">
            <Key size="12" class="text-[var(--primary)]" /> Thời gian hiệu lực JWT Token (Giờ)
          </label>
          <input 
            type="number" 
            v-model.number="settings.jwt_token_expiry_hours"
            min="1" max="168"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
          <span class="text-[10px] text-[var(--text-secondary)] mt-1 block">Mặc định: 24 giờ. Token hết hạn sẽ tự động làm mới bằng Refresh Cookie.</span>
        </div>

        <!-- Admin 2FA Toggle -->
        <div class="bg-[var(--surface)] p-4 rounded-xl border border-[var(--border)] flex items-center justify-between">
          <div>
            <h4 class="font-bold text-[var(--text-main)] text-sm flex items-center gap-1.5">
              <Shield size="14" class="text-[var(--success)]" /> Xác thực 2 bước (2FA Admin)
            </h4>
            <p class="text-xs text-[var(--text-secondary)] mt-0.5">
              Yêu cầu OTP email khi truy cập các tính năng nâng cao.
            </p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" v-model="settings.admin_2fa_required" class="sr-only peer">
            <div class="w-14 h-7 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all dark:border-slate-600 peer-checked:bg-[var(--primary)]"></div>
          </label>
        </div>
      </div>

      <div class="pt-3 border-t border-[var(--border)] flex justify-end">
        <button 
          @click="saveSettings" 
          :disabled="saving"
          class="px-4 py-1.5 bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white font-bold text-sm rounded-lg shadow-xs transition-all disabled:opacity-50 flex items-center gap-2"
        >
          <Save size="14" />
          {{ saving ? 'Đang lưu...' : 'Lưu thay đổi' }}
        </button>
      </div>
    </Card>

    <!-- Tab 4: Notification Settings -->
    <Card v-if="activeTab === 'notifications'" class="rounded-2xl border border-[var(--border)] shadow-sm p-6 space-y-4 animate-fade-in">
      <div class="border-b border-[var(--border)] pb-3">
        <h3 class="text-base font-bold text-[var(--text-main)] flex items-center gap-2">
          <Bell class="text-[var(--primary)]" size="18" />
          Cấu hình Thông báo &amp; Redis TTL
        </h3>
        <p class="text-xs text-[var(--text-secondary)] mt-1">
          Quản lý thời gian lưu trữ Redis, giới hạn số lượng và kiểm soát bật/tắt các loại sự kiện thông báo.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- Redis TTL -->
        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Thời gian lưu trữ Redis / TTL (Ngày)
          </label>
          <input 
            type="number" 
            v-model.number="settings.notification_ttl_days"
            min="1" max="365"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
          <span class="text-[10px] text-[var(--text-secondary)] mt-1 block">Mặc định: 30 ngày. Thông báo cũ hơn TTL sẽ tự động xóa.</span>
        </div>

        <!-- Max items per user -->
        <div>
          <label class="block text-[11px] font-bold text-[var(--text-secondary)] uppercase tracking-wider mb-1">
            Giới hạn tối đa mỗi User (Items)
          </label>
          <input 
            type="number" 
            v-model.number="settings.notification_max_per_user"
            min="20" max="1000"
            class="w-full px-3 py-2 bg-[var(--surface)] border border-[var(--border)] rounded-lg text-sm text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
          <span class="text-[10px] text-[var(--text-secondary)] mt-1 block">Mặc định: 200 thông báo gần nhất.</span>
        </div>
      </div>

      <div class="border-t border-[var(--border)] pt-4">
        <h4 class="font-bold text-[var(--text-main)] text-sm mb-3">Bật / Tắt theo sự kiện</h4>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <!-- Notify on new applicant -->
          <div class="bg-[var(--surface)] px-4 py-3 rounded-xl border border-[var(--border)] flex items-center justify-between">
            <div>
              <h5 class="font-semibold text-[var(--text-main)] text-sm">Ứng viên mới nộp CV</h5>
              <p class="text-[11px] text-[var(--text-secondary)] mt-0.5">Khi ứng viên mới apply vào tin tuyển dụng.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer shrink-0">
              <input type="checkbox" v-model="settings.notify_on_new_applicant" class="sr-only peer">
              <div class="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-slate-600 peer-checked:bg-[var(--primary)]"></div>
            </label>
          </div>

          <!-- Notify on report ready -->
          <div class="bg-[var(--surface)] px-4 py-3 rounded-xl border border-[var(--border)] flex items-center justify-between">
            <div>
              <h5 class="font-semibold text-[var(--text-main)] text-sm">Báo cáo AI sẵn sàng</h5>
              <p class="text-[11px] text-[var(--text-secondary)] mt-0.5">Khi AI hoàn tất chấm điểm và xuất báo cáo.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer shrink-0">
              <input type="checkbox" v-model="settings.notify_on_report_ready" class="sr-only peer">
              <div class="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-slate-600 peer-checked:bg-[var(--primary)]"></div>
            </label>
          </div>

          <!-- Notify on interview cancelled -->
          <div class="bg-[var(--surface)] px-4 py-3 rounded-xl border border-[var(--border)] flex items-center justify-between">
            <div>
              <h5 class="font-semibold text-[var(--text-main)] text-sm">Lịch phỏng vấn bị hủy</h5>
              <p class="text-[11px] text-[var(--text-secondary)] mt-0.5">Khi ứng viên hoặc HR hủy lịch phỏng vấn.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer shrink-0">
              <input type="checkbox" v-model="settings.notify_on_interview_cancelled" class="sr-only peer">
              <div class="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-slate-600 peer-checked:bg-[var(--primary)]"></div>
            </label>
          </div>

          <!-- Enable email notifications -->
          <div class="bg-[var(--surface)] px-4 py-3 rounded-xl border border-[var(--border)] flex items-center justify-between">
            <div>
              <h5 class="font-semibold text-[var(--text-main)] text-sm">Gửi song song qua Email</h5>
              <p class="text-[11px] text-[var(--text-secondary)] mt-0.5">Gửi email kèm thông báo cho các sự kiện quan trọng.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer shrink-0">
              <input type="checkbox" v-model="settings.enable_email_notifications" class="sr-only peer">
              <div class="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-slate-600 peer-checked:bg-[var(--primary)]"></div>
            </label>
          </div>
        </div>
      </div>

      <div class="pt-3 border-t border-[var(--border)] flex justify-end">
        <button 
          @click="saveSettings" 
          :disabled="saving"
          class="px-5 py-2 bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white font-bold text-sm rounded-xl shadow-xs transition-all disabled:opacity-50 flex items-center gap-2"
        >
          <Save size="16" />
          {{ saving ? 'Đang lưu...' : 'Lưu Cấu hình Thông báo' }}
        </button>
      </div>
    </Card>

    <!-- Modal for Create/Edit Prompt -->

    <div v-if="showPromptModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs animate-fade-in">
      <div class="bg-[var(--surface)] rounded-3xl max-w-2xl w-full p-8 shadow-2xl border border-[var(--border)] space-y-6 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-[var(--border)] pb-4">
          <h3 class="text-xl font-bold text-[var(--text-main)] flex items-center gap-2">
            <Cpu class="text-[var(--primary)]" size="24" />
            Cấu hình Prompt Template AI
          </h3>
          <button @click="showPromptModal = false" class="text-[var(--text-secondary)] hover:text-[var(--text-main)] text-sm font-bold">✕ Đóng</button>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-[var(--text-main)] uppercase tracking-wider mb-2">Tên định danh Prompt Template</label>
            <input 
              type="text" 
              v-model="editingPrompt.name" 
              placeholder="e.g. RUBRIC_EVALUATION, INTERVIEW_CONDUCTOR"
              class="w-full px-4 py-3 font-mono text-sm bg-[var(--background)] border border-[var(--border)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)]"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[var(--text-main)] uppercase tracking-wider mb-2">Mô hình AI xử lý</label>
            <select 
              v-model="editingPrompt.model"
              class="w-full px-4 py-3 bg-[var(--background)] border border-[var(--border)] rounded-xl text-sm font-semibold text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)]"
            >
              <option value="gemini-3.1-flash-lite" class="text-[var(--primary)]">Gemini 3.1 Flash Lite</option>
              <option value="gemini-2.5-flash-lite" class="text-[var(--primary)]">Gemini 2.5 Flash Lite</option>
              <option value="gemini-2.5-flash" class="text-[var(--primary)]">Gemini 2.5 Flash</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-bold text-[var(--text-main)] uppercase tracking-wider mb-2">
              Nội dung System Prompt (Có thể sử dụng biến {{ '{' + '{ variable_name }' + '}' }})
            </label>
            <textarea 
              v-model="editingPrompt.content" 
              rows="8"
              placeholder="Nhập hướng dẫn chi tiết cho mô hình AI..."
              class="w-full p-4 font-mono text-xs bg-[var(--background)] border border-[var(--border)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] leading-relaxed"
            ></textarea>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-4 border-t border-[var(--border)]">
          <button @click="showPromptModal = false" class="px-5 py-2.5 bg-[var(--surface-soft)] hover:bg-[var(--border)] text-[var(--text-main)] font-bold text-sm rounded-xl transition-colors">
            Hủy bỏ
          </button>
          <button @click="savePromptTemplate" :disabled="saving" class="px-6 py-2.5 bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white font-bold text-sm rounded-xl shadow-xs transition-all disabled:opacity-50 flex items-center gap-2">
            <Save size="16" /> {{ saving ? 'Đang lưu phiên bản mới...' : 'Lưu & Khởi tạo Version mới' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
