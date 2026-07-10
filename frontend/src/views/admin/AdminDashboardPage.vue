<script setup>
import { ref, onMounted } from 'vue'
import { Users, Briefcase, FileText, CheckCircle } from 'lucide-vue-next'
import Card from '../../components/common/AppCard.vue'
import { apiService } from '../../services/api.service'

const stats = ref({
  total_users: 0,
  total_companies: 0,
  total_jobs: 0,
  total_interviews: 0
})
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await apiService.get('/admin/stats')
    stats.value = res || stats.value
  } catch (error) {
    console.error('Failed to load admin stats:', error)
  } finally {
    loading.value = false
  }
})

const statCards = [
  { key: 'total_users', label: 'Tổng số Người dùng', icon: Users, color: 'bg-blue-500' },
  { key: 'total_companies', label: 'Tổng số Công ty', icon: Briefcase, color: 'bg-indigo-500' },
  { key: 'total_jobs', label: 'Tổng số Việc làm', icon: FileText, color: 'bg-emerald-500' },
  { key: 'total_interviews', label: 'Tổng số Phỏng vấn', icon: CheckCircle, color: 'bg-purple-500' }
]
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Tổng quan Hệ thống (Admin)</h1>
      <p class="text-slate-500 dark:text-slate-400 mt-1">Theo dõi hoạt động của nền tảng.</p>
    </div>

    <div v-if="loading" class="flex justify-center p-12">
      <div class="w-8 h-8 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card v-for="stat in statCards" :key="stat.key" class="p-6 flex items-center gap-4 hover:shadow-lg transition-shadow">
        <div :class="[stat.color, 'w-12 h-12 rounded-xl flex items-center justify-center text-white shadow-lg']">
          <component :is="stat.icon" size="24" />
        </div>
        <div>
          <div class="text-sm font-medium text-slate-500 dark:text-slate-400 mb-1">{{ stat.label }}</div>
          <div class="text-2xl font-bold text-slate-800 dark:text-slate-100">{{ stats[stat.key] || 0 }}</div>
        </div>
      </Card>
    </div>
  </div>
</template>
