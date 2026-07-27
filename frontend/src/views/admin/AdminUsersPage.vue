<script setup>
import { ref, onMounted, computed } from 'vue'
import Card from '../../components/common/AppCard.vue'
import { apiService } from '../../services/api.service'
import { fileService } from '../../services/file.service'
import { Users, Search, CheckCircle2, XCircle, Shield, FileText, Clock, RefreshCw, Filter, UserCheck } from 'lucide-vue-next'
import { isOneOf, isValidId, maxLength, normalizeText } from '../../utils/validators.js'

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

const openingFile = ref(null)

const viewFile = async (fileId) => {
  if (!fileId) return
  openingFile.value = fileId
  try {
    const res = await fileService.getDownloadUrl(fileId)
    const url = res?.url
    if (url) {
      window.open(url, '_blank', 'noopener')
    } else {
      alert('Không lấy được đường dẫn file.')
    }
  } catch (error) {
    console.error('Failed to get file signed url', error)
    alert('Lỗi mở file: ' + (error.message || error))
  } finally {
    openingFile.value = null
  }
}

const approveUser = async (userId) => {
  if (!isValidId(userId) || approving.value) return
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
  if (!isValidId(user?.id) || approving.value) return
  if (isOneOf(user.role, ['candidate', 'recruiter', 'admin'], 'Vai trò người dùng không hợp lệ.')) return
  if (isOneOf(user.status, ['active', 'blocked', 'inactive'], 'Trạng thái người dùng không hợp lệ.')) return
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
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm border border-[var(--border)]">
      <div>
        <h1 class="text-h1 flex items-center gap-2.5">
          <Users size="26" class="text-[var(--primary)]" />
          Quản lý Tài khoản & Người dùng
        </h1>
        <p class="text-[var(--text-secondary)] text-sm mt-1">
          Xét duyệt giấy phép nhà tuyển dụng, phân quyền và quản lý trạng thái hoạt động của toàn bộ user.
        </p>
      </div>
      <button 
        @click="fetchUsers" 
        :disabled="loading"
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-[var(--surface-soft)] hover:bg-[var(--border)] text-[var(--text-main)] rounded-xl font-semibold text-sm transition-colors shrink-0 disabled:opacity-50"
      >
        <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới danh sách
      </button>
    </Card>

    <!-- Filters Bar & Tabs -->
    <Card class="p-4 rounded-2xl shadow-sm border border-[var(--border)] flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <!-- Tabs -->
      <div class="flex items-center gap-2 bg-[var(--surface)] p-1.5 rounded-xl border border-[var(--border)]">
        <button 
          @click="activeTab = 'pending'"
          class="flex items-center gap-2 px-4 py-2 rounded-lg font-semibold text-xs sm:text-sm transition-all"
          :class="activeTab === 'pending' ? 'bg-[var(--background)] text-[var(--primary)] shadow-xs border border-[var(--border)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-main)]'"
        >
          <Clock size="16" />
          <span>Chờ duyệt</span>
          <span class="px-2 py-0.5 rounded-full text-[11px]" :class="activeTab === 'pending' ? 'bg-[var(--primary-light)] text-[var(--primary)]' : 'bg-[var(--surface-soft)] text-[var(--text-secondary)]'">
            {{ pendingUsers.length }}
          </span>
        </button>

        <button 
          @click="activeTab = 'all'"
          class="flex items-center gap-2 px-4 py-2 rounded-lg font-semibold text-xs sm:text-sm transition-all"
          :class="activeTab === 'all' ? 'bg-[var(--background)] text-[var(--primary)] shadow-xs border border-[var(--border)]' : 'text-[var(--text-secondary)] hover:text-[var(--text-main)]'"
        >
          <Users size="16" />
          <span>Tất cả thành viên</span>
          <span class="px-2 py-0.5 rounded-full text-[11px]" :class="activeTab === 'all' ? 'bg-[var(--primary-light)] text-[var(--primary)]' : 'bg-[var(--surface-soft)] text-[var(--text-secondary)]'">
            {{ allUsers.length }}
          </span>
        </button>
      </div>

      <!-- Search and Role filter -->
      <div class="flex items-center gap-3">
        <div class="relative flex-1 sm:w-64">
          <Search size="16" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-secondary)]" />
          <input 
            type="text" 
            v-model="searchQuery"
            placeholder="Tìm theo tên, email..." 
            class="w-full pl-9 pr-4 py-2 text-sm bg-[var(--surface)] border border-[var(--border)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
          />
        </div>

        <select 
          v-if="activeTab === 'all'"
          v-model="roleFilter"
          class="px-3.5 py-2 text-sm bg-[var(--surface)] border border-[var(--border)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] font-medium"
        >
          <option value="all">Tất cả vai trò</option>
          <option value="recruiter">Nhà tuyển dụng</option>
          <option value="candidate">Ứng viên</option>
          <option value="admin">Quản trị viên</option>
        </select>
      </div>
    </Card>

    <!-- Data Table Card -->
    <Card class="rounded-2xl shadow-sm border border-[var(--border)] overflow-hidden">
      <!-- Loading Skeleton -->
      <div v-if="loading" class="p-12 text-center space-y-4">
        <div class="inline-block w-8 h-8 border-4 border-[var(--primary)] border-t-transparent rounded-full animate-spin"></div>
        <p class="text-[var(--text-secondary)] text-sm font-medium">Đang tải danh sách người dùng...</p>
      </div>

      <!-- Table View -->
      <div v-else-if="displayedUsers.length > 0" class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-[var(--surface)] border-b border-[var(--border)] text-xs font-bold text-[var(--text-secondary)] uppercase tracking-wider">
              <th class="py-4 px-6">Thành viên</th>
              <th class="py-4 px-6">Vai trò</th>
              <th class="py-4 px-6">Trạng thái</th>
              <th class="py-4 px-6">Ngày tham gia</th>
              <th v-if="activeTab === 'pending'" class="py-4 px-6">Giấy tờ xác thực</th>
              <th class="py-4 px-6 text-right">Thao tác</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[var(--border)] text-sm">
            <tr 
              v-for="user in displayedUsers" 
              :key="user.id" 
              class="hover:bg-[var(--surface)]/50 transition-colors group"
            >
              <!-- Name & Avatar -->
              <td class="py-4 px-6">
                <div class="flex items-center gap-3.5">
                  <div class="w-10 h-10 rounded-full bg-[var(--primary)] text-white font-bold flex items-center justify-center text-sm shadow-xs shrink-0">
                    {{ user.full_name ? user.full_name.charAt(0).toUpperCase() : 'U' }}
                  </div>
                  <div>
                    <div class="font-bold text-[var(--text-main)] group-hover:text-[var(--primary-hover)] transition-colors">
                      {{ user.full_name || 'Chưa đặt tên' }}
                    </div>
                    <div class="text-xs text-[var(--text-secondary)] mt-0.5">{{ user.email }}</div>
                  </div>
                </div>
              </td>

              <!-- Role -->
              <td class="py-4 px-6">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-bold uppercase tracking-wide"
                  :class="{
                    'bg-[var(--primary-light)] text-[var(--primary)] border border-blue-200/40': user.role === 'recruiter' || user.role === 'candidate',
                    'bg-[var(--highlight-bg)] text-[var(--highlight-hover)] border border-[var(--highlight)]/20': user.role === 'admin'
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
                    'bg-[var(--success)]/10 text-[var(--success)] border border-[var(--success)]/20': user.status === 'active',
                    'bg-[var(--warning)]/10 text-[var(--warning)] border border-[var(--warning)]/20': user.status === 'pending',
                    'bg-[var(--danger)]/10 text-[var(--danger)] border border-[var(--danger)]/20': user.status === 'blocked' || user.status === 'inactive'
                  }">
                  <span class="w-1.5 h-1.5 rounded-full animate-pulse" :class="{
                    'bg-[var(--success)]': user.status === 'active',
                    'bg-[var(--warning)]': user.status === 'pending',
                    'bg-[var(--danger)]': user.status === 'blocked' || user.status === 'inactive'
                  }"></span>
                  {{ user.status === 'active' ? 'Hoạt động' : (user.status === 'pending' ? 'Chờ duyệt' : 'Đã khóa') }}
                </span>
              </td>

              <!-- Date -->
              <td class="py-4 px-6 text-[var(--text-secondary)] text-xs">
                {{ formatDate(user.created_at) }}
              </td>

              <!-- Verification File -->
              <td v-if="activeTab === 'pending'" class="py-4 px-6">
                <button
                  v-if="user.verification_file_id"
                  :disabled="openingFile === user.verification_file_id"
                  @click="viewFile(user.verification_file_id)"
                  class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-[var(--primary-light)] text-[var(--primary)] hover:bg-[var(--primary-light)]/80 rounded-lg text-xs font-semibold transition-colors border border-blue-200/40 disabled:opacity-50"
                >
                  <FileText size="14" /> {{ openingFile === user.verification_file_id ? 'Đang mở...' : 'Xem giấy phép' }}
                </button>
                <span v-else class="inline-flex items-center gap-1 text-[var(--text-secondary)] italic text-xs">
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
                    class="px-3.5 py-2 bg-[var(--success)] hover:bg-[var(--success)]/90 text-white font-bold text-xs rounded-xl shadow-xs transition-all disabled:opacity-50 flex items-center gap-1.5"
                  >
                    <CheckCircle2 size="14" />
                    {{ approving === user.id ? 'Đang duyệt...' : 'Phê duyệt' }}
                  </button>

                  <button 
                    v-if="user.status === 'active' && user.role !== 'admin'" 
                    @click="toggleUserStatus(user)"
                    class="px-3 py-1.5 border border-[var(--border)] hover:bg-[var(--danger)]/15 hover:text-[var(--danger)] hover:border-[var(--danger)]/20 rounded-xl text-[var(--text-secondary)] font-medium text-xs transition-colors"
                  >
                    Khóa tài khoản
                  </button>

                  <button 
                    v-if="(user.status === 'blocked' || user.status === 'inactive') && user.role !== 'admin'" 
                    @click="toggleUserStatus(user)"
                    class="px-3 py-1.5 border border-[var(--success)]/20 hover:bg-[var(--success)]/15 hover:text-[var(--success)] hover:border-[var(--success)]/25 rounded-xl text-[var(--text-secondary)] font-medium text-xs transition-colors"
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
        <div class="w-16 h-16 bg-[var(--surface)] rounded-full flex items-center justify-center mx-auto mb-4 text-[var(--text-secondary)]">
          <UserCheck size="32" />
        </div>
        <h3 class="font-bold text-[var(--text-main)] text-base">Không tìm thấy người dùng nào</h3>
        <p class="text-[var(--text-secondary)] text-sm mt-1">
          {{ activeTab === 'pending' ? 'Tất cả các nhà tuyển dụng đăng ký mới đã được xét duyệt xong!' : 'Chưa có tài khoản nào khớp với từ khóa tìm kiếm của bạn.' }}
        </p>
      </div>
    </Card>
  </div>
</template>
