<script setup>
import { ref, onMounted, computed } from 'vue'
import { apiService } from '../../services/api.service'
import { Activity, Search, RefreshCw, Clock, Shield, Globe, Terminal, User, AlertCircle } from 'lucide-vue-next'

const logs = ref([])
const loading = ref(false)
const searchQuery = ref('')
const actionFilter = ref('all')

const displayedLogs = computed(() => {
  if (!logs.value || !Array.isArray(logs.value)) return []
  return logs.value.filter(l => {
    const q = searchQuery.value.toLowerCase()
    const matchSearch = !searchQuery.value ||
      (l.action && l.action.toLowerCase().includes(q)) ||
      (l.resource_type && l.resource_type.toLowerCase().includes(q)) ||
      (l.actor_user_id && l.actor_user_id.toLowerCase().includes(q)) ||
      (l.ip_address && l.ip_address.toLowerCase().includes(q))
    
    const matchAction = actionFilter.value === 'all' || 
      (l.action && l.action.toUpperCase().includes(actionFilter.value))
    
    return matchSearch && matchAction
  })
})

const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await apiService.get('/admin/logs')
    logs.value = Array.isArray(res) ? res : (res.data || [])
  } catch (error) {
    console.error('Failed to fetch logs', error)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  return d.toLocaleString('vi-VN', {
    hour: '2-digit', minute: '2-digit', second: '2-digit',
    day: '2-digit', month: '2-digit', year: 'numeric'
  })
}

const getActionColor = (action) => {
  if (!action) return 'bg-slate-100 text-slate-700 dark:bg-slate-700 dark:text-slate-300'
  const a = action.toUpperCase()
  if (a.includes('LOGIN') || a.includes('AUTH')) return 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300 border border-blue-200 dark:border-blue-800'
  if (a.includes('APPROVE') || a.includes('VERIFY') || a.includes('CREATE')) return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800'
  if (a.includes('DELETE') || a.includes('BLOCK') || a.includes('REJECT')) return 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300 border border-red-200 dark:border-red-800'
  if (a.includes('UPDATE') || a.includes('EDIT')) return 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300 border border-amber-200 dark:border-amber-800'
  return 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-300 border border-purple-200 dark:border-purple-800'
}

onMounted(() => {
  fetchLogs()
})
</script>

<template>
  <div class="space-y-6 animate-fade-in">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-white dark:bg-slate-800 p-6 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-white flex items-center gap-2.5">
          <Activity size="26" class="text-indigo-600 dark:text-indigo-400" />
          Audit Logs - Nhật ký Hoạt động Hệ thống
        </h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
          Ghi nhận tự động toàn bộ thao tác bảo mật, xác thực, chỉnh sửa và truy xuất dữ liệu từ tất cả người dùng.
        </p>
      </div>
      <button 
        @click="fetchLogs" 
        :disabled="loading"
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-200 rounded-xl font-medium text-sm transition-colors shrink-0 disabled:opacity-50"
      >
        <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới nhật ký
      </button>
    </div>

    <!-- Filters Bar -->
    <div class="bg-white dark:bg-slate-800 p-4 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div class="relative flex-1 max-w-md">
        <Search size="16" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
        <input 
          type="text" 
          v-model="searchQuery"
          placeholder="Tìm theo thao tác, ID tài nguyên, IP..." 
          class="w-full pl-9 pr-4 py-2 text-sm bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all"
        />
      </div>

      <div class="flex items-center gap-3">
        <select 
          v-model="actionFilter"
          class="px-3.5 py-2 text-sm bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-700 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 font-medium"
        >
          <option value="all">Tất cả hành động</option>
          <option value="LOGIN">Đăng nhập (Auth)</option>
          <option value="CREATE">Tạo mới</option>
          <option value="UPDATE">Cập nhật / Sửa</option>
          <option value="APPROVE">Phê duyệt</option>
          <option value="DELETE">Xóa</option>
        </select>
        <div class="text-xs font-semibold text-slate-500 dark:text-slate-400">
          Hiển thị: <span class="text-indigo-600 dark:text-indigo-400 font-bold text-sm">{{ displayedLogs.length }}</span> dòng
        </div>
      </div>
    </div>

    <!-- Logs Table Card -->
    <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 overflow-hidden">
      <!-- Loading Skeleton -->
      <div v-if="loading" class="p-12 text-center space-y-4">
        <div class="inline-block w-8 h-8 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-slate-500 dark:text-slate-400 text-sm font-medium">Đang truy xuất dữ liệu nhật ký hệ thống...</p>
      </div>

      <!-- Logs Table -->
      <div v-else-if="displayedLogs.length > 0" class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-sm">
          <thead>
            <tr class="bg-slate-50/80 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
              <th class="py-4 px-6 w-48">Thời gian</th>
              <th class="py-4 px-6">Người thực hiện (Actor)</th>
              <th class="py-4 px-6">Hành động</th>
              <th class="py-4 px-6">Tài nguyên & ID</th>
              <th class="py-4 px-6 text-right">IP / Client Device</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-700/50">
            <tr 
              v-for="log in displayedLogs" 
              :key="log.id" 
              class="hover:bg-slate-50/80 dark:hover:bg-slate-700/30 transition-colors group"
            >
              <!-- Time -->
              <td class="py-4 px-6 font-mono text-xs text-slate-500 dark:text-slate-400 whitespace-nowrap">
                <div class="flex items-center gap-1.5 font-medium">
                  <Clock size="13" class="text-slate-400" />
                  {{ formatDate(log.created_at) }}
                </div>
              </td>

              <!-- Actor -->
              <td class="py-4 px-6">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-lg bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 font-bold flex items-center justify-center text-xs shrink-0">
                    <User size="14" />
                  </div>
                  <div>
                    <div class="font-bold text-slate-800 dark:text-white text-xs truncate max-w-[180px]" :title="log.actor_user_id">
                      {{ log.actor_user_id || 'Hệ thống tự động' }}
                    </div>
                    <span class="inline-block px-2 py-0.5 mt-0.5 rounded text-[10px] font-bold uppercase bg-slate-100 dark:bg-slate-700 text-slate-600 dark:text-slate-300">
                      {{ log.actor_role || 'system' }}
                    </span>
                  </div>
                </div>
              </td>

              <!-- Action Badge -->
              <td class="py-4 px-6">
                <span :class="`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-extrabold tracking-wide uppercase ${getActionColor(log.action)}`">
                  <Terminal size="13" />
                  {{ log.action || 'UNKNOWN' }}
                </span>
              </td>

              <!-- Resource -->
              <td class="py-4 px-6">
                <div class="font-bold text-slate-800 dark:text-slate-200 text-xs flex items-center gap-1.5">
                  <Shield size="13" class="text-slate-400" />
                  {{ log.resource_type || 'system' }}
                </div>
                <div class="font-mono text-[11px] text-slate-400 truncate max-w-[180px] mt-0.5" :title="log.resource_id">
                  ID: {{ log.resource_id || 'N/A' }}
                </div>
              </td>

              <!-- IP Address / Device -->
              <td class="py-4 px-6 text-right font-mono text-xs text-slate-500 dark:text-slate-400">
                <div class="flex items-center justify-end gap-1.5 text-slate-700 dark:text-slate-300 font-semibold">
                  <Globe size="13" class="text-indigo-500" />
                  {{ log.ip_address || 'localhost' }}
                </div>
                <div class="text-[11px] text-slate-400 truncate max-w-[200px] ml-auto mt-0.5" :title="log.user_agent">
                  {{ log.user_agent || 'WeMake Client Engine' }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty State -->
      <div v-else class="p-16 text-center">
        <div class="w-16 h-16 bg-slate-100 dark:bg-slate-700/50 rounded-full flex items-center justify-center mx-auto mb-4 text-slate-400">
          <Activity size="32" />
        </div>
        <h3 class="font-bold text-slate-800 dark:text-white text-base">Chưa ghi nhận nhật ký nào</h3>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
          {{ searchQuery || actionFilter !== 'all' ? 'Không có dòng log nào khớp với bộ lọc tìm kiếm.' : 'Hệ thống hiện chưa có thao tác nào cần lưu lại trong Audit Logs.' }}
        </p>
      </div>
    </div>
  </div>
</template>
