<script setup>
import { ref, onMounted } from 'vue'
import { Search, Building, CheckCircle, XCircle } from 'lucide-vue-next'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import { apiService } from '../../services/api.service'

const companies = ref([])
const loading = ref(true)
const searchQuery = ref('')

const fetchCompanies = async () => {
  loading.value = true
  try {
    const res = await apiService.get('/admin/companies')
    companies.value = Array.isArray(res) ? res : (res.data || [])
  } catch (error) {
    console.error('Failed to load companies:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCompanies()
})

const handleToggleStatus = async (company) => {
  if (!confirm(`Bạn có chắc muốn ${company.status === 'active' ? 'khoá' : 'kích hoạt'} công ty này?`)) return
  try {
    const newStatus = company.status === 'active' ? 'inactive' : 'active'
    await apiService.put(`/admin/companies/${company.id}/status`, { status: newStatus })
    company.status = newStatus
  } catch (error) {
    console.error('Failed to update status:', error)
  }
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Quản lý Công ty</h1>
        <p class="text-slate-500 dark:text-slate-400 mt-1">Danh sách công ty trên nền tảng.</p>
      </div>
    </div>

    <Card class="p-6">
      <div class="flex gap-4 mb-6">
        <div class="relative w-72">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size="18" />
          <Input v-model="searchQuery" placeholder="Tìm tên công ty..." class="pl-10 w-full" />
        </div>
      </div>

      <div v-if="loading" class="flex justify-center p-12">
        <div class="w-8 h-8 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
      </div>

      <div v-else-if="companies.length === 0" class="text-center p-12 text-slate-500">
        Không tìm thấy công ty nào.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-700 text-slate-500 font-medium">
              <th class="pb-3 pl-4">Tên công ty</th>
              <th class="pb-3">Website</th>
              <th class="pb-3">Lĩnh vực</th>
              <th class="pb-3">Trạng thái</th>
              <th class="pb-3 text-right pr-4">Hành động</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="company in companies" :key="company.id" class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50">
              <td class="py-4 pl-4 font-medium text-slate-800 dark:text-slate-200 flex items-center gap-2">
                <Building size="16" class="text-slate-400" />
                {{ company.name }}
              </td>
              <td class="py-4 text-slate-600 dark:text-slate-400">{{ company.website || '--' }}</td>
              <td class="py-4 text-slate-600 dark:text-slate-400">{{ company.industry || '--' }}</td>
              <td class="py-4">
                <Badge :variant="company.status === 'active' ? 'success' : 'danger'">
                  {{ company.status === 'active' ? 'Hoạt động' : 'Đã khoá' }}
                </Badge>
              </td>
              <td class="py-4 pr-4 flex justify-end gap-2">
                <Button :variant="company.status === 'active' ? 'danger' : 'success'" size="sm" @click="handleToggleStatus(company)" :title="company.status === 'active' ? 'Khoá công ty' : 'Kích hoạt'">
                  <component :is="company.status === 'active' ? XCircle : CheckCircle" size="16" />
                </Button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </Card>
  </div>
</template>
