<script setup>
import { ref, onMounted } from 'vue'
import { Search, UserX, UserCheck, Shield } from 'lucide-vue-next'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import { apiService } from '../../services/api.service'

const users = ref([])
const loading = ref(true)
const searchQuery = ref('')

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await apiService.get('/admin/users')
    users.value = Array.isArray(res) ? res : (res.data || [])
  } catch (error) {
    console.error('Failed to load users:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUsers()
})

const handleToggleStatus = async (user) => {
  if (!confirm(`Bạn có chắc muốn ${user.is_active ? 'khoá' : 'kích hoạt'} tài khoản này?`)) return
  try {
    await apiService.put(`/admin/users/${user.id}/status`, { is_active: !user.is_active })
    user.is_active = !user.is_active
  } catch (error) {
    console.error('Failed to update status:', error)
  }
}

const handleToggleRole = async (user) => {
  if (!confirm(`Bạn có chắc muốn đổi quyền của tài khoản này thành ${user.role === 'admin' ? 'người dùng' : 'Admin'}?`)) return
  try {
    const newRole = user.role === 'admin' ? 'candidate' : 'admin'
    await apiService.put(`/admin/users/${user.id}/role`, { role: newRole })
    user.role = newRole
  } catch (error) {
    console.error('Failed to update role:', error)
  }
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Quản lý Người dùng</h1>
        <p class="text-slate-500 dark:text-slate-400 mt-1">Quản lý toàn bộ tài khoản trên hệ thống.</p>
      </div>
    </div>

    <Card class="p-6">
      <div class="flex gap-4 mb-6">
        <div class="relative w-72">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size="18" />
          <Input v-model="searchQuery" placeholder="Tìm tên, email..." class="pl-10 w-full" />
        </div>
      </div>

      <div v-if="loading" class="flex justify-center p-12">
        <div class="w-8 h-8 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
      </div>

      <div v-else-if="users.length === 0" class="text-center p-12 text-slate-500">
        Không tìm thấy người dùng nào.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700 text-slate-500 font-medium">
              <th class="pb-3 pl-4">Họ và tên</th>
              <th class="pb-3">Email</th>
              <th class="pb-3">Quyền</th>
              <th class="pb-3">Trạng thái</th>
              <th class="pb-3 text-right pr-4">Hành động</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id" class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50">
              <td class="py-4 pl-4 font-medium text-slate-800 dark:text-slate-200">{{ user.full_name }}</td>
              <td class="py-4 text-slate-600 dark:text-slate-400">{{ user.email }}</td>
              <td class="py-4">
                <Badge :variant="user.role === 'admin' ? 'primary' : 'secondary'">{{ user.role }}</Badge>
              </td>
              <td class="py-4">
                <Badge :variant="user.is_active ? 'success' : 'danger'">
                  {{ user.is_active ? 'Hoạt động' : 'Đã khoá' }}
                </Badge>
              </td>
              <td class="py-4 pr-4 flex justify-end gap-2">
                <Button variant="outline" size="sm" @click="handleToggleRole(user)" :title="user.role === 'admin' ? 'Hủy Admin' : 'Cấp Admin'">
                  <Shield size="16" />
                </Button>
                <Button :variant="user.is_active ? 'danger' : 'success'" size="sm" @click="handleToggleStatus(user)" :title="user.is_active ? 'Khoá tài khoản' : 'Mở khoá'">
                  <component :is="user.is_active ? UserX : UserCheck" size="16" />
                </Button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </Card>
  </div>
</template>
