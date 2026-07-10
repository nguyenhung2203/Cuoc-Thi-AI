<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import WelcomeAlert from '../../components/common/WelcomeAlert.vue'
import { candidatePortalService } from '../../services/candidate-portal.service'
import { authStore } from '../../stores/auth.store'
import { Calendar, CalendarPlus, Bot, Star, UserCheck, ArrowRight, Building2, Play, ShieldCheck, FileText, Award, Sparkles, CheckCircle2, AlertCircle } from 'lucide-vue-next'

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

  if (!authStore.isAuthenticated && !localStorage.getItem('access_token')) {
    stats.value = {
      upcoming_interviews: 0,
      completed_mock_tests: 0,
      average_mock_score: 0,
      profile_completeness: 0
    }
    upcomingInterviews.value = []
    loading.value = false
    return
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
  <div class="space-y-6 pb-10">
    <WelcomeAlert
      v-if="entryToast"
      role="candidate"
      title="Thành công!"
      :message="entryToast.message"
      @close="entryToast = null"
    />

    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8 animate-rise">
      <div>
        <h1 class="page-title">
          Chào mừng trở lại, {{ authStore.user?.full_name || 'Ứng viên' }}
        </h1>
        <p class="page-subtitle">
          Theo dõi lịch phỏng vấn, luyện tập với AI và cải thiện kỹ năng trả lời của bạn.
        </p>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="dash-spinner mb-4"></div>
      <p class="text-helper">Đang tải dữ liệu dashboard...</p>
    </div>

    <template v-else>
      <!-- Stats Overview Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 stagger">
        <!-- Card 1 -->
        <Card class="card-elevate tilt-3d kpi-card">
          <div class="kpi-icon">
            <Calendar :size="22" />
          </div>
          <h3 class="kpi-label">Lịch sắp tới</h3>
          <p class="kpi-value">{{ stats.upcoming_interviews }}</p>
        </Card>

        <!-- Card 2 -->
        <Card class="card-elevate tilt-3d kpi-card">
          <div class="kpi-icon is-accent">
            <Bot :size="22" />
          </div>
          <h3 class="kpi-label">Luyện tập AI đã xong</h3>
          <p class="kpi-value">{{ stats.completed_mock_tests }}</p>
        </Card>

        <!-- Card 3 -->
        <Card class="card-elevate tilt-3d kpi-card">
          <div class="kpi-icon is-warning">
            <Star :size="22" />
          </div>
          <h3 class="kpi-label">Điểm AI trung bình</h3>
          <p class="kpi-value">{{ stats.average_mock_score.toFixed(1) }}</p>
        </Card>

        <!-- Card 4 -->
        <Card class="card-elevate tilt-3d kpi-card">
          <div class="kpi-icon is-success">
            <UserCheck :size="22" />
          </div>
          <h3 class="kpi-label">Mức độ hoàn thiện CV</h3>
          <p class="kpi-value">{{ stats.profile_completeness }}%</p>
          <div class="kpi-progress">
            <div class="kpi-progress-bar" :style="`width: ${stats.profile_completeness}%`"></div>
          </div>
        </Card>
      </div>
      
      <div class="grid grid-cols-1 lg:grid-cols-3 mt-8 gap-8">
        <!-- Lịch phỏng vấn sắp tới -->
        <div class="lg:col-span-2 space-y-8">
          <Card class="card-elevate section-card animate-rise">
            <div class="flex justify-between items-center mb-6">
              <h3 class="section-heading">
                <span class="kpi-icon"><Calendar :size="18" /></span>
                Lịch phỏng vấn sắp tới
              </h3>
              <button class="link-more" @click="router.push('/my-interviews')">
                Xem tất cả <ArrowRight :size="15" />
              </button>
            </div>

            <div v-if="upcomingInterviews.length === 0" class="empty-box">
              <div class="empty-icon"><CalendarPlus :size="30" /></div>
              <p class="text-secondary-strong">Bạn chưa có lịch phỏng vấn nào sắp tới.</p>
              <Button variant="outline" class="mt-4" @click="router.push('/job-board')">Tìm việc ngay</Button>
            </div>

            <div class="space-y-4">
              <div v-for="iv in upcomingInterviews" :key="iv.id"
                class="iv-item hover-rail">
                <div>
                  <h4 class="iv-title">{{ iv.job_title || iv.title }}</h4>
                  <p class="iv-company">
                    <Building2 :size="15" />
                    {{ iv.company_name || 'Công ty ẩn danh' }}
                  </p>
                  <div class="flex flex-wrap gap-2 mt-3">
                    <span class="badge badge-info">{{ formatDate(iv.scheduled_at) }}</span>
                    <span class="badge" :class="iv.mode === 'real' ? 'badge-danger' : 'badge-neutral'">
                      {{ iv.mode === 'real' ? 'Phỏng vấn thật' : 'Phỏng vấn thử' }}
                    </span>
                  </div>
                </div>
                <div>
                  <Button v-if="iv.mode === 'real'" variant="primary" class="sheen" @click="iv.join_link ? router.push(iv.join_link) : null">
                    Tham gia ngay
                  </Button>
                </div>
              </div>
            </div>
          </Card>

          <!-- Banner: AI Practice -->
          <div class="brand-banner sheen practice-banner animate-rise">
            <div class="banner-glyph"><Bot :size="180" /></div>
            <div class="banner-content">
              <h3 class="banner-title">Sẵn sàng vượt qua mọi câu hỏi phỏng vấn?</h3>
              <p class="banner-desc">Trải nghiệm phỏng vấn 1-kèm-1 với AI Interviewer. Luyện tập không giới hạn, nhận phản hồi ngay lập tức.</p>
              <button class="banner-cta" @click="router.push('/mock-setup')">
                <Play :size="18" /> Bắt đầu luyện tập
              </button>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-8">
          <Card class="card-elevate section-card animate-rise">
            <h3 class="section-heading mb-5">
              <span class="kpi-icon is-success"><ShieldCheck :size="18" /></span>
              Hành trang ứng viên
            </h3>
            <div class="space-y-3">
              <div class="prep-row" @click="router.push('/profile')">
                <div class="flex items-center gap-3">
                  <div class="kpi-icon"><FileText :size="18" /></div>
                  <div>
                    <span class="prep-title">Tải lên CV</span>
                    <span class="prep-sub">Bắt buộc để AI phân tích</span>
                  </div>
                </div>
                <CheckCircle2 v-if="stats.profile_completeness > 50" :size="22" class="text-success" />
                <AlertCircle v-else :size="22" class="text-warning" />
              </div>

              <div class="prep-row" @click="router.push('/profile')">
                <div class="flex items-center gap-3">
                  <div class="kpi-icon is-accent"><Award :size="18" /></div>
                  <div>
                    <span class="prep-title">Thêm Kỹ năng</span>
                    <span class="prep-sub">Giúp nhà tuyển dụng tìm thấy bạn</span>
                  </div>
                </div>
                <AlertCircle :size="22" class="text-warning" />
              </div>
            </div>
          </Card>

          <Card class="card-elevate coach-card animate-rise">
            <h3 class="section-heading mb-3">
              <span class="kpi-icon is-accent"><Sparkles :size="18" /></span>
              AI Career Coach
            </h3>
            <div class="ai-block coach-note">
              <p class="coach-text">
                Dựa trên kết quả phỏng vấn gần đây, tốc độ nói của bạn rất tốt, tuy nhiên bạn nên luyện tập thêm cách trả lời rành mạch các câu hỏi về <strong>Kỹ năng chuyên môn sâu</strong>.
              </p>
            </div>
            <Button variant="outline" class="w-full mt-4" @click="router.push('/mock-setup')">
              Luyện chủ đề này
            </Button>
          </Card>
        </div>
      </div>
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

/* Practice banner */
.practice-banner { padding: 30px; }
.banner-glyph { position: absolute; top: -30px; right: -20px; color: rgba(255,255,255,0.12); pointer-events: none; }
.banner-content { position: relative; z-index: 1; }
.banner-title { font-size: 22px; font-weight: 700; margin-bottom: 8px; }
.banner-desc { color: rgba(255,255,255,0.88); max-width: 34rem; margin-bottom: 22px; font-size: 15px; }
.banner-cta { display: inline-flex; align-items: center; gap: 8px; background: #fff; color: var(--primary); font-weight: 700; padding: 12px 22px; border: none; border-radius: var(--radius-full); cursor: pointer; box-shadow: var(--shadow-md); transition: transform 0.2s ease; }
.banner-cta:hover { transform: translateY(-2px); }

/* Prep rows */
.prep-row { display: flex; align-items: center; justify-content: space-between; padding: 12px; border-radius: var(--radius); cursor: pointer; transition: background 0.2s ease; }
.prep-row:hover { background: var(--surface-soft); }
.prep-title { display: block; font-size: 14px; font-weight: 600; color: var(--text-main); }
.prep-sub { display: block; font-size: 12px; color: var(--text-muted); }
.text-success { color: var(--success); }
.text-warning { color: var(--warning); }

.coach-card { padding: 24px; border-top: 3px solid var(--accent); }
.coach-note { padding: 16px; }
.coach-text { font-size: 14px; line-height: 1.6; color: var(--text-secondary); }
.coach-text strong { color: var(--accent); font-weight: 700; }
</style>
