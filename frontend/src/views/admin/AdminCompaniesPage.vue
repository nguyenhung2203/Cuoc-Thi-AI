<script setup>
import { ref, onMounted, computed } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Modal from '../../components/common/AppModal.vue'
import { apiService } from '../../services/api.service'
import { Building, Search, Globe, Users, Briefcase, RefreshCw, ExternalLink, ShieldCheck } from 'lucide-vue-next'

const companies = ref([])
const loading = ref(false)
const fieldValue = (value) => value && typeof value === 'object' && 'String' in value
  ? (value.Valid ? String(value.String || '') : '')
  : (value == null ? '' : String(value))

const normalizeCompany = (company) => ({
  ...company,
  name: fieldValue(company?.name),
  industry: fieldValue(company?.industry),
  size: fieldValue(company?.size),
  website: fieldValue(company?.website)
})

const searchQuery = ref('')
const selectedCompany = ref(null)
const showCompanyModal = ref(false)
const openCompanyProfile = (company) => {
  selectedCompany.value = company
  showCompanyModal.value = true
}
const closeCompanyProfile = () => {
  showCompanyModal.value = false
  selectedCompany.value = null
}

const displayedCompanies = computed(() => {
  if (!companies.value || !Array.isArray(companies.value)) return []
  if (!searchQuery.value) return companies.value
  const q = searchQuery.value.toLowerCase()
  return q ? companies.value.filter(c => [c.name, c.industry, c.website].some(v => String(v || '').toLowerCase().includes(q))) : companies.value
})

const fetchCompanies = async () => {
  loading.value = true
  try {
    const res = await apiService.get('/admin/companies')
    companies.value = (Array.isArray(res) ? res : (res?.data || [])).map(normalizeCompany)
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
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm border border-[var(--border)]">
      <div>
        <h1 class="text-h1 flex items-center gap-2.5">
          <Building size="26" class="text-[var(--primary)]" />
          Quản lý Doanh nghiệp & Công ty
        </h1>
        <p class="text-[var(--text-secondary)] text-sm mt-1">
          Theo dõi danh sách các tổ chức, doanh nghiệp tham gia sử dụng nền tảng ViệcLàmAI.
        </p>
      </div>
      <button 
        @click="fetchCompanies" 
        :disabled="loading"
        class="inline-flex items-center gap-2 px-4 py-2.5 bg-[var(--surface-soft)] hover:bg-[var(--border)] text-[var(--text-main)] rounded-xl font-semibold text-sm transition-colors shrink-0 disabled:opacity-50"
      >
        <RefreshCw size="16" :class="{ 'animate-spin': loading }" /> Làm mới danh sách
      </button>
    </Card>

    <!-- Search Bar -->
    <Card class="p-4 rounded-2xl shadow-sm border border-[var(--border)] flex items-center justify-between gap-4">
      <div class="relative flex-1 max-w-md">
        <Search size="16" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-secondary)]" />
        <input 
          type="text" 
          v-model="searchQuery"
          placeholder="Tìm theo tên công ty, ngành nghề, website..." 
          class="w-full pl-9 pr-4 py-2 text-sm bg-[var(--surface)] border border-[var(--border)] rounded-xl text-[var(--text-main)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] transition-all"
        />
      </div>
      <div class="text-xs font-semibold text-[var(--text-secondary)]">
        Tổng số: <span class="text-[var(--primary)] font-bold text-sm">{{ displayedCompanies.length }}</span> doanh nghiệp
      </div>
    </Card>

    <!-- Loading -->
    <Card v-if="loading" class="rounded-2xl p-12 text-center border border-[var(--border)] shadow-sm space-y-4">
      <div class="inline-block w-8 h-8 border-4 border-[var(--primary)] border-t-transparent rounded-full animate-spin"></div>
      <p class="text-[var(--text-secondary)] text-sm font-medium">Đang tải hồ sơ các doanh nghiệp...</p>
    </Card>

    <!-- Grid View -->
    <div v-else-if="displayedCompanies.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="company in displayedCompanies" 
        :key="company.id"
        class="bg-[var(--surface)] rounded-2xl p-6 shadow-sm border border-[var(--border)] hover:shadow-md hover:border-[var(--primary)]/30 transition-all flex flex-col justify-between group"
      >
        <div class="space-y-4">
          <div class="flex items-start justify-between gap-3">
            <div class="w-12 h-12 rounded-2xl bg-[var(--primary-light)] border border-blue-200/50 text-[var(--primary)] font-extrabold flex items-center justify-center text-lg shrink-0 shadow-xs group-hover:scale-105 transition-transform">
              {{ company.name ? company.name.charAt(0).toUpperCase() : 'C' }}
            </div>
            <span class="px-2.5 py-1 rounded-full text-[11px] font-bold bg-[var(--success)]/10 text-[var(--success)] border border-[var(--success)]/20 flex items-center gap-1">
              <ShieldCheck size="12" /> Verified
            </span>
          </div>

          <div>
            <h3 class="font-bold text-[var(--text-main)] text-base group-hover:text-[var(--primary)] transition-colors line-clamp-1" :title="company.name">
              {{ company.name || 'Công ty chưa đặt tên' }}
            </h3>
            <div class="text-xs text-[var(--text-secondary)] flex items-center gap-1.5 mt-1">
              <Briefcase size="13" class="text-[var(--text-secondary)] shrink-0" />
              <span>Ngành: <strong class="text-[var(--text-main)]">{{ company.industry || 'Công nghệ / IT' }}</strong></span>
            </div>
          </div>

          <div class="pt-3 border-t border-[var(--border)] grid grid-cols-2 gap-2 text-xs">
            <div class="bg-[var(--background)] p-2.5 rounded-xl border border-[var(--border)]">
              <div class="text-[var(--text-secondary)] text-[10px] uppercase font-bold tracking-wider mb-0.5">Quy mô nhân sự</div>
              <div class="font-bold text-[var(--text-main)] flex items-center gap-1">
                <Users size="13" class="text-[var(--primary)]" />
                {{ company.size || '1 - 50' }} nhân viên
              </div>
            </div>
            <div class="bg-[var(--background)] p-2.5 rounded-xl border border-[var(--border)]">
              <div class="text-[var(--text-secondary)] text-[10px] uppercase font-bold tracking-wider mb-0.5">Lượt phỏng vấn</div>
              <div class="font-bold text-[var(--text-main)] flex items-center gap-1">
                <span class="w-2 h-2 rounded-full bg-[var(--success)] animate-pulse"></span> Active AI
              </div>
            </div>
          </div>
        </div>

        <div class="pt-4 mt-4 border-t border-[var(--border)] flex items-center justify-between gap-2">
          <a 
            v-if="company.website" 
            :href="/^https?:\/\//i.test(company.website) ? company.website : `https://${company.website}`" 
            target="_blank" 
            class="text-xs font-semibold text-[var(--accent)] hover:underline flex items-center gap-1 overflow-hidden text-ellipsis whitespace-nowrap max-w-[180px]"
          >
            <Globe size="13" class="shrink-0" />
            <span class="truncate">{{ company.website }}</span>
            <ExternalLink size="11" class="shrink-0" />
          </a>
          <span v-else class="text-xs text-[var(--text-secondary)] italic">Chưa có website</span>

          <button @click="openCompanyProfile(company)" class="inline-flex items-center gap-1.5 rounded-xl border border-blue-200 bg-blue-50 px-3.5 py-2 text-xs font-bold text-blue-700 shadow-sm transition-all hover:-translate-y-0.5 hover:border-blue-300 hover:bg-blue-100 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-300">
            <Building size="14" /> Xem hồ sơ
          </button>
        </div>
      </div>
    </div>

    <Modal :isOpen="showCompanyModal" title="Hồ sơ công ty" size="lg" @close="closeCompanyProfile">
      <div v-if="selectedCompany" class="-m-1 overflow-hidden rounded-2xl bg-white">
        <div class="relative overflow-hidden bg-slate-950 px-7 py-7 text-white">
          <div class="absolute -right-10 -top-16 h-48 w-48 rounded-full bg-blue-500/30 blur-2xl"></div>
          <div class="absolute -bottom-24 left-1/2 h-48 w-48 rounded-full bg-cyan-400/20 blur-3xl"></div>
          <div class="relative flex items-center gap-5">
            <div class="flex h-20 w-20 shrink-0 items-center justify-center rounded-2xl border border-white/20 bg-white text-3xl font-extrabold text-blue-700 shadow-xl">
              {{ selectedCompany.name?.charAt(0)?.toUpperCase() || 'C' }}
            </div>
            <div class="min-w-0">
              <span class="mb-2 inline-flex items-center rounded-full bg-emerald-400/15 px-3 py-1 text-[11px] font-bold text-emerald-300">✓ Đã xác minh</span>
              <h3 class="truncate text-2xl font-extrabold tracking-tight">{{ selectedCompany.name || 'Công ty chưa đặt tên' }}</h3>
              <p class="mt-1 text-sm text-slate-300">Thông tin doanh nghiệp trên nền tảng ViệcLàmAI</p>
            </div>
          </div>
        </div>
        <div class="grid gap-6 p-7 md:grid-cols-[1fr_210px]">
          <div>
            <h4 class="mb-4 text-xs font-extrabold uppercase tracking-[0.18em] text-slate-400">Thông tin doanh nghiệp</h4>
            <div class="grid gap-3 sm:grid-cols-2">
              <div class="rounded-2xl border border-slate-200 p-4"><p class="text-xs text-slate-400">Ngành nghề</p><p class="mt-1 font-bold text-slate-800">{{ selectedCompany.industry || 'Chưa cập nhật' }}</p></div>
              <div class="rounded-2xl border border-slate-200 p-4"><p class="text-xs text-slate-400">Quy mô nhân sự</p><p class="mt-1 font-bold text-slate-800">{{ selectedCompany.size || 'Chưa cập nhật' }}</p></div>
              <div class="rounded-2xl border border-slate-200 p-4 sm:col-span-2"><p class="text-xs text-slate-400">Website</p><a v-if="selectedCompany.website" :href="/^https?:\/\//i.test(selectedCompany.website) ? selectedCompany.website : `https://${selectedCompany.website}`" target="_blank" class="mt-1 inline-flex max-w-full items-center gap-1 truncate font-bold text-blue-600 hover:underline">{{ selectedCompany.website }} <ExternalLink size="13" /></a><p v-else class="mt-1 font-bold text-slate-800">Chưa cập nhật</p></div>
            </div>
          </div>
          <div class="rounded-2xl bg-blue-50 p-5">
            <h4 class="mb-4 text-xs font-extrabold uppercase tracking-wider text-blue-900">Tổng quan</h4>
            <div class="space-y-4 text-sm"><div><p class="text-xs text-blue-500">Trạng thái</p><span class="mt-1 inline-flex rounded-full bg-emerald-100 px-2.5 py-1 text-xs font-bold text-emerald-700">Đang hoạt động</span></div><div><p class="text-xs text-blue-500">Mã định danh</p><p class="mt-1 truncate font-bold text-slate-700" :title="selectedCompany.slug">{{ selectedCompany.slug || 'Chưa cập nhật' }}</p></div></div>
          </div>
        </div>
      </div>
    </Modal>

    <Card v-if="!loading && displayedCompanies.length === 0" class="rounded-2xl p-16 text-center border border-[var(--border)] shadow-sm">
      <div class="w-16 h-16 bg-[var(--surface)] rounded-full flex items-center justify-center mx-auto mb-4 text-[var(--text-secondary)]">
        <Building size="32" />
      </div>
      <h3 class="font-bold text-[var(--text-main)] text-base">Chưa có công ty nào</h3>
      <p class="text-[var(--text-secondary)] text-sm mt-1">
        {{ searchQuery ? 'Không có công ty nào khớp với từ khóa tìm kiếm.' : 'Hệ thống hiện chưa ghi nhận công ty hoặc doanh nghiệp nào được tạo.' }}
      </p>
    </Card>
  </div>
</template>
