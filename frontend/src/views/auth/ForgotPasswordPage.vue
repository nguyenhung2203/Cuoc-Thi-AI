<script setup>
import { ref, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Input from '../../components/common/AppInput.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mail, KeyRound, Lock, ArrowLeft } from 'lucide-vue-next'

import { authService } from '../../services/auth.service'

const router = useRouter()
const step = ref(1)
const email = ref('')
const loading = ref(false)
const toast = ref(null)

const handleSendEmail = async (e) => {
  e.preventDefault()
  if (!email.value) return
  loading.value = true
  try {
    await authService.forgotPassword(email.value)
    toast.value = { type: 'success', message: 'Mã xác nhận đã được gửi!' }
    setTimeout(() => {
      router.push({ path: '/reset-password', query: { email: email.value } })
    }, 1500)
  } catch (err) {
    toast.value = { type: 'error', message: err.response?.data?.message || 'Có lỗi xảy ra, vui lòng thử lại.' }
  } finally {
    loading.value = false
  }
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

      <!-- Step 2: Thành công -->
      <div v-else-if="step === 2" class="auth-form" style="text-align: center; animation: fadeIn 0.4s ease-out">
        <div class="icon-circle" style="margin: 0 auto 24px; width: 56px; height: 56px; background-color: rgba(37, 99, 235, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center;">
          <Mail size="24" class="text-primary" />
        </div>
        <h2 class="text-h2" style="margin-bottom: 12px">Kiểm tra email của bạn</h2>
        <p class="text-body text-secondary" style="margin-bottom: 24px; color: var(--text-muted)">
          Chúng tôi đã gửi một liên kết đặt lại mật khẩu đến<br>
          <strong>{{ email }}</strong>
        </p>

        <Button variant="outline" style="width: 100%; margin-top: 16px" @click="$router.push('/login')">Quay lại Đăng nhập</Button>
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
