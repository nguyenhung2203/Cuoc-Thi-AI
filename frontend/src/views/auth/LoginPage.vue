<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Toast from '../../components/common/AppToast.vue'
import AppLogo from '../../components/common/AppLogo.vue'
import { authStore } from '../../stores/auth.store'
import { isEmail, normalizeEmail, requiredTrim, validateForm } from '../../utils/validators.js'
import { useGoogleAuth } from '../../composables/useGoogleAuth'
import { Mail, Lock, ArrowRight, ShieldCheck, AlertCircle, Sparkles, KeyRound, FileCode, TrendingUp } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const error = ref('')
const fieldErrors = ref({})
const loading = ref(false)
const entryToast = ref(history.state?.message ? { type: history.state.type || 'success', message: history.state.message } : null)

const { available: googleAvailable, signIn: googleSignIn } = useGoogleAuth()

onMounted(() => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }
})

const redirectAfterAuth = (user, googleFlow = false) => {
  const redirectPath = route.query.redirect
  const prefix = googleFlow ? 'Đăng nhập Google thành công! ' : ''
  if (redirectPath) {
    router.push(redirectPath)
  } else if (user.role === 'recruiter' || user.role === 'admin' || user.role === 'owner') {
    router.push({ path: '/dashboard', state: { message: `${prefix}Chào mừng ${user.full_name} quay trở lại Bảng điều khiển Quản lý!` } })
  } else {
    router.push({ path: '/home', state: { message: `${prefix}Đăng nhập thành công! Chào mừng trở lại, ${user.full_name || 'Ứng viên'}.` } })
  }
}

const handleLogin = async (e) => {
  e.preventDefault()
  if (loading.value) return

  const normalizedEmail = normalizeEmail(email.value)
  const validation = validateForm(
    { email: normalizedEmail, password: password.value },
    {
      email: [requiredTrim, isEmail],
      password: [(value) => requiredTrim(value, 'Vui lòng nhập mật khẩu.')],
    },
  )
  fieldErrors.value = validation.errors
  if (!validation.isValid) {
    error.value = ''
    return
  }

  email.value = normalizedEmail
  error.value = ''
  loading.value = true

  try {
    const user = await authStore.login(email.value, password.value)
    redirectAfterAuth(user)
  } catch (err) {
    error.value = err.code === 'INVALID_CREDENTIALS'
      ? 'Email hoặc mật khẩu không chính xác. Vui lòng kiểm tra lại!'
      : (err.message || 'Đăng nhập thất bại, vui lòng kiểm tra lại kết nối hoặc tài khoản.')
  } finally {
    loading.value = false
  }
}

const handleGoogleLogin = async () => {
  error.value = ''
  loading.value = true
  try {
    const idToken = await googleSignIn()
    const user = await authStore.loginWithGoogle(idToken, 'candidate')
    redirectAfterAuth(user, true)
  } catch (err) {
    error.value = 'Đăng nhập Google không thành công: ' + (err.message || 'Lỗi xác thực OAuth')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-950 px-4 py-6 relative overflow-hidden" style="background-image: url('/images/auth_bg.png'); background-size: cover; background-position: center;">
    <!-- Grid Overlay & Blur -->
    <div class="absolute inset-0 bg-slate-950/40 backdrop-blur-[3px] pointer-events-none"></div>
    <div class="absolute inset-0 bg-[linear-gradient(to_right,rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:3.5rem_3.5rem] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_50%,#000_70%,transparent_100%)] pointer-events-none"></div>

    <!-- Premium Animated Glows -->
    <div class="absolute -top-40 -left-40 w-[600px] h-[600px] bg-blue-600/10 rounded-full blur-[140px] pointer-events-none animate-float-1"></div>
    <div class="absolute -bottom-40 -right-40 w-[600px] h-[600px] bg-cyan-500/10 rounded-full blur-[140px] pointer-events-none animate-float-2"></div>
    <div class="absolute top-1/4 right-1/4 w-[400px] h-[400px] bg-indigo-500/5 rounded-full blur-[100px] pointer-events-none animate-float-3"></div>

    <!-- Floating Decorative AI Badges (Desktop Only) -->
    <div class="hidden lg:block absolute top-16 left-16 p-4 bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-white/10 shadow-2xl animate-float-1 pointer-events-none">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-blue-500/20 flex items-center justify-center text-blue-400">
          <Sparkles size="18" />
        </div>
        <div>
          <div class="text-[11px] font-bold text-slate-200">AI Recruiter Engine</div>
          <div class="text-[9px] text-emerald-400 font-semibold">Active & Online</div>
        </div>
      </div>
    </div>

    <div class="hidden lg:block absolute bottom-20 left-24 p-4 bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-white/10 shadow-2xl animate-float-2 pointer-events-none">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-indigo-500/20 flex items-center justify-center text-indigo-400">
          <ShieldCheck size="18" />
        </div>
        <div>
          <div class="text-[11px] font-bold text-slate-200">Biometric Verification</div>
          <div class="text-[9px] text-indigo-300 font-semibold">2FA Secured</div>
        </div>
      </div>
    </div>

    <div class="hidden lg:block absolute top-1/4 right-16 p-4 bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-white/10 shadow-2xl animate-float-3 pointer-events-none">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-amber-500/20 flex items-center justify-center text-amber-400">
          <KeyRound size="18" />
        </div>
        <div>
          <div class="text-[11px] font-bold text-slate-200">LiveKit Audio Room</div>
          <div class="text-[9px] text-amber-300 font-semibold">Ready to connect</div>
        </div>
      </div>
    </div>

    <!-- Extra Floating elements for Rich Aesthetics -->
    <div class="hidden lg:block absolute bottom-[18%] right-[10%] p-4 bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-white/10 shadow-2xl animate-float-1 pointer-events-none" style="animation-delay: 1.5s;">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-full bg-emerald-500/20 flex items-center justify-center text-emerald-400 font-extrabold text-xs">
          94%
        </div>
        <div>
          <div class="text-[11px] font-bold text-slate-200">Resume Match Score</div>
          <div class="text-[9px] text-slate-400">Calibrated to Senior React JD</div>
        </div>
      </div>
    </div>

    <div class="hidden lg:block absolute top-[10%] right-[30%] p-3 bg-slate-900/60 backdrop-blur-xl rounded-xl border border-white/10 shadow-2xl animate-float-2 pointer-events-none" style="animation-delay: 3s;">
      <div class="flex items-center gap-2">
        <span class="w-2.5 h-2.5 rounded-full bg-red-500 animate-ping shrink-0"></span>
        <span class="text-[10px] text-slate-300 font-bold">AI Transcript: "I scaled DB performance by 40%..."</span>
      </div>
    </div>

    <!-- Extra Floating Widget 6: AI Code Analysis -->
    <div class="hidden lg:block absolute top-[45%] left-[8%] p-4 bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-white/10 shadow-2xl animate-float-3 pointer-events-none" style="animation-delay: 0.5s;">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-purple-500/20 flex items-center justify-center text-purple-400">
          <FileCode size="18" />
        </div>
        <div>
          <div class="text-[11px] font-bold text-slate-200">AI Code Analysis</div>
          <div class="text-[9px] text-purple-300 font-semibold">Complexity: O(N log N)</div>
        </div>
      </div>
    </div>

    <!-- Extra Floating Widget 7: AI Speech Insights -->
    <div class="hidden lg:block absolute bottom-[8%] right-[25%] p-4 bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-white/10 shadow-2xl animate-float-2 pointer-events-none" style="animation-delay: 2.2s;">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-cyan-500/20 flex items-center justify-center text-cyan-400">
          <TrendingUp size="18" />
        </div>
        <div>
          <div class="text-[11px] font-bold text-slate-200">Speech Insights</div>
          <div class="text-[9px] text-cyan-300 font-semibold">Confidence Level: 92%</div>
        </div>
      </div>
    </div>

    <!-- Extra Floating Widget 8: Skills Badge -->
    <div class="hidden lg:block absolute top-[12%] left-[32%] p-3.5 bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-white/10 shadow-2xl animate-float-1 pointer-events-none" style="animation-delay: 1.1s;">
      <div class="flex items-center gap-2">
        <span class="px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 text-[9px] font-bold">React</span>
        <span class="px-2 py-0.5 rounded bg-purple-500/10 text-purple-400 text-[9px] font-bold">Go</span>
        <span class="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 text-[9px] font-bold">System Design</span>
      </div>
    </div>

    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />

    <!-- Main Login Box -->
    <div class="w-full max-w-md bg-white/85 dark:bg-slate-900/85 backdrop-blur-2xl rounded-3xl shadow-2xl border border-white/20 dark:border-slate-800/80 p-6 sm:p-8 relative z-10 transition-all duration-500">
      
      <!-- Brand Logo Header -->
      <div class="text-center space-y-3 mb-6 flex flex-col items-center">
        <div @click="router.push('/')" class="inline-block cursor-pointer transform hover:scale-105 transition-transform">
          <AppLogo size="lg" showSubtitle subtitle="Enterprise AI Platform" />
        </div>
        <h1 class="text-xl sm:text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight pt-1">
          Chào mừng trở lại
        </h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm">
          Đăng nhập vào nền tảng tuyển dụng thông minh AI
        </p>
      </div>

      <!-- Error Alert -->
      <div v-if="error" class="mb-6 p-4 rounded-2xl bg-rose-50 dark:bg-rose-950/50 border border-rose-200 dark:border-rose-800/60 flex items-start gap-3 text-rose-600 dark:text-rose-400 text-sm animate-shake">
        <AlertCircle size="20" class="shrink-0 mt-0.5" />
        <span class="font-medium flex-1">{{ error }}</span>
      </div>

      <!-- ================= EMAIL & GOOGLE LOGIN ================= -->
      <div class="space-y-4 animate-fade-in">
        <form @submit="handleLogin" novalidate class="space-y-3">
          <!-- Email Input -->
          <div>
            <label class="block text-xs font-bold text-slate-600 dark:text-slate-300 uppercase tracking-wider mb-1.5">Địa chỉ Email</label>
            <div class="relative">
              <Mail size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
              <input v-model="email"
                     type="email"
                     required
                     placeholder="nhapemail@congty.com"
                     :class="fieldErrors.email ? 'border-rose-500 focus:ring-rose-500' : ''"
                     class="w-full pl-11 pr-4 py-3 bg-slate-50 dark:bg-slate-900/80 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 font-medium text-sm transition-all" />
            </div>
            <p v-if="fieldErrors.email" class="mt-1.5 text-xs font-medium text-rose-600">{{ fieldErrors.email }}</p>
          </div>

          <!-- Password Input -->
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <label class="block text-xs font-bold text-slate-600 dark:text-slate-300 uppercase tracking-wider">Mật khẩu</label>
              <span @click="router.push('/forgot-password')" class="text-xs text-[var(--accent)] font-bold hover:underline cursor-pointer">
                Quên mật khẩu?
              </span>
            </div>
            <div class="relative">
              <Lock size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
              <input v-model="password"
                     type="password"
                     required
                     placeholder="••••••••"
                     :class="fieldErrors.password ? 'border-rose-500 focus:ring-rose-500' : ''"
                     class="w-full pl-11 pr-4 py-3 bg-slate-50 dark:bg-slate-900/80 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 font-medium text-sm transition-all" />
            </div>
            <p v-if="fieldErrors.password" class="mt-1.5 text-xs font-medium text-rose-600">{{ fieldErrors.password }}</p>
          </div>

          <!-- Submit Button -->
          <button type="submit"
                  :disabled="loading"
                  class="w-full mt-2 py-3.5 px-6 rounded-2xl bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white font-bold text-sm shadow-md transition-all flex items-center justify-center gap-2.5 disabled:opacity-50 transform hover:-translate-y-0.5">
            <span v-if="loading" class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            <template v-else>
              <span>Đăng nhập vào hệ thống</span>
              <ArrowRight size="18" />
            </template>
          </button>
        </form>

        <div v-if="googleAvailable" class="relative flex py-1 items-center">
          <div class="flex-grow border-t border-slate-200 dark:border-slate-700"></div>
          <span class="flex-shrink mx-4 text-xs font-semibold text-slate-400 uppercase tracking-wider">hoặc đăng nhập nhanh qua google</span>
          <div class="flex-grow border-t border-slate-200 dark:border-slate-700"></div>
        </div>

        <!-- Google OAuth Button (real Google Identity Services) -->
        <button v-if="googleAvailable"
                @click="handleGoogleLogin"
                type="button"
                :disabled="loading"
                class="w-full py-3.5 px-4 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 hover:border-blue-500/50 dark:hover:border-blue-500/50 hover:bg-slate-50 dark:hover:bg-slate-800/80 text-slate-700 dark:text-slate-200 font-semibold text-sm shadow-sm transition-all flex items-center justify-center gap-3 group disabled:opacity-50">
          <svg class="w-5 h-5 transition-transform group-hover:scale-110" viewBox="0 0 24 24">
            <path fill="#EA4335" d="M12 5c1.6 0 3 .6 4.1 1.7l3.1-3.1C17.3 1.8 14.8 1 12 1 7.4 1 3.5 3.6 1.6 7.4l3.7 2.8C6.2 7.3 8.9 5 12 5z"/>
            <path fill="#4285F4" d="M23.5 12.3c0-.8-.1-1.7-.2-2.3H12v4.6h6.5c-.3 1.5-1.1 2.8-2.4 3.7l3.7 2.9c2.2-2 3.7-5 3.7-8.9z"/>
            <path fill="#FBBC05" d="M5.3 14.8c-.2-.7-.4-1.5-.4-2.3s.2-1.5.4-2.3L1.6 7.4C.6 9.4 0 11.6 0 14s.6 4.6 1.6 6.6l3.7-2.8z"/>
            <path fill="#34A853" d="M12 23c3.2 0 6-1.1 8-3l-3.7-2.9c-1.1.7-2.5 1.2-4.3 1.2-3.1 0-5.8-2.3-6.7-5.2L1.6 15.9C3.5 19.7 7.4 23 12 23z"/>
          </svg>
          <span>Đăng nhập bằng tài khoản Google</span>
        </button>

        <!-- Footer Link -->
        <div class="pt-2 text-center border-t border-slate-100 dark:border-slate-700/60">
          <p class="text-sm text-slate-500 dark:text-slate-400">
            Chưa có tài khoản ViệcLàm AI?
            <span @click="router.push('/register')" class="text-[var(--primary)] font-bold hover:underline cursor-pointer ml-1">
              Đăng ký tài khoản mới
            </span>
          </p>
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

@keyframes float1 {
  0%, 100% { transform: translateY(0px) rotate(0deg) scale(1); }
  50% { transform: translateY(-15px) rotate(2deg) scale(1.03); }
}
@keyframes float2 {
  0%, 100% { transform: translateY(0px) rotate(0deg) scale(1); }
  50% { transform: translateY(-22px) rotate(-3deg) scale(1.05); }
}
@keyframes float3 {
  0%, 100% { transform: translateY(0px) rotate(0deg) scale(1); }
  50% { transform: translateY(18px) rotate(1.5deg) scale(0.97); }
}

.animate-float-1 {
  animation: float1 8s infinite ease-in-out;
}
.animate-float-2 {
  animation: float2 12s infinite ease-in-out;
}
.animate-float-3 {
  animation: float3 10s infinite ease-in-out;
}
</style>
