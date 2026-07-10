<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Users, Building, Video, Activity, ShieldCheck, UserCheck, Clock, ArrowRight, CheckCircle2, AlertCircle, Sparkles } from 'lucide-vue-next'
import { apiService } from '../../services/api.service'

const router = useRouter()
const loading = ref(true)

const stats = ref([
  { title: 'Tổng người dùng', value: 0, icon: Users, colorClass: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-100 dark:border-blue-800/30', badge: '+Mới' },
  { title: 'Doanh nghiệp & Cty', value: 0, icon: Building, colorClass: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-100 dark:border-purple-800/30', badge: 'Hoạt động' },
  { title: 'Lượt phỏng vấn AI', value: 0, icon: Video, colorClass: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-100 dark:border-emerald-800/30', badge: 'Thực tế' },
  { title: 'Tài khoản chờ duyệt', value: 0, icon: Clock, colorClass: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-100 dark:border-amber-800/30', badge: 'Khẩn cấp' },
])

const pendingUsersCount = ref(0)
const pendingUsersList = ref([])

const fetchDashboardData = async () => {
  loading.value = true
  try {
    const data = await apiService.get('/admin/dashboard-stats')
    stats.value[0].value = data.total_users || 0
    stats.value[1].value = data.total_companies || 0
    stats.value[2].value = data.total_interviews || 0
    stats.value[3].value = data.pending_users || 0
    pendingUsersCount.value = data.pending_users || 0

    // Fetch up to 5 pending users for quick approval directly from dashboard
    if (data.pending_users > 0) {
      const pendingRes = await apiService.get('/admin/users/pending')
      pendingUsersList.value = (pendingRes || []).slice(0, 5)
    } else {
      pendingUsersList.value = []
    }
  } catch (error) {
    console.error('Failed to fetch admin stats', error)
  } finally {
    loading.value = false
  }
}

const handleApproveUser = async (userId) => {
  try {
    await apiService.put(`/admin/users/${userId}/approve`)
    // Refresh stats
    fetchDashboardData()
  } catch (error) {
    alert(error.message || 'Phê duyệt thất bại')
  }
}

const formatDate = (isoStr) => {
  if (!isoStr) return 'Vừa xong'
  return new Date(isoStr).toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

onMounted(() => {
  fetchDashboardData()
})
</script>

<template>
  <div class="space-y-8 animate-fade-in">
    <!-- Hero Banner -->
    <div class="relative overflow-hidden rounded-3xl bg-gradient-to-r from-slate-900 via-indigo-950 to-slate-900 p-8 text-white shadow-xl border border-slate-800">
      <div class="absolute -right-10 -bottom-10 opacity-10 pointer-events-none">
        <ShieldCheck size="240" />
      </div>
      <div class="relative z-10 max-w-3xl space-y-3">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-500/20 border border-indigo-400/30 text-indigo-300 text-xs font-semibold uppercase tracking-wider">
          <Sparkles size="14" /> Trung tâm điều khiển WeMake AI
        </div>
        <h1 class="text-3xl sm:text-4xl font-extrabold tracking-tight">
          Xin chào Admin, chúc một ngày làm việc hiệu quả!
        </h1>
        <p class="text-slate-300 text-sm sm:text-base leading-relaxed">
          Quản lý toàn bộ hệ sinh thái nhà tuyển dụng, ứng viên, công ty và giám sát các cuộc phỏng vấn tự động hóa AI trong thời gian thực.
        </p>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div 
        v-for="stat in stats" 
        :key="stat.title" 
        class="bg-white dark:bg-slate-800 p-6 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 relative overflow-hidden group hover:shadow-md hover:border-indigo-500/50 transition-all duration-300"
      >
        <div class="flex items-center justify-between mb-4">
          <div :class="`w-12 h-12 rounded-xl flex items-center justify-center border ${stat.colorClass} shadow-sm group-hover:scale-110 transition-transform duration-300`">
            <component :is="stat.icon" size="22" />
          </div>
          <span class="text-[11px] font-bold px-2.5 py-1 rounded-full bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-300">
            {{ stat.badge }}
          </span>
        </div>
        <div class="text-slate-500 dark:text-slate-400 text-sm font-medium">{{ stat.title }}</div>
        <div class="text-3xl font-extrabold text-slate-800 dark:text-white mt-1">
          <span v-if="loading" class="inline-block w-16 h-8 bg-slate-200 dark:bg-slate-700 animate-pulse rounded"></span>
          <span v-else>{{ stat.value }}</span>
        </div>
      </div>
    </div>

    <!-- Main Content Split -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Pending Approvals (Takes 2 Columns) -->
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
          <div class="flex items-center justify-between mb-6">
            <div>
              <h2 class="text-lg font-bold text-slate-800 dark:text-white flex items-center gap-2">
                <UserCheck size="20" class="text-indigo-600 dark:text-indigo-400" />
                Nhà tuyển dụng cần phê duyệt
              </h2>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Xác minh giấy phép kinh doanh trước khi cho phép đăng tuyển</p>
            </div>
            <button 
              @click="router.push('/admin/users')" 
              class="text-xs font-semibold text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 flex items-center gap-1 hover:underline"
            >
              Xem tất cả <ArrowRight size="14" />
            </button>
          </div>

          <!-- Loading State -->
          <div v-if="loading" class="space-y-3">
            <div v-for="n in 3" :key="n" class="h-16 bg-slate-100 dark:bg-slate-700/50 animate-pulse rounded-xl"></div>
          </div>

          <!-- Empty State -->
          <div v-else-if="pendingUsersList.length === 0" class="py-12 text-center bg-slate-50 dark:bg-slate-800/50 rounded-2xl border-2 border-dashed border-slate-200 dark:border-slate-700">
            <div class="w-12 h-12 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400 flex items-center justify-center mx-auto mb-3">
              <CheckCircle2 size="24" />
            </div>
            <h3 class="font-bold text-slate-800 dark:text-white">Không có yêu cầu chờ duyệt</h3>
            <p class="text-sm text-slate-500 dark:text-slate-400 max-w-sm mx-auto mt-1">
              Tất cả các tài khoản nhà tuyển dụng đăng ký mới đều đã được xử lý xong!
            </p>
          </div>

          <!-- Pending List -->
          <div v-else class="space-y-3">
            <div 
              v-for="user in pendingUsersList" 
              :key="user.id"
              class="flex flex-col sm:flex-row sm:items-center justify-between p-4 rounded-xl bg-slate-50 dark:bg-slate-700/40 border border-slate-200/80 dark:border-slate-700 hover:border-indigo-500/40 transition-all gap-4"
            >
              <div class="flex items-center gap-3.5">
                <div class="w-11 h-11 rounded-full bg-indigo-100 dark:bg-indigo-900/40 text-indigo-600 dark:text-indigo-400 font-bold flex items-center justify-center text-base shrink-0 shadow-inner">
                  {{ user.full_name ? user.full_name.charAt(0).toUpperCase() : 'U' }}
                </div>
                <div>
                  <div class="font-bold text-slate-800 dark:text-white text-sm flex items-center gap-2">
                    {{ user.full_name || 'Nhà tuyển dụng mới' }}
                    <span class="px-2 py-0.5 text-[10px] bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400 rounded-md font-semibold">Pending</span>
                  </div>
                  <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{{ user.email }}</div>
                  <div class="text-[11px] text-slate-400 mt-1 flex items-center gap-1">
                    <Clock size="12" /> Đăng ký ngày: {{ formatDate(user.created_at) }}
                  </div>
                </div>
              </div>

              <div class="flex items-center gap-2.5 shrink-0 self-end sm:self-center">
                <button 
                  @click="router.push('/admin/users')"
                  class="px-3 py-2 text-xs font-semibold text-slate-600 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600 rounded-lg transition-colors"
                >
                  Xem chi tiết file
                </button>
                <button 
                  @click="handleApproveUser(user.id)"
                  class="px-4 py-2 text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg shadow-sm shadow-emerald-600/20 transition-colors flex items-center gap-1.5"
                >
                  <CheckCircle2 size="14" /> Duyệt ngay
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Quick Links & System Health -->
      <div class="space-y-6">
        <!-- Quick Actions Panel -->
        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
          <h2 class="text-lg font-bold text-slate-800 dark:text-white mb-4">Thao tác Quản trị</h2>
          <div class="grid grid-cols-1 gap-3">
            <button 
              @click="router.push('/admin/users')" 
              class="flex items-center justify-between p-3.5 bg-slate-50 hover:bg-indigo-50/50 dark:bg-slate-700/50 dark:hover:bg-slate-700 rounded-xl border border-slate-200/80 dark:border-slate-700 group transition-all"
            >
              <div class="flex items-center gap-3">
                <div class="p-2 rounded-lg bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 group-hover:bg-indigo-600 group-hover:text-white transition-colors">
                  <Users size="18" />
                </div>
                <div class="text-left">
                  <div class="text-sm font-bold text-slate-800 dark:text-white">Danh sách tài khoản</div>
                  <div class="text-xs text-slate-500 dark:text-slate-400">Khóa, phân quyền, phê duyệt</div>
                </div>
              </div>
              <ArrowRight size="16" class="text-slate-400 group-hover:text-indigo-600 group-hover:translate-x-0.5 transition-all" />
            </button>

            <button 
              @click="router.push('/admin/companies')" 
              class="flex items-center justify-between p-3.5 bg-slate-50 hover:bg-indigo-50/50 dark:bg-slate-700/50 dark:hover:bg-slate-700 rounded-xl border border-slate-200/80 dark:border-slate-700 group transition-all"
            >
              <div class="flex items-center gap-3">
                <div class="p-2 rounded-lg bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 group-hover:bg-indigo-600 group-hover:text-white transition-colors">
                  <Building size="18" />
                </div>
                <div class="text-left">
                  <div class="text-sm font-bold text-slate-800 dark:text-white">Hồ sơ công ty</div>
                  <div class="text-xs text-slate-500 dark:text-slate-400">Xem doanh nghiệp tham gia</div>
                </div>
              </div>
              <ArrowRight size="16" class="text-slate-400 group-hover:text-indigo-600 group-hover:translate-x-0.5 transition-all" />
            </button>

            <button 
              @click="router.push('/admin/logs')" 
              class="flex items-center justify-between p-3.5 bg-slate-50 hover:bg-indigo-50/50 dark:bg-slate-700/50 dark:hover:bg-slate-700 rounded-xl border border-slate-200/80 dark:border-slate-700 group transition-all"
            >
              <div class="flex items-center gap-3">
                <div class="p-2 rounded-lg bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400 group-hover:bg-indigo-600 group-hover:text-white transition-colors">
                  <Activity size="18" />
                </div>
                <div class="text-left">
                  <div class="text-sm font-bold text-slate-800 dark:text-white">Audit Logs hệ thống</div>
                  <div class="text-xs text-slate-500 dark:text-slate-400">Theo dõi bảo mật & lịch sử</div>
                </div>
              </div>
              <ArrowRight size="16" class="text-slate-400 group-hover:text-indigo-600 group-hover:translate-x-0.5 transition-all" />
            </button>
          </div>
        </div>

        <!-- System Health Widget -->
        <div class="bg-gradient-to-br from-slate-900 to-indigo-950 rounded-2xl p-6 text-white shadow-md border border-slate-800 relative overflow-hidden">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-bold text-base flex items-center gap-2">
              <Activity size="18" class="text-emerald-400 animate-pulse" />
              Trạng thái máy chủ
            </h3>
            <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">ONLINE</span>
          </div>
          <div class="space-y-3 text-xs text-slate-300">
            <div class="flex justify-between items-center py-1 border-b border-slate-800">
              <span>Database Connection (Supabase)</span>
              <span class="font-semibold text-emerald-400">Connected</span>
            </div>
            <div class="flex justify-between items-center py-1 border-b border-slate-800">
              <span>Realtime LiveKit Engine</span>
              <span class="font-semibold text-emerald-400">Active</span>
            </div>
            <div class="flex justify-between items-center py-1">
              <span>AI Evaluation Service</span>
              <span class="font-semibold text-indigo-300">Ready</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
