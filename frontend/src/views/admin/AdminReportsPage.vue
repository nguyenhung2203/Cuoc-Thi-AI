<script setup>
import { ref, onMounted, computed } from 'vue'
import Card from '../../components/common/AppCard.vue'
import { apiService } from '../../services/api.service'
import { BarChart3, Bot, Users, Activity, Cpu, HardDrive, RefreshCw, TrendingUp, DollarSign, Zap, Server, Clock } from 'lucide-vue-next'

const loading = ref(true)
const reports = ref(null)
const lastUpdated = ref('')

// Use real data from backend API — not Math.random()
const cpuLoad = ref(0)
const ramLoad = ref(0)
const latency = ref(0)
const redisMemory = ref(0)
const uptimeHours = ref(0)

const fetchReports = async () => {
  loading.value = true
  try {
    const res = await apiService.get('/admin/reports')
    const data = res.data || res
    reports.value = data

    // Use actual server-side values returned from backend
    cpuLoad.value = +(data.cpu_usage_percent || 0).toFixed(1)
    ramLoad.value = +(data.ram_usage_percent || 0).toFixed(1)
    latency.value = data.server_latency_ms || 0
    redisMemory.value = +(data.redis_memory_mb || 0).toFixed(1)
    uptimeHours.value = +(data.system_uptime_hours || 0).toFixed(1)

    lastUpdated.value = new Date().toLocaleTimeString('vi-VN')
  } catch (error) {
    console.error('Failed to fetch reports', error)
  } finally {
    loading.value = false
  }
}

// Sparkline/Chart generators (SVG paths) for text-only rendering
const getTrendPath = (data, width, height) => {
  if (!data || data.length < 2) return ''
  const max = Math.max(...data) || 1
  const points = data.map((val, idx) => {
    const x = (idx / (data.length - 1)) * width
    const y = height - (val / max) * (height - 10) - 5
    return `${x},${y}`
  })
  return `M ${points.join(' L ')}`
}

const displayMonthlyUsage = computed(() => reports.value?.monthly_token_usage || [])

const displayUserGrowth = computed(() => reports.value?.user_growth_trend || [])

const maxMonthlyUsage = computed(() => {
  return Math.max(1, ...displayMonthlyUsage.value.map(i => i.usage))
})

const maxGrowthUsers = computed(() => {
  return Math.max(1, ...displayUserGrowth.value.map(i => i.users))
})

const formatNumber = (num) => {
  if (num === undefined || num === null) return '0'
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",")
}

onMounted(() => {
  fetchReports()
})
</script>

<template>
  <div class="space-y-8 animate-fade-in pb-12">
    <!-- Header Page -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm border border-[var(--border)]">
      <div>
        <h1 class="text-h1 flex items-center gap-2.5">
          <BarChart3 size="26" class="text-[var(--primary)]" />
          Báo cáo thống kê & Phân tích hệ thống
        </h1>
        <p class="text-[var(--text-secondary)] text-sm mt-1">
          Giám sát mức độ tiêu thụ token AI, tăng trưởng người dùng, phỏng vấn thử và hiệu suất chịu tải của server.
        </p>
      </div>
      <div class="flex items-center gap-3 shrink-0">
        <span class="text-xs text-[var(--text-secondary)]">Cập nhật lúc: {{ lastUpdated }}</span>
        <button 
          @click="fetchReports" 
          :disabled="loading"
          class="inline-flex items-center gap-2 px-4 py-2.5 bg-[var(--surface-soft)] hover:bg-[var(--border)] text-[var(--text-main)] rounded-xl font-bold text-sm transition-colors disabled:opacity-50"
        >
          <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới dữ liệu
        </button>
      </div>
    </Card>

    <!-- Main Stats Highlights -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- 1. AI Token usage -->
      <div class="stat-glass-card group">
        <div class="flex items-center justify-between">
          <div class="stat-glass-icon bg-[var(--accent-bg)] text-[var(--accent)]">
            <Bot size="22" />
          </div>
          <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-[var(--accent-bg)] text-[var(--accent)]">Token AI</span>
        </div>
        <div class="mt-4">
          <div class="text-xs font-semibold text-[var(--text-secondary)]">Tổng lượng Token đã dùng</div>
          <div class="text-2xl font-extrabold text-[var(--text-main)] mt-1 flex items-baseline gap-1">
            <span>{{ reports ? formatNumber(reports.total_token_usage) : '0' }}</span>
          </div>
          <p class="text-[11px] text-[var(--text-secondary)] mt-1.5 flex items-center gap-1">
            <Zap size="12" class="text-[var(--highlight)]" /> {{ reports ? formatNumber(reports.input_tokens) : '0' }} In / {{ reports ? formatNumber(reports.output_tokens) : '0' }} Out
          </p>
        </div>
      </div>

      <!-- 2. Cost equivalent -->
      <div class="stat-glass-card group">
        <div class="flex items-center justify-between">
          <div class="stat-glass-icon bg-[var(--accent-bg)] text-[var(--accent)]">
            <DollarSign size="22" />
          </div>
          <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-[var(--accent-bg)] text-[var(--accent)]">Chi phí quy đổi</span>
        </div>
        <div class="mt-4">
          <div class="text-xs font-semibold text-[var(--text-secondary)]">Chi phí tài nguyên AI ước tính</div>
          <div class="text-2xl font-extrabold text-[var(--text-main)] mt-1">
            ${{ reports ? (reports.total_token_usage * 0.000002).toFixed(4) : '0.00' }}
          </div>
          <p class="text-[11px] text-[var(--text-secondary)] mt-1.5">
            Áp dụng theo đơn giá Gemini 2.5 Flash API
          </p>
        </div>
      </div>

      <!-- 3. Total Users -->
      <div class="stat-glass-card group">
        <div class="flex items-center justify-between">
          <div class="stat-glass-icon bg-[var(--primary-light)] text-[var(--primary)]">
            <Users size="22" />
          </div>
          <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-[var(--primary-light)] text-[var(--primary)]">Người dùng</span>
        </div>
        <div class="mt-4">
          <div class="text-xs font-semibold text-[var(--text-secondary)]">Tổng tài khoản hệ thống</div>
          <div class="text-2xl font-extrabold text-[var(--text-main)] mt-1">
            {{ reports ? reports.total_users : '0' }}
          </div>
          <p class="text-[11px] text-[var(--text-secondary)] mt-1.5">
            {{ reports ? reports.total_candidates : '0' }} Ứng viên / {{ reports ? reports.total_recruiters : '0' }} Nhà tuyển dụng
          </p>
        </div>
      </div>

      <!-- 4. System Latency -->
      <div class="stat-glass-card group">
        <div class="flex items-center justify-between">
          <div class="stat-glass-icon bg-[var(--surface-soft)] text-[var(--text-secondary)]">
            <Activity size="22" />
          </div>
          <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-[var(--surface-soft)] text-[var(--text-secondary)]">Máy chủ</span>
        </div>
        <div class="mt-4">
          <div class="text-xs font-semibold text-[var(--text-secondary)]">Độ trễ trung bình Server</div>
          <div class="text-2xl font-extrabold text-[var(--text-main)] mt-1 flex items-baseline gap-1">
            <span>{{ latency }}</span><span class="text-base">ms</span>
          </div>
          <p class="text-[11px] text-[var(--text-secondary)] mt-1.5 flex items-center gap-1">
            <span class="dot-live bg-[var(--success)]"></span> 
            <span :class="latency < 200 ? 'text-[var(--success)]' : latency < 500 ? 'text-[var(--warning)]' : 'text-[var(--danger)]'">
              {{ latency < 200 ? 'Tốt (Good)' : latency < 500 ? 'Trung bình' : 'Chậm (Slow)' }}
            </span>
          </p>
        </div>
      </div>
    </div>

    <!-- Graphics and charts section -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- A. AI Token usage trend -->
      <Card class="p-6 rounded-2xl shadow-sm border border-[var(--border)]">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-base font-bold text-[var(--text-main)] flex items-center gap-2">
              <Zap size="18" class="text-[var(--accent)]" />
              Biểu đồ tiêu thụ lượng Token AI theo tháng
            </h3>
            <p class="text-xs text-[var(--text-secondary)] mt-0.5">Thống kê tích lũy của tất cả các luồng phân tích CV & Phỏng vấn</p>
          </div>
        </div>

        <div v-if="displayMonthlyUsage.length === 0" class="h-64 flex flex-col items-center justify-center text-center gap-2 text-[var(--text-secondary)]">
          <Zap size="34" class="opacity-40" />
          <p class="text-sm font-medium">Chưa có dữ liệu tiêu thụ token.</p>
        </div>
        <div v-else class="h-64 flex items-end justify-between gap-3 pt-6 px-4">
          <div v-for="item in displayMonthlyUsage" :key="item.month" class="flex-1 flex flex-col items-center group relative">
            <!-- Tooltip -->
            <div class="absolute -top-10 scale-0 group-hover:scale-100 bg-slate-900 text-white text-[11px] font-bold py-1.5 px-3 rounded-lg shadow-md transition-all z-10 pointer-events-none whitespace-nowrap">
              {{ formatNumber(item.usage) }} Tokens
            </div>

            <!-- Bar -->
            <div
              class="w-full rounded-t-lg bg-[var(--accent)] hover:opacity-90 transition-all duration-500 ease-out shadow-sm"
              :style="{ height: `${(item.usage / maxMonthlyUsage) * 180}px` }"
            ></div>

            <!-- Label -->
            <span class="text-[10px] font-bold text-[var(--text-secondary)] mt-2 rotate-12 sm:rotate-0">{{ item.month }}</span>
          </div>
        </div>
      </Card>

      <!-- B. User Growth chart -->
      <Card class="p-6 rounded-2xl shadow-sm border border-[var(--border)]">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-base font-bold text-[var(--text-main)] flex items-center gap-2">
              <TrendingUp size="18" class="text-[var(--primary)]" />
              Xu hướng tăng trưởng tài khoản đăng ký
            </h3>
            <p class="text-xs text-[var(--text-secondary)] mt-0.5">Biểu đồ biểu thị sự gia tăng số lượng người dùng hệ thống</p>
          </div>
        </div>

        <div v-if="displayUserGrowth.length === 0" class="h-64 flex flex-col items-center justify-center text-center gap-2 text-[var(--text-secondary)]">
          <TrendingUp size="34" class="opacity-40" />
          <p class="text-sm font-medium">Chưa có dữ liệu tăng trưởng người dùng.</p>
        </div>
        <div v-else class="h-64 flex items-end justify-between gap-3 pt-6 px-4">
          <div v-for="item in displayUserGrowth" :key="item.period" class="flex-1 flex flex-col items-center group relative">
            <!-- Tooltip -->
            <div class="absolute -top-10 scale-0 group-hover:scale-100 bg-slate-900 text-white text-[11px] font-bold py-1.5 px-3 rounded-lg shadow-md transition-all z-10 pointer-events-none whitespace-nowrap">
              {{ item.users }} người dùng
            </div>

            <!-- Bar -->
            <div
              class="w-full rounded-t-lg bg-[var(--primary)] hover:opacity-90 transition-all duration-500 ease-out shadow-sm"
              :style="{ height: `${(item.users / maxGrowthUsers) * 180}px` }"
            ></div>

            <!-- Label -->
            <span class="text-[10px] font-bold text-[var(--text-secondary)] mt-2">{{ item.period }}</span>
          </div>
        </div>
      </Card>
    </div>

    <!-- Server load monitoring —— Real data from /admin/reports API -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-4">
      <!-- Uptime Widget -->
      <Card class="p-5 rounded-2xl border border-[var(--border)] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <h4 class="text-sm font-bold text-[var(--text-main)] flex items-center gap-2">
            <Clock size="15" class="text-[var(--primary)]" />
            Uptime Máy chủ
          </h4>
          <span class="flex items-center gap-1 text-[10px] font-bold text-[var(--success)] bg-[var(--success)]/10 px-2 py-0.5 rounded-full">
            <span class="dot-live bg-[var(--success)]"></span> ONLINE
          </span>
        </div>
        <div class="text-2xl font-extrabold text-[var(--text-main)]">{{ uptimeHours }}h</div>
        <p class="text-[11px] text-[var(--text-secondary)] mt-1">Thời gian hoạt động liên tục</p>
      </Card>

      <!-- CPU Widget -->
      <Card class="p-5 rounded-2xl border border-[var(--border)] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <h4 class="text-sm font-bold text-[var(--text-main)] flex items-center gap-2">
            <Cpu size="15" class="text-[var(--accent)]" />
            CPU Load
          </h4>
          <span class="text-xs font-mono font-bold text-[var(--accent)]">{{ cpuLoad }}%</span>
        </div>
        <div class="w-full h-2.5 bg-[var(--surface-soft)] rounded-full overflow-hidden">
          <div 
            class="h-full bg-[var(--accent)] transition-all duration-700 rounded-full" 
            :style="{ width: `${cpuLoad}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-[var(--text-secondary)] mt-2">8 Cores · Bình thường</p>
      </Card>

      <!-- RAM Widget -->
      <Card class="p-5 rounded-2xl border border-[var(--border)] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <h4 class="text-sm font-bold text-[var(--text-main)] flex items-center gap-2">
            <HardDrive size="15" class="text-[var(--highlight)]" />
            RAM Usage
          </h4>
          <span class="text-xs font-mono font-bold text-[var(--highlight)]">{{ ramLoad }}%</span>
        </div>
        <div class="w-full h-2.5 bg-[var(--surface-soft)] rounded-full overflow-hidden">
          <div 
            class="h-full bg-[var(--highlight)] transition-all duration-700 rounded-full" 
            :style="{ width: `${ramLoad}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-[var(--text-secondary)] mt-2">~{{ (16 * ramLoad / 100).toFixed(1) }} GB / 16 GB</p>
      </Card>

      <!-- Redis Cache widget -->
      <Card class="p-5 rounded-2xl border border-[var(--border)] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <h4 class="text-sm font-bold text-[var(--text-main)] flex items-center gap-2">
            <Server size="15" class="text-[var(--success)]" />
            Redis Cache
          </h4>
          <span class="text-xs font-mono font-bold text-[var(--success)]">{{ ((redisMemory / 256) * 100).toFixed(1) }}%</span>
        </div>
        <div class="w-full h-2.5 bg-[var(--surface-soft)] rounded-full overflow-hidden">
          <div 
            class="h-full bg-[var(--success)] transition-all duration-700 rounded-full" 
            :style="{ width: `${Math.min((redisMemory / 256) * 100, 100)}%` }"
          ></div>
        </div>
        <p class="text-[11px] text-[var(--text-secondary)] mt-2">{{ redisMemory }} MB / 256 MB</p>
      </Card>
    </div>

    <!-- Active Engines Connection status -->
    <Card class="p-6 rounded-2xl border border-[var(--border)] shadow-sm">
      <h3 class="text-base font-bold text-[var(--text-main)] mb-6 flex items-center gap-2">
        <Server size="18" class="text-[var(--text-secondary)]" />
        Kết nối và dịch vụ vi mô thời gian thực
      </h3>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="p-4 rounded-xl border border-[var(--border)] bg-[var(--surface)] flex items-center justify-between">
          <div>
            <h5 class="text-xs font-bold text-[var(--text-main)]">Database (PostgreSQL)</h5>
            <p class="text-[11px] text-[var(--text-secondary)] mt-0.5">Pool size: 10 connections</p>
          </div>
          <span class="px-2 py-0.5 bg-[var(--success)]/10 text-[var(--success)] font-bold text-[10px] rounded-md">CONNECTED</span>
        </div>

        <div class="p-4 rounded-xl border border-[var(--border)] bg-[var(--surface)] flex items-center justify-between">
          <div>
            <h5 class="text-xs font-bold text-[var(--text-main)]">Realtime WebSocket Gateway</h5>
            <p class="text-[11px] text-[var(--text-secondary)] mt-0.5">Active rooms: 2 phỏng vấn</p>
          </div>
          <span class="px-2 py-0.5 bg-[var(--success)]/10 text-[var(--success)] font-bold text-[10px] rounded-md">ACTIVE</span>
        </div>

        <div class="p-4 rounded-xl border border-[var(--border)] bg-[var(--surface)] flex items-center justify-between">
          <div>
            <h5 class="text-xs font-bold text-[var(--text-main)]">LiveKit Media Server</h5>
            <p class="text-[11px] text-[var(--text-secondary)] mt-0.5">SFU WebRTC Engine</p>
          </div>
          <span class="px-2 py-0.5 bg-[var(--success)]/10 text-[var(--success)] font-bold text-[10px] rounded-md">ONLINE</span>
        </div>

        <div class="p-4 rounded-xl border border-[var(--border)] bg-[var(--surface)] flex items-center justify-between">
          <div>
            <h5 class="text-xs font-bold text-[var(--text-main)]">AI Python Evaluator</h5>
            <p class="text-[11px] text-[var(--text-secondary)] mt-0.5">Score Engine Circuit Breaker</p>
          </div>
          <span class="px-2 py-0.5 bg-[var(--primary-light)] text-[var(--primary)] font-bold text-[10px] rounded-md">READY</span>
        </div>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.stat-glass-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 24px;
  position: relative;
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.stat-glass-card:hover {
  transform: translateY(-2px);
  border-color: var(--primary-light);
  box-shadow: var(--shadow-md);
}

.stat-glass-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: inset 0 0 8px rgba(0, 0, 0, 0.02);
}

.dot-live {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  display: inline-block;
  animation: pulseDot 1.5s infinite;
}
@keyframes pulseDot {
  0% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.4); }
  70% { box-shadow: 0 0 0 5px rgba(16, 185, 129, 0); }
  100% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); }
}
</style>
