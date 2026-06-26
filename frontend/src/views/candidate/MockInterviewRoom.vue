<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from '../../components/common/AppButton.vue'
import Card from '../../components/common/AppCard.vue'
import Badge from '../../components/common/AppBadge.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mic, ArrowRight, Play, Square, RefreshCw, Send, CheckCircle } from 'lucide-vue-next'

const router = useRouter()
const recording = ref(false)
const questionIndex = ref(1)
const analyzing = ref(false)
const feedback = ref(null)
const textAnswer = ref('')
const showEndModal = ref(false)

const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)

onMounted(() => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }
})

const questions = [
  "Hãy giới thiệu ngắn gọn về bản thân và kinh nghiệm làm việc của bạn.",
  "Bạn đã từng sử dụng React trong dự án nào? Hãy mô tả một thử thách lớn nhất bạn gặp phải.",
  "Làm thế nào để bạn quản lý state trong một ứng dụng React lớn?"
]

const handleRecord = () => {
  if (!recording.value) {
    recording.value = true
    feedback.value = null
  } else {
    recording.value = false
    analyzing.value = true
    setTimeout(() => {
      analyzing.value = false
      feedback.value = {
        score: 'Tốt',
        message: 'Câu trả lời rõ ràng, cấu trúc tốt. Đã đề cập được số năm kinh nghiệm và công nghệ chính.',
        improvement: 'Có thể thêm một ví dụ ngắn về dự án gần nhất để tăng tính thuyết phục.'
      }
    }, 2000)
  }
}

const handleSendText = () => {
  if (!textAnswer.value) return
  analyzing.value = true
  feedback.value = null
  setTimeout(() => {
    analyzing.value = false
    feedback.value = {
      score: 'Khá',
      message: 'Bạn đã nêu được các ý chính, câu văn mạch lạc.',
      improvement: 'Cố gắng trả lời bằng giọng nói để rèn luyện sự tự tin tốt hơn nhé!'
    }
  }, 1500)
}

const handleNext = () => {
  if (questionIndex.value < questions.length) {
    questionIndex.value++
    feedback.value = null
    textAnswer.value = ''
  } else {
    router.push({ path: '/mock-results', state: { message: 'Hoàn thành bài thi thử!' } })
  }
}
</script>

<template>
  <div style="height: 100vh; display: flex; flex-direction: column; background-color: var(--background)">
    <!-- Header -->
    <div style="height: 64px; background-color: var(--surface); border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 24px">
      <div style="display: flex; align-items: center; gap: 16px">
        <div>
          <h1 class="text-h2">Luyện tập AI: Frontend Developer</h1>
          <p class="text-helper" style="margin-top: 4px">Câu hỏi {{ questionIndex }} / {{ questions.length }}</p>
        </div>
      </div>
      <Button variant="ghost" style="color: var(--danger)" @click="showEndModal = true">Kết thúc sớm</Button>
    </div>

    <Modal :isOpen="showEndModal" @close="showEndModal = false" title="Kết thúc sớm">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn kết thúc bài thi sớm? Kết quả sẽ được tính trên những câu bạn đã trả lời.</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showEndModal = false">Huỷ</Button>
        <Button variant="primary" style="background-color: var(--danger); color: white; border-color: var(--danger)" @click="router.push({ path: '/home', state: { message: 'Đã hủy bài thi thử' } })">Xác nhận</Button>
      </div>
    </Modal>

    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />

    <!-- Main Content -->
    <div style="display: flex; flex: 1; overflow: hidden">
      
      <!-- Left: AI Interviewer -->
      <div style="flex: 4; display: flex; flex-direction: column; padding: 24px; gap: 24px; border-right: 1px solid var(--border); background-color: var(--surface)">
        <!-- AI Avatar Area -->
        <div style="flex: 1; background-color: var(--surface-soft); border-radius: var(--radius-lg); display: flex; flex-direction: column; align-items: center; justify-content: center; border: 1px solid var(--border); position: relative">
          <div style="width: 120px; height: 120px; border-radius: 50%; background-color: var(--accent); display: flex; align-items: center; justify-content: center; color: white; margin-bottom: 24px; box-shadow: 0 0 0 12px rgba(8, 145, 178, 0.1)">
            <span style="font-size: 32px; font-weight: bold">AI</span>
          </div>
          <h3 class="text-h2" style="margin-bottom: 8px">AI Interviewer</h3>
          
          <div :style="{ padding: '8px 16px', backgroundColor: analyzing ? 'rgba(8, 145, 178, 0.1)' : 'var(--surface)', borderRadius: '20px', border: '1px solid var(--border)', display: 'flex', alignItems: 'center', gap: '8px' }">
            <template v-if="analyzing">
              <RefreshCw size="14" class="spin" color="var(--accent)" /> <span class="text-helper" style="color: var(--accent)">Đang phân tích câu trả lời...</span>
            </template>
            <template v-else>
              <Play size="14" color="var(--success)" /> <span class="text-helper">Sẵn sàng lắng nghe</span>
            </template>
          </div>
        </div>

        <!-- Current Question -->
        <Card style="background-color: var(--surface)">
          <div style="display: flex; justify-content: space-between; margin-bottom: 12px">
            <Badge type="info">Câu hỏi hiện tại</Badge>
            <Button variant="ghost" style="padding: 4px 8px; height: auto"><Play size="14" style="margin-right: 4px"/> Nghe lại</Button>
          </div>
          <p class="text-body" style="font-size: 16px; line-height: 1.6; font-weight: 500">
            {{ questions[questionIndex - 1] }}
          </p>
        </Card>
      </div>

      <!-- Right: Candidate Answer -->
      <div style="flex: 6; display: flex; flex-direction: column; padding: 24px; gap: 24px">
        
        <div style="display: flex; flex-direction: column; gap: 16px; flex: 1">
          <!-- Audio Recording UI -->
          <Card style="display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 48px 24px; flex: 1">
            <div 
              :style="{ 
                width: '80px', height: '80px', borderRadius: '50%', 
                backgroundColor: recording ? 'rgba(239, 68, 68, 0.1)' : 'rgba(37, 99, 235, 0.1)', 
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                cursor: 'pointer', transition: 'all 0.3s',
                border: `2px solid ${recording ? 'var(--danger)' : 'var(--primary)'}`
              }"
              @click="handleRecord"
            >
              <Square v-if="recording" size="32" color="var(--danger)" />
              <Mic v-else size="32" color="var(--primary)" />
            </div>
            <h3 class="text-body" style="margin-top: 24px; font-weight: 600">
              {{ recording ? 'Đang thu âm...' : 'Nhấn để trả lời bằng giọng nói' }}
            </h3>
            <p class="text-helper" style="margin-top: 8px">
              {{ recording ? '00:15' : 'Tối đa 3 phút cho mỗi câu trả lời' }}
            </p>

            <!-- Fake Waveform when recording -->
            <div v-if="recording" style="display: flex; gap: 4px; margin-top: 24px; height: 32px; align-items: center">
              <div v-for="(h, i) in [1, 2, 3, 4, 5, 4, 3, 2, 1, 2, 4, 5, 3, 2]" :key="i" :style="{ width: '4px', height: `${h * 20}%`, backgroundColor: 'var(--danger)', borderRadius: '2px', animation: `pulse ${0.5 + (i%3)*0.1}s infinite alternate` }"></div>
            </div>
          </Card>

          <!-- OR Text input -->
          <div style="position: relative">
            <textarea 
              class="input-field" 
              placeholder="Hoặc nhập câu trả lời bằng văn bản tại đây..."
              style="width: 100%; height: 120px"
              v-model="textAnswer"
              :disabled="recording || analyzing"
            ></textarea>
            <Button 
              style="position: absolute; bottom: 16px; right: 16px"
              :disabled="!textAnswer || recording || analyzing"
              @click="handleSendText"
            >
              <Send size="14" style="margin-right: 8px" /> Gửi
            </Button>
          </div>
        </div>

        <!-- Feedback Area -->
        <Card v-if="feedback" style="border-top: 4px solid var(--accent); background-color: var(--surface); animation: fadeIn 0.5s ease-out">
          <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 16px">
            <h3 class="text-h2" style="display: flex; align-items: center; gap: 8px">
              <CheckCircle size="20" color="var(--success)" /> Phản hồi từ AI
            </h3>
            <Badge type="success">Điểm: {{ feedback.score }}</Badge>
          </div>
          <p class="text-body" style="margin-bottom: 12px">{{ feedback.message }}</p>
          <div style="padding: 12px; background-color: rgba(217, 119, 6, 0.05); border-radius: var(--radius); border-left: 3px solid var(--warning)">
            <p class="text-helper" style="font-weight: 600; color: var(--warning); margin-bottom: 4px">Gợi ý cải thiện:</p>
            <p class="text-body" style="color: var(--text-secondary)">{{ feedback.improvement }}</p>
          </div>

          <div style="display: flex; justify-content: flex-end; margin-top: 24px">
            <Button @click="handleNext">
              {{ questionIndex < questions.length ? 'Câu hỏi tiếp theo' : 'Xem báo cáo tổng hợp' }} <ArrowRight size="16" />
            </Button>
          </div>
        </Card>

      </div>
    </div>
  </div>
</template>

<style>
@keyframes pulse { 0% { height: 20%; } 100% { height: 100%; } }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
.spin { animation: spin 2s linear infinite; }
@keyframes spin { 100% { transform: rotate(360deg); } }
</style>
