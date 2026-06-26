<script setup>
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import { Clock, Mail, ArrowLeft } from 'lucide-vue-next'

const router = useRouter()
</script>

<template>
  <div style="min-height: 100vh; background-color: var(--background); display: flex; align-items: center; justify-content: center; padding: 24px">
    <div style="max-width: 480px; width: 100%">
      <Card style="padding: 40px 32px; text-align: center">
        <div style="width: 80px; height: 80px; background-color: rgba(220, 38, 38, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 24px">
          <Clock size="40" color="var(--danger)" />
        </div>
        
        <h1 class="text-h1" style="margin-bottom: 12px">Link phỏng vấn đã hết hạn</h1>
        <p class="text-body" style="color: var(--text-secondary); margin-bottom: 32px; line-height: 1.6">
          Buổi phỏng vấn này đã kết thúc hoặc đường dẫn bạn truy cập không còn hiệu lực. Vui lòng kiểm tra lại thời gian trong email mời.
        </p>

        <div style="display: flex; flex-direction: column; gap: 12px">
          <Button variant="primary" style="width: 100%" @click="() => window.location.href = 'mailto:hr@techcorp.com'">
            <Mail size="18" /> Liên hệ với Bộ phận Nhân sự
          </Button>
          <Button variant="ghost" style="width: 100%" @click="router.push('/home')">
            <ArrowLeft size="18" /> Quay về Trang chủ
          </Button>
        </div>
      </Card>
    </div>
  </div>
</template>
