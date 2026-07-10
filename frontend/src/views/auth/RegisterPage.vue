<script setup>
import { ref, reactive, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { ShieldCheck, ArrowRight, ArrowLeft, Mail, Lock, User, CheckCircle2, RefreshCw, Sparkles, Building2, UserCheck, AlertCircle } from 'lucide-vue-next'

const router = useRouter()

// Registration Step (1: Form, 2: OTP Verification)
const step = ref(1)

// Step 1 Form Data
const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: 'candidate'
})

const error = ref('')
const loading = ref(false)

// Google Auth Simulation State
const showGoogleModal = ref(false)

// Step 2 OTP Data (6 digits)
const otpDigits = ref(['', '', '', '', '', ''])
const otpInputs = ref([])
const timer = ref(59)
let timerInterval = null

const startOtpTimer = () => {
  clearInterval(timerInterval)
  timer.value = 59
  timerInterval = setInterval(() => {
    if (timer.value > 0) {
      timer.value--
    } else {
      clearInterval(timerInterval)
    }
  }, 1000)
}

const handleStep1Submit = async (e) => {
  e.preventDefault()
  if (!form.name || !form.email || !form.password || !form.confirmPassword) {
    error.value = 'Vui lòng nhập đầy đủ thông tin cá nhân và mật khẩu!'
    return
  }
  
  if (form.password.length < 6) {
    error.value = 'Mật khẩu phải có ít nhất 6 ký tự!'
    return
  }

  if (form.password !== form.confirmPassword) {
    error.value = 'Mật khẩu xác nhận không trùng khớp!'
    return
  }
  
  error.value = ''
  loading.value = true

  // Simulate sending real OTP verification email via backend mailer
  setTimeout(() => {
    loading.value = false
    step.value = 2
    startOtpTimer()
    nextTick(() => {
      if (otpInputs.value[0]) otpInputs.value[0].focus()
    })
  }, 700)
}

// Handle OTP digit inputs
const handleOtpInput = (index, e) => {
  const value = e.target.value
  if (!/^\d*$/.test(value)) {
    otpDigits.value[index] = ''
    return
  }
  
  if (value.length > 1) {
    // Paste handled separately, but just in case takes last char
    otpDigits.value[index] = value.slice(-1)
  } else {
    otpDigits.value[index] = value
  }

  // Move to next digit if typed
  if (value !== '' && index < 5) {
    nextTick(() => {
      if (otpInputs.value[index + 1]) otpInputs.value[index + 1].focus()
    })
  }

  // Auto verify if all 6 digits entered
  if (otpDigits.value.every(d => d !== '') && index === 5) {
    handleVerifyOtp()
  }
}

const handleOtpKeyDown = (index, e) => {
  if (e.key === 'Backspace' && !otpDigits.value[index] && index > 0) {
    nextTick(() => {
      if (otpInputs.value[index - 1]) {
        otpInputs.value[index - 1].focus()
        otpDigits.value[index - 1] = ''
      }
    })
  }
}

const handleOtpPaste = (e) => {
  e.preventDefault()
  const pastedData = e.clipboardData.getData('text').replace(/\D/g, '').slice(0, 6)
  if (pastedData) {
    for (let i = 0; i < 6; i++) {
      otpDigits.value[i] = pastedData[i] || ''
    }
    nextTick(() => {
      const focusIndex = Math.min(pastedData.length, 5)
      if (otpInputs.value[focusIndex]) otpInputs.value[focusIndex].focus()
      if (pastedData.length === 6) {
        handleVerifyOtp()
      }
    })
  }
}

const resendOtp = () => {
  if (timer.value > 0) return
  startOtpTimer()
  otpDigits.value = ['', '', '', '', '', '']
  nextTick(() => {
    if (otpInputs.value[0]) otpInputs.value[0].focus()
  })
}

const handleVerifyOtp = async () => {
  const code = otpDigits.value.join('')
  if (code.length < 6) {
    error.value = 'Vui lòng nhập đủ 6 chữ số mã xác nhận OTP!'
    return
  }

  error.value = ''
  loading.value = true

  try {
    // 1. Register user
    await authStore.register(form.email, form.password, form.name, form.role)
    
    // 2. Auto login right after successful OTP verification
    await authStore.login(form.email, form.password)
    
    if (form.role === 'recruiter') {
      router.push({ path: '/dashboard', state: { message: `Xác thực OTP thành công! Chào mừng Nhà tuyển dụng ${form.name}. Vui lòng chuẩn bị giấy phép kinh doanh để xác minh doanh nghiệp.` } })
    } else {
      router.push({ path: '/', state: { message: `Xác thực OTP thành công! Chào mừng ${form.name} đến với WeMake AI.` } })
    }
  } catch (err) {
    if (err.message === 'email already exists') {
      error.value = 'Email này đã được đăng ký trong hệ thống. Vui lòng quay lại và dùng email khác hoặc đăng nhập.'
    } else {
      error.value = err.message || 'Mã xác nhận OTP hoặc đăng ký không thành công, vui lòng thử lại.'
    }
  } finally {
    loading.value = false
  }
}

const handleGoogleSelect = async (roleType, emailChoice) => {
  showGoogleModal.value = false
  loading.value = true
  error.value = ''
  try {
    const user = await authStore.loginWithGoogle(roleType, emailChoice)
    if (user.role === 'recruiter') {
      router.push({ path: '/dashboard', state: { message: `Đăng nhập Google thành công! Chào mừng Nhà tuyển dụng ${user.full_name}.` } })
    } else {
      router.push({ path: '/', state: { message: `Đăng nhập Google thành công! Chào mừng ${user.full_name}.` } })
    }
  } catch (err) {
    error.value = 'Đăng nhập Google không thành công: ' + (err.message || 'Lỗi kết nối OAuth')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-900 px-4 py-12 relative overflow-hidden">
    <!-- Ambient Background Glows -->
    <div class="absolute top-1/4 -left-32 w-96 h-96 bg-indigo-500/20 rounded-full blur-3xl pointer-events-none animate-pulse"></div>
    <div class="absolute bottom-1/4 -right-32 w-96 h-96 bg-purple-500/20 rounded-full blur-3xl pointer-events-none animate-pulse" style="animation-delay: 2s;"></div>

    <!-- Google OAuth Selection Modal -->
    <div v-if="showGoogleModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4 animate-fade-in">
      <div class="bg-white dark:bg-slate-800 rounded-2xl max-w-md w-full p-6 shadow-2xl border border-slate-200 dark:border-slate-700 space-y-6">
        <div class="flex items-center gap-3 pb-4 border-b border-slate-100 dark:border-slate-700">
          <svg class="w-8 h-8" viewBox="0 0 24 24">
            <path fill="#EA4335" d="M12 5c1.6 0 3 .6 4.1 1.7l3.1-3.1C17.3 1.8 14.8 1 12 1 7.4 1 3.5 3.6 1.6 7.4l3.7 2.8C6.2 7.3 8.9 5 12 5z"/>
            <path fill="#4285F4" d="M23.5 12.3c0-.8-.1-1.7-.2-2.3H12v4.6h6.5c-.3 1.5-1.1 2.8-2.4 3.7l3.7 2.9c2.2-2 3.7-5 3.7-8.9z"/>
            <path fill="#FBBC05" d="M5.3 14.8c-.2-.7-.4-1.5-.4-2.3s.2-1.5.4-2.3L1.6 7.4C.6 9.4 0 11.6 0 14s.6 4.6 1.6 6.6l3.7-2.8z"/>
            <path fill="#34A853" d="M12 23c3.2 0 6-1.1 8-3l-3.7-2.9c-1.1.7-2.5 1.2-4.3 1.2-3.1 0-5.8-2.3-6.7-5.2L1.6 15.9C3.5 19.7 7.4 23 12 23z"/>
          </svg>
          <div>
            <h3 class="font-bold text-slate-800 dark:text-white text-lg">Đăng ký bằng tài khoản Google</h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">Chọn tài khoản Google để tiếp tục với WeMake AI</p>
          </div>
        </div>

        <div class="space-y-3">
          <button @click="handleGoogleSelect('candidate', 'candidate@wemake.vn')" 
                  class="w-full flex items-center justify-between p-3.5 rounded-xl border border-slate-200 dark:border-slate-700 hover:border-indigo-500 hover:bg-indigo-50/50 dark:hover:bg-slate-700/60 transition-all text-left group">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-blue-600 text-white font-bold flex items-center justify-center text-sm shadow-sm">
                N
              </div>
              <div>
                <div class="font-bold text-slate-800 dark:text-white text-sm group-hover:text-indigo-600 dark:group-hover:text-indigo-400">Nguyễn Văn A (Ứng viên)</div>
                <div class="text-xs text-slate-500 dark:text-slate-400">candidate@wemake.vn</div>
              </div>
            </div>
            <span class="text-xs font-semibold px-2.5 py-1 bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300 rounded-full">Candidate</span>
          </button>

          <button @click="handleGoogleSelect('recruiter', 'recruiter@wemake.vn')" 
                  class="w-full flex items-center justify-between p-3.5 rounded-xl border border-slate-200 dark:border-slate-700 hover:border-purple-500 hover:bg-purple-50/50 dark:hover:bg-slate-700/60 transition-all text-left group">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-purple-600 text-white font-bold flex items-center justify-center text-sm shadow-sm">
                H
              </div>
              <div>
                <div class="font-bold text-slate-800 dark:text-white text-sm group-hover:text-purple-600 dark:group-hover:text-purple-400">HR Manager (Nhà tuyển dụng)</div>
                <div class="text-xs text-slate-500 dark:text-slate-400">recruiter@wemake.vn</div>
              </div>
            </div>
            <span class="text-xs font-semibold px-2.5 py-1 bg-purple-100 dark:bg-purple-900/40 text-purple-700 dark:text-purple-300 rounded-full">Recruiter</span>
          </button>
        </div>

        <div class="pt-2 flex justify-end">
          <button @click="showGoogleModal = false" class="px-4 py-2 text-sm font-semibold text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-white transition-colors">
            Hủy bỏ
          </button>
        </div>
      </div>
    </div>

    <!-- Main Registration Box -->
    <div class="w-full max-w-lg bg-white/95 dark:bg-slate-800/95 backdrop-blur-2xl rounded-3xl shadow-2xl border border-white/20 dark:border-slate-700/50 p-8 sm:p-10 relative z-10 transition-all duration-500">
      
      <!-- Brand Logo Header -->
      <div class="text-center space-y-2 mb-8">
        <div @click="router.push('/')" class="inline-block cursor-pointer transform hover:scale-105 transition-transform">
          <img src="/images/logo.png" alt="WeMake AI Logo" class="h-16 mx-auto object-contain drop-shadow-md" />
        </div>
        <h1 class="text-2xl sm:text-3xl font-extrabold text-slate-900 dark:text-white tracking-tight">
          {{ step === 1 ? 'Khởi tạo tài khoản AI' : 'Xác minh bảo mật OTP' }}
        </h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm">
          {{ step === 1 ? 'Nền tảng phỏng vấn & đánh giá năng lực AI hàng đầu' : `Mã 6 chữ số đã được gửi tới ${form.email}` }}
        </p>
      </div>

      <!-- Error Alert -->
      <div v-if="error" class="mb-6 p-4 rounded-2xl bg-rose-50 dark:bg-rose-950/50 border border-rose-200 dark:border-rose-800/60 flex items-start gap-3 text-rose-600 dark:text-rose-400 text-sm animate-shake">
        <AlertCircle size="20" class="shrink-0 mt-0.5" />
        <span class="font-medium flex-1">{{ error }}</span>
      </div>

      <!-- ================= STEP 1: FORM INFORMATION ================= -->
      <div v-if="step === 1" class="space-y-6 animate-fade-in">
        
        <form @submit="handleStep1Submit" class="space-y-4">
          <!-- Role Selector Tabs -->
          <div>
            <label class="block text-xs font-bold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-2">Bạn là ai?</label>
            <div class="grid grid-cols-2 gap-3">
              <button type="button"
                      @click="form.role = 'candidate'"
                      :class="form.role === 'candidate' ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30 border-indigo-600 font-bold' : 'bg-slate-100 dark:bg-slate-900/60 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-700 font-medium hover:bg-slate-200 dark:hover:bg-slate-800'"
                      class="flex items-center justify-center gap-2.5 py-3 px-4 rounded-xl border text-sm transition-all">
                <UserCheck size="18" />
                <span>Ứng viên</span>
              </button>

              <button type="button"
                      @click="form.role = 'recruiter'"
                      :class="form.role === 'recruiter' ? 'bg-purple-600 text-white shadow-md shadow-purple-600/30 border-purple-600 font-bold' : 'bg-slate-100 dark:bg-slate-900/60 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-700 font-medium hover:bg-slate-200 dark:hover:bg-slate-800'"
                      class="flex items-center justify-center gap-2.5 py-3 px-4 rounded-xl border text-sm transition-all">
                <Building2 size="18" />
                <span>Nhà tuyển dụng</span>
              </button>
            </div>
            <p v-if="form.role === 'recruiter'" class="text-[11px] text-purple-600 dark:text-purple-400 font-medium mt-1.5 flex items-center gap-1">
              <ShieldCheck size="14" /> Nhà tuyển dụng cần xác thực doanh nghiệp sau bước đăng ký.
            </p>
          </div>

          <!-- Name Input -->
          <div>
            <label class="block text-xs font-bold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1.5">Họ và tên</label>
            <div class="relative">
              <User size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
              <input v-model="form.name"
                     type="text"
                     required
                     placeholder="Nguyễn Văn A"
                     class="w-full pl-11 pr-4 py-3 bg-slate-50 dark:bg-slate-900/80 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 font-medium text-sm transition-all" />
            </div>
          </div>

          <!-- Email Input -->
          <div>
            <label class="block text-xs font-bold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1.5">Địa chỉ Email</label>
            <div class="relative">
              <Mail size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
              <input v-model="form.email"
                     type="email"
                     required
                     placeholder="nhapemail@congty.com"
                     class="w-full pl-11 pr-4 py-3 bg-slate-50 dark:bg-slate-900/80 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 font-medium text-sm transition-all" />
            </div>
          </div>

          <!-- Password Grid -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1.5">Mật khẩu</label>
              <div class="relative">
                <Lock size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
                <input v-model="form.password"
                       type="password"
                       required
                       placeholder="••••••••"
                       class="w-full pl-11 pr-4 py-3 bg-slate-50 dark:bg-slate-900/80 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 font-medium text-sm transition-all" />
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1.5">Xác nhận mật khẩu</label>
              <div class="relative">
                <Lock size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
                <input v-model="form.confirmPassword"
                       type="password"
                       required
                       placeholder="••••••••"
                       class="w-full pl-11 pr-4 py-3 bg-slate-50 dark:bg-slate-900/80 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 font-medium text-sm transition-all" />
              </div>
            </div>
          </div>

          <!-- Submit Button -->
          <button type="submit"
                  :disabled="loading"
                  class="w-full mt-2 py-3.5 px-6 rounded-2xl bg-gradient-to-r from-indigo-600 via-blue-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-bold text-sm shadow-lg shadow-indigo-600/30 hover:shadow-indigo-600/50 transition-all flex items-center justify-center gap-2.5 disabled:opacity-50 transform hover:-translate-y-0.5">
            <span v-if="loading" class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <template v-else>
              <span>Tiếp tục bước xác minh OTP</span>
              <ArrowRight size="18" />
            </template>
          </button>
        </form>

        <div class="relative flex py-1 items-center">
          <div class="flex-grow border-t border-slate-200 dark:border-slate-700"></div>
          <span class="flex-shrink mx-4 text-xs font-semibold text-slate-400 uppercase tracking-wider">hoặc đăng nhập nhanh qua google</span>
          <div class="flex-grow border-t border-slate-200 dark:border-slate-700"></div>
        </div>

        <!-- Google Sign Up Button -->
        <button @click="showGoogleModal = true"
                type="button"
                class="w-full py-3 px-4 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 hover:border-indigo-500/50 dark:hover:border-indigo-500/50 hover:bg-slate-50 dark:hover:bg-slate-800/80 text-slate-700 dark:text-slate-200 font-semibold text-sm shadow-sm transition-all flex items-center justify-center gap-3 group">
          <svg class="w-5 h-5 transition-transform group-hover:scale-110" viewBox="0 0 24 24">
            <path fill="#EA4335" d="M12 5c1.6 0 3 .6 4.1 1.7l3.1-3.1C17.3 1.8 14.8 1 12 1 7.4 1 3.5 3.6 1.6 7.4l3.7 2.8C6.2 7.3 8.9 5 12 5z"/>
            <path fill="#4285F4" d="M23.5 12.3c0-.8-.1-1.7-.2-2.3H12v4.6h6.5c-.3 1.5-1.1 2.8-2.4 3.7l3.7 2.9c2.2-2 3.7-5 3.7-8.9z"/>
            <path fill="#FBBC05" d="M5.3 14.8c-.2-.7-.4-1.5-.4-2.3s.2-1.5.4-2.3L1.6 7.4C.6 9.4 0 11.6 0 14s.6 4.6 1.6 6.6l3.7-2.8z"/>
            <path fill="#34A853" d="M12 23c3.2 0 6-1.1 8-3l-3.7-2.9c-1.1.7-2.5 1.2-4.3 1.2-3.1 0-5.8-2.3-6.7-5.2L1.6 15.9C3.5 19.7 7.4 23 12 23z"/>
          </svg>
          <span>Tiếp tục với tài khoản Google</span>
        </button>

        <!-- Footer Link -->
        <div class="pt-2 text-center border-t border-slate-100 dark:border-slate-700/60">
          <p class="text-sm text-slate-500 dark:text-slate-400">
            Đã có tài khoản WeMake AI?
            <span @click="router.push('/login')" class="text-indigo-600 dark:text-indigo-400 font-bold hover:underline cursor-pointer ml-1">
              Đăng nhập ngay
            </span>
          </p>
        </div>
      </div>

      <!-- ================= STEP 2: SLEEK 6-DIGIT OTP SCREEN ================= -->
      <div v-else-if="step === 2" class="space-y-6 animate-fade-in">
        
        <div class="bg-indigo-50 dark:bg-indigo-950/40 p-4 rounded-2xl border border-indigo-200/60 dark:border-indigo-800/40 flex items-start gap-3">
          <ShieldCheck size="24" class="text-indigo-600 dark:text-indigo-400 shrink-0 mt-0.5" />
          <div class="text-xs text-indigo-950 dark:text-indigo-200 leading-relaxed">
            Hệ thống bảo mật đã gửi mã OTP 6 chữ số tới hộp thư <strong>{{ form.email }}</strong>. Vui lòng nhập mã để kích hoạt tài khoản và mở khóa đặc quyền AI.
          </div>
        </div>

        <!-- 6-Digit OTP Inputs -->
        <div class="space-y-3">
          <label class="block text-center text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-widest">
            Nhập mã xác thực 6 chữ số
          </label>
          <div class="flex items-center justify-center gap-2 sm:gap-3" @paste="handleOtpPaste">
            <input v-for="(digit, index) in otpDigits"
                   :key="index"
                   :ref="el => otpInputs[index] = el"
                   v-model="otpDigits[index]"
                   type="text"
                   maxlength="1"
                   @input="handleOtpInput(index, $event)"
                   @keydown="handleOtpKeyDown(index, $event)"
                   class="w-12 h-14 sm:w-14 sm:h-16 text-center text-xl sm:text-2xl font-black bg-slate-50 dark:bg-slate-900 border-2 border-slate-200 dark:border-slate-700 rounded-2xl text-indigo-600 dark:text-indigo-400 focus:outline-none focus:border-indigo-600 focus:ring-4 focus:ring-indigo-500/20 transition-all shadow-inner" />
          </div>
        </div>

        <!-- Timer & Resend Link -->
        <div class="text-center text-xs text-slate-500 dark:text-slate-400 space-y-2">
          <div v-if="timer > 0" class="flex items-center justify-center gap-1.5 font-medium">
            <span>Mã OTP hết hạn sau:</span>
            <span class="font-bold text-indigo-600 dark:text-indigo-400">00:{{ timer < 10 ? '0' + timer : timer }}s</span>
          </div>
          <div v-else>
            <span>Chưa nhận được email?</span>
            <button type="button" @click="resendOtp" class="ml-1 text-indigo-600 dark:text-indigo-400 font-bold hover:underline inline-flex items-center gap-1">
              <RefreshCw size="13" /> Gửi lại mã OTP ngay
            </button>
          </div>
        </div>

        <!-- Actions -->
        <div class="space-y-3 pt-2">
          <button @click="handleVerifyOtp"
                  :disabled="loading || otpDigits.some(d => d === '')"
                  class="w-full py-3.5 px-6 rounded-2xl bg-gradient-to-r from-emerald-600 via-teal-600 to-indigo-600 hover:from-emerald-700 hover:to-indigo-700 text-white font-bold text-sm shadow-lg shadow-emerald-600/30 hover:shadow-emerald-600/50 transition-all flex items-center justify-center gap-2 disabled:opacity-50 transform hover:-translate-y-0.5">
            <span v-if="loading" class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <template v-else>
              <CheckCircle2 size="18" />
              <span>Xác nhận & Hoàn tất Đăng ký</span>
            </template>
          </button>

          <button @click="step = 1; error = ''"
                  type="button"
                  class="w-full py-2.5 px-4 rounded-xl bg-slate-100 dark:bg-slate-700/60 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-600 dark:text-slate-300 font-semibold text-xs transition-colors flex items-center justify-center gap-1.5">
            <ArrowLeft size="14" />
            <span>Quay lại đổi thông tin đăng ký</span>
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(12px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-4px); }
  75% { transform: translateX(4px); }
}
.animate-fade-in {
  animation: fadeIn 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
.animate-shake {
  animation: shake 0.3s ease-in-out;
}
</style>
