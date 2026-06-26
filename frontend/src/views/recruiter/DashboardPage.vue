<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import WelcomeAlert from '../../components/common/WelcomeAlert.vue'
import { Users, Briefcase, Calendar, Plus, FileText, Sparkles, TrendingUp } from 'lucide-vue-next'

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
    role="recruiter"
    title="Đăng nhập thành công!"
    :message="entryToast.message" 
    @close="entryToast = null" 
  />
  <div class="page-header">
    <div>
      <h1 class="text-h1">Tổng quan</h1>
      <p class="text-body" style="color: var(--text-secondary); margin-top: 8px">
        Theo dõi lịch phỏng vấn, ứng viên và báo cáo AI hôm nay.
      </p>
    </div>
    <div style="display: flex; gap: 12px">
      <Button @click="router.push('/interviews/new')"><Plus size="16" /> Tạo lịch phỏng vấn</Button>
    </div>
  </div>
  
  <div class="grid" style="grid-template-columns: repeat(4, 1fr); margin-top: 24px; margin-bottom: 24px">
    <Card>
      <div class="text-helper" style="display: flex; align-items: center; gap: 8px; color: var(--text-secondary)">
        <Briefcase size="16" color="var(--primary)" /> Jobs đang mở
      </div>
      <div class="text-h1" style="margin-top: 8px">12</div>
    </Card>
    <Card>
      <div class="text-helper" style="display: flex; align-items: center; gap: 8px; color: var(--text-secondary)">
        <Users size="16" color="var(--accent)" /> Ứng viên mới
      </div>
      <div class="text-h1" style="margin-top: 8px">48</div>
    </Card>
    <Card>
      <div class="text-helper" style="display: flex; align-items: center; gap: 8px; color: var(--text-secondary)">
        <Calendar size="16" color="var(--warning)" /> Phỏng vấn hôm nay
      </div>
      <div class="text-h1" style="margin-top: 8px">5</div>
    </Card>
    <Card>
      <div class="text-helper" style="display: flex; align-items: center; gap: 8px; color: var(--text-secondary)">
        <FileText size="16" color="var(--success)" /> Báo cáo chờ xem
      </div>
      <div class="text-h1" style="margin-top: 8px">3</div>
    </Card>
  </div>

  <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px">
    <div style="display: flex; flex-direction: column; gap: 24px">
      <Card title="Lịch phỏng vấn sắp tới">
        <div style="padding: 32px; text-align: center">
          <Calendar size="32" color="var(--text-muted)" style="margin: 0 auto 16px" />
          <p class="text-body" style="color: var(--text-muted)">Chưa có lịch phỏng vấn nào sắp tới.</p>
        </div>
      </Card>
      <Card title="Ứng viên mới nhất">
        <div style="padding: 32px; text-align: center">
          <Users size="32" color="var(--text-muted)" style="margin: 0 auto 16px" />
          <p class="text-body" style="color: var(--text-muted)">Chưa có ứng viên mới ứng tuyển.</p>
        </div>
      </Card>
    </div>

    <div style="display: flex; flex-direction: column; gap: 24px">
      <Card title="AI Insights" style="border-top: 4px solid var(--accent)">
        <div style="display: flex; flex-direction: column; gap: 16px">
          <div style="display: flex; gap: 12px; align-items: flex-start">
            <Sparkles size="16" color="var(--accent)" style="margin-top: 2px; flex-shrink: 0" />
            <div>
              <p class="text-body" style="font-weight: 500">3 ứng viên có mức phù hợp cao</p>
              <p class="text-helper">Vừa nộp đơn vào vị trí Frontend Developer.</p>
            </div>
          </div>
          <div style="display: flex; gap: 12px; align-items: flex-start">
            <TrendingUp size="16" color="var(--warning)" style="margin-top: 2px; flex-shrink: 0" />
            <div>
              <p class="text-body" style="font-weight: 500">2 buổi phỏng vấn cần review</p>
              <p class="text-helper">Báo cáo AI đã sẵn sàng để bạn đưa ra quyết định.</p>
            </div>
          </div>
        </div>
      </Card>
      
      <Card title="Thao tác nhanh">
        <div style="display: flex; flex-direction: column; gap: 8px">
          <Button variant="secondary" style="justify-content: flex-start" @click="router.push('/jobs/new')"><Plus size="16" /> Tạo job mới</Button>
          <Button variant="secondary" style="justify-content: flex-start" @click="router.push('/candidates/new')"><Plus size="16" /> Thêm ứng viên</Button>
        </div>
      </Card>
    </div>
  </div>
</template>
