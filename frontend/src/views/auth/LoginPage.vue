<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Input from '../../components/common/AppInput.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'

const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)

onMounted(() => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }
})

import { authStore } from '../../stores/auth.store'

const handleLogin = async (e) => {
  e.preventDefault()
  if (!email.value || !password.value) {
    error.value = 'Vui lòng nhập đầy đủ email và mật khẩu'
    return
  }
  
  error.value = ''
  loading.value = true

  try {
    const user = await authStore.login(email.value, password.value)
    
    // Check for redirect query param (e.g. from Career Site)
    const redirectPath = route.query.redirect
    if (redirectPath) {
      router.push(redirectPath)
    } else if (user.role === 'recruiter' || user.role === 'admin' || user.role === 'owner') {
      router.push({ path: '/dashboard', state: { message: `Chào mừng ${user.full_name} quay trở lại màn hình quản lý!` } })
    } else {
      router.push({ path: '/home', state: { message: `Đăng nhập thành công! Chào mừng ${user.full_name}.` } })
    }
  } catch (err) {
    if (err.message === 'invalid email or password') {
      error.value = 'Email hoặc mật khẩu không chính xác.'
    } else {
      error.value = err.message || 'Đăng nhập thất bại, vui lòng kiểm tra lại email hoặc mật khẩu.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background-color: var(--background)">
    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />
    
    <div style="display: flex; flex-direction: column; gap: 32px; width: 100%; max-width: 400px">
      <div style="text-align: center">
        <h1 class="text-h1" style="margin-bottom: 8px; color: var(--primary)">Interview AI</h1>
        <p class="text-body" style="color: var(--text-secondary)">Đăng nhập vào hệ thống tuyển dụng</p>
      </div>
      
      <Card>
        <form @submit="handleLogin">
          <Input 
            label="Email" 
            type="email" 
            v-model="email"
            required
            placeholder="nhapemail@congty.com"
          />
          <Input 
            label="Mật khẩu" 
            type="password" 
            v-model="password"
            required
            placeholder="••••••••"
          />
          <div style="display: flex; justify-content: flex-end; margin-top: -8px; margin-bottom: 16px">
            <span @click="router.push('/forgot-password')" style="font-size: 13px; color: var(--primary); font-weight: 500; cursor: pointer">Quên mật khẩu?</span>
          </div>
          
          <div v-if="error" class="error-text" style="margin-bottom: 16px">{{ error }}</div>
          
          <Button type="submit" style="width: 100%; margin-top: 16px" :disabled="loading">
            {{ loading ? 'Đang xử lý...' : 'Đăng nhập' }}
          </Button>
        </form>
        
        <div style="margin-top: 24px; text-align: center">
          <p class="text-body" style="color: var(--text-secondary)">
            Chưa có tài khoản? <span style="color: var(--primary); cursor: pointer; font-weight: 500" @click="router.push('/register')">Đăng ký ngay</span>
          </p>
        </div>
      </Card>
    </div>
  </div>
</template>
