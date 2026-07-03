<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import WelcomeAlert from '../../components/common/WelcomeAlert.vue'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { authStore } from '../../stores/auth.store'

const router = useRouter()
const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)

const stats = ref({
  upcoming_interviews: 0,
  completed_mock_tests: 0,
  average_mock_score: 0,
  profile_completeness: 0
})
const upcomingInterviews = ref([])
const loading = ref(true)

onMounted(async () => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }

  try {
    const [statsData, interviewsData] = await Promise.all([
      candidatePortalService.getDashboardStats(),
      candidatePortalService.getInterviews()
    ])
    
    stats.value = statsData
    // Filter only future interviews or recently active ones
    upcomingInterviews.value = interviewsData.filter(i => i.status === 'scheduled' || i.status === 'active').slice(0, 3)
  } catch (err) {
    console.error('Lỗi tải dữ liệu dashboard:', err)
  } finally {
    loading.value = false
  }
})

const formatDate = (dateString) => {
  if (!dateString) return ''
  const d = new Date(dateString)
  return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' }) + ', ' + d.toLocaleDateString('vi-VN')
}
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
      <h1 class="text-h1">Chào mừng trở lại, {{ authStore.user?.full_name || 'Ứng viên' }}</h1>
      <p class="text-body" style="color: var(--text-secondary); margin-top: 8px">
        Theo dõi lịch phỏng vấn, luyện tập với AI và cải thiện kỹ năng trả lời.
      </p>
    </div>
  </div>
  
  <div v-if="loading" style="padding: 48px; text-align: center; color: var(--text-muted)">
    Đang tải dữ liệu...
  </div>
  
  <template v-else>
    <div class="grid" style="grid-template-columns: repeat(4, 1fr); margin-top: 24px">
      <Card>
        <div class="text-helper" style="color: var(--text-secondary)">Lịch sắp tới</div>
        <div class="text-h1" style="margin-top: 8px">{{ stats.upcoming_interviews }}</div>
      </Card>
      <Card>
        <div class="text-helper" style="color: var(--text-secondary)">Luyện tập đã xong</div>
        <div class="text-h1" style="margin-top: 8px">{{ stats.completed_mock_tests }}</div>
      </Card>
      <Card>
        <div class="text-helper" style="color: var(--text-secondary)">Điểm trung bình</div>
        <div class="text-h1" style="margin-top: 8px">{{ stats.average_mock_score.toFixed(1) }}</div>
      </Card>
      <Card>
        <div class="text-helper" style="color: var(--text-secondary)">Hồ sơ hoàn thiện</div>
        <div class="text-h1" style="margin-top: 8px">{{ stats.profile_completeness }}%</div>
      </Card>
    </div>
    
    <div class="grid" style="grid-template-columns: 2fr 1fr; margin-top: 24px; gap: 24px">
      <div style="display: flex; flex-direction: column; gap: 24px">
        <Card>
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px">
            <h3 class="text-h3">Lịch phỏng vấn sắp tới</h3>
            <span class="text-helper" style="color: var(--primary); cursor: pointer; font-weight: 500" @click="router.push('/my-interviews')">Xem tất cả</span>
          </div>
          
          <div v-if="upcomingInterviews.length === 0" style="padding: 24px; text-align: center; border: 1px dashed var(--border); border-radius: 8px; color: var(--text-secondary)">
            Bạn chưa có lịch phỏng vấn nào sắp tới.
          </div>
          
          <div v-for="iv in upcomingInterviews" :key="iv.id" style="border: 1px solid var(--border); border-radius: 8px; padding: 16px; margin-bottom: 12px; display: flex; justify-content: space-between; align-items: center">
            <div>
              <div style="font-weight: 600; margin-bottom: 4px">{{ iv.job_title || iv.title }}</div>
              <div class="text-body" style="color: var(--text-secondary)">{{ iv.company_name || 'Công ty ẩn danh' }}</div>
              <div style="display: flex; gap: 12px; margin-top: 8px">
                <Badge type="primary">{{ formatDate(iv.scheduled_at) }}</Badge>
                <Badge :type="iv.mode === 'real' ? 'danger' : 'default'">{{ iv.mode === 'real' ? 'Phỏng vấn thật' : 'Phỏng vấn thử' }}</Badge>
              </div>
            </div>
            <div style="display: flex; gap: 8px">
              <Button v-if="iv.mode === 'real'" variant="primary" @click="iv.join_link ? router.push(iv.join_link) : null">Tham gia</Button>
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
              <Badge :type="stats.profile_completeness > 50 ? 'success' : 'warning'">{{ stats.profile_completeness > 50 ? 'Đã hoàn thành' : 'Cần cập nhật' }}</Badge>
            </div>
            <div style="display: flex; justify-content: space-between">
              <span class="text-body">Thêm Kỹ năng</span>
              <Badge type="warning">Cần cập nhật</Badge>
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
            Dựa trên kết quả phỏng vấn gần đây, bạn nên luyện tập thêm các câu hỏi về <strong style="color: var(--primary)">Kỹ năng chuyên môn</strong>.
          </p>
          <Button variant="outline" style="width: 100%; margin-top: 16px" @click="router.push('/mock-setup')">
            Luyện chủ đề này
          </Button>
        </Card>
      </div>
    </div>
  </template>
</template>
