<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiService } from '../../services/api.service'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Briefcase, MapPin, Clock, Building2, Search, ChevronLeft, ChevronRight } from 'lucide-vue-next'

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
  <div class="job-board-page">
    <div style="margin-bottom: 32px;">
      <h1 style="font-size: 28px; margin: 0 0 8px 0;">Tìm việc làm</h1>
      <p style="color: var(--text-secondary);">Khám phá các vị trí đang tuyển dụng trên hệ thống.</p>
    </div>

    <!-- Search -->
    <div style="display: flex; gap: 12px; margin-bottom: 32px;">
      <div style="flex: 1; position: relative;">
        <Search size="16" color="var(--text-muted)" style="position: absolute; left: 12px; top: 50%; transform: translateY(-50%);" />
        <input 
          v-model="keyword"
          type="text"
          placeholder="Tìm theo tên công việc, mô tả..."
          style="width: 100%; padding: 10px 12px 10px 36px; border: 1px solid var(--border); border-radius: var(--radius-md); font-size: 14px; background: var(--surface); color: var(--text-main); outline: none;"
          @keydown.enter="handleSearch"
        />
      </div>
      <Button @click="handleSearch">Tìm kiếm</Button>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="text-align: center; padding: 40px;">
      <div class="spinner" style="margin: 0 auto;"></div>
    </div>

    <!-- Empty -->
    <div v-else-if="jobs.length === 0" style="text-align: center; padding: 60px; background: var(--surface); border-radius: var(--radius-lg); border: 1px solid var(--border);">
      <Briefcase size="48" color="var(--border)" style="margin-bottom: 16px;" />
      <h3 style="font-size: 18px; margin-bottom: 8px;">Chưa có vị trí nào đang tuyển</h3>
      <p style="color: var(--text-secondary);">Vui lòng quay lại sau nhé.</p>
    </div>

    <!-- Job List -->
    <div v-else>
      <div class="job-grid">
        <Card v-for="job in jobs" :key="job.id" class="job-item" hover @click="viewJob(job)" style="cursor: pointer; display: flex; flex-direction: column;">
          <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px;">
            <h3 style="font-size: 18px; margin: 0; color: var(--primary); display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; height: 42px;">{{ job.title }}</h3>
            <Badge v-if="job.department" variant="primary" style="flex-shrink: 0; margin-left: 12px;">{{ job.department }}</Badge>
          </div>
          
          <div style="display: flex; flex-direction: column; gap: 8px; color: var(--text-secondary); font-size: 13px; margin-bottom: 16px; flex: 1;">
            <span style="display: flex; align-items: center; gap: 6px;" v-if="job.company_name">
              <Building2 size="14" style="flex-shrink: 0;" /> <span style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ job.company_name }}</span>
            </span>
            <span style="display: flex; align-items: center; gap: 6px;">
              <MapPin size="14" style="flex-shrink: 0;" /> <span style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ job.location || 'Bất kỳ' }}</span>
            </span>
            <span style="display: flex; align-items: center; gap: 6px;">
              <Clock size="14" style="flex-shrink: 0;" /> <span style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ job.employment_type || 'Full-time' }}</span>
            </span>
          </div>
          
          <div style="border-top: 1px solid var(--border); padding-top: 12px;">
            <p style="color: var(--text-main); font-size: 13px; line-height: 1.5; margin: 0; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;">
              {{ job.description }}
            </p>
          </div>
        </Card>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" style="display: flex; justify-content: center; align-items: center; gap: 16px; margin-top: 32px;">
        <Button variant="outline" :disabled="currentPage === 1" @click="fetchJobs(currentPage - 1)">
          <ChevronLeft size="16" /> Trang trước
        </Button>
        <span style="color: var(--text-secondary); font-size: 14px;">Trang {{ currentPage }} / {{ totalPages }}</span>
        <Button variant="outline" :disabled="currentPage === totalPages" @click="fetchJobs(currentPage + 1)">
          Trang sau <ChevronRight size="16" />
        </Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.job-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}
@media (max-width: 1024px) {
  .job-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
@media (max-width: 768px) {
  .job-grid {
    grid-template-columns: 1fr;
  }
}
.job-item {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  height: 100%;
}
.job-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}
</style>
