<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import Card from '../../components/common/AppCard.vue'
import { apiService } from '../../services/api.service'
import { Activity, Search, RefreshCw, Clock, Shield, Globe, Terminal, User, AlertCircle, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-vue-next'

const logs = ref([])
const loading = ref(false)
const searchQuery = ref('')
const actionFilter = ref('all')

// CẤU HÌNH PHÂN TRANG (PAGINATION STATE)
const currentPage = ref(1)
const pageSize = ref(15) // Mặc định 15 dòng mỗi trang

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

// Tự động quay về trang 1 khi tìm kiếm, đổi bộ lọc hoặc đổi số dòng mỗi trang
watch([searchQuery, actionFilter, pageSize], () => {
  currentPage.value = 1
})

// Tổng số trang
const totalPages = computed(() => {
  return Math.max(1, Math.ceil(displayedLogs.value.length / pageSize.value))
})

// Danh sách log hiển thị ở trang hiện tại
const paginatedLogs = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return displayedLogs.value.slice(start, end)
})

// Danh sách các số trang hiển thị trên thanh phân trang (Smart page window)
const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const cur = currentPage.value

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i)
  } else {
    if (cur <= 4) {
      pages.push(1, 2, 3, 4, 5, '...', total)
    } else if (cur >= total - 3) {
      pages.push(1, '...', total - 4, total - 3, total - 2, total - 1, total)
    } else {
      pages.push(1, '...', cur - 1, cur, cur + 1, '...', total)
    }
  }
  return pages
})

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await apiService.get('/admin/logs')
    logs.value = Array.isArray(res) ? res : (res.data || [])
    currentPage.value = 1
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
  return 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300 border border-blue-200 dark:border-blue-800'
}

onMounted(() => {
  fetchLogs()
})
</script>

<template>
  <div class="space-y-6 animate-fade-in">
    <!-- Header -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm border border-[var(--border)]">
      <div>
        <h1 class="text-h1 flex items-center gap-2.5">
          <Activity size="26" class="text-[var(--primary)]" />
          Audit Logs - Nhật ký Hoạt động Hệ thống
        </h1>
        <p class="text-[var(--text-secondary)] text-sm mt-1">
          Ghi nhận tự động toàn bộ thao tác bảo mật, xác thực, chỉnh sửa và truy xuất dữ liệu từ tất cả người dùng.
        </p>
      </div>
      <button 
        @click="fetchLogs" 
        :disabled="loading"
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-[var(--surface-soft)] hover:bg-[var(--border)] text-[var(--text-main)] rounded-xl font-semibold text-sm transition-colors shrink-0 disabled:opacity-50"
      >
        <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới nhật ký
      </button>
    </Card>

    <!-- Filters & Pagination Toolbar -->
    <Card class="p-4 rounded-2xl shadow-sm border border-[var(--border)] flex flex-col lg:flex-row lg:items-center justify-between gap-4">
      <div class="relative flex-1 max-w-md">
        <Search size="16" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-secondary)]" />
        <input 
          type="text" 
          v-model="searchQuery"
          placeholder="Tìm theo thao tác, ID tài nguyên, IP..." 
          class="w-full pl-9 pr-4 py-2 text-sm bg-[var(--surface)] border border-[var(--border)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
        />
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <!-- Bộ lọc hành động -->
        <select 
          v-model="actionFilter"
          class="px-3.5 py-2 text-sm bg-[var(--surface)] border border-[var(--border)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] font-medium"
        >
          <option value="all">Tất cả hành động</option>
          <option value="LOGIN">Đăng nhập (Auth)</option>
          <option value="CREATE">Tạo mới</option>
          <option value="UPDATE">Cập nhật / Sửa</option>
          <option value="APPROVE">Phê duyệt</option>
          <option value="DELETE">Xóa</option>
        </select>

        <!-- Chọn số dòng mỗi trang -->
        <div class="flex items-center gap-2 bg-[var(--surface)] border border-[var(--border)] px-3 py-1.5 rounded-xl text-xs font-semibold text-[var(--text-secondary)]">
          <span>Số dòng:</span>
          <select 
            v-model.number="pageSize"
            class="bg-transparent text-[var(--primary)] font-bold focus:outline-none cursor-pointer text-sm"
          >
            <option :value="10">10</option>
            <option :value="15">15</option>
            <option :value="25">25</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </div>

        <!-- Thông tin số lượng -->
        <div class="text-xs font-semibold text-[var(--text-secondary)] pl-2 border-l border-[var(--border)]">
          Tổng: <span class="text-[var(--primary)] font-bold text-sm">{{ displayedLogs.length }}</span> nhật ký
        </div>
      </div>
    </Card>

    <!-- Logs Table Card -->
    <Card class="rounded-2xl shadow-sm border border-[var(--border)] overflow-hidden flex flex-col justify-between min-h-[500px]">
      <div>
        <!-- Loading Skeleton -->
        <div v-if="loading" class="p-16 text-center space-y-4">
          <div class="inline-block w-8 h-8 border-4 border-[var(--primary)] border-t-transparent rounded-full animate-spin"></div>
          <p class="text-[var(--text-secondary)] text-sm font-medium">Đang truy xuất dữ liệu nhật ký hệ thống...</p>
        </div>

        <!-- Logs Table -->
        <div v-else-if="paginatedLogs.length > 0" class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-sm">
            <thead>
              <tr class="bg-[var(--surface)] border-b border-[var(--border)] text-xs font-bold text-[var(--text-secondary)] uppercase tracking-wider">
                <th class="py-4 px-6 w-48">Thời gian</th>
                <th class="py-4 px-6">Người thực hiện (Actor)</th>
                <th class="py-4 px-6">Hành động</th>
                <th class="py-4 px-6">Tài nguyên & ID</th>
                <th class="py-4 px-6 text-right">IP / Client Device</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[var(--border)]">
              <tr 
                v-for="log in paginatedLogs" 
                :key="log.id" 
                class="hover:bg-[var(--surface)]/50 transition-colors group"
              >
                <!-- Time -->
                <td class="py-4 px-6 font-mono text-xs text-[var(--text-secondary)] whitespace-nowrap">
                  <div class="flex items-center gap-1.5 font-medium">
                    <Clock size="13" class="text-[var(--text-secondary)] shrink-0" />
                    <span>{{ formatDate(log.created_at) }}</span>
                  </div>
                </td>

                <!-- Actor -->
                <td class="py-4 px-6">
                  <div class="flex items-center gap-2.5">
                    <div class="w-8 h-8 rounded-lg bg-[var(--primary-light)] text-[var(--primary)] font-bold flex items-center justify-center text-xs shrink-0">
                      <User size="14" />
                    </div>
                    <div>
                      <div class="font-bold text-[var(--text-main)] text-xs truncate max-w-[180px]" :title="log.actor_user_id">
                        {{ log.actor_user_id || 'Hệ thống tự động' }}
                      </div>
                      <span class="inline-block px-2 py-0.5 mt-0.5 rounded text-[10px] font-bold uppercase bg-[var(--surface-soft)] text-[var(--text-secondary)]">
                        {{ log.actor_role || 'system' }}
                      </span>
                    </div>
                  </div>
                </td>

                <!-- Action Badge -->
                <td class="py-4 px-6">
                  <span :class="`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-extrabold tracking-wide uppercase ${getActionColor(log.action)}`">
                    <Terminal size="13" class="shrink-0" />
                    <span>{{ log.action || 'UNKNOWN' }}</span>
                  </span>
                </td>

                <!-- Resource -->
                <td class="py-4 px-6">
                  <div class="font-bold text-[var(--text-main)] text-xs flex items-center gap-1.5">
                    <Shield size="13" class="text-[var(--text-secondary)] shrink-0" />
                    <span>{{ log.resource_type || 'system' }}</span>
                  </div>
                  <div class="font-mono text-[11px] text-[var(--text-secondary)] truncate max-w-[180px] mt-0.5" :title="log.resource_id">
                    ID: {{ log.resource_id || 'N/A' }}
                  </div>
                </td>

                <!-- IP Address / Device -->
                <td class="py-4 px-6 text-right font-mono text-xs text-[var(--text-secondary)]">
                  <div class="flex items-center justify-end gap-1.5 text-[var(--text-main)] font-semibold">
                    <Globe size="13" class="text-[var(--primary)] shrink-0" />
                    <span>{{ log.ip_address || 'localhost' }}</span>
                  </div>
                  <div class="text-[11px] text-[var(--text-secondary)] truncate max-w-[200px] ml-auto mt-0.5" :title="log.user_agent">
                    {{ log.user_agent || 'Client Engine' }}
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Empty State -->
        <div v-else class="p-16 text-center">
          <div class="w-16 h-16 bg-[var(--surface)] rounded-full flex items-center justify-center mx-auto mb-4 text-[var(--text-secondary)]">
            <Activity size="32" />
          </div>
          <h3 class="font-bold text-[var(--text-main)] text-base">Chưa ghi nhận nhật ký nào</h3>
          <p class="text-[var(--text-secondary)] text-sm mt-1">
            {{ searchQuery || actionFilter !== 'all' ? 'Không có dòng log nào khớp với bộ lọc tìm kiếm.' : 'Hệ thống hiện chưa có thao tác nào cần lưu lại trong Audit Logs.' }}
          </p>
        </div>
      </div>

      <!-- THANH PHÂN TRANG (PAGINATION CONTROLS BAR) -->
      <div v-if="displayedLogs.length > 0" class="border-t border-[var(--border)] px-6 py-4 bg-[var(--surface)] flex flex-col sm:flex-row items-center justify-between gap-4">
        <!-- Thông tin trang & Phạm vi hiển thị có thể chọn trực tiếp (Interactive Inline Selector) -->
        <div class="flex flex-wrap items-center gap-1.5 text-xs font-medium text-[var(--text-secondary)]">
          <span>Hiển thị từ</span>
          <span class="font-bold text-[var(--text-main)]">{{ (currentPage - 1) * pageSize + 1 }}</span>
          <span>đến</span>
          
          <!-- Hộp chọn số dòng/trang trực tiếp ngay tại chân bảng -->
          <div class="inline-flex items-center bg-[var(--background)] border border-[var(--border)] px-2 py-1 rounded-lg shadow-2xs hover:border-[var(--accent)] transition-colors">
            <select 
              v-model.number="pageSize"
              class="bg-transparent font-bold text-[var(--primary)] focus:outline-none cursor-pointer text-xs pr-1"
              title="Bấm để đổi số dòng hiển thị mỗi trang"
            >
              <option :value="10">10 dòng/trang</option>
              <option :value="15">15 dòng/trang</option>
              <option :value="25">25 dòng/trang</option>
              <option :value="50">50 dòng/trang</option>
              <option :value="100">100 dòng/trang</option>
            </select>
          </div>

          <span>trong tổng số</span>
          <span class="font-bold text-[var(--primary)]">{{ displayedLogs.length }}</span>
          <span>nhật ký</span>
        </div>

        <!-- Các nút chuyển trang -->
        <div class="flex items-center gap-1.5">
          <!-- Trang đầu -->
          <button 
            @click="goToPage(1)" 
            :disabled="currentPage === 1"
            class="p-2 rounded-lg border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--surface-soft)] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            title="Trang đầu"
          >
            <ChevronsLeft size="15" />
          </button>

          <!-- Trang trước -->
          <button 
            @click="goToPage(currentPage - 1)" 
            :disabled="currentPage === 1"
            class="p-2 rounded-lg border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--surface-soft)] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            title="Trang trước"
          >
            <ChevronLeft size="15" />
          </button>

          <!-- Danh sách số trang (Smart Page Pills) -->
          <template v-for="(page, idx) in visiblePages" :key="idx">
            <span v-if="page === '...'" class="px-2 py-1 text-[var(--text-muted)] font-bold text-xs">...</span>
            <button 
              v-else 
              @click="goToPage(page)"
              :class="[
                'min-w-[34px] h-[34px] px-2.5 rounded-lg text-xs font-bold transition-all flex items-center justify-center',
                currentPage === page 
                  ? 'bg-[var(--primary)] text-white shadow-xs border border-[var(--primary)]' 
                  : 'bg-[var(--background)] border border-[var(--border)] text-[var(--text-main)] hover:bg-[var(--surface-soft)]'
              ]"
            >
              {{ page }}
            </button>
          </template>

          <!-- Trang tiếp -->
          <button 
            @click="goToPage(currentPage + 1)" 
            :disabled="currentPage === totalPages"
            class="p-2 rounded-lg border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--surface-soft)] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            title="Trang tiếp"
          >
            <ChevronRight size="15" />
          </button>

          <!-- Trang cuối -->
          <button 
            @click="goToPage(totalPages)" 
            :disabled="currentPage === totalPages"
            class="p-2 rounded-lg border border-[var(--border)] text-[var(--text-secondary)] hover:bg-[var(--surface-soft)] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            title="Trang cuối"
          >
            <ChevronsRight size="15" />
          </button>
        </div>
      </div>
    </Card>
  </div>
</template>
