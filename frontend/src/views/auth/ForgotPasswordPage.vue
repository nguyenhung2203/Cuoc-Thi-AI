<script setup>
import { ref, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Input from '../../components/common/AppInput.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mail, KeyRound, Lock, ArrowLeft, Sparkles, ShieldCheck, FileCode, TrendingUp } from 'lucide-vue-next'
import { authService } from '../../services/auth.service'
import {
  confirmPassword as confirmPasswordRule,
  isEmail,
  normalizeEmail,
  requiredTrim,
  validateForm,
  validatePassword,
} from '../../utils/validators.js'

const router = useRouter()
const step = ref(1)
const email = ref('')
const otp = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const toast = ref(null)
const fieldErrors = ref({})

// Timer state
const timeLeft = ref(60)
const canResend = ref(false)

let timer = null

watch([step, timeLeft], ([newStep, newTimeLeft]) => {
  if (newStep === 2 && newTimeLeft > 0) {
    if (!timer) {
      timer = setInterval(() => {
        timeLeft.value--
      }, 1000)
    }
  } else if (newTimeLeft === 0) {
    canResend.value = true
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }
}, { immediate: true })

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const handleSendEmail = async (e) => {
  e.preventDefault()
  if (loading.value) return
  const normalizedEmail = normalizeEmail(email.value)
  const validation = validateForm({ email: normalizedEmail }, {
    email: [(value) => requiredTrim(value, 'Vui lòng nhập địa chỉ email.'), isEmail],
  })
  fieldErrors.value = validation.errors
  if (!validation.isValid) return

  email.value = normalizedEmail
  loading.value = true
  try {
    await authService.forgotPassword(email.value)
    step.value = 2
    timeLeft.value = 60
    canResend.value = false
    toast.value = { type: 'success', message: 'Nếu email tồn tại, mã xác nhận đã được gửi đến hộp thư của bạn.' }
  } catch (err) {
    toast.value = { type: 'error', message: err.message || 'Không thể gửi mã xác nhận. Vui lòng thử lại.' }
  } finally {
    loading.value = false
  }
}

const handleResendOTP = async () => {
  if (!canResend.value) return
  try {
    await authService.forgotPassword(email.value)
    toast.value = { type: 'success', message: 'Đã gửi lại mã xác nhận mới!' }
    timeLeft.value = 60
    canResend.value = false
  } catch (err) {
    toast.value = { type: 'error', message: err.message || 'Không thể gửi lại mã.' }
  }
}

const handleVerifyOTP = async (e) => {
  e.preventDefault()
  if (loading.value) return
  fieldErrors.value = /^\d{6}$/.test(otp.value)
    ? {}
    : { otp: 'Mã OTP phải gồm đúng 6 chữ số.' }
  if (fieldErrors.value.otp) return

  loading.value = true
  try {
    await authService.verifyResetOtp({ email: email.value, otp: otp.value })
    step.value = 3
  } catch (err) {
    toast.value = { type: 'error', message: err.message || 'Mã xác nhận không hợp lệ hoặc đã hết hạn.' }
  } finally {
    loading.value = false
  }
}

const handleResetPassword = async (e) => {
  e.preventDefault()
  if (loading.value) return
  const validation = validateForm(
    { password: password.value, confirmPassword: confirmPassword.value },
    {
      password: [(value) => requiredTrim(value, 'Vui lòng nhập mật khẩu mới.'), validatePassword],
      confirmPassword: [
        (value) => requiredTrim(value, 'Vui lòng xác nhận mật khẩu mới.'),
        (value, values) => confirmPasswordRule(value, values.password),
      ],
    },
  )
  fieldErrors.value = validation.errors
  if (!validation.isValid) return

  loading.value = true
  try {
    await authService.resetPassword({ email: email.value, otp: otp.value, new_password: password.value })
    router.push({ path: '/login', state: { message: 'Đổi mật khẩu thành công! Vui lòng đăng nhập lại.' } })
  } catch (err) {
    toast.value = { type: 'error', message: err.message || 'Không thể đổi mật khẩu. Vui lòng thử lại.' }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-start justify-center bg-slate-950 px-4 pt-8 sm:pt-14 md:pt-16 pb-12 relative overflow-hidden" style="background-image: url('/images/auth_bg.png'); background-size: cover; background-position: center;">
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

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <Card class="w-full max-w-md bg-white/80 dark:bg-slate-900/80 backdrop-blur-2xl rounded-3xl shadow-2xl border border-white/20 dark:border-slate-800/80 p-8 sm:p-10 relative z-10 transition-all duration-500">
      
      <!-- Back to Login Action Header -->
      <div class="flex items-center gap-2 mb-6 pb-4 border-b border-slate-100 dark:border-slate-800/60">
        <button @click="router.push('/login')" class="flex items-center gap-1.5 text-xs font-bold transition-colors text-[var(--primary)] hover:text-[var(--primary-hover)] bg-transparent border-none cursor-pointer">
          <ArrowLeft size="14" />
          Quay lại Đăng nhập
        </button>
      </div>

      <!-- STEP 1: Nhập Email -->
      <div v-if="step === 1" style="animation: fadeIn 0.4s ease-out">
        <div style="text-align: center; margin-bottom: 32px">
          <div style="width: 56px; height: 56px; background-color: rgba(37, 99, 235, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; color: var(--primary)">
            <KeyRound size="28" />
          </div>
          <h1 class="text-h1" style="margin-bottom: 8px">Quên mật khẩu?</h1>
          <p class="text-helper" style="line-height: 1.5">Nhập địa chỉ email được liên kết với tài khoản của bạn để nhận mã xác nhận.</p>
        </div>

        <form @submit="handleSendEmail" novalidate style="display: flex; flex-direction: column; gap: 20px">
          <Input 
            label="Địa chỉ Email" 
            type="email" 
            placeholder="Ví dụ: candidate@test.com"
            v-model="email"
            :error="fieldErrors.email"
            required
          />
          <Button type="submit" :disabled="!email || loading" style="height: 44px">
            {{ loading ? 'Đang gửi...' : 'Gửi mã xác nhận' }}
          </Button>
        </form>
      </div>

      <!-- STEP 2: Xác thực OTP -->
      <div v-if="step === 2" style="animation: fadeIn 0.4s ease-out">
        <div style="text-align: center; margin-bottom: 32px">
          <div style="width: 56px; height: 56px; background-color: rgba(8, 145, 178, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; color: var(--accent)">
            <Mail size="28" />
          </div>
          <h1 class="text-h1" style="margin-bottom: 8px">Nhập mã xác nhận</h1>
          <p class="text-helper" style="line-height: 1.5">
            Chúng tôi đã gửi mã xác nhận 6 chữ số tới email<br/>
            <strong style="color: var(--text-main)">{{ email }}</strong>
          </p>
        </div>

        <form @submit="handleVerifyOTP" novalidate style="display: flex; flex-direction: column; gap: 24px">
          <div>
            <input 
              class="input-field"
              placeholder="Nhập mã 6 chữ số (VD: 123456)" 
              :value="otp"
              @input="e => otp = e.target.value.replace(/[^0-9]/g, '').slice(0, 6)"
              style="text-align: center; font-size: 20px; letter-spacing: 4px; font-weight: 600; width: 100%"
              :class="{ 'input-error': fieldErrors.otp }"
              required
            />
            <span v-if="fieldErrors.otp" class="error-text" style="display: block; margin-top: 6px; text-align: center">{{ fieldErrors.otp }}</span>
          </div>
          
          <Button type="submit" :disabled="otp.length < 6 || loading" style="height: 44px">
            {{ loading ? 'Đang xác thực...' : 'Xác nhận mã' }}
          </Button>

          <div style="text-align: center">
            <p class="text-helper">
              Chưa nhận được mã? 
              <span 
                @click="handleResendOTP" 
                :style="{ 
                  color: canResend ? 'var(--primary)' : 'var(--text-muted)', 
                  cursor: canResend ? 'pointer' : 'not-allowed',
                  fontWeight: 500
                }"
              >
                Gửi lại mã {{ canResend ? '' : `(${timeLeft}s)` }}
              </span>
            </p>
          </div>
        </form>
      </div>

      <!-- STEP 3: Đặt lại mật khẩu -->
      <div v-if="step === 3" style="animation: fadeIn 0.4s ease-out">
        <div style="text-align: center; margin-bottom: 32px">
          <div style="width: 56px; height: 56px; background-color: rgba(16, 185, 129, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; color: var(--success)">
            <Lock size="28" />
          </div>
          <h1 class="text-h1" style="margin-bottom: 8px">Tạo mật khẩu mới</h1>
          <p class="text-helper" style="line-height: 1.5">Mật khẩu mới của bạn phải khác với mật khẩu sử dụng trước đó.</p>
        </div>

        <form @submit="handleResetPassword" novalidate style="display: flex; flex-direction: column; gap: 20px">
          <Input 
            label="Mật khẩu mới" 
            type="password" 
            placeholder="Tối thiểu 8 ký tự, có chữ hoa, chữ thường và số"
            v-model="password"
            :error="fieldErrors.password"
            required
          />
          <Input 
            label="Xác nhận mật khẩu mới" 
            type="password" 
            placeholder="Nhập lại mật khẩu mới"
            v-model="confirmPassword"
            :error="fieldErrors.confirmPassword"
            required
          />
          
          <Button type="submit" :disabled="!password || !confirmPassword || loading" style="height: 44px; margin-top: 8px">
            {{ loading ? 'Đang đổi mật khẩu...' : 'Xác nhận đổi mật khẩu' }}
          </Button>
        </form>
      </div>

    </Card>
  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(12px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.animate-fade-in {
  animation: fadeIn 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
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
