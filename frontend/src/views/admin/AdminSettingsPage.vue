<script setup>
import { ref, reactive, onMounted } from 'vue'
import Card from '../../components/common/AppCard.vue'
import { apiService } from '../../services/api.service'
import { usePlatformStore } from '../../stores/platform.store'
import { Settings, Sliders, Bot, Shield, Save, RefreshCw, Plus, Edit3, CheckCircle2, AlertCircle, Trash2, Cpu, Key, FileText, Bell } from 'lucide-vue-next'

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
  default_ai_model: 'gemini-2.5-flash',
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
const showPromptModal = ref(false)
const editingPrompt = reactive({
  name: '',
  model: 'gemini-2.5-flash',
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

const fetchPromptTemplates = async () => {
  loadingPrompts.value = true
  try {
    const res = await apiService.get('/admin/ai-prompts')
    promptTemplates.value = Array.isArray(res) ? res : (res.data || [])
  } catch (err) {
    console.error('Failed to load AI prompts:', err)
  } finally {
    loadingPrompts.value = false
  }
}

const openNewPromptModal = () => {
  editingPrompt.name = 'RUBRIC_EVALUATION_V2'
  editingPrompt.model = settings.default_ai_model || 'gemini-2.5-flash'
  editingPrompt.content = 'Bạn là chuyên gia nhân sự AI. Hãy đánh giá câu trả lời sau dựa trên tiêu chí...'
  showPromptModal.value = true
}

const openEditPromptModal = (tmpl) => {
  editingPrompt.name = tmpl.name
  editingPrompt.model = tmpl.model || 'gemini-2.5-flash'
  editingPrompt.content = tmpl.content || ''
  showPromptModal.value = true
}

const savePromptTemplate = async () => {
  if (!editingPrompt.name || !editingPrompt.content) {
    alert('Vui lòng nhập tên và nội dung prompt')
    return
  }
  saving.value = true
  try {
    await apiService.post('/admin/ai-prompts', {
      name: editingPrompt.name,
      model: editingPrompt.model,
      content: editingPrompt.content
    })
    showPromptModal.value = false
    await fetchPromptTemplates()
    alert('Đã lưu phiên bản mới của Prompt Template thành công!')
  } catch (err) {
    alert('Lỗi lưu Prompt: ' + (err.response?.data?.message || err.message))
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchSettings()
  fetchPromptTemplates()
})
</script>

<template>
  <div class="space-y-6 animate-fade-in pb-12">
    <!-- Page Header -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-white flex items-center gap-2.5">
          <Settings size="26" class="text-blue-600 dark:text-blue-400" />
          Cài đặt & Cấu hình Hệ thống
        </h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
          Quản lý toàn bộ thông số kỹ thuật, mô hình AI, Prompt Templates và chính sách bảo mật hệ sinh thái.
        </p>
      </div>
      <button 
        @click="fetchSettings(); fetchPromptTemplates()" 
        :disabled="loading"
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-200 rounded-xl font-medium text-sm transition-colors shrink-0 disabled:opacity-50"
      >
        <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới cấu hình
      </button>
    </Card>

    <!-- Alert Notifications -->
    <div v-if="successMessage" class="p-4 bg-emerald-500/10 border border-emerald-500/20 rounded-2xl text-emerald-600 dark:text-emerald-400 font-semibold text-sm flex items-center gap-2">
      <CheckCircle2 size="18" class="shrink-0 text-emerald-500" />
      <span>{{ successMessage }}</span>
    </div>
    <div v-if="errorMessage" class="p-4 bg-red-500/10 border border-red-500/20 rounded-2xl text-red-600 dark:text-red-400 font-semibold text-sm flex items-center gap-2">
      <AlertCircle size="18" class="shrink-0 text-red-500" />
      <span>{{ errorMessage }}</span>
    </div>

    <!-- Tabs Navigation -->
    <div class="flex items-center gap-2 border-b border-slate-200 dark:border-slate-700 pb-2">
      <button 
        @click="activeTab = 'general'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :class="activeTab === 'general' ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'"
      >
        <Sliders size="18" />
        <span>Cấu hình chung</span>
      </button>

      <button 
        @click="activeTab = 'ai'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :class="activeTab === 'ai' ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'"
      >
        <Bot size="18" />
        <span>Cấu hình AI & Prompts</span>
        <span class="px-2 py-0.5 rounded-full text-xs font-semibold bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300">
          {{ promptTemplates.length }}
        </span>
      </button>

      <button 
        @click="activeTab = 'security'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :class="activeTab === 'security' ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'"
      >
        <Shield size="18" />
        <span>Bảo mật & Phiên làm việc</span>
      </button>

      <button 
        @click="activeTab = 'notifications'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :class="activeTab === 'notifications' ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'"
      >
        <Bell size="18" />
        <span>Cấu hình Thông báo</span>
      </button>
    </div>

    <!-- Tab 1: General Settings -->
    <Card v-if="activeTab === 'general'" class="rounded-2xl shadow-sm p-8 space-y-6">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-4">
        <h3 class="text-lg font-bold text-slate-800 dark:text-white flex items-center gap-2">
          <Sliders class="text-blue-600 dark:text-blue-400" size="20" />
          Thông số hoạt động nền tảng
        </h3>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          Các thiết lập sẽ áp dụng tức thì cho toàn bộ các dịch vụ tuyển dụng và phỏng vấn AI.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- System Name -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-1">
            Tên Hệ sinh thái / Platform Name
          </label>
          <span class="block text-[11px] text-slate-500 dark:text-slate-400 mb-2">
            Tên đầy đủ của cả nền tảng máy chủ (Dùng trong báo cáo nội bộ, log kỹ thuật server & email hệ thống)
          </span>
          <input 
            type="text" 
            v-model="settings.system_name"
            placeholder="Ví dụ: ViệcLàm AI Platform"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>

        <!-- Brand Name -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-1">
            Tên Thương hiệu chính / Brand Name
          </label>
          <span class="block text-[11px] text-blue-600 dark:text-blue-400 font-semibold mb-2">
            ⭐ Tên hiển thị trực tiếp ra giao diện (Navbar Logo, Footer, Tab Website, Admin Sidebar)
          </span>
          <input 
            type="text" 
            v-model="settings.brand_name"
            placeholder="Ví dụ: ViệcLàm, Talent, TuyểnDụng"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>

        <!-- Brand Badge -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
            Huy hiệu Logo / Brand Badge
          </label>
          <input 
            type="text" 
            v-model="settings.brand_badge"
            placeholder="Ví dụ: AI, PRO, VN"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>

        <!-- Support Email -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
            Email hỗ trợ hệ thống / Support Email
          </label>
          <input 
            type="email" 
            v-model="settings.support_email"
            placeholder="Ví dụ: support@vieclam.ai"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>

        <!-- Brand Slogan (Full width) -->
        <div class="md:col-span-2">
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
            Slogan & Tiêu đề phụ Website / Brand Slogan
          </label>
          <input 
            type="text" 
            v-model="settings.brand_slogan"
            placeholder="Ví dụ: Nền tảng Phỏng vấn & Tuyển dụng Thông minh"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>

        <!-- Custom Logo URL (Full width) -->
        <div class="md:col-span-2">
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2 flex items-center justify-between">
            <span>Đường dẫn ảnh Logo riêng / Custom Logo URL (PNG, JPG, SVG)</span>
            <span class="text-[11px] font-normal text-blue-600 dark:text-blue-400">Để trống nếu muốn dùng Icon Vector AI mặc định</span>
          </label>
          <input 
            type="text" 
            v-model="settings.brand_logo_url"
            placeholder="Ví dụ: /images/logo.png hoặc https://example.com/logo.svg"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>

        <!-- Default AI Model -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
            Mô hình AI Phỏng vấn mặc định
          </label>
          <select 
            v-model="settings.default_ai_model"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          >
            <option value="gemini-2.5-flash">Gemini 2.5 Flash (Tốc độ siêu nhanh - Khuyên dùng)</option>
            <option value="gemini-2.5-pro">Gemini 2.5 Pro (Phân tích chuyên sâu cao cấp)</option>
            <option value="gpt-4o">OpenAI GPT-4o (Dự phòng hệ thống)</option>
          </select>
        </div>

        <!-- Max File Size -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
            Dung lượng file CV/Giấy phép tối đa (MB)
          </label>
          <input 
            type="number" 
            v-model.number="settings.max_upload_size_mb"
            min="1" max="100"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>

        <!-- Default Passing Score -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
            Điểm Đạt (Pass Score) Phỏng vấn AI tối thiểu / 100
          </label>
          <input 
            type="number" 
            v-model.number="settings.default_passing_score"
            min="10" max="100"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>
      </div>

      <!-- Maintenance Mode Toggle -->
      <div class="pt-4 border-t border-slate-200 dark:border-slate-700 flex items-center justify-between">
        <div>
          <h4 class="font-bold text-slate-800 dark:text-white text-sm">Chế độ bảo trì hệ thống (Maintenance Mode)</h4>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
            Khi bật, toàn bộ ứng viên và nhà tuyển dụng sẽ nhận thông báo nâng cấp. Chỉ Admin mới có quyền truy cập.
          </p>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" v-model="settings.maintenance_mode" class="sr-only peer">
          <div class="w-14 h-7 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all dark:border-slate-600 peer-checked:bg-amber-600"></div>
        </label>
      </div>

      <!-- Save Button -->
      <div class="pt-6 border-t border-slate-200 dark:border-slate-700 flex justify-end">
        <button 
          @click="saveSettings" 
          :disabled="saving"
          class="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm rounded-xl shadow-lg shadow-blue-600/30 transition-all disabled:opacity-50 flex items-center gap-2"
        >
          <Save size="18" />
          {{ saving ? 'Đang lưu thiết lập...' : 'Lưu Thay đổi Cấu hình' }}
        </button>
      </div>
    </Card>

    <!-- Tab 2: AI & Prompt Templates -->
    <div v-if="activeTab === 'ai'" class="space-y-6">
      <!-- Action Bar -->
      <Card class="p-6 rounded-2xl shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h3 class="text-lg font-bold text-slate-800 dark:text-white flex items-center gap-2">
            <Cpu class="text-blue-600 dark:text-blue-400" size="20" />
            Thư viện Prompt Templates AI
          </h3>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
            Mỗi lần cập nhật nội dung prompt sẽ tự động tạo một phiên bản mới (Version history) trong cơ sở dữ liệu.
          </p>
        </div>
        <button 
          @click="openNewPromptModal"
          class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm rounded-xl shadow-md shadow-blue-600/20 transition-all flex items-center gap-2 shrink-0"
        >
          <Plus size="18" /> Thêm Prompt Template mới
        </button>
      </Card>

      <!-- Prompts Grid -->
      <div v-if="loadingPrompts" class="py-12 text-center text-slate-500">
        <RefreshCw size="24" class="animate-spin mx-auto mb-2 text-blue-600" />
        Đang tải danh sách Prompt Templates...
      </div>

      <Card v-else-if="promptTemplates.length === 0" class="p-12 rounded-2xl text-center border">
        <Bot size="36" class="mx-auto mb-3 text-slate-400" />
        <h4 class="font-bold text-slate-800 dark:text-white">Chưa có Prompt Template nào được tùy biến</h4>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 mb-4">Hệ thống đang sử dụng Prompt AI mặc định được nạp từ mã nguồn.</p>
        <button @click="openNewPromptModal" class="px-4 py-2 bg-blue-600 text-white font-bold text-xs rounded-xl">Tạo Prompt mẫu ngay</button>
      </Card>

      <div v-else class="grid grid-cols-1 gap-4">
        <div 
          v-for="tmpl in promptTemplates" 
          :key="tmpl.id"
          class="bg-[var(--surface)] p-6 rounded-2xl border border-[var(--border)] shadow-sm hover:border-blue-400/60 transition-all flex flex-col md:flex-row md:items-center justify-between gap-4"
        >
          <div class="space-y-2 flex-1">
            <div class="flex items-center gap-3">
              <span class="px-3 py-1 rounded-lg bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300 font-extrabold font-mono text-xs">
                v{{ tmpl.version || 1 }}
              </span>
              <h4 class="text-base font-extrabold text-slate-800 dark:text-white">{{ tmpl.name }}</h4>
              <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400">
                Active
              </span>
              <span class="text-xs text-slate-400 font-mono">Model: {{ tmpl.model || 'gemini-2.5-flash' }}</span>
            </div>
            <p class="text-xs text-slate-600 dark:text-slate-300 font-mono line-clamp-2 bg-slate-50 dark:bg-slate-900/60 p-3 rounded-xl border border-slate-200/50 dark:border-slate-700/50">
              {{ tmpl.content }}
            </p>
          </div>
          
          <div class="flex items-center gap-2 shrink-0">
            <button 
              @click="openEditPromptModal(tmpl)"
              class="px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-200 font-bold text-xs rounded-xl transition-colors flex items-center gap-1.5"
            >
              <Edit3 size="14" /> Chỉnh sửa / Tạo bản v{{ (tmpl.version || 1) + 1 }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tab 3: Security & Sessions -->
    <Card v-if="activeTab === 'security'" class="rounded-2xl shadow-sm p-8 space-y-6">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-4">
        <h3 class="text-lg font-bold text-slate-800 dark:text-white flex items-center gap-2">
          <Shield class="text-blue-600 dark:text-blue-400" size="20" />
          Chính sách bảo mật phiên làm việc (Sessions & Auth)
        </h3>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          Thiết lập tiêu chuẩn an toàn truy cập cho các tài khoản Quản trị và Nhà tuyển dụng.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- JWT Expiry -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2 flex items-center gap-1.5">
            <Key size="14" class="text-blue-500" /> Thời gian hiệu lực JWT Token (Giờ)
          </label>
          <input 
            type="number" 
            v-model.number="settings.jwt_token_expiry_hours"
            min="1" max="168"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
          <span class="text-[11px] text-slate-400 mt-1.5 block">Mặc định: 24 giờ. Token hết hạn sẽ tự động làm mới bằng Refresh Cookie.</span>
        </div>

        <!-- Admin 2FA Toggle -->
        <div class="bg-slate-50 dark:bg-slate-900/60 p-5 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 flex items-center justify-between">
          <div>
            <h4 class="font-bold text-slate-800 dark:text-white text-sm flex items-center gap-1.5">
              <Shield size="16" class="text-emerald-500" /> Xác thực 2 bước (2FA Admin Access)
            </h4>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
              Yêu cầu xác nhận qua OTP email khi truy cập các tính năng nâng cao.
            </p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" v-model="settings.admin_2fa_required" class="sr-only peer">
            <div class="w-14 h-7 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all dark:border-slate-600 peer-checked:bg-blue-600"></div>
          </label>
        </div>
      </div>

      <div class="pt-6 border-t border-slate-200 dark:border-slate-700 flex justify-end">
        <button 
          @click="saveSettings" 
          :disabled="saving"
          class="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm rounded-xl shadow-lg shadow-blue-600/30 transition-all disabled:opacity-50 flex items-center gap-2"
        >
          <Save size="18" />
          {{ saving ? 'Đang lưu thiết lập...' : 'Lưu Thay đổi Bảo mật' }}
        </button>
      </div>
    </Card>

    <!-- Tab 4: Notification Settings -->
    <Card v-if="activeTab === 'notifications'" class="rounded-2xl shadow-sm p-8 space-y-6 animate-fade-in">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-4">
        <h3 class="text-lg font-bold text-slate-800 dark:text-white flex items-center gap-2">
          <Bell class="text-blue-600 dark:text-blue-400" size="20" />
          Cấu hình Thông báo & Vòng đời trên Redis
        </h3>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
          Quản lý thời gian lưu trữ trong bộ nhớ Redis, giới hạn số lượng và kiểm soát bật/tắt các loại sự kiện thông báo hệ thống.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Redis TTL -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2 flex items-center gap-1.5">
            Thời gian lưu trữ trong Redis / TTL (Ngày)
          </label>
          <input 
            type="number" 
            v-model.number="settings.notification_ttl_days"
            min="1" max="365"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
          <span class="text-[11px] text-slate-400 mt-1.5 block">Mặc định: 30 ngày. Thông báo cũ hơn TTL sẽ tự động được xóa khỏi bộ nhớ Redis.</span>
        </div>

        <!-- Max items per user -->
        <div>
          <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2 flex items-center gap-1.5">
            Giới hạn số lượng tối đa mỗi User (Items)
          </label>
          <input 
            type="number" 
            v-model.number="settings.notification_max_per_user"
            min="20" max="1000"
            class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
          <span class="text-[11px] text-slate-400 mt-1.5 block">Mặc định: 200 thông báo gần nhất. Dùng lệnh LTRIM để tối ưu dung lượng RAM Redis.</span>
        </div>
      </div>

      <div class="border-t border-slate-200 dark:border-slate-700 pt-6">
        <h4 class="font-bold text-slate-800 dark:text-white text-sm mb-4">Bật / Tắt theo sự kiện (Event Toggles)</h4>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Notify on new applicant -->
          <div class="bg-slate-50 dark:bg-slate-900/60 p-4 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 flex items-center justify-between">
            <div>
              <h5 class="font-bold text-slate-800 dark:text-white text-sm">Ứng viên mới nộp CV</h5>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Thông báo khi có ứng viên mới apply vào tin tuyển dụng.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="settings.notify_on_new_applicant" class="sr-only peer">
              <div class="w-12 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-slate-600 peer-checked:bg-blue-600"></div>
            </label>
          </div>

          <!-- Notify on report ready -->
          <div class="bg-slate-50 dark:bg-slate-900/60 p-4 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 flex items-center justify-between">
            <div>
              <h5 class="font-bold text-slate-800 dark:text-white text-sm">Báo cáo AI phỏng vấn sẵn sàng</h5>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Thông báo khi AI hoàn tất chấm điểm và xuất báo cáo.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="settings.notify_on_report_ready" class="sr-only peer">
              <div class="w-12 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-slate-600 peer-checked:bg-blue-600"></div>
            </label>
          </div>

          <!-- Notify on interview cancelled -->
          <div class="bg-slate-50 dark:bg-slate-900/60 p-4 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 flex items-center justify-between">
            <div>
              <h5 class="font-bold text-slate-800 dark:text-white text-sm">Lịch phỏng vấn bị hủy/thay đổi</h5>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Thông báo khi ứng viên hoặc HR hủy lịch phỏng vấn.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="settings.notify_on_interview_cancelled" class="sr-only peer">
              <div class="w-12 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-slate-600 peer-checked:bg-blue-600"></div>
            </label>
          </div>

          <!-- Enable email notifications -->
          <div class="bg-slate-50 dark:bg-slate-900/60 p-4 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 flex items-center justify-between">
            <div>
              <h5 class="font-bold text-slate-800 dark:text-white text-sm">Gửi song song qua Email</h5>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Gửi email kèm thông báo cho các sự kiện quan trọng.</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="settings.enable_email_notifications" class="sr-only peer">
              <div class="w-12 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer dark:bg-slate-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-slate-600 peer-checked:bg-blue-600"></div>
            </label>
          </div>
        </div>
      </div>

      <div class="pt-6 border-t border-slate-200 dark:border-slate-700 flex justify-end">
        <button 
          @click="saveSettings" 
          :disabled="saving"
          class="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm rounded-xl shadow-lg shadow-blue-600/30 transition-all disabled:opacity-50 flex items-center gap-2"
        >
          <Save size="18" />
          {{ saving ? 'Đang lưu thiết lập...' : 'Lưu Thay đổi Cấu hình Thông báo' }}
        </button>
      </div>
    </Card>

    <!-- Modal for Create/Edit Prompt -->
    <div v-if="showPromptModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/70 backdrop-blur-sm animate-fade-in">
      <div class="bg-[var(--surface)] rounded-3xl max-w-2xl w-full p-8 shadow-2xl border border-[var(--border)] space-y-6 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between border-b border-slate-200 dark:border-slate-700 pb-4">
          <h3 class="text-xl font-bold text-slate-800 dark:text-white flex items-center gap-2">
            <Cpu class="text-blue-600 dark:text-blue-400" size="24" />
            Cấu hình Prompt Template AI
          </h3>
          <button @click="showPromptModal = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 text-sm font-bold">✕ Đóng</button>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">Tên định danh Prompt Template</label>
            <input 
              type="text" 
              v-model="editingPrompt.name" 
              placeholder="e.g. RUBRIC_EVALUATION, INTERVIEW_CONDUCTOR"
              class="w-full px-4 py-3 font-mono text-sm bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">Mô hình AI xử lý</label>
            <select 
              v-model="editingPrompt.model"
              class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-semibold text-slate-800 dark:text-white"
            >
              <option value="gemini-2.5-flash">Gemini 2.5 Flash (Tốc độ tối ưu)</option>
              <option value="gemini-2.5-pro">Gemini 2.5 Pro (Độ chính xác cao)</option>
              <option value="gpt-4o">OpenAI GPT-4o</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider mb-2">
              Nội dung System Prompt (Có thể sử dụng biến {{ '{' + '{ variable_name }' + '}' }})
            </label>
            <textarea 
              v-model="editingPrompt.content" 
              rows="8"
              placeholder="Nhập hướng dẫn chi tiết cho mô hình AI..."
              class="w-full p-4 font-mono text-xs bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 leading-relaxed"
            ></textarea>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-200 dark:border-slate-700">
          <button @click="showPromptModal = false" class="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-300 font-bold text-sm rounded-xl transition-colors">
            Hủy bỏ
          </button>
          <button @click="savePromptTemplate" :disabled="saving" class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm rounded-xl shadow-lg shadow-blue-600/30 transition-all disabled:opacity-50 flex items-center gap-2">
            <Save size="16" /> {{ saving ? 'Đang lưu phiên bản mới...' : 'Lưu & Khởi tạo Version mới' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
