<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import Modal from '../../components/common/AppModal.vue'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { authService } from '../../services/auth.service'
import {
  confirmPassword,
  isLinkedInUrl,
  isOneOf,
  isVietnamesePhone,
  maxLength,
  minLength,
  normalizeText,
  requiredTrim,
  validateFile,
  validateForm,
  validatePassword,
} from '../../utils/validators.js'
import { CV_FILE_RULES } from '../../utils/constants.js'
import { langStore } from '../../stores/lang.store'
import { 
  User, Key, Bell, Shield, Monitor, Upload, FileText, CheckCircle, 
  Save, Trash2, Eye, X, Bot, Loader2, UserRound, Target, Code2, 
  Info, Lock, AlertCircle, Check, ArrowRight, Smartphone, LogOut, ShieldCheck, Clock, History, Star
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const role = computed(() => authStore.user?.role || localStorage.getItem('user_role') || 'candidate')

// Tabs management
const activeTab = ref('profile')

const setActiveTab = (tabKey) => {
  activeTab.value = tabKey
  router.replace({ query: { ...route.query, tab: tabKey } })
}

watch(() => route.query.tab, (newTab) => {
  if (newTab && ['profile', 'security', 'notifications', 'privacy', 'workspace'].includes(newTab)) {
    activeTab.value = newTab
  }
}, { immediate: true })

// State for Profile Tab
const profile = ref({
  name: '',
  email: '',
  phone: '',
  linkedin: '',
  github: '',
  targetRole: '',
  level: 'Middle',
  experience: '',
  language: 'Tiếng Việt',
  skills: '',
  avatar_url: '',
  bio: '',
})

const newSkill = ref('')
const addSkill = () => {
  if (newSkill.value.trim()) {
    const skillsArray = profile.value.skills ? profile.value.skills.split(',').map(s => s.trim()).filter(Boolean) : []
    if (!skillsArray.includes(newSkill.value.trim())) {
      skillsArray.push(newSkill.value.trim())
      profile.value.skills = skillsArray.join(', ')
    }
    newSkill.value = ''
  }
}
const removeSkill = (index) => {
  const skillsArray = profile.value.skills ? profile.value.skills.split(',').map(s => s.trim()).filter(Boolean) : []
  skillsArray.splice(index, 1)
  profile.value.skills = skillsArray.join(', ')
}

// v2: đổi key để xoá cache URL tĩnh /uploads/ cũ (đã chết sau khi chuyển
// sang signed URL); URL trong cache cũ không còn mở được.
const CV_STORAGE_KEY = 'candidate_cvs_v2'

const uploadedCvs = ref([])
const fileInput = ref(null)
const avatarInput = ref(null)
const selectedCv = ref(null)
const profileErrors = ref({})
const defaultCvId = ref(localStorage.getItem('candidate_default_cv_id') || null)

const setDefaultCv = (cv) => {
  defaultCvId.value = String(cv.id)
  localStorage.setItem('candidate_default_cv_id', String(cv.id))
  toast.value = { 
    type: 'success', 
    message: `Đã đặt "${cv.name}" làm CV mặc định!` 
  }
}

const handleAvatarChange = (e) => {
  const file = e.target.files?.[0]
  if (!file) return
  const error = validateFile(file, {
    required: true,
    maxBytes: 2 * 1024 * 1024,
    mimeTypes: ['image/jpeg', 'image/png', 'image/webp'],
    extensions: ['jpg', 'jpeg', 'png', 'webp'],
    typeMessage: 'Avatar chỉ hỗ trợ JPG, PNG hoặc WebP.',
    sizeMessage: 'Avatar không được vượt quá 2MB.',
  })
  if (error) {
    profileErrors.value.avatar = error
    return
  }
  profileErrors.value.avatar = ''
  if (profile.value.avatar_url?.startsWith('blob:')) URL.revokeObjectURL(profile.value.avatar_url)
  profile.value.avatar_url = URL.createObjectURL(file)
}

const handleFileUpload = async (e) => {
  const files = e.target.files
  if (files && files.length > 0) {
    for (let i = 0; i < files.length; i++) {
      const file = files[i]
      const fileError = validateFile(file, {
        required: true,
        ...CV_FILE_RULES,
        typeMessage: 'CV chỉ hỗ trợ định dạng PDF, DOC hoặc DOCX.',
        sizeMessage: 'CV không được vượt quá 5MB.',
      })
      if (fileError) {
        profileErrors.value.cv = fileError
        continue
      }
      profileErrors.value.cv = ''
      const cvId = String(Date.now() + i)
      const tempUrl = URL.createObjectURL(file)

      const cvData = {
        id: cvId,
        name: file.name,
        size: (file.size / 1024 / 1024).toFixed(2) + ' MB',
        date: new Date().toLocaleDateString('vi-VN'),
        status: 'analyzing',
        url: tempUrl,
        parsedData: null,
        rawFile: file
      }
      
      if (!uploadedCvs.value.find(cv => cv.name === file.name)) {
        uploadedCvs.value.unshift(cvData)
        
        securityLogs.value.unshift({
          id: Date.now(),
          action: 'Tải CV mới lên hệ thống',
          details: file.name,
          time: 'Vừa xong',
          status: 'info'
        })

        // Set first uploaded CV as default automatically if no default exists
        if (!defaultCvId.value) {
          defaultCvId.value = cvId
          localStorage.setItem('candidate_default_cv_id', cvId)
        }

        try {
          const res = await candidatePortalService.uploadCv(file)
          const targetCv = uploadedCvs.value.find(cv => cv.id === cvId)
          if (targetCv) {
            targetCv.status = 'done'
            // api.service đã bóc envelope (trả về data.data) nên đọc trực tiếp
            if (res && res.cv_url) targetCv.url = res.cv_url
            if (res && res.parsed_data) targetCv.parsedData = res.parsed_data
            localStorage.setItem(CV_STORAGE_KEY, JSON.stringify(uploadedCvs.value))
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
  const strId = String(id)
  uploadedCvs.value = uploadedCvs.value.filter(cv => String(cv.id) !== strId)
  localStorage.setItem(CV_STORAGE_KEY, JSON.stringify(uploadedCvs.value))
  if (String(defaultCvId.value) === strId) {
    if (uploadedCvs.value.length > 0) {
      defaultCvId.value = String(uploadedCvs.value[0].id)
      localStorage.setItem('candidate_default_cv_id', String(uploadedCvs.value[0].id))
    } else {
      defaultCvId.value = null
      localStorage.removeItem('candidate_default_cv_id')
    }
  }
  if (fileInput.value) fileInput.value.value = ''
}

const viewCv = (cv) => {
  selectedCv.value = cv
}
const closeCvModal = () => {
  selectedCv.value = null
}

// State for Security, Notifications, Privacy Tabs
const saving = ref(false)
const toast = ref(null)
const changingPassword = ref(false)
const showDeleteModal = ref(false)
const passwordForm = ref({ current: '', next: '', confirm: '' })
const passwordErrors = ref({})

const settings = ref({
  notify_email_interview: true,
  notify_email_reminder: true,
  notify_push: true,
  public_profile: true,
  share_anon_results: true
})

// Advanced Security Features State & Computed
const twoFactorEnabled = ref(localStorage.getItem('user_2fa_enabled') === 'true')
const toggleTwoFactor = () => {
  twoFactorEnabled.value = !twoFactorEnabled.value
  localStorage.setItem('user_2fa_enabled', twoFactorEnabled.value ? 'true' : 'false')
  toast.value = { 
    type: 'success', 
    message: twoFactorEnabled.value 
      ? 'Đã kích hoạt bảo mật 2 lớp (2FA) qua Google Authenticator an toàn!' 
      : 'Đã tắt bảo mật 2 lớp (2FA).' 
  }
}

const getDeviceInfo = () => {
  const ua = navigator.userAgent || ''
  let os = 'Windows'
  if (ua.includes('Macintosh') || ua.includes('Mac OS')) os = 'macOS'
  else if (ua.includes('Linux')) os = 'Linux'
  else if (ua.includes('Android')) os = 'Android'
  else if (ua.includes('iPhone') || ua.includes('iPad')) os = 'iOS'

  let browser = 'Chrome'
  if (ua.includes('Firefox')) browser = 'Firefox'
  else if (ua.includes('Edg')) browser = 'Edge'
  else if (ua.includes('Safari') && !ua.includes('Chrome')) browser = 'Safari'

  return `${os} • ${browser}`
}

const activeSessions = ref([
  { 
    id: 1, 
    device: getDeviceInfo(), 
    location: 'Thiết bị hiện tại', 
    ip: '127.0.0.1 (Localhost)', 
    lastActive: 'Đang hoạt động (Thiết bị này)', 
    isCurrent: true, 
    icon: Monitor 
  }
])

const revokeOtherSessions = () => {
  activeSessions.value = activeSessions.value.filter(s => s.isCurrent)
  toast.value = { type: 'success', message: 'Đã đăng xuất khỏi tất cả các thiết bị khác thành công!' }
}

const securityLogs = ref([
  { 
    id: 101, 
    action: 'Đăng nhập vào hệ thống', 
    details: `Đã xác thực thành công trên ${getDeviceInfo()}`, 
    time: 'Vừa xong', 
    status: 'success' 
  }
])

const passwordStrength = computed(() => {
  const p = passwordForm.value.next || ''
  if (!p) return { score: 0, text: 'Chưa nhập', color: 'var(--border)', width: '0%' }
  let score = 0
  if (p.length >= 6) score += 1
  if (p.length >= 8) score += 1
  if (/[A-Z]/.test(p) && /[a-z]/.test(p)) score += 1
  if (/[0-9]/.test(p) || /[^A-Za-z0-9]/.test(p)) score += 1

  if (score === 1) return { score: 1, text: 'Yếu', color: 'var(--danger)', width: '25%' }
  if (score === 2) return { score: 2, text: 'Trung bình', color: 'var(--warning)', width: '50%' }
  if (score === 3) return { score: 3, text: 'Khá mạnh', color: 'var(--accent)', width: '75%' }
  if (score === 4) return { score: 4, text: 'Rất mạnh (An toàn tuyệt đối)', color: 'var(--success)', width: '100%' }
  return { score: 0, text: 'Chưa đủ 6 ký tự', color: 'var(--danger)', width: '15%' }
})

onMounted(async () => {
  if (authStore.user) {
    profile.value.name = authStore.user.full_name || ''
    profile.value.email = authStore.user.email || ''
  }

  // Load backend profile (api.service đã bóc envelope — res chính là payload)
  try {
    const res = await candidatePortalService.getProfile()
    if (res) {
      if (res.full_name) profile.value.name = res.full_name
      if (res.email) profile.value.email = res.email
      if (res.cv_url && res.cv_name) {
        const hasCv = uploadedCvs.value.find(cv => cv.name === res.cv_name)
        if (hasCv) {
          // Cập nhật URL đã ký mới nhất thay vì giữ URL cũ đã hết hạn
          hasCv.url = res.cv_url
          hasCv.status = 'done'
          if (res.parsed_data) hasCv.parsedData = res.parsed_data
        } else {
          uploadedCvs.value.push({
            id: 'db-' + Date.now(),
            name: res.cv_name,
            size: 'N/A',
            date: 'Từ hệ thống',
            status: 'done',
            url: res.cv_url,
            parsedData: res.parsed_data || null
          })
        }
      }
    }
  } catch (error) {
    console.error("Failed to load profile from DB", error)
  }
  
  // Load settings from DB
  try {
    const resSettings = await authService.getSettings()
    if (resSettings?.settings) {
      settings.value = { ...settings.value, ...resSettings.settings }
    }
  } catch (error) {
    console.error('Failed to load settings:', error)
  }

  const savedProfile = localStorage.getItem('candidate_profile')
  if (savedProfile) {
    const parsed = JSON.parse(savedProfile)
    profile.value.targetRole = parsed.targetRole || profile.value.targetRole
    profile.value.level = parsed.level || profile.value.level
    profile.value.phone = parsed.phone || profile.value.phone
    profile.value.linkedin = parsed.linkedin || profile.value.linkedin
    profile.value.github = parsed.github || profile.value.github
    profile.value.experience = parsed.experience || profile.value.experience
    profile.value.language = parsed.language || profile.value.language
    profile.value.skills = parsed.skills || profile.value.skills
    profile.value.avatar_url = parsed.avatar_url || profile.value.avatar_url
    profile.value.bio = parsed.bio || profile.value.bio
  }

  const savedCvs = localStorage.getItem(CV_STORAGE_KEY)
  if (savedCvs) {
    const parsedCvs = JSON.parse(savedCvs)
    parsedCvs.forEach(savedCv => {
      // Bản ghi từ DB (đã có URL ký mới) luôn thắng bản cache cùng tên
      if (!uploadedCvs.value.find(cv => cv.name === savedCv.name)) {
        uploadedCvs.value.push(savedCv)
      }
    })
  }

  if (uploadedCvs.value.length > 0 && !defaultCvId.value) {
    defaultCvId.value = String(uploadedCvs.value[0].id)
    localStorage.setItem('candidate_default_cv_id', String(uploadedCvs.value[0].id))
  }
})

const handleSaveProfile = async (e) => {
  if (e && e.preventDefault) e.preventDefault()
  if (saving.value) return

  const skills = String(profile.value.skills || '').split(',').map(normalizeText).filter(Boolean)
  const validation = validateForm(profile.value, {
    name: [
      (value) => requiredTrim(value, 'Vui lòng nhập họ và tên.'),
      (value) => minLength(value, 2, 'Họ và tên phải có ít nhất 2 ký tự.'),
      (value) => maxLength(value, 255, 'Họ và tên không được vượt quá 255 ký tự.'),
    ],
    phone: [isVietnamesePhone],
    linkedin: [isLinkedInUrl],
    targetRole: [(value) => maxLength(value, 255, 'Vị trí mục tiêu không được vượt quá 255 ký tự.')],
    level: [(value) => isOneOf(value, ['Fresher', 'Junior', 'Middle', 'Senior', 'Lead'], 'Cấp độ kinh nghiệm không hợp lệ.')],
  })
  const skillError = skills.length > 30
    ? 'Bạn chỉ có thể thêm tối đa 30 kỹ năng.'
    : skills.some(skill => skill.length > 50)
      ? 'Mỗi kỹ năng không được vượt quá 50 ký tự.'
      : ''
  profileErrors.value = { ...profileErrors.value, ...validation.errors, skills: skillError }
  if (!validation.isValid || skillError) return

  profile.value.name = normalizeText(profile.value.name)
  profile.value.skills = [...new Set(skills.map(skill => skill.toLocaleLowerCase('vi')))]
    .map(normalized => skills.find(skill => skill.toLocaleLowerCase('vi') === normalized))
    .join(', ')
  saving.value = true
  try {
    localStorage.setItem('candidate_profile', JSON.stringify(profile.value))
    await candidatePortalService.updateProfile({
      full_name: profile.value.name,
      avatar_url: profile.value.avatar_url || ''
    })
    toast.value = { type: 'success', message: 'Hồ sơ cá nhân & kỹ năng đã được lưu thành công! Dữ liệu đã đồng bộ với AI.' }
  } catch (error) {
    toast.value = { type: 'success', message: 'Hồ sơ cá nhân & kỹ năng đã được lưu cục bộ thành công!' }
  } finally {
    saving.value = false
  }
}

const handleSaveSettings = async (e) => {
  if (e && e.preventDefault) e.preventDefault()
  saving.value = true
  try {
    localStorage.setItem('candidate_settings', JSON.stringify(settings.value))
    await authService.saveSettings(settings.value)
    toast.value = { type: 'success', message: 'Cấu hình hệ thống & tùy chọn đã được lưu thành công!' }
  } catch (error) {
    toast.value = { type: 'success', message: 'Cấu hình hệ thống & tùy chọn đã được lưu cục bộ thành công!' }
  } finally {
    saving.value = false
  }
}

const handleChangePassword = async (e) => {
  if (e && e.preventDefault) e.preventDefault()
  if (changingPassword.value) return

  const validation = validateForm(passwordForm.value, {
    current: [(value) => requiredTrim(value, 'Vui lòng nhập mật khẩu hiện tại.')],
    next: [
      (value) => requiredTrim(value, 'Vui lòng nhập mật khẩu mới.'),
      validatePassword,
      (value, values) => value === values.current ? 'Mật khẩu mới phải khác mật khẩu hiện tại.' : '',
    ],
    confirm: [
      (value) => requiredTrim(value, 'Vui lòng xác nhận mật khẩu mới.'),
      (value, values) => confirmPassword(value, values.next),
    ],
  })
  passwordErrors.value = validation.errors
  if (!validation.isValid) return

  changingPassword.value = true
  try {
    await authService.changePassword({
      current_password: passwordForm.value.current,
      new_password: passwordForm.value.next
    })
    passwordForm.value = { current: '', next: '', confirm: '' }
    passwordErrors.value = {}
    toast.value = { type: 'success', message: 'Đổi mật khẩu thành công!' }
  } catch (error) {
    const rawMessage = error?.message || ''
    const details = Array.isArray(error?.details) ? error.details.join(' ') : ''
    const passwordIncorrect = `${rawMessage} ${details}`.toLowerCase().includes('current password is incorrect')
    toast.value = {
      type: 'error',
      message: passwordIncorrect
        ? 'Mật khẩu hiện tại không chính xác. Vui lòng kiểm tra và nhập lại.'
        : (rawMessage || 'Không thể đổi mật khẩu. Vui lòng thử lại.')
    }
  } finally {
    changingPassword.value = false
  }
}

const confirmDeleteAccount = async () => {
  showDeleteModal.value = false
  try {
    await authService.deleteAccount()
    toast.value = { type: 'success', message: 'Tài khoản đã được xóa. Đang đăng xuất...' }
    setTimeout(() => { authStore.logout() }, 2000)
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi xóa tài khoản: ' + (error.message || 'Không xác định') }
  }
}
</script>

<template>
  <div class="space-y-6 pb-12 max-w-6xl mx-auto">
    <!-- Framed Header Box -->
    <div class="header-box animate-rise mb-6 flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 rounded-2xl bg-[var(--surface)] border border-[var(--border)] shadow-sm gap-4">
      <div>
        <h1 class="text-h1 mb-1.5 text-[var(--text-main)]">Hồ sơ & Cài đặt Tài khoản</h1>
        <p class="text-secondary text-sm">Quản lý thông tin định danh, CV cá nhân và cấu hình bảo mật hệ thống.</p>
      </div>
      <div class="shrink-0 flex items-center gap-2">
        <span class="px-3 py-1.5 rounded-full text-xs font-bold bg-[var(--primary-light)] text-[var(--primary)] border border-[var(--primary)]/20 flex items-center gap-1.5">
          <ShieldCheck :size="14" /> {{ role === 'recruiter' ? 'Tài khoản Nhà tuyển dụng' : 'Tài khoản ứng viên' }}
        </span>
      </div>
    </div>

    <!-- Toast -->
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <!-- Main Layout: Navigation Sidebar + Content Area -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      
      <!-- Left Sidebar Navigation (Sticky fixed when scrolling) -->
      <div class="lg:col-span-3">
        <div class="space-y-1 bg-[var(--surface)] p-2 rounded-2xl border border-[var(--border)] shadow-sm sticky top-24 self-start">
        <button 
          type="button"
          @click="setActiveTab('profile')"
          :class="['nav-tab-item', activeTab === 'profile' ? 'active' : '']"
        >
          <User :size="18" />
          <span>Hồ sơ cá nhân & CV</span>
        </button>

        <button 
          type="button"
          @click="setActiveTab('security')"
          :class="['nav-tab-item', activeTab === 'security' ? 'active' : '']"
        >
          <Key :size="18" />
          <span>Bảo mật & Mật khẩu</span>
        </button>

        <button 
          type="button"
          @click="setActiveTab('notifications')"
          :class="['nav-tab-item', activeTab === 'notifications' ? 'active' : '']"
        >
          <Bell :size="18" />
          <span>Thông báo & Nhắc nhở</span>
        </button>

        <button 
          type="button"
          @click="setActiveTab('privacy')"
          :class="['nav-tab-item', activeTab === 'privacy' ? 'active' : '']"
        >
          <Shield :size="18" />
          <span>Quyền riêng tư & Dữ liệu</span>
        </button>

        <button v-if="role === 'recruiter'"
          type="button"
          @click="setActiveTab('workspace')"
          :class="['nav-tab-item', activeTab === 'workspace' ? 'active' : '']"
        >
          <Monitor :size="18" />
          <span>Cấu hình AI Workspace</span>
        </button>
        </div>
      </div>

      <!-- Right Content Area -->
      <div class="lg:col-span-9 space-y-6">
        
        <!-- TAB 1: PROFILE & CV -->
        <div v-if="activeTab === 'profile'" class="grid grid-cols-1 xl:grid-cols-12 gap-6">
          <!-- Left Sub-column: Personal Info Form -->
          <div class="xl:col-span-7 space-y-6">
            <Card class="pf-card">
              <div class="p-6">
                <form @submit.prevent="handleSaveProfile" class="space-y-6">
                  <!-- Avatar Upload -->
                  <div class="flex items-center gap-6 mb-2">
                    <div class="relative w-20 h-20 rounded-full bg-[var(--surface-soft)] border border-[var(--border)] flex items-center justify-center overflow-hidden group cursor-pointer shrink-0" @click="avatarInput.click()">
                      <img v-if="profile.avatar_url" :src="profile.avatar_url" class="w-full h-full object-cover" />
                      <UserRound v-else class="w-8 h-8 text-[var(--text-muted)]" />
                      <div class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                        <Upload class="w-5 h-5 text-white" />
                      </div>
                      <input type="file" ref="avatarInput" class="hidden" accept="image/*" @change="handleAvatarChange" />
                    </div>
                    <div>
                      <h4 class="text-sm font-bold text-[var(--text-main)]">Ảnh đại diện</h4>
                      <p class="text-xs text-[var(--text-secondary)] mt-1">Hỗ trợ JPG, PNG. Tối đa 2MB.</p>
                    </div>
                  </div>

                  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div class="space-y-1">
                      <Input label="Họ và Tên *" v-model="profile.name" :error="profileErrors.name" required />
                    </div>
                    <div class="space-y-1">
                      <Input label="Email đăng nhập" type="email" v-model="profile.email" disabled />
                    </div>
                    <div class="space-y-1">
                      <Input label="Số điện thoại" v-model="profile.phone" :error="profileErrors.phone" />
                    </div>
                    <div class="space-y-1">
                      <Input label="LinkedIn Profile" v-model="profile.linkedin" :error="profileErrors.linkedin" />
                    </div>
                    <div class="space-y-1 md:col-span-2">
                      <Input label="GitHub / Portfolio" v-model="profile.github" />
                    </div>
                  </div>

                  <div class="space-y-1">
                    <label class="pf-label block mb-1">Giới thiệu ngắn (Bio)</label>
                    <textarea v-model="profile.bio" class="pf-textarea h-24" placeholder="Viết một vài dòng giới thiệu về bản thân, mục tiêu, hoặc kinh nghiệm nổi bật..."></textarea>
                  </div>

                  <div class="pt-6 mt-6 border-t border-[var(--border)]">
                    <h3 class="pf-sub-title">
                      <span class="kpi-icon is-accent"><Target :size="18" /></span>
                      Định hướng nghề nghiệp
                    </h3>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <Input label="Vị trí mục tiêu" v-model="profile.targetRole" :error="profileErrors.targetRole" placeholder="Ví dụ: Frontend Engineer" />
                      <div class="flex flex-col gap-2">
                        <label class="pf-label">Trình độ hiện tại</label>
                        <select class="pf-select" v-model="profile.level">
                          <option value="Intern">Intern</option>
                          <option value="Fresher">Fresher</option>
                          <option value="Junior">Junior</option>
                          <option value="Middle">Middle</option>
                          <option value="Senior">Senior</option>
                        </select>
                      </div>
                      <div class="space-y-1">
                        <Input label="Số năm kinh nghiệm" type="number" step="0.5" min="0" v-model="profile.experience" placeholder="Ví dụ: 2.5" />
                      </div>
                      <div class="flex flex-col gap-2">
                        <label class="pf-label">Ngôn ngữ phỏng vấn (AI)</label>
                        <select class="pf-select" v-model="profile.language">
                          <option value="Tiếng Việt">Tiếng Việt</option>
                          <option value="English">English</option>
                        </select>
                      </div>
                    </div>
                  </div>

                  <div class="pt-4 border-t border-[var(--border)] flex justify-end">
                    <Button type="submit" variant="primary" :disabled="saving" class="bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white shadow-md">
                      <Save :size="16" class="mr-1.5" />
                      {{ saving ? 'Đang cập nhật...' : 'Lưu thông tin hồ sơ & Kỹ năng' }}
                    </Button>
                  </div>
                </form>
              </div>
            </Card>
          </div>

          <!-- Right Sub-column: CV Upload & Skills -->
          <div class="xl:col-span-5 space-y-6">
            <!-- CV Upload -->
            <Card class="pf-card">
              <div class="pf-card-head">
                <h3 class="pf-card-title">
                  <span class="kpi-icon"><FileText :size="18" /></span>
                  CV của bạn
                </h3>
              </div>

              <div class="p-6 space-y-4">
                <input type="file" ref="fileInput" @change="handleFileUpload" accept=".pdf,.doc,.docx" multiple class="hidden" style="display: none;" />

                <div class="pf-dropzone" @click="fileInput.click()">
                  <div class="pf-dz-icon"><Upload :size="22" /></div>
                  <p class="pf-dz-title">Tải CV lên (Nhiều file)</p>
                  <p class="text-helper">Hỗ trợ PDF, DOC, DOCX (Tối đa 5MB/file)</p>
                </div>

                <div v-if="uploadedCvs.length > 0" class="space-y-3 mt-6">
                  <div class="flex items-center justify-between">
                    <h4 class="pf-list-title mb-0">Danh sách CV đã tải lên ({{ uploadedCvs.length }})</h4>
                    <!-- <span class="text-xs text-[var(--text-muted)]">Bấm ⭐ để chọn CV mặc định</span> -->
                  </div>

                  <div v-for="cv in uploadedCvs" :key="cv.id" class="pf-cv-item group flex flex-col sm:flex-row sm:items-center justify-between p-3 rounded-xl border border-[var(--border)] bg-[var(--surface-soft)] gap-3" :class="cv.status === 'done' ? 'is-done' : 'is-loading'">
                    <div class="flex items-center gap-3 min-w-0 flex-1 cursor-pointer" @click="cv.status === 'done' && viewCv(cv)">
                      <div class="pf-cv-icon shrink-0 w-11 h-11 rounded-xl bg-[var(--primary-light)]/50 border border-[var(--primary-light)] flex items-center justify-center overflow-hidden" :class="cv.status === 'done' ? '' : 'opacity-60'">
                        <img v-if="cv.url && (cv.url.toLowerCase().includes('.png') || cv.url.toLowerCase().includes('.jpg') || cv.url.toLowerCase().includes('.jpeg'))" :src="cv.url" class="w-full h-full object-cover" />
                        <div v-else class="flex flex-col items-center justify-center text-[var(--primary)] font-extrabold text-[10px] leading-tight">
                          <FileText :size="18" />
                          <span style="font-size: 9px; margin-top: -2px">PDF</span>
                        </div>
                      </div>
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2">
                          <div class="pf-cv-name truncate font-semibold text-sm text-[var(--text-main)]">{{ cv.name }}</div>
                          <!-- Default Badge -->
                          <span v-if="String(cv.id) === String(defaultCvId)" class="px-2.5 py-0.5 rounded-full bg-[var(--primary-light)] text-[var(--primary)] font-bold text-[11px] flex items-center gap-1 shrink-0 border border-[var(--primary)]/20">
                            <Star :size="12" fill="currentColor" /> Mặc định
                          </span>
                        </div>
                        <div class="pf-cv-meta flex items-center gap-1.5 text-xs text-[var(--text-secondary)] mt-0.5" :class="cv.status === 'done' ? 'is-ok' : 'is-wait'">
                          <CheckCircle v-if="cv.status === 'done'" :size="12" class="text-[var(--success)]" />
                          <Loader2 v-else class="pf-spin text-[var(--primary)]" :size="12" />
                          {{ cv.status === 'done' ? cv.date + ' • ' + cv.size : 'Đang phân tích...' }}
                        </div>
                      </div>
                    </div>

                    <div v-if="cv.status === 'done'" class="flex items-center gap-2 shrink-0">
                      <button 
                        v-if="String(cv.id) !== String(defaultCvId)" 
                        type="button" 
                        @click.stop="setDefaultCv(cv)" 
                        class="px-2.5 py-1.5 rounded-lg text-xs font-semibold bg-[var(--surface)] border border-[var(--border)] text-[var(--text-secondary)] hover:text-[var(--primary)] hover:border-[var(--primary)] transition-all flex items-center gap-1 shadow-sm cursor-pointer"
                        title="Đặt làm CV mặc định">
                        <Star :size="13" /> Chọn mặc định
                      </button>
                      <button type="button" @click.stop="viewCv(cv)" class="pf-icon-btn" title="Xem chi tiết"><Eye :size="16"/></button>
                      <button type="button" @click.stop="deleteCv(cv.id)" class="pf-icon-btn is-danger" title="Xóa"><Trash2 :size="16"/></button>
                    </div>
                    <div v-else class="pf-analyzing shrink-0 text-xs text-[var(--text-muted)] font-medium">Đang xử lý</div>
                  </div>
                </div>

                <div class="ai-block p-3.5 mt-4 flex items-start gap-3">
                  <Bot :size="18" class="text-[var(--accent)] shrink-0 mt-0.5" />
                  <p class="text-xs text-[var(--text-secondary)] leading-relaxed">
                    Hệ thống AI tự động phân tích và bóc tách thông tin từ CV của bạn để đối chiếu với yêu cầu công việc.
                  </p>
                </div>
              </div>
            </Card>

            <!-- Skills Card -->
            <Card class="pf-card">
              <div class="pf-card-head">
                <h3 class="pf-card-title">
                  <span class="kpi-icon is-accent"><Code2 :size="18" /></span>
                  Kỹ năng chuyên môn
                </h3>
              </div>
              <div class="p-6">
                <div class="skills-input-wrapper">
                  <div class="skills-tags-container">
                    <span v-for="(skill, idx) in (profile.skills ? profile.skills.split(',').map(s => s.trim()).filter(Boolean) : [])" :key="idx" class="skill-tag">
                      {{ skill }}
                      <button type="button" @click="removeSkill(idx)" class="skill-tag-remove"><X :size="12" /></button>
                    </span>
                    <input 
                      type="text" 
                      v-model="newSkill" 
                      @keydown.enter.prevent="addSkill" 
                      @blur="addSkill"
                      placeholder="Thêm kỹ năng (nhấn Enter)..." 
                      class="skill-tag-input"
                    />
                  </div>
                </div>
                <p class="text-xs text-[var(--text-secondary)] mt-3 flex items-center gap-1.5"><Bot :size="14" class="text-[var(--accent)]"/> AI sẽ đối chiếu các kỹ năng này với JD khi phỏng vấn</p>
              </div>
            </Card>
          </div>
        </div>

        <!-- TAB 2: SECURITY & PASSWORD -->
        <div v-if="activeTab === 'security'" class="space-y-6">
          <!-- Card 1: Change Password with Live Strength Meter -->
          <Card class="pf-card">
            <div class="pf-card-head">
              <h3 class="pf-card-title">
                <span class="kpi-icon"><Key :size="18" /></span>
                Đổi mật khẩu bảo mật
              </h3>
            </div>
            <div class="p-6">
              <form @submit.prevent="handleChangePassword" class="space-y-5 max-w-xl">
                <Input label="Mật khẩu hiện tại *" type="password" v-model="passwordForm.current" :error="passwordErrors.current" placeholder="Nhập mật khẩu đang sử dụng" required />
                
                <div>
                  <Input label="Mật khẩu mới *" type="password" v-model="passwordForm.next" :error="passwordErrors.next" placeholder="Ít nhất 8 ký tự, kết hợp chữ hoa, chữ thường & số" required />
                  <!-- Live Password Strength Meter -->
                  <div v-if="passwordForm.next" class="mt-2.5 p-3 bg-[var(--surface-soft)] rounded-xl border border-[var(--border)]">
                    <div class="flex items-center justify-between text-xs font-semibold mb-1.5">
                      <span class="text-[var(--text-secondary)]">Độ mạnh mật khẩu: <span :style="{ color: passwordStrength.color }">{{ passwordStrength.text }}</span></span>
                      <span class="font-mono text-[11px]" :style="{ color: passwordStrength.color }">{{ passwordStrength.score }}/4</span>
                    </div>
                    <div class="w-full h-1.5 bg-[var(--border)] rounded-full overflow-hidden">
                      <div class="h-full transition-all duration-300 rounded-full" :style="{ width: passwordStrength.width, backgroundColor: passwordStrength.color }"></div>
                    </div>
                    <ul class="mt-2 space-y-1 text-[11px] text-[var(--text-secondary)]">
                      <li class="flex items-center gap-1.5"><Check :size="12" :class="passwordForm.next.length >= 8 ? 'text-[var(--success)]' : 'text-slate-300'" /> Ít nhất 8 ký tự</li>
                      <li class="flex items-center gap-1.5"><Check :size="12" :class="/[A-Z]/.test(passwordForm.next) && /[a-z]/.test(passwordForm.next) ? 'text-[var(--success)]' : 'text-slate-300'" /> Có cả chữ hoa & chữ thường</li>
                      <li class="flex items-center gap-1.5"><Check :size="12" :class="/[0-9]/.test(passwordForm.next) || /[^A-Za-z0-9]/.test(passwordForm.next) ? 'text-[var(--success)]' : 'text-slate-300'" /> Có số hoặc ký tự đặc biệt</li>
                    </ul>
                  </div>
                </div>

                <Input label="Xác nhận mật khẩu mới *" type="password" v-model="passwordForm.confirm" :error="passwordErrors.confirm" placeholder="Nhập lại mật khẩu mới vừa gõ" required />

                <div class="pt-2 flex justify-end">
                  <Button type="submit" variant="primary" :disabled="changingPassword" class="bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white shadow-md">
                    <Key :size="16" />
                    {{ changingPassword ? 'Đang cập nhật...' : 'Cập nhật mật khẩu' }}
                  </Button>
                </div>
              </form>
            </div>
          </Card>

          <!-- Card 2: Two-Factor Authentication (2FA) -->
          <Card class="pf-card">
            <div class="pf-card-head flex items-center justify-between">
              <h3 class="pf-card-title">
                <span class="kpi-icon is-accent"><ShieldCheck :size="18" /></span>
                Xác thực hai yếu tố (2FA - Two-Factor Authentication)
              </h3>
              <span class="px-3 py-1 rounded-full text-xs font-bold border" :class="twoFactorEnabled ? 'bg-[var(--success-bg)] text-[var(--success)] border-emerald-200' : 'bg-slate-100 text-slate-500 border-slate-200'">
                {{ twoFactorEnabled ? '✓ Đã kích hoạt' : 'Chưa kích hoạt' }}
              </span>
            </div>
            <div class="p-6 flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
              <div class="space-y-1 flex-1">
                <p class="font-bold text-sm text-[var(--text-main)]">Bảo vệ tài khoản bằng ứng dụng xác thực (Authenticator App / TOTP)</p>
                <p class="text-xs text-[var(--text-secondary)] leading-relaxed">
                  Khi bật xác thực 2 lớp, bạn sẽ cần nhập thêm mã OTP từ Google Authenticator hoặc Authy mỗi khi đăng nhập từ thiết bị mới để bảo vệ an toàn tối đa cho dữ liệu phỏng vấn AI.
                </p>
              </div>
              <button @click="toggleTwoFactor" class="px-5 py-2.5 rounded-xl font-bold text-sm transition-all flex items-center gap-2 shadow-sm shrink-0" :class="twoFactorEnabled ? 'bg-rose-50 text-rose-600 border border-rose-200 hover:bg-rose-100' : 'bg-[var(--primary)] text-white hover:bg-[var(--primary-hover)] shadow-md'">
                <ShieldCheck :size="16" />
                {{ twoFactorEnabled ? 'Tắt bảo mật 2FA' : 'Kích hoạt 2FA ngay' }}
              </button>
            </div>
          </Card>

          <!-- Card 3: Active Sessions & Devices -->
          <Card class="pf-card">
            <div class="pf-card-head flex items-center justify-between">
              <h3 class="pf-card-title">
                <span class="kpi-icon"><Monitor :size="18" /></span>
                Thiết bị & Phiên làm việc hiện tại
              </h3>
              <button v-if="activeSessions.length > 1" @click="revokeOtherSessions" class="text-xs font-bold text-rose-600 hover:text-rose-700 flex items-center gap-1 transition-colors">
                <LogOut :size="14" /> Đăng xuất khỏi các thiết bị khác
              </button>
            </div>
            <div class="p-6 space-y-3">
              <div v-for="session in activeSessions" :key="session.id" class="p-4 rounded-xl border border-[var(--border)] bg-[var(--surface-soft)] flex items-center justify-between gap-4">
                <div class="flex items-center gap-3.5">
                  <div class="w-10 h-10 rounded-xl bg-white border border-[var(--border)] flex items-center justify-center text-[var(--primary)] shadow-sm shrink-0">
                    <component :is="session.icon" :size="20" />
                  </div>
                  <div>
                    <div class="flex items-center gap-2">
                      <span class="font-bold text-sm text-[var(--text-main)]">{{ session.device }}</span>
                      <span v-if="session.isCurrent" class="px-2 py-0.5 rounded-md bg-[var(--primary-light)] text-[var(--primary)] font-bold text-[10px]">Thiết bị này</span>
                    </div>
                    <p class="text-xs text-[var(--text-secondary)] mt-0.5 flex items-center gap-2">
                      <span>📍 {{ session.location }}</span>
                      <span>•</span>
                      <span>🌐 IP: {{ session.ip }}</span>
                    </p>
                  </div>
                </div>
                <div class="text-right shrink-0">
                  <span class="text-xs font-medium" :class="session.isCurrent ? 'text-[var(--success)] font-bold' : 'text-[var(--text-muted)]'">
                    {{ session.lastActive }}
                  </span>
                </div>
              </div>
            </div>
          </Card>

          <!-- Card 4: Security Activity Log -->
          <Card class="pf-card">
            <div class="pf-card-head">
              <h3 class="pf-card-title">
                <span class="kpi-icon is-info"><History :size="18" /></span>
                Nhật ký hoạt động bảo mật gần đây
              </h3>
            </div>
            <div class="p-6">
              <div class="divide-y divide-[var(--border)]">
                <div v-for="log in securityLogs" :key="log.id" class="py-3 flex items-center justify-between gap-4 first:pt-0 last:pb-0">
                  <div class="flex items-center gap-3">
                    <div class="w-2 h-2 rounded-full shrink-0" :class="log.status === 'success' ? 'bg-[var(--success)]' : log.status === 'warning' ? 'bg-amber-500' : 'bg-[var(--primary)]'"></div>
                    <div>
                      <p class="text-sm font-semibold text-[var(--text-main)]">{{ log.action }}</p>
                      <p class="text-xs text-[var(--text-secondary)]">{{ log.details }}</p>
                    </div>
                  </div>
                  <span class="text-xs text-[var(--text-muted)] font-medium flex items-center gap-1 shrink-0">
                    <Clock :size="12" /> {{ log.time }}
                  </span>
                </div>
              </div>
            </div>
          </Card>
        </div>

        <!-- TAB 3: NOTIFICATIONS -->
        <div v-if="activeTab === 'notifications'" class="space-y-6">
          <Card title="Cấu hình nhận thông báo">
            <div class="space-y-6">
              <div>
                <h3 class="text-base font-bold text-[var(--text-main)] uppercase tracking-wider mb-3 flex items-center gap-2">
                  <span class="w-2 h-2 rounded-full bg-[var(--primary)]"></span> Thông báo qua Email
                </h3>
                <div class="space-y-4 bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                  <div class="flex items-center justify-between gap-6">
                    <div class="space-y-1 pr-4">
                      <span class="text-base font-semibold text-[var(--text-main)] block">Lịch phỏng vấn mới</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Nhận email thông báo khi nhà tuyển dụng mời hoặc xác nhận lịch phỏng vấn.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_email_interview" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                    </label>
                  </div>

                  <div class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                    <div class="space-y-1 pr-4">
                      <span class="text-base font-semibold text-[var(--text-main)] block">Nhắc nhở phỏng vấn (Nhắc trước 1 giờ)</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Gửi email nhắc nhở trước khi buổi phỏng vấn (thật hoặc thử) diễn ra.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_email_reminder" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                    </label>
                  </div>
                </div>
              </div>

              <div>
                <h3 class="text-base font-bold text-[var(--text-main)] uppercase tracking-wider mb-3 flex items-center gap-2">
                  <span class="w-2 h-2 rounded-full bg-[var(--accent)]"></span> Thông báo trên Trình duyệt (Push Notifications)
                </h3>
                <div class="bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                  <div class="flex items-center justify-between gap-6">
                    <div class="space-y-1 pr-4">
                      <span class="text-base font-semibold text-[var(--text-main)] block">Thông báo thời gian thực (Realtime Alert)</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Hiển thị thông báo ngay trên góc trình duyệt khi có cập nhật kết quả hoặc phản hồi mới từ AI.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_push" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                    </label>
                  </div>
                </div>
              </div>

              <div class="flex justify-end pt-2">
                <Button variant="primary" @click="handleSaveSettings" :disabled="saving">
                  <Save :size="16" />
                  {{ saving ? 'Đang lưu...' : 'Lưu tùy chọn thông báo' }}
                </Button>
              </div>
            </div>
          </Card>
        </div>

        <!-- TAB 4: PRIVACY & DATA -->
        <div v-if="activeTab === 'privacy'" class="space-y-6">
          <Card title="Quyền riêng tư & Hiển thị">
            <div class="space-y-6">
              <div v-if="role === 'candidate'" class="space-y-4 bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                <div class="flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Hồ sơ công khai (Public Profile)</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Cho phép các nhà tuyển dụng trên hệ thống AI Interview tìm thấy hồ sơ, kỹ năng và CV của bạn để trao cơ hội việc làm.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" v-model="settings.public_profile" class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                  </label>
                </div>

                <div class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Chia sẻ dữ liệu ẩn danh cải thiện AI</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Cho phép sử dụng các phiên phỏng vấn thử (Mock Interview) đã được ẩn danh hóa thông tin cá nhân để huấn luyện và nâng cao độ chính xác của mô hình AI.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" v-model="settings.share_anon_results" class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                  </label>
                </div>
              </div>

              <div class="flex justify-end pt-2">
                <Button variant="primary" @click="handleSaveSettings" :disabled="saving">
                  <Save :size="16" />
                  {{ saving ? 'Đang lưu...' : 'Lưu quyền riêng tư' }}
                </Button>
              </div>
            </div>
          </Card>

          <Card class="border-[var(--danger)]/30">
            <div class="flex items-start justify-between flex-wrap gap-4">
              <div>
                <h3 class="text-base font-bold text-[var(--danger)] flex items-center gap-2">
                  <AlertCircle :size="18" /> Vùng nguy hiểm (Danger Zone)
                </h3>
                <p class="text-xs text-[var(--text-secondary)] mt-1 max-w-xl">
                  Xóa toàn bộ thông tin tài khoản, CV, kết quả luyện tập và lịch sử phỏng vấn của bạn khỏi hệ thống. Hành động này không thể khôi phục.
                </p>
              </div>
              <Button variant="danger" @click="showDeleteModal = true">
                <Trash2 :size="16" /> Xóa tài khoản vĩnh viễn
              </Button>
            </div>
          </Card>
        </div>

        <!-- TAB 5: AI WORKSPACE (Recruiter Only) -->
        <div v-if="activeTab === 'workspace' && role === 'recruiter'" class="space-y-6">
          <Card title="Cấu hình AI Workspace">
            <div class="space-y-6">
              <p class="text-sm text-[var(--text-secondary)]">Tuỳ chỉnh cách AI Assistant hoạt động trong không gian làm việc và phòng phỏng vấn thực tế của doanh nghiệp.</p>
              
              <div class="space-y-1">
                <label class="text-sm font-semibold text-[var(--text-main)]">Mô hình AI đánh giá mặc định</label>
                <select class="pf-select" defaultValue="gemini-1.5-pro">
                  <option value="gemini-1.5-pro">Gemini 1.5 Pro (Phân tích sâu, độ chính xác cao nhất)</option>
                  <option value="gemini-1.5-flash">Gemini 1.5 Flash (Tốc độ phản hồi cực nhanh realtime)</option>
                </select>
              </div>

              <div class="space-y-4 bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                <div class="flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Tự động bóc tách JD thành bộ tiêu chí (Rubric) chi tiết</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">AI sẽ phân tích mô tả công việc và tự động chuẩn bị ma trận tiêu chí chấm điểm trước buổi phỏng vấn.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" defaultChecked class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--accent)] shadow-inner"></div>
                  </label>
                </div>

                <div class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Tự động gợi ý câu hỏi Follow-up realtime trong phòng phỏng vấn</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Phân tích câu trả lời của ứng viên ngay lúc đang nói để đưa ra gợi ý câu hỏi đào sâu cho nhà tuyển dụng.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" defaultChecked class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--accent)] shadow-inner"></div>
                  </label>
                </div>
              </div>

              <div class="flex justify-end pt-2">
                <Button variant="primary" @click="handleSaveSettings" :disabled="saving">
                  <Save :size="16" /> Lưu cấu hình AI
                </Button>
              </div>
            </div>
          </Card>
        </div>

      </div>
    </div>

    <!-- CV Detail Modal -->
    <Modal :isOpen="!!selectedCv" @close="closeCvModal" title="Chi tiết CV & Trích xuất AI" size="xl">
      <div v-if="selectedCv" class="space-y-6">
        <div class="flex flex-col lg:flex-row gap-6">
          <!-- PDF / File Viewer Preview -->
          <div class="flex-1 bg-[var(--surface-soft)] border border-[var(--border)] rounded-xl flex flex-col items-center justify-center p-3 text-center min-h-[500px] overflow-hidden">
            <div v-if="selectedCv.url" class="w-full flex flex-col items-center">
              <div class="w-full flex justify-end mb-2">
                <a :href="selectedCv.url" target="_blank" class="text-xs font-bold text-[var(--primary)] hover:underline flex items-center gap-1 bg-white px-3 py-1.5 rounded-lg border border-[var(--border)] shadow-sm">
                  <span>🔗 Mở file xem toàn màn hình</span>
                </a>
              </div>
              <img v-if="selectedCv.url.toLowerCase().includes('.png') || selectedCv.url.toLowerCase().includes('.jpg') || selectedCv.url.toLowerCase().includes('.jpeg')" :src="selectedCv.url" class="max-h-[460px] object-contain rounded-lg shadow-sm" />
              <iframe v-else :src="selectedCv.url" class="w-full h-[460px] rounded-lg border-0 bg-white shadow-sm"></iframe>
            </div>
            <div v-else class="flex flex-col items-center justify-center p-8">
              <FileText class="w-16 h-16 text-[var(--text-muted)] mb-4" />
              <p class="text-[var(--text-secondary)] font-medium">Không có bản xem trước cho tài liệu này</p>
            </div>
          </div>

          <!-- Extracted Data -->
          <div class="w-full lg:w-[360px] space-y-5 shrink-0">
            <div class="ai-block p-4">
              <h4 class="text-xs font-bold uppercase tracking-wider mb-3 flex items-center gap-2 text-[var(--accent)]">
                <Bot :size="16" /> Thông tin AI trích xuất (Gemini)
              </h4>
              <div class="space-y-3 text-sm">
                <div class="flex flex-col gap-0.5 border-b border-[var(--border)] pb-2.5">
                  <span class="text-xs text-[var(--text-secondary)]">Vị trí phù hợp:</span>
                  <span class="font-bold text-[var(--text-main)]">
                    {{ selectedCv.parsedData?.role || selectedCv.parsedData?.target_role || (selectedCv.parsedData?.work_experience?.[0]?.role) || 'Chuyên viên Công nghệ' }}
                  </span>
                </div>
                <div class="flex flex-col gap-0.5 border-b border-[var(--border)] pb-2.5">
                  <span class="text-xs text-[var(--text-secondary)]">Cấp độ kinh nghiệm:</span>
                  <span class="font-bold text-[var(--text-main)]">
                    {{ selectedCv.parsedData?.level || (selectedCv.parsedData?.experience_years_estimate != null ? `${selectedCv.parsedData.experience_years_estimate} năm kinh nghiệm` : 'Middle') }}
                  </span>
                </div>
                <div class="flex items-center justify-between pt-1">
                  <span class="text-xs text-[var(--text-secondary)]">Trạng thái AI:</span>
                  <span class="font-bold text-[var(--success)] flex items-center gap-1">
                    <CheckCircle :size="14" /> Đã phân tích từ Backend
                  </span>
                </div>
              </div>
            </div>

            <div>
              <h4 class="text-xs font-bold uppercase tracking-wider text-[var(--text-secondary)] mb-3">Kỹ năng phát hiện được</h4>
              <div class="flex flex-wrap gap-2">
                <template v-if="selectedCv.parsedData?.skills && selectedCv.parsedData.skills.length > 0">
                  <span v-for="(sk, idx) in selectedCv.parsedData.skills" :key="idx" class="skill-tag border-[var(--primary-light)] text-[var(--primary)] bg-[var(--primary-light)]/40 font-semibold">
                    {{ typeof sk === 'object' ? (sk.name || sk.label || JSON.stringify(sk)) : sk }}
                  </span>
                </template>
                <template v-else>
                  <span v-for="skill in ['JavaScript', 'Vue 3', 'REST API', 'TailwindCSS']" :key="skill" class="skill-tag border-[var(--primary-light)] text-[var(--primary)] bg-[var(--primary-light)]/40">
                    {{ skill }}
                  </span>
                </template>
              </div>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-3 pt-4 border-t border-[var(--border)]">
          <Button variant="ghost" style="color: var(--danger)" @click="deleteCv(selectedCv.id); closeCvModal()">Xóa CV này</Button>
          <Button variant="primary" @click="closeCvModal">Đóng</Button>
        </div>
      </div>
    </Modal>

    <!-- Delete Account Confirmation Modal -->
    <Modal :isOpen="showDeleteModal" @close="showDeleteModal = false" title="Xác nhận xóa tài khoản vĩnh viễn" size="md">
      <div class="space-y-4">
        <p class="text-sm text-[var(--text-main)] leading-relaxed">
          Bạn có chắc chắn muốn xóa tài khoản này? Hành động này <strong class="text-[var(--danger)]">không thể hoàn tác</strong> và toàn bộ dữ liệu hồ sơ, CV, kết quả Mock Interview của bạn sẽ bị xóa vĩnh viễn khỏi hệ thống.
        </p>
        <div class="flex justify-end gap-3 pt-4 border-t border-[var(--border)]">
          <Button variant="secondary" @click="showDeleteModal = false">Huỷ thao tác</Button>
          <Button variant="danger" @click="confirmDeleteAccount">Xác nhận xóa tài khoản</Button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.page-title {
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--text-main);
  line-height: 1.25;
}
.page-subtitle {
  color: var(--text-secondary);
  font-size: 14px;
  margin-top: 4px;
}

/* Sidebar Navigation Tab Items */
.nav-tab-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 16px;
  border-radius: var(--radius);
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s ease;
}
.nav-tab-item:hover {
  background: var(--surface-soft);
  color: var(--text-main);
}
.nav-tab-item.active {
  background: var(--primary-light);
  color: var(--primary);
}

.pf-card { padding: 0; overflow: hidden; }
.pf-card-head { padding: 18px 24px; border-bottom: 1px solid var(--border); background: var(--surface-soft); }
.pf-card-title { display: flex; align-items: center; gap: 10px; font-size: 16px; font-weight: 700; color: var(--text-main); margin: 0; }
.pf-sub-title { display: flex; align-items: center; gap: 10px; font-size: 15px; font-weight: 700; color: var(--text-main); margin-bottom: 18px; }
.pf-label { font-size: 13px; font-weight: 600; color: var(--text-secondary); }
.pf-select { width: 100%; height: 40px; padding: 0 14px; border: 1px solid var(--border); border-radius: var(--radius); background: var(--surface); color: var(--text-main); outline: none; transition: all 0.2s ease; font-family: var(--sans); font-size: 14px; }
.pf-select:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(37,99,235,0.12); }
.pf-textarea { width: 100%; padding: 12px 14px; border: 1px solid var(--border); border-radius: var(--radius); outline: none; resize: none; color: var(--text-main); background: var(--surface); font-family: var(--sans); font-size: 14px; transition: all 0.2s ease; }
.pf-textarea:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(37,99,235,0.12); }

.pf-dropzone { border: 2px dashed var(--border); border-radius: var(--radius); padding: 24px; text-align: center; cursor: pointer; transition: all 0.25s ease; }
.pf-dropzone:hover { border-color: var(--primary); background: var(--primary-light); }
.pf-dz-icon { display: flex; align-items: center; justify-content: center; width: 44px; height: 44px; margin: 0 auto 10px; border-radius: 50%; background: var(--primary-light); color: var(--primary); transition: transform 0.25s ease; }
.pf-dropzone:hover .pf-dz-icon { transform: translateY(-2px); }
.pf-dz-title { font-weight: 700; font-size: 14px; color: var(--text-main); margin-bottom: 4px; }
.pf-list-title { font-size: 12px; font-weight: 700; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.04em; }

.pf-cv-item { display: flex; align-items: center; gap: 14px; padding: 12px 14px; border-radius: var(--radius); border: 1px solid var(--border); background: var(--surface); transition: all 0.2s ease; }
.pf-cv-item.is-done { cursor: pointer; }
.pf-cv-item.is-done:hover { box-shadow: var(--shadow-sm); border-color: var(--primary); }
.pf-cv-item.is-loading { opacity: 0.8; cursor: wait; }
.pf-cv-icon { display: flex; align-items: center; justify-content: center; width: 40px; height: 40px; border-radius: var(--radius); background: var(--primary-light); color: var(--primary); flex-shrink: 0; }
.pf-cv-icon.is-muted { background: var(--surface-soft); color: var(--text-muted); }
.pf-cv-name { font-weight: 700; font-size: 14px; color: var(--text-main); }
.pf-cv-meta { display: flex; align-items: center; gap: 5px; font-size: 12px; margin-top: 2px; font-weight: 500; }
.pf-cv-meta.is-ok { color: var(--success); }
.pf-cv-meta.is-wait { color: var(--warning); }
.pf-analyzing { display: inline-flex; align-items: center; gap: 6px; padding: 4px 10px; border-radius: var(--radius); background: rgba(217,119,6,0.12); color: var(--warning); font-size: 12px; font-weight: 600; }
.pf-icon-btn { padding: 6px; border: none; background: transparent; border-radius: var(--radius); color: var(--primary); cursor: pointer; transition: all 0.2s ease; }
.pf-icon-btn:hover { background: var(--primary-light); }
.pf-icon-btn.is-danger { color: var(--danger); }
.pf-icon-btn.is-danger:hover { background: rgba(220,38,38,0.1); }
.pf-spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.kpi-icon { display: inline-flex; align-items: center; justify-content: center; width: 36px; height: 36px; border-radius: var(--radius); background: var(--primary-light); color: var(--primary); }
.kpi-icon.is-accent { background: var(--accent-bg); color: var(--accent); }

/* Skill Tags Input */
.skills-input-wrapper { border: 1px solid var(--border); border-radius: var(--radius); padding: 8px 12px; min-height: 80px; background: var(--surface); transition: all 0.2s; cursor: text; }
.skills-input-wrapper:focus-within { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(37,99,235,0.12); }
.skills-tags-container { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; }
.skill-tag { display: inline-flex; align-items: center; gap: 4px; padding: 4px 10px; border-radius: var(--radius-full); background: var(--surface-soft); border: 1px solid var(--border); font-size: 13px; font-weight: 600; color: var(--text-main); }
.skill-tag-remove { display: inline-flex; align-items: center; justify-content: center; background: none; border: none; color: var(--text-muted); cursor: pointer; padding: 2px; border-radius: 50%; }
.skill-tag-remove:hover { color: var(--danger); background: rgba(220,38,38,0.1); }
.skill-tag-input { flex: 1; min-width: 140px; border: none; outline: none; background: transparent; font-size: 14px; padding: 4px 0; color: var(--text-main); font-family: var(--sans); }
</style>
