<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Input from '../../components/common/AppInput.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { Lock } from 'lucide-vue-next'
import { authService } from '../../services/auth.service'

const router = useRouter()
const route = useRoute()
const otp = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const toast = ref(null)

const email = route.query.email

if (!email) {
  router.push('/login')
}

const handleResetPassword = async (e) => {
  e.preventDefault()
  if (!otp.value || !password.value || !confirmPassword.value) return
  if (password.value !== confirmPassword.value) {
    toast.value = { type: 'error', message: 'Mật khẩu xác nhận không khớp!' }
    return
  }
  if (password.value.length < 6) {
    toast.value = { type: 'error', message: 'Mật khẩu phải có ít nhất 6 ký tự!' }
    return
  }
  
  loading.value = true
  try {
    await authService.resetPassword(email, otp.value, password.value)
    toast.value = { type: 'success', message: 'Đổi mật khẩu thành công! Đang chuyển về Đăng nhập...' }
    setTimeout(() => {
      router.push('/login')
    }, 2000)
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

    <Card style="width: 100%; max-width: 440px; padding: 40px; position: relative; overflow: hidden">
      <div style="animation: fadeIn 0.4s ease-out">
        <div style="text-align: center; margin-bottom: 32px">
          <div style="width: 56px; height: 56px; background-color: rgba(16, 185, 129, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; color: var(--success)">
            <Lock size="28" />
          </div>
          <h1 class="text-h1" style="margin-bottom: 8px">Tạo mật khẩu mới</h1>
          <p class="text-helper" style="line-height: 1.5">Mật khẩu mới của bạn phải khác với mật khẩu sử dụng trước đó.</p>
        </div>

        <form @submit="handleResetPassword" style="display: flex; flex-direction: column; gap: 20px">
          <Input 
            label="Mã OTP" 
            type="text" 
            placeholder="Nhập mã 6 số từ email"
            v-model="otp"
            required
            maxlength="6"
          />
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
          
          <Button type="submit" :disabled="!otp || !password || !confirmPassword || loading" style="height: 44px; margin-top: 8px">
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
