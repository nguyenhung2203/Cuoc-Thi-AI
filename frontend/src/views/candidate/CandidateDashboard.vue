<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import WelcomeAlert from '../../components/common/WelcomeAlert.vue'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { apiService } from '../../services/api.service'
import { authStore } from '../../stores/auth.store'
import { langStore } from '../../stores/lang.store'
import { Calendar, CalendarPlus, Bot, Star, UserCheck, ArrowRight, Building2, Play, ShieldCheck, FileText, Award, Sparkles, CheckCircle2, AlertCircle, Briefcase, MapPin, Banknote, Clock } from 'lucide-vue-next'
import { formatSalaryTrieu } from '../../utils/formatters'

const router = useRouter()
const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)

const stats = ref({
  upcoming_interviews: 0,
  completed_mock_tests: 0,
  average_mock_score: 0,
  profile_completeness: 0
})
const upcomingInterviews = ref([])
const recommendedJobs = ref([])
const loading = ref(true)

const unwrap = (val) => {
  if (!val) return ''
  if (typeof val === 'object') {
    if ('String' in val) return val.Valid ? val.String : ''
    if ('Int64' in val) return val.Valid ? val.Int64 : ''
    if ('Float64' in val) return val.Valid ? val.Float64 : ''
  }
  return val
}

const getSalaryDisplay = (job) => formatSalaryTrieu(unwrap(job.salary_min), unwrap(job.salary_max))

const loadError = ref('')

const loadDashboard = async () => {
  loading.value = true
  loadError.value = ''

  if (!authStore.isAuthenticated && !localStorage.getItem('access_token')) {
    stats.value = {
      upcoming_interviews: 0,
      completed_mock_tests: 0,
      average_mock_score: 0,
      profile_completeness: 0
    }
    upcomingInterviews.value = []
    recommendedJobs.value = []
    loading.value = false
    return
  }

  try {
    const [statsData, interviewsData, jobsRes] = await Promise.all([
      candidatePortalService.getDashboardStats(),
      candidatePortalService.getInterviews(),
      apiService.getWithMeta('/public/all-jobs?page=1&page_size=3')
    ])

    stats.value = {
      upcoming_interviews: Number(statsData?.upcoming_interviews || 0),
      completed_mock_tests: Number(statsData?.completed_mock_tests || 0),
      average_mock_score: Number(statsData?.average_mock_score || 0),
      profile_completeness: Number(statsData?.profile_completeness || 0),
    }

    upcomingInterviews.value = (Array.isArray(interviewsData) ? interviewsData : [])
      .filter(i => i?.status === 'scheduled' || i?.status === 'active')
      .slice(0, 3)

    recommendedJobs.value = (Array.isArray(jobsRes?.data) ? jobsRes.data : []).map(j => ({
      ...j, company_name: unwrap(j.company_name), location: unwrap(j.location), employment_type: unwrap(j.employment_type), department: unwrap(j.department), level: unwrap(j.level)
    }))
  } catch (err) {
    console.error('Lỗi tải dữ liệu dashboard:', err)
    loadError.value = err?.message || 'Không tải được trang chủ ứng viên. Vui lòng thử lại.'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }
  await loadDashboard()
})

const formatDate = (dateString) => {
  if (!dateString) return ''
  const d = new Date(dateString)
  return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' }) + ', ' + d.toLocaleDateString('vi-VN')
}
</script>

<template>
  <div class="space-y-6 pb-10">
    <WelcomeAlert
      v-if="entryToast"
      role="candidate"
      title="Thành công!"
      :message="entryToast.message"
      @close="entryToast = null"
    />

    <!-- Framed Welcome Header (Exact style as MockSetup & PracticeHistory) -->
    <div class="header-box animate-rise mb-8 flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 rounded-2xl bg-[var(--surface)] border border-[var(--border)] shadow-sm gap-4">
      <div>
        <h1 class="text-h1 mb-1.5 text-[var(--text-main)]">
          {{ langStore.t('dashboard', 'welcome') }}, {{ authStore.user?.full_name || 'Ứng viên' }}!
        </h1>
        <p class="text-secondary text-sm">
          Hệ thống AI đã sẵn sàng hỗ trợ bạn đánh giá kỹ năng và tìm kiếm cơ hội phù hợp.
        </p>
      </div>
      <div class="shrink-0">
        <Button variant="primary" @click="router.push('/mock-setup')" class="sheen">
          <Play :size="16" class="mr-1.5 shrink-0" />
          <span>Luyện tập phỏng vấn ngay</span>
        </Button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="dash-spinner mb-4"></div>
      <p class="text-helper">Đang tải dữ liệu dashboard...</p>
    </div>

    <div v-else-if="loadError" class="empty-state bg-[var(--surface)] p-10 rounded-2xl border border-dashed border-[var(--border)] text-center">
      <AlertCircle :size="36" class="mx-auto mb-3 text-[var(--danger)]" />
      <h3 class="text-base font-bold text-[var(--text-main)] mb-1">Không tải được dữ liệu</h3>
      <p class="text-sm text-[var(--text-secondary)] mb-4">{{ loadError }}</p>
      <Button variant="primary" @click="loadDashboard">Thử lại</Button>
    </div>

    <template v-else>
      <!-- Stats Overview Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 stagger">
        <!-- Card 1 -->
        <Card class="card-elevate tilt-3d kpi-card">
          <div class="kpi-icon">
            <Calendar :size="22" />
          </div>
          <h3 class="kpi-label">{{ langStore.t('dashboard', 'upcoming') }}</h3>
          <p class="kpi-value">{{ stats.upcoming_interviews }}</p>
        </Card>

        <!-- Card 2 -->
        <Card class="card-elevate tilt-3d kpi-card">
          <div class="kpi-icon is-accent">
            <Bot :size="22" />
          </div>
          <h3 class="kpi-label">{{ langStore.t('dashboard', 'completedAI') }}</h3>
          <p class="kpi-value">{{ stats.completed_mock_tests }}</p>
        </Card>

        <!-- Card 3 -->
        <Card class="card-elevate tilt-3d kpi-card">
          <div class="kpi-icon is-warning">
            <Star :size="22" />
          </div>
          <h3 class="kpi-label">{{ langStore.t('dashboard', 'avgScore') }}</h3>
          <p class="kpi-value">{{ stats.average_mock_score ? stats.average_mock_score.toFixed(1) : '0' }}</p>
        </Card>

        <!-- Card 4 -->
        <Card class="card-elevate tilt-3d kpi-card">
          <div class="kpi-icon is-success">
            <UserCheck :size="22" />
          </div>
          <h3 class="kpi-label">{{ langStore.t('dashboard', 'profileComplete') }}</h3>
          <p class="kpi-value">{{ stats.profile_completeness }}%</p>
          <div class="kpi-progress">
            <div class="kpi-progress-bar" :style="`width: ${stats.profile_completeness}%`"></div>
          </div>
        </Card>
      </div>
      
      <!-- 2-Column Balanced Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 mt-8 gap-6 items-stretch">
        <!-- Main Column (2 cols) -->
        <div class="lg:col-span-2 flex flex-col gap-6">
          <!-- Lịch phỏng vấn sắp tới -->
          <Card class="card-elevate section-card animate-rise flex-1 flex flex-col justify-between">
            <div>
              <div class="flex justify-between items-center mb-5">
                <h3 class="section-heading">
                  <span class="kpi-icon"><Calendar :size="18" /></span>
                  {{ langStore.t('dashboard', 'upcomingSection') }}
                </h3>
                <button class="link-more" @click="router.push('/my-interviews')">
                  {{ langStore.t('dashboard', 'viewAll') }} <ArrowRight :size="15" />
                </button>
              </div>

              <div v-if="upcomingInterviews.length === 0" class="empty-box my-auto">
                <div class="empty-icon"><CalendarPlus :size="28" /></div>
                <p class="text-secondary-strong">{{ langStore.t('dashboard', 'noUpcoming') }}</p>
                <Button variant="outline" class="mt-3" @click="router.push('/job-board')">{{ langStore.t('dashboard', 'findJobsNow') }}</Button>
              </div>

              <div v-else class="space-y-3">
                <div v-for="iv in upcomingInterviews" :key="iv.id" class="iv-item hover-rail">
                  <div>
                    <h4 class="iv-title">{{ iv.job_title || iv.title }}</h4>
                    <p class="iv-company">
                      <Building2 :size="15" />
                      {{ iv.company_name || 'Công ty ẩn danh' }}
                    </p>
                    <div class="flex flex-wrap gap-2 mt-2">
                      <span class="badge badge-info">{{ formatDate(iv.scheduled_at) }}</span>
                      <span class="badge" :class="iv.mode === 'real' ? 'badge-danger' : 'badge-neutral'">
                        {{ iv.mode === 'real' ? langStore.t('dashboard', 'realMode') : langStore.t('dashboard', 'mockMode') }}
                      </span>
                    </div>
                  </div>
                  <div>
                    <Button v-if="iv.mode === 'real'" variant="primary" class="sheen" @click="iv.join_link ? router.push(iv.join_link) : null">
                      {{ langStore.t('dashboard', 'joinNow') }}
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </Card>

          <!-- Gợi ý việc làm -->
          <Card class="card-elevate section-card animate-rise flex-1 flex flex-col justify-between">
            <div>
              <div class="flex justify-between items-center mb-5">
                <h3 class="section-heading">
                  <span class="kpi-icon is-accent"><Briefcase :size="18" /></span>
                  Gợi ý việc làm phù hợp
                </h3>
                <button class="link-more" @click="router.push('/job-board')">
                  {{ langStore.t('dashboard', 'viewAll') }} <ArrowRight :size="15" />
                </button>
              </div>

              <div v-if="recommendedJobs.length === 0" class="empty-box my-auto" style="padding: 24px;">
                <div class="empty-icon"><Briefcase :size="28" /></div>
                <p class="text-secondary-strong">Hiện tại chưa có công việc gợi ý phù hợp.</p>
                <Button variant="outline" class="mt-3" @click="router.push('/job-board')">Khám phá tất cả việc làm</Button>
              </div>
              <div v-else class="space-y-3">
                <div v-for="job in recommendedJobs" :key="job.id" class="iv-item hover-rail !p-3.5" style="cursor: pointer;" @click="router.push(`/careers/${job.company_id}/jobs/${job.id}`)">
                  <div>
                    <h4 class="iv-title !text-[15px]">{{ job.title }}</h4>
                    <p class="iv-company !mt-1">
                      <Building2 :size="14" />
                      {{ job.company_name || 'Công ty ẩn danh' }}
                    </p>
                    <div class="flex flex-wrap gap-2 mt-2">
                      <span class="badge badge-info"><MapPin :size="13" class="mr-1"/> {{ job.location || 'Bất kỳ' }}</span>
                      <span class="badge badge-success"><Banknote :size="13" class="mr-1"/> {{ getSalaryDisplay(job) }}</span>
                    </div>
                  </div>
                  <div>
                    <Button variant="primary" class="sheen !px-3 !py-1.5 !text-xs" @click.stop="router.push(`/careers/${job.company_id}/jobs/${job.id}`)">
                      Ứng tuyển
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </Card>
        </div>

        <!-- Sidebar Column (1 col) -->
        <div class="flex flex-col gap-6">
          <!-- Hành trang ứng viên -->
          <Card class="card-elevate section-card animate-rise">
            <h3 class="section-heading mb-4">
              <span class="kpi-icon is-success"><ShieldCheck :size="18" /></span>
              {{ langStore.t('dashboard', 'prepTitle') }}
            </h3>
            <div class="space-y-2.5">
              <div class="prep-row" @click="router.push('/profile')">
                <div class="flex items-center gap-3">
                  <div class="kpi-icon"><FileText :size="18" /></div>
                  <div>
                    <span class="prep-title">Tải lên CV</span>
                    <span class="prep-sub">Bắt buộc để AI phân tích</span>
                  </div>
                </div>
                <CheckCircle2 v-if="stats.profile_completeness > 50" :size="20" class="text-success" />
                <AlertCircle v-else :size="20" class="text-warning" />
              </div>

              <div class="prep-row" @click="router.push('/profile')">
                <div class="flex items-center gap-3">
                  <div class="kpi-icon is-accent"><Award :size="18" /></div>
                  <div>
                    <span class="prep-title">Thêm Kỹ năng</span>
                    <span class="prep-sub">Giúp nhà tuyển dụng tìm thấy bạn</span>
                  </div>
                </div>
                <AlertCircle :size="20" class="text-warning" />
              </div>
            </div>
          </Card>

          <!-- Hoạt động gần đây -->
          <Card class="card-elevate section-card animate-rise flex-1 flex flex-col justify-between">
            <div>
              <h3 class="section-heading mb-4">
                <span class="kpi-icon"><Clock :size="18" /></span>
                Hoạt động gần đây
              </h3>
              <div class="empty-box" style="padding: 24px;">
                <div class="empty-icon"><Clock :size="26" /></div>
                <p class="text-secondary-strong">Chưa có hoạt động nào gần đây.</p>
              </div>
            </div>
          </Card>
        </div>
      </div>

      <!-- Full-Width AI Coach Banner at the Bottom -->
      <Card class="card-elevate animate-rise mt-6 flex flex-col sm:flex-row items-center justify-between p-6 rounded-2xl border border-[rgba(37,99,235,0.2)] bg-gradient-to-r from-[var(--primary-light)] to-[rgba(236,254,255,0.7)] gap-6">
        <div class="space-y-2 max-w-2xl">
          <div class="flex items-center gap-2 text-xs font-bold text-[var(--primary)] uppercase tracking-wider">
            <Sparkles :size="14" /> AI Career Coach
          </div>
          <h3 class="text-xl font-bold text-slate-900">{{ langStore.t('dashboard', 'bannerTitle') }}</h3>
          <p class="text-sm text-slate-600">Trải nghiệm phỏng vấn giả lập 1-kèm-1 với Trợ lý AI. Luyện tập không giới hạn và nhận phản hồi Rubric ngay lập tức.</p>
        </div>
        <div class="shrink-0 flex items-center gap-4">
          <Button variant="primary" class="sheen !px-5 !py-2.5 !text-sm font-semibold shadow-md" @click="router.push('/mock-setup')">
            <Play :size="16" class="mr-2" /> {{ langStore.t('dashboard', 'bannerCta') }}
          </Button>
        </div>
      </Card>
    </template>
  </div>
</template>

<style scoped>
.dash-spinner {
  width: 44px; height: 44px; border-radius: 50%;
  border: 3px solid var(--primary-light);
  border-top-color: var(--primary);
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* KPI cards */
.kpi-card { padding: 22px; }
.kpi-label { color: var(--text-secondary); font-size: 13px; font-weight: 500; margin-top: 14px; }
.kpi-value { font-size: 30px; font-weight: 700; color: var(--text-main); margin-top: 4px; letter-spacing: -0.02em; }
.kpi-progress { width: 100%; height: 6px; border-radius: 999px; background: var(--surface-soft); margin-top: 14px; overflow: hidden; }
.kpi-progress-bar { height: 100%; border-radius: 999px; background: var(--gradient-brand); transition: width 1s cubic-bezier(0.2,0.8,0.2,1); }

.section-card { padding: 24px; }
.section-heading { display: flex; align-items: center; gap: 10px; font-size: 17px; font-weight: 700; color: var(--text-main); }
.link-more { display: inline-flex; align-items: center; gap: 4px; color: var(--primary); font-weight: 600; font-size: 14px; }
.link-more:hover { color: var(--primary-hover); }

.empty-box { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 6px; padding: 40px 16px; border: 2px dashed var(--border); border-radius: var(--radius); background: var(--surface-soft); }
.empty-icon { display: flex; align-items: center; justify-content: center; width: 60px; height: 60px; border-radius: 50%; background: var(--primary-light); color: var(--primary); margin-bottom: 6px; }
.text-secondary-strong { color: var(--text-secondary); font-weight: 500; }

.iv-item { display: flex; flex-direction: column; gap: 16px; padding: 18px; padding-left: 22px; border: 1px solid var(--border); border-radius: var(--radius); background: var(--surface); transition: box-shadow 0.3s ease, border-color 0.3s ease; }
.iv-item:hover { box-shadow: var(--shadow-md); border-color: var(--primary-light); }
.iv-title { font-size: 16px; font-weight: 700; color: var(--text-main); }
.iv-company { display: flex; align-items: center; gap: 6px; color: var(--text-secondary); font-size: 14px; font-weight: 500; margin-top: 4px; }
@media (min-width: 640px) { .iv-item { flex-direction: row; align-items: center; justify-content: space-between; } }


.prep-row { display: flex; align-items: center; justify-content: space-between; padding: 12px; border-radius: var(--radius); cursor: pointer; transition: background 0.2s ease; }
.prep-row:hover { background: var(--surface-soft); }
.prep-title { display: block; font-size: 14px; font-weight: 600; color: var(--text-main); }
.prep-sub { display: block; font-size: 12px; color: var(--text-muted); }
.text-success { color: var(--success); }
.text-warning { color: var(--warning); }

.coach-card { padding: 24px; border-top: 3px solid var(--accent); }

.coach-card { padding: 24px; border-top: 3px solid var(--accent); }
</style>
