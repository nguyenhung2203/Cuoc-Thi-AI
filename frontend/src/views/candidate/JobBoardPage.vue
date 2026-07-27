<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../../services/api.service'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Briefcase, MapPin, Clock, Building2, Search, ChevronLeft, ChevronRight, Sparkles, Heart, Banknote, Filter, ChevronDown, Info } from 'lucide-vue-next'
import { langStore } from '../../stores/lang.store'

const router = useRouter()
const jobs = ref([])
const loading = ref(true)
const keyword = ref('')

const currentPage = ref(1)
const totalPages = ref(1)

const selectedLocation = ref('')
const selectedSalary = ref('')
const selectedDepartment = ref('')
const selectedExperience = ref('')
const selectedWorkType = ref('')
const quickFilterActive = ref('Ngẫu nhiên')

// Compute stats from real data
const locationStats = computed(() => {
  const counts = {}
  jobs.value.forEach(job => {
    const loc = unwrap(job.location)
    if (loc) counts[loc] = (counts[loc] || 0) + 1
  })
  return Object.keys(counts).map(name => ({ name, count: counts[name] })).sort((a, b) => b.count - a.count)
})

const departmentStats = computed(() => {
  const counts = {}
  jobs.value.forEach(job => {
    const dep = unwrap(job.department)
    if (dep) counts[dep] = (counts[dep] || 0) + 1
  })
  return Object.keys(counts).map(name => ({ name, count: counts[name] })).sort((a, b) => b.count - a.count)
})

const experiences = computed(() => {
  const set = new Set()
  jobs.value.forEach(job => {
    const level = unwrap(job.level)
    if (level) set.add(level)
  })
  return Array.from(set).sort()
})

const workTypes = computed(() => {
  const set = new Set()
  jobs.value.forEach(job => {
    const type = unwrap(job.employment_type)
    if (type) set.add(type)
  })
  return Array.from(set).sort()
})

const getSalaryDisplay = (job) => {
  const min = unwrap(job.salary_min)
  const max = unwrap(job.salary_max)
  const curr = unwrap(job.currency) || 'VND'
  if (!min && !max) return 'Thỏa thuận'
  if (min && !max) return `Từ ${min} ${curr}`
  if (!min && max) return `Đến ${max} ${curr}`
  return `${min} - ${max} ${curr}`
}

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

const savedJobIds = ref([])
const loadSavedJobIds = () => {
  const list = localStorage.getItem('candidate_saved_jobs')
  if (list) {
    try {
      const parsed = JSON.parse(list)
      savedJobIds.value = parsed.map(j => j.id)
    } catch (e) {
      savedJobIds.value = []
    }
  } else {
    savedJobIds.value = []
  }
}

const toggleSaveJob = (job) => {
  const listStr = localStorage.getItem('candidate_saved_jobs')
  let list = []
  if (listStr) {
    try { list = JSON.parse(listStr) } catch (e) { list = [] }
  }
  
  const isSaved = list.some(item => item.id === job.id)
  if (isSaved) {
    list = list.filter(item => item.id !== job.id)
    savedJobIds.value = savedJobIds.value.filter(id => id !== job.id)
  } else {
    const newItem = {
      id: job.id,
      company_id: job.company_id,
      title: job.title,
      company_name: job.company_name,
      location: job.location,
      salary_min: { Valid: job.salary_min != null, Int64: job.salary_min || 0 },
      salary_max: { Valid: job.salary_max != null, Int64: job.salary_max || 0 },
      currency: { Valid: true, String: job.currency || 'VND' },
      saved_at: new Date().toISOString()
    }
    list.push(newItem)
    savedJobIds.value.push(job.id)
  }
  localStorage.setItem('candidate_saved_jobs', JSON.stringify(list))
}

onMounted(() => {
  fetchJobs(1)
  loadSavedJobIds()
})

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
    <div class="hero-container animate-rise">
      <div class="hero-layout">
        <div class="hero-content">
          <h1 class="page-title text-h1" style="margin-bottom: 12px">Tìm việc cùng Trợ lý AI</h1>
          <p class="page-subtitle" style="max-width: 500px; margin-bottom: 32px; font-size: 16px;">Trợ lý AI sẽ giúp bạn phân tích kỹ năng và tìm ra công việc phù hợp nhất với định hướng của bạn.</p>
  
          <div class="search-box">
            <div class="search-input-wrap">
              <Sparkles class="search-icon text-accent" :size="20" />
              <input type="text" v-model="keyword" class="search-input" placeholder="Hỏi AI việc làm (VD: Frontend ít áp lực)..." @keyup.enter="handleSearch" />
            </div>
            <button class="search-btn sheen" @click="handleSearch">Tìm kiếm thông minh</button>
          </div>
        </div>
        <div class="hero-image-wrapper">
          <img src="/images/hero_network.png" alt="AI Match" class="hero-illustration tilt-3d" />
        </div>
      </div>
    </div>

    <!-- Filter Bar & Quick Pills -->
    <div class="jb-filter-bar">
      <div class="jb-filter-row">
        <!-- Dropdowns -->
        <div class="jb-filter-label">
          <Filter :size="16" class="text-muted" />
          <span>Lọc theo:</span>
        </div>
        
        <div class="jb-select-wrap">
          <select v-model="selectedLocation" class="jb-select">
            <option value="">Địa điểm</option>
            <option v-for="loc in locationStats" :key="loc.name" :value="loc.name">{{ loc.name }} ({{ loc.count }})</option>
          </select>
          <ChevronDown :size="16" class="jb-select-icon" />
        </div>

        <div class="jb-select-wrap">
          <select v-model="selectedDepartment" class="jb-select">
            <option value="">Ngành nghề</option>
            <option v-for="dep in departmentStats" :key="dep.name" :value="dep.name">{{ dep.name }} ({{ dep.count }})</option>
          </select>
          <ChevronDown :size="16" class="jb-select-icon" />
        </div>

        <div class="jb-select-wrap">
          <select v-model="selectedSalary" class="jb-select">
            <option value="">Mức lương</option>
            <option value="Dưới 10 triệu">Dưới 10 triệu</option>
            <option value="10 - 15 triệu">10 - 15 triệu</option>
            <option value="15 - 20 triệu">15 - 20 triệu</option>
            <option value="Trên 20 triệu">Trên 20 triệu</option>
            <option value="Thỏa thuận">Thỏa thuận</option>
          </select>
          <ChevronDown :size="16" class="jb-select-icon" />
        </div>

        <div class="jb-select-wrap">
          <select v-model="selectedExperience" class="jb-select">
            <option value="">Kinh nghiệm</option>
            <option v-for="exp in experiences" :key="exp" :value="exp">{{ exp }}</option>
          </select>
          <ChevronDown :size="16" class="jb-select-icon" />
        </div>

        <div class="jb-select-wrap">
          <select v-model="selectedWorkType" class="jb-select">
            <option value="">Hình thức</option>
            <option v-for="type in workTypes" :key="type" :value="type">{{ type }}</option>
          </select>
          <ChevronDown :size="16" class="jb-select-icon" />
        </div>
      </div>

      <!-- Quick Pills -->
      <div class="jb-quick-pills">
        <button 
          @click="quickFilterActive = 'Ngẫu nhiên'; selectedLocation = ''; handleSearch()"
          :class="['jb-pill', quickFilterActive === 'Ngẫu nhiên' ? 'active' : '']">
          Ngẫu nhiên
        </button>
        <button 
          v-for="loc in locationStats.slice(0, 8)" :key="loc.name"
          @click="quickFilterActive = loc.name; selectedLocation = loc.name; handleSearch()"
          :class="['jb-pill', quickFilterActive === loc.name ? 'active' : '']">
          {{ loc.name }}
        </button>
      </div>
    </div>

    <!-- Hint Bar -->
    <div class="jb-hint-bar">
      <Info :size="16" class="hint-icon" /> Gợi ý: Bấm vào xem công việc để Trợ lý AI phân tích chi tiết mức độ phù hợp của bạn!
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="jb-spinner mb-4"></div>
      <p class="text-helper">{{ langStore.t('jobs', 'loadingText') }}</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="jobs.length === 0" class="jb-empty">
      <div class="jb-empty-icon"><Briefcase :size="42" /></div>
      <h3 class="jb-empty-title">{{ langStore.t('jobs', 'emptyTitle') }}</h3>
      <p class="jb-empty-desc">{{ langStore.t('jobs', 'emptyDesc') }}</p>
      <button @click="keyword = ''; handleSearch()" class="link-more mt-5">{{ langStore.t('jobs', 'clearSearch') }}</button>
    </div>

    <!-- Job List -->
    <div v-else class="space-y-8">
      <div class="flex items-center justify-between flex-wrap gap-3">
        <h2 class="list-heading text-h1">
          <span v-if="keyword">Kết quả tìm kiếm cho "<span class="text-primary">{{ keyword }}</span>"</span>
          <span v-else class="flex items-center gap-2">Việc làm tốt nhất <span class="text-muted font-normal">|</span> Đề xuất bởi <Sparkles :size="18" class="text-accent" /> <span class="text-accent font-extrabold tracking-tight">WEMAKE AI</span></span>
        </h2>
        <span class="page-chip">{{ langStore.t('jobs', 'pageChip') }} {{ currentPage }} / {{ totalPages }}</span>
      </div>

      <div class="jb-grid stagger">
        <!-- Job Card -->
        <div v-for="job in jobs" :key="job.id"
          @click="viewJob(job)"
          class="job-card-new">
          
          <!-- Save Job Button -->
          <button class="jb-save-btn" @click.stop="toggleSaveJob(job)" :class="{ 'is-saved': savedJobIds.includes(job.id) }">
            <Heart :size="20" :fill="savedJobIds.includes(job.id) ? 'currentColor' : 'none'" />
          </button>

          <!-- Company Logo -->
          <div class="jb-logo">
            <img src="/images/logo.png" alt="Company Logo" />
          </div>

          <!-- Job Info -->
          <div class="jb-info">
            <h3 class="jb-title">
              {{ job.title }}
            </h3>
            <div class="jb-company">
              {{ unwrap(job.company_name) || 'Công ty TNHH WeMake' }}
            </div>
            
            <div class="jb-tags">
              <span class="jb-tag">
                <Banknote :size="12" /> {{ getSalaryDisplay(job) }}
              </span>
              <span class="jb-tag">
                <MapPin :size="12" /> {{ unwrap(job.location) || 'Không xác định' }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center items-center gap-3 mt-12">
        <button :disabled="currentPage === 1" @click="fetchJobs(currentPage - 1)" class="page-btn">
          <ChevronLeft :size="16" /> {{ langStore.t('jobs', 'prevPage') }}
        </button>
        <span class="page-current">{{ currentPage }}</span>
        <span class="text-muted-sep">/ {{ totalPages }}</span>
        <button :disabled="currentPage === totalPages" @click="fetchJobs(currentPage + 1)" class="page-btn">
          {{ langStore.t('jobs', 'nextPage') }} <ChevronRight :size="16" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ---- New Hero Styles ---- */
.hero-container { background: var(--surface-soft); border-radius: var(--radius-lg); padding: 40px 32px; margin-bottom: 32px; border: 1px solid var(--border); position: relative; overflow: hidden; }
.hero-layout { display: flex; flex-direction: column; gap: 32px; align-items: center; }
@media (min-width: 992px) { .hero-layout { flex-direction: row; justify-content: space-between; } }
.hero-content { flex: 1; z-index: 10; position: relative; width: 100%; }
.hero-image-wrapper { flex: 1; display: none; justify-content: center; position: relative; z-index: 1; max-width: 320px; }
@media (min-width: 768px) { .hero-image-wrapper { display: flex; } }
.hero-illustration { width: 100%; height: auto; object-fit: contain; filter: drop-shadow(0 20px 30px rgba(8, 145, 178, 0.15)); }

.search-box { display: flex; gap: 12px; max-width: 700px; background: var(--surface); padding: 8px; border-radius: var(--radius-lg); border: 1px solid var(--border); box-shadow: var(--shadow-sm); }
@media (min-width: 768px) { .search-box { flex-direction: row; } }
.search-input-wrap { position: relative; flex: 1; display: flex; align-items: center; }
.search-icon { position: absolute; left: 14px; color: var(--text-muted); }
.search-input { width: 100%; border: none; outline: none; background: transparent; padding: 12px 14px 12px 42px; font-size: 15px; color: var(--text-main); font-family: var(--sans); }
.search-input::placeholder { color: var(--text-muted); }
.search-btn { background: var(--primary); color: #fff; border: none; font-weight: 700; padding: 12px 28px; border-radius: var(--radius); cursor: pointer; white-space: nowrap; box-shadow: var(--shadow-sm); transition: all 0.2s ease; }
.search-btn:hover { background: var(--primary-hover); box-shadow: var(--shadow-md); transform: translateY(-1px); }

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

/* ---- New Filter Bar Styles ---- */
.jb-filter-bar { background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); padding: 16px; margin-bottom: 16px; box-shadow: var(--shadow-sm); }
.jb-filter-row { display: flex; flex-wrap: wrap; align-items: center; gap: 16px; margin-bottom: 16px; padding-bottom: 16px; border-bottom: 1px solid var(--border); }
.jb-filter-label { display: flex; align-items: center; gap: 8px; font-size: 14px; font-weight: 500; color: var(--text-main); }
.jb-select-wrap { position: relative; min-width: 150px; flex: 1; }
@media (min-width: 640px) { .jb-select-wrap { flex: none; } }
.jb-select { width: 100%; appearance: none; background: var(--surface-soft); border: 1px solid var(--border); color: var(--text-main); font-size: 14px; border-radius: var(--radius); padding: 10px 16px; outline: none; cursor: pointer; transition: border-color 0.2s ease; }
.jb-select:focus { border-color: var(--primary); }
.jb-select-icon { position: absolute; right: 12px; top: 12px; color: var(--text-muted); pointer-events: none; }

.jb-quick-pills { display: flex; align-items: center; gap: 8px; overflow-x: auto; padding-bottom: 4px; }
.jb-quick-pills::-webkit-scrollbar { display: none; }
.jb-pill { padding: 6px 16px; border-radius: var(--radius-full); font-size: 14px; font-weight: 500; white-space: nowrap; border: 1px solid var(--border); background: var(--surface-soft); color: var(--text-secondary); cursor: pointer; transition: all 0.2s ease; }
.jb-pill:hover { background: var(--border); color: var(--text-main); }
.jb-pill.active { background: var(--primary); border-color: var(--primary); color: #fff; box-shadow: var(--shadow-sm); font-weight: 600; }

.jb-hint-bar { display: flex; align-items: center; gap: 8px; background: var(--surface-soft); border: 1px solid var(--border); color: var(--text-secondary); font-size: 13px; font-weight: 500; padding: 12px 16px; border-radius: var(--radius); margin-bottom: 32px; }
.hint-icon { color: var(--accent); flex-shrink: 0; }

/* ---- New Job Card Styles ---- */
.jb-grid { display: grid; grid-template-columns: 1fr; gap: 20px; }
@media (min-width: 768px) { .jb-grid { grid-template-columns: repeat(2, 1fr); } }
@media (min-width: 1024px) { .jb-grid { grid-template-columns: repeat(3, 1fr); } }

.job-card-new { display: flex; gap: 16px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); padding: 16px; cursor: pointer; position: relative; transition: all 0.2s ease; }
.job-card-new:hover { border-color: var(--primary); box-shadow: var(--shadow-md); transform: translateY(-2px); }

.jb-save-btn { position: absolute; bottom: 16px; right: 16px; color: var(--text-muted); cursor: pointer; transition: color 0.2s ease; background: transparent; border: none; z-index: 10; }
.jb-save-btn:hover, .jb-save-btn.is-saved { color: var(--danger); }

.jb-logo { width: 64px; height: 64px; flex-shrink: 0; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); padding: 4px; display: flex; align-items: center; justify-content: center; }
.jb-logo img { max-width: 100%; max-height: 100%; object-fit: contain; }

.jb-info { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.jb-title { font-size: 15px; font-weight: 700; color: var(--text-main); margin-bottom: 4px; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; transition: color 0.2s ease; padding-right: 24px; }
.job-card-new:hover .jb-title { color: var(--primary); }
.jb-company { font-size: 13px; color: var(--text-secondary); margin-bottom: 12px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.jb-tags { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 8px; margin-top: auto; }
.jb-tag { display: inline-flex; align-items: center; gap: 4px; background: var(--surface-soft); color: var(--text-secondary); font-size: 12px; font-weight: 500; padding: 4px 8px; border-radius: 4px; }
.jb-tag.ai-tag { background: var(--accent-bg); color: var(--accent); border: 1px solid rgba(6,182,212,0.25); font-weight: 700; font-size: 11px; }
.jb-ai-badge { display: flex; align-items: center; gap: 8px; margin-top: 4px; }
</style>
