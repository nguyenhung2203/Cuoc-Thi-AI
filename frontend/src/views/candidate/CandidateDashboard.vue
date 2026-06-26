<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import WelcomeAlert from '../../components/common/WelcomeAlert.vue'

const router = useRouter()
const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)

onMounted(() => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }
})
</script>

<template>
  <WelcomeAlert 
    v-if="entryToast" 
    role="candidate"
    title="Thành công!"
    :message="entryToast.message" 
    @close="entryToast = null" 
  />
  <div class="page-header" style="margin-bottom: 24px">
    <div>
      <h1 class="text-h1">Chào mừng trở lại, Nguyễn Văn A</h1>
      <p class="text-body" style="color: var(--text-secondary); margin-top: 8px">
        Theo dõi lịch phỏng vấn, luyện tập với AI và cải thiện kỹ năng trả lời.
      </p>
    </div>
  </div>
  
  <div class="grid" style="grid-template-columns: repeat(4, 1fr); margin-top: 24px">
    <Card>
      <div class="text-helper" style="color: var(--text-secondary)">Lịch sắp tới</div>
      <div class="text-h1" style="margin-top: 8px">2</div>
    </Card>
    <Card>
      <div class="text-helper" style="color: var(--text-secondary)">Luyện tập đã xong</div>
      <div class="text-h1" style="margin-top: 8px">5</div>
    </Card>
    <Card>
      <div class="text-helper" style="color: var(--text-secondary)">Điểm trung bình</div>
      <div class="text-h1" style="margin-top: 8px">7.8</div>
    </Card>
    <Card>
      <div class="text-helper" style="color: var(--text-secondary)">Hồ sơ hoàn thiện</div>
      <div class="text-h1" style="margin-top: 8px">80%</div>
    </Card>
  </div>
  
  <div class="grid" style="grid-template-columns: 2fr 1fr; margin-top: 24px; gap: 24px">
    <div style="display: flex; flex-direction: column; gap: 24px">
      <Card>
        <h3 class="text-h3" style="margin-bottom: 16px">Lịch phỏng vấn sắp tới</h3>
        <div style="border: 1px solid var(--border); border-radius: 8px; padding: 16px; margin-bottom: 12px; display: flex; justify-content: space-between; align-items: center">
          <div>
            <div style="font-weight: 600; margin-bottom: 4px">Frontend Developer</div>
            <div class="text-body" style="color: var(--text-secondary)">Công ty TNHH ABC</div>
            <div style="display: flex; gap: 12px; margin-top: 8px">
              <Badge type="primary">14:00 - 15:00, 25/06</Badge>
              <Badge type="default">Phỏng vấn Kỹ thuật</Badge>
            </div>
          </div>
          <div style="display: flex; gap: 8px">
            <Button variant="outline" @click="router.push('/mock-setup')">Luyện tập</Button>
            <Button variant="primary" @click="router.push('/candidate-room')">Tham gia</Button>
          </div>
        </div>
      </Card>
      
      <Card style="background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%); border: 1px solid #bae6fd">
        <h3 class="text-h3" style="margin-bottom: 8px; color: #0369a1">Sẵn sàng luyện phỏng vấn với AI?</h3>
        <p class="text-body" style="margin-bottom: 16px; color: #0c4a6e">Trải nghiệm phỏng vấn thực tế với AI Interviewer của chúng tôi để chuẩn bị tốt nhất.</p>
        <Button variant="primary" @click="router.push('/mock-setup')">Bắt đầu luyện tập</Button>
      </Card>
    </div>
    
    <div style="display: flex; flex-direction: column; gap: 24px">
      <Card>
        <h3 class="text-h3" style="margin-bottom: 16px">Trạng thái hồ sơ</h3>
        <div style="display: flex; flex-direction: column; gap: 12px">
          <div style="display: flex; justify-content: space-between">
            <span class="text-body">Tải lên CV</span>
            <Badge type="success">Đã hoàn thành</Badge>
          </div>
          <div style="display: flex; justify-content: space-between">
            <span class="text-body">Thêm Kỹ năng</span>
            <Badge type="success">Đã hoàn thành</Badge>
          </div>
          <div style="display: flex; justify-content: space-between">
            <span class="text-body">Kinh nghiệm</span>
            <Badge type="warning">Cần cập nhật</Badge>
          </div>
        </div>
      </Card>
      
      <Card>
        <h3 class="text-h3" style="margin-bottom: 16px; display: flex; align-items: center; gap: 8px">
          <span style="color: var(--primary)">✨</span> AI Career Coach
        </h3>
        <p class="text-body" style="line-height: 1.6">
          Dựa trên kết quả phỏng vấn gần đây, bạn nên luyện tập thêm các câu hỏi về <strong style="color: var(--primary)">React Performance Optimization</strong>.
        </p>
        <Button variant="outline" style="width: 100%; margin-top: 16px" @click="router.push('/mock-setup')">
          Luyện chủ đề này
        </Button>
      </Card>
    </div>
  </div>
</template>
