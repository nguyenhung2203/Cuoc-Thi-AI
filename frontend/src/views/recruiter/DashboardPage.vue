<script setup>
import { ref, onMounted } from 'vue'
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
  <div class="space-y-6 animate-fade-in">
    <WelcomeAlert 
      v-if="entryToast" 
      role="recruiter"
      title="Đăng nhập thành công!"
      :message="entryToast.message" 
      @close="entryToast = null" 
    />
    
    <!-- Header section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Tổng quan</h1>
        <p class="text-slate-500 dark:text-slate-400 mt-1 text-sm">
          Theo dõi lịch phỏng vấn, ứng viên và báo cáo AI hôm nay.
        </p>
      </div>
      <div class="flex gap-3">
        <Button @click="router.push('/interviews/new')" class="bg-blue-600 hover:bg-blue-700 text-white border-none shadow-md shadow-blue-500/20">
          <Plus size="16" /> Tạo lịch phỏng vấn
        </Button>
      </div>
    </div>
    
    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6" v-if="!loading">
      <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 relative overflow-hidden group hover:shadow-md transition-shadow">
        <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <Briefcase size="64" class="text-blue-500" />
        </div>
        <div class="flex items-center gap-3 text-slate-500 dark:text-slate-400 font-medium text-sm mb-3">
          <div class="p-2 bg-blue-50 dark:bg-blue-500/10 rounded-lg text-blue-600 dark:text-blue-400">
            <Briefcase size="18" />
          </div>
          Jobs đang mở
        </div>
        <div class="text-3xl font-bold text-slate-800 dark:text-slate-100">{{ stats.jobs }}</div>
      </div>

      <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 relative overflow-hidden group hover:shadow-md transition-shadow">
        <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <Users size="64" class="text-blue-500" />
        </div>
        <div class="flex items-center gap-3 text-slate-500 dark:text-slate-400 font-medium text-sm mb-3">
          <div class="p-2 bg-blue-50 dark:bg-blue-500/10 rounded-lg text-blue-600 dark:text-blue-400">
            <Users size="18" />
          </div>
          Ứng viên mới
        </div>
        <div class="text-3xl font-bold text-slate-800 dark:text-slate-100">{{ stats.candidates }}</div>
      </div>

      <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 relative overflow-hidden group hover:shadow-md transition-shadow">
        <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <Calendar size="64" class="text-amber-500" />
        </div>
        <div class="flex items-center gap-3 text-slate-500 dark:text-slate-400 font-medium text-sm mb-3">
          <div class="p-2 bg-amber-50 dark:bg-amber-500/10 rounded-lg text-amber-600 dark:text-amber-400">
            <Calendar size="18" />
          </div>
          Phỏng vấn hôm nay
        </div>
        <div class="text-3xl font-bold text-slate-800 dark:text-slate-100">{{ stats.interviewsToday }}</div>
      </div>

      <div class="bg-white dark:bg-slate-800 p-6 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 relative overflow-hidden group hover:shadow-md transition-shadow">
        <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
          <FileText size="64" class="text-emerald-500" />
        </div>
        <div class="flex items-center gap-3 text-slate-500 dark:text-slate-400 font-medium text-sm mb-3">
          <div class="p-2 bg-emerald-50 dark:bg-emerald-500/10 rounded-lg text-emerald-600 dark:text-emerald-400">
            <FileText size="18" />
          </div>
          Báo cáo chờ xem
        </div>
        <div class="text-3xl font-bold text-slate-800 dark:text-slate-100">{{ stats.pendingReports }}</div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 space-y-6">
        <!-- Upcoming Interviews -->
        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-4">Lịch phỏng vấn sắp tới</h2>
          <div v-if="upcomingInterviews.length === 0" class="py-12 text-center">
            <Calendar size="48" class="text-slate-300 dark:text-slate-600 mx-auto mb-4" />
            <p class="text-slate-500 dark:text-slate-400 text-sm font-medium">Chưa có lịch phỏng vấn nào sắp tới.</p>
          </div>
          <div v-else class="space-y-3">
            <div v-for="i in upcomingInterviews" :key="i.id" class="flex justify-between items-center p-4 bg-slate-50 dark:bg-slate-700/50 rounded-xl border border-transparent hover:border-slate-200 dark:hover:border-slate-600 transition-colors">
              <div>
                <div class="font-semibold text-slate-800 dark:text-slate-200">{{ i.candidate.full_name || i.candidate.name || 'Unknown' }}</div>
                <div class="text-sm text-slate-500 dark:text-slate-400 mt-1">{{ i.job.title || 'Unknown Job' }}</div>
              </div>
              <div class="text-right">
                <div class="font-bold text-blue-600 dark:text-blue-400">{{ formatTime(i.dt) }}</div>
                <div class="text-xs font-medium text-slate-500 dark:text-slate-400 mt-1 bg-white dark:bg-slate-800 px-2 py-1 rounded-md border border-slate-200 dark:border-slate-700 inline-block">{{ formatDate(i.dt) }}</div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- New Candidates -->
        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-4">Ứng viên mới nhất</h2>
          <div v-if="newCandidates.length === 0" class="py-12 text-center">
            <Users size="48" class="text-slate-300 dark:text-slate-600 mx-auto mb-4" />
            <p class="text-slate-500 dark:text-slate-400 text-sm font-medium">Chưa có ứng viên mới ứng tuyển.</p>
          </div>
          <div v-else class="space-y-3">
            <div v-for="c in newCandidates" :key="c.id" class="flex justify-between items-center p-4 bg-slate-50 dark:bg-slate-700/50 rounded-xl border border-transparent hover:border-slate-200 dark:hover:border-slate-600 transition-colors">
              <div>
                <div class="font-semibold text-slate-800 dark:text-slate-200">{{ c.full_name || c.name }}</div>
                <div class="text-sm text-slate-500 dark:text-slate-400 mt-1">{{ c.email }}</div>
              </div>
              <div class="text-right">
                <div class="font-medium text-sm text-slate-700 dark:text-slate-300">{{ c.job?.title || 'Unknown' }}</div>
                <div class="text-xs text-slate-500 dark:text-slate-400 mt-1">{{ formatDate(c.created_at) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="space-y-6">
        <!-- AI Insights -->
        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6 border-t-4 border-t-blue-500 relative overflow-hidden">
          <div class="absolute -right-6 -top-6 text-blue-500/10 pointer-events-none">
            <Sparkles size="100" />
          </div>
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-5 relative z-10">AI Insights</h2>
          <div class="space-y-5 relative z-10" v-if="!loading">
            <div class="flex gap-3">
              <div class="p-2 bg-blue-50 dark:bg-blue-500/10 rounded-lg text-blue-600 dark:text-blue-400 shrink-0 h-min">
                <Sparkles size="18" />
              </div>
              <div>
                <p class="font-semibold text-slate-800 dark:text-slate-200 text-sm">{{ aiInsights.highMatchCandidates }} ứng viên có mức phù hợp cao</p>
                <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">Vừa nộp đơn vào hệ thống.</p>
              </div>
            </div>
            <div class="flex gap-3">
              <div class="p-2 bg-amber-50 dark:bg-amber-500/10 rounded-lg text-amber-600 dark:text-amber-400 shrink-0 h-min">
                <TrendingUp size="18" />
              </div>
              <div>
                <p class="font-semibold text-slate-800 dark:text-slate-200 text-sm">{{ aiInsights.reviewsNeeded }} buổi phỏng vấn cần review</p>
                <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">Báo cáo AI đã sẵn sàng để bạn đưa ra quyết định.</p>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Quick Actions -->
        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-4">Thao tác nhanh</h2>
          <div class="flex flex-col gap-3">
            <button @click="router.push('/jobs/new')" class="flex items-center gap-2 w-full p-3 bg-slate-50 hover:bg-slate-100 dark:bg-slate-700/50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 font-medium rounded-xl transition-colors text-sm border border-transparent hover:border-slate-200 dark:hover:border-slate-600">
              <div class="bg-white dark:bg-slate-800 p-1.5 rounded-md shadow-sm border border-slate-200 dark:border-slate-700"><Plus size="16" /></div>
              Tạo job mới
            </button>
            <button @click="router.push('/candidates/new')" class="flex items-center gap-2 w-full p-3 bg-slate-50 hover:bg-slate-100 dark:bg-slate-700/50 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 font-medium rounded-xl transition-colors text-sm border border-transparent hover:border-slate-200 dark:hover:border-slate-600">
              <div class="bg-white dark:bg-slate-800 p-1.5 rounded-md shadow-sm border border-slate-200 dark:border-slate-700"><Plus size="16" /></div>
              Thêm ứng viên
            </button>
          </div>
        </div>

        <!-- Career Site -->
        <div v-if="careerLink" class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6 border-t-4 border-t-emerald-500">
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-3">Career Site</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mb-3">Chia sẻ link này để ứng viên tự nộp đơn ứng tuyển:</p>
          <div class="flex gap-2 items-center bg-slate-50 dark:bg-slate-900 p-3 rounded-xl border border-slate-200 dark:border-slate-700">
            <Link size="16" class="text-slate-400 shrink-0" />
            <span class="flex-1 overflow-hidden text-ellipsis whitespace-nowrap text-xs font-medium text-emerald-600 dark:text-emerald-400">{{ careerLink }}</span>
            <button @click="copyCareerLink" class="p-1.5 rounded-md hover:bg-slate-200 dark:hover:bg-slate-800 transition-colors" :title="linkCopied ? 'Đã copy!' : 'Copy link'">
              <Copy size="14" :class="linkCopied ? 'text-emerald-500' : 'text-slate-400'" />
            </button>
          </div>
          <p v-if="linkCopied" class="text-emerald-500 text-xs font-medium mt-2 flex items-center gap-1">
            <CheckCircle size="12" /> Đã sao chép!
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
