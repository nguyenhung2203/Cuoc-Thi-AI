<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { onMounted } from 'vue'
import { Play, FileText, Settings2, Bot, Info, CheckCircle2 } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'
import { candidatePortalService } from '../../services/candidate-portal.service'

const router = useRouter()
const loading = ref(false)
const toast = ref(null)
const cvFileId = ref('')

const setup = ref({
  jobRole: 'frontend',
  level: 'junior',
  type: 'tech',
  style: 'friendly',
  useCurrentCv: false
})

// Lấy CV file id thật từ hồ sơ candidate để dùng cho phiên luyện tập
onMounted(async () => {
  try {
    const res = await candidatePortalService.getProfile()
    const profile = res?.data || res || {}
    cvFileId.value = profile.cv_file_id || ''
  } catch (err) {
    cvFileId.value = ''
  }
})

const handleStart = async (e) => {
  e.preventDefault()
  if (setup.value.useCurrentCv && !cvFileId.value) {
    toast.value = { type: 'warning', message: 'Bạn chưa có CV trong hồ sơ. Vui lòng tải CV lên trước hoặc bỏ chọn tùy chọn này.' }
    return
  }
  loading.value = true
  try {
    // Bước 1: Tạo session mới (status=draft) — API_SPEC §11.1
    const session = await mockService.createMockInterview({
      target_role: setup.value.jobRole,
      target_level: setup.value.level,
      cv_file_id: setup.value.useCurrentCv ? cvFileId.value : undefined
    })
    // Bước 2: Bắt đầu session → nhận câu hỏi đầu tiên — API_SPEC §11.2
    await mockService.startMockInterview(session.id)
    // Chuyển vào phòng phỏng vấn mock (kèm role/level cho phiên thoại AI)
    router.push({
      path: '/mock-room',
      query: {
        mock_id: session.id,
        role: setup.value.jobRole,
        level: setup.value.level,
      },
    })
  } catch (error) {
    console.error(error)
    toast.value = { type: 'error', message: 'Không thể khởi động phỏng vấn. Backend đang được kết nối.' }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="space-y-8 pb-12 max-w-5xl mx-auto">
    <div class="animate-rise">
      <h1 class="page-title">Luyện phỏng vấn cùng AI</h1>
      <p class="page-subtitle">AI sẽ đóng vai người phỏng vấn thật, đặt câu hỏi theo vị trí ứng tuyển và đưa feedback chi tiết.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Form Setup -->
      <div class="lg:col-span-2 space-y-6">
        <Card class="card-elevate setup-card animate-rise">
          <h2 class="section-heading mb-7">
            <span class="kpi-icon"><Settings2 :size="18" /></span>
            Cấu hình buổi phỏng vấn
          </h2>

          <form @submit="handleStart" class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div class="field">
                <label class="field-label">Vị trí ứng tuyển</label>
                <select class="field-select" required v-model="setup.jobRole">
                  <option value="frontend">Frontend Developer</option>
                  <option value="backend">Backend Developer</option>
                  <option value="fullstack">Fullstack Developer</option>
                  <option value="pm">Product Manager</option>
                </select>
              </div>

              <div class="field">
                <label class="field-label">Cấp độ</label>
                <select class="field-select" required v-model="setup.level">
                  <option value="fresher">Fresher</option>
                  <option value="junior">Junior</option>
                  <option value="middle">Middle</option>
                  <option value="senior">Senior</option>
                </select>
              </div>

              <div class="field">
                <label class="field-label">Loại phỏng vấn</label>
                <select class="field-select" required v-model="setup.type">
                  <option value="tech">Phỏng vấn Kỹ thuật (Technical)</option>
                  <option value="behavior">Phỏng vấn Hành vi (Behavioral)</option>
                  <option value="hr">Phỏng vấn Nhân sự (HR)</option>
                </select>
              </div>

              <div class="field">
                <label class="field-label">Phong cách AI</label>
                <select class="field-select" required v-model="setup.style">
                  <option value="friendly">Thân thiện, gợi mở</option>
                  <option value="professional">Chuyên nghiệp, tiêu chuẩn</option>
                  <option value="challenging">Khó tính, hay hỏi xoáy</option>
                </select>
              </div>
            </div>

            <div class="setup-foot">
              <label class="cv-toggle">
                <input type="checkbox" v-model="setup.useCurrentCv" :disabled="!cvFileId" />
                <FileText :size="18" />
                <span>Sử dụng CV hiện tại trong Hồ sơ</span>
                <span v-if="!cvFileId" class="cv-hint">(chưa có CV)</span>
              </label>

              <button type="submit" :disabled="loading" class="btn btn-primary sheen start-btn">
                <Play :size="18" />
                {{ loading ? 'Đang khởi tạo phòng...' : 'Bắt đầu ngay' }}
              </button>
            </div>
          </form>
        </Card>
      </div>

      <!-- Right Column -->
      <div class="lg:col-span-1 space-y-6">
        <!-- AI Interviewer Card -->
        <div class="ai-tile animate-rise">
          <div class="ai-avatar">
            <Bot :size="44" />
            <span class="ai-ring"></span>
          </div>
          <h3 class="ai-name">AI Interviewer</h3>
          <p class="ai-role">Sẵn sàng hỗ trợ bạn</p>
          <div class="ai-status">
            <span class="ai-dot"><span class="ai-dot-ping"></span></span>
            Trực tuyến
          </div>
        </div>

        <!-- Checklist -->
        <Card class="card-elevate setup-card animate-rise">
          <h3 class="section-heading mb-4">
            <span class="kpi-icon is-warning"><Info :size="18" /></span>
            Lưu ý trước khi bắt đầu
          </h3>
          <ul class="tips">
            <li><CheckCircle2 :size="18" /> Chuẩn bị micro và ngồi trong không gian yên tĩnh.</li>
            <li><CheckCircle2 :size="18" /> Sẽ có khoảng 3-5 câu hỏi tùy thuộc vào chức danh.</li>
            <li><CheckCircle2 :size="18" /> Bạn có thể trả lời bằng Giọng nói (khuyên dùng) hoặc Văn bản.</li>
          </ul>
        </Card>
      </div>
    </div>
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
  </div>
</template>

<style scoped>
.setup-card { padding: 28px; }
.field { display: flex; flex-direction: column; gap: 6px; }
.field-label { font-size: 13px; font-weight: 600; color: var(--text-secondary); }
.field-select {
  width: 100%; height: 44px; padding: 0 14px;
  border: 1px solid var(--border); border-radius: var(--radius);
  background: var(--surface-soft); color: var(--text-main); font-size: 14px; font-weight: 500;
  outline: none; transition: all 0.2s ease; cursor: pointer; font-family: var(--sans);
}
.field-select:focus { border-color: var(--primary); background: var(--surface); box-shadow: 0 0 0 3px rgba(37,99,235,0.12); }

.setup-foot { display: flex; flex-direction: column; gap: 18px; margin-top: 24px; padding-top: 24px; border-top: 1px solid var(--border); }
@media (min-width: 768px) { .setup-foot { flex-direction: row; align-items: center; justify-content: space-between; } }
.cv-toggle { display: flex; align-items: center; gap: 10px; padding: 12px 16px; border-radius: var(--radius); background: var(--primary-light); border: 1px solid rgba(37,99,235,0.15); color: var(--primary); font-weight: 500; font-size: 14px; cursor: pointer; }
.cv-toggle input { width: 18px; height: 18px; accent-color: var(--primary); cursor: pointer; }
.cv-hint { color: var(--text-muted); font-size: 12px; font-weight: 400; }
.start-btn { height: 46px; padding: 0 28px; }

/* AI tile — focal dark presence, brand blue/cyan (no purple) */
.ai-tile {
  position: relative; overflow: hidden;
  padding: 32px 24px; text-align: center;
  border-radius: var(--radius-lg);
  background: radial-gradient(120% 120% at 50% 0%, #0f2f5c 0%, #0b1b34 60%, #0a1526 100%);
  box-shadow: var(--shadow-lg);
}
.ai-avatar { position: relative; width: 88px; height: 88px; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; border-radius: 50%; background: var(--gradient-brand); color: #fff; box-shadow: 0 0 36px rgba(37,99,235,0.5); }
.ai-ring { position: absolute; inset: -6px; border-radius: 50%; border: 2px solid rgba(56,189,248,0.35); animation: aiPulse 2.4s ease-out infinite; }
@keyframes aiPulse { 0% { transform: scale(1); opacity: 0.8; } 100% { transform: scale(1.25); opacity: 0; } }
.ai-name { font-size: 20px; font-weight: 700; color: #fff; margin-bottom: 4px; }
.ai-role { color: #93c5fd; font-weight: 500; margin-bottom: 20px; }
.ai-status { display: inline-flex; align-items: center; gap: 8px; padding: 7px 16px; border-radius: var(--radius-full); background: rgba(255,255,255,0.08); border: 1px solid rgba(255,255,255,0.14); color: #6ee7b7; font-size: 13px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; }
.ai-dot { position: relative; width: 10px; height: 10px; border-radius: 50%; background: #10b981; }
.ai-dot-ping { position: absolute; inset: 0; border-radius: 50%; background: #10b981; animation: aiPulse 1.8s ease-out infinite; }

.tips { display: flex; flex-direction: column; gap: 14px; }
.tips li { display: flex; align-items: flex-start; gap: 10px; color: var(--text-secondary); font-size: 14px; font-weight: 500; }
.tips li :deep(svg) { color: var(--success); flex-shrink: 0; margin-top: 1px; }

@media (prefers-reduced-motion: reduce) {
  .ai-ring, .ai-dot-ping { animation: none !important; }
}
</style>
