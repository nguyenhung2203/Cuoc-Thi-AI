<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../../services/api.service'
import Button from '../../components/common/AppButton.vue'
import { Briefcase, MapPin, ChevronLeft, ChevronRight, Sparkles, Heart, Banknote, Filter, ChevronDown, Info, AlertCircle } from 'lucide-vue-next'
import { langStore } from '../../stores/lang.store'
import { loadSavedJobIds, toggleSavedJob } from '../../utils/savedJobs'
import { formatSalaryTrieu, formatExperience, experienceFilterOptions } from '../../utils/formatters'

const router = useRouter()
const PAGE_SIZE = 12

const allJobs = ref([])
const loading = ref(true)
const loadError = ref('')
const keyword = ref('')
const currentPage = ref(1)

const selectedLocation = ref('')
const selectedSalary = ref('')
const selectedDepartment = ref('')
const selectedExperience = ref('')
const selectedWorkType = ref('')
const quickFilterActive = ref('Ngẫu nhiên')

const WORK_TYPE_LABELS = {
  full_time: 'Toàn thời gian',
  part_time: 'Bán thời gian',
  contract: 'Hợp đồng',
  intern: 'Thực tập',
  freelance: 'Freelance',
}

/** Chuẩn hóa địa điểm tin → cấp tỉnh/thành (không hiện huyện/xã). */
const PROVINCE_ALIASES = [
  // Gồm cả biến thể thiếu dấu: "Dăk Lăk", "Dak Lak", ...
  { match: /[đd][aăắấầậẫàáạãâằặẵ]k\s*l[aăắấầậẫàáạãâằặẵ]k|buôn\s*ma\s*thuột|buon\s*ma\s*thuot/i, name: 'Đắk Lắk' },
  { match: /lâm\s*đồng|lam\s*dong|da\s*lat|đà\s*lạt/i, name: 'Lâm Đồng' },
  { match: /gia\s*lai|pleiku/i, name: 'Gia Lai' },
  { match: /khánh\s*hòa|khanh\s*hoa|nha\s*trang/i, name: 'Khánh Hòa' },
  { match: /phú\s*yên|phu\s*yen|tuy\s*hòa/i, name: 'Phú Yên' },
  { match: /[đd][aăắấầậẫàáạãâằặẵ]k\s*n[oôóòọõơờớợỡ]ng/i, name: 'Đắk Nông' },
  { match: /kon\s*tum/i, name: 'Kon Tum' },
  { match: /hồ\s*chí\s*minh|tp\.?\s*hcm|ho\s*chi\s*minh|sài\s*gòn|sai\s*gon/i, name: 'TP. Hồ Chí Minh' },
  { match: /hà\s*nội|ha\s*noi/i, name: 'Hà Nội' },
  { match: /đà\s*nẵng|da\s*nang/i, name: 'Đà Nẵng' },
  { match: /cần\s*thơ|can\s*tho/i, name: 'Cần Thơ' },
  { match: /hải\s*phòng|hai\s*phong/i, name: 'Hải Phòng' },
  { match: /bình\s*định|binh\s*dinh/i, name: 'Bình Định' },
  { match: /quảng\s*nam|quang\s*nam/i, name: 'Quảng Nam' },
  { match: /quảng\s*ngãi|quang\s*ngai/i, name: 'Quảng Ngãi' },
]

const matchProvinceAlias = (text) => {
  const raw = String(text || '').normalize('NFC').trim()
  if (!raw) return ''
  for (const item of PROVINCE_ALIASES) {
    if (item.match.test(raw)) return item.name
  }
  return ''
}

const toProvince = (location) => {
  const raw = String(location || '').normalize('NFC').trim()
  if (!raw) return ''
  const direct = matchProvinceAlias(raw)
  if (direct) return direct
  // "Huyện X, Tỉnh Y" → lấy phần sau dấu phẩy cuối, rồi chuẩn hóa lại
  const parts = raw.split(',').map((p) => p.trim()).filter(Boolean)
  if (parts.length >= 2) {
    const tail = parts[parts.length - 1]
    return matchProvinceAlias(tail) || tail
  }
  return raw
}

const unwrap = (val) => {
  if (!val && val !== 0) return ''
  if (typeof val === 'object') {
    if ('String' in val) return val.Valid ? val.String : ''
    if ('Int64' in val) return val.Valid ? val.Int64 : ''
    if ('Float64' in val) return val.Valid ? val.Float64 : ''
    if ('Int32' in val) return val.Valid ? val.Int32 : ''
    if ('Bool' in val) return val.Valid ? val.Bool : ''
  }
  return val
}

const normalizeJob = (j) => {
  const location = unwrap(j.location)
  return {
    ...j,
    company_name: unwrap(j.company_name),
    company_logo_url: unwrap(j.company_logo_url),
    location,
    province: toProvince(location),
    employment_type: unwrap(j.employment_type),
    department: unwrap(j.department),
    level: String(unwrap(j.level) || '').toLowerCase(),
    salary_min: Number(unwrap(j.salary_min) || 0) || 0,
    salary_max: Number(unwrap(j.salary_max) || 0) || 0,
    currency: unwrap(j.currency),
  }
}

const countBy = (list, getter) => {
  const counts = {}
  list.forEach((job) => {
    const key = getter(job)
    if (!key) return
    counts[key] = (counts[key] || 0) + 1
  })
  return Object.keys(counts)
    .map((name) => ({ name, count: counts[name] }))
    .sort((a, b) => b.count - a.count || a.name.localeCompare(b.name, 'vi'))
}

const matchesSalary = (job, bucket) => {
  const min = job.salary_min
  const max = job.salary_max
  const hasSalary = min > 0 || max > 0
  if (bucket === 'Thỏa thuận') return !hasSalary
  if (!hasSalary) return false
  const lo = min > 0 ? min : max
  const hi = max > 0 ? max : min
  if (bucket === 'Dưới 10 triệu') return hi > 0 && hi < 10_000_000
  if (bucket === '10 - 15 triệu') return lo < 15_000_000 && hi >= 10_000_000
  if (bucket === '15 - 20 triệu') return lo < 20_000_000 && hi >= 15_000_000
  if (bucket === 'Trên 20 triệu') return hi >= 20_000_000
  return true
}

/** Facets tỉnh/thành — gộp huyện/xã vào cùng tỉnh. */
const locationStats = computed(() => countBy(allJobs.value, (j) => j.province))
const departmentStats = computed(() => countBy(allJobs.value, (j) => j.department))
const experiences = experienceFilterOptions()
const workTypeStats = computed(() =>
  countBy(allJobs.value, (j) => j.employment_type).map((item) => ({
    ...item,
    label: WORK_TYPE_LABELS[item.name] || item.name,
  }))
)

/** Đẩy tin Sơn TOA xuống cuối danh sách (tránh chiếm hết trang đầu). */
const isSonToaJob = (job) => /s[ơo]n\s*toa/i.test(String(job.company_name || ''))

const filteredJobs = computed(() => {
  const q = keyword.value.trim().toLowerCase()
  return allJobs.value
    .filter((job) => {
      if (selectedLocation.value && job.province !== selectedLocation.value) return false
      if (selectedDepartment.value && job.department !== selectedDepartment.value) return false
      if (selectedExperience.value && job.level !== selectedExperience.value) return false
      if (selectedWorkType.value && job.employment_type !== selectedWorkType.value) return false
      if (selectedSalary.value && !matchesSalary(job, selectedSalary.value)) return false
      if (q) {
        const hay = `${job.title} ${job.company_name} ${job.location} ${job.province} ${job.department} ${job.description || ''}`.toLowerCase()
        if (!hay.includes(q)) return false
      }
      return true
    })
    .slice()
    .sort((a, b) => Number(isSonToaJob(a)) - Number(isSonToaJob(b)))
})

const totalPages = computed(() => Math.max(1, Math.ceil(filteredJobs.value.length / PAGE_SIZE)))
const jobs = computed(() => {
  const page = Math.min(currentPage.value, totalPages.value)
  const start = (page - 1) * PAGE_SIZE
  return filteredJobs.value.slice(start, start + PAGE_SIZE)
})

const getSalaryDisplay = (job) => formatSalaryTrieu(job.salary_min, job.salary_max)
const getExperienceDisplay = (job) => formatExperience(job.level)

const loadCatalog = async () => {
  loading.value = true
  loadError.value = ''
  try {
    const collected = []
    let page = 1
    let pages = 1
    do {
      const res = await apiService.getWithMeta(`/public/all-jobs?page=${page}&page_size=100`)
      collected.push(...(res.data || []).map(normalizeJob))
      pages = res.meta?.total_pages || 1
      page += 1
    } while (page <= pages)
    allJobs.value = collected
    currentPage.value = 1
  } catch (error) {
    console.error('Failed to load jobs', error)
    allJobs.value = []
    loadError.value = error?.message || 'Không tải được danh sách việc làm. Vui lòng thử lại.'
  } finally {
    loading.value = false
  }
}

const savedJobIds = ref([])
const refreshSavedIds = () => {
  savedJobIds.value = loadSavedJobIds()
}

const handleToggleSaveJob = (job) => {
  const { list } = toggleSavedJob(job)
  savedJobIds.value = list.map((j) => j.id)
}

const goPage = (page) => {
  currentPage.value = Math.min(Math.max(1, page), totalPages.value)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const clearFilters = () => {
  selectedLocation.value = ''
  selectedSalary.value = ''
  selectedDepartment.value = ''
  selectedExperience.value = ''
  selectedWorkType.value = ''
  quickFilterActive.value = 'Ngẫu nhiên'
  keyword.value = ''
  currentPage.value = 1
}

const applyLocation = (name) => {
  selectedLocation.value = name
  quickFilterActive.value = name || 'Ngẫu nhiên'
  currentPage.value = 1
}

const handleSearch = () => {
  currentPage.value = 1
}

watch(
  [selectedLocation, selectedDepartment, selectedSalary, selectedExperience, selectedWorkType],
  () => {
    currentPage.value = 1
    if (selectedLocation.value) quickFilterActive.value = selectedLocation.value
    else if (quickFilterActive.value !== 'Ngẫu nhiên') quickFilterActive.value = 'Ngẫu nhiên'
  }
)

onMounted(() => {
  loadCatalog()
  refreshSavedIds()
})

const viewJob = (job) => {
  router.push(`/careers/${job.company_id}/jobs/${job.id}`)
}
</script>

<template>
  <div class="jb-page">
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
          <select v-model="selectedLocation" class="jb-select" @change="applyLocation(selectedLocation)">
            <option value="">Tỉnh / Thành ({{ allJobs.length }})</option>
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
            <option v-for="exp in experiences" :key="exp.value" :value="exp.value">{{ exp.label }}</option>
          </select>
          <ChevronDown :size="16" class="jb-select-icon" />
        </div>

        <div class="jb-select-wrap">
          <select v-model="selectedWorkType" class="jb-select">
            <option value="">Hình thức</option>
            <option v-for="type in workTypeStats" :key="type.name" :value="type.name">{{ type.label }} ({{ type.count }})</option>
          </select>
          <ChevronDown :size="16" class="jb-select-icon" />
        </div>
      </div>

      <!-- Quick Pills -->
      <div class="jb-quick-pills">
        <button 
          type="button"
          @click="clearFilters()"
          :class="['jb-pill', quickFilterActive === 'Ngẫu nhiên' && !selectedLocation ? 'active' : '']">
          Tất cả ({{ allJobs.length }})
        </button>
        <button 
          type="button"
          v-for="loc in locationStats.slice(0, 10)" :key="loc.name"
          @click="applyLocation(loc.name)"
          :class="['jb-pill', quickFilterActive === loc.name ? 'active' : '']">
          <span :title="`${loc.name} (${loc.count})`">{{ loc.name }} ({{ loc.count }})</span>
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
    <div v-else-if="loadError" class="jb-empty">
      <div class="jb-empty-icon"><AlertCircle :size="42" /></div>
      <h3 class="jb-empty-title">Không tải được việc làm</h3>
      <p class="jb-empty-desc">{{ loadError }}</p>
      <Button variant="primary" class="mt-4" @click="loadCatalog">Thử lại</Button>
    </div>

    <div v-else-if="filteredJobs.length === 0" class="jb-empty">
      <div class="jb-empty-icon"><Briefcase :size="42" /></div>
      <h3 class="jb-empty-title">Không có việc làm phù hợp bộ lọc</h3>
      <p class="jb-empty-desc">Thử bỏ bớt điều kiện hoặc chọn địa điểm khác.</p>
      <button @click="clearFilters()" class="link-more mt-5">Xóa bộ lọc</button>
    </div>

    <!-- Job List -->
    <div v-else class="jb-list-wrap space-y-8">
      <div class="flex items-center justify-between flex-wrap gap-3" style="min-width:0;max-width:100%">
        <h2 class="list-heading text-h1" :title="keyword ? `Kết quả cho ${keyword}` : (selectedLocation ? `Việc làm tại ${selectedLocation}` : 'Việc làm tốt nhất')">
          <span v-if="keyword">Kết quả cho "<span class="text-primary">{{ keyword }}</span>"</span>
          <span v-else-if="selectedLocation">Việc làm tại <span class="text-primary">{{ selectedLocation }}</span></span>
          <span v-else class="inline-flex items-center gap-2 min-w-0">Việc làm tốt nhất <span class="text-muted font-normal">|</span> Đề xuất bởi <Sparkles :size="18" class="text-accent" /> <span class="text-accent font-extrabold tracking-tight">ViệcLàmAI</span></span>
        </h2>
        <span class="page-chip">{{ filteredJobs.length }} việc làm · trang {{ currentPage }}/{{ totalPages }}</span>
      </div>

      <div class="jb-grid stagger">
        <!-- Job Card -->
        <div v-for="job in jobs" :key="job.id"
          @click="viewJob(job)"
          class="job-card-new">
          
          <!-- Save Job Button -->
          <button class="jb-save-btn" @click.stop="handleToggleSaveJob(job)" :class="{ 'is-saved': savedJobIds.includes(job.id) }">
            <Heart :size="20" :fill="savedJobIds.includes(job.id) ? 'currentColor' : 'none'" />
          </button>

          <div class="jb-card-top">
            <div class="jb-logo">
              <img :src="job.company_logo_url || '/images/logo.png'" :alt="job.company_name || 'Company'" @error="(e) => { e.target.src = '/images/logo.png' }" />
            </div>
            <div class="jb-info">
              <h3 class="jb-title" :title="job.title">
                {{ job.title }}
              </h3>
              <div class="jb-company" :title="job.company_name || ''">
                {{ job.company_name || 'Chưa cập nhật' }}
              </div>
            </div>
          </div>

          <div class="jb-tags">
            <span class="jb-tag jb-tag-salary" :title="getSalaryDisplay(job)">
              <Banknote :size="12" /><span class="jb-tag-text">{{ getSalaryDisplay(job) }}</span>
            </span>
            <span class="jb-tag jb-tag-loc" :title="job.location || job.province || ''">
              <MapPin :size="12" /><span class="jb-tag-text">{{ job.province || job.location || 'Không xác định' }}</span>
            </span>
            <span class="jb-tag" :title="getExperienceDisplay(job)">
              <Briefcase :size="12" /><span class="jb-tag-text">{{ getExperienceDisplay(job) }}</span>
            </span>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center items-center gap-3 mt-12">
        <button :disabled="currentPage === 1" @click="goPage(currentPage - 1)" class="page-btn">
          <ChevronLeft :size="16" /> {{ langStore.t('jobs', 'prevPage') }}
        </button>
        <span class="page-current">{{ currentPage }}</span>
        <span class="text-muted-sep">/ {{ totalPages }}</span>
        <button :disabled="currentPage === totalPages" @click="goPage(currentPage + 1)" class="page-btn">
          {{ langStore.t('jobs', 'nextPage') }} <ChevronRight :size="16" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.jb-page {
  width: 100%;
  max-width: 1120px;
  min-width: 0;
  margin: 0 auto;
  padding: 0 0 48px;
  box-sizing: border-box;
  overflow-x: hidden;
}
.jb-page *,
.jb-page *::before,
.jb-page *::after { box-sizing: border-box; }

/* ---- New Hero Styles ---- */
.hero-container { background: var(--surface-soft); border-radius: var(--radius-lg); padding: 40px 32px; margin-bottom: 32px; border: 1px solid var(--border); position: relative; overflow: hidden; }
.hero-layout { display: flex; flex-direction: column; gap: 32px; align-items: center; }
@media (min-width: 992px) { .hero-layout { flex-direction: row; justify-content: space-between; } }
.hero-content { flex: 1; z-index: 10; position: relative; width: 100%; }
.hero-image-wrapper { flex: 1; display: none; justify-content: center; position: relative; z-index: 1; max-width: 320px; }
@media (min-width: 768px) { .hero-image-wrapper { display: flex; } }
.hero-illustration { width: 100%; height: auto; object-fit: contain; filter: drop-shadow(0 20px 30px rgba(8, 145, 178, 0.15)); }

.search-box { display: flex; flex-wrap: wrap; gap: 12px; width: 100%; max-width: 700px; min-width: 0; background: var(--surface); padding: 8px; border-radius: var(--radius-lg); border: 1px solid var(--border); box-shadow: var(--shadow-sm); }
.search-input-wrap { position: relative; flex: 1 1 200px; min-width: 0; display: flex; align-items: center; }
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
.jb-filter-bar {
  width: 100%;
  max-width: 100%;
  overflow: hidden;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 16px;
  margin-bottom: 16px;
  box-shadow: var(--shadow-sm);
}
.jb-filter-row { display: flex; flex-wrap: wrap; align-items: center; gap: 12px; margin-bottom: 16px; padding-bottom: 16px; border-bottom: 1px solid var(--border); }
.jb-filter-label { display: flex; align-items: center; gap: 8px; font-size: 14px; font-weight: 500; color: var(--text-main); flex-shrink: 0; }
.jb-select-wrap { position: relative; min-width: 0; flex: 1 1 140px; max-width: 100%; }
@media (min-width: 900px) { .jb-select-wrap { flex: 1 1 0; } }
.jb-select {
  width: 100%;
  max-width: 100%;
  appearance: none;
  background: var(--surface-soft);
  border: 1px solid var(--border);
  color: var(--text-main);
  font-size: 14px;
  border-radius: var(--radius);
  padding: 10px 32px 10px 12px;
  outline: none;
  cursor: pointer;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
}
.jb-select:focus { border-color: var(--primary); }
.jb-select-icon { position: absolute; right: 12px; top: 12px; color: var(--text-muted); pointer-events: none; }

.jb-quick-pills {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 100%;
}
.jb-pill {
  flex: 0 1 auto;
  max-width: 11.5rem;
  padding: 6px 14px;
  border-radius: var(--radius-full);
  font-size: 13px;
  font-weight: 500;
  border: 1px solid var(--border);
  background: var(--surface-soft);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.jb-pill:first-child { max-width: none; }
.jb-pill:hover { background: var(--border); color: var(--text-main); }
.jb-pill.active { background: var(--primary); border-color: var(--primary); color: #fff; box-shadow: var(--shadow-sm); font-weight: 600; }

.jb-hint-bar { display: flex; align-items: flex-start; gap: 8px; width: 100%; max-width: 100%; background: var(--surface-soft); border: 1px solid var(--border); color: var(--text-secondary); font-size: 13px; font-weight: 500; padding: 12px 16px; border-radius: var(--radius); margin-bottom: 32px; overflow: hidden; }
.hint-icon { color: var(--accent); flex-shrink: 0; margin-top: 1px; }
.jb-hint-bar > :not(.hint-icon) { min-width: 0; }

/* ---- New Job Card Styles ---- */
.jb-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 16px;
  width: 100%;
}
@media (min-width: 768px) { .jb-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (min-width: 1024px) { .jb-grid { grid-template-columns: repeat(3, minmax(0, 1fr)); } }

.job-card-new {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 0;
  width: 100%;
  max-width: 100%;
  overflow: hidden;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 14px;
  cursor: pointer;
  position: relative;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}
.job-card-new:hover { border-color: var(--primary); box-shadow: var(--shadow-md); }

.jb-save-btn { position: absolute; bottom: 12px; right: 12px; color: var(--text-muted); cursor: pointer; transition: color 0.2s ease; background: transparent; border: none; z-index: 10; }
.jb-save-btn:hover, .jb-save-btn.is-saved { color: var(--danger); }

.jb-card-top {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  min-width: 0;
  width: 100%;
}

.jb-logo { width: 52px; height: 52px; flex-shrink: 0; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius); padding: 4px; display: flex; align-items: center; justify-content: center; overflow: hidden; }
.jb-logo img { max-width: 100%; max-height: 100%; object-fit: contain; }

.jb-info { flex: 1; min-width: 0; max-width: 100%; display: flex; flex-direction: column; overflow: hidden; }
.jb-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 4px;
  padding-right: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: color 0.2s ease;
}
.job-card-new:hover .jb-title { color: var(--primary); }
.jb-company {
  font-size: 13px;
  color: var(--text-secondary);
  padding-right: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.jb-tags {
  display: flex;
  flex-wrap: nowrap;
  align-items: center;
  gap: 4px;
  width: 100%;
  min-width: 0;
  max-width: 100%;
  padding-right: 28px; /* chỗ cho nút tim */
  overflow: hidden;
}
.jb-tag {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  min-width: 0;
  flex: 1 1 0;
  background: var(--surface-soft);
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 500;
  padding: 3px 6px;
  border-radius: 4px;
  overflow: hidden;
  white-space: nowrap;
}
.jb-tag-text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.jb-tag :deep(svg) { flex-shrink: 0; width: 11px; height: 11px; }
.jb-tag-salary { color: #e11d48; background: #fff1f2; font-weight: 600; }
.jb-tag.ai-tag { background: var(--accent-bg); color: var(--accent); border: 1px solid rgba(6,182,212,0.25); font-weight: 700; font-size: 11px; }
.jb-ai-badge { display: flex; align-items: center; gap: 8px; margin-top: 4px; }

.list-heading {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.page-chip { flex-shrink: 0; }
.jb-list-wrap { width: 100%; min-width: 0; max-width: 100%; overflow: hidden; }
</style>
