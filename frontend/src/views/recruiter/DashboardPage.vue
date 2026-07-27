<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import WelcomeAlert from '../../components/common/WelcomeAlert.vue'
import { Users, Briefcase, Calendar, Plus, FileText, Sparkles, TrendingUp, Link, Copy, CheckCircle } from 'lucide-vue-next'

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
    const newState = { ...history.state }
    delete newState.message
    window.history.replaceState(newState, '')
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
    
    // Count candidates whose AI-computed fit_score clears the high-match threshold.
    // fit_score is only present once AI matching has run; absent/0 counts as not high-match.
    aiInsights.value.highMatchCandidates = candidates.filter(c => {
      const fs = typeof c.fit_score === 'object' && c.fit_score !== null ? c.fit_score.Float64 : c.fit_score
      return Number(fs) >= 75
    }).length
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
  <div class="space-y-6 animate-fade-in">
    <WelcomeAlert 
      v-if="entryToast" 
      role="recruiter"
      title="Đăng nhập thành công!"
      :message="entryToast.message" 
      @close="entryToast = null" 
    />
    
    <!-- Header section -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6">
      <div>
        <h1 class="text-2xl font-bold" style="color: var(--text-main)">Tổng quan</h1>
        <p class="mt-1 text-sm" style="color: var(--text-secondary)">
          Theo dõi lịch phỏng vấn, ứng viên và báo cáo AI hôm nay.
        </p>
      </div>
      <div class="flex gap-3">
        <Button @click="router.push('/interviews/new')" class="bg-blue-600 hover:bg-blue-700 text-white border-none shadow-md shadow-blue-500/20">
          <Plus size="16" /> Tạo lịch phỏng vấn
        </Button>
      </div>
    </Card>
    
    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6" v-if="!loading">
      <Card class="p-6 relative group">
        <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <Briefcase size="64" class="text-blue-500" />
        </div>
        <div class="flex items-center gap-3 font-medium text-sm mb-3" style="color: var(--text-secondary)">
          <div class="p-2 rounded-lg" style="background-color: var(--primary-light); color: var(--primary)">
            <Briefcase size="18" />
          </div>
          Jobs đang mở
        </div>
        <div class="text-[28px] font-bold" style="color: var(--text-main)">{{ stats.jobs }}</div>
      </Card>

      <Card class="p-6 relative group">
        <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <Users size="64" class="text-blue-500" />
        </div>
        <div class="flex items-center gap-3 font-medium text-sm mb-3" style="color: var(--text-secondary)">
          <div class="p-2 rounded-lg" style="background-color: var(--primary-light); color: var(--primary)">
            <Users size="18" />
          </div>
          Ứng viên mới
        </div>
        <div class="text-[28px] font-bold" style="color: var(--text-main)">{{ stats.candidates }}</div>
      </Card>

      <Card class="p-6 relative group">
        <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <Calendar size="64" class="text-amber-500" />
        </div>
        <div class="flex items-center gap-3 font-medium text-sm mb-3" style="color: var(--text-secondary)">
          <div class="p-2 rounded-lg" style="background-color: rgba(217, 119, 6, 0.12); color: var(--warning)">
            <Calendar size="18" />
          </div>
          Phỏng vấn hôm nay
        </div>
        <div class="text-[28px] font-bold" style="color: var(--text-main)">{{ stats.interviewsToday }}</div>
      </Card>

      <Card class="p-6 relative group">
        <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <FileText size="64" class="text-emerald-500" />
        </div>
        <div class="flex items-center gap-3 font-medium text-sm mb-3" style="color: var(--text-secondary)">
          <div class="p-2 rounded-lg" style="background-color: rgba(16, 185, 129, 0.12); color: var(--success)">
            <FileText size="18" />
          </div>
          Báo cáo chờ xem
        </div>
        <div class="text-[28px] font-bold" style="color: var(--text-main)">{{ stats.pendingReports }}</div>
      </Card>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 space-y-6">
        <!-- Upcoming Interviews -->
        <Card class="p-6">
          <h2 class="text-lg font-bold mb-4" style="color: var(--text-main)">Lịch phỏng vấn sắp tới</h2>
          <div v-if="upcomingInterviews.length === 0" class="py-12 text-center">
            <Calendar size="48" style="color: var(--border); margin: 0 auto 16px auto;" />
            <p class="text-sm font-medium" style="color: var(--text-muted)">Chưa có lịch phỏng vấn nào sắp tới.</p>
          </div>
          <div v-else class="space-y-3">
            <div v-for="i in upcomingInterviews" :key="i.id" class="flex justify-between items-center p-4 rounded-xl transition-colors" style="background-color: var(--surface-soft); border: 1px solid transparent;" onmouseover="this.style.borderColor='var(--border)'" onmouseout="this.style.borderColor='transparent'">
              <div>
                <div class="font-semibold" style="color: var(--text-main)">{{ i.candidate.full_name || i.candidate.name || 'Unknown' }}</div>
                <div class="text-sm mt-1" style="color: var(--text-secondary)">{{ i.job.title || 'Unknown Job' }}</div>
              </div>
              <div class="text-right">
                <div class="font-bold" style="color: var(--primary)">{{ formatTime(i.dt) }}</div>
                <div class="text-xs font-medium mt-1 px-2 py-1 rounded-md inline-block" style="background-color: var(--surface); color: var(--text-secondary); border: 1px solid var(--border)">{{ formatDate(i.dt) }}</div>
              </div>
            </div>
          </div>
        </Card>
        
        <!-- New Candidates -->
        <Card class="p-6">
          <h2 class="text-lg font-bold mb-4" style="color: var(--text-main)">Ứng viên mới nhất</h2>
          <div v-if="newCandidates.length === 0" class="py-12 text-center">
            <Users size="48" style="color: var(--border); margin: 0 auto 16px auto;" />
            <p class="text-sm font-medium" style="color: var(--text-muted)">Chưa có ứng viên mới ứng tuyển.</p>
          </div>
          <div v-else class="space-y-3">
            <div v-for="c in newCandidates" :key="c.id" class="flex justify-between items-center p-4 rounded-xl transition-colors" style="background-color: var(--surface-soft); border: 1px solid transparent;" onmouseover="this.style.borderColor='var(--border)'" onmouseout="this.style.borderColor='transparent'">
              <div>
                <div class="font-semibold" style="color: var(--text-main)">{{ c.full_name || c.name }}</div>
                <div class="text-sm mt-1" style="color: var(--text-secondary)">{{ c.email }}</div>
              </div>
              <div class="text-right">
                <div class="font-medium text-sm" style="color: var(--text-main)">{{ c.job?.title || 'Unknown' }}</div>
                <div class="text-xs mt-1" style="color: var(--text-secondary)">{{ formatDate(c.created_at) }}</div>
              </div>
            </div>
          </div>
        </Card>
      </div>

      <div class="space-y-6">
        <!-- AI Insights -->
        <Card class="p-6 relative overflow-hidden" style="border-top: 4px solid var(--primary)">
          <div class="absolute -right-6 -top-6 pointer-events-none" style="color: var(--primary-light); opacity: 0.5;">
            <Sparkles size="100" />
          </div>
          <h2 class="text-lg font-bold mb-5 relative z-10" style="color: var(--text-main)">AI Insights</h2>
          <div class="space-y-5 relative z-10" v-if="!loading">
            <div class="flex gap-3">
              <div class="p-2 rounded-lg shrink-0 h-min" style="background-color: var(--primary-light); color: var(--primary)">
                <Sparkles size="18" />
              </div>
              <div>
                <p class="font-semibold text-sm" style="color: var(--text-main)">{{ aiInsights.highMatchCandidates }} ứng viên có mức phù hợp cao</p>
                <p class="text-xs mt-1" style="color: var(--text-secondary)">Vừa nộp đơn vào hệ thống.</p>
              </div>
            </div>
            <div class="flex gap-3">
              <div class="p-2 rounded-lg shrink-0 h-min" style="background-color: rgba(217, 119, 6, 0.12); color: var(--warning)">
                <TrendingUp size="18" />
              </div>
              <div>
                <p class="font-semibold text-sm" style="color: var(--text-main)">{{ aiInsights.reviewsNeeded }} buổi phỏng vấn cần review</p>
                <p class="text-xs mt-1" style="color: var(--text-secondary)">Báo cáo AI đã sẵn sàng để bạn đưa ra quyết định.</p>
              </div>
            </div>
          </div>
        </Card>
        
        <!-- Quick Actions -->
        <Card class="p-6">
          <h2 class="text-lg font-bold mb-4" style="color: var(--text-main)">Thao tác nhanh</h2>
          <div class="flex flex-col gap-3">
            <button @click="router.push('/jobs/new')" class="flex items-center gap-2 w-full p-3 font-medium rounded-xl transition-colors text-sm" style="background-color: var(--surface-soft); color: var(--text-main); border: 1px solid transparent;" onmouseover="this.style.borderColor='var(--border)'" onmouseout="this.style.borderColor='transparent'">
              <div class="p-1.5 rounded-md shadow-sm" style="background-color: var(--surface); border: 1px solid var(--border)"><Plus size="16" /></div>
              Tạo job mới
            </button>
            <button @click="router.push('/candidates/new')" class="flex items-center gap-2 w-full p-3 font-medium rounded-xl transition-colors text-sm" style="background-color: var(--surface-soft); color: var(--text-main); border: 1px solid transparent;" onmouseover="this.style.borderColor='var(--border)'" onmouseout="this.style.borderColor='transparent'">
              <div class="p-1.5 rounded-md shadow-sm" style="background-color: var(--surface); border: 1px solid var(--border)"><Plus size="16" /></div>
              Thêm ứng viên
            </button>
          </div>
        </Card>

        <!-- Career Site -->
        <Card v-if="careerLink" class="p-6" style="border-top: 4px solid var(--success)">
          <h2 class="text-lg font-bold mb-3" style="color: var(--text-main)">Career Site</h2>
          <p class="text-xs mb-3" style="color: var(--text-secondary)">Chia sẻ link này để ứng viên tự nộp đơn ứng tuyển:</p>
          <div class="flex gap-2 items-center p-3 rounded-xl" style="background-color: var(--surface-soft); border: 1px solid var(--border)">
            <Link size="16" style="color: var(--text-muted)" class="shrink-0" />
            <span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap text-xs font-medium" style="color: var(--success)">{{ careerLink }}</span>
            <button @click="copyCareerLink" class="p-1.5 rounded-md transition-colors" style="background-color: transparent;" onmouseover="this.style.backgroundColor='var(--border)'" onmouseout="this.style.backgroundColor='transparent'" :title="linkCopied ? 'Đã copy!' : 'Copy link'">
              <Copy size="14" :style="{ color: linkCopied ? 'var(--success)' : 'var(--text-muted)' }" />
            </button>
          </div>
          <p v-if="linkCopied" class="text-xs font-medium mt-2 flex items-center gap-1" style="color: var(--success)">
            <CheckCircle size="12" /> Đã sao chép!
          </p>
        </Card>
      </div>
    </div>
  </div>
</template>
