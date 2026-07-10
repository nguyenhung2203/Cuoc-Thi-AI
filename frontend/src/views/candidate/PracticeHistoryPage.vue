<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Table from '../../components/common/AppTable.vue'
import { Target, TrendingUp, Eye, Zap } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'

const router = useRouter()
const history = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    // API_SPEC §11.5 — GET /mock-interviews/my
    const data = await mockService.listMyMockInterviews()
    // Map API response fields sang format hiển thị
    history.value = (Array.isArray(data) ? data : []).map(item => ({
      id: item.id,
      role: item.target_role || 'Không rõ',
      level: item.target_level || '—',
      date: item.created_at ? new Date(item.created_at).toLocaleDateString('vi-VN') : '—',
      score: item.final_score != null ? Number(item.final_score).toFixed(1) : '—',
      status: item.status || 'completed'
    }))
  } catch (err) {
    console.error('Lỗi tải lịch sử mock', err)
    // Giữ nguyên empty array — Backend chưa sẵn sàng
    history.value = []
  } finally {
    loading.value = false
  }
})

const averageScore = computed(() => {
  const scored = history.value.filter(h => h.score !== '—')
  if (scored.length === 0) return 0
  const sum = scored.reduce((acc, curr) => acc + Number(curr.score), 0)
  return (sum / scored.length).toFixed(1)
})

const columns = [
  { header: 'Vị trí luyện tập', key: 'role' },
  { header: 'Cấp độ', key: 'level' },
  { header: 'Ngày thực hiện', key: 'date' },
  { header: 'Điểm số', key: 'score' },
  { header: 'Hành động', key: 'action' }
]
</script>

<template>
  <div class="space-y-8 pb-12 max-w-6xl mx-auto">
    <div class="animate-rise">
      <h1 class="page-title">Kết quả & Lịch sử luyện tập</h1>
      <p class="page-subtitle">Theo dõi sự tiến bộ của bạn qua các bài luyện tập với AI.</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 stagger">
      <Card class="card-elevate tilt-3d stat-card">
        <div class="kpi-icon is-success"><Target :size="22" /></div>
        <div>
          <p class="stat-label">Điểm trung bình</p>
          <h2 class="stat-num">{{ averageScore }}<span class="stat-unit">/10</span></h2>
        </div>
      </Card>

      <Card class="card-elevate tilt-3d stat-card">
        <div class="kpi-icon"><TrendingUp :size="22" /></div>
        <div>
          <p class="stat-label">Số lần luyện tập</p>
          <h2 class="stat-num">{{ history.length }}<span class="stat-unit">lần</span></h2>
        </div>
      </Card>

      <Card class="card-elevate improve-card">
        <p class="stat-label is-warning"><Zap :size="16" /> Điểm cần khắc phục</p>
        <ul class="improve-list">
          <li><span class="dot"></span> Trình bày cấu trúc câu trả lời (STAR)</li>
          <li><span class="dot"></span> Đưa ra thêm nhiều ví dụ thực tế</li>
        </ul>
      </Card>
    </div>

    <!-- History Table Area -->
    <Card class="card-elevate history-card animate-rise">
      <div class="history-head">
        <h2 class="section-heading">Lịch sử bài luyện tập</h2>
      </div>

      <div v-if="loading" class="p-12 text-center flex flex-col items-center">
        <div class="ph-spinner mb-4"></div>
        <p class="text-helper">Đang tải lịch sử...</p>
      </div>

      <div v-else-if="history.length === 0" class="p-12 text-center flex flex-col items-center">
        <div class="kpi-icon" style="width:60px;height:60px;margin-bottom:14px;"><TrendingUp :size="26" /></div>
        <p class="text-secondary-strong">Bạn chưa có bài luyện tập nào. Hãy bắt đầu ngay!</p>
        <Button variant="outline" class="mt-4" @click="router.push('/mock-setup')">Luyện tập với AI</Button>
      </div>

      <div v-else class="overflow-x-auto">
        <div class="modern-table-container">
          <Table :columns="columns" :data="history">
            <template #date="{ row }">
              <span class="cell-muted">{{ row.date }}</span>
            </template>
            <template #score="{ row }">
              <span class="score-pill" :class="row.score >= 8 ? 'is-high' : row.score >= 7 ? 'is-mid' : 'is-low'">
                {{ row.score }}/10
              </span>
            </template>
            <template #action="{ row }">
              <button @click="router.push(`/mock-results/${row.id}`)" class="detail-btn">
                <Eye :size="16" /> Xem chi tiết
              </button>
            </template>
          </Table>
        </div>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.stat-card { display: flex; align-items: center; gap: 16px; padding: 22px; }
.stat-label { display: flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.04em; margin-bottom: 6px; }
.stat-label.is-warning { color: var(--warning); margin-bottom: 12px; }
.stat-num { font-size: 28px; font-weight: 700; color: var(--text-main); display: flex; align-items: baseline; gap: 4px; letter-spacing: -0.02em; }
.stat-unit { font-size: 15px; font-weight: 600; color: var(--text-muted); }
.improve-card { padding: 22px; }
.improve-list { display: flex; flex-direction: column; gap: 8px; }
.improve-list li { display: flex; align-items: flex-start; gap: 8px; color: var(--text-secondary); font-size: 14px; font-weight: 500; }
.improve-list .dot { width: 6px; height: 6px; border-radius: 50%; background: var(--warning); margin-top: 7px; flex-shrink: 0; }

.history-card { padding: 0; overflow: hidden; }
.history-head { padding: 20px 24px; border-bottom: 1px solid var(--border); background: var(--surface-soft); }
.ph-spinner { width: 40px; height: 40px; border-radius: 50%; border: 3px solid var(--primary-light); border-top-color: var(--primary); animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.text-secondary-strong { color: var(--text-secondary); font-weight: 500; }
.cell-muted { color: var(--text-secondary); font-weight: 500; }
.detail-btn { display: inline-flex; align-items: center; gap: 6px; padding: 7px 14px; font-size: 13px; font-weight: 600; color: var(--primary); border: 1px solid var(--primary-light); border-radius: var(--radius); background: transparent; cursor: pointer; transition: all 0.2s ease; }
.detail-btn:hover { background: var(--primary); color: #fff; border-color: var(--primary); }
</style>

<style>
/* Non-scoped: style the AppTable internals rendered inside this page */
.modern-table-container table { width: 100%; border-collapse: separate; border-spacing: 0; }
.modern-table-container th {
  background-color: var(--surface-soft);
  color: var(--text-secondary);
  font-weight: 700; text-transform: uppercase; font-size: 0.72rem; letter-spacing: 0.05em;
  padding: 14px 24px; border-bottom: 1px solid var(--border);
}
.modern-table-container td { padding: 16px 24px; color: var(--text-main); border-bottom: 1px solid var(--surface-soft); font-weight: 500; }
.modern-table-container tr:last-child td { border-bottom: none; }
.modern-table-container tbody tr { transition: background 0.2s ease; }
.modern-table-container tbody tr:hover { background-color: var(--surface-soft); }
</style>
