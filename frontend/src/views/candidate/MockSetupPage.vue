<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import { Play, FileText, Sparkles } from 'lucide-vue-next'

const router = useRouter()
const loading = ref(false)

const handleStart = (e) => {
  e.preventDefault()
  loading.value = true
  setTimeout(() => {
    router.push({ path: '/mock-room', state: { message: 'Bắt đầu phiên luyện tập thành công!' } })
  }, 800)
}
</script>

<template>
  <div style="max-width: 800px; margin: 0 auto">
    <div style="margin-bottom: 32px">
      <h1 class="text-h1">Luyện phỏng vấn cùng AI</h1>
      <p class="text-helper" style="margin-top: 4px">AI sẽ đóng vai người phỏng vấn thật, đặt câu hỏi theo vị trí ứng tuyển và đưa feedback chi tiết.</p>
    </div>

    <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px">
      <Card title="Cấu hình buổi phỏng vấn">
        <form @submit="handleStart">
          <div class="input-group">
            <label class="input-label">Vị trí ứng tuyển (Target Role)</label>
            <select class="input-field" required>
              <option value="frontend">Frontend Developer</option>
              <option value="backend">Backend Developer</option>
              <option value="fullstack">Fullstack Developer</option>
              <option value="pm">Product Manager</option>
            </select>
          </div>

          <div class="input-group">
            <label class="input-label">Cấp độ (Level)</label>
            <select class="input-field" required>
              <option value="fresher">Fresher</option>
              <option value="junior">Junior</option>
              <option value="middle">Middle</option>
              <option value="senior">Senior</option>
            </select>
          </div>

          <div class="input-group">
            <label class="input-label">Loại phỏng vấn</label>
            <select class="input-field" required>
              <option value="tech">Phỏng vấn Kỹ thuật (Technical)</option>
              <option value="behavior">Phỏng vấn Hành vi (Behavioral)</option>
              <option value="hr">Phỏng vấn Nhân sự (HR)</option>
            </select>
          </div>

          <div class="input-group">
            <label class="input-label">Phong cách AI (Interviewer Style)</label>
            <select class="input-field" required>
              <option value="friendly">Thân thiện, gợi mở</option>
              <option value="professional">Chuyên nghiệp, tiêu chuẩn</option>
              <option value="challenging">Khó tính, hay hỏi xoáy</option>
            </select>
          </div>

          <div style="margin-top: 24px; padding-top: 24px; border-top: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center">
            <span class="text-helper" style="display: flex; align-items: center; gap: 8px">
              <FileText size="16" /> Sử dụng CV hiện tại trong Hồ sơ
            </span>
            <Button type="submit" size="large" :disabled="loading">
              <Play size="16" /> {{ loading ? 'Đang khởi tạo...' : 'Bắt đầu ngay' }}
            </Button>
          </div>
        </form>
      </Card>

      <div style="display: flex; flex-direction: column; gap: 24px">
        <Card style="background-color: var(--surface-soft)">
          <div style="text-align: center; padding: 16px 0">
            <div style="width: 80px; height: 80px; border-radius: 50%; background-color: var(--accent); margin: 0 auto 16px; display: flex; align-items: center; justify-content: center; color: white">
              <Sparkles size="32" />
            </div>
            <h3 class="text-h2">AI Interviewer</h3>
            <p class="text-helper" style="margin-top: 8px">Sẵn sàng hỗ trợ bạn</p>
          </div>
        </Card>
        
        <Card title="Lưu ý">
          <ul class="text-helper" style="padding-left: 20px; display: flex; flex-direction: column; gap: 8px">
            <li>Chuẩn bị micro và không không gian yên tĩnh.</li>
            <li>Sẽ có khoảng 3-5 câu hỏi tùy thuộc vào chức danh.</li>
            <li>Bạn có thể trả lời bằng Giọng nói hoặc Văn bản.</li>
          </ul>
        </Card>
      </div>
    </div>
  </div>
</template>
