<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../../services/api.service'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Briefcase, MapPin, Clock, Building2, Search, ChevronLeft, ChevronRight, Sparkles } from 'lucide-vue-next'

const router = useRouter()
const jobs = ref([])
const loading = ref(true)
const keyword = ref('')

const currentPage = ref(1)
const totalPages = ref(1)

const unwrap = (val) => {
  if (!val) return ''
  if (typeof val === 'object' && 'String' in val) {
    return val.Valid ? val.String : ''
  }
  return val
}

const fetchJobs = async (page = 1) => {
  loading.value = true
  currentPage.value = page
  try {
    const params = new URLSearchParams()
    params.append('page', page)
    params.append('page_size', 12)
    if (keyword.value) params.append('keyword', keyword.value)
    
    const res = await apiService.getWithMeta(`/public/all-jobs?${params.toString()}`)
    
    // Map to unwrap sql.NullString objects
    jobs.value = (res.data || []).map(j => ({
      ...j,
      company_name: unwrap(j.company_name),
      location: unwrap(j.location),
      employment_type: unwrap(j.employment_type),
      department: unwrap(j.department),
      level: unwrap(j.level),
    }))
    
    if (res.meta) {
      totalPages.value = res.meta.total_pages || 1
    }
  } catch (error) {
    console.error('Failed to load jobs', error)
    jobs.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchJobs(1))

const handleSearch = () => {
  fetchJobs(1)
}

const viewJob = (job) => {
  router.push(`/careers/${job.company_id}/jobs/${job.id}`)
}
</script>

<template>
  <div class="space-y-8 pb-12 max-w-7xl mx-auto">
    <!-- Hero -->
    <div class="brand-banner job-hero animate-rise">
      <div class="hero-content">
        <h1 class="hero-title">Khám phá cơ hội nghề nghiệp</h1>
        <p class="hero-desc">Tìm kiếm hàng ngàn việc làm phù hợp với kỹ năng và định hướng phát triển của bạn.</p>

        <div class="search-box">
          <div class="search-input-wrap">
            <Search :size="20" class="search-icon" />
            <input
              v-model="keyword"
              type="text"
              placeholder="Nhập chức danh, từ khóa hoặc công ty..."
              class="search-input"
              @keydown.enter="handleSearch"
            />
          </div>
          <button @click="handleSearch" class="search-btn sheen">Tìm việc ngay</button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="jb-spinner mb-4"></div>
      <p class="text-helper">Đang tìm kiếm việc làm phù hợp...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="jobs.length === 0" class="jb-empty">
      <div class="jb-empty-icon"><Briefcase :size="42" /></div>
      <h3 class="jb-empty-title">Chưa tìm thấy kết quả</h3>
      <p class="jb-empty-desc">Rất tiếc, chúng tôi không tìm thấy vị trí nào phù hợp với từ khóa của bạn lúc này. Vui lòng thử lại với từ khóa khác.</p>
      <button @click="keyword = ''; handleSearch()" class="link-more mt-5">Xóa tìm kiếm</button>
    </div>

    <!-- Job List -->
    <div v-else class="space-y-8">
      <div class="flex items-center justify-between flex-wrap gap-3">
        <h2 class="list-heading">
          <Sparkles v-if="!keyword" :size="18" class="text-accent" />
          <span v-if="keyword">Kết quả cho "<span class="text-primary">{{ keyword }}</span>"</span>
          <span v-else>Việc làm đề xuất cho bạn</span>
        </h2>
        <span class="page-chip">Trang {{ currentPage }} / {{ totalPages }}</span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 stagger">
        <!-- Job Card -->
        <div v-for="job in jobs" :key="job.id"
          @click="viewJob(job)"
          class="job-card card-elevate">

          <div class="flex justify-between items-start mb-4 gap-4">
            <h3 class="job-title">{{ job.title }}</h3>
            <span v-if="job.department" class="badge badge-info whitespace-nowrap">{{ job.department }}</span>
          </div>

          <div class="job-meta">
            <div v-if="job.company_name" class="job-meta-row">
              <Building2 :size="16" /> <span class="truncate">{{ job.company_name }}</span>
            </div>
            <div class="job-meta-row">
              <MapPin :size="16" /> <span class="truncate">{{ job.location || 'Bất kỳ' }}</span>
            </div>
            <div class="job-meta-row">
              <Clock :size="16" /> <span class="truncate">{{ job.employment_type || 'Full-time' }}</span>
            </div>
          </div>

          <div class="job-foot">
            <p class="job-desc">{{ job.description }}</p>
            <div class="job-link">Xem chi tiết <ChevronRight :size="16" /></div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center items-center gap-3 mt-12">
        <button :disabled="currentPage === 1" @click="fetchJobs(currentPage - 1)" class="page-btn">
          <ChevronLeft :size="16" /> Trang trước
        </button>
        <span class="page-current">{{ currentPage }}</span>
        <span class="text-muted-sep">/ {{ totalPages }}</span>
        <button :disabled="currentPage === totalPages" @click="fetchJobs(currentPage + 1)" class="page-btn">
          Trang sau <ChevronRight :size="16" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.job-hero { padding: 40px; }
@media (min-width: 768px) { .job-hero { padding: 48px; } }
.hero-content { position: relative; z-index: 1; max-width: 44rem; }
.hero-title { font-size: 32px; font-weight: 700; letter-spacing: -0.02em; margin-bottom: 12px; }
.hero-desc { color: rgba(255,255,255,0.9); font-size: 16px; margin-bottom: 28px; }

.search-box { display: flex; flex-direction: column; gap: 8px; background: #fff; padding: 8px; border-radius: var(--radius-lg); box-shadow: var(--shadow-lg); max-width: 46rem; }
@media (min-width: 768px) { .search-box { flex-direction: row; } }
.search-input-wrap { position: relative; flex: 1; display: flex; align-items: center; }
.search-icon { position: absolute; left: 14px; color: var(--text-muted); }
.search-input { width: 100%; border: none; outline: none; background: transparent; padding: 12px 14px 12px 42px; font-size: 15px; color: var(--text-main); font-family: var(--sans); }
.search-input::placeholder { color: var(--text-muted); }
.search-btn { background: var(--primary); color: #fff; border: none; font-weight: 700; padding: 12px 28px; border-radius: var(--radius); cursor: pointer; white-space: nowrap; transition: background 0.2s ease; }
.search-btn:hover { background: var(--primary-hover); }

.jb-spinner { width: 44px; height: 44px; border-radius: 50%; border: 3px solid var(--primary-light); border-top-color: var(--primary); animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.jb-empty { display: flex; flex-direction: column; align-items: center; text-align: center; padding: 72px 16px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); }
.jb-empty-icon { display: flex; align-items: center; justify-content: center; width: 88px; height: 88px; border-radius: 50%; background: var(--surface-soft); color: var(--text-muted); margin-bottom: 22px; }
.jb-empty-title { font-size: 22px; font-weight: 700; color: var(--text-main); margin-bottom: 8px; }
.jb-empty-desc { color: var(--text-secondary); max-width: 30rem; }

.list-heading { display: flex; align-items: center; gap: 8px; font-size: 18px; font-weight: 700; color: var(--text-main); }
.text-accent { color: var(--accent); }
.text-primary { color: var(--primary); }
.page-chip { font-size: 13px; color: var(--text-secondary); font-weight: 500; background: var(--surface-soft); padding: 5px 14px; border-radius: var(--radius-full); }

.job-card { display: flex; flex-direction: column; height: 100%; padding: 22px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); cursor: pointer; }
.job-title { font-size: 17px; font-weight: 700; color: var(--text-main); line-height: 1.35; flex: 1; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; transition: color 0.2s ease; }
.job-card:hover .job-title { color: var(--primary); }
.job-meta { display: flex; flex-direction: column; gap: 10px; margin-bottom: 20px; }
.job-meta-row { display: flex; align-items: center; gap: 10px; color: var(--text-secondary); font-size: 14px; }
.job-meta-row :deep(svg) { color: var(--text-muted); flex-shrink: 0; }
.job-foot { margin-top: auto; padding-top: 18px; border-top: 1px solid var(--border); }
.job-desc { color: var(--text-secondary); font-size: 14px; line-height: 1.55; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.job-link { display: flex; align-items: center; gap: 4px; margin-top: 14px; color: var(--primary); font-weight: 600; font-size: 14px; opacity: 0; transform: translateX(-4px); transition: all 0.25s ease; }
.job-card:hover .job-link { opacity: 1; transform: translateX(0); }

.page-btn { display: inline-flex; align-items: center; gap: 6px; padding: 10px 18px; border-radius: var(--radius); border: 1px solid var(--border); background: var(--surface); color: var(--text-secondary); font-weight: 500; cursor: pointer; transition: all 0.2s ease; }
.page-btn:hover:not(:disabled) { border-color: var(--primary); color: var(--primary); }
.page-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.page-current { display: inline-flex; align-items: center; padding: 8px 16px; border-radius: var(--radius); background: var(--primary-light); color: var(--primary); font-weight: 700; }
.text-muted-sep { color: var(--text-muted); }
</style>
