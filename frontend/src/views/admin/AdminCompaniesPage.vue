<script setup>
import { ref, onMounted, computed } from 'vue'
import Card from '../../components/common/AppCard.vue'
import { apiService } from '../../services/api.service'
import { Building, Search, Globe, Users, Briefcase, RefreshCw, ExternalLink, ShieldCheck } from 'lucide-vue-next'

const companies = ref([])
const loading = ref(false)
const searchQuery = ref('')

const displayedCompanies = computed(() => {
  if (!companies.value || !Array.isArray(companies.value)) return []
  if (!searchQuery.value) return companies.value
  const q = searchQuery.value.toLowerCase()
  return companies.value.filter(c => 
    (c.name && c.name.toLowerCase().includes(q)) ||
    (c.industry && c.industry.toLowerCase().includes(q)) ||
    (c.website && c.website.toLowerCase().includes(q))
  )
})

const fetchCompanies = async () => {
  loading.value = true
  try {
    const res = await apiService.get('/admin/companies')
    companies.value = Array.isArray(res) ? res : (res.data || [])
  } catch (error) {
    console.error('Failed to fetch companies', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCompanies()
})
</script>

<template>
  <div class="space-y-6 animate-fade-in">
    <!-- Header -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-white flex items-center gap-2.5">
          <Building size="26" class="text-blue-600 dark:text-blue-400" />
          Quản lý Doanh nghiệp & Công ty
        </h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
          Theo dõi danh sách các tổ chức, doanh nghiệp tham gia sử dụng nền tảng WeMake AI.
        </p>
      </div>
      <button 
        @click="fetchCompanies" 
        :disabled="loading"
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-200 rounded-xl font-medium text-sm transition-colors shrink-0 disabled:opacity-50"
      >
        <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới danh sách
      </button>
    </Card>

    <!-- Search Bar -->
    <Card class="p-4 rounded-2xl shadow-sm flex items-center justify-between gap-4">
      <div class="relative flex-1 max-w-md">
        <Search size="16" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
        <input 
          type="text" 
          v-model="searchQuery"
          placeholder="Tìm theo tên công ty, ngành nghề, website..." 
          class="w-full pl-9 pr-4 py-2 text-sm bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-800 dark:text-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
        />
      </div>
      <div class="text-xs font-semibold text-slate-500 dark:text-slate-400">
        Tổng số: <span class="text-blue-600 dark:text-blue-400 font-bold text-sm">{{ displayedCompanies.length }}</span> doanh nghiệp
      </div>
    </Card>

    <!-- Loading -->
    <Card v-if="loading" class="rounded-2xl p-12 text-center shadow-sm space-y-4">
      <div class="inline-block w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
      <p class="text-slate-500 dark:text-slate-400 text-sm font-medium">Đang tải hồ sơ các doanh nghiệp...</p>
    </Card>

    <!-- Grid View -->
    <div v-else-if="displayedCompanies.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="company in displayedCompanies" 
        :key="company.id"
        class="bg-[var(--surface)] rounded-2xl p-6 shadow-sm border border-[var(--border)] hover:shadow-md hover:border-blue-500/40 transition-all flex flex-col justify-between group"
      >
        <div class="space-y-4">
          <div class="flex items-start justify-between gap-3">
            <div class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-blue-500/10 to-blue-500/10 dark:from-blue-500/20 dark:to-blue-500/20 border border-blue-500/20 text-blue-600 dark:text-blue-400 font-extrabold flex items-center justify-center text-lg shrink-0 shadow-sm group-hover:scale-105 transition-transform">
              {{ company.name ? company.name.charAt(0).toUpperCase() : 'C' }}
            </div>
            <span class="px-2.5 py-1 rounded-full text-[11px] font-bold bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 border border-emerald-200/60 dark:border-emerald-800/40 flex items-center gap-1">
              <ShieldCheck size="12" /> Verified
            </span>
          </div>

          <div>
            <h3 class="font-bold text-slate-800 dark:text-white text-base group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors line-clamp-1" :title="company.name">
              {{ company.name || 'Công ty chưa đặt tên' }}
            </h3>
            <div class="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1.5 mt-1">
              <Briefcase size="13" class="text-slate-400 shrink-0" />
              <span>Ngành: <strong class="text-slate-700 dark:text-slate-300">{{ company.industry || 'Công nghệ / IT' }}</strong></span>
            </div>
          </div>

          <div class="pt-3 border-t border-slate-100 dark:border-slate-700/60 grid grid-cols-2 gap-2 text-xs">
            <div class="bg-slate-50 dark:bg-slate-900/50 p-2.5 rounded-xl border border-slate-200/50 dark:border-slate-700/50">
              <div class="text-slate-400 text-[10px] uppercase font-bold tracking-wider mb-0.5">Quy mô nhân sự</div>
              <div class="font-bold text-slate-700 dark:text-slate-300 flex items-center gap-1">
                <Users size="13" class="text-blue-500" />
                {{ company.size || '1 - 50' }} nhân viên
              </div>
            </div>
            <div class="bg-slate-50 dark:bg-slate-900/50 p-2.5 rounded-xl border border-slate-200/50 dark:border-slate-700/50">
              <div class="text-slate-400 text-[10px] uppercase font-bold tracking-wider mb-0.5">Lượt phỏng vấn</div>
              <div class="font-bold text-slate-700 dark:text-slate-300 flex items-center gap-1">
                <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span> Active AI
              </div>
            </div>
          </div>
        </div>

        <div class="pt-4 mt-4 border-t border-slate-100 dark:border-slate-700/60 flex items-center justify-between gap-2">
          <a 
            v-if="company.website" 
            :href="company.website.startsWith('http') ? company.website : `https://${company.website}`" 
            target="_blank" 
            class="text-xs font-semibold text-blue-600 dark:text-blue-400 hover:underline flex items-center gap-1 overflow-hidden text-ellipsis whitespace-nowrap max-w-[180px]"
          >
            <Globe size="13" class="shrink-0" />
            <span class="truncate">{{ company.website }}</span>
            <ExternalLink size="11" class="shrink-0" />
          </a>
          <span v-else class="text-xs text-slate-400 italic">Chưa có website</span>

          <button class="px-3 py-1.5 text-xs font-semibold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors shrink-0">
            Xem hồ sơ
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <Card v-else class="rounded-2xl p-16 text-center shadow-sm">
      <div class="w-16 h-16 bg-slate-100 dark:bg-slate-700/50 rounded-full flex items-center justify-center mx-auto mb-4 text-slate-400">
        <Building size="32" />
      </div>
      <h3 class="font-bold text-slate-800 dark:text-white text-base">Chưa có công ty nào</h3>
      <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
        {{ searchQuery ? 'Không có công ty nào khớp với từ khóa tìm kiếm.' : 'Hệ thống hiện chưa ghi nhận công ty hoặc doanh nghiệp nào được tạo.' }}
      </p>
    </Card>
  </div>
</template>
