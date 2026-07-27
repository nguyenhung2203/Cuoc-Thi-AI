<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authStore } from '../../stores/auth.store'
import { isEmail, normalizeEmail, requiredTrim, validateForm } from '../../utils/validators.js'
import { ShieldCheck, Lock, Mail, ArrowLeft, RefreshCw, AlertCircle } from 'lucide-vue-next'

const router = useRouter()

const email = ref('')
const password = ref('')
const error = ref('')
const fieldErrors = ref({})
const loading = ref(false)

const handleLogin = async (e) => {
  e.preventDefault()
  if (loading.value) return

  const normalizedEmail = normalizeEmail(email.value)
  const validation = validateForm(
    { email: normalizedEmail, password: password.value },
    {
      email: [(value) => requiredTrim(value, 'Vui lòng nhập email quản trị.'), isEmail],
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
    
    if (user.role === 'admin') {
      router.push('/admin/dashboard')
    } else {
      await authStore.logout()
      error.value = 'Tài khoản của bạn không có quyền truy cập vào Cổng Quản trị viên (Admin Portal).'
    }
  } catch (err) {
    error.value = err.code === 'INVALID_CREDENTIALS'
      ? 'Email hoặc mật khẩu không chính xác.'
      : (err.message || 'Đăng nhập thất bại, vui lòng kiểm tra lại.')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-950 via-slate-900 to-blue-950 px-4 relative overflow-hidden font-sans">
    <!-- Background Glow -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-blue-500/10 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-blue-500/10 rounded-full blur-3xl pointer-events-none"></div>

    <div class="w-full max-w-md relative z-10">
      <!-- Brand & Title -->
      <div class="text-center mb-8">
        <div class="w-16 h-16 bg-gradient-to-tr from-blue-600 to-blue-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-xl shadow-blue-500/30 border border-blue-400/20">
          <ShieldCheck size="36" class="text-white" />
        </div>
        <h1 class="text-3xl font-extrabold text-white tracking-tight">WeMake <span class="text-blue-400">Admin</span></h1>
        <p class="text-slate-400 text-sm mt-1">Cổng Quản trị Hệ sinh thái Tuyển dụng AI</p>
      </div>
      
      <!-- Login Card -->
      <div class="bg-slate-900/80 backdrop-blur-xl border border-slate-800 rounded-3xl p-8 shadow-2xl">
        <form @submit="handleLogin" novalidate class="space-y-5">
          <div>
            <label class="block text-xs font-bold text-slate-300 uppercase tracking-wider mb-2">Email Quản trị</label>
            <div class="relative">
              <Mail size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" />
              <input 
                type="email"
                v-model="email"
                required
                :class="fieldErrors.email ? 'border-red-500 focus:ring-red-500' : ''"
                class="w-full pl-11 pr-4 py-3 bg-slate-950/80 border border-slate-800 rounded-xl text-white placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm transition-all"
                placeholder="admin@wemake.vn"
              />
            </div>
            <p v-if="fieldErrors.email" class="mt-1.5 text-xs font-medium text-red-400">{{ fieldErrors.email }}</p>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-300 uppercase tracking-wider mb-2">Mật khẩu</label>
            <div class="relative">
              <Lock size="18" class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" />
              <input 
                type="password"
                v-model="password"
                required
                :class="fieldErrors.password ? 'border-red-500 focus:ring-red-500' : ''"
                class="w-full pl-11 pr-4 py-3 bg-slate-950/80 border border-slate-800 rounded-xl text-white placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm transition-all"
                placeholder="••••••••"
              />
            </div>
            <p v-if="fieldErrors.password" class="mt-1.5 text-xs font-medium text-red-400">{{ fieldErrors.password }}</p>
          </div>

          <div v-if="error" class="p-3.5 bg-red-500/10 border border-red-500/20 rounded-xl text-red-400 text-xs flex items-center gap-2">
            <AlertCircle size="16" class="shrink-0 text-red-400" />
            <span>{{ error }}</span>
          </div>
          
          <button 
            type="submit" 
            :disabled="loading"
            class="w-full py-3.5 bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-700 hover:to-blue-600 text-white font-bold text-sm rounded-xl transition-all duration-200 disabled:opacity-70 flex items-center justify-center gap-2 shadow-lg shadow-blue-600/30 hover:shadow-blue-600/50"
          >
            <RefreshCw v-if="loading" size="18" class="animate-spin" />
            {{ loading ? 'Đang xác thực bảo mật...' : 'Đăng nhập vào Hệ thống' }}
          </button>
        </form>
        
        <div class="mt-6 pt-6 border-t border-slate-800/80 text-center">
          <button @click="router.push('/login')" class="text-xs font-medium text-slate-500 hover:text-slate-300 transition-colors flex items-center justify-center gap-1.5 mx-auto">
            <ArrowLeft size="14" /> Quay lại trang đăng nhập User
          </button>
        </div>
      </div>

      <!-- Quick Info / Hint for Dev/Testing -->
      <div class="mt-6 text-center">
        <span class="inline-block px-3 py-1 bg-slate-900/60 border border-slate-800 rounded-full text-[11px] text-slate-500 font-mono">
          Mặc định: admin@wemake.vn / admin123
        </span>
      </div>
    </div>
  </div>
</template>
