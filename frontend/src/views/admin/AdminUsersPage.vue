<script setup>
import { ref, onMounted, computed } from 'vue'
import Card from '../../components/common/AppCard.vue'
import { apiService } from '../../services/api.service'
import { Users, Search, CheckCircle2, XCircle, Shield, FileText, Clock, RefreshCw, Filter, UserCheck } from 'lucide-vue-next'

const pendingUsers = ref([])
const allUsers = ref([])
const loading = ref(false)
const approving = ref(null)
const activeTab = ref('pending') // 'pending' | 'all'
const searchQuery = ref('')
const roleFilter = ref('all')

const displayedUsers = computed(() => {
  const source = activeTab.value === 'pending' ? pendingUsers.value : allUsers.value
  if (!source || !Array.isArray(source)) return []
  
  return source.filter(u => {
    const matchSearch = !searchQuery.value || 
      (u.full_name && u.full_name.toLowerCase().includes(searchQuery.value.toLowerCase())) ||
      (u.email && u.email.toLowerCase().includes(searchQuery.value.toLowerCase()))
    
    const matchRole = roleFilter.value === 'all' || u.role === roleFilter.value
    return matchSearch && matchRole
  })
})

const fetchUsers = async () => {
  loading.value = true
  try {
    const [pendingRes, allRes] = await Promise.all([
      apiService.get('/admin/users/pending'),
      apiService.get('/admin/users')
    ])
    // apiService.get unwraps { success: true, data: [...] } and directly returns [...]
    pendingUsers.value = Array.isArray(pendingRes) ? pendingRes : (pendingRes.data || [])
    allUsers.value = Array.isArray(allRes) ? allRes : (allRes.data || [])
  } catch (error) {
    console.error('Failed to fetch users', error)
  } finally {
    loading.value = false
  }
}

const getFileUrl = (fileId) => {
  const apiUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  return `${apiUrl}/files/${fileId}/download`
}

const approveUser = async (userId) => {
  if (!confirm('Bạn có chắc chắn muốn phê duyệt nhà tuyển dụng này? Họ sẽ có thể đăng tin tuyển dụng ngay lập tức.')) return
  
  approving.value = userId
  try {
    await apiService.put(`/admin/users/${userId}/approve`)
    await fetchUsers()
    alert('Đã phê duyệt tài khoản thành công!')
  } catch (error) {
    console.error('Failed to approve user', error)
    alert('Lỗi phê duyệt: ' + (error.response?.data?.message || error.message || error))
  } finally {
    approving.value = null
  }
}

const toggleUserStatus = async (user) => {
  const newStatus = user.status === 'active' ? 'blocked' : 'active'
  const actionText = newStatus === 'blocked' ? 'khóa' : 'mở khóa'
  if (!confirm(`Bạn có chắc chắn muốn ${actionText} tài khoản của ${user.full_name || user.email}?`)) return
  
  try {
    await apiService.put(`/admin/users/${user.id}/status`, { status: newStatus })
    await fetchUsers()
    alert(`Đã ${actionText} tài khoản thành công!`)
  } catch (error) {
    console.error('Failed to update user status', error)
    alert(`Lỗi khi ${actionText} tài khoản: ` + (error.response?.data?.message || error.message || error))
  }
}

const formatDate = (isoStr) => {
  if (!isoStr) return '—'
  return new Date(isoStr).toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="space-y-6 animate-fade-in">
    <!-- Page Header -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-white flex items-center gap-2.5">
          <Users size="26" class="text-blue-600 dark:text-blue-400" />
          Quản lý Tài khoản & Người dùng
        </h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
          Xét duyệt giấy phép nhà tuyển dụng, phân quyền và quản lý trạng thái hoạt động của toàn bộ user.
        </p>
      </div>
      <button 
        @click="fetchUsers" 
        :disabled="loading"
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-200 rounded-xl font-medium text-sm transition-colors shrink-0 disabled:opacity-50"
      >
        <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới danh sách
      </button>
    </Card>

    <!-- Filters Bar & Tabs -->
    <Card class="p-4 rounded-2xl shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <!-- Tabs -->
      <div class="flex items-center gap-2 bg-slate-100 dark:bg-slate-900/60 p-1.5 rounded-xl border border-slate-200/60 dark:border-slate-700/60">
        <button 
          @click="activeTab = 'pending'"
          class="flex items-center gap-2 px-4 py-2 rounded-lg font-semibold text-xs sm:text-sm transition-all"
          :class="activeTab === 'pending' ? 'bg-white dark:bg-slate-800 text-blue-600 dark:text-blue-400 shadow-sm' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'"
        >
          <Clock size="16" />
          <span>Chờ duyệt</span>
          <span class="px-2 py-0.5 rounded-full text-[11px]" :class="activeTab === 'pending' ? 'bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'">
            {{ pendingUsers.length }}
          </span>
        </button>

        <button 
          @click="activeTab = 'all'"
          class="flex items-center gap-2 px-4 py-2 rounded-lg font-semibold text-xs sm:text-sm transition-all"
          :class="activeTab === 'all' ? 'bg-white dark:bg-slate-800 text-blue-600 dark:text-blue-400 shadow-sm' : 'text-slate-500 hover:text-slate-800 dark:hover:text-slate-200'"
        >
          <Users size="16" />
          <span>Tất cả thành viên</span>
          <span class="px-2 py-0.5 rounded-full text-[11px]" :class="activeTab === 'all' ? 'bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'">
            {{ allUsers.length }}
          </span>
        </button>
      </div>

      <!-- Search and Role filter -->
      <div class="flex items-center gap-3">
        <div class="relative flex-1 sm:w-64">
          <Search size="16" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
          <input 
            type="text" 
            v-model="searchQuery"
            placeholder="Tìm theo tên, email..." 
            class="w-full pl-9 pr-4 py-2 text-sm bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
          />
        </div>

        <select 
          v-if="activeTab === 'all'"
          v-model="roleFilter"
          class="px-3.5 py-2 text-sm bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-700 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-500 font-medium"
        >
          <option value="all">Tất cả vai trò</option>
          <option value="recruiter">Nhà tuyển dụng</option>
          <option value="candidate">Ứng viên</option>
          <option value="admin">Quản trị viên</option>
        </select>
      </div>
    </Card>

    <!-- Data Table Card -->
    <Card class="rounded-2xl shadow-sm overflow-hidden">
      <!-- Loading Skeleton -->
      <div v-if="loading" class="p-12 text-center space-y-4">
        <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-slate-500 dark:text-slate-400 text-sm font-medium">Đang tải danh sách người dùng...</p>
      </div>

      <!-- Table View -->
      <div v-else-if="displayedUsers.length > 0" class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50/80 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-700 text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
              <th class="py-4 px-6">Thành viên</th>
              <th class="py-4 px-6">Vai trò</th>
              <th class="py-4 px-6">Trạng thái</th>
              <th class="py-4 px-6">Ngày tham gia</th>
              <th v-if="activeTab === 'pending'" class="py-4 px-6">Giấy tờ xác thực</th>
              <th class="py-4 px-6 text-right">Thao tác</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-700/50 text-sm">
            <tr 
              v-for="user in displayedUsers" 
              :key="user.id" 
              class="hover:bg-slate-50/80 dark:hover:bg-slate-700/30 transition-colors group"
            >
              <!-- Name & Avatar -->
              <td class="py-4 px-6">
                <div class="flex items-center gap-3.5">
                  <div class="w-10 h-10 rounded-full bg-gradient-to-tr from-blue-500 to-blue-600 text-white font-bold flex items-center justify-center text-sm shadow-sm shrink-0">
                    {{ user.full_name ? user.full_name.charAt(0).toUpperCase() : 'U' }}
                  </div>
                  <div>
                    <div class="font-bold text-slate-800 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                      {{ user.full_name || 'Chưa đặt tên' }}
                    </div>
                    <div class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">{{ user.email }}</div>
                  </div>
                </div>
              </td>

              <!-- Role -->
              <td class="py-4 px-6">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-bold uppercase tracking-wide"
                  :class="{
                    'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800/40': user.role === 'recruiter',
                    'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800/40': user.role === 'candidate',
                    'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300 border border-amber-200 dark:border-amber-800/40': user.role === 'admin'
                  }">
                  <Shield v-if="user.role === 'admin'" size="13" />
                  <Users v-else size="13" />
                  {{ user.role === 'recruiter' ? 'Nhà tuyển dụng' : (user.role === 'candidate' ? 'Ứng viên' : 'Quản trị viên') }}
                </span>
              </td>

              <!-- Status -->
              <td class="py-4 px-6">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold"
                  :class="{
                    'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800/40': user.status === 'active',
                    'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400 border border-amber-200 dark:border-amber-800/40': user.status === 'pending',
                    'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 border border-red-200 dark:border-red-800/40': user.status === 'blocked' || user.status === 'inactive'
                  }">
                  <span class="w-1.5 h-1.5 rounded-full" :class="{
                    'bg-emerald-500': user.status === 'active',
                    'bg-amber-500': user.status === 'pending',
                    'bg-red-500': user.status === 'blocked' || user.status === 'inactive'
                  }"></span>
                  {{ user.status === 'active' ? 'Hoạt động' : (user.status === 'pending' ? 'Chờ duyệt' : 'Đã khóa') }}
                </span>
              </td>

              <!-- Date -->
              <td class="py-4 px-6 text-slate-500 dark:text-slate-400 text-xs">
                {{ formatDate(user.created_at) }}
              </td>

              <!-- Verification File -->
              <td v-if="activeTab === 'pending'" class="py-4 px-6">
                <a 
                  v-if="user.verification_file_id" 
                  :href="getFileUrl(user.verification_file_id)" 
                  target="_blank" 
                  class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 hover:bg-blue-100 rounded-lg text-xs font-semibold transition-colors border border-blue-200/60 dark:border-blue-800/40"
                >
                  <FileText size="14" /> Xem giấy phép
                </a>
                <span v-else class="inline-flex items-center gap-1 text-slate-400 italic text-xs">
                  <Clock size="14" /> Chưa tải file
                </span>
              </td>

              <!-- Actions -->
              <td class="py-4 px-6 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button 
                    v-if="user.status === 'pending'"
                    @click="approveUser(user.id)" 
                    :disabled="approving === user.id"
                    class="px-3.5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-xl shadow-sm shadow-emerald-600/20 transition-all disabled:opacity-50 flex items-center gap-1.5"
                  >
                    <CheckCircle2 size="14" />
                    {{ approving === user.id ? 'Đang duyệt...' : 'Phê duyệt' }}
                  </button>

                  <button 
                    v-if="user.status === 'active' && user.role !== 'admin'" 
                    @click="toggleUserStatus(user)"
                    class="px-3 py-1.5 border border-slate-200 dark:border-slate-700 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600 dark:hover:text-red-400 rounded-xl text-slate-600 dark:text-slate-400 font-medium text-xs transition-colors"
                  >
                    Khóa tài khoản
                  </button>

                  <button 
                    v-if="(user.status === 'blocked' || user.status === 'inactive') && user.role !== 'admin'" 
                    @click="toggleUserStatus(user)"
                    class="px-3 py-1.5 border border-emerald-200 dark:border-emerald-700 hover:bg-emerald-50 dark:hover:bg-emerald-900/20 hover:text-emerald-600 dark:hover:text-emerald-400 rounded-xl text-emerald-600 dark:text-emerald-400 font-medium text-xs transition-colors"
                  >
                    Mở khóa
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty State -->
      <div v-else class="p-16 text-center">
        <div class="w-16 h-16 bg-slate-100 dark:bg-slate-700/50 rounded-full flex items-center justify-center mx-auto mb-4 text-slate-400">
          <UserCheck size="32" />
        </div>
        <h3 class="font-bold text-slate-800 dark:text-white text-base">Không tìm thấy người dùng nào</h3>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
          {{ activeTab === 'pending' ? 'Tất cả các nhà tuyển dụng đăng ký mới đã được xét duyệt xong!' : 'Chưa có tài khoản nào khớp với từ khóa tìm kiếm của bạn.' }}
        </p>
      </div>
    </Card>
  </div>
</template>
