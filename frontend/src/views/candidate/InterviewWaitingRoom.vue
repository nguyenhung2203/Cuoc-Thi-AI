<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { AlertTriangle, Video, Mic, ShieldCheck } from 'lucide-vue-next'
import { interviewService } from '../../services/interview.service'

const router = useRouter()
const route = useRoute()
const agreed = ref(false)
const loading = ref(true)
const errorMsg = ref('')
const interviewInfo = ref(null)
const inviteToken = ref(route.query.token || '')

onMounted(async () => {
  if (!inviteToken.value) {
    errorMsg.value = 'Link mời phỏng vấn không hợp lệ hoặc bị thiếu token.'
    loading.value = false
    return
  }

  try {
    const data = await interviewService.joinByToken(inviteToken.value)
    interviewInfo.value = data
  } catch (error) {
    errorMsg.value = error.message || 'Không thể xác thực link mời phỏng vấn. Link có thể đã hết hạn.'
  } finally {
    loading.value = false
  }
})

const handleJoin = () => {
  if (!agreed.value) return
  // Pass token and details to candidate room
  router.push({ 
    path: '/candidate-room', 
    query: { token: inviteToken.value },
    state: { message: 'Vào phòng phỏng vấn thành công!', interviewInfo: interviewInfo.value } 
  })
}
</script>

<template>
  <div style="min-height: 100vh; background-color: var(--background); display: flex; align-items: center; justify-content: center; padding: 24px">
    <div style="max-width: 600px; width: 100%">
      <div style="text-align: center; margin-bottom: 32px">
        <div style="width: 64px; height: 64px; background-color: var(--primary); color: white; border-radius: 16px; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px">
          <ShieldCheck size="32" />
        </div>
        <h1 class="text-h1">Chuẩn bị vào phòng phỏng vấn</h1>
        <div v-if="loading" class="text-helper" style="margin-top: 8px">Đang xác thực thông tin...</div>
        <div v-else-if="errorMsg" class="text-helper" style="margin-top: 8px; color: var(--danger)">{{ errorMsg }}</div>
        <p v-else class="text-helper" style="margin-top: 8px; font-size: 15px">Vị trí: {{ interviewInfo?.job_title }} - {{ interviewInfo?.company_name }}</p>
      </div>

      <Card v-if="!loading && !errorMsg" style="padding: 32px">
        <div style="background-color: rgba(217, 119, 6, 0.1); border: 1px solid rgba(217, 119, 6, 0.3); border-radius: var(--radius); padding: 20px; margin-bottom: 24px; display: flex; gap: 16px">
          <AlertTriangle size="24" color="var(--warning)" style="flex-shrink: 0" />
          <div>
            <h3 class="text-body" style="font-weight: 600; color: var(--warning); margin-bottom: 8px">Thông báo về Quyền riêng tư & AI</h3>
            <p class="text-body" style="color: var(--text-secondary); line-height: 1.6">
              Buổi phỏng vấn này sẽ có sự hỗ trợ của Trí tuệ nhân tạo (AI). Nhằm mục đích đánh giá công bằng và cung cấp báo cáo chi tiết cho Nhà tuyển dụng, toàn bộ diễn biến bao gồm:
            </p>
            <ul class="text-body" style="color: var(--text-secondary); margin-top: 8px; padding-left: 24px; line-height: 1.6">
              <li>Dữ liệu âm thanh (Giọng nói) sẽ được ghi lại và chuyển ngữ thành văn bản.</li>
              <li>Dữ liệu văn bản sẽ được AI phân tích tự động.</li>
              <li>Hình ảnh từ Camera sẽ được truyền trực tiếp nhưng KHÔNG lưu trữ.</li>
            </ul>
          </div>
        </div>

        <div style="display: flex; align-items: center; gap: 12px; padding: 16px; background-color: var(--surface-soft); border-radius: var(--radius); margin-bottom: 24px">
          <div style="display: flex; gap: 8px; color: var(--primary)">
            <Video size="20" /> <Mic size="20" />
          </div>
          <p class="text-body">Hệ thống sẽ yêu cầu quyền truy cập Camera và Micro ở bước tiếp theo.</p>
        </div>

        <label style="display: flex; align-items: flex-start; gap: 12px; cursor: pointer; margin-bottom: 32px">
          <input 
            type="checkbox" 
            v-model="agreed"
            style="width: 20px; height: 20px; margin-top: 2px; accent-color: var(--primary)" 
          />
          <span class="text-body" style="font-weight: 500">
            Tôi đã đọc, hiểu rõ và đồng ý với việc sử dụng hệ thống AI phân tích và ghi âm trong buổi phỏng vấn này.
          </span>
        </label>

        <div style="display: flex; gap: 16px">
          <Button variant="secondary" style="flex: 1" @click="router.push('/home')">Từ chối & Quay lại</Button>
          <Button variant="primary" style="flex: 1" :disabled="!agreed" @click="handleJoin">
            Tham gia phỏng vấn
          </Button>
        </div>
      </Card>
      
      <Card v-if="errorMsg" style="padding: 32px; text-align: center">
        <p class="text-body" style="color: var(--danger); margin-bottom: 24px">{{ errorMsg }}</p>
        <Button variant="primary" @click="router.push('/home')">Quay về trang chủ</Button>
      </Card>
    </div>
  </div>
</template>
