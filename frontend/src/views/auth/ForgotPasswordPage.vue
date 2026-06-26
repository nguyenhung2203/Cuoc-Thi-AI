<script setup>
import { ref, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Input from '../../components/common/AppInput.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mail, KeyRound, Lock, ArrowLeft } from 'lucide-vue-next'

const router = useRouter()
const step = ref(1)
const email = ref('')
const otp = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const toast = ref(null)

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

const handleSendEmail = (e) => {
  e.preventDefault()
  if (!email.value) return
  loading.value = true
  setTimeout(() => {
    loading.value = false
    step.value = 2
    timeLeft.value = 60
    canResend.value = false
    toast.value = { type: 'success', message: 'Mã xác nhận đã được gửi đến email của bạn!' }
  }, 1000)
}

const handleResendOTP = () => {
  if (!canResend.value) return
  toast.value = { type: 'success', message: 'Đã gửi lại mã xác nhận mới!' }
  timeLeft.value = 60
  canResend.value = false
}

const handleVerifyOTP = (e) => {
  e.preventDefault()
  if (!otp.value) return
  if (otp.value !== '123456') {
    toast.value = { type: 'error', message: 'Mã xác nhận không hợp lệ. Vui lòng thử lại!' }
    return
  }
  loading.value = true
  setTimeout(() => {
    loading.value = false
    step.value = 3
  }, 800)
}

const handleResetPassword = (e) => {
  e.preventDefault()
  if (!password.value || !confirmPassword.value) return
  if (password.value !== confirmPassword.value) {
    toast.value = { type: 'error', message: 'Mật khẩu xác nhận không khớp!' }
    return
  }
  if (password.value.length < 6) {
    toast.value = { type: 'error', message: 'Mật khẩu phải có ít nhất 6 ký tự!' }
    return
  }
  loading.value = true
  setTimeout(() => {
    loading.value = false
    router.push({ path: '/login', state: { message: 'Đổi mật khẩu thành công! Vui lòng đăng nhập lại.' } })
  }, 1000)
}
</script>

<template>
  <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background-color: var(--background)">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
    
    <div style="position: absolute; top: 32px; left: 32px">
      <Button variant="ghost" @click="router.push('/login')">
        <ArrowLeft size="18" style="margin-right: 8px" /> Quay lại Đăng nhập
      </Button>
    </div>

    <Card style="width: 100%; max-width: 440px; padding: 40px; position: relative; overflow: hidden">
      
      <!-- STEP 1: Nhập Email -->
      <div v-if="step === 1" style="animation: fadeIn 0.4s ease-out">
        <div style="text-align: center; margin-bottom: 32px">
          <div style="width: 56px; height: 56px; background-color: rgba(37, 99, 235, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; color: var(--primary)">
            <KeyRound size="28" />
          </div>
          <h1 class="text-h1" style="margin-bottom: 8px">Quên mật khẩu?</h1>
          <p class="text-helper" style="line-height: 1.5">Nhập địa chỉ email được liên kết với tài khoản của bạn để nhận mã xác nhận.</p>
        </div>

        <form @submit="handleSendEmail" style="display: flex; flex-direction: column; gap: 20px">
          <Input 
            label="Địa chỉ Email" 
            type="email" 
            placeholder="Ví dụ: candidate@test.com"
            v-model="email"
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

        <form @submit="handleVerifyOTP" style="display: flex; flex-direction: column; gap: 24px">
          <div>
            <input 
              class="input-field"
              placeholder="Nhập mã 6 chữ số (VD: 123456)" 
              :value="otp"
              @input="e => otp = e.target.value.replace(/[^0-9]/g, '').slice(0, 6)"
              style="text-align: center; font-size: 20px; letter-spacing: 4px; font-weight: 600; width: 100%"
              required
            />
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

        <form @submit="handleResetPassword" style="display: flex; flex-direction: column; gap: 20px">
          <Input 
            label="Mật khẩu mới" 
            type="password" 
            placeholder="Tối thiểu 6 ký tự"
            v-model="password"
            required
          />
          <Input 
            label="Xác nhận mật khẩu mới" 
            type="password" 
            placeholder="Nhập lại mật khẩu mới"
            v-model="confirmPassword"
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

<style>
@keyframes fadeIn {
  from { opacity: 0; transform: translateX(20px); }
  to { opacity: 1; transform: translateX(0); }
}
</style>
