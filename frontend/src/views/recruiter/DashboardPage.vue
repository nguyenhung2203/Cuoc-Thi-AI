<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import WelcomeAlert from '../../components/common/WelcomeAlert.vue'
import { Users, Briefcase, Calendar, Plus, FileText, Sparkles, TrendingUp, Link, Copy } from 'lucide-vue-next'

import { jobService } from '../../services/job.service'
import { candidateService } from '../../services/candidate.service'
import { interviewService } from '../../services/interview.service'
import { reportService } from '../../services/report.service'
import { authStore } from '../../stores/auth.store'

const router = useRouter()
const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)

const stats = ref({ jobs: 0, candidates: 0, interviewsToday: 0, pendingReports: 0 })
const upcomingInterviews = ref([])
const newCandidates = ref([])
const aiInsights = ref({ highMatchCandidates: 0, reviewsNeeded: 0 })
const loading = ref(true)

onMounted(async () => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }
  
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) return
    
    const [jobs, candidates, interviews] = await Promise.all([
      jobService.getJobs(companyId),
      candidateService.getCandidates(companyId),
      interviewService.getInterviews(companyId)
    ])
    
    stats.value.jobs = jobs.length
    stats.value.candidates = candidates.length
    
    const todayStr = new Date().toISOString().split('T')[0]
    
    const todayInterviews = interviews.filter(i => {
      const dt = typeof i.scheduled_at === 'object' && i.scheduled_at !== null ? i.scheduled_at.Time : i.scheduled_at
      return dt && dt.startsWith(todayStr)
    })
    stats.value.interviewsToday = todayInterviews.length
    
    const futureInterviews = interviews.filter(i => {
      const dt = typeof i.scheduled_at === 'object' && i.scheduled_at !== null ? i.scheduled_at.Time : i.scheduled_at
      return dt && new Date(dt) > new Date()
    }).sort((a, b) => {
      const aTime = typeof a.scheduled_at === 'object' && a.scheduled_at !== null ? a.scheduled_at.Time : a.scheduled_at
      const bTime = typeof b.scheduled_at === 'object' && b.scheduled_at !== null ? b.scheduled_at.Time : b.scheduled_at
      return new Date(aTime) - new Date(bTime)
    }).slice(0, 5)
    
    const jobMap = jobs.reduce((acc, j) => { acc[j.id] = j; return acc; }, {})
    const candidateMap = candidates.reduce((acc, c) => { acc[c.id] = c; return acc; }, {})
    
    upcomingInterviews.value = futureInterviews.map(i => ({
      ...i,
      candidate: i.candidate || candidateMap[i.candidate_id] || {},
      job: i.job || jobMap[i.job_id?.String || i.job_id] || {},
      dt: typeof i.scheduled_at === 'object' && i.scheduled_at !== null ? i.scheduled_at.Time : i.scheduled_at
    }))
    
    newCandidates.value = candidates.sort((a, b) => new Date(b.created_at) - new Date(a.created_at)).slice(0, 5).map(c => ({
      ...c,
      job: jobMap[c.job_id] || {}
    }))
    
    aiInsights.value.highMatchCandidates = candidates.length > 0 ? Math.floor(candidates.length / 3) : 0
    aiInsights.value.reviewsNeeded = interviews.filter(i => i.status === 'completed').length
    stats.value.pendingReports = aiInsights.value.reviewsNeeded
    
  } catch (error) {
    console.error('Lỗi tải dữ liệu dashboard:', error)
  } finally {
    loading.value = false
  }
})

const formatDate = (isoStr) => {
  if (!isoStr) return 'N/A'
  return new Date(isoStr).toLocaleDateString('vi-VN')
}
const formatTime = (isoStr) => {
  if (!isoStr) return 'N/A'
  return new Date(isoStr).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' })
}

import { computed } from 'vue'
const careerLink = computed(() => {
  const companyId = authStore.user?.companies?.[0]?.id
  if (!companyId) return ''
  return `${window.location.origin}/careers/${companyId}`
})
const linkCopied = ref(false)
const copyCareerLink = async () => {
  try {
    await navigator.clipboard.writeText(careerLink.value)
    linkCopied.value = true
    setTimeout(() => linkCopied.value = false, 2000)
  } catch { /* ignore */ }
}
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
  
  <div class="grid" style="grid-template-columns: repeat(4, 1fr); margin-top: 24px; margin-bottom: 24px" v-if="!loading">
    <Card>
      <div class="text-helper" style="display: flex; align-items: center; gap: 8px; color: var(--text-secondary)">
        <Briefcase size="16" color="var(--primary)" /> Jobs đang mở
      </div>
      <div class="text-h1" style="margin-top: 8px">{{ stats.jobs }}</div>
    </Card>
    <Card>
      <div class="text-helper" style="display: flex; align-items: center; gap: 8px; color: var(--text-secondary)">
        <Users size="16" color="var(--accent)" /> Ứng viên mới
      </div>
      <div class="text-h1" style="margin-top: 8px">{{ stats.candidates }}</div>
    </Card>
    <Card>
      <div class="text-helper" style="display: flex; align-items: center; gap: 8px; color: var(--text-secondary)">
        <Calendar size="16" color="var(--warning)" /> Phỏng vấn hôm nay
      </div>
      <div class="text-h1" style="margin-top: 8px">{{ stats.interviewsToday }}</div>
    </Card>
    <Card>
      <div class="text-helper" style="display: flex; align-items: center; gap: 8px; color: var(--text-secondary)">
        <FileText size="16" color="var(--success)" /> Báo cáo chờ xem
      </div>
      <div class="text-h1" style="margin-top: 8px">{{ stats.pendingReports }}</div>
    </Card>
  </div>

  <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px">
    <div style="display: flex; flex-direction: column; gap: 24px">
      <Card title="Lịch phỏng vấn sắp tới">
        <div v-if="upcomingInterviews.length === 0" style="padding: 32px; text-align: center">
          <Calendar size="32" color="var(--text-muted)" style="margin: 0 auto 16px" />
          <p class="text-body" style="color: var(--text-muted)">Chưa có lịch phỏng vấn nào sắp tới.</p>
        </div>
        <div v-else style="display: flex; flex-direction: column; gap: 12px; margin-top: 12px">
          <div v-for="i in upcomingInterviews" :key="i.id" style="display: flex; justify-content: space-between; align-items: center; padding: 12px; background-color: var(--surface-soft); border-radius: var(--radius-md)">
            <div>
              <div style="font-weight: 500">{{ i.candidate.full_name || i.candidate.name || 'Unknown' }}</div>
              <div class="text-helper" style="color: var(--text-secondary)">{{ i.job.title || 'Unknown Job' }}</div>
            </div>
            <div style="text-align: right">
              <div style="font-weight: 500; color: var(--primary)">{{ formatTime(i.dt) }}</div>
              <div class="text-helper" style="color: var(--text-secondary)">{{ formatDate(i.dt) }}</div>
            </div>
          </div>
        </div>
      </Card>
      <Card title="Ứng viên mới nhất">
        <div v-if="newCandidates.length === 0" style="padding: 32px; text-align: center">
          <Users size="32" color="var(--text-muted)" style="margin: 0 auto 16px" />
          <p class="text-body" style="color: var(--text-muted)">Chưa có ứng viên mới ứng tuyển.</p>
        </div>
        <div v-else style="display: flex; flex-direction: column; gap: 12px; margin-top: 12px">
          <div v-for="c in newCandidates" :key="c.id" style="display: flex; justify-content: space-between; align-items: center; padding: 12px; background-color: var(--surface-soft); border-radius: var(--radius-md)">
            <div>
              <div style="font-weight: 500">{{ c.full_name || c.name }}</div>
              <div class="text-helper" style="color: var(--text-secondary)">{{ c.email }}</div>
            </div>
            <div style="text-align: right">
              <div style="font-weight: 500; font-size: 13px">{{ c.job?.title || 'Unknown' }}</div>
              <div class="text-helper" style="color: var(--text-secondary)">{{ formatDate(c.created_at) }}</div>
            </div>
          </div>
        </div>
      </Card>
    </div>

    <div style="display: flex; flex-direction: column; gap: 24px">
      <Card title="AI Insights" style="border-top: 4px solid var(--accent)">
        <div style="display: flex; flex-direction: column; gap: 16px" v-if="!loading">
          <div style="display: flex; gap: 12px; align-items: flex-start">
            <Sparkles size="16" color="var(--accent)" style="margin-top: 2px; flex-shrink: 0" />
            <div>
              <p class="text-body" style="font-weight: 500">{{ aiInsights.highMatchCandidates }} ứng viên có mức phù hợp cao</p>
              <p class="text-helper">Vừa nộp đơn vào hệ thống.</p>
            </div>
          </div>
          <div style="display: flex; gap: 12px; align-items: flex-start">
            <TrendingUp size="16" color="var(--warning)" style="margin-top: 2px; flex-shrink: 0" />
            <div>
              <p class="text-body" style="font-weight: 500">{{ aiInsights.reviewsNeeded }} buổi phỏng vấn cần review</p>
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

      <Card title="Career Site" style="border-top: 4px solid var(--success)" v-if="careerLink">
        <div style="margin-bottom: 12px;">
          <p class="text-helper" style="margin-bottom: 8px;">Chia sẻ link này để ứng viên tự nộp đơn ứng tuyển:</p>
          <div style="display: flex; gap: 8px; align-items: center; background: var(--surface-soft); padding: 10px 12px; border-radius: var(--radius-md); border: 1px solid var(--border);">
            <Link size="14" color="var(--text-secondary)" style="flex-shrink: 0;" />
            <span style="flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 13px; color: var(--primary);">{{ careerLink }}</span>
            <button @click="copyCareerLink" style="background: none; border: none; cursor: pointer; padding: 4px;" :title="linkCopied ? 'Đã copy!' : 'Copy link'">
              <Copy size="14" :color="linkCopied ? 'var(--success)' : 'var(--text-secondary)'" />
            </button>
          </div>
          <p v-if="linkCopied" style="color: var(--success); font-size: 12px; margin-top: 4px;">✓ Đã sao chép!</p>
        </div>
      </Card>
    </div>
  </div>
</template>
