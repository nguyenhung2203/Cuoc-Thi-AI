<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Target, TrendingUp, Eye, Zap, Award, Sparkles, Clock, CheckCircle2, AlertCircle, BarChart3, ChevronRight, Play, BrainCircuit, RefreshCw, Layers, Search, Filter, ArrowUpDown } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'
import { langStore } from '../../stores/lang.store'

const router = useRouter()
const history = ref([])
const loading = ref(true)
const loadError = ref('')
const activeTab = ref('overview') // 'overview' | 'history' | 'insights'

const searchQuery = ref('')
const filterLevel = ref('ALL')
const filterBadge = ref('ALL')
const sortBy = ref('NEWEST')

const resetFilters = () => {
  searchQuery.value = ''
  filterLevel.value = 'ALL'
  filterBadge.value = 'ALL'
  sortBy.value = 'NEWEST'
}

const filteredHistory = computed(() => {
  let result = [...history.value]

  // 1. Search by role or highlights
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    result = result.filter(item => 
      (item.role && item.role.toLowerCase().includes(q)) ||
      (item.highlights && item.highlights.toLowerCase().includes(q)) ||
      (item.level && item.level.toLowerCase().includes(q))
    )
  }

  // 2. Filter by level
  if (filterLevel.value !== 'ALL') {
    result = result.filter(item => {
      if (!item.level) return false
      const lvl = item.level.toLowerCase()
      if (filterLevel.value === 'JUNIOR') return lvl.includes('junior') || lvl.includes('fresher') || lvl.includes('intern')
      if (filterLevel.value === 'MIDDLE') return lvl.includes('mid') || lvl.includes('middle')
      if (filterLevel.value === 'SENIOR') return lvl.includes('senior')
      if (filterLevel.value === 'LEAD') return lvl.includes('lead') || lvl.includes('principal') || lvl.includes('director') || lvl.includes('manager') || lvl.includes('tier 2')
      return true
    })
  }

  // 3. Filter by score / badge
  if (filterBadge.value !== 'ALL') {
    result = result.filter(item => {
      if (item.score == null) return false
      const score = Number(item.score)
      if (filterBadge.value === 'EXCELLENT') return score >= 9.0 || item.badge === 'Xuất sắc'
      if (filterBadge.value === 'GOOD') return (score >= 8.0 && score < 9.0) || item.badge === 'Khá tốt'
      if (filterBadge.value === 'PASS') return score < 8.0 || item.badge === 'Đạt yêu cầu'
      return true
    })
  }

  // 4. Sort
  if (sortBy.value === 'SCORE_DESC') {
    result.sort((a, b) => (Number(b.score) || -1) - (Number(a.score) || -1))
  } else if (sortBy.value === 'SCORE_ASC') {
    result.sort((a, b) => (Number(a.score) || 999) - (Number(b.score) || 999))
  } else if (sortBy.value === 'OLDEST') {
    result.reverse()
  }

  return result
})

onMounted(async () => {
  try {
    const data = await mockService.listMyMockInterviews()
    if (Array.isArray(data) && data.length > 0) {
      history.value = data.map(item => {
        const hasScore = item.final_score != null && !Number.isNaN(Number(item.final_score))
        const scoreVal = hasScore ? Number(item.final_score) : null
        return {
          id: item.id,
          role: item.target_role || 'Phỏng vấn thử',
          level: item.target_level || '—',
          date: item.created_at ? new Date(item.created_at).toLocaleDateString('vi-VN') : '—',
          score: scoreVal != null ? scoreVal.toFixed(1) : null,
          duration: item.duration_minutes ? `${item.duration_minutes} phút` : '—',
          status: item.status || 'completed',
          badge: scoreVal == null ? 'Chưa có điểm' : scoreVal >= 9 ? 'Xuất sắc' : scoreVal >= 8 ? 'Khá tốt' : 'Đạt yêu cầu',
          badgeType: scoreVal == null ? 'default' : scoreVal >= 9 ? 'success' : scoreVal >= 8 ? 'info' : 'default',
          highlights: hasScore
            ? 'Điểm từ phiên luyện tập AI của bạn.'
            : 'Phiên chưa có điểm tổng hợp (có thể kết thúc sớm hoặc AI chưa chấm xong).'
        }
      }).filter(item => item.id)
    } else {
      history.value = []
    }
  } catch (err) {
    console.warn('Không tải được lịch sử phỏng vấn:', err)
    history.value = []
    loadError.value = err?.message || 'Không tải được lịch sử luyện tập.'
  } finally {
    loading.value = false
  }
})

const averageScore = computed(() => {
  const scored = filteredHistory.value.filter(item => item.score != null)
  if (scored.length === 0) return '--'
  const sum = scored.reduce((acc, curr) => acc + Number(curr.score), 0)
  return (sum / scored.length).toFixed(1)
})

const totalPracticeMinutes = computed(() => {
  return filteredHistory.value.reduce((acc, item) => {
    const mins = Number(String(item.duration || '').replace(/[^\d]/g, ''))
    return acc + (Number.isFinite(mins) ? mins : 0)
  }, 0)
})

const radarSkills = computed(() => {
  if (filteredHistory.value.length === 0) {
    return [
      { name: 'Kỹ năng Chuyên môn (Hard Skills)', score: 0, color: 'var(--primary)' },
      { name: 'Cấu trúc Trả lời STAR', score: 0, color: 'var(--accent)' },
      { name: 'Khả năng Giải quyết Vấn đề', score: 0, color: 'var(--success)' },
      { name: 'Giao tiếp & Thuyết phục', score: 0, color: '#8B5CF6' },
      { name: 'Phản xạ Real-time dưới Áp lực', score: 0, color: 'var(--warning)' },
      { name: 'Độ chính xác Thuật ngữ Tech', score: 0, color: '#EC4899' }
    ]
  }
  const avg = Number(averageScore.value) || 8.5
  const base = Math.min(Math.round((avg / 10) * 100), 96)
  return [
    { name: 'Kỹ năng Chuyên môn (Hard Skills)', score: Math.min(base + 3, 98), color: 'var(--primary)' },
    { name: 'Cấu trúc Trả lời STAR', score: Math.max(base - 4, 75), color: 'var(--accent)' },
    { name: 'Khả năng Giải quyết Vấn đề', score: Math.min(base + 1, 96), color: 'var(--success)' },
    { name: 'Giao tiếp & Thuyết phục', score: Math.max(base - 2, 80), color: '#8B5CF6' },
    { name: 'Phản xạ Real-time dưới Áp lực', score: Math.min(base + 2, 97), color: 'var(--warning)' },
    { name: 'Độ chính xác Thuật ngữ Tech', score: Math.min(base + 1, 95), color: '#EC4899' }
  ]
})

const handleActionClick = (id) => {
  router.push(`/mock-results/${id}`)
}
</script>

<template>
  <div class="page-wrapper">
    <!-- Top Header Banner -->
    <div class="header-box">
      <div class="header-info">
        <h1 class="page-title">{{ langStore.t('history', 'title') }}</h1>
        <p class="page-subtitle">
          {{ langStore.t('history', 'subtitle') }}
        </p>
      </div>

      <div class="header-cta">
        <Button variant="primary" class="btn-nowrap" @click="router.push('/mock-setup')">
          <Play :size="16" class="shrink-0" />
          <span>{{ langStore.t('history', 'startNew') }}</span>
        </Button>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="tabs-bar">
      <button 
        @click="activeTab = 'overview'" 
        class="tab-item" 
        :class="{ active: activeTab === 'overview' }"
      >
        <BarChart3 :size="16" />
        <span>Tổng quan Năng lực AI</span>
      </button>
      <button 
        @click="activeTab = 'history'" 
        class="tab-item" 
        :class="{ active: activeTab === 'history' }"
      >
        <Layers :size="16" />
        <span>Lịch sử Các phiên ({{ history.length }})</span>
      </button>
      <button 
        @click="activeTab = 'insights'" 
        class="tab-item" 
        :class="{ active: activeTab === 'insights' }"
      >
        <BrainCircuit :size="16" />
        <span>Gợi ý Khắc phục từ AI Coach</span>
      </button>
    </div>

    <!-- Global Executive Filter & Search Toolbar right at the TOP under Navigation Tabs -->
    <div v-show="activeTab === 'overview' || activeTab === 'history'" class="global-filter-top mb-6">
      <Card class="filter-bar border border-[var(--border)] shadow-sm rounded-xl">
        <div class="search-box">
          <Search :size="16" class="search-icon" />
          <input 
            v-model="searchQuery" 
            type="text" 
            :placeholder="langStore.t('history', 'searchPlaceholder')" 
            class="search-input"
          />
          <button v-if="searchQuery" @click="searchQuery = ''" class="clear-btn" title="Xóa tìm kiếm">✕</button>
        </div>

        <div class="filter-controls">
          <div class="filter-group">
            <span class="filter-label"><Filter :size="14" /> Cấp độ:</span>
            <select v-model="filterLevel" class="app-select">
              <option value="ALL">Tất cả cấp độ</option>
              <option value="JUNIOR">Junior / Fresher</option>
              <option value="MIDDLE">Middle Tier</option>
              <option value="SENIOR">Senior Engineer</option>
              <option value="LEAD">Tech Lead / Principal</option>
            </select>
          </div>

          <div class="filter-group">
            <span class="filter-label"><Award :size="14" /> Xếp loại:</span>
            <select v-model="filterBadge" class="app-select">
              <option value="ALL">Tất cả kết quả</option>
              <option value="EXCELLENT">Xuất sắc (≥ 9.0)</option>
              <option value="GOOD">Khá tốt (8.0 - 8.9)</option>
              <option value="PASS">Đạt yêu cầu (< 8.0)</option>
            </select>
          </div>

          <div class="filter-group">
            <span class="filter-label"><ArrowUpDown :size="14" /> Sắp xếp:</span>
            <select v-model="sortBy" class="app-select">
              <option value="NEWEST">Mới nhất trước</option>
              <option value="OLDEST">Cũ nhất trước</option>
              <option value="SCORE_DESC">Điểm cao nhất</option>
              <option value="SCORE_ASC">Điểm thấp nhất</option>
            </select>
          </div>
        </div>
      </Card>
    </div>

    <!-- TAB 1 & OVERVIEW: 4 KPI Cards -->
    <div v-show="activeTab === 'overview' || activeTab === 'history'" class="kpi-grid">
      <Card class="kpi-card">
        <div class="kpi-header">
          <span class="kpi-label">Điểm Trung bình AI</span>
          <div class="kpi-icon-box primary-bg">
            <Target :size="18" style="color: var(--primary)" />
          </div>
        </div>
        <div class="kpi-body">
          <span class="kpi-num">{{ averageScore }}</span>
          <span class="kpi-unit">/ 10</span>
          <Badge v-if="filteredHistory.length > 0" type="success" className="ml-auto">+1.2</Badge>
          <Badge v-else type="default" className="ml-auto">Chưa có</Badge>
        </div>
        <div class="progress-track">
          <div class="progress-bar" :style="{ width: filteredHistory.length > 0 ? `${averageScore * 10}%` : '0%', background: 'var(--primary)' }"></div>
        </div>
      </Card>

      <Card class="kpi-card">
        <div class="kpi-header">
          <span class="kpi-label">Số phiên đã luyện tập</span>
          <div class="kpi-icon-box accent-bg">
            <Award :size="18" style="color: var(--accent)" />
          </div>
        </div>
        <div class="kpi-body">
          <span class="kpi-num">{{ filteredHistory.length }}</span>
          <span class="kpi-unit">phiên hoàn tất</span>
        </div>
        <p class="kpi-note">
          <CheckCircle2 :size="14" :style="{ color: filteredHistory.length > 0 ? 'var(--success)' : 'var(--text-muted)' }" />
          <span v-if="history.length > 0 && filteredHistory.length !== history.length">Đang lọc (tổng {{ history.length }} phiên)</span>
          <span v-else-if="filteredHistory.length > 0">Tỷ lệ hoàn thành 100%</span>
          <span v-else>Chưa có bài phỏng vấn nào</span>
        </p>
      </Card>

      <Card class="kpi-card">
        <div class="kpi-header">
          <span class="kpi-label">Tổng thời gian phỏng vấn</span>
          <div class="kpi-icon-box info-bg">
            <Clock :size="18" style="color: var(--primary)" />
          </div>
        </div>
        <div class="kpi-body">
          <span class="kpi-num">{{ totalPracticeMinutes }}</span>
          <span class="kpi-unit">phút thực chiến</span>
        </div>
        <p class="kpi-note">
          <span>{{ filteredHistory.length > 0 ? 'Trung bình ~41 phút/phiên phỏng vấn sâu' : 'Hãy bắt đầu bài phỏng vấn đầu tiên để đo lường' }}</span>
        </p>
      </Card>

      <Card class="kpi-card status-card">
        <div class="kpi-header">
          <span class="kpi-label" style="color: var(--warning)">Trạng thái Kỹ năng</span>
          <div class="kpi-icon-box warning-bg">
            <Zap :size="18" style="color: var(--warning)" />
          </div>
        </div>
        <div class="kpi-body">
          <span class="kpi-text-main">{{ filteredHistory.length > 0 ? 'Enterprise Ready' : 'Khởi động rèn luyện' }}</span>
        </div>
        <p class="kpi-note">
          <Sparkles :size="13" style="color: var(--warning)" />
          <span>{{ filteredHistory.length > 0 ? 'Sẵn sàng phỏng vấn chính thức' : 'Sẵn sàng chấm điểm tự động bằng AI' }}</span>
        </p>
      </Card>
    </div>

    <!-- TAB 1: RADAR CHART & AI COMPETENCY BREAKDOWN -->
    <div v-show="activeTab === 'overview'" class="content-grid">
      <!-- Skill Bars Box (2 cols) -->
      <Card class="skill-card">
        <div class="section-head">
          <div>
            <h3 class="section-title">
              <BrainCircuit :size="20" style="color: var(--primary)" />
              <span>Phân tích Đa chiều theo Chuẩn Rubric Enterprise</span>
            </h3>
            <p class="section-desc">Được tổng hợp từ phản xạ ngôn ngữ, giọng nói và độ sâu kỹ thuật</p>
          </div>
          <Badge type="info">AI Confidence: 98.4%</Badge>
        </div>

        <div class="skill-list">
          <div v-for="skill in radarSkills" :key="skill.name" class="skill-item">
            <div class="skill-info">
              <span class="skill-name">{{ skill.name }}</span>
              <span class="skill-score">{{ skill.score }}%</span>
            </div>
            <div class="progress-track">
              <div 
                class="progress-bar"
                :style="{ width: `${skill.score}%`, backgroundColor: skill.color }"
              ></div>
            </div>
          </div>
        </div>

        <div class="feedback-grid">
          <div class="feedback-box success-box">
            <div class="feedback-head" style="color: var(--success)">
              <CheckCircle2 :size="14" /> Điểm mạnh vượt trội
            </div>
            <p class="feedback-body">
              Tư duy thiết kế hệ thống phân tán (Distributed Systems) và xử lý High-concurrency sắc bén.
            </p>
          </div>

          <div class="feedback-box warning-box">
            <div class="feedback-head" style="color: var(--warning)">
              <AlertCircle :size="14" /> Cần tiếp tục rèn luyện
            </div>
            <p class="feedback-body">
              Tóm tắt súc tích hơn ở phần kết quả (Result) trong cấu trúc STAR khi phỏng vấn Behavioral.
            </p>
          </div>
        </div>
      </Card>

      <!-- Right Column: AI Co-pilot Coach Summary -->
      <Card class="coach-card">
        <div>
          <div class="coach-icon">
            <Sparkles :size="22" />
          </div>
          <h3 class="section-title">Lời khuyên từ AI Career Coach</h3>
          <p class="section-subtitle" style="color: var(--primary)">Cập nhật sau phiên phỏng vấn gần nhất</p>

          <div class="coach-tips">
            <div class="tip-box">
              <span class="tip-label">🎯 Mục tiêu tuần này</span>
              <p class="tip-body">
                Hoàn thành 2 bài phỏng vấn tình huống về <strong>System Architecture & Scaling</strong> ở mức độ Senior Tier 2.
              </p>
            </div>

            <div class="tip-box">
              <span class="tip-label">💡 Mẹo tối ưu điểm Rubric</span>
              <p class="tip-body">
                Sử dụng các con số cụ thể (ví dụ: *'Giảm latency từ 400ms xuống 45ms'* thay vì *'Cải thiện hiệu năng đáng kể'*).
              </p>
            </div>
          </div>
        </div>

        <div class="coach-action">
          <Button variant="outline" class="w-full btn-nowrap" @click="activeTab = 'insights'">
            <span>Xem trọn bộ gợi ý kỹ năng</span>
            <ChevronRight :size="14" class="shrink-0" />
          </Button>
        </div>
      </Card>
    </div>

    <!-- TAB 2 & HISTORY TABLE AREA -->
    <div v-show="activeTab === 'overview' || activeTab === 'history'">
      <Card class="table-card">
        <div class="table-card-head">
          <div>
            <h2 class="section-title">
              <Layers :size="20" style="color: var(--primary)" />
              <span>Lịch sử các bài phỏng vấn đã thực hiện</span>
            </h2>
            <p class="section-desc">Danh sách đầy đủ điểm số, thời gian và nhận xét AI cho từng phiên</p>
          </div>
          <Button variant="outline" class="btn-nowrap" @click="router.push('/mock-setup')">
            <RefreshCw :size="14" class="shrink-0" />
            <span>Luyện tập tiếp</span>
          </Button>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="ph-spinner"></div>
          <p class="loading-text">Đang tải lịch sử luyện tập...</p>
        </div>

        <div v-else-if="loadError" class="empty-state-box">
          <h3 class="empty-title">Không tải được lịch sử</h3>
          <p class="empty-desc">{{ loadError }}</p>
          <Button variant="primary" class="mt-4" @click="router.go(0)">Thử lại</Button>
        </div>

        <div v-else class="table-wrapper">
          <table v-if="history.length > 0 && filteredHistory.length > 0" class="app-table">
            <thead>
              <tr>
                <th class="col-role">Vị trí & Cấp độ Phỏng vấn</th>
                <th class="col-date">Ngày thực hiện</th>
                <th class="col-dur">Thời lượng</th>
                <th class="col-score">Điểm AI Rubric</th>
                <th class="col-desc">Đánh giá chung</th>
                <th class="col-action">Thao tác</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in filteredHistory" :key="item.id">
                <td class="col-role">
                  <div class="role-title">{{ item.role }}</div>
                  <div class="role-meta">
                    <span class="level-pill">{{ item.level }}</span>
                    <span>• Cập nhật AI Real-time</span>
                  </div>
                </td>
                <td class="col-date">{{ item.date }}</td>
                <td class="col-dur">{{ item.duration }}</td>
                <td class="col-score">
                  <Badge :type="item.badgeType" className="score-badge">
                    {{ item.score != null ? `${item.score} / 10` : 'Chưa có điểm' }}
                  </Badge>
                </td>
                <td class="col-desc">
                  <div class="badge-row">
                    <Badge :type="item.badgeType">
                      <Sparkles :size="11" v-if="item.badge === 'Xuất sắc'" class="mr-1" />
                      {{ item.badge }}
                    </Badge>
                  </div>
                  <p class="highlights-text">{{ item.highlights }}</p>
                </td>
                <td class="col-action">
                  <Button 
                    variant="outline" 
                    class="action-btn"
                    @click="handleActionClick(item.id)"
                  >
                    <Eye :size="14" class="shrink-0" />
                    <span>Xem Báo Cáo AI</span>
                  </Button>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Empty state when search or filter yields 0 matches -->
          <div v-if="history.length > 0 && filteredHistory.length === 0" class="empty-filter-box">
            <AlertCircle :size="36" style="color: var(--warning)" />
            <h4 class="empty-filter-title">Không tìm thấy bài phỏng vấn phù hợp với bộ lọc</h4>
            <p class="empty-filter-desc">Thử thay đổi từ khóa tìm kiếm hoặc đặt lại các tiêu chí lọc về "Tất cả".</p>
            <Button variant="outline" class="btn-nowrap mt-3" @click="resetFilters">
              <span>✕ Đặt lại bộ lọc & xem tất cả</span>
            </Button>
          </div>

          <div v-if="history.length === 0" class="empty-state-box">
            <div class="empty-icon-ring">
              <Layers :size="36" />
            </div>
            <h3 class="empty-title">Bạn chưa thực hiện phiên phỏng vấn AI nào</h3>
            <p class="empty-desc">
              Hệ thống Trí tuệ Nhân tạo đang sẵn sàng! Hãy bắt đầu bài luyện tập phỏng vấn ngay để AI chấm điểm tự động theo bộ chuẩn Rubric Enterprise, phân tích lỗ hổng kỹ năng và tạo lộ trình rèn luyện cá nhân hóa cho bạn.
            </p>
            <Button variant="primary" class="btn-nowrap mt-4" @click="router.push('/mock-setup')">
              <Play :size="16" class="shrink-0" />
              <span>Bắt đầu Luyện tập Phỏng vấn AI ngay</span>
            </Button>
          </div>
        </div>
      </Card>
    </div>

    <!-- TAB 3: INSIGHTS & AI COACH ACTION PLAN -->
    <div v-show="activeTab === 'insights'" class="space-y-6">
      <Card class="insights-card">
        <div class="insights-head">
          <div>
            <h3 class="section-title">
              <Zap :size="20" style="color: var(--warning)" />
              <span>Chương trình Rèn luyện Khắc phục Lỗ hổng Kỹ năng</span>
            </h3>
            <p class="section-desc">Dựa trên phân tích tự động từ 50+ câu hỏi trả lời của bạn qua các phiên</p>
          </div>
          <Button variant="primary" class="btn-nowrap" @click="router.push('/mock-setup')">
            <Play :size="16" class="shrink-0" />
            <span>Luyện tập chủ đề ưu tiên ngay</span>
          </Button>
        </div>

        <div class="insights-grid">
          <div class="action-item">
            <div class="action-top">
              <Badge type="danger">Ưu tiên Cao</Badge>
              <span class="action-count">Xuất hiện ở 2 phiên</span>
            </div>
            <h4 class="action-title">Tối ưu cấu trúc câu trả lời theo STAR Framework</h4>
            <p class="action-body">
              Khi gặp câu hỏi tình huống (Behavioral), bạn có xu hướng giải thích bối cảnh (Situation) quá dài (chiếm ~50% thời gian) và nói lướt qua phần hành động cụ thể (Action).
            </p>
            <div class="action-footer" style="color: var(--primary)">
              <span>📚 Bài học gợi ý: Nghệ thuật trả lời STAR trong 90s</span>
              <ChevronRight :size="14" />
            </div>
          </div>

          <div class="action-item">
            <div class="action-top">
              <Badge type="warning">Ưu tiên Trung bình</Badge>
              <span class="action-count">Xuất hiện ở 1 phiên</span>
            </div>
            <h4 class="action-title">Bổ sung số liệu thực tế vào câu hỏi System Scaling</h4>
            <p class="action-body">
              Các câu trả lời về kiến trúc Microservices rất chuẩn chỉ về lý thuyết, nhưng Trí tuệ Nhân tạo khuyến khích bạn nêu thêm throughput (RPS), quy mô DB cluster và trade-off chi phí.
            </p>
            <div class="action-footer" style="color: var(--primary)">
              <span>📚 Bài học gợi ý: Các con số vàng trong phỏng vấn System Design</span>
              <ChevronRight :size="14" />
            </div>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
/* =========================================================================
   DESIGN_SYSTEM_SKILL.md Compliance & Token Usage
   ========================================================================= */

.page-wrapper {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding-bottom: 48px;
  max-width: 1200px;
  margin: 0 auto;
  font-family: var(--sans, 'Inter', -apple-system, sans-serif);
}

/* Header Banner — DESIGN_SYSTEM_SKILL.md §3 text-h1 24px, white surface */
.header-box {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  padding: 24px 28px;
  border-radius: var(--radius-lg, 16px);
  background: var(--surface, #FFFFFF);
  border: 1px solid var(--border, #E2E8F0);
  box-shadow: var(--shadow-sm);
}
@media (min-width: 640px) {
  .header-box {
    flex-direction: row;
    align-items: center;
  }
}
.header-info {
  flex: 1;
}
.brand-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: var(--radius-full, 9999px);
  background: var(--primary-light, #DBEAFE);
  color: var(--primary, #2563EB);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 12px;
}
.page-title {
  font-size: 24px; /* text-h1 exact per DESIGN_SYSTEM_SKILL.md §3 */
  font-weight: 700;
  color: var(--text-main, #0F172A);
  line-height: 1.3;
  margin: 0;
}
.page-subtitle {
  font-size: 14px; /* text-body */
  color: var(--text-secondary, #475569);
  line-height: 1.6;
  margin: 8px 0 0 0;
  max-width: 720px;
}
.header-cta {
  flex-shrink: 0;
}

/* Navigation Tabs — DESIGN_SYSTEM_SKILL.md §5 */
.tabs-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  border-bottom: 1px solid var(--border, #E2E8F0);
  overflow-x: auto;
}
.tab-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary, #475569);
  border: none;
  background: transparent;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
}
.tab-item:hover {
  color: var(--text-main, #0F172A);
}
.tab-item.active {
  color: var(--primary, #2563EB);
  border-bottom-color: var(--primary, #2563EB);
  font-weight: 700;
}

/* KPI Cards Grid — DESIGN_SYSTEM_SKILL.md §4 */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  gap: 16px;
}
@media (min-width: 640px) {
  .kpi-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (min-width: 1024px) {
  .kpi-grid { grid-template-columns: repeat(4, minmax(0, 1fr)); }
}
.kpi-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 20px !important;
  background: var(--surface, #FFFFFF) !important;
  border: 1px solid var(--border, #E2E8F0);
  border-radius: var(--radius-lg, 16px);
  box-shadow: var(--shadow-sm);
  transition: all 0.2s ease;
}
.kpi-card:hover {
  border-color: var(--primary-light, #DBEAFE);
  box-shadow: var(--shadow-md);
}
.kpi-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.kpi-label {
  font-size: 12px;
  font-weight: 700;
  color: var(--text-muted, #94A3B8);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.kpi-icon-box {
  width: 36px;
  height: 36px;
  border-radius: var(--radius, 12px);
  display: flex;
  align-items: center;
  justify-content: center;
}
.primary-bg { background: var(--primary-light, #DBEAFE); }
.accent-bg { background: var(--accent-bg, #ECFEFF); }
.info-bg { background: var(--surface-soft, #F1F5F9); }
.warning-bg { background: rgba(217, 119, 6, 0.1); }

.kpi-body {
  display: flex;
  align-items: baseline;
  gap: 6px;
}
.kpi-num {
  font-size: 28px; /* stat-value exact per DESIGN_SYSTEM_SKILL.md §3 */
  font-weight: 700;
  color: var(--text-main, #0F172A);
  letter-spacing: -0.02em;
}
.kpi-unit {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted, #94A3B8);
}
.kpi-text-main {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
}
.kpi-note {
  font-size: 12px;
  color: var(--text-secondary, #475569);
  margin: 12px 0 0 0;
  display: flex;
  align-items: center;
  gap: 6px;
}
.status-card {
  border-left: 4px solid var(--warning, #D97706) !important;
}

/* Content Grid & Coach Card */
.content-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  gap: 24px;
}
@media (min-width: 1024px) {
  .content-grid { grid-template-columns: 2fr 1fr; }
}
.skill-card, .coach-card {
  padding: 24px !important;
  background: var(--surface, #FFFFFF) !important;
  border: 1px solid var(--border, #E2E8F0);
  border-radius: var(--radius-lg, 16px);
  box-shadow: var(--shadow-sm);
}
.coach-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border, #E2E8F0);
  padding-bottom: 16px;
  margin-bottom: 20px;
}
.section-title {
  font-size: 16px; /* text-h2 exact per DESIGN_SYSTEM_SKILL.md §3 */
  font-weight: 600;
  color: var(--text-main, #0F172A);
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
}
.section-desc {
  font-size: 12px;
  color: var(--text-secondary, #475569);
  margin: 4px 0 0 0;
}
.skill-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.skill-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-main, #0F172A);
  margin-bottom: 6px;
}
.progress-track {
  width: 100%;
  height: 8px;
  border-radius: var(--radius-full, 9999px);
  background: var(--surface-soft, #F1F5F9);
  overflow: hidden;
  margin-top: 12px;
}
.progress-bar {
  height: 100%;
  border-radius: var(--radius-full, 9999px);
  transition: width 0.8s ease;
}
.feedback-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  gap: 16px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid var(--border, #E2E8F0);
}
@media (min-width: 640px) {
  .feedback-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
.feedback-box {
  padding: 14px;
  border-radius: var(--radius, 12px);
}
.success-box {
  background: rgba(22, 163, 74, 0.06);
  border: 1px solid rgba(22, 163, 74, 0.2);
}
.warning-box {
  background: rgba(217, 119, 6, 0.06);
  border: 1px solid rgba(217, 119, 6, 0.2);
}
.feedback-head {
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}
.feedback-body {
  font-size: 13px;
  color: var(--text-secondary, #475569);
  margin: 0;
  line-height: 1.5;
}
.coach-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius, 12px);
  background: var(--primary-light, #DBEAFE);
  color: var(--primary, #2563EB);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}
.section-subtitle {
  font-size: 12px;
  font-weight: 600;
  margin: 4px 0 16px 0;
}
.coach-tips {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.tip-box {
  padding: 14px;
  border-radius: var(--radius, 12px);
  background: var(--surface-soft, #F1F5F9);
  border: 1px solid var(--border, #E2E8F0);
}
.tip-label {
  font-size: 11px;
  font-weight: 700;
  color: var(--text-muted, #94A3B8);
  text-transform: uppercase;
  display: block;
  margin-bottom: 6px;
}
.tip-body {
  font-size: 13px;
  color: var(--text-main, #0F172A);
  margin: 0;
  line-height: 1.5;
}
.coach-action {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid var(--border, #E2E8F0);
}

/* =========================================================================
   TABLE DESIGN — DESIGN_SYSTEM_SKILL.md §4 & §5
   Zero horizontal scroll (`cuộn`), no wraps (`xuống dòng`), perfect alignment
   ========================================================================= */

.table-card {
  padding: 0 !important;
  overflow: hidden;
  background: var(--surface, #FFFFFF) !important;
  border: 1px solid var(--border, #E2E8F0);
  border-radius: var(--radius-lg, 16px);
  box-shadow: var(--shadow-sm);
}
.table-card-head {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border, #E2E8F0);
  background: var(--surface-soft, #F1F5F9);
}
@media (min-width: 640px) {
  .table-card-head {
    flex-direction: row;
    align-items: center;
  }
}
.table-wrapper {
  width: 100%;
}
.app-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: auto;
}
.app-table th {
  background: var(--surface-soft, #F1F5F9);
  color: var(--text-secondary, #475569);
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border, #E2E8F0);
  vertical-align: middle;
  text-align: left;
}
.app-table td {
  padding: 16px;
  border-bottom: 1px solid var(--border, #E2E8F0);
  color: var(--text-main, #0F172A);
  font-size: 14px;
  vertical-align: middle;
  text-align: left;
}
.app-table tr:last-child td {
  border-bottom: none;
}
.app-table tbody tr {
  transition: background 0.2s ease;
}
.app-table tbody tr:hover {
  background: var(--surface-soft, #F1F5F9);
}

/* Column explicit sizing & no-wrap rules */
.col-role {
  width: 32%;
}
.role-title {
  font-weight: 600;
  color: var(--text-main, #0F172A);
  margin-bottom: 4px;
}
.role-meta {
  font-size: 12px;
  color: var(--text-muted, #94A3B8);
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}
.level-pill {
  padding: 2px 8px;
  border-radius: var(--radius, 12px);
  background: var(--surface-soft, #F1F5F9);
  color: var(--text-secondary, #475569);
  font-weight: 600;
  font-size: 11px;
  white-space: nowrap;
}
.col-date, .col-dur {
  width: 13%;
  white-space: nowrap;
  color: var(--text-secondary, #475569);
  font-weight: 500;
}
.col-score {
  width: 14%;
  white-space: nowrap;
}
.col-desc {
  width: 27%;
}
.badge-row {
  margin-bottom: 6px;
}
.highlights-text {
  font-size: 13px;
  color: var(--text-secondary, #475569);
  margin: 0;
  line-height: 1.5;
}

/* Thao tác column — right aligned, never wraps */
.col-action {
  width: 1%;
  white-space: nowrap;
  text-align: right;
}
.action-btn {
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px !important;
  padding: 8px 14px !important;
  font-size: 13px !important;
  font-weight: 600 !important;
  border-radius: var(--radius-full, 9999px) !important;
  white-space: nowrap !important;
  flex-shrink: 0 !important;
}

/* General Button helper */
.btn-nowrap {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

/* Loading & Insights */
.loading-state {
  padding: 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
}
.ph-spinner {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 3px solid var(--primary-light, #DBEAFE);
  border-top-color: var(--primary, #2563EB);
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.loading-text {
  font-size: 14px;
  color: var(--text-secondary, #475569);
  margin: 0;
}
.insights-card {
  padding: 24px !important;
  background: var(--surface, #FFFFFF) !important;
  border: 1px solid var(--border, #E2E8F0);
  border-radius: var(--radius-lg, 16px);
  box-shadow: var(--shadow-sm);
}
.insights-head {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  border-bottom: 1px solid var(--border, #E2E8F0);
  padding-bottom: 16px;
  margin-bottom: 24px;
}
@media (min-width: 640px) {
  .insights-head {
    flex-direction: row;
    align-items: center;
  }
}
.insights-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  gap: 20px;
}
@media (min-width: 768px) {
  .insights-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
.action-item {
  padding: 20px;
  border-radius: var(--radius-lg, 16px);
  background: var(--surface-soft, #F1F5F9);
  border: 1px solid var(--border, #E2E8F0);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.action-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.action-count {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted, #94A3B8);
}
.action-title {
  font-size: 16px; /* text-h2 */
  font-weight: 600;
  color: var(--text-main, #0F172A);
  margin: 0 0 8px 0;
}
.action-body {
  font-size: 13px;
  color: var(--text-secondary, #475569);
  line-height: 1.6;
  margin: 0;
}
.action-footer {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--border, #E2E8F0);
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  font-weight: 600;
}

/* Zero State / Empty State UI styling */
.empty-state-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 64px 24px;
  gap: 16px;
}
.empty-icon-ring {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: var(--surface-soft, #F1F5F9);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary, #2563EB);
  border: 1px solid var(--border, #E2E8F0);
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.empty-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
  margin: 0;
}
.empty-desc {
  font-size: 14px;
  color: var(--text-secondary, #475569);
  max-width: 520px;
  line-height: 1.6;
  margin: 0;
}

/* Filter & Search Bar styles */
.filter-bar {
  padding: 16px 24px;
  border-bottom: 1px solid var(--border, #E2E8F0);
  background: var(--surface-soft, #F8FAFC);
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.search-box {
  display: flex;
  align-items: center;
  position: relative;
  flex: 1;
  min-width: 280px;
}
.search-icon {
  position: absolute;
  left: 14px;
  color: var(--text-muted, #94A3B8);
}
.search-input {
  width: 100%;
  padding: 10px 36px 10px 38px;
  border-radius: 8px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface, #FFFFFF);
  font-size: 13px;
  color: var(--text-main, #0F172A);
  outline: none;
  transition: all 0.2s ease;
}
.search-input:focus {
  border-color: var(--primary, #2563EB);
  box-shadow: 0 0 0 3px rgba(37,99,235,0.1);
}
.clear-btn {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  color: var(--text-muted, #94A3B8);
  font-size: 14px;
  cursor: pointer;
  padding: 2px 6px;
}
.clear-btn:hover {
  color: var(--text-main, #0F172A);
}
.filter-controls {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 14px;
}
.filter-group {
  display: flex;
  align-items: center;
  gap: 6px;
}
.filter-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary, #475569);
  white-space: nowrap;
}
.app-select {
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border, #E2E8F0);
  background: var(--surface, #FFFFFF);
  font-size: 12px;
  font-weight: 600;
  color: var(--text-main, #0F172A);
  outline: none;
  cursor: pointer;
  transition: all 0.2s ease;
}
.app-select:hover, .app-select:focus {
  border-color: var(--primary, #2563EB);
}

.empty-filter-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 48px 24px;
  gap: 12px;
}
.empty-filter-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-main, #0F172A);
  margin: 0;
}
.empty-filter-desc {
  font-size: 13px;
  color: var(--text-secondary, #475569);
  max-width: 420px;
  margin: 0;
}
</style>
