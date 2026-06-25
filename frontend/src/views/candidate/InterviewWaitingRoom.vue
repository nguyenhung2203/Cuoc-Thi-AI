<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import { AlertTriangle, Video, Mic, ShieldCheck } from 'lucide-vue-next'

const router = useRouter()
const agreed = ref(false)
</script>

<template>
  <div style="min-height: 100vh; background-color: var(--background); display: flex; align-items: center; justify-content: center; padding: 24px">
    <div style="max-width: 600px; width: 100%">
      <div style="text-align: center; margin-bottom: 32px">
        <div style="width: 64px; height: 64px; background-color: var(--primary); color: white; border-radius: 16px; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px">
          <ShieldCheck size="32" />
        </div>
        <h1 class="text-h1">Chuẩn bị vào phòng phỏng vấn</h1>
        <p class="text-helper" style="margin-top: 8px; font-size: 15px">Vị trí: Frontend Developer - TechCorp Inc.</p>
      </div>

      <Card style="padding: 32px">
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
          <Button variant="primary" style="flex: 1" :disabled="!agreed" @click="router.push({ path: '/candidate-room', state: { message: 'Vào phòng phỏng vấn thành công!' } })">
            Tham gia phỏng vấn
          </Button>
        </div>
      </Card>
    </div>
  </div>
</template>
