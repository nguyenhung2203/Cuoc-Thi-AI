<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Input from '../../components/common/AppInput.vue'
import Button from '../../components/common/AppButton.vue'

const router = useRouter()
const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const role = ref('candidate')
const error = ref('')
const loading = ref(false)

const handleRegister = (e) => {
  e.preventDefault()
  if (!name.value || !email.value || !password.value || !confirmPassword.value) {
    error.value = 'Vui lòng nhập đầy đủ thông tin'
    return
  }
  
  if (password.value !== confirmPassword.value) {
    error.value = 'Mật khẩu xác nhận không khớp'
    return
  }
  
  error.value = ''
  loading.value = true

  setTimeout(() => {
    loading.value = false
    router.push({ path: '/login', state: { message: 'Đăng ký thành công! Vui lòng đăng nhập.' } })
  }, 1000)
}
</script>

<template>
  <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background-color: var(--background)">
    <div style="display: flex; flex-direction: column; gap: 32px; width: 100%; max-width: 450px">
      <div style="text-align: center">
        <h1 class="text-h1" style="margin-bottom: 8px; color: var(--primary)">Interview AI</h1>
        <p class="text-body" style="color: var(--text-secondary)">Tạo tài khoản mới</p>
      </div>
      
      <Card>
        <form @submit="handleRegister">
          <Input 
            label="Họ và tên" 
            v-model="name"
            required
            placeholder="Nguyễn Văn A"
          />
          <Input 
            label="Email" 
            type="email" 
            v-model="email"
            required
            placeholder="nhapemail@congty.com"
          />
          
          <div style="margin-bottom: 16px">
            <label class="input-label">Bạn là:</label>
            <div style="display: flex; gap: 16px; margin-top: 8px">
              <label style="display: flex; align-items: center; gap: 8px; cursor: pointer">
                <input type="radio" value="candidate" v-model="role" name="role" />
                <span>Ứng viên</span>
              </label>
              <label style="display: flex; align-items: center; gap: 8px; cursor: pointer">
                <input type="radio" value="recruiter" v-model="role" name="role" />
                <span>Nhà tuyển dụng</span>
              </label>
            </div>
          </div>
          
          <Input 
            label="Mật khẩu" 
            type="password" 
            v-model="password"
            required
            placeholder="••••••••"
          />
          <Input 
            label="Xác nhận mật khẩu" 
            type="password" 
            v-model="confirmPassword"
            required
            placeholder="••••••••"
          />
          
          <div v-if="error" class="error-text" style="margin-bottom: 16px">{{ error }}</div>
          
          <Button type="submit" style="width: 100%; margin-top: 16px" :disabled="loading">
            {{ loading ? 'Đang xử lý...' : 'Đăng ký tài khoản' }}
          </Button>
        </form>
        
        <div style="margin-top: 24px; text-align: center">
          <p class="text-body" style="color: var(--text-secondary)">
            Đã có tài khoản? <span style="color: var(--primary); cursor: pointer; font-weight: 500" @click="router.push('/login')">Đăng nhập</span>
          </p>
        </div>
      </Card>
    </div>
  </div>
</template>
