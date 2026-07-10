<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Input from '../../components/common/AppInput.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { MailCheck, ArrowLeft } from 'lucide-vue-next'

import { authService } from '../../services/auth.service'

const router = useRouter()
const route = useRoute()
const email = route.query.email || ''
const otp = ref('')
const loading = ref(false)
const resending = ref(false)
const toast = ref(null)

if (!email) {
  router.push('/login')
}

const handleVerify = async (e) => {
  e.preventDefault()
  if (!otp.value) return
  loading.value = true
  try {
    await authService.verifyEmail(email, otp.value)
    toast.value = { type: 'success', message: 'Xác minh thành công! Đang chuyển hướng...' }
    setTimeout(() => {
      // Logic from login: if they just registered and verified, send them to login so they can actually log in,
      // or if they already have tokens in localStorage, route to home.
      // Usually, after verifying email, redirect to login to ensure fresh state.
      router.push('/login')
    }, 1500)
  } catch (err) {
    toast.value = { type: 'error', message: err.response?.data?.message || 'Có lỗi xảy ra, vui lòng thử lại.' }
  } finally {
    loading.value = false
  }
}

const handleResend = async () => {
  if (countdown.value > 0) return
  resending.value = true
  try {
    await authService.resendOTP(email, 'register')
    toast.value = { type: 'success', message: 'Mã xác nhận mới đã được gửi!' }
    startCountdown()
  } catch (err) {
    toast.value = { type: 'error', message: err.response?.data?.message || 'Có lỗi xảy ra, vui lòng thử lại.' }
  } finally {
    resending.value = false
  }
}

const countdown = ref(0)
let timer = null

const startCountdown = () => {
  countdown.value = 60
  if (timer) clearInterval(timer)
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
    }
  }, 1000)
}

// Optionally start countdown immediately when page loads since email was just sent
startCountdown()

import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (timer) clearInterval(timer)
})
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
      <div style="animation: fadeIn 0.4s ease-out">
        <div style="text-align: center; margin-bottom: 32px">
          <div style="width: 56px; height: 56px; background-color: rgba(16, 185, 129, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; color: var(--success)">
            <MailCheck size="28" />
          </div>
          <h1 class="text-h1" style="margin-bottom: 8px">Xác minh Email</h1>
          <p class="text-helper" style="line-height: 1.5">
            Chúng tôi đã gửi một mã xác nhận gồm 6 chữ số đến email<br>
            <strong>{{ email }}</strong>
          </p>
        </div>

        <form @submit="handleVerify" style="display: flex; flex-direction: column; gap: 20px">
          <Input 
            label="Mã OTP" 
            type="text" 
            placeholder="Nhập mã 6 chữ số"
            v-model="otp"
            required
            maxlength="6"
            style="text-align: center; letter-spacing: 4px; font-size: 1.25rem;"
          />
          <Button type="submit" :disabled="!otp || loading" style="height: 44px">
            {{ loading ? 'Đang xác minh...' : 'Xác minh' }}
          </Button>
        </form>

        <div style="margin-top: 24px; text-align: center;">
          <p class="text-helper">
            Chưa nhận được mã? 
            <a href="#" @click.prevent="handleResend" 
               :style="{ color: countdown > 0 ? 'var(--text-muted)' : 'var(--primary)', fontWeight: '500', textDecoration: 'none', cursor: countdown > 0 ? 'not-allowed' : 'pointer', pointerEvents: countdown > 0 ? 'none' : 'auto' }">
              {{ countdown > 0 ? `Gửi lại mã sau ${countdown}s` : (resending ? 'Đang gửi...' : 'Gửi lại mã') }}
            </a>
          </p>
        </div>
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
